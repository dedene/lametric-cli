package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

const defaultTimeout = 10 * time.Second

// Client is the LaMetric local API client.
type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
	ip         string
}

// NewClient creates a Client for a LaMetric device at the given IP.
// It auto-detects the protocol: HTTPS:4343 first, fallback HTTP:8080.
func NewClient(ip, apiKey string) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // LaMetric uses self-signed certs
		},
	}

	c := &Client{
		httpClient: &http.Client{
			Transport: NewRetryTransport(transport),
			Timeout:   defaultTimeout,
		},
		apiKey: apiKey,
		ip:     ip,
	}

	c.baseURL = detectBaseURL(ip)

	return c
}

// detectBaseURL probes HTTPS:4343, falls back to HTTP:8080.
func detectBaseURL(ip string) string {
	httpsURL := fmt.Sprintf("https://%s:4343", ip)
	httpURL := fmt.Sprintf("http://%s:8080", ip)

	// Quick TCP dial to HTTPS port.
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, "4343"), 2*time.Second)
	if err == nil {
		_ = conn.Close()
		return httpsURL
	}

	// Try HTTP port.
	conn, err = net.DialTimeout("tcp", net.JoinHostPort(ip, "8080"), 2*time.Second)
	if err == nil {
		_ = conn.Close()
		return httpURL
	}

	// Default to HTTPS.
	return httpsURL
}

// do executes an API request with basic auth and error handling.
func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.SetBasicAuth("dev", c.apiKey)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return NewAPIError(resp.StatusCode, respBody)
	}

	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, path string, out any) error {
	return c.do(ctx, http.MethodGet, path, nil, out)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, path string, body, out any) error {
	return c.do(ctx, http.MethodPost, path, body, out)
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, path string, body, out any) error {
	return c.do(ctx, http.MethodPut, path, body, out)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// IP returns the device IP address.
func (c *Client) IP() string {
	return c.ip
}

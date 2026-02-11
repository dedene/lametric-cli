package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestClient returns a Client pointed at the given httptest.Server.
func newTestClient(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()
	return &Client{
		httpClient: srv.Client(),
		apiKey:     "test-api-key",
		baseURL:    srv.URL,
	}
}

func TestClient_Get(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/api/v2/device" {
			t.Errorf("path = %s, want /api/v2/device", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"name": "LM1234"})
	}))
	defer srv.Close()

	client := newTestClient(t, srv)
	var out map[string]string
	if err := client.Get(context.Background(), "/api/v2/device", &out); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if out["name"] != "LM1234" {
		t.Errorf("name = %q, want %q", out["name"], "LM1234")
	}
}

func TestClient_Get_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"message":"not found"}]}`))
	}))
	defer srv.Close()

	client := newTestClient(t, srv)
	err := client.Get(context.Background(), "/api/v2/device/notifications/999", nil)
	if err == nil {
		t.Fatal("expected error for 404")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
}

func TestClient_Post(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": "42"})
	}))
	defer srv.Close()

	client := newTestClient(t, srv)
	var out map[string]string
	err := client.Post(context.Background(), "/api/v2/device/notifications", map[string]string{"text": "hello"}, &out)
	if err != nil {
		t.Fatalf("Post: %v", err)
	}
	if out["id"] != "42" {
		t.Errorf("id = %q, want %q", out["id"], "42")
	}
}

func TestClient_Delete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := newTestClient(t, srv)
	if err := client.Delete(context.Background(), "/api/v2/device/notifications/1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestClient_BasicAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			t.Fatal("no basic auth header")
		}
		if user != "dev" {
			t.Errorf("user = %q, want %q", user, "dev")
		}
		if pass != "test-api-key" {
			t.Errorf("pass = %q, want %q", pass, "test-api-key")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := newTestClient(t, srv)
	if err := client.Get(context.Background(), "/api/v2/device", nil); err != nil {
		t.Fatalf("Get: %v", err)
	}
}

func TestClient_Put(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer srv.Close()

	client := newTestClient(t, srv)
	var out map[string]string
	err := client.Put(context.Background(), "/api/v2/device/display", map[string]int{"brightness": 80}, &out)
	if err != nil {
		t.Fatalf("Put: %v", err)
	}
	if out["status"] != "ok" {
		t.Errorf("status = %q, want %q", out["status"], "ok")
	}
}

func TestClient_AuthError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"errors":[{"message":"invalid api key"}]}`))
	}))
	defer srv.Close()

	client := newTestClient(t, srv)
	err := client.Get(context.Background(), "/api/v2/device", nil)
	if err == nil {
		t.Fatal("expected error for 401")
	}

	_, ok := err.(*AuthError)
	if !ok {
		t.Errorf("expected *AuthError, got %T", err)
	}
}

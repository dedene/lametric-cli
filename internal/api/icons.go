package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const cloudIconsURL = "https://developer.lametric.com/api/v2/icons"

// CloudIcon represents an icon from the LaMetric cloud icon library.
type CloudIcon struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Code  string `json:"code"`
	Type  string `json:"type"` // "picture" or "movie"
}

// IconsResponse is the API response for the icons endpoint.
type IconsResponse struct {
	Meta struct {
		TotalIconCount int `json:"total_icon_count"`
		Page           int `json:"page"`
		PageSize       int `json:"page_size"`
		PageCount      int `json:"page_count"`
	} `json:"meta"`
	Data []CloudIcon `json:"data"`
}

// GetPopularIcons returns popular icons from the LaMetric cloud library.
func GetPopularIcons(ctx context.Context, limit int) ([]CloudIcon, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	icons, _, err := fetchIconPage(ctx, client, 0, limit)
	return icons, err
}

// SearchIcons searches the LaMetric cloud icon library.
// It fetches icons and filters by title containing the query (case-insensitive).
func SearchIcons(ctx context.Context, query string, limit int) ([]CloudIcon, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	// Fetch pages until we have enough matches or run out of pages.
	var matches []CloudIcon
	queryLower := strings.ToLower(query)
	pageSize := 200
	maxPages := 10 // Don't fetch more than 2000 icons

	for page := 0; page < maxPages && len(matches) < limit; page++ {
		icons, hasMore, err := fetchIconPage(ctx, client, page, pageSize)
		if err != nil {
			return nil, err
		}

		for _, icon := range icons {
			if strings.Contains(strings.ToLower(icon.Title), queryLower) {
				matches = append(matches, icon)
				if len(matches) >= limit {
					break
				}
			}
		}

		if !hasMore {
			break
		}
	}

	return matches, nil
}

// fetchIconPage fetches a single page of icons from the cloud API.
func fetchIconPage(ctx context.Context, client *http.Client, page, pageSize int) ([]CloudIcon, bool, error) {
	u, _ := url.Parse(cloudIconsURL)
	q := u.Query()
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("page_size", fmt.Sprintf("%d", pageSize))
	q.Set("order", "popular")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, false, fmt.Errorf("create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("fetch icons: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("icons API returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, fmt.Errorf("read response: %w", err)
	}

	var result IconsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, false, fmt.Errorf("parse response: %w", err)
	}

	hasMore := page < result.Meta.PageCount-1
	return result.Data, hasMore, nil
}

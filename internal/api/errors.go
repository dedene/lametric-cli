package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

// APIError represents a LaMetric API error response.
// LaMetric returns: {"errors":[{"message":"..."}]}
type APIError struct {
	StatusCode int
	Messages   []string
}

func (e *APIError) Error() string {
	if len(e.Messages) == 0 {
		return fmt.Sprintf("api error (%d)", e.StatusCode)
	}
	return fmt.Sprintf("api error (%d): %s", e.StatusCode, strings.Join(e.Messages, "; "))
}

// AuthError represents a 401 authentication error.
type AuthError struct {
	APIError
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("authentication failed (%d): %s", e.StatusCode, strings.Join(e.Messages, "; "))
}

// apiErrorResponse mirrors the LaMetric JSON error format.
type apiErrorResponse struct {
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// NewAPIError parses a JSON error response body and returns the appropriate typed error.
func NewAPIError(statusCode int, body []byte) error {
	var msgs []string

	var resp apiErrorResponse
	if err := json.Unmarshal(body, &resp); err == nil {
		for _, e := range resp.Errors {
			if e.Message != "" {
				msgs = append(msgs, e.Message)
			}
		}
	}

	if len(msgs) == 0 {
		msgs = []string{"unknown error"}
	}

	base := APIError{
		StatusCode: statusCode,
		Messages:   msgs,
	}

	switch statusCode {
	case 401, 403:
		return &AuthError{APIError: base}
	default:
		return &base
	}
}

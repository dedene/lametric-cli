package api

import (
	"strings"
	"testing"
)

func TestNewAPIError_SingleMessage(t *testing.T) {
	body := []byte(`{"errors":[{"message":"device not found"}]}`)
	err := NewAPIError(404, body)

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
	if len(apiErr.Messages) != 1 || apiErr.Messages[0] != "device not found" {
		t.Errorf("Messages = %v, want [device not found]", apiErr.Messages)
	}
}

func TestNewAPIError_MultipleMessages(t *testing.T) {
	body := []byte(`{"errors":[{"message":"field required"},{"message":"invalid value"}]}`)
	err := NewAPIError(400, body)

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if len(apiErr.Messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(apiErr.Messages))
	}
	if apiErr.Messages[0] != "field required" || apiErr.Messages[1] != "invalid value" {
		t.Errorf("Messages = %v", apiErr.Messages)
	}
}

func TestNewAPIError_InvalidJSON(t *testing.T) {
	body := []byte(`not json`)
	err := NewAPIError(500, body)

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if len(apiErr.Messages) != 1 || apiErr.Messages[0] != "unknown error" {
		t.Errorf("Messages = %v, want [unknown error]", apiErr.Messages)
	}
}

func TestNewAPIError_AuthError(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{"401", 401},
		{"403", 403},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := []byte(`{"errors":[{"message":"unauthorized"}]}`)
			err := NewAPIError(tt.status, body)

			_, ok := err.(*AuthError)
			if !ok {
				t.Errorf("expected *AuthError for %d, got %T", tt.status, err)
			}
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{StatusCode: 500, Messages: []string{"boom"}}
	got := err.Error()
	if !strings.Contains(got, "500") || !strings.Contains(got, "boom") {
		t.Errorf("Error() = %q, want to contain 500 and boom", got)
	}
}

func TestAPIError_Error_NoMessages(t *testing.T) {
	err := &APIError{StatusCode: 500}
	got := err.Error()
	if !strings.Contains(got, "500") {
		t.Errorf("Error() = %q, want to contain 500", got)
	}
}

func TestAuthError_Error(t *testing.T) {
	err := &AuthError{APIError: APIError{StatusCode: 401, Messages: []string{"bad key"}}}
	got := err.Error()
	if !strings.Contains(got, "authentication failed") || !strings.Contains(got, "bad key") {
		t.Errorf("Error() = %q", got)
	}
}

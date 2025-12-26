package birdeye

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		StatusCode: 404,
		Message:    "Token not found",
		Path:       "/defi/price",
	}

	expected := "birdeye api error: /defi/price returned status 404: Token not found"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

func TestAPIError_IsNotFound(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{http.StatusNotFound, true},
		{http.StatusOK, false},
		{http.StatusInternalServerError, false},
		{http.StatusTooManyRequests, false},
	}

	for _, tt := range tests {
		err := &APIError{StatusCode: tt.statusCode}
		if err.IsNotFound() != tt.expected {
			t.Errorf("IsNotFound() for status %d: expected %v, got %v",
				tt.statusCode, tt.expected, err.IsNotFound())
		}
	}
}

func TestAPIError_IsRateLimited(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{http.StatusTooManyRequests, true},
		{http.StatusOK, false},
		{http.StatusNotFound, false},
		{http.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		err := &APIError{StatusCode: tt.statusCode}
		if err.IsRateLimited() != tt.expected {
			t.Errorf("IsRateLimited() for status %d: expected %v, got %v",
				tt.statusCode, tt.expected, err.IsRateLimited())
		}
	}
}

func TestAPIError_IsServerError(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{500, true},
		{502, true},
		{503, true},
		{504, true},
		{599, true},
		{400, false},
		{404, false},
		{200, false},
	}

	for _, tt := range tests {
		err := &APIError{StatusCode: tt.statusCode}
		if err.IsServerError() != tt.expected {
			t.Errorf("IsServerError() for status %d: expected %v, got %v",
				tt.statusCode, tt.expected, err.IsServerError())
		}
	}
}

func TestAPIError_IsClientError(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{400, true},
		{401, true},
		{403, true},
		{404, true},
		{499, true},
		{500, false},
		{200, false},
	}

	for _, tt := range tests {
		err := &APIError{StatusCode: tt.statusCode}
		if err.IsClientError() != tt.expected {
			t.Errorf("IsClientError() for status %d: expected %v, got %v",
				tt.statusCode, tt.expected, err.IsClientError())
		}
	}
}

func TestIsAPIError(t *testing.T) {
	t.Run("direct APIError", func(t *testing.T) {
		err := &APIError{StatusCode: 404, Message: "not found"}
		apiErr, ok := IsAPIError(err)
		if !ok {
			t.Error("expected IsAPIError to return true")
		}
		if apiErr.StatusCode != 404 {
			t.Errorf("expected status 404, got %d", apiErr.StatusCode)
		}
	})

	t.Run("wrapped APIError", func(t *testing.T) {
		original := &APIError{StatusCode: 429, Message: "rate limited"}
		wrapped := fmt.Errorf("failed to get price: %w", original)

		apiErr, ok := IsAPIError(wrapped)
		if !ok {
			t.Error("expected IsAPIError to return true for wrapped error")
		}
		if apiErr.StatusCode != 429 {
			t.Errorf("expected status 429, got %d", apiErr.StatusCode)
		}
	})

	t.Run("non-APIError", func(t *testing.T) {
		err := errors.New("some other error")
		_, ok := IsAPIError(err)
		if ok {
			t.Error("expected IsAPIError to return false for non-APIError")
		}
	})

	t.Run("nil error", func(t *testing.T) {
		_, ok := IsAPIError(nil)
		if ok {
			t.Error("expected IsAPIError to return false for nil")
		}
	})
}

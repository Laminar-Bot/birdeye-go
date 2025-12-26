package birdeye

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError represents an error response from the Birdeye API.
type APIError struct {
	// StatusCode is the HTTP status code returned.
	StatusCode int

	// Message is the error message from the API response body.
	Message string

	// Path is the API endpoint that returned the error.
	Path string
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("birdeye api error: %s returned status %d: %s",
		e.Path, e.StatusCode, e.Message)
}

// IsNotFound returns true if the error indicates the resource was not found.
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsRateLimited returns true if the error indicates rate limiting.
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// IsServerError returns true if the error is a server-side error (5xx).
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500
}

// IsClientError returns true if the error is a client-side error (4xx).
func (e *APIError) IsClientError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

// IsAPIError checks if an error is a Birdeye API error and returns it.
// This correctly handles wrapped errors using errors.As.
func IsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}

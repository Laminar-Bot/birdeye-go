package birdeye

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testServer creates a test HTTP server that responds with the provided responses.
// The responses map uses path as key and JSON response body as value.
func testServer(t *testing.T, responses map[string]interface{}) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check API key header
		if r.Header.Get("X-API-KEY") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "missing api key"}`))
			return
		}

		// Look up response for path
		response, ok := responses[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error": "not found"}`))
			return
		}

		// Handle special response types
		switch resp := response.(type) {
		case int:
			// Integer means return that status code with empty body
			w.WriteHeader(resp)
		case error:
			// Error means return 500 with error message
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(resp.Error()))
		default:
			// Otherwise, marshal as JSON
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}
	}))
}

// testClient creates a test client pointing to the test server.
func testClient(t *testing.T, serverURL string) *Client {
	t.Helper()

	client, err := NewClient("test-api-key",
		WithBaseURL(serverURL),
		WithMaxRetries(0), // Disable retries for tests
	)
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}
	return client
}

// wrapResponse wraps data in Birdeye's standard response format.
func wrapResponse(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": true,
		"data":    data,
	}
}

// wrapFailure creates a failed Birdeye response.
func wrapFailure() map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"data":    nil,
	}
}

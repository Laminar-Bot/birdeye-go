package birdeye

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient_RequiresAPIKey(t *testing.T) {
	_, err := NewClient("")
	if err == nil {
		t.Error("expected error for empty API key")
	}
}

func TestNewClient_DefaultConfig(t *testing.T) {
	client, err := NewClient("test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client.apiKey != "test-key" {
		t.Errorf("expected apiKey 'test-key', got '%s'", client.apiKey)
	}
	if client.baseURL != DefaultBaseURL {
		t.Errorf("expected baseURL '%s', got '%s'", DefaultBaseURL, client.baseURL)
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	customURL := "https://custom.api.example.com"
	client, err := NewClient("test-key",
		WithBaseURL(customURL),
		WithTimeout(30*time.Second),
		WithMaxRetries(5),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client.baseURL != customURL {
		t.Errorf("expected baseURL '%s', got '%s'", customURL, client.baseURL)
	}
}

func TestNewClient_WithLogger(t *testing.T) {
	logger := &testLogger{}
	client, err := NewClient("test-key", WithLogger(logger))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Make a request to trigger logging
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success": true, "data": {"value": 1.5}}`))
	}))
	defer server.Close()

	client.baseURL = server.URL
	_, _ = client.GetPrice(context.Background(), "test-address")

	if !logger.debugCalled {
		t.Error("expected Debug to be called")
	}
}

func TestNewClient_WithHTTPClient(t *testing.T) {
	customClient := &http.Client{Timeout: 60 * time.Second}
	client, err := NewClient("test-key", WithHTTPClient(customClient))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client.httpClient != customClient {
		t.Error("expected custom HTTP client to be used")
	}
}

func TestClient_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient("test-key",
		WithBaseURL(server.URL),
		WithMaxRetries(0),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.GetPrice(ctx, "test-address")
	if err == nil {
		t.Error("expected error due to cancelled context")
	}
}

// testLogger implements Logger for testing.
type testLogger struct {
	debugCalled bool
	infoCalled  bool
	warnCalled  bool
	errorCalled bool
}

func (l *testLogger) Debug(_ string, _ ...interface{}) { l.debugCalled = true }
func (l *testLogger) Info(_ string, _ ...interface{})  { l.infoCalled = true }
func (l *testLogger) Warn(_ string, _ ...interface{})  { l.warnCalled = true }
func (l *testLogger) Error(_ string, _ ...interface{}) { l.errorCalled = true }

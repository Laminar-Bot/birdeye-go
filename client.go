package birdeye

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// API configuration defaults.
const (
	// DefaultBaseURL is the Birdeye public API endpoint.
	DefaultBaseURL = "https://public-api.birdeye.so"

	// DefaultTimeout for HTTP requests.
	DefaultTimeout = 10 * time.Second

	// DefaultMaxRetries before giving up on a request.
	DefaultMaxRetries = 3

	// DefaultRetryWaitMin is the minimum wait time between retries.
	DefaultRetryWaitMin = 500 * time.Millisecond

	// DefaultRetryWaitMax is the maximum wait time between retries.
	DefaultRetryWaitMax = 3 * time.Second

	// chainSolana is the Solana chain identifier for Birdeye API.
	chainSolana = "solana"
)

// Logger is an optional interface for structured logging.
// Implement this to integrate with your logging library.
type Logger interface {
	// Debug logs a debug message with optional key-value pairs.
	Debug(msg string, keysAndValues ...interface{})

	// Info logs an info message with optional key-value pairs.
	Info(msg string, keysAndValues ...interface{})

	// Warn logs a warning message with optional key-value pairs.
	Warn(msg string, keysAndValues ...interface{})

	// Error logs an error message with optional key-value pairs.
	Error(msg string, keysAndValues ...interface{})
}

// noopLogger is the default logger that discards all log messages.
type noopLogger struct{}

func (noopLogger) Debug(_ string, _ ...interface{}) {}
func (noopLogger) Info(_ string, _ ...interface{})  {}
func (noopLogger) Warn(_ string, _ ...interface{})  {}
func (noopLogger) Error(_ string, _ ...interface{}) {}

// Client provides methods for interacting with the Birdeye API.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	logger     Logger
}

// config holds internal configuration built from options.
type config struct {
	baseURL      string
	timeout      time.Duration
	maxRetries   int
	retryWaitMin time.Duration
	retryWaitMax time.Duration
	logger       Logger
	httpClient   *http.Client
}

// Option configures the Client.
type Option func(*config)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(url string) Option {
	return func(c *config) {
		c.baseURL = url
	}
}

// WithTimeout sets the HTTP request timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *config) {
		c.timeout = d
	}
}

// WithMaxRetries sets the maximum number of retries for failed requests.
func WithMaxRetries(n int) Option {
	return func(c *config) {
		c.maxRetries = n
	}
}

// WithRetryWait sets the minimum and maximum wait times between retries.
func WithRetryWait(min, max time.Duration) Option {
	return func(c *config) {
		c.retryWaitMin = min
		c.retryWaitMax = max
	}
}

// WithLogger sets a custom logger for the client.
// If not set, logging is disabled (noop logger is used).
func WithLogger(l Logger) Option {
	return func(c *config) {
		c.logger = l
	}
}

// WithHTTPClient sets a custom HTTP client.
// This overrides the default retryable client. Use with caution.
func WithHTTPClient(client *http.Client) Option {
	return func(c *config) {
		c.httpClient = client
	}
}

// NewClient creates a new Birdeye API client.
//
// The apiKey is required. Additional options can be provided to customize
// the client behavior.
//
// Example:
//
//	client, err := birdeye.NewClient("your-api-key",
//	    birdeye.WithTimeout(30 * time.Second),
//	    birdeye.WithMaxRetries(5),
//	)
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("api key is required")
	}

	// Apply defaults.
	cfg := &config{
		baseURL:      DefaultBaseURL,
		timeout:      DefaultTimeout,
		maxRetries:   DefaultMaxRetries,
		retryWaitMin: DefaultRetryWaitMin,
		retryWaitMax: DefaultRetryWaitMax,
		logger:       noopLogger{},
	}

	// Apply options.
	for _, opt := range opts {
		opt(cfg)
	}

	// Use custom HTTP client if provided.
	var httpClient *http.Client
	if cfg.httpClient != nil {
		httpClient = cfg.httpClient
	} else {
		// Configure retryable HTTP client with exponential backoff.
		retryClient := retryablehttp.NewClient()
		retryClient.RetryMax = cfg.maxRetries
		retryClient.RetryWaitMin = cfg.retryWaitMin
		retryClient.RetryWaitMax = cfg.retryWaitMax
		retryClient.HTTPClient.Timeout = cfg.timeout

		// Disable retryablehttp's default logging.
		retryClient.Logger = nil

		// Custom retry policy: retry on 429 (rate limit) and 5xx errors.
		retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
			// Don't retry on context cancellation.
			if ctx.Err() != nil {
				return false, ctx.Err()
			}

			// Retry on connection errors.
			if err != nil {
				return true, err
			}

			// Retry on rate limiting.
			if resp.StatusCode == http.StatusTooManyRequests {
				return true, nil
			}

			// Retry on server errors.
			if resp.StatusCode >= 500 {
				return true, nil
			}

			return false, nil
		}

		httpClient = retryClient.StandardClient()
	}

	return &Client{
		apiKey:     apiKey,
		baseURL:    cfg.baseURL,
		httpClient: httpClient,
		logger:     cfg.logger,
	}, nil
}

// doGet performs a GET request to the Birdeye API.
func (c *Client) doGet(ctx context.Context, path string, params url.Values) ([]byte, error) {
	// Build request URL.
	reqURL := c.baseURL + path
	if len(params) > 0 {
		reqURL = reqURL + "?" + params.Encode()
	}

	// Create request with context for cancellation support.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set required headers.
	req.Header.Set("X-API-KEY", c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-chain", chainSolana)

	c.logger.Debug("birdeye api request", "method", http.MethodGet, "path", path)

	// Execute request.
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("birdeye api request failed", "path", path, "error", err)
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn("failed to close response body", "error", closeErr)
		}
	}()

	// Read response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	// Handle non-OK status codes.
	if resp.StatusCode != http.StatusOK {
		c.logger.Error("birdeye api error response",
			"path", path,
			"status_code", resp.StatusCode,
			"body", truncateForLog(string(body), 500),
		)

		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
			Path:       path,
		}
	}

	return body, nil
}

// parseResponse unmarshals a Birdeye API response and checks the success flag.
//
// Birdeye responses follow this structure:
//
//	{
//	  "success": true,
//	  "data": { ... }
//	}
func parseResponse[T any](body []byte) (*T, error) {
	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
		Data    T      `json:"data"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if !resp.Success {
		if resp.Message != "" {
			return nil, fmt.Errorf("birdeye api error: %s", resp.Message)
		}
		return nil, errors.New("birdeye api returned success=false")
	}

	return &resp.Data, nil
}

// truncateForLog truncates a string for safe logging.
func truncateForLog(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "...(truncated)"
}

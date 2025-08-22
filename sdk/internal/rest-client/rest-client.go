package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alexandernovytsky/iac-assignment/sdk/config"
	"github.com/alexandernovytsky/iac-assignment/sdk/errors"
)

// RestClient handles HTTP communication with the API
type RestClient struct {
	base       string
	apiKey     string
	httpClient *http.Client
	maxRetries int
	backoff    time.Duration
	headers    map[string]string
}

// RestClientOption defines function type for options
type RestClientOption func(*RestClient)

// WithMaxRetries sets the maximum number of retries
func WithMaxRetries(maxRetries int) RestClientOption {
	return func(rc *RestClient) {
		if maxRetries >= 0 {
			rc.maxRetries = maxRetries
		}
	}
}

// WithBackoff sets the backoff duration between retries
func WithBackoff(backoff time.Duration) RestClientOption {
	return func(rc *RestClient) {
		if backoff > 0 {
			rc.backoff = backoff
		}
	}
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) RestClientOption {
	return func(rc *RestClient) {
		if timeout > 0 {
			rc.httpClient.Timeout = timeout
		}
	}
}

// WithHeader adds a custom header to all requests
func WithHeader(key, value string) RestClientOption {
	return func(rc *RestClient) {
		rc.headers[key] = value
	}
}

// NewRestClient creates a new REST client with the given configuration
func NewRestClient(baseURL, apiKey string, options ...RestClientOption) *RestClient {
	client := &RestClient{
		base:       baseURL,
		apiKey:     apiKey,
		maxRetries: config.Defaults.MaxRetries,
		backoff:    config.Defaults.Backoff,
		httpClient: &http.Client{
			Timeout: config.Defaults.Timeout,
		},
		headers: make(map[string]string),
	}

	// Apply options
	for _, option := range options {
		option(client)
	}

	return client
}

// Request executes an HTTP request with the given method
func (rc *RestClient) Request(ctx context.Context, path, method string, in, out any) error {
	if ctx == nil {
		return errors.NewInputError(fmt.Errorf("context cannot be nil"))
	}

	var lastErr error
	for attempt := 0; attempt <= rc.maxRetries; attempt++ {
		if attempt > 0 {
			// Apply exponential backoff
			sleep := rc.backoff * time.Duration(1<<uint(attempt-1))

			// Create a timer for sleeping
			timer := time.NewTimer(sleep)

			// Wait for either the timer or context cancellation
			select {
			case <-ctx.Done():
				timer.Stop()
				return errors.NewSDKError(0, "", fmt.Errorf("context canceled during retry: %w", ctx.Err()))
			case <-timer.C:
				// Continue with retry
			}
		}

		lastErr = rc.doRequest(ctx, path, method, in, out)
		if lastErr == nil {
			return nil // Success
		}

		// Only retry for server errors (5xx) or 429 Too Many Requests
		if sdkErr, ok := lastErr.(*errors.SDKError); ok {
			if !(sdkErr.StatusCode == 429 || (sdkErr.StatusCode >= 500 && sdkErr.StatusCode < 600)) {
				return lastErr // Do not retry for other errors
			}
		} else {
			return lastErr // Not an SDKError, do not retry
		}

		// Check if context was canceled
		select {
		case <-ctx.Done():
			return errors.NewSDKError(0, "", fmt.Errorf("context canceled: %w", ctx.Err()))
		default:
			// Continue with retry
		}
	}

	return fmt.Errorf("failed after %d retries: %w", rc.maxRetries, lastErr)
}

// doRequest performs a single HTTP request
func (rc *RestClient) doRequest(ctx context.Context, path, method string, in, out any) error {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return errors.NewInputError(fmt.Errorf("marshal payload: %w", err))
		}
		body = bytes.NewBuffer(b)
	}

	url := rc.base + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return errors.NewInputError(fmt.Errorf("create request: %w", err))
	}

	// Set default headers
	req.Header.Set(config.HeaderKeys.ContentType, config.HeaderValues.JSONContentType)
	req.Header.Set(config.HeaderKeys.Authorization, fmt.Sprintf(config.HeaderValues.AuthFormat, rc.apiKey))

	// Set custom headers
	for key, value := range rc.headers {
		req.Header.Set(key, value)
	}

	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return errors.NewSDKError(0, "", fmt.Errorf("do request: %w", err))
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.NewSDKError(0, "", fmt.Errorf("read response body: %w", err))
	}

	// Check if response status code indicates success (2xx)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.NewSDKError(
			resp.StatusCode,
			string(bodyBytes),
			fmt.Errorf("unexpected status code: %d", resp.StatusCode),
		)
	}

	// If out is nil, the caller doesn't expect a response body
	if out == nil {
		return nil
	}

	if err := json.Unmarshal(bodyBytes, out); err != nil {
		return errors.NewSDKError(
			resp.StatusCode,
			string(bodyBytes),
			fmt.Errorf("unmarshal response: %w", err),
		)
	}

	return nil
}

// Convenience methods for HTTP verbs
func (rc *RestClient) Get(ctx context.Context, path string, out any) error {
	return rc.Request(ctx, path, config.Method.GET, nil, out)
}

func (rc *RestClient) Post(ctx context.Context, path string, in, out any) error {
	return rc.Request(ctx, path, config.Method.POST, in, out)
}

func (rc *RestClient) Put(ctx context.Context, path string, in, out any) error {
	return rc.Request(ctx, path, config.Method.PUT, in, out)
}

func (rc *RestClient) Delete(ctx context.Context, path string, out any) error {
	return rc.Request(ctx, path, config.Method.DELETE, nil, out)
}

func (rc *RestClient) Patch(ctx context.Context, path string, in, out any) error {
	return rc.Request(ctx, path, config.Method.PATCH, in, out)
}

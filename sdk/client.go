package sdk

import (
	"time"

	alerts "github.com/alexandernovytsky/iac-assignment/sdk/alerts/v3"
	"github.com/alexandernovytsky/iac-assignment/sdk/config"
	internal "github.com/alexandernovytsky/iac-assignment/sdk/internal/rest-client"
	webhooks "github.com/alexandernovytsky/iac-assignment/sdk/webhooks/v1"
)

// ClientCreator helps create resource-specific clients
type ClientCreator struct {
	region     string
	apiKey     string
	baseURL    string
	timeout    time.Duration
	maxRetries int
	backoff    time.Duration
	headers    map[string]string
	restClient *internal.RestClient
}

// ClientOption defines function type for client options
type ClientOption func(*ClientCreator)

// WithTimeout sets HTTP client timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(cc *ClientCreator) {
		if timeout > 0 {
			cc.timeout = timeout
		}
	}
}

// WithMaxRetries sets maximum number of retries
func WithMaxRetries(maxRetries int) ClientOption {
	return func(cc *ClientCreator) {
		if maxRetries >= 0 {
			cc.maxRetries = maxRetries
		}
	}
}

// WithBackoff sets backoff duration between retries
func WithBackoff(backoff time.Duration) ClientOption {
	return func(cc *ClientCreator) {
		if backoff > 0 {
			cc.backoff = backoff
		}
	}
}

// WithHeader adds a custom header to all requests
func WithHeader(key, value string) ClientOption {
	return func(cc *ClientCreator) {
		cc.headers[key] = value
	}
}

// NewClientCreator initializes a client creator with common configuration
func NewClientCreator(region, apiKey string, options ...ClientOption) *ClientCreator {
	cc := &ClientCreator{
		region:     region,
		apiKey:     apiKey,
		baseURL:    config.GetBaseURL(region),
		timeout:    config.Defaults.Timeout,
		maxRetries: config.Defaults.MaxRetries,
		backoff:    config.Defaults.Backoff,
		headers:    make(map[string]string),
	}

	// Apply options
	for _, option := range options {
		option(cc)
	}

	return cc
}

// createRestClient creates a REST client for a specific resource
func (cc *ClientCreator) createRestClient() *internal.RestClient {
	if cc.restClient != nil {
		return cc.restClient
	}

	// Create a slice of internal options based on our configuration
	var options []internal.RestClientOption

	// Only add options if they differ from defaults
	if cc.timeout != config.Defaults.Timeout {
		options = append(options, internal.WithTimeout(cc.timeout))
	}

	if cc.maxRetries != config.Defaults.MaxRetries {
		options = append(options, internal.WithMaxRetries(cc.maxRetries))
	}

	if cc.backoff != config.Defaults.Backoff {
		options = append(options, internal.WithBackoff(cc.backoff))
	}

	// Add all headers
	for k, v := range cc.headers {
		options = append(options, internal.WithHeader(k, v))
	}

	cc.restClient = internal.NewRestClient(cc.baseURL, cc.apiKey, options...)
	return cc.restClient
}

// Alerts returns an alerts client
func (cc *ClientCreator) Alerts() *alerts.AlertsClient {
	return alerts.NewAlertsV3Client(cc.createRestClient())
}

// Webhooks returns a webhooks client
func (cc *ClientCreator) Webhooks() *webhooks.WebhooksClient {
	return webhooks.NewWebhooksV1Client(cc.createRestClient())
}

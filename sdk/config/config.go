package config

import (
	"fmt"
	"time"
)

// API endpoints
const (
	BaseURLFormat = "https://api.%s.coralogix.com/mgmt/openapi"
)

const (
	RegionEU2 = "eu2"
)

// API versions
const (
	V1 = "v1"
	V2 = "v2"
	V3 = "v3"
)

// Resource paths
const (
	AlertsPath   = "alert-defs"
	WebhooksPath = "outgoing-webhooks"
	LogsPath     = "logs"
)

// HeaderKeys contains common HTTP header keys
var HeaderKeys = struct {
	ContentType   string
	Authorization string
	Accept        string
	CorrelationID string
}{
	ContentType:   "Content-Type",
	Authorization: "Authorization",
	Accept:        "Accept",
	CorrelationID: "X-Correlation-ID",
}

// HeaderValues contains common HTTP header values
var HeaderValues = struct {
	JSONContentType string
	AuthFormat      string
}{
	JSONContentType: "application/json",
	AuthFormat:      "Bearer %s",
}

// HTTP methods
var Method = struct {
	GET    string
	POST   string
	PUT    string
	DELETE string
	PATCH  string
}{
	GET:    "GET",
	POST:   "POST",
	PUT:    "PUT",
	DELETE: "DELETE",
	PATCH:  "PATCH",
}

// Defaults for client options
var Defaults = struct {
	MaxRetries int
	Backoff    time.Duration
	Timeout    time.Duration
}{
	MaxRetries: 3,
	Backoff:    100 * time.Millisecond,
	Timeout:    5 * time.Second,
}

// GetBaseURL returns the base URL for a specific region
func GetBaseURL(region string) string {
	return fmt.Sprintf(BaseURLFormat, region)
}

// GetResourcePath builds a path with version and resource
func GetResourcePath(version, resource string) string {
	return fmt.Sprintf("/%s/%s", version, resource)
}

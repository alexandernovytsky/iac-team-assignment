package errors

import (
	"fmt"
	"net/http"
)

// Common error types
var (
	ErrBadRequest              = "bad_request"
	ErrUnauthorized            = "unauthorized"
	ErrForbidden               = "forbidden"
	ErrNotFound                = "not_found"
	ErrRateLimit               = "rate_limit"
	ErrServerError             = "server_error"
	ErrNetworkError            = "network_error"
	ErrInvalidInput            = "invalid_input"
	ErrContextCanceled         = "context_canceled"
	ErrContextDeadlineExceeded = "context_deadline_exceeded"
)

// SDKError represents an error returned by the API
type SDKError struct {
	StatusCode int
	Body       string
	Type       string
	Message    string
	Err        error
}

func (e *SDKError) Error() string {
	return fmt.Sprintf("coralogix SDK error: type=%s, status=%d, message=%s, err=%v",
		e.Type, e.StatusCode, e.Message, e.Err)
}

func (e *SDKError) Unwrap() error {
	return e.Err
}

// NewSDKError creates a new SDKError from an HTTP response
func NewSDKError(statusCode int, body string, err error) *SDKError {
	errorType := getErrorTypeFromStatus(statusCode)

	return &SDKError{
		StatusCode: statusCode,
		Body:       body,
		Type:       errorType,
		Message:    fmt.Sprintf("API request failed with status %d", statusCode),
		Err:        err,
	}
}

// NewInputError creates a new SDKError for input validation errors
func NewInputError(err error) *SDKError {
	return &SDKError{
		StatusCode: 0,
		Type:       ErrInvalidInput,
		Message:    "Invalid input parameters",
		Err:        err,
	}
}

// getErrorTypeFromStatus maps HTTP status codes to error types
func getErrorTypeFromStatus(statusCode int) string {
	switch {
	case statusCode >= 400 && statusCode < 500:
		switch statusCode {
		case http.StatusBadRequest:
			return ErrBadRequest
		case http.StatusUnauthorized:
			return ErrUnauthorized
		case http.StatusForbidden:
			return ErrForbidden
		case http.StatusNotFound:
			return ErrNotFound
		case http.StatusTooManyRequests:
			return ErrRateLimit
		default:
			return ErrBadRequest
		}
	case statusCode >= 500:
		return ErrServerError
	default:
		return ErrNetworkError
	}
}

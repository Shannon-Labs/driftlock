package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError represents a structured API error response
type APIError struct {
	Code       int               `json:"code"`
	HTTPStatus int               `json:"-"`  // For compatibility with tests
	Message    string            `json:"message"`
	Details    string            `json:"details,omitempty"`
	FieldErrors map[string]string `json:"field_errors,omitempty"`
	RequestID  string            `json:"request_id,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// SetHTTPStatus sets the HTTP status code
func (e *APIError) SetHTTPStatus(statusCode int) *APIError {
	e.Code = statusCode
	e.HTTPStatus = statusCode
	return e
}

// GetHTTPStatus returns the HTTP status code
func (e *APIError) GetHTTPStatus() int {
	return e.HTTPStatus
}

// WithDetails adds details to an API error
func (e *APIError) WithDetails(details string) *APIError {
	e.Details = details
	return e
}

// WithFieldError adds a field-specific error
func (e *APIError) WithFieldError(field, message string) *APIError {
	if e.FieldErrors == nil {
		e.FieldErrors = make(map[string]string)
	}
	e.FieldErrors[field] = message
	return e
}

// WithRequestID adds a request ID for tracking
func (e *APIError) WithRequestID(requestID string) *APIError {
	e.RequestID = requestID
	return e
}

// Common API errors
var (
	// 4xx Client Errors
	ErrBadRequest = &APIError{
		Code:       http.StatusBadRequest,
		HTTPStatus: http.StatusBadRequest,
		Message:    "Bad request",
	}

	ErrInvalidJSON = &APIError{
		Code:       http.StatusBadRequest,
		HTTPStatus: http.StatusBadRequest,
		Message:    "Invalid JSON payload",
	}

	ErrValidationFailed = &APIError{
		Code:       http.StatusBadRequest,
		HTTPStatus: http.StatusBadRequest,
		Message:    "Validation failed",
	}

	ErrUnauthorized = &APIError{
		Code:       http.StatusUnauthorized,
		HTTPStatus: http.StatusUnauthorized,
		Message:    "Unauthorized",
		Details:    "Valid API key or authentication token required",
	}

	ErrForbidden = &APIError{
		Code:       http.StatusForbidden,
		HTTPStatus: http.StatusForbidden,
		Message:    "Forbidden",
		Details:    "You don't have permission to access this resource",
	}

	ErrNotFound = &APIError{
		Code:       http.StatusNotFound,
		HTTPStatus: http.StatusNotFound,
		Message:    "Resource not found",
	}

	ErrMethodNotAllowed = &APIError{
		Code:       http.StatusMethodNotAllowed,
		HTTPStatus: http.StatusMethodNotAllowed,
		Message:    "Method not allowed",
	}

	ErrConflict = &APIError{
		Code:       http.StatusConflict,
		HTTPStatus: http.StatusConflict,
		Message:    "Resource conflict",
	}

	ErrPayloadTooLarge = &APIError{
		Code:       http.StatusRequestEntityTooLarge,
		HTTPStatus: http.StatusRequestEntityTooLarge,
		Message:    "Payload too large",
		Details:    "Request body exceeds maximum allowed size",
	}

	ErrTooManyRequests = &APIError{
		Code:       http.StatusTooManyRequests,
		HTTPStatus: http.StatusTooManyRequests,
		Message:    "Rate limit exceeded",
		Details:    "Too many requests, please try again later",
	}

	// 5xx Server Errors
	ErrInternalServer = &APIError{
		Code:       http.StatusInternalServerError,
		HTTPStatus: http.StatusInternalServerError,
		Message:    "Internal server error",
		Details:    "An unexpected error occurred",
	}

	ErrNotImplemented = &APIError{
		Code:       http.StatusNotImplemented,
		HTTPStatus: http.StatusNotImplemented,
		Message:    "Not implemented",
	}

	ErrServiceUnavailable = &APIError{
		Code:       http.StatusServiceUnavailable,
		HTTPStatus: http.StatusServiceUnavailable,
		Message:    "Service unavailable",
		Details:    "The service is temporarily unavailable",
	}

	ErrGatewayTimeout = &APIError{
		Code:       http.StatusGatewayTimeout,
		HTTPStatus: http.StatusGatewayTimeout,
		Message:    "Gateway timeout",
		Details:    "Upstream service did not respond in time",
	}
)

// New creates a new API error with a custom message
func New(code int, message string) *APIError {
	return &APIError{
		Code:       code,
		HTTPStatus: code,
		Message:    message,
	}
}

// Wrap wraps a standard error into an APIError
func Wrap(err error, code int, message string) *APIError {
	return &APIError{
		Code:       code,
		HTTPStatus: code,
		Message:    message,
		Details:    err.Error(),
	}
}

// WriteJSON writes an APIError as a JSON response
func WriteJSON(w http.ResponseWriter, err *APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(err)
}

// HTTPError is a convenience function to write an error response
func HTTPError(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(*APIError); ok {
		WriteJSON(w, apiErr)
		return
	}

	// Default to internal server error for unknown errors
	WriteJSON(w, ErrInternalServer.WithDetails(err.Error()))
}

// ValidationError creates a validation error with field-specific messages
func ValidationError(fieldErrors map[string]string) *APIError {
	return &APIError{
		Code:        http.StatusBadRequest,
		Message:     "Validation failed",
		FieldErrors: fieldErrors,
	}
}

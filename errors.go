package controld

import (
	"net/http"
)

const (
	errEmptyAPIToken        = "invalid credentials: API Token must not be empty" //nolint:gosec,unused
	errInternalServiceError = "internal service error"
	errMakeRequestError     = "error from makeRequest"
	errMarshalError         = "error marshalling the object"
	errUnmarshalError       = "error unmarshalling the JSON response"
	errTypeError            = "error verifying the type of the response"
	errUnmarshalErrorBody   = "error unmarshalling the JSON response error body"
)

type ErrorType string

const (
	ErrorTypeRequest        ErrorType = "request"
	ErrorTypeAuthentication ErrorType = "authentication"
	ErrorTypeAuthorization  ErrorType = "authorization"
	ErrorTypeNotFound       ErrorType = "not_found"
	ErrorTypeRateLimit      ErrorType = "rate_limit"
)

type Error struct {
	// The classification of error encountered.
	Type ErrorType

	// StatusCode is the HTTP status code from the response.
	StatusCode int

	// Errors is all of the error messages and codes, combined.
	Error ResponseInfo
}

// RequestError is for 4xx errors that we encounter not covered elsewhere
// (generally bad payloads).
type RequestError struct {
	controldError *Error
}

func (e RequestError) Error() string {
	return e.controldError.Error.Message
}

func (e RequestError) InternalErrorCodeIs(code int) bool {
	return e.controldError.InternalErrorCodeIs(code)
}

func (e RequestError) Type() ErrorType {
	return e.controldError.Type
}

func NewRequestError(e *Error) RequestError {
	return RequestError{
		controldError: e,
	}
}

// RatelimitError is for HTTP 429s where the service is telling the client to
// slow down.
type RatelimitError struct {
	controldError *Error
}

func (e RatelimitError) Error() string {
	return e.controldError.Error.Message
}

func (e RatelimitError) InternalErrorCodeIs(code int) bool {
	return e.controldError.InternalErrorCodeIs(code)
}

func (e RatelimitError) Type() ErrorType {
	return e.controldError.Type
}

func NewRatelimitError(e *Error) RatelimitError {
	return RatelimitError{
		controldError: e,
	}
}

// ServiceError is a handler for 5xx errors returned to the client.
type ServiceError struct {
	controldError *Error
}

func (e ServiceError) Error() string {
	return e.controldError.Error.Message
}

func (e ServiceError) InternalErrorCodeIs(code int) bool {
	return e.controldError.InternalErrorCodeIs(code)
}

func (e ServiceError) Type() ErrorType {
	return e.controldError.Type
}

func NewServiceError(e *Error) ServiceError {
	return ServiceError{
		controldError: e,
	}
}

// AuthenticationError is for HTTP 401 responses.
type AuthenticationError struct {
	controldError *Error
}

func (e AuthenticationError) Error() string {
	return e.controldError.Error.Message
}

func (e AuthenticationError) InternalErrorCodeIs(code int) bool {
	return e.controldError.InternalErrorCodeIs(code)
}

func (e AuthenticationError) Type() ErrorType {
	return e.controldError.Type
}

func NewAuthenticationError(e *Error) AuthenticationError {
	return AuthenticationError{
		controldError: e,
	}
}

// AuthorizationError is for HTTP 403 responses.
type AuthorizationError struct {
	controldError *Error
}

func (e AuthorizationError) Error() string {
	return e.controldError.Error.Message
}

func (e AuthorizationError) InternalErrorCodeIs(code int) bool {
	return e.controldError.InternalErrorCodeIs(code)
}

func (e AuthorizationError) Type() ErrorType {
	return e.controldError.Type
}

func NewAuthorizationError(e *Error) AuthorizationError {
	return AuthorizationError{
		controldError: e,
	}
}

// NotFoundError is for HTTP 404 responses.
type NotFoundError struct {
	controldError *Error
}

func (e NotFoundError) Error() string {
	return e.controldError.Error.Message
}

func (e NotFoundError) InternalErrorCodeIs(code int) bool {
	return e.controldError.InternalErrorCodeIs(code)
}

func (e NotFoundError) Type() ErrorType {
	return e.controldError.Type
}

func NewNotFoundError(e *Error) NotFoundError {
	return NotFoundError{
		controldError: e,
	}
}

// ClientError returns a boolean whether or not the raised error was caused by
// something client side.
func (e *Error) ClientError() bool {
	return e.StatusCode >= http.StatusBadRequest &&
		e.StatusCode < http.StatusInternalServerError
}

// ClientRateLimited returns a boolean whether or not the raised error was
// caused by too many requests from the client.
func (e *Error) ClientRateLimited() bool {
	return e.Type == ErrorTypeRateLimit
}

// InternalErrorCodeIs returns a boolean whether or not the desired internal
// error code is present in `e.InternalErrorCodes`.
func (e *Error) InternalErrorCodeIs(code int) bool {
	return e.StatusCode == code
}

package icheck

import "encoding/json"

// ErrorType is the list of allowed values for the error's type.
type ErrorType string

// ErrorCode is the list of allowed values for the error's code.
type ErrorCode string

const (
	// ErrorTypeAPI ...
	ErrorTypeAPI ErrorType = "api_error"
	// ErrorTypeAPIConnection ...
	ErrorTypeAPIConnection ErrorType = "api_connection_error"
	// ErrorTypeAuthentication ...
	ErrorTypeAuthentication ErrorType = "authentication_error"
	// ErrorTypeCard ...
	ErrorTypeCard ErrorType = "card_error"
	// ErrorTypeInvalidRequest ...
	ErrorTypeInvalidRequest ErrorType = "invalid_request_error"
	// ErrorTypePermission ...
	ErrorTypePermission ErrorType = "more_permissions_required"
	// ErrorTypeRateLimit ...
	ErrorTypeRateLimit ErrorType = "rate_limit_error"

	// IncorrectNum ...
	IncorrectNum ErrorCode = "incorrect_number"
	// InvalidNum ...
	InvalidNum ErrorCode = "invalid_number"
	// InvalidExpM ...
	InvalidExpM ErrorCode = "invalid_expiry_month"
	// InvalidExpY ...
	InvalidExpY ErrorCode = "invalid_expiry_year"
	// InvalidCvc ...
	InvalidCvc ErrorCode = "invalid_cvc"
	// ExpiredCard ...
	ExpiredCard ErrorCode = "expired_card"
	// IncorrectCvc ...
	IncorrectCvc ErrorCode = "incorrect_cvc"
	// IncorrectZip ...
	IncorrectZip ErrorCode = "incorrect_zip"
	// CardDeclined ...
	CardDeclined ErrorCode = "card_declined"
	// Missing ...
	Missing ErrorCode = "missing"
	// ProcessingErr ...
	ProcessingErr ErrorCode = "processing_error"
	// RateLimit ...
	RateLimit ErrorCode = "rate_limit"

	// APIErr ...
	APIErr ErrorType = ErrorTypeAPI
	// CardErr ...
	CardErr ErrorType = ErrorTypeCard
	// InvalidRequest ...
	InvalidRequest ErrorType = ErrorTypeInvalidRequest
)

// Error is the response returned when a call is unsuccessful.
type Error struct {
	Type           ErrorType `json:"type"`
	Msg            string    `json:"message"`
	Code           ErrorCode `json:"code,omitempty"`
	Param          string    `json:"param,omitempty"`
	RequestID      string    `json:"request_id,omitempty"`
	HTTPStatusCode int       `json:"status,omitempty"`
	ChargeID       string    `json:"charge,omitempty"`

	// Err contains an internal error with an additional level of granularity
	// that can be used in some cases to get more detailed information about
	// what went wrong. For example, Err may hold a ChargeError that indicates
	// exactly what went wrong during a charge.
	Err error `json:"-"`
}

// Error serializes the error object to JSON and returns it as a string.
func (e *Error) Error() string {
	ret, _ := json.Marshal(e)
	return string(ret)
}

// APIConnectionError is a failure to connect to the Icheck API.
type APIConnectionError struct {
	icheckErr *Error
}

// Error serializes the error object to JSON and returns it as a string.
func (e *APIConnectionError) Error() string {
	return e.icheckErr.Error()
}

// APIError is a catch all for any errors not covered by other types (and
// should be extremely uncommon).
type APIError struct {
	icheckErr *Error
}

// Error serializes the error object to JSON and returns it as a string.
func (e *APIError) Error() string {
	return e.icheckErr.Error()
}

// AuthenticationError is a failure to properly authenticate during a request.
type AuthenticationError struct {
	icheckErr *Error
}

// Error serializes the error object to JSON and returns it as a string.
func (e *AuthenticationError) Error() string {
	return e.icheckErr.Error()
}

// PermissionError results when you attempt to make an API request
// for which your API key doesn't have the right permissions.
type PermissionError struct {
	icheckErr *Error
}

// Error serializes the error object to JSON and returns it as a string.
func (e *PermissionError) Error() string {
	return e.icheckErr.Error()
}

// CardError are the most common type of error you should expect to handle.
// They result when the user enters a card that can't be charged for some
// reason.
type CardError struct {
	icheckErr   *Error
	DeclineCode string `json:"decline_code,omitempty"`
}

// Error serializes the error object to JSON and returns it as a string.
func (e *CardError) Error() string {
	return e.icheckErr.Error()
}

// InvalidRequestError is an error that occurs when a request contains invalid
// parameters.
type InvalidRequestError struct {
	icheckErr *Error
}

// Error serializes the error object to JSON and returns it as a string.
func (e *InvalidRequestError) Error() string {
	return e.icheckErr.Error()
}

// RateLimitError occurs when the Icheck API is hit to with too many requests
// too quickly and indicates that the current request has been rate limited.
type RateLimitError struct {
	icheckErr *Error
}

// Error serializes the error object to JSON and returns it as a string.
func (e *RateLimitError) Error() string {
	return e.icheckErr.Error()
}

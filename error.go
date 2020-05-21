package shoptypewooCommerce

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Error struct {
	ActualError error  `json:"actual_error,omitempty"`
	StatusCode  int    `json:"status_code,omitempty"`
	Message     string `json:"message,omitempty"`
}

// Error implements the error interface.
func (e Error) Error() string {
	if v, err := json.Marshal(e); err == nil {
		return string(v)
	}

	return "shoptype-wooCommerce error"
}

// Ref: https://woocommerce.github.io/woocommerce-rest-api-docs/?shell#errors
//400 Bad Request	Invalid request, e.g. using an unsupported HTTP method
//401 Unauthorized	Authentication or permission error, e.g. incorrect API keys
//404 Not Found	Requests to resources that don't exist or are missing
//500 Internal Server Error	Server error
func NewError(err error, statusCode int, messages ...string) Error {
	if statusCode == http.StatusBadRequest {
		return ErrorBadRequest(err, messages...)
	}

	if statusCode == http.StatusUnauthorized {
		return ErrorUnauthorized(err, messages...)
	}

	if statusCode == http.StatusNotFound {
		return ErrorNotFound(err, messages...)
	}

	return ErrorInternal(err, messages...)

}

// ErrorNotFound returns a not found error.
func ErrorNotFound(err error, messages ...string) Error {
	return Error{
		ActualError: err,
		StatusCode:  http.StatusNotFound,
		Message:     strings.Join(messages, " "),
	}
}

// ErrorInternal returns an internal error.
func ErrorInternal(err error, messages ...string) Error {
	return Error{
		ActualError: err,
		StatusCode:  http.StatusInternalServerError,
		Message:     strings.Join(messages, " "),
	}
}

// ErrorBadRequest returns an bad request error.
func ErrorBadRequest(err error, messages ...string) Error {
	return Error{
		ActualError: err,
		StatusCode:  http.StatusBadRequest,
		Message:     strings.Join(messages, " "),
	}
}

// ErrorUnauthorized returns a unauthorized request error.
func ErrorUnauthorized(err error, messages ...string) Error {
	return Error{
		ActualError: err,
		StatusCode:  http.StatusUnauthorized,
		Message:     strings.Join(messages, " "),
	}
}

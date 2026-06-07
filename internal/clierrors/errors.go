package clierrors

import (
	"errors"
	"fmt"
)

type Error struct {
	Code       string
	Message    string
	Exit       int
	HTTPStatus int
	RequestID  string
	Cause      error
}

func (e *Error) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("Error: %s\nCode: %s\nHTTP status: %d\nRequest ID: %s", e.Message, e.Code, e.HTTPStatus, e.RequestID)
	}
	if e.HTTPStatus != 0 {
		return fmt.Sprintf("Error: %s\nCode: %s\nHTTP status: %d", e.Message, e.Code, e.HTTPStatus)
	}
	return fmt.Sprintf("Error: %s\nCode: %s", e.Message, e.Code)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func New(code, message string, exit int) error {
	return &Error{Code: code, Message: message, Exit: exit}
}

func Validation(message string) error {
	return New("validation_failed", message, 2)
}

func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var cliErr *Error
	if errors.As(err, &cliErr) && cliErr.Exit != 0 {
		return cliErr.Exit
	}
	return 1
}

func FromHTTP(status int, message, requestID string) error {
	code, exit := "unexpected_response", 1
	switch status {
	case 400, 422:
		code, exit = "validation_failed", 2
	case 401:
		code, exit = "authentication_required", 3
	case 403:
		code, exit = "permission_denied", 4
	case 404:
		code, exit = "not_found", 5
	case 409:
		code, exit = "conflict", 6
	case 408:
		code, exit = "timeout", 7
	case 429:
		code, exit = "rate_limited", 9
	case 502, 503, 504:
		code, exit = "backend_unavailable", 8
	case 500:
		code, exit = "internal_error", 1
	}
	return &Error{Code: code, Message: message, Exit: exit, HTTPStatus: status, RequestID: requestID}
}

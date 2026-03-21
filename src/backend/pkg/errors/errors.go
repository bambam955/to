// Package errors provides error types and categorization for the to tool.
package errors

import (
	"fmt"
	"os"
)

// ErrorType categorizes the kind of error that occurred.
type ErrorType string

const (
	// ErrorTypeNotFound indicates a resource (alias, directory) was not found.
	ErrorTypeNotFound ErrorType = "not_found"
	// ErrorTypeExists indicates a resource already exists.
	ErrorTypeExists ErrorType = "exists"
	// ErrorTypeInvalid indicates invalid input (format, validation).
	ErrorTypeInvalid ErrorType = "invalid"
	// ErrorTypeDatabase indicates database operation failures.
	ErrorTypeDatabase ErrorType = "database"
	// ErrorTypePermission indicates permission or access issues.
	ErrorTypePermission ErrorType = "permission"
	// ErrorTypeCorrupted indicates corrupted or malformed data.
	ErrorTypeCorrupted ErrorType = "corrupted"
	// ErrorTypeInternal indicates internal errors (should not happen).
	ErrorTypeInternal ErrorType = "internal"
)

// Error is a structured error type for the to tool.
type Error struct {
	Type    ErrorType
	Message string
	Cause   error
}

// Error implements the error interface.
func (e *Error) Error() string {
	return e.Message
}

// Unwrap returns the underlying error cause.
func (e *Error) Unwrap() error {
	return e.Cause
}

// New creates a new Error with the given type and message.
func New(errType ErrorType, message string) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Cause:   nil,
	}
}

// Wrap creates a new Error wrapping an existing error.
func Wrap(errType ErrorType, message string, cause error) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Cause:   cause,
	}
}

// Newf creates a new Error with formatted message.
func Newf(errType ErrorType, format string, args ...interface{}) *Error {
	return &Error{
		Type:    errType,
		Message: fmt.Sprintf(format, args...),
		Cause:   nil,
	}
}

// NotFound creates an error for when a resource is not found.
func NotFound(resource string) *Error {
	return New(ErrorTypeNotFound, fmt.Sprintf("%s not found", resource))
}

// AlreadyExists creates an error for when a resource already exists.
func AlreadyExists(resource string) *Error {
	return New(ErrorTypeExists, fmt.Sprintf("%s already exists", resource))
}

// InvalidInput creates an error for invalid input.
func InvalidInput(reason string) *Error {
	return New(ErrorTypeInvalid, fmt.Sprintf("invalid input: %s", reason))
}

// PermissionDenied creates an error for permission issues.
func PermissionDenied(path string) *Error {
	return New(ErrorTypePermission, fmt.Sprintf("permission denied: cannot access %s", path))
}

// Fail logs an error and exits with code 1.
// The message is printed to stderr with the "error: " prefix.
func Fail(err error) {
	message := err.Error()
	fmt.Fprintf(os.Stderr, "error: %s\n", message)
	os.Exit(1)
}

// Failf logs a formatted error and exits with code 1.
func Failf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: %s\n", fmt.Sprintf(format, args...))
	os.Exit(1)
}

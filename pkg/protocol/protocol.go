// Package protocol defines the communication protocol between the bash wrapper
// and the Go backend.
package protocol

import (
	"fmt"
)

// NavigationResponse outputs the special format for successful navigation.
// Format: [to] <absolute_path>
// This is parsed by the bash wrapper to extract the directory path for cd.
func NavigationResponse(path string) string {
	return fmt.Sprintf("[to] %s", path)
}

// ErrorResponse formats error messages with the "error:" prefix.
// These are printed to stderr and not parsed by the bash wrapper.
func ErrorResponse(message string) string {
	return fmt.Sprintf("error: %s", message)
}

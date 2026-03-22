// Package protocol defines the communication protocol between the bash wrapper
// and the Go backend.
package protocol

import (
	"fmt"
	"io"
	"strings"
)

const (
	// NavigationTag marks a control-frame line that carries a directory path
	// for shell wrappers to `cd` into after backend command completion.
	NavigationTag = "NAV"
)

// NavigationControlFrame formats the navigation payload emitted on fd 3.
// Format: NAV <absolute_path>
func NavigationControlFrame(path string) string {
	return fmt.Sprintf("%s %s", NavigationTag, path)
}

// ParseNavigationControlFrame parses a single control-frame line and returns
// the path payload when the line is a valid navigation frame.
func ParseNavigationControlFrame(line string) (string, bool) {
	line = strings.TrimSpace(line)
	prefix := NavigationTag + " "
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}

	path := strings.TrimSpace(strings.TrimPrefix(line, prefix))
	if path == "" {
		return "", false
	}

	return path, true
}

// WriteNavigationControlFrame writes a navigation control frame to the given
// writer, including a trailing newline so shell wrappers can read line-by-line.
func WriteNavigationControlFrame(w io.Writer, path string) error {
	_, err := io.WriteString(w, NavigationControlFrame(path)+"\n")
	return err
}

// ErrorResponse formats error messages with the "error:" prefix.
// These are printed to stderr and not parsed by the bash wrapper.
func ErrorResponse(message string) string {
	return fmt.Sprintf("error: %s", message)
}

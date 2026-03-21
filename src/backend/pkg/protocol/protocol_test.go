package protocol

import (
	"testing"
)

func TestNavigationResponse(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple path",
			path:     "/home/user",
			expected: "[to] /home/user",
		},
		{
			name:     "root path",
			path:     "/",
			expected: "[to] /",
		},
		{
			name:     "nested path",
			path:     "/home/user/projects/myproject",
			expected: "[to] /home/user/projects/myproject",
		},
		{
			name:     "path with spaces",
			path:     "/home/user/my documents",
			expected: "[to] /home/user/my documents",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NavigationResponse(tt.path)
			if result != tt.expected {
				t.Errorf("NavigationResponse(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "simple error",
			message:  "alias not found",
			expected: "error: alias not found",
		},
		{
			name:     "error with details",
			message:  "directory /nonexistent does not exist",
			expected: "error: directory /nonexistent does not exist",
		},
		{
			name:     "permission error",
			message:  "permission denied: cannot access /root",
			expected: "error: permission denied: cannot access /root",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ErrorResponse(tt.message)
			if result != tt.expected {
				t.Errorf("ErrorResponse(%q) = %q, want %q", tt.message, result, tt.expected)
			}
		})
	}
}

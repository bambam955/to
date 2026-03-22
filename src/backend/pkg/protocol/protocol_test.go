package protocol

import (
	"testing"
)

func TestNavigationControlFrame(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple path",
			path:     "/home/user",
			expected: "NAV /home/user",
		},
		{
			name:     "root path",
			path:     "/",
			expected: "NAV /",
		},
		{
			name:     "nested path",
			path:     "/home/user/projects/myproject",
			expected: "NAV /home/user/projects/myproject",
		},
		{
			name:     "path with spaces",
			path:     "/home/user/my documents",
			expected: "NAV /home/user/my documents",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NavigationControlFrame(tt.path)
			if result != tt.expected {
				t.Errorf("NavigationControlFrame(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestParseNavigationControlFrame(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		ok       bool
	}{
		{
			name:     "valid frame",
			input:    "NAV /tmp",
			expected: "/tmp",
			ok:       true,
		},
		{
			name:     "valid frame with trailing newline",
			input:    "NAV /tmp\n",
			expected: "/tmp",
			ok:       true,
		},
		{
			name:  "invalid tag",
			input: "NOPE /tmp",
			ok:    false,
		},
		{
			name:  "missing path",
			input: "NAV",
			ok:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, ok := ParseNavigationControlFrame(tt.input)
			if ok != tt.ok {
				t.Fatalf("ParseNavigationControlFrame(%q) ok = %v, want %v", tt.input, ok, tt.ok)
			}
			if path != tt.expected {
				t.Fatalf("ParseNavigationControlFrame(%q) path = %q, want %q", tt.input, path, tt.expected)
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

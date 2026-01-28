package errors

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(ErrorTypeNotFound, "alias not found")
	if err.Type != ErrorTypeNotFound {
		t.Errorf("Error type = %v, want %v", err.Type, ErrorTypeNotFound)
	}
	if err.Message != "alias not found" {
		t.Errorf("Error message = %q, want %q", err.Message, "alias not found")
	}
	if err.Error() != "alias not found" {
		t.Errorf("Error() = %q, want %q", err.Error(), "alias not found")
	}
}

func TestWrap(t *testing.T) {
	cause := errors.New("file not found")
	err := Wrap(ErrorTypeDatabase, "failed to read database", cause)
	if err.Type != ErrorTypeDatabase {
		t.Errorf("Error type = %v, want %v", err.Type, ErrorTypeDatabase)
	}
	if err.Cause != cause {
		t.Errorf("Error cause not preserved")
	}
	if !errors.Is(err, cause) {
		t.Errorf("Error chain broken, errors.Is failed")
	}
}

func TestNewf(t *testing.T) {
	err := Newf(ErrorTypeInvalid, "invalid value: %d", 42)
	if err.Type != ErrorTypeInvalid {
		t.Errorf("Error type = %v, want %v", err.Type, ErrorTypeInvalid)
	}
	if err.Message != "invalid value: 42" {
		t.Errorf("Error message = %q, want %q", err.Message, "invalid value: 42")
	}
}

func TestNotFound(t *testing.T) {
	err := NotFound("alias 'myalias'")
	if err.Type != ErrorTypeNotFound {
		t.Errorf("Error type = %v, want %v", err.Type, ErrorTypeNotFound)
	}
	if err.Message != "alias 'myalias' not found" {
		t.Errorf("Error message = %q, want %q", err.Message, "alias 'myalias' not found")
	}
}

func TestAlreadyExists(t *testing.T) {
	err := AlreadyExists("alias 'work'")
	if err.Type != ErrorTypeExists {
		t.Errorf("Error type = %v, want %v", err.Type, ErrorTypeExists)
	}
	if err.Message != "alias 'work' already exists" {
		t.Errorf("Error message = %q, want %q", err.Message, "alias 'work' already exists")
	}
}

func TestInvalidInput(t *testing.T) {
	err := InvalidInput("alias name contains invalid characters")
	if err.Type != ErrorTypeInvalid {
		t.Errorf("Error type = %v, want %v", err.Type, ErrorTypeInvalid)
	}
	if err.Message != "invalid input: alias name contains invalid characters" {
		t.Errorf("Error message = %q, want starting with 'invalid input:'", err.Message)
	}
}

func TestPermissionDenied(t *testing.T) {
	err := PermissionDenied("/root")
	if err.Type != ErrorTypePermission {
		t.Errorf("Error type = %v, want %v", err.Type, ErrorTypePermission)
	}
	if err.Message != "permission denied: cannot access /root" {
		t.Errorf("Error message = %q, want %q", err.Message, "permission denied: cannot access /root")
	}
}

func TestErrorInterface(t *testing.T) {
	var e error
	e = New(ErrorTypeNotFound, "test error")
	if e.Error() != "test error" {
		t.Errorf("Error interface not properly implemented")
	}
}

func TestErrorWrapping(t *testing.T) {
	cause := errors.New("original error")
	err := Wrap(ErrorTypeDatabase, "wrapped error", cause)
	if errors.Unwrap(err) != cause {
		t.Errorf("Error wrapping not working properly")
	}
}

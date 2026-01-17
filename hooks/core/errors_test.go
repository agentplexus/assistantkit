package core

import (
	"errors"
	"testing"
)

func TestHookValidationError(t *testing.T) {
	innerErr := ErrNoCommandOrPrompt
	err := &HookValidationError{
		Event:      BeforeCommand,
		EntryIndex: 0,
		HookIndex:  1,
		Err:        innerErr,
	}

	// Test Error() string
	errStr := err.Error()
	if errStr == "" {
		t.Error("Error() returned empty string")
	}
	// Should contain event name, indices, and inner error
	if !containsAll(errStr, "before_command", "entry 0", "hook 1") {
		t.Errorf("Error() missing expected content: %s", errStr)
	}

	// Test Unwrap()
	unwrapped := err.Unwrap()
	if unwrapped != innerErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, innerErr)
	}

	// Test errors.Is
	if !errors.Is(err, ErrNoCommandOrPrompt) {
		t.Error("errors.Is should match inner error")
	}
}

func TestParseError(t *testing.T) {
	innerErr := errors.New("json syntax error")

	// Test with path
	errWithPath := &ParseError{
		Format: "claude",
		Path:   "/path/to/file.json",
		Err:    innerErr,
	}
	errStr := errWithPath.Error()
	if !containsAll(errStr, "claude", "/path/to/file.json", "json syntax error") {
		t.Errorf("Error() with path missing expected content: %s", errStr)
	}

	// Test without path
	errNoPath := &ParseError{
		Format: "cursor",
		Err:    innerErr,
	}
	errStr = errNoPath.Error()
	if !containsAll(errStr, "cursor", "json syntax error") {
		t.Errorf("Error() without path missing expected content: %s", errStr)
	}
	// Should not contain "from" when no path
	if containsAll(errStr, "from /") {
		t.Error("Error() without path should not contain path separator")
	}

	// Test Unwrap()
	if errWithPath.Unwrap() != innerErr {
		t.Errorf("Unwrap() = %v, want %v", errWithPath.Unwrap(), innerErr)
	}
}

func TestWriteError(t *testing.T) {
	innerErr := errors.New("permission denied")
	err := &WriteError{
		Format: "windsurf",
		Path:   "/tmp/hooks.json",
		Err:    innerErr,
	}

	errStr := err.Error()
	if !containsAll(errStr, "windsurf", "/tmp/hooks.json", "permission denied") {
		t.Errorf("Error() missing expected content: %s", errStr)
	}

	// Test Unwrap()
	if err.Unwrap() != innerErr {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), innerErr)
	}

	// Test errors.Is
	if !errors.Is(err, innerErr) {
		t.Error("errors.Is should match inner error")
	}
}

func TestConversionError(t *testing.T) {
	innerErr := ErrUnsupportedEvent

	// Test with event
	errWithEvent := &ConversionError{
		From:  "claude",
		To:    "cursor",
		Event: OnSessionStart,
		Err:   innerErr,
	}
	errStr := errWithEvent.Error()
	if !containsAll(errStr, "claude", "cursor", "on_session_start") {
		t.Errorf("Error() with event missing expected content: %s", errStr)
	}

	// Test without event
	errNoEvent := &ConversionError{
		From: "claude",
		To:   "windsurf",
		Err:  innerErr,
	}
	errStr = errNoEvent.Error()
	if !containsAll(errStr, "claude", "windsurf") {
		t.Errorf("Error() without event missing expected content: %s", errStr)
	}

	// Test Unwrap()
	if errWithEvent.Unwrap() != innerErr {
		t.Errorf("Unwrap() = %v, want %v", errWithEvent.Unwrap(), innerErr)
	}
}

func TestCommonErrors(t *testing.T) {
	// Just ensure the common errors are defined and have messages
	errs := []error{
		ErrNoCommandOrPrompt,
		ErrBothCommandAndPrompt,
		ErrUnsupportedEvent,
		ErrInvalidMatcher,
		ErrEmptyConfig,
	}

	for _, err := range errs {
		if err == nil {
			t.Error("Common error should not be nil")
		}
		if err.Error() == "" {
			t.Errorf("Error %v has empty message", err)
		}
	}
}

// Helper function to check if string contains all substrings
func containsAll(s string, substrs ...string) bool {
	for _, substr := range substrs {
		found := false
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

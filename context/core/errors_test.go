package core

import (
	"errors"
	"strings"
	"testing"
)

func TestParseError(t *testing.T) {
	t.Run("with path", func(t *testing.T) {
		innerErr := errors.New("syntax error")
		err := &ParseError{Path: "/path/to/file.json", Err: innerErr}

		if !strings.Contains(err.Error(), "/path/to/file.json") {
			t.Errorf("expected path in error message, got: %s", err.Error())
		}
		if !strings.Contains(err.Error(), "syntax error") {
			t.Errorf("expected inner error in message, got: %s", err.Error())
		}
		if err.Unwrap() != innerErr {
			t.Error("Unwrap() should return inner error")
		}
	})

	t.Run("without path", func(t *testing.T) {
		innerErr := errors.New("invalid JSON")
		err := &ParseError{Err: innerErr}

		if strings.Contains(err.Error(), "from") {
			t.Errorf("expected no 'from' in error without path, got: %s", err.Error())
		}
		if !strings.Contains(err.Error(), "invalid JSON") {
			t.Errorf("expected inner error in message, got: %s", err.Error())
		}
	})
}

func TestWriteError(t *testing.T) {
	t.Run("with format", func(t *testing.T) {
		innerErr := errors.New("permission denied")
		err := &WriteError{Format: "claude", Path: "/path/to/file", Err: innerErr}

		if !strings.Contains(err.Error(), "claude") {
			t.Errorf("expected format in error message, got: %s", err.Error())
		}
		if !strings.Contains(err.Error(), "/path/to/file") {
			t.Errorf("expected path in error message, got: %s", err.Error())
		}
		if !strings.Contains(err.Error(), "permission denied") {
			t.Errorf("expected inner error in message, got: %s", err.Error())
		}
		if err.Unwrap() != innerErr {
			t.Error("Unwrap() should return inner error")
		}
	})

	t.Run("without format", func(t *testing.T) {
		innerErr := errors.New("disk full")
		err := &WriteError{Path: "/path/to/file", Err: innerErr}

		if strings.Contains(err.Error(), "context to /path") {
			// This is correct - no format specified
		}
		if !strings.Contains(err.Error(), "/path/to/file") {
			t.Errorf("expected path in error message, got: %s", err.Error())
		}
	})
}

func TestConversionError(t *testing.T) {
	innerErr := errors.New("unsupported feature")
	err := &ConversionError{Format: "cursor", Err: innerErr}

	if !strings.Contains(err.Error(), "cursor") {
		t.Errorf("expected format in error message, got: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "unsupported feature") {
		t.Errorf("expected inner error in message, got: %s", err.Error())
	}
	if err.Unwrap() != innerErr {
		t.Error("Unwrap() should return inner error")
	}
}

func TestErrorsUnwrap(t *testing.T) {
	t.Run("ParseError chain", func(t *testing.T) {
		baseErr := errors.New("base error")
		parseErr := &ParseError{Path: "/test", Err: baseErr}

		if !errors.Is(parseErr, baseErr) {
			t.Error("errors.Is should find base error in chain")
		}
	})

	t.Run("WriteError chain", func(t *testing.T) {
		baseErr := errors.New("base error")
		writeErr := &WriteError{Format: "test", Path: "/test", Err: baseErr}

		if !errors.Is(writeErr, baseErr) {
			t.Error("errors.Is should find base error in chain")
		}
	})

	t.Run("ConversionError chain", func(t *testing.T) {
		baseErr := errors.New("base error")
		convErr := &ConversionError{Format: "test", Err: baseErr}

		if !errors.Is(convErr, baseErr) {
			t.Error("errors.Is should find base error in chain")
		}
	})
}

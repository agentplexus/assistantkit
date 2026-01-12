package core

import (
	"errors"
	"fmt"
)

// Common errors for context operations.
var (
	// ErrEmptyContext is returned when the context is empty.
	ErrEmptyContext = errors.New("context is empty")

	// ErrMissingName is returned when the context name is missing.
	ErrMissingName = errors.New("context name is required")

	// ErrUnsupportedFormat is returned when a format is not supported.
	ErrUnsupportedFormat = errors.New("unsupported output format")
)

// ParseError represents an error parsing a context file.
type ParseError struct {
	Path string
	Err  error
}

func (e *ParseError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("failed to parse context from %s: %v", e.Path, e.Err)
	}
	return fmt.Sprintf("failed to parse context: %v", e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// WriteError represents an error writing a context file.
type WriteError struct {
	Format string
	Path   string
	Err    error
}

func (e *WriteError) Error() string {
	if e.Format != "" {
		return fmt.Sprintf("failed to write %s context to %s: %v", e.Format, e.Path, e.Err)
	}
	return fmt.Sprintf("failed to write context to %s: %v", e.Path, e.Err)
}

func (e *WriteError) Unwrap() error {
	return e.Err
}

// ConversionError represents an error converting to a specific format.
type ConversionError struct {
	Format string
	Err    error
}

func (e *ConversionError) Error() string {
	return fmt.Sprintf("failed to convert context to %s format: %v", e.Format, e.Err)
}

func (e *ConversionError) Unwrap() error {
	return e.Err
}

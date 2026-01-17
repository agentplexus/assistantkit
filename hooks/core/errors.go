package core

import (
	"errors"
	"fmt"
)

// Common errors for hooks configuration.
var (
	// ErrNoCommandOrPrompt is returned when a hook has neither command nor prompt.
	ErrNoCommandOrPrompt = errors.New("hook must have either command or prompt")

	// ErrBothCommandAndPrompt is returned when a hook has both command and prompt.
	ErrBothCommandAndPrompt = errors.New("hook cannot have both command and prompt")

	// ErrUnsupportedEvent is returned when an event is not supported by a tool.
	ErrUnsupportedEvent = errors.New("event not supported by this tool")

	// ErrInvalidMatcher is returned when a matcher pattern is invalid.
	ErrInvalidMatcher = errors.New("invalid matcher pattern")

	// ErrEmptyConfig is returned when configuration is empty.
	ErrEmptyConfig = errors.New("configuration is empty")
)

// HookValidationError wraps a validation error with context.
type HookValidationError struct {
	Event      Event
	EntryIndex int
	HookIndex  int
	Err        error
}

func (e *HookValidationError) Error() string {
	return fmt.Sprintf("hook validation error for event %q (entry %d, hook %d): %v",
		e.Event, e.EntryIndex, e.HookIndex, e.Err)
}

func (e *HookValidationError) Unwrap() error {
	return e.Err
}

// ParseError represents an error parsing a configuration file.
type ParseError struct {
	Format string
	Path   string
	Err    error
}

func (e *ParseError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("failed to parse %s hooks config from %s: %v", e.Format, e.Path, e.Err)
	}
	return fmt.Sprintf("failed to parse %s hooks config: %v", e.Format, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// WriteError represents an error writing a configuration file.
type WriteError struct {
	Format string
	Path   string
	Err    error
}

func (e *WriteError) Error() string {
	return fmt.Sprintf("failed to write %s hooks config to %s: %v", e.Format, e.Path, e.Err)
}

func (e *WriteError) Unwrap() error {
	return e.Err
}

// ConversionError represents an error converting between formats.
type ConversionError struct {
	From  string
	To    string
	Event Event
	Err   error
}

func (e *ConversionError) Error() string {
	if e.Event != "" {
		return fmt.Sprintf("failed to convert event %q from %s to %s: %v",
			e.Event, e.From, e.To, e.Err)
	}
	return fmt.Sprintf("failed to convert from %s to %s: %v", e.From, e.To, e.Err)
}

func (e *ConversionError) Unwrap() error {
	return e.Err
}

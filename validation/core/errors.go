package core

import "fmt"

// ReadError represents an error reading a file.
type ReadError struct {
	Path string
	Err  error
}

func (e *ReadError) Error() string {
	return fmt.Sprintf("failed to read %s: %v", e.Path, e.Err)
}

func (e *ReadError) Unwrap() error {
	return e.Err
}

// WriteError represents an error writing a file.
type WriteError struct {
	Path string
	Err  error
}

func (e *WriteError) Error() string {
	return fmt.Sprintf("failed to write %s: %v", e.Path, e.Err)
}

func (e *WriteError) Unwrap() error {
	return e.Err
}

// ParseError represents an error parsing a file format.
type ParseError struct {
	Format string
	Path   string
	Err    error
}

func (e *ParseError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("failed to parse %s format in %s: %v", e.Format, e.Path, e.Err)
	}
	return fmt.Sprintf("failed to parse %s format: %v", e.Format, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// MarshalError represents an error marshaling to a format.
type MarshalError struct {
	Format string
	Err    error
}

func (e *MarshalError) Error() string {
	return fmt.Sprintf("failed to marshal to %s format: %v", e.Format, e.Err)
}

func (e *MarshalError) Unwrap() error {
	return e.Err
}

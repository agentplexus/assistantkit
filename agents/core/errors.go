package core

import "fmt"

// ReadError indicates a failure to read a file.
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

// WriteError indicates a failure to write a file.
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

// ParseError indicates a failure to parse agent data.
type ParseError struct {
	Format string
	Path   string
	Err    error
}

func (e *ParseError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("failed to parse %s format from %s: %v", e.Format, e.Path, e.Err)
	}
	return fmt.Sprintf("failed to parse %s format: %v", e.Format, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// MarshalError indicates a failure to marshal agent data.
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

// AdapterError indicates an unknown adapter was requested.
type AdapterError struct {
	Name string
}

func (e *AdapterError) Error() string {
	return fmt.Sprintf("unknown adapter: %s", e.Name)
}

package bundle

import "fmt"

// GenerateError represents an error during bundle generation.
type GenerateError struct {
	Tool      string
	Component string
	Err       error
}

func (e *GenerateError) Error() string {
	if e.Component != "" {
		return fmt.Sprintf("bundle generate %s/%s: %v", e.Tool, e.Component, e.Err)
	}
	return fmt.Sprintf("bundle generate %s: %v", e.Tool, e.Err)
}

func (e *GenerateError) Unwrap() error {
	return e.Err
}

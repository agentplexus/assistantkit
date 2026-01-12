package core

import (
	"io/fs"
	"os"
)

// DefaultFileMode is the default permission mode for generated files.
// This can be used by converters or overridden with WriteFileWithDataAndMode.
const DefaultFileMode fs.FileMode = 0600

// Converter defines the interface for converting project context
// to tool-specific formats.
type Converter interface {
	// Name returns the converter name (e.g., "claude", "cursor").
	Name() string

	// OutputFileName returns the default output file name (e.g., "CLAUDE.md").
	OutputFileName() string

	// Convert converts the context to the tool-specific format.
	Convert(ctx *Context) ([]byte, error)

	// WriteFile writes the converted context to a file.
	WriteFile(ctx *Context, path string) error
}

// ConverterRegistry holds registered converters for different tools.
type ConverterRegistry struct {
	converters map[string]Converter
}

// NewConverterRegistry creates a new converter registry.
func NewConverterRegistry() *ConverterRegistry {
	return &ConverterRegistry{
		converters: make(map[string]Converter),
	}
}

// Register adds a converter to the registry.
func (r *ConverterRegistry) Register(converter Converter) {
	r.converters[converter.Name()] = converter
}

// Get returns a converter by name.
func (r *ConverterRegistry) Get(name string) (Converter, bool) {
	converter, ok := r.converters[name]
	return converter, ok
}

// Names returns the names of all registered converters.
func (r *ConverterRegistry) Names() []string {
	names := make([]string, 0, len(r.converters))
	for name := range r.converters {
		names = append(names, name)
	}
	return names
}

// Convert converts a context to a specific format.
func (r *ConverterRegistry) Convert(ctx *Context, format string) ([]byte, error) {
	converter, ok := r.Get(format)
	if !ok {
		return nil, &ConversionError{Format: format, Err: ErrUnsupportedFormat}
	}
	return converter.Convert(ctx)
}

// WriteFile writes a context to a file in a specific format.
func (r *ConverterRegistry) WriteFile(ctx *Context, format, path string) error {
	converter, ok := r.Get(format)
	if !ok {
		return &ConversionError{Format: format, Err: ErrUnsupportedFormat}
	}
	return converter.WriteFile(ctx, path)
}

// GenerateAll generates all supported formats in the given directory.
func (r *ConverterRegistry) GenerateAll(ctx *Context, dir string) error {
	for _, converter := range r.converters {
		var path string
		if dir != "" {
			path = dir + "/" + converter.OutputFileName()
		} else {
			path = converter.OutputFileName()
		}
		if err := converter.WriteFile(ctx, path); err != nil {
			return err
		}
	}
	return nil
}

// DefaultRegistry is the global converter registry.
var DefaultRegistry = NewConverterRegistry()

// RegisterConverter adds a converter to the default registry.
func RegisterConverter(converter Converter) {
	DefaultRegistry.Register(converter)
}

// GetConverter returns a converter from the default registry.
func GetConverter(name string) (Converter, bool) {
	return DefaultRegistry.Get(name)
}

// ConvertTo converts a context to a specific format using the default registry.
func ConvertTo(ctx *Context, format string) ([]byte, error) {
	return DefaultRegistry.Convert(ctx, format)
}

// BaseConverter provides common functionality for converters.
type BaseConverter struct {
	name       string
	outputFile string
}

// NewBaseConverter creates a new base converter.
func NewBaseConverter(name, outputFile string) BaseConverter {
	return BaseConverter{name: name, outputFile: outputFile}
}

// Name returns the converter name.
func (c *BaseConverter) Name() string {
	return c.name
}

// OutputFileName returns the default output file name.
func (c *BaseConverter) OutputFileName() string {
	return c.outputFile
}

// WriteFileWithData writes data to a file with proper error wrapping using DefaultFileMode.
func (c *BaseConverter) WriteFileWithData(data []byte, path string) error {
	return c.WriteFileWithDataAndMode(data, path, DefaultFileMode)
}

// WriteFileWithDataAndMode writes data to a file with proper error wrapping using the specified permission mode.
func (c *BaseConverter) WriteFileWithDataAndMode(data []byte, path string, mode fs.FileMode) error {
	if err := os.WriteFile(path, data, mode); err != nil {
		return &WriteError{Format: c.name, Path: path, Err: err}
	}
	return nil
}

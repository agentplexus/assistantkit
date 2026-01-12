package core

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// DefaultFileMode is the default permission for generated files.
const DefaultFileMode fs.FileMode = 0600

// DefaultDirMode is the default permission for generated directories.
const DefaultDirMode fs.FileMode = 0700

// Adapter converts between canonical ValidationArea and tool-specific formats.
type Adapter interface {
	// Name returns the adapter identifier (e.g., "claude", "gemini").
	Name() string

	// FileExtension returns the file extension for validation files.
	FileExtension() string

	// DefaultDir returns the default directory name for validation areas.
	DefaultDir() string

	// Parse converts tool-specific bytes to canonical ValidationArea.
	Parse(data []byte) (*ValidationArea, error)

	// Marshal converts canonical ValidationArea to tool-specific bytes.
	Marshal(area *ValidationArea) ([]byte, error)

	// ReadFile reads from path and returns canonical ValidationArea.
	ReadFile(path string) (*ValidationArea, error)

	// WriteFile writes canonical ValidationArea to path.
	WriteFile(area *ValidationArea, path string) error
}

// Registry manages adapter registration and lookup.
type Registry struct {
	mu       sync.RWMutex
	adapters map[string]Adapter
}

// NewRegistry creates a new adapter registry.
func NewRegistry() *Registry {
	return &Registry{
		adapters: make(map[string]Adapter),
	}
}

// Register adds an adapter to the registry.
func (r *Registry) Register(adapter Adapter) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adapters[adapter.Name()] = adapter
}

// GetAdapter returns an adapter by name.
func (r *Registry) GetAdapter(name string) (Adapter, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	adapter, ok := r.adapters[name]
	return adapter, ok
}

// AdapterNames returns all registered adapter names sorted alphabetically.
func (r *Registry) AdapterNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// DefaultRegistry is the global adapter registry.
var DefaultRegistry = NewRegistry()

// Register adds an adapter to the default registry.
func Register(adapter Adapter) {
	DefaultRegistry.Register(adapter)
}

// GetAdapter returns an adapter from the default registry.
func GetAdapter(name string) (Adapter, bool) {
	return DefaultRegistry.GetAdapter(name)
}

// AdapterNames returns adapter names from the default registry.
func AdapterNames() []string {
	return DefaultRegistry.AdapterNames()
}

// ReadCanonicalFile reads a canonical validation-area.json file.
func ReadCanonicalFile(path string) (*ValidationArea, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &ReadError{Path: path, Err: err}
	}

	var area ValidationArea
	if err := json.Unmarshal(data, &area); err != nil {
		return nil, &ParseError{Format: "canonical", Path: path, Err: err}
	}

	return &area, nil
}

// WriteCanonicalFile writes a canonical validation-area.json file.
func WriteCanonicalFile(area *ValidationArea, path string) error {
	data, err := json.MarshalIndent(area, "", "  ")
	if err != nil {
		return &MarshalError{Format: "canonical", Err: err}
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, DefaultDirMode); err != nil {
		return &WriteError{Path: path, Err: err}
	}

	if err := os.WriteFile(path, append(data, '\n'), DefaultFileMode); err != nil {
		return &WriteError{Path: path, Err: err}
	}

	return nil
}

// ReadCanonicalDir reads all validation-area.json files from a directory.
func ReadCanonicalDir(dir string) ([]*ValidationArea, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, &ReadError{Path: dir, Err: err}
	}

	var areas []*ValidationArea
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		area, err := ReadCanonicalFile(path)
		if err != nil {
			return nil, err
		}
		areas = append(areas, area)
	}

	return areas, nil
}

// WriteAreasToDir writes multiple validation areas to a directory using the specified adapter.
func WriteAreasToDir(areas []*ValidationArea, dir string, adapterName string) error {
	adapter, ok := GetAdapter(adapterName)
	if !ok {
		return fmt.Errorf("unknown adapter: %s", adapterName)
	}

	if err := os.MkdirAll(dir, DefaultDirMode); err != nil {
		return &WriteError{Path: dir, Err: err}
	}

	for _, area := range areas {
		filename := area.Name + adapter.FileExtension()
		path := filepath.Join(dir, filename)
		if err := adapter.WriteFile(area, path); err != nil {
			return err
		}
	}

	return nil
}

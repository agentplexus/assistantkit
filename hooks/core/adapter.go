package core

// Adapter defines the interface for converting between the canonical
// Config format and tool-specific formats.
type Adapter interface {
	// Name returns the name of the tool/format (e.g., "claude", "cursor").
	Name() string

	// DefaultPaths returns the default configuration file paths for this tool.
	DefaultPaths() []string

	// Parse parses tool-specific data into the canonical Config format.
	Parse(data []byte) (*Config, error)

	// Marshal converts the canonical Config to tool-specific format.
	Marshal(cfg *Config) ([]byte, error)

	// ReadFile reads a tool-specific config file and returns canonical Config.
	ReadFile(path string) (*Config, error)

	// WriteFile writes the canonical Config to a tool-specific file.
	WriteFile(cfg *Config, path string) error

	// SupportedEvents returns the events supported by this tool.
	SupportedEvents() []Event
}

// AdapterRegistry holds registered adapters for different tools.
type AdapterRegistry struct {
	adapters map[string]Adapter
}

// NewAdapterRegistry creates a new adapter registry.
func NewAdapterRegistry() *AdapterRegistry {
	return &AdapterRegistry{
		adapters: make(map[string]Adapter),
	}
}

// Register adds an adapter to the registry.
func (r *AdapterRegistry) Register(adapter Adapter) {
	r.adapters[adapter.Name()] = adapter
}

// Get returns an adapter by name.
func (r *AdapterRegistry) Get(name string) (Adapter, bool) {
	adapter, ok := r.adapters[name]
	return adapter, ok
}

// Names returns the names of all registered adapters.
func (r *AdapterRegistry) Names() []string {
	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}
	return names
}

// Convert converts a config from one format to another.
func (r *AdapterRegistry) Convert(data []byte, from, to string) ([]byte, error) {
	fromAdapter, ok := r.Get(from)
	if !ok {
		return nil, &ConversionError{From: from, To: to, Err: ErrUnsupportedEvent}
	}
	toAdapter, ok := r.Get(to)
	if !ok {
		return nil, &ConversionError{From: from, To: to, Err: ErrUnsupportedEvent}
	}

	cfg, err := fromAdapter.Parse(data)
	if err != nil {
		return nil, err
	}

	// Filter to only events supported by the target tool
	filtered := cfg.FilterByTool(to)

	return toAdapter.Marshal(filtered)
}

// DefaultRegistry is the global adapter registry.
var DefaultRegistry = NewAdapterRegistry()

// Register adds an adapter to the default registry.
func Register(adapter Adapter) {
	DefaultRegistry.Register(adapter)
}

// GetAdapter returns an adapter from the default registry.
func GetAdapter(name string) (Adapter, bool) {
	return DefaultRegistry.Get(name)
}

// Convert converts between formats using the default registry.
func Convert(data []byte, from, to string) ([]byte, error) {
	return DefaultRegistry.Convert(data, from, to)
}

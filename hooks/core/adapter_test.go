package core

import (
	"testing"
)

// mockAdapter is a test adapter
type mockAdapter struct {
	name   string
	events []Event
}

func (m *mockAdapter) Name() string             { return m.name }
func (m *mockAdapter) DefaultPaths() []string   { return []string{".test/hooks.json"} }
func (m *mockAdapter) SupportedEvents() []Event { return m.events }

func (m *mockAdapter) Parse(data []byte) (*Config, error) {
	cfg := NewConfig()
	cfg.AddHook(BeforeCommand, NewCommandHook("echo test"))
	return cfg, nil
}

func (m *mockAdapter) Marshal(cfg *Config) ([]byte, error) {
	return []byte(`{"hooks": {}}`), nil
}

func (m *mockAdapter) ReadFile(path string) (*Config, error) {
	return m.Parse(nil)
}

func (m *mockAdapter) WriteFile(cfg *Config, path string) error {
	return nil
}

func TestNewAdapterRegistry(t *testing.T) {
	registry := NewAdapterRegistry()
	if registry == nil {
		t.Fatal("NewAdapterRegistry returned nil")
	}
	if registry.adapters == nil {
		t.Error("adapters map should be initialized")
	}
}

func TestAdapterRegistryRegister(t *testing.T) {
	registry := NewAdapterRegistry()
	adapter := &mockAdapter{name: "test", events: []Event{BeforeCommand}}

	registry.Register(adapter)

	got, ok := registry.Get("test")
	if !ok {
		t.Fatal("Registered adapter not found")
	}
	if got.Name() != "test" {
		t.Errorf("Adapter name = %q, want 'test'", got.Name())
	}
}

func TestAdapterRegistryGet(t *testing.T) {
	registry := NewAdapterRegistry()
	adapter := &mockAdapter{name: "test"}
	registry.Register(adapter)

	// Existing adapter
	got, ok := registry.Get("test")
	if !ok {
		t.Error("Expected to find 'test' adapter")
	}
	if got == nil {
		t.Error("Got nil adapter")
	}

	// Non-existing adapter
	_, ok = registry.Get("nonexistent")
	if ok {
		t.Error("Should not find 'nonexistent' adapter")
	}
}

func TestAdapterRegistryNames(t *testing.T) {
	registry := NewAdapterRegistry()
	registry.Register(&mockAdapter{name: "alpha"})
	registry.Register(&mockAdapter{name: "beta"})
	registry.Register(&mockAdapter{name: "gamma"})

	names := registry.Names()
	if len(names) != 3 {
		t.Errorf("Expected 3 names, got %d", len(names))
	}

	nameSet := make(map[string]bool)
	for _, n := range names {
		nameSet[n] = true
	}

	for _, expected := range []string{"alpha", "beta", "gamma"} {
		if !nameSet[expected] {
			t.Errorf("Expected name %q not found", expected)
		}
	}
}

func TestAdapterRegistryConvert(t *testing.T) {
	registry := NewAdapterRegistry()
	registry.Register(&mockAdapter{
		name:   "source",
		events: []Event{BeforeCommand, AfterCommand},
	})
	registry.Register(&mockAdapter{
		name:   "target",
		events: []Event{BeforeCommand}, // Only supports BeforeCommand
	})

	data := []byte(`{}`)
	result, err := registry.Convert(data, "source", "target")
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}
	if len(result) == 0 {
		t.Error("Convert() returned empty result")
	}
}

func TestAdapterRegistryConvertUnknownSource(t *testing.T) {
	registry := NewAdapterRegistry()
	registry.Register(&mockAdapter{name: "target"})

	_, err := registry.Convert([]byte(`{}`), "unknown", "target")
	if err == nil {
		t.Error("Convert() should return error for unknown source")
	}

	convErr, ok := err.(*ConversionError)
	if !ok {
		t.Errorf("Expected ConversionError, got %T", err)
	}
	if convErr.From != "unknown" || convErr.To != "target" {
		t.Errorf("ConversionError fields incorrect: from=%q, to=%q", convErr.From, convErr.To)
	}
}

func TestAdapterRegistryConvertUnknownTarget(t *testing.T) {
	registry := NewAdapterRegistry()
	registry.Register(&mockAdapter{name: "source"})

	_, err := registry.Convert([]byte(`{}`), "source", "unknown")
	if err == nil {
		t.Error("Convert() should return error for unknown target")
	}

	convErr, ok := err.(*ConversionError)
	if !ok {
		t.Errorf("Expected ConversionError, got %T", err)
	}
	if convErr.From != "source" || convErr.To != "unknown" {
		t.Errorf("ConversionError fields incorrect: from=%q, to=%q", convErr.From, convErr.To)
	}
}

func TestDefaultRegistryFunctions(t *testing.T) {
	// Test Register adds to default registry
	Register(&mockAdapter{name: "test-default"})
	adapter, ok := GetAdapter("test-default")
	if !ok {
		t.Error("Register() should add adapter to default registry")
	}
	if adapter.Name() != "test-default" {
		t.Errorf("Adapter name = %q, want 'test-default'", adapter.Name())
	}

	// Test Names() returns registered adapters
	names := DefaultRegistry.Names()
	found := false
	for _, name := range names {
		if name == "test-default" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Names() should include registered adapter")
	}
}

func TestConvertFunction(t *testing.T) {
	// Register mock adapters for testing
	Register(&mockAdapter{name: "mock-source", events: []Event{BeforeCommand}})
	Register(&mockAdapter{name: "mock-target", events: []Event{BeforeCommand}})

	// Test the package-level Convert function with mock adapters
	testData := []byte(`{}`)

	result, err := Convert(testData, "mock-source", "mock-target")
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}
	if len(result) == 0 {
		t.Error("Convert() returned empty data")
	}
}

func TestConvertFunctionUnknownAdapter(t *testing.T) {
	// Test Convert with unknown adapters
	_, err := Convert([]byte(`{}`), "unknown-source", "unknown-target")
	if err == nil {
		t.Error("Convert() should return error for unknown adapters")
	}
}

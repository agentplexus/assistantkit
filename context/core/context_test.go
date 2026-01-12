package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext("test-project")

	if ctx.Name != "test-project" {
		t.Errorf("expected name 'test-project', got '%s'", ctx.Name)
	}
	if ctx.Commands == nil {
		t.Error("expected Commands to be initialized")
	}
}

func TestContextAddPackage(t *testing.T) {
	ctx := NewContext("test")
	ctx.AddPackage("pkg/core", "Core functionality")

	if len(ctx.Packages) != 1 {
		t.Fatalf("expected 1 package, got %d", len(ctx.Packages))
	}
	if ctx.Packages[0].Path != "pkg/core" {
		t.Errorf("expected path 'pkg/core', got '%s'", ctx.Packages[0].Path)
	}
	if ctx.Packages[0].Purpose != "Core functionality" {
		t.Errorf("expected purpose 'Core functionality', got '%s'", ctx.Packages[0].Purpose)
	}
}

func TestContextAddConvention(t *testing.T) {
	ctx := NewContext("test")
	ctx.AddConvention("Use gofmt for formatting")

	if len(ctx.Conventions) != 1 {
		t.Fatalf("expected 1 convention, got %d", len(ctx.Conventions))
	}
	if ctx.Conventions[0] != "Use gofmt for formatting" {
		t.Errorf("unexpected convention: %s", ctx.Conventions[0])
	}
}

func TestContextAddNote(t *testing.T) {
	ctx := NewContext("test")
	ctx.AddNote("Simple note")
	ctx.AddNoteWithSeverity("Warning", "This is a warning", "warning")

	if len(ctx.Notes) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(ctx.Notes))
	}
	if ctx.Notes[0].Content != "Simple note" {
		t.Errorf("unexpected note content: %s", ctx.Notes[0].Content)
	}
	if ctx.Notes[1].Severity != "warning" {
		t.Errorf("expected severity 'warning', got '%s'", ctx.Notes[1].Severity)
	}
}

func TestContextSetCommand(t *testing.T) {
	ctx := NewContext("test")
	ctx.SetCommand("build", "go build ./...")
	ctx.SetCommand("test", "go test ./...")

	if len(ctx.Commands) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(ctx.Commands))
	}
	if ctx.Commands["build"] != "go build ./..." {
		t.Errorf("unexpected build command: %s", ctx.Commands["build"])
	}
}

func TestPackageIsPublic(t *testing.T) {
	tests := []struct {
		name     string
		public   *bool
		expected bool
	}{
		{"nil defaults to true", nil, true},
		{"explicit true", boolPtr(true), true},
		{"explicit false", boolPtr(false), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := Package{Path: "test", Purpose: "test", Public: tt.public}
			if pkg.IsPublic() != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, pkg.IsPublic())
			}
		})
	}
}

func TestNoteGetSeverity(t *testing.T) {
	tests := []struct {
		name     string
		severity string
		expected string
	}{
		{"empty defaults to info", "", "info"},
		{"info", "info", "info"},
		{"warning", "warning", "warning"},
		{"critical", "critical", "critical"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := Note{Content: "test", Severity: tt.severity}
			if note.GetSeverity() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, note.GetSeverity())
			}
		})
	}
}

func TestContextMarshalParse(t *testing.T) {
	ctx := NewContext("test-project")
	ctx.Description = "A test project"
	ctx.Language = "go"
	ctx.AddPackage("pkg/core", "Core functionality")
	ctx.SetCommand("build", "go build ./...")

	data, err := ctx.Marshal()
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	parsed, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if parsed.Name != ctx.Name {
		t.Errorf("expected name '%s', got '%s'", ctx.Name, parsed.Name)
	}
	if parsed.Description != ctx.Description {
		t.Errorf("expected description '%s', got '%s'", ctx.Description, parsed.Description)
	}
	if len(parsed.Packages) != len(ctx.Packages) {
		t.Errorf("expected %d packages, got %d", len(ctx.Packages), len(parsed.Packages))
	}
}

func TestContextReadWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "context.json")

	ctx := NewContext("test-project")
	ctx.Description = "A test project"
	ctx.Language = "go"

	if err := ctx.WriteFile(path); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	loaded, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if loaded.Name != ctx.Name {
		t.Errorf("expected name '%s', got '%s'", ctx.Name, loaded.Name)
	}
	if loaded.Description != ctx.Description {
		t.Errorf("expected description '%s', got '%s'", ctx.Description, loaded.Description)
	}
}

func TestReadFileNotFound(t *testing.T) {
	_, err := ReadFile("/nonexistent/path/context.json")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}

	if _, ok := err.(*ParseError); !ok {
		t.Errorf("expected ParseError, got %T", err)
	}
}

func TestParseInvalidJSON(t *testing.T) {
	_, err := Parse([]byte("invalid json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}

	if _, ok := err.(*ParseError); !ok {
		t.Errorf("expected ParseError, got %T", err)
	}
}

func TestWriteFileError(t *testing.T) {
	ctx := NewContext("test")
	err := ctx.WriteFile("/nonexistent/directory/context.json")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func TestContextWithFullData(t *testing.T) {
	ctx := &Context{
		Name:        "full-project",
		Description: "A project with all fields",
		Version:     "1.0.0",
		Language:    "go",
		Architecture: &Architecture{
			Pattern: "adapter",
			Summary: "Uses adapter pattern for extensibility",
			Diagrams: []Diagram{
				{Title: "Overview", Type: "ascii", Content: "A -> B -> C"},
			},
		},
		Packages: []Package{
			{Path: "pkg/core", Purpose: "Core types"},
			{Path: "pkg/adapter", Purpose: "Adapters"},
		},
		Commands: map[string]string{
			"build": "go build ./...",
			"test":  "go test ./...",
		},
		Conventions: []string{
			"Use gofmt",
			"Follow Go idioms",
		},
		Dependencies: &Dependencies{
			Runtime: []Dependency{
				{Name: "go-toml/v2", Purpose: "TOML parsing"},
			},
		},
		Testing: &Testing{
			Framework: "go test",
			Coverage:  "80%",
			Patterns:  []string{"Table-driven tests"},
		},
		Files: &Files{
			EntryPoints: []string{"main.go"},
			Config:      []string{"go.mod", "go.sum"},
		},
		Notes: []Note{
			{Content: "Simple note"},
			{Title: "Important", Content: "Critical info", Severity: "critical"},
		},
		Related: []Related{
			{Name: "OmniLLM", URL: "https://github.com/example/omnillm", Description: "LLM abstraction"},
		},
	}

	data, err := ctx.Marshal()
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	parsed, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify nested structures
	if parsed.Architecture == nil {
		t.Fatal("Architecture should not be nil")
	}
	if parsed.Architecture.Pattern != "adapter" {
		t.Errorf("expected pattern 'adapter', got '%s'", parsed.Architecture.Pattern)
	}
	if len(parsed.Architecture.Diagrams) != 1 {
		t.Errorf("expected 1 diagram, got %d", len(parsed.Architecture.Diagrams))
	}
	if parsed.Dependencies == nil {
		t.Fatal("Dependencies should not be nil")
	}
	if len(parsed.Dependencies.Runtime) != 1 {
		t.Errorf("expected 1 runtime dependency, got %d", len(parsed.Dependencies.Runtime))
	}
}

// Test the file written can be read by os.ReadFile
func TestWriteFileContents(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "context.json")

	ctx := NewContext("test")
	ctx.Description = "Test project"

	if err := ctx.WriteFile(path); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Read raw file
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile failed: %v", err)
	}

	// Verify it's valid JSON
	parsed, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if parsed.Name != "test" {
		t.Errorf("expected name 'test', got '%s'", parsed.Name)
	}
}

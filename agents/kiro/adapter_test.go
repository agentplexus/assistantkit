package kiro

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grokify/aiassistkit/agents/core"
)

func TestAdapter_Name(t *testing.T) {
	adapter := &Adapter{}
	if got := adapter.Name(); got != "kiro" {
		t.Errorf("Name() = %q, want %q", got, "kiro")
	}
}

func TestAdapter_FileExtension(t *testing.T) {
	adapter := &Adapter{}
	if got := adapter.FileExtension(); got != ".json" {
		t.Errorf("FileExtension() = %q, want %q", got, ".json")
	}
}

func TestAdapter_DefaultDir(t *testing.T) {
	adapter := &Adapter{}
	if got := adapter.DefaultDir(); got != "agents" {
		t.Errorf("DefaultDir() = %q, want %q", got, "agents")
	}
}

func TestAdapter_Parse(t *testing.T) {
	adapter := &Adapter{}

	input := `{
  "name": "release-agent",
  "description": "Automates software releases",
  "tools": ["read", "write", "shell"],
  "allowedTools": ["read"],
  "resources": ["file://README.md"],
  "prompt": "You are a release automation specialist.",
  "model": "claude-sonnet-4"
}`

	agent, err := adapter.Parse([]byte(input))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if agent.Name != "release-agent" {
		t.Errorf("Name = %q, want %q", agent.Name, "release-agent")
	}

	if agent.Description != "Automates software releases" {
		t.Errorf("Description = %q, want %q", agent.Description, "Automates software releases")
	}

	if agent.Model != "sonnet" {
		t.Errorf("Model = %q, want %q", agent.Model, "sonnet")
	}

	if agent.Instructions != "You are a release automation specialist." {
		t.Errorf("Instructions = %q, want %q", agent.Instructions, "You are a release automation specialist.")
	}

	// Check tools mapping
	expectedTools := []string{"Read", "Write", "Bash"}
	if len(agent.Tools) != len(expectedTools) {
		t.Errorf("Tools count = %d, want %d", len(agent.Tools), len(expectedTools))
	}
	for i, tool := range expectedTools {
		if i < len(agent.Tools) && agent.Tools[i] != tool {
			t.Errorf("Tools[%d] = %q, want %q", i, agent.Tools[i], tool)
		}
	}
}

func TestAdapter_Marshal(t *testing.T) {
	adapter := &Adapter{}

	agent := &core.Agent{
		Name:         "test-agent",
		Description:  "A test agent",
		Model:        "sonnet",
		Tools:        []string{"Read", "Write", "Bash", "Grep"},
		Skills:       []string{"version-analysis"},
		Instructions: "You are a helpful assistant.",
	}

	data, err := adapter.Marshal(agent)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	output := string(data)

	// Check key fields
	if !strings.Contains(output, `"name": "test-agent"`) {
		t.Error("Output should contain name field")
	}

	if !strings.Contains(output, `"description": "A test agent"`) {
		t.Error("Output should contain description field")
	}

	if !strings.Contains(output, `"model": "claude-sonnet-4"`) {
		t.Error("Output should contain model field with Kiro model name")
	}

	if !strings.Contains(output, `"prompt": "You are a helpful assistant."`) {
		t.Error("Output should contain prompt field")
	}

	// Check tools mapping
	if !strings.Contains(output, `"read"`) {
		t.Error("Output should contain read tool")
	}
	if !strings.Contains(output, `"shell"`) {
		t.Error("Output should contain shell tool (mapped from Bash)")
	}

	// Check skills mapped to resources
	if !strings.Contains(output, `"resources"`) {
		t.Error("Output should contain resources field")
	}
	if !strings.Contains(output, `"file://.kiro/steering/version-analysis.md"`) {
		t.Error("Output should map skills to steering files")
	}
}

func TestAdapter_RoundTrip(t *testing.T) {
	adapter := &Adapter{}

	original := &core.Agent{
		Name:         "round-trip-agent",
		Description:  "Tests round-trip conversion",
		Model:        "opus",
		Tools:        []string{"Read", "Write"},
		Instructions: "System instructions here.",
	}

	// Marshal to Kiro format
	data, err := adapter.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Parse back to canonical
	parsed, err := adapter.Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Verify fields preserved
	if parsed.Name != original.Name {
		t.Errorf("Name = %q, want %q", parsed.Name, original.Name)
	}
	if parsed.Description != original.Description {
		t.Errorf("Description = %q, want %q", parsed.Description, original.Description)
	}
	if parsed.Model != original.Model {
		t.Errorf("Model = %q, want %q", parsed.Model, original.Model)
	}
	if parsed.Instructions != original.Instructions {
		t.Errorf("Instructions = %q, want %q", parsed.Instructions, original.Instructions)
	}
}

func TestAdapter_WriteFile_ReadFile(t *testing.T) {
	adapter := &Adapter{}

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "kiro-agent-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	agent := &core.Agent{
		Name:         "file-test-agent",
		Description:  "Tests file operations",
		Model:        "haiku",
		Tools:        []string{"Read", "Grep", "Glob"},
		Instructions: "You help with file operations.",
	}

	// Write to file
	path := filepath.Join(tmpDir, "file-test-agent.json")
	if err := adapter.WriteFile(agent, path); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("WriteFile() did not create file")
	}

	// Read back
	readAgent, err := adapter.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	// Verify content
	if readAgent.Name != agent.Name {
		t.Errorf("Name = %q, want %q", readAgent.Name, agent.Name)
	}
	if readAgent.Description != agent.Description {
		t.Errorf("Description = %q, want %q", readAgent.Description, agent.Description)
	}
}

func TestModelMapping(t *testing.T) {
	tests := []struct {
		kiroModel      string
		canonicalModel string
	}{
		{"claude-sonnet-4", "sonnet"},
		{"claude-4-sonnet", "sonnet"},
		{"claude-opus-4", "opus"},
		{"claude-4-opus", "opus"},
		{"claude-haiku", "haiku"},
		{"claude-3-haiku", "haiku"},
		{"unknown-model", "unknown-model"},
	}

	for _, tt := range tests {
		got := mapKiroModelToCanonical(tt.kiroModel)
		if got != tt.canonicalModel {
			t.Errorf("mapKiroModelToCanonical(%q) = %q, want %q", tt.kiroModel, got, tt.canonicalModel)
		}
	}
}

func TestToolMapping(t *testing.T) {
	kiroTools := []string{"read", "write", "shell", "web_search", "grep"}
	expected := []string{"Read", "Write", "Bash", "WebSearch", "Grep"}

	got := mapKiroToolsToCanonical(kiroTools)

	if len(got) != len(expected) {
		t.Fatalf("Tool count = %d, want %d", len(got), len(expected))
	}

	for i, tool := range expected {
		if got[i] != tool {
			t.Errorf("Tool[%d] = %q, want %q", i, got[i], tool)
		}
	}
}

func TestReverseToolMapping(t *testing.T) {
	canonicalTools := []string{"Read", "Write", "Bash", "WebFetch", "Edit"}
	expected := []string{"read", "write", "shell", "web_fetch", "write"}

	got := mapCanonicalToolsToKiro(canonicalTools)

	if len(got) != len(expected) {
		t.Fatalf("Tool count = %d, want %d", len(got), len(expected))
	}

	for i, tool := range expected {
		if got[i] != tool {
			t.Errorf("Tool[%d] = %q, want %q", i, got[i], tool)
		}
	}
}

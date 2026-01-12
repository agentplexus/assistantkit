package agents

import (
	"strings"
	"testing"
)

func TestAdapterRegistry(t *testing.T) {
	names := AdapterNames()

	// Should have at least Claude adapter
	if len(names) < 1 {
		t.Errorf("expected at least 1 adapter, got %d", len(names))
	}

	// Check Claude adapter exists
	claude, ok := GetAdapter("claude")
	if !ok {
		t.Error("expected Claude adapter to be registered")
	}
	if claude.Name() != "claude" {
		t.Errorf("expected adapter name 'claude', got '%s'", claude.Name())
	}
}

func TestClaudeAdapter(t *testing.T) {
	adapter, ok := GetAdapter("claude")
	if !ok {
		t.Fatal("Claude adapter not found")
	}

	// Test marshal
	agent := NewAgent("release-coordinator", "Orchestrates software releases")
	agent.SetModel("sonnet")
	agent.AddTools("Read", "Write", "Bash")
	agent.AddSkills("version-analysis", "commit-classification")
	agent.Instructions = "You are a release coordinator agent."

	data, err := adapter.Marshal(agent)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	content := string(data)

	// Check frontmatter
	if !strings.HasPrefix(content, "---") {
		t.Error("expected Markdown to start with frontmatter")
	}
	if !strings.Contains(content, "name: release-coordinator") {
		t.Error("expected name in frontmatter")
	}
	if !strings.Contains(content, "model: sonnet") {
		t.Error("expected model in frontmatter")
	}
	if !strings.Contains(content, "tools:") {
		t.Error("expected tools in frontmatter")
	}
	if !strings.Contains(content, "skills:") {
		t.Error("expected skills in frontmatter")
	}

	// Check instructions are in body (after frontmatter)
	if !strings.Contains(content, "You are a release coordinator agent.") {
		t.Error("expected instructions in body")
	}

	// Test round-trip
	parsed, err := adapter.Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if parsed.Name != agent.Name {
		t.Errorf("round-trip: expected Name '%s', got '%s'", agent.Name, parsed.Name)
	}
	if parsed.Model != agent.Model {
		t.Errorf("round-trip: expected Model '%s', got '%s'", agent.Model, parsed.Model)
	}
	if len(parsed.Tools) != len(agent.Tools) {
		t.Errorf("round-trip: expected %d tools, got %d", len(agent.Tools), len(parsed.Tools))
	}
	if len(parsed.Skills) != len(agent.Skills) {
		t.Errorf("round-trip: expected %d skills, got %d", len(agent.Skills), len(parsed.Skills))
	}
}

func TestClaudeAdapterMinimal(t *testing.T) {
	adapter, ok := GetAdapter("claude")
	if !ok {
		t.Fatal("Claude adapter not found")
	}

	// Test with minimal agent (no tools, skills, or model)
	agent := NewAgent("simple-agent", "A simple agent")
	agent.Instructions = "Do something simple."

	data, err := adapter.Marshal(agent)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	content := string(data)

	// Should not have model, tools, or skills in frontmatter
	if strings.Contains(content, "model:") {
		t.Error("should not have model when not set")
	}
	if strings.Contains(content, "tools:") {
		t.Error("should not have tools when empty")
	}
	if strings.Contains(content, "skills:") {
		t.Error("should not have skills when empty")
	}
}

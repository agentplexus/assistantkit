package core

import "testing"

func TestNewAgent(t *testing.T) {
	agent := NewAgent("release-coordinator", "Orchestrates releases")

	if agent.Name != "release-coordinator" {
		t.Errorf("expected Name 'release-coordinator', got '%s'", agent.Name)
	}
	if agent.Description != "Orchestrates releases" {
		t.Errorf("expected Description 'Orchestrates releases', got '%s'", agent.Description)
	}
}

func TestAgentWithModel(t *testing.T) {
	agent := NewAgent("test", "test").WithModel(ModelSonnet)

	if agent.Model != ModelSonnet {
		t.Errorf("expected Model 'sonnet', got '%s'", agent.Model)
	}
}

func TestAgentWithTools(t *testing.T) {
	agent := NewAgent("test", "test").WithTools("Read", "Write")

	if len(agent.Tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(agent.Tools))
	}
}

func TestAgentWithInstructions(t *testing.T) {
	agent := NewAgent("test", "test").WithInstructions("Do the thing")

	if agent.Instructions != "Do the thing" {
		t.Errorf("expected Instructions 'Do the thing', got '%s'", agent.Instructions)
	}
}

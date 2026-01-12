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

func TestAgentSetModel(t *testing.T) {
	agent := NewAgent("test", "test")

	agent.SetModel("sonnet")

	if agent.Model != "sonnet" {
		t.Errorf("expected Model 'sonnet', got '%s'", agent.Model)
	}
}

func TestAgentAddTool(t *testing.T) {
	agent := NewAgent("test", "test")

	agent.AddTool("Read")
	agent.AddTool("Write")

	if len(agent.Tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(agent.Tools))
	}
}

func TestAgentAddTools(t *testing.T) {
	agent := NewAgent("test", "test")

	agent.AddTools("Read", "Write", "Bash", "Glob")

	if len(agent.Tools) != 4 {
		t.Errorf("expected 4 tools, got %d", len(agent.Tools))
	}
}

func TestAgentAddSkill(t *testing.T) {
	agent := NewAgent("test", "test")

	agent.AddSkill("version-analysis")

	if len(agent.Skills) != 1 {
		t.Errorf("expected 1 skill, got %d", len(agent.Skills))
	}
}

func TestAgentAddSkills(t *testing.T) {
	agent := NewAgent("test", "test")

	agent.AddSkills("version-analysis", "commit-classification")

	if len(agent.Skills) != 2 {
		t.Errorf("expected 2 skills, got %d", len(agent.Skills))
	}
}

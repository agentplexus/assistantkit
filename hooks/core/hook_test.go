package core

import "testing"

func TestNewCommandHook(t *testing.T) {
	hook := NewCommandHook("echo hello")
	if hook.Type != HookTypeCommand {
		t.Errorf("Expected type %q, got %q", HookTypeCommand, hook.Type)
	}
	if hook.Command != "echo hello" {
		t.Errorf("Expected command 'echo hello', got %q", hook.Command)
	}
}

func TestNewPromptHook(t *testing.T) {
	hook := NewPromptHook("Check if valid")
	if hook.Type != HookTypePrompt {
		t.Errorf("Expected type %q, got %q", HookTypePrompt, hook.Type)
	}
	if hook.Prompt != "Check if valid" {
		t.Errorf("Expected prompt 'Check if valid', got %q", hook.Prompt)
	}
}

func TestHookWithTimeout(t *testing.T) {
	hook := NewCommandHook("echo test").WithTimeout(30)
	if hook.Timeout != 30 {
		t.Errorf("Expected timeout 30, got %d", hook.Timeout)
	}
}

func TestHookWithShowOutput(t *testing.T) {
	hook := NewCommandHook("echo test").WithShowOutput(true)
	if !hook.ShowOutput {
		t.Error("Expected ShowOutput to be true")
	}
}

func TestHookWithWorkingDir(t *testing.T) {
	hook := NewCommandHook("echo test").WithWorkingDir("/tmp")
	if hook.WorkingDir != "/tmp" {
		t.Errorf("Expected WorkingDir '/tmp', got %q", hook.WorkingDir)
	}
}

func TestHookIsCommand(t *testing.T) {
	cmdHook := NewCommandHook("echo test")
	if !cmdHook.IsCommand() {
		t.Error("Command hook should return true for IsCommand")
	}

	promptHook := NewPromptHook("Check")
	if promptHook.IsCommand() {
		t.Error("Prompt hook should return false for IsCommand")
	}

	// Test inference from fields
	inferredHook := Hook{Command: "echo test"}
	if !inferredHook.IsCommand() {
		t.Error("Hook with command should be inferred as command type")
	}
}

func TestHookIsPrompt(t *testing.T) {
	promptHook := NewPromptHook("Check")
	if !promptHook.IsPrompt() {
		t.Error("Prompt hook should return true for IsPrompt")
	}

	cmdHook := NewCommandHook("echo test")
	if cmdHook.IsPrompt() {
		t.Error("Command hook should return false for IsPrompt")
	}
}

func TestHookValidate(t *testing.T) {
	tests := []struct {
		name      string
		hook      Hook
		wantError bool
	}{
		{
			name:      "valid command hook",
			hook:      NewCommandHook("echo test"),
			wantError: false,
		},
		{
			name:      "valid prompt hook",
			hook:      NewPromptHook("Check"),
			wantError: false,
		},
		{
			name:      "no command or prompt",
			hook:      Hook{},
			wantError: true,
		},
		{
			name:      "both command and prompt",
			hook:      Hook{Command: "echo", Prompt: "check"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.hook.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

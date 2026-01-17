package hooks

import (
	"testing"
)

func TestGetAdapter(t *testing.T) {
	adapters := []string{"claude", "cursor", "windsurf"}

	for _, name := range adapters {
		t.Run(name, func(t *testing.T) {
			adapter, ok := GetAdapter(name)
			if !ok {
				t.Errorf("Adapter %q not found", name)
				return
			}
			if adapter.Name() != name {
				t.Errorf("Adapter name mismatch: expected %q, got %q", name, adapter.Name())
			}
		})
	}
}

func TestAdapterNames(t *testing.T) {
	names := AdapterNames()
	if len(names) < 3 {
		t.Errorf("Expected at least 3 adapters, got %d", len(names))
	}
}

func TestSupportedTools(t *testing.T) {
	tools := SupportedTools()
	expected := []string{"claude", "cursor", "windsurf"}

	if len(tools) != len(expected) {
		t.Errorf("Expected %d tools, got %d", len(expected), len(tools))
	}

	for i, tool := range expected {
		if tools[i] != tool {
			t.Errorf("Tool mismatch at index %d: expected %q, got %q", i, tool, tools[i])
		}
	}
}

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	if cfg == nil {
		t.Fatal("NewConfig returned nil")
	}
	if cfg.Hooks == nil {
		t.Error("Hooks map should be initialized")
	}
}

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

func TestConfigAddHook(t *testing.T) {
	cfg := NewConfig()
	hook := NewCommandHook("echo test")

	cfg.AddHook(BeforeCommand, hook)

	entries := cfg.GetHooks(BeforeCommand)
	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}
	if len(entries[0].Hooks) != 1 {
		t.Fatalf("Expected 1 hook, got %d", len(entries[0].Hooks))
	}
	if entries[0].Hooks[0].Command != "echo test" {
		t.Errorf("Expected command 'echo test', got %q", entries[0].Hooks[0].Command)
	}
}

func TestConfigAddHookWithMatcher(t *testing.T) {
	cfg := NewConfig()
	hook := NewCommandHook("echo bash")

	cfg.AddHookWithMatcher(BeforeCommand, "Bash", hook)

	entries := cfg.GetHooks(BeforeCommand)
	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}
	if entries[0].Matcher != "Bash" {
		t.Errorf("Expected matcher 'Bash', got %q", entries[0].Matcher)
	}
}

func TestConvertClaudeToCursor(t *testing.T) {
	claudeJSON := []byte(`{
		"hooks": {
			"PreToolUse": [
				{
					"matcher": "Bash",
					"hooks": [
						{
							"type": "command",
							"command": "echo before shell"
						}
					]
				}
			]
		}
	}`)

	cursorData, err := Convert(claudeJSON, "claude", "cursor")
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if len(cursorData) == 0 {
		t.Error("Convert returned empty data")
	}

	// Parse back to verify
	cursorAdapter, _ := GetAdapter("cursor")
	cfg, err := cursorAdapter.Parse(cursorData)
	if err != nil {
		t.Fatalf("Failed to parse converted data: %v", err)
	}

	hooks := cfg.GetAllHooksForEvent(BeforeCommand)
	if len(hooks) != 1 {
		t.Fatalf("Expected 1 hook, got %d", len(hooks))
	}
	if hooks[0].Command != "echo before shell" {
		t.Errorf("Expected command 'echo before shell', got %q", hooks[0].Command)
	}
}

func TestEventCanBlock(t *testing.T) {
	blockableEvents := []Event{
		BeforeFileRead, BeforeFileWrite, BeforeCommand, BeforeMCP, BeforePrompt,
	}
	for _, event := range blockableEvents {
		if !event.CanBlock() {
			t.Errorf("Expected %q to be blockable", event)
		}
	}

	nonBlockableEvents := []Event{
		AfterFileRead, AfterFileWrite, AfterCommand, AfterMCP, OnStop,
	}
	for _, event := range nonBlockableEvents {
		if event.CanBlock() {
			t.Errorf("Expected %q to not be blockable", event)
		}
	}
}

func TestEventToolSupport(t *testing.T) {
	// BeforeCommand should be supported by all three
	support := BeforeCommand.GetToolSupport()
	if !support.Claude || !support.Cursor || !support.Windsurf {
		t.Error("BeforeCommand should be supported by Claude, Cursor, and Windsurf")
	}

	// OnSessionStart is Claude-only
	support = OnSessionStart.GetToolSupport()
	if !support.Claude || support.Cursor || support.Windsurf {
		t.Error("OnSessionStart should only be supported by Claude")
	}

	// AfterResponse is Cursor-only
	support = AfterResponse.GetToolSupport()
	if support.Claude || !support.Cursor || support.Windsurf {
		t.Error("AfterResponse should only be supported by Cursor")
	}
}

func TestAllEvents(t *testing.T) {
	events := AllEvents()
	if len(events) < 10 {
		t.Errorf("Expected at least 10 events, got %d", len(events))
	}
}

package core

import (
	"encoding/json"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	if cfg == nil {
		t.Fatal("NewConfig returned nil")
	}
	if cfg.Hooks == nil {
		t.Error("Hooks map should be initialized")
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
}

func TestConfigAddHookWithMatcher(t *testing.T) {
	cfg := NewConfig()
	hook1 := NewCommandHook("echo bash1")
	hook2 := NewCommandHook("echo bash2")

	cfg.AddHookWithMatcher(BeforeCommand, "Bash", hook1)
	cfg.AddHookWithMatcher(BeforeCommand, "Bash", hook2)

	entries := cfg.GetHooks(BeforeCommand)
	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry (same matcher), got %d", len(entries))
	}
	if len(entries[0].Hooks) != 2 {
		t.Fatalf("Expected 2 hooks, got %d", len(entries[0].Hooks))
	}
}

func TestConfigGetAllHooksForEvent(t *testing.T) {
	cfg := NewConfig()
	cfg.AddHookWithMatcher(BeforeCommand, "Bash", NewCommandHook("echo 1"))
	cfg.AddHookWithMatcher(BeforeCommand, "Write", NewCommandHook("echo 2"))

	hooks := cfg.GetAllHooksForEvent(BeforeCommand)
	if len(hooks) != 2 {
		t.Errorf("Expected 2 hooks, got %d", len(hooks))
	}
}

func TestConfigEvents(t *testing.T) {
	cfg := NewConfig()
	cfg.AddHook(BeforeCommand, NewCommandHook("echo 1"))
	cfg.AddHook(AfterCommand, NewCommandHook("echo 2"))

	events := cfg.Events()
	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}
}

func TestConfigHasHooks(t *testing.T) {
	cfg := NewConfig()
	if cfg.HasHooks() {
		t.Error("Empty config should not have hooks")
	}

	cfg.AddHook(BeforeCommand, NewCommandHook("echo test"))
	if !cfg.HasHooks() {
		t.Error("Config with hook should have hooks")
	}
}

func TestConfigHookCount(t *testing.T) {
	cfg := NewConfig()
	if cfg.HookCount() != 0 {
		t.Error("Empty config should have 0 hooks")
	}

	cfg.AddHook(BeforeCommand, NewCommandHook("echo 1"))
	cfg.AddHook(BeforeCommand, NewCommandHook("echo 2"))
	cfg.AddHook(AfterCommand, NewCommandHook("echo 3"))

	if cfg.HookCount() != 3 {
		t.Errorf("Expected 3 hooks, got %d", cfg.HookCount())
	}
}

func TestConfigRemoveHooks(t *testing.T) {
	cfg := NewConfig()
	cfg.AddHook(BeforeCommand, NewCommandHook("echo test"))
	cfg.RemoveHooks(BeforeCommand)

	if len(cfg.GetHooks(BeforeCommand)) != 0 {
		t.Error("Hooks should be removed")
	}
}

func TestConfigMerge(t *testing.T) {
	cfg1 := NewConfig()
	cfg1.AddHook(BeforeCommand, NewCommandHook("echo 1"))

	cfg2 := NewConfig()
	cfg2.AddHook(AfterCommand, NewCommandHook("echo 2"))
	cfg2.DisableAllHooks = true

	cfg1.Merge(cfg2)

	if cfg1.HookCount() != 2 {
		t.Errorf("Expected 2 hooks after merge, got %d", cfg1.HookCount())
	}
	if !cfg1.DisableAllHooks {
		t.Error("DisableAllHooks should be true after merge")
	}
}

func TestConfigFilterByTool(t *testing.T) {
	cfg := NewConfig()
	cfg.AddHook(BeforeCommand, NewCommandHook("echo cmd"))      // All tools
	cfg.AddHook(OnSessionStart, NewCommandHook("echo session")) // Claude only
	cfg.AddHook(AfterResponse, NewCommandHook("echo response")) // Cursor only

	claudeCfg := cfg.FilterByTool("claude")
	if len(claudeCfg.Hooks) != 2 {
		t.Errorf("Claude config should have 2 events, got %d", len(claudeCfg.Hooks))
	}

	cursorCfg := cfg.FilterByTool("cursor")
	if len(cursorCfg.Hooks) != 2 {
		t.Errorf("Cursor config should have 2 events, got %d", len(cursorCfg.Hooks))
	}

	windsurfCfg := cfg.FilterByTool("windsurf")
	if len(windsurfCfg.Hooks) != 1 {
		t.Errorf("Windsurf config should have 1 event, got %d", len(windsurfCfg.Hooks))
	}
}

func TestConfigJSON(t *testing.T) {
	cfg := NewConfig()
	cfg.AddHook(BeforeCommand, NewCommandHook("echo test"))

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded Config
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.HookCount() != 1 {
		t.Errorf("Expected 1 hook after round-trip, got %d", decoded.HookCount())
	}
}

func TestConfigValidate(t *testing.T) {
	// Valid config
	cfg := NewConfig()
	cfg.AddHook(BeforeCommand, NewCommandHook("echo test"))
	if err := cfg.Validate(); err != nil {
		t.Errorf("Valid config should not return error: %v", err)
	}

	// Invalid config - hook with neither command nor prompt
	invalidCfg := NewConfig()
	invalidCfg.Hooks[BeforeCommand] = []HookEntry{
		{Hooks: []Hook{{}}}, // Empty hook
	}
	err := invalidCfg.Validate()
	if err == nil {
		t.Error("Invalid config should return error")
	}
	if _, ok := err.(*HookValidationError); !ok {
		t.Errorf("Expected HookValidationError, got %T", err)
	}

	// Invalid config - hook with both command and prompt
	invalidCfg2 := NewConfig()
	invalidCfg2.Hooks[BeforeCommand] = []HookEntry{
		{Hooks: []Hook{{Command: "echo", Prompt: "check"}}},
	}
	err = invalidCfg2.Validate()
	if err == nil {
		t.Error("Config with both command and prompt should return error")
	}
}

func TestConfigWriteReadFile(t *testing.T) {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "config-test-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// Write config
	cfg := NewConfig()
	cfg.Version = 2
	cfg.DisableAllHooks = true
	cfg.AddHook(BeforeCommand, NewCommandHook("echo test"))

	if err := cfg.WriteFile(tmpPath); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Read config
	readCfg, err := ReadFile(tmpPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	// Verify
	if readCfg.Version != 2 {
		t.Errorf("Version = %d, want 2", readCfg.Version)
	}
	if !readCfg.DisableAllHooks {
		t.Error("DisableAllHooks should be true")
	}
	if readCfg.HookCount() != 1 {
		t.Errorf("HookCount() = %d, want 1", readCfg.HookCount())
	}
}

func TestConfigReadFileNotFound(t *testing.T) {
	_, err := ReadFile("/nonexistent/path/config.json")
	if err == nil {
		t.Error("ReadFile() should return error for nonexistent file")
	}
}

func TestConfigReadFileInvalidJSON(t *testing.T) {
	// Create temp file with invalid JSON
	tmpFile, err := os.CreateTemp("", "config-invalid-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmpPath := tmpFile.Name()
	if _, err := tmpFile.WriteString("{invalid json}"); err != nil {
		panic(err)
	}
	tmpFile.Close()
	defer os.Remove(tmpPath)

	_, err = ReadFile(tmpPath)
	if err == nil {
		t.Error("ReadFile() should return error for invalid JSON")
	}
}

func TestConfigMergeNil(t *testing.T) {
	cfg := NewConfig()
	cfg.AddHook(BeforeCommand, NewCommandHook("echo test"))

	// Merge nil should not panic
	cfg.Merge(nil)

	if cfg.HookCount() != 1 {
		t.Errorf("HookCount after merging nil = %d, want 1", cfg.HookCount())
	}
}

func TestConfigMergeAllowManagedHooksOnly(t *testing.T) {
	cfg1 := NewConfig()
	cfg2 := NewConfig()
	cfg2.AllowManagedHooksOnly = true

	cfg1.Merge(cfg2)

	if !cfg1.AllowManagedHooksOnly {
		t.Error("AllowManagedHooksOnly should be true after merge")
	}
}

func TestConfigAddHookWithMatcherNilHooks(t *testing.T) {
	cfg := &Config{} // Hooks is nil
	cfg.AddHookWithMatcher(BeforeCommand, "Bash", NewCommandHook("echo test"))

	if cfg.HookCount() != 1 {
		t.Errorf("HookCount() = %d, want 1", cfg.HookCount())
	}
}

func TestConfigFilterByToolUnknown(t *testing.T) {
	cfg := NewConfig()
	cfg.AddHook(BeforeCommand, NewCommandHook("echo test"))

	// Unknown tool should get no hooks (none match)
	filtered := cfg.FilterByTool("unknown")
	if filtered.HookCount() != 0 {
		t.Errorf("Unknown tool filter should have 0 hooks, got %d", filtered.HookCount())
	}
}

func TestConfigFilterPreservesSettings(t *testing.T) {
	cfg := NewConfig()
	cfg.Version = 5
	cfg.DisableAllHooks = true
	cfg.AllowManagedHooksOnly = true
	cfg.AddHook(BeforeCommand, NewCommandHook("echo test"))

	filtered := cfg.FilterByTool("claude")

	if filtered.Version != 5 {
		t.Errorf("Version not preserved: got %d", filtered.Version)
	}
	if !filtered.DisableAllHooks {
		t.Error("DisableAllHooks not preserved")
	}
	if !filtered.AllowManagedHooksOnly {
		t.Error("AllowManagedHooksOnly not preserved")
	}
}

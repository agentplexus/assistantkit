package core

import "testing"

func TestEventString(t *testing.T) {
	event := BeforeCommand
	if event.String() != "before_command" {
		t.Errorf("Expected 'before_command', got %q", event.String())
	}
}

func TestEventIsBeforeEvent(t *testing.T) {
	beforeEvents := []Event{
		BeforeFileRead, BeforeFileWrite, BeforeCommand, BeforeMCP,
		BeforePrompt, BeforeCompact, BeforeTabRead,
	}
	for _, event := range beforeEvents {
		if !event.IsBeforeEvent() {
			t.Errorf("Expected %q to be a before event", event)
		}
	}

	nonBeforeEvents := []Event{
		AfterFileRead, AfterCommand, OnStop, OnSessionStart,
	}
	for _, event := range nonBeforeEvents {
		if event.IsBeforeEvent() {
			t.Errorf("Expected %q to not be a before event", event)
		}
	}
}

func TestEventIsAfterEvent(t *testing.T) {
	afterEvents := []Event{
		AfterFileRead, AfterFileWrite, AfterCommand, AfterMCP,
		AfterResponse, AfterThought, AfterTabEdit,
	}
	for _, event := range afterEvents {
		if !event.IsAfterEvent() {
			t.Errorf("Expected %q to be an after event", event)
		}
	}

	nonAfterEvents := []Event{
		BeforeCommand, OnStop, OnSessionStart,
	}
	for _, event := range nonAfterEvents {
		if event.IsAfterEvent() {
			t.Errorf("Expected %q to not be an after event", event)
		}
	}
}

func TestEventCanBlock(t *testing.T) {
	blockableEvents := []Event{
		BeforeFileRead, BeforeFileWrite, BeforeCommand, BeforeMCP,
		BeforePrompt, BeforeCompact, BeforeTabRead, OnPermission,
	}
	for _, event := range blockableEvents {
		if !event.CanBlock() {
			t.Errorf("Expected %q to be blockable", event)
		}
	}

	nonBlockableEvents := []Event{
		AfterFileRead, AfterCommand, OnStop, OnSessionStart, AfterResponse,
	}
	for _, event := range nonBlockableEvents {
		if event.CanBlock() {
			t.Errorf("Expected %q to not be blockable", event)
		}
	}
}

func TestEventGetToolSupport(t *testing.T) {
	tests := []struct {
		event    Event
		claude   bool
		cursor   bool
		windsurf bool
	}{
		{BeforeCommand, true, true, true},
		{AfterCommand, true, true, true},
		{BeforeMCP, true, true, true},
		{OnSessionStart, true, false, false},
		{OnSessionEnd, true, false, false},
		{AfterResponse, false, true, false},
		{AfterThought, false, true, false},
		{OnPermission, true, false, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.event), func(t *testing.T) {
			support := tt.event.GetToolSupport()
			if support.Claude != tt.claude {
				t.Errorf("Claude support: expected %v, got %v", tt.claude, support.Claude)
			}
			if support.Cursor != tt.cursor {
				t.Errorf("Cursor support: expected %v, got %v", tt.cursor, support.Cursor)
			}
			if support.Windsurf != tt.windsurf {
				t.Errorf("Windsurf support: expected %v, got %v", tt.windsurf, support.Windsurf)
			}
		})
	}
}

func TestAllEvents(t *testing.T) {
	events := AllEvents()
	if len(events) < 15 {
		t.Errorf("Expected at least 15 events, got %d", len(events))
	}

	// Check that some key events are present
	eventSet := make(map[Event]bool)
	for _, e := range events {
		eventSet[e] = true
	}

	requiredEvents := []Event{
		BeforeFileRead, AfterFileWrite, BeforeCommand, AfterCommand,
		BeforeMCP, AfterMCP, BeforePrompt, OnStop,
	}
	for _, required := range requiredEvents {
		if !eventSet[required] {
			t.Errorf("Expected event %q to be in AllEvents()", required)
		}
	}
}

func TestEventGetToolSupportComprehensive(t *testing.T) {
	// Test all events to ensure GetToolSupport doesn't panic and returns valid data
	allEvents := AllEvents()
	for _, event := range allEvents {
		support := event.GetToolSupport()
		// At least one tool should support each event
		if !support.Claude && !support.Cursor && !support.Windsurf {
			t.Errorf("Event %q is not supported by any tool", event)
		}
	}
}

func TestEventGetToolSupportFileEvents(t *testing.T) {
	// Test file events with their expected support
	tests := []struct {
		event    Event
		claude   bool
		cursor   bool
		windsurf bool
	}{
		{BeforeFileRead, true, true, true},
		{AfterFileRead, true, false, true},   // Cursor doesn't support AfterFileRead
		{BeforeFileWrite, true, false, true}, // Cursor doesn't support BeforeFileWrite
		{AfterFileWrite, true, true, true},
	}
	for _, tt := range tests {
		support := tt.event.GetToolSupport()
		if support.Claude != tt.claude {
			t.Errorf("Event %q Claude support: got %v, want %v", tt.event, support.Claude, tt.claude)
		}
		if support.Cursor != tt.cursor {
			t.Errorf("Event %q Cursor support: got %v, want %v", tt.event, support.Cursor, tt.cursor)
		}
		if support.Windsurf != tt.windsurf {
			t.Errorf("Event %q Windsurf support: got %v, want %v", tt.event, support.Windsurf, tt.windsurf)
		}
	}
}

func TestEventGetToolSupportMCPEvents(t *testing.T) {
	// MCP events should be supported by all three tools
	mcpEvents := []Event{BeforeMCP, AfterMCP}
	for _, event := range mcpEvents {
		support := event.GetToolSupport()
		if !support.Claude || !support.Cursor || !support.Windsurf {
			t.Errorf("Event %q should be supported by all tools", event)
		}
	}
}

func TestEventGetToolSupportClaudeOnly(t *testing.T) {
	// Claude-only events
	claudeOnlyEvents := []Event{
		OnSessionStart, OnSessionEnd,
		OnPermission, OnNotification,
		BeforeCompact, OnSubagentStop,
	}
	for _, event := range claudeOnlyEvents {
		support := event.GetToolSupport()
		if !support.Claude {
			t.Errorf("Event %q should be supported by Claude", event)
		}
		if support.Cursor {
			t.Errorf("Event %q should NOT be supported by Cursor", event)
		}
		if support.Windsurf {
			t.Errorf("Event %q should NOT be supported by Windsurf", event)
		}
	}
}

func TestEventGetToolSupportCursorOnly(t *testing.T) {
	// Cursor-only events
	cursorOnlyEvents := []Event{
		AfterResponse, AfterThought,
		BeforeTabRead, AfterTabEdit,
	}
	for _, event := range cursorOnlyEvents {
		support := event.GetToolSupport()
		if support.Claude {
			t.Errorf("Event %q should NOT be supported by Claude", event)
		}
		if !support.Cursor {
			t.Errorf("Event %q should be supported by Cursor", event)
		}
		if support.Windsurf {
			t.Errorf("Event %q should NOT be supported by Windsurf", event)
		}
	}
}

func TestEventGetToolSupportUnknownEvent(t *testing.T) {
	// Unknown event should return empty support (default case)
	unknownEvent := Event("unknown_event")
	support := unknownEvent.GetToolSupport()
	if support.Claude || support.Cursor || support.Windsurf {
		t.Error("Unknown event should not be supported by any tool")
	}
}

func TestEventIsBeforeEventComprehensive(t *testing.T) {
	beforeEvents := []Event{
		BeforeFileRead, BeforeFileWrite, BeforeCommand,
		BeforeMCP, BeforePrompt, BeforeCompact, BeforeTabRead,
	}
	for _, event := range beforeEvents {
		if !event.IsBeforeEvent() {
			t.Errorf("Event %q should be a before event", event)
		}
	}
}

func TestEventIsAfterEventComprehensive(t *testing.T) {
	afterEvents := []Event{
		AfterFileRead, AfterFileWrite, AfterCommand,
		AfterMCP, AfterResponse, AfterThought, AfterTabEdit,
	}
	for _, event := range afterEvents {
		if !event.IsAfterEvent() {
			t.Errorf("Event %q should be an after event", event)
		}
	}
}

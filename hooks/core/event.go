// Package core provides the canonical types for hook configuration
// that can be converted to/from various AI assistant formats.
package core

// Event represents the canonical hook event types.
// Different tools use different naming conventions, but these map to common concepts.
type Event string

const (
	// File operations
	BeforeFileRead  Event = "before_file_read"
	AfterFileRead   Event = "after_file_read"
	BeforeFileWrite Event = "before_file_write"
	AfterFileWrite  Event = "after_file_write"

	// Shell/Command operations
	BeforeCommand Event = "before_command"
	AfterCommand  Event = "after_command"

	// MCP operations
	BeforeMCP Event = "before_mcp"
	AfterMCP  Event = "after_mcp"

	// Prompt/Input operations
	BeforePrompt Event = "before_prompt"

	// Agent lifecycle
	OnStop         Event = "on_stop"
	OnSessionStart Event = "on_session_start"
	OnSessionEnd   Event = "on_session_end"

	// Response events (Cursor-specific)
	AfterResponse Event = "after_response"
	AfterThought  Event = "after_thought"

	// Permission events (Claude-specific)
	OnPermission Event = "on_permission"

	// Other events
	OnNotification Event = "on_notification"
	BeforeCompact  Event = "before_compact"
	OnSubagentStop Event = "on_subagent_stop"

	// Tab/Completion events (Cursor-specific)
	BeforeTabRead Event = "before_tab_read"
	AfterTabEdit  Event = "after_tab_edit"
)

// String returns the string representation of the event.
func (e Event) String() string {
	return string(e)
}

// IsBeforeEvent returns true if this is a "before" event that can block actions.
func (e Event) IsBeforeEvent() bool {
	switch e {
	case BeforeFileRead, BeforeFileWrite, BeforeCommand, BeforeMCP,
		BeforePrompt, BeforeCompact, BeforeTabRead:
		return true
	default:
		return false
	}
}

// IsAfterEvent returns true if this is an "after" event for observation.
func (e Event) IsAfterEvent() bool {
	switch e {
	case AfterFileRead, AfterFileWrite, AfterCommand, AfterMCP,
		AfterResponse, AfterThought, AfterTabEdit:
		return true
	default:
		return false
	}
}

// CanBlock returns true if hooks for this event can block the action.
func (e Event) CanBlock() bool {
	return e.IsBeforeEvent() || e == OnPermission
}

// AllEvents returns all defined canonical events.
func AllEvents() []Event {
	return []Event{
		BeforeFileRead, AfterFileRead,
		BeforeFileWrite, AfterFileWrite,
		BeforeCommand, AfterCommand,
		BeforeMCP, AfterMCP,
		BeforePrompt,
		OnStop, OnSessionStart, OnSessionEnd,
		AfterResponse, AfterThought,
		OnPermission, OnNotification,
		BeforeCompact, OnSubagentStop,
		BeforeTabRead, AfterTabEdit,
	}
}

// ToolSupport indicates which tools support which events.
type ToolSupport struct {
	Claude   bool
	Cursor   bool
	Windsurf bool
}

// GetToolSupport returns which tools support the given event.
func (e Event) GetToolSupport() ToolSupport {
	switch e {
	case BeforeFileRead:
		return ToolSupport{Claude: true, Cursor: true, Windsurf: true}
	case AfterFileRead:
		return ToolSupport{Claude: true, Cursor: false, Windsurf: true}
	case BeforeFileWrite:
		return ToolSupport{Claude: true, Cursor: false, Windsurf: true}
	case AfterFileWrite:
		return ToolSupport{Claude: true, Cursor: true, Windsurf: true}
	case BeforeCommand:
		return ToolSupport{Claude: true, Cursor: true, Windsurf: true}
	case AfterCommand:
		return ToolSupport{Claude: true, Cursor: true, Windsurf: true}
	case BeforeMCP:
		return ToolSupport{Claude: true, Cursor: true, Windsurf: true}
	case AfterMCP:
		return ToolSupport{Claude: true, Cursor: true, Windsurf: true}
	case BeforePrompt:
		return ToolSupport{Claude: true, Cursor: true, Windsurf: true}
	case OnStop:
		return ToolSupport{Claude: true, Cursor: true, Windsurf: false}
	case OnSessionStart:
		return ToolSupport{Claude: true, Cursor: false, Windsurf: false}
	case OnSessionEnd:
		return ToolSupport{Claude: true, Cursor: false, Windsurf: false}
	case AfterResponse:
		return ToolSupport{Claude: false, Cursor: true, Windsurf: false}
	case AfterThought:
		return ToolSupport{Claude: false, Cursor: true, Windsurf: false}
	case OnPermission:
		return ToolSupport{Claude: true, Cursor: false, Windsurf: false}
	case OnNotification:
		return ToolSupport{Claude: true, Cursor: false, Windsurf: false}
	case BeforeCompact:
		return ToolSupport{Claude: true, Cursor: false, Windsurf: false}
	case OnSubagentStop:
		return ToolSupport{Claude: true, Cursor: false, Windsurf: false}
	case BeforeTabRead:
		return ToolSupport{Claude: false, Cursor: true, Windsurf: false}
	case AfterTabEdit:
		return ToolSupport{Claude: false, Cursor: true, Windsurf: false}
	default:
		return ToolSupport{}
	}
}

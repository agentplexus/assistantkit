package core

// HookType represents the type of hook execution.
type HookType string

const (
	// HookTypeCommand executes a shell command.
	HookTypeCommand HookType = "command"

	// HookTypePrompt uses an LLM for context-aware decisions (Claude-specific).
	HookTypePrompt HookType = "prompt"
)

// Hook represents a single hook definition that can be triggered by an event.
type Hook struct {
	// Type specifies how the hook is executed (command or prompt).
	Type HookType `json:"type"`

	// Command is the shell command to execute (for command type).
	Command string `json:"command,omitempty"`

	// Prompt is the LLM prompt for context-aware decisions (Claude-specific).
	Prompt string `json:"prompt,omitempty"`

	// Timeout in seconds for hook execution.
	Timeout int `json:"timeout,omitempty"`

	// ShowOutput displays hook output in the UI (Windsurf-specific).
	ShowOutput bool `json:"showOutput,omitempty"`

	// WorkingDir is the working directory for command execution.
	WorkingDir string `json:"workingDir,omitempty"`
}

// HookEntry represents a collection of hooks for a specific event,
// with optional filtering by tool/matcher.
type HookEntry struct {
	// Matcher filters which tools trigger this hook (Claude-specific).
	// Examples: "Bash", "Write", "Edit", "Read", "Bash|Write"
	Matcher string `json:"matcher,omitempty"`

	// Hooks is the list of hooks to execute for this entry.
	Hooks []Hook `json:"hooks"`
}

// NewCommandHook creates a new command-type hook.
func NewCommandHook(command string) Hook {
	return Hook{
		Type:    HookTypeCommand,
		Command: command,
	}
}

// NewPromptHook creates a new prompt-type hook (Claude-specific).
func NewPromptHook(prompt string) Hook {
	return Hook{
		Type:   HookTypePrompt,
		Prompt: prompt,
	}
}

// WithTimeout sets the timeout for a hook.
func (h Hook) WithTimeout(seconds int) Hook {
	h.Timeout = seconds
	return h
}

// WithShowOutput sets whether to show output (Windsurf-specific).
func (h Hook) WithShowOutput(show bool) Hook {
	h.ShowOutput = show
	return h
}

// WithWorkingDir sets the working directory for command execution.
func (h Hook) WithWorkingDir(dir string) Hook {
	h.WorkingDir = dir
	return h
}

// IsCommand returns true if this is a command-type hook.
func (h *Hook) IsCommand() bool {
	return h.Type == HookTypeCommand || (h.Type == "" && h.Command != "")
}

// IsPrompt returns true if this is a prompt-type hook (Claude-specific).
func (h *Hook) IsPrompt() bool {
	return h.Type == HookTypePrompt
}

// Validate checks if the hook is valid.
func (h *Hook) Validate() error {
	if h.Command == "" && h.Prompt == "" {
		return ErrNoCommandOrPrompt
	}
	if h.Command != "" && h.Prompt != "" {
		return ErrBothCommandAndPrompt
	}
	return nil
}

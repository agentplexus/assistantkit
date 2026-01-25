// Package core provides canonical power definition types.
//
// Powers are capability packages that bundle MCP servers, steering files,
// and hooks into a single installable unit for AI coding assistants.
//
// The canonical Power type maps to platform-specific formats like Kiro IDE Powers.
package core

import (
	"os"
	"path/filepath"
)

// Power represents a canonical power definition.
// Powers bundle MCP servers, steering files, and optional hooks
// into a package that can be dynamically activated based on keywords.
type Power struct {
	// Name is the unique identifier for the power (e.g., "prdtool").
	Name string `json:"name" yaml:"name"`

	// DisplayName is the user-facing title (e.g., "PRD Tool").
	DisplayName string `json:"displayName,omitempty" yaml:"displayName,omitempty"`

	// Description explains what the power provides.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Version is the semantic version (e.g., "1.0.0").
	Version string `json:"version,omitempty" yaml:"version,omitempty"`

	// Keywords trigger power activation when mentioned in conversation.
	// These should be developer-relevant terms (e.g., "database", "auth").
	Keywords []string `json:"keywords,omitempty" yaml:"keywords,omitempty"`

	// MCPServers defines the MCP server configurations for this power.
	MCPServers map[string]MCPServer `json:"mcpServers,omitempty" yaml:"mcpServers,omitempty"`

	// Onboarding contains setup instructions run when the power is first activated.
	// This can include dependency checks, configuration steps, and hook creation.
	Onboarding string `json:"onboarding,omitempty" yaml:"onboarding,omitempty"`

	// Instructions are the main steering content for the power.
	// For simple powers, this contains all guidance.
	Instructions string `json:"instructions,omitempty" yaml:"instructions,omitempty"`

	// SteeringFiles maps workflow names to steering file paths.
	// These are loaded conditionally based on the current task.
	SteeringFiles map[string]SteeringFile `json:"steeringFiles,omitempty" yaml:"steeringFiles,omitempty"`

	// Hooks defines automation triggers for IDE events.
	Hooks []Hook `json:"hooks,omitempty" yaml:"hooks,omitempty"`

	// Repository is the GitHub URL for distribution.
	Repository string `json:"repository,omitempty" yaml:"repository,omitempty"`

	// Author is the power creator.
	Author string `json:"author,omitempty" yaml:"author,omitempty"`

	// License specifies the distribution license.
	License string `json:"license,omitempty" yaml:"license,omitempty"`
}

// MCPServer defines an MCP server configuration.
type MCPServer struct {
	// Command is the executable to launch for stdio servers.
	Command string `json:"command,omitempty" yaml:"command,omitempty"`

	// Args are command-line arguments for the server.
	Args []string `json:"args,omitempty" yaml:"args,omitempty"`

	// Env contains environment variables for the server process.
	// Use ${VAR_NAME} syntax for user-provided values.
	Env map[string]string `json:"env,omitempty" yaml:"env,omitempty"`

	// URL is the endpoint for remote HTTP/SSE servers.
	URL string `json:"url,omitempty" yaml:"url,omitempty"`

	// Description explains what tools this server provides.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// SteeringFile represents a conditional steering file.
type SteeringFile struct {
	// Path is the relative path to the steering file.
	Path string `json:"path" yaml:"path"`

	// Keywords trigger loading of this specific steering file.
	Keywords []string `json:"keywords,omitempty" yaml:"keywords,omitempty"`

	// Description explains when this steering file is used.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Content is the inline steering content (alternative to Path).
	Content string `json:"content,omitempty" yaml:"content,omitempty"`
}

// Hook defines an automation trigger.
type Hook struct {
	// Name is the hook identifier.
	Name string `json:"name" yaml:"name"`

	// Event is the trigger event (e.g., "pre-commit", "on-save").
	Event string `json:"event" yaml:"event"`

	// Command is the action to execute.
	Command string `json:"command,omitempty" yaml:"command,omitempty"`

	// Prompt is the AI prompt to execute.
	Prompt string `json:"prompt,omitempty" yaml:"prompt,omitempty"`

	// Condition determines when the hook runs.
	Condition string `json:"condition,omitempty" yaml:"condition,omitempty"`
}

// NewPower creates a new Power with the given name and description.
func NewPower(name, description string) *Power {
	return &Power{
		Name:        name,
		Description: description,
		MCPServers:  make(map[string]MCPServer),
	}
}

// AddKeyword adds an activation keyword to the power.
func (p *Power) AddKeyword(keyword string) *Power {
	p.Keywords = append(p.Keywords, keyword)
	return p
}

// AddKeywords adds multiple activation keywords to the power.
func (p *Power) AddKeywords(keywords ...string) *Power {
	p.Keywords = append(p.Keywords, keywords...)
	return p
}

// AddMCPServer adds an MCP server configuration.
func (p *Power) AddMCPServer(name string, server MCPServer) *Power {
	if p.MCPServers == nil {
		p.MCPServers = make(map[string]MCPServer)
	}
	p.MCPServers[name] = server
	return p
}

// AddSteeringFile adds a conditional steering file.
func (p *Power) AddSteeringFile(name string, sf SteeringFile) *Power {
	if p.SteeringFiles == nil {
		p.SteeringFiles = make(map[string]SteeringFile)
	}
	p.SteeringFiles[name] = sf
	return p
}

// AddHook adds an automation hook.
func (p *Power) AddHook(hook Hook) *Power {
	p.Hooks = append(p.Hooks, hook)
	return p
}

// WithVersion sets the power version.
func (p *Power) WithVersion(version string) *Power {
	p.Version = version
	return p
}

// WithDisplayName sets the display name.
func (p *Power) WithDisplayName(displayName string) *Power {
	p.DisplayName = displayName
	return p
}

// WithOnboarding sets the onboarding instructions.
func (p *Power) WithOnboarding(onboarding string) *Power {
	p.Onboarding = onboarding
	return p
}

// WithInstructions sets the main steering instructions.
func (p *Power) WithInstructions(instructions string) *Power {
	p.Instructions = instructions
	return p
}

// Validate checks if the power has required fields.
func (p *Power) Validate() error {
	if p.Name == "" {
		return &ValidationError{Field: "name", Message: "name is required"}
	}
	if len(p.Keywords) == 0 {
		return &ValidationError{Field: "keywords", Message: "at least one keyword is required"}
	}
	return nil
}

// ValidationError represents a power validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// DefaultFileMode is the default permission for created files.
const DefaultFileMode = 0600

// DefaultDirMode is the default permission for created directories.
const DefaultDirMode = 0755

// WriteTo writes the power to a directory in the platform-specific format.
// This is a convenience method that uses the default Kiro format.
func (p *Power) WriteTo(dir string) error {
	if err := os.MkdirAll(dir, DefaultDirMode); err != nil {
		return err
	}

	// Create steering directory if needed
	if len(p.SteeringFiles) > 0 {
		steeringDir := filepath.Join(dir, "steering")
		if err := os.MkdirAll(steeringDir, DefaultDirMode); err != nil {
			return err
		}
	}

	return nil
}

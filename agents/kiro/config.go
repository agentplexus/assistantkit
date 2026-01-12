// Package kiro provides the AWS Kiro CLI agent adapter.
package kiro

// AgentConfig represents a Kiro CLI custom agent configuration.
// File location: ~/.kiro/agents/[agent-name].json
type AgentConfig struct {
	// Name is the agent identifier.
	Name string `json:"name"`

	// Description is a human-readable description of the agent's purpose.
	Description string `json:"description,omitempty"`

	// Tools lists the tools available to this agent.
	// Built-in tools: read, write, shell, web_search, web_fetch, etc.
	Tools []string `json:"tools,omitempty"`

	// AllowedTools lists tools that can execute without user confirmation.
	AllowedTools []string `json:"allowedTools,omitempty"`

	// Resources lists file paths or glob patterns for context.
	// Uses file:// prefix, e.g., "file://README.md", "file://.kiro/steering/**/*.md"
	Resources []string `json:"resources,omitempty"`

	// Prompt contains the system instructions for the agent.
	Prompt string `json:"prompt,omitempty"`

	// Model specifies the Claude model to use (e.g., "claude-sonnet-4").
	Model string `json:"model,omitempty"`

	// MCPServers defines MCP server configurations for this agent.
	MCPServers map[string]MCPServerConfig `json:"mcpServers,omitempty"`

	// IncludeMcpJson determines whether to inherit servers from workspace/user config.
	IncludeMcpJson bool `json:"includeMcpJson,omitempty"`
}

// MCPServerConfig represents an MCP server configuration within an agent.
type MCPServerConfig struct {
	// Command is the executable to launch for stdio servers.
	Command string `json:"command,omitempty"`

	// Args are command-line arguments for the server.
	Args []string `json:"args,omitempty"`

	// Env contains environment variables for the server process.
	Env map[string]string `json:"env,omitempty"`

	// URL is the endpoint for remote HTTP/SSE servers.
	URL string `json:"url,omitempty"`

	// Headers contains HTTP headers for authentication.
	Headers map[string]string `json:"headers,omitempty"`
}

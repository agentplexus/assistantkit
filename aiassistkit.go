// Package aiassistkit provides a unified interface for managing configuration files
// across multiple AI coding assistants including Claude Code, Cursor, Windsurf,
// VS Code, OpenAI Codex CLI, Cline, and Roo Code.
//
// AI Assist Kit supports multiple configuration types:
//
//   - MCP (Model Context Protocol) server configurations
//   - Settings (permissions, sandbox, general settings) - coming soon
//   - Rules (team rules, coding guidelines) - coming soon
//   - Memory (CLAUDE.md, .cursorrules, etc.) - coming soon
//
// # MCP Configuration
//
// The mcp subpackage provides adapters for reading, writing, and converting
// MCP server configurations between different AI assistant formats.
//
// Example usage:
//
//	import (
//	    "github.com/grokify/aiassistkit/mcp"
//	    "github.com/grokify/aiassistkit/mcp/claude"
//	    "github.com/grokify/aiassistkit/mcp/vscode"
//	)
//
//	// Read Claude config and write to VS Code format
//	cfg, _ := claude.ReadProjectConfig()
//	vscode.WriteWorkspaceConfig(cfg)
//
//	// Or use dynamic conversion
//	data, _ := mcp.Convert(jsonData, "claude", "vscode")
//
// # Related Projects
//
// AI Assist Kit is part of the AgentPlexus family of Go modules:
//   - AI Assist Kit - AI coding assistant configuration management
//   - OmniVault - Unified secrets management
//   - OmniLLM - Multi-provider LLM abstraction
//   - OmniSerp - Search engine abstraction
//   - OmniObserve - LLM observability abstraction
package aiassistkit

// Version is the current version of AI Assist Kit.
const Version = "0.1.0"

// ConfigType represents the type of configuration.
type ConfigType string

const (
	// ConfigTypeMCP represents MCP server configuration.
	ConfigTypeMCP ConfigType = "mcp"

	// ConfigTypeSettings represents general settings configuration.
	ConfigTypeSettings ConfigType = "settings"

	// ConfigTypeRules represents team rules configuration.
	ConfigTypeRules ConfigType = "rules"

	// ConfigTypeMemory represents memory/context configuration.
	ConfigTypeMemory ConfigType = "memory"
)

// SupportedConfigTypes returns a list of configuration types that AI Assist Kit supports.
func SupportedConfigTypes() []ConfigType {
	return []ConfigType{
		ConfigTypeMCP,
		ConfigTypeSettings,
		ConfigTypeRules,
		ConfigTypeMemory,
	}
}

// SupportedTools returns a list of AI coding tools that AI Assist Kit supports.
func SupportedTools() []string {
	return []string{
		"claude",   // Claude Code / Claude Desktop
		"cursor",   // Cursor IDE
		"windsurf", // Windsurf (Codeium)
		"vscode",   // VS Code / GitHub Copilot
		"codex",    // OpenAI Codex CLI
		"cline",    // Cline VS Code extension
		"roo",      // Roo Code VS Code extension
		"kiro",     // AWS Kiro CLI
	}
}

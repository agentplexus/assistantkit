// Package agents provides adapters for AI assistant agent definitions.
//
// Supported tools:
//   - Claude Code: agents/<name>.md (Markdown with YAML frontmatter)
//   - AWS Kiro CLI: ~/.kiro/agents/<name>.json (JSON format)
//
// Example usage:
//
//	package main
//
//	import (
//	    "github.com/grokify/aiassistkit/agents"
//	)
//
//	func main() {
//	    // Create a new agent
//	    agent := agents.NewAgent("release-coordinator", "Orchestrates software releases")
//	    agent.SetModel("sonnet")
//	    agent.AddTools("Read", "Write", "Bash", "Glob", "Grep")
//	    agent.AddSkills("version-analysis", "commit-classification")
//	    agent.Instructions = "You are a release coordinator agent..."
//
//	    // Write to Claude format
//	    claudeAdapter, _ := agents.GetAdapter("claude")
//	    claudeAdapter.WriteFile(agent, "./agents/release-coordinator.md")
//
//	    // Write to Kiro format
//	    kiroAdapter, _ := agents.GetAdapter("kiro")
//	    kiroAdapter.WriteFile(agent, "~/.kiro/agents/release-coordinator.json")
//	}
package agents

import (
	"github.com/grokify/aiassistkit/agents/core"

	// Import adapters for side-effect registration
	_ "github.com/grokify/aiassistkit/agents/claude"
	_ "github.com/grokify/aiassistkit/agents/codex"
	_ "github.com/grokify/aiassistkit/agents/gemini"
	_ "github.com/grokify/aiassistkit/agents/kiro"
)

// Re-export core types for convenience
type (
	Agent   = core.Agent
	Adapter = core.Adapter
)

// Re-export core functions
var (
	NewAgent             = core.NewAgent
	GetAdapter           = core.GetAdapter
	AdapterNames         = core.AdapterNames
	ReadCanonicalFile    = core.ReadCanonicalFile
	WriteCanonicalFile   = core.WriteCanonicalFile
	WriteCanonicalJSON   = core.WriteCanonicalJSON
	ReadCanonicalDir     = core.ReadCanonicalDir
	WriteAgentsToDir     = core.WriteAgentsToDir
	ParseMarkdownAgent   = core.ParseMarkdownAgent
	MarshalMarkdownAgent = core.MarshalMarkdownAgent
)

// Re-export error types
type (
	ParseError   = core.ParseError
	MarshalError = core.MarshalError
	ReadError    = core.ReadError
	WriteError   = core.WriteError
)

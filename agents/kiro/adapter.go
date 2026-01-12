package kiro

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/grokify/aiassistkit/agents/core"
)

const (
	// AdapterName is the identifier for this adapter.
	AdapterName = "kiro"

	// AgentsDir is the agents directory name.
	AgentsDir = "agents"

	// ProjectConfigDir is the project config directory.
	ProjectConfigDir = ".kiro"
)

func init() {
	core.Register(&Adapter{})
}

// Adapter converts between canonical Agent and Kiro CLI agent format.
type Adapter struct{}

// Name returns the adapter identifier.
func (a *Adapter) Name() string {
	return AdapterName
}

// FileExtension returns the file extension for Kiro agents.
func (a *Adapter) FileExtension() string {
	return ".json"
}

// DefaultDir returns the default directory name for Kiro agents.
func (a *Adapter) DefaultDir() string {
	return AgentsDir
}

// Parse converts Kiro agent JSON bytes to canonical Agent.
func (a *Adapter) Parse(data []byte) (*core.Agent, error) {
	var kiroCfg AgentConfig
	if err := json.Unmarshal(data, &kiroCfg); err != nil {
		return nil, &core.ParseError{Format: AdapterName, Err: err}
	}

	return a.ToCore(&kiroCfg), nil
}

// Marshal converts canonical Agent to Kiro agent JSON bytes.
func (a *Adapter) Marshal(agent *core.Agent) ([]byte, error) {
	kiroCfg := a.FromCore(agent)
	return json.MarshalIndent(kiroCfg, "", "  ")
}

// ReadFile reads a Kiro agent JSON file and returns canonical Agent.
func (a *Adapter) ReadFile(path string) (*core.Agent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &core.ReadError{Path: path, Err: err}
	}

	agent, err := a.Parse(data)
	if err != nil {
		if pe, ok := err.(*core.ParseError); ok {
			pe.Path = path
		}
		return nil, err
	}

	// Infer name from filename if not set
	if agent.Name == "" {
		base := filepath.Base(path)
		agent.Name = strings.TrimSuffix(base, filepath.Ext(base))
	}

	return agent, nil
}

// WriteFile writes canonical Agent to a Kiro agent JSON file.
func (a *Adapter) WriteFile(agent *core.Agent, path string) error {
	data, err := a.Marshal(agent)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, core.DefaultDirMode); err != nil {
		return &core.WriteError{Path: path, Err: err}
	}

	if err := os.WriteFile(path, data, core.DefaultFileMode); err != nil {
		return &core.WriteError{Path: path, Err: err}
	}

	return nil
}

// ToCore converts Kiro agent config to canonical Agent.
func (a *Adapter) ToCore(kiroCfg *AgentConfig) *core.Agent {
	agent := &core.Agent{
		Name:         kiroCfg.Name,
		Description:  kiroCfg.Description,
		Instructions: kiroCfg.Prompt,
	}

	// Map Kiro model names to canonical model names
	if kiroCfg.Model != "" {
		agent.Model = mapKiroModelToCanonical(kiroCfg.Model)
	}

	// Map Kiro tools to canonical tools
	if len(kiroCfg.Tools) > 0 {
		agent.Tools = mapKiroToolsToCanonical(kiroCfg.Tools)
	}

	// Store allowed tools in metadata (Kiro-specific)
	// Note: AllowedTools is a Kiro-specific concept for auto-approval

	// Store resources as skills (closest mapping)
	// Note: Resources in Kiro load context files, similar to skill dependencies

	return agent
}

// FromCore converts canonical Agent to Kiro agent config.
func (a *Adapter) FromCore(agent *core.Agent) *AgentConfig {
	kiroCfg := &AgentConfig{
		Name:        agent.Name,
		Description: agent.Description,
		Prompt:      agent.Instructions,
	}

	// Map canonical model to Kiro model name
	if agent.Model != "" {
		kiroCfg.Model = mapCanonicalModelToKiro(agent.Model)
	}

	// Map canonical tools to Kiro tools
	if len(agent.Tools) > 0 {
		kiroCfg.Tools = mapCanonicalToolsToKiro(agent.Tools)
	}

	// Map skills to resources (steering files)
	if len(agent.Skills) > 0 {
		kiroCfg.Resources = mapSkillsToResources(agent.Skills)
	}

	return kiroCfg
}

// mapKiroModelToCanonical maps Kiro model names to canonical names.
func mapKiroModelToCanonical(kiroModel string) string {
	switch kiroModel {
	case "claude-sonnet-4", "claude-4-sonnet":
		return "sonnet"
	case "claude-opus-4", "claude-4-opus":
		return "opus"
	case "claude-haiku", "claude-3-haiku":
		return "haiku"
	default:
		// Return as-is if not recognized
		return kiroModel
	}
}

// mapCanonicalModelToKiro maps canonical model names to Kiro names.
func mapCanonicalModelToKiro(model string) string {
	switch strings.ToLower(model) {
	case "sonnet":
		return "claude-sonnet-4"
	case "opus":
		return "claude-opus-4"
	case "haiku":
		return "claude-haiku"
	default:
		return model
	}
}

// mapKiroToolsToCanonical maps Kiro tool names to canonical names.
func mapKiroToolsToCanonical(kiroTools []string) []string {
	toolMap := map[string]string{
		"read":       "Read",
		"write":      "Write",
		"shell":      "Bash",
		"web_search": "WebSearch",
		"web_fetch":  "WebFetch",
		"grep":       "Grep",
		"glob":       "Glob",
	}

	var canonical []string
	for _, tool := range kiroTools {
		if mapped, ok := toolMap[tool]; ok {
			canonical = append(canonical, mapped)
		} else {
			// Capitalize first letter for unknown tools
			if len(tool) > 0 {
				canonical = append(canonical, strings.ToUpper(tool[:1])+tool[1:])
			}
		}
	}
	return canonical
}

// mapCanonicalToolsToKiro maps canonical tool names to Kiro names.
func mapCanonicalToolsToKiro(tools []string) []string {
	toolMap := map[string]string{
		"Read":      "read",
		"Write":     "write",
		"Bash":      "shell",
		"WebSearch": "web_search",
		"WebFetch":  "web_fetch",
		"Grep":      "grep",
		"Glob":      "glob",
		"Edit":      "write", // Edit maps to write in Kiro
	}

	var kiroTools []string
	for _, tool := range tools {
		if mapped, ok := toolMap[tool]; ok {
			kiroTools = append(kiroTools, mapped)
		} else {
			// Lowercase for unknown tools
			kiroTools = append(kiroTools, strings.ToLower(tool))
		}
	}
	return kiroTools
}

// mapSkillsToResources converts skill names to Kiro resource paths.
func mapSkillsToResources(skills []string) []string {
	var resources []string
	for _, skill := range skills {
		// Map skills to steering files
		resources = append(resources, "file://.kiro/steering/"+skill+".md")
	}
	return resources
}

// UserAgentsPath returns the path to the user's agents directory.
func UserAgentsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ProjectConfigDir, AgentsDir), nil
}

// UserAgentPath returns the path to a specific user agent config.
func UserAgentPath(name string) (string, error) {
	dir, err := UserAgentsPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name+".json"), nil
}

// ReadUserAgent reads a user-level agent configuration.
func ReadUserAgent(name string) (*core.Agent, error) {
	path, err := UserAgentPath(name)
	if err != nil {
		return nil, err
	}
	adapter := &Adapter{}
	return adapter.ReadFile(path)
}

// WriteUserAgent writes an agent to the user's agents directory.
func WriteUserAgent(agent *core.Agent) error {
	path, err := UserAgentPath(agent.Name)
	if err != nil {
		return err
	}
	adapter := &Adapter{}
	return adapter.WriteFile(agent, path)
}

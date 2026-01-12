// Package core provides the canonical agent definition types.
// Agent definitions use Description/Instructions as the canonical form,
// which maps losslessly to Claude Code, Kiro CLI, and OpenAI Codex.
package core

// Agent represents a canonical agent/subagent definition.
// This structure maps directly to Claude Code, Kiro CLI, and Codex agents.
type Agent struct {
	// Name is the unique identifier for the agent (e.g., "release-coordinator").
	Name string `json:"name"`

	// Description is a brief summary of what the agent does and when to use it.
	// Example: "Orchestrates software releases including versioning and tagging."
	Description string `json:"description,omitempty"`

	// Instructions are the detailed system prompt for the agent.
	// This is the full guidance on how the agent should behave.
	Instructions string `json:"instructions,omitempty"`

	// Model is the preferred AI model (e.g., "haiku", "sonnet", "opus").
	Model string `json:"model,omitempty"`

	// Tools are the tools available to the agent (e.g., "Read", "Write", "Bash").
	Tools []string `json:"tools,omitempty"`

	// Skills are capabilities or skills the agent can invoke.
	Skills []string `json:"skills,omitempty"`

	// Dependencies are external CLI tools required by this agent.
	Dependencies []string `json:"dependencies,omitempty"`
}

// NewAgent creates a new Agent with the given name and description.
func NewAgent(name, description string) *Agent {
	return &Agent{
		Name:        name,
		Description: description,
	}
}

// SetModel sets the model for the agent.
func (a *Agent) SetModel(model string) {
	a.Model = model
}

// AddTool adds a tool to the agent's allowed tools list.
func (a *Agent) AddTool(tool string) {
	a.Tools = append(a.Tools, tool)
}

// AddTools adds multiple tools to the agent's allowed tools list.
func (a *Agent) AddTools(tools ...string) {
	a.Tools = append(a.Tools, tools...)
}

// AddSkill adds a skill to the agent.
func (a *Agent) AddSkill(skill string) {
	a.Skills = append(a.Skills, skill)
}

// AddSkills adds multiple skills to the agent.
func (a *Agent) AddSkills(skills ...string) {
	a.Skills = append(a.Skills, skills...)
}

// AddDependency adds a dependency to the agent.
func (a *Agent) AddDependency(dep string) {
	a.Dependencies = append(a.Dependencies, dep)
}

// WithInstructions sets the agent's instructions and returns the agent for chaining.
func (a *Agent) WithInstructions(instructions string) *Agent {
	a.Instructions = instructions
	return a
}

// WithTools sets the agent's tools and returns the agent for chaining.
func (a *Agent) WithTools(tools ...string) *Agent {
	a.Tools = tools
	return a
}

// WithModel sets the agent's preferred model and returns the agent for chaining.
func (a *Agent) WithModel(model string) *Agent {
	a.Model = model
	return a
}

// WithSkills sets the agent's skills and returns the agent for chaining.
func (a *Agent) WithSkills(skills ...string) *Agent {
	a.Skills = skills
	return a
}

// Package core provides canonical types for validation area definitions.
// Validation areas represent departments or areas of responsibility in the
// release process (e.g., QA, Documentation, Release Management, Security).
package core

// ValidationArea represents a canonical validation area definition.
// Each area can be converted to tool-specific formats:
//   - Claude Code: Sub-agents (agents/*.md)
//   - Gemini CLI: Commands or prompts
//   - Codex: Prompts
type ValidationArea struct {
	// Metadata
	Name        string `json:"name"`        // Area identifier (e.g., "qa", "documentation")
	Description string `json:"description"` // Brief description of the area's responsibility

	// Sign-off criteria
	SignOffCriteria string `json:"sign_off_criteria"` // What must pass for GO status

	// Checks to perform
	Checks []Check `json:"checks"` // Individual validation checks

	// Dependencies
	Dependencies []string `json:"dependencies,omitempty"` // Required CLI tools

	// Instructions for AI agents
	Instructions string `json:"instructions"` // Full instructions/system prompt

	// Claude-specific (used when generating agent)
	Model  string   `json:"model,omitempty"`  // Model for agent (sonnet, opus, haiku)
	Tools  []string `json:"tools,omitempty"`  // Allowed tools for agent
	Skills []string `json:"skills,omitempty"` // Skills to load for agent
}

// Check represents an individual validation check within an area.
type Check struct {
	Name        string `json:"name"`                   // Check identifier
	Description string `json:"description,omitempty"`  // What this check validates
	Command     string `json:"command,omitempty"`      // CLI command to execute
	Pattern     string `json:"pattern,omitempty"`      // Regex pattern to search for (failure if found)
	FilePattern string `json:"file_pattern,omitempty"` // Glob pattern for files to check
	Required    bool   `json:"required"`               // If true, failure blocks release (NO-GO)
}

// CheckStatus represents the result of a check.
type CheckStatus string

const (
	StatusGo   CheckStatus = "GO"
	StatusNoGo CheckStatus = "NO-GO"
	StatusWarn CheckStatus = "WARN"
	StatusSkip CheckStatus = "SKIP"
)

// NewValidationArea creates a new ValidationArea with the given name and description.
func NewValidationArea(name, description string) *ValidationArea {
	return &ValidationArea{
		Name:        name,
		Description: description,
	}
}

// AddCheck adds a check to the validation area.
func (v *ValidationArea) AddCheck(check Check) {
	v.Checks = append(v.Checks, check)
}

// AddDependency adds a CLI tool dependency.
func (v *ValidationArea) AddDependency(dep string) {
	v.Dependencies = append(v.Dependencies, dep)
}

// SetModel sets the model for the area (used in Claude agent generation).
func (v *ValidationArea) SetModel(model string) {
	v.Model = model
}

// AddTool adds a tool to the allowed tools list (used in Claude agent generation).
func (v *ValidationArea) AddTool(tool string) {
	v.Tools = append(v.Tools, tool)
}

// AddTools adds multiple tools to the allowed tools list.
func (v *ValidationArea) AddTools(tools ...string) {
	v.Tools = append(v.Tools, tools...)
}

// AddSkill adds a skill (used in Claude agent generation).
func (v *ValidationArea) AddSkill(skill string) {
	v.Skills = append(v.Skills, skill)
}

// Predefined validation areas for software releases.
var (
	// AreaQA is the Quality Assurance validation area.
	AreaQA = "qa"

	// AreaDocumentation is the Documentation validation area.
	AreaDocumentation = "documentation"

	// AreaRelease is the Release Management validation area.
	AreaRelease = "release"

	// AreaSecurity is the Security/Compliance validation area.
	AreaSecurity = "security"
)

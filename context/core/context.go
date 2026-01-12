// Package core provides the canonical types for project context that can be
// converted to various AI assistant formats (CLAUDE.md, .cursorrules, etc.).
package core

import (
	"encoding/json"
	"os"
)

// Context represents the canonical project context that can be
// converted to various AI assistant formats.
type Context struct {
	// Schema is the JSON Schema reference for validation.
	Schema string `json:"$schema,omitempty"`

	// Name is the project name.
	Name string `json:"name"`

	// Description is a brief description of the project's purpose.
	Description string `json:"description,omitempty"`

	// Version is the current project version.
	Version string `json:"version,omitempty"`

	// Language is the primary programming language.
	Language string `json:"language,omitempty"`

	// Architecture describes the high-level architecture.
	Architecture *Architecture `json:"architecture,omitempty"`

	// Packages lists key packages/modules and their purposes.
	Packages []Package `json:"packages,omitempty"`

	// Commands contains common commands for working with the project.
	Commands map[string]string `json:"commands,omitempty"`

	// Conventions lists coding conventions and patterns.
	Conventions []string `json:"conventions,omitempty"`

	// Dependencies describes key dependencies.
	Dependencies *Dependencies `json:"dependencies,omitempty"`

	// Testing describes the testing strategy.
	Testing *Testing `json:"testing,omitempty"`

	// Files describes important files.
	Files *Files `json:"files,omitempty"`

	// Notes contains additional notes and gotchas.
	Notes []Note `json:"notes,omitempty"`

	// Related lists related projects or resources.
	Related []Related `json:"related,omitempty"`
}

// Architecture describes the high-level architecture of the project.
type Architecture struct {
	// Pattern is the primary architectural pattern (e.g., "adapter", "hexagonal").
	Pattern string `json:"pattern,omitempty"`

	// Summary is a brief summary of the architecture.
	Summary string `json:"summary,omitempty"`

	// Diagrams contains ASCII or mermaid diagrams.
	Diagrams []Diagram `json:"diagrams,omitempty"`
}

// Diagram represents an architecture diagram.
type Diagram struct {
	// Title is the diagram title.
	Title string `json:"title,omitempty"`

	// Type is the diagram type ("ascii" or "mermaid").
	Type string `json:"type,omitempty"`

	// Content is the diagram content.
	Content string `json:"content"`
}

// Package describes a package or module in the project.
type Package struct {
	// Path is the relative path to the package.
	Path string `json:"path"`

	// Purpose is a brief description of the package's purpose.
	Purpose string `json:"purpose"`

	// Public indicates whether this is part of the public API.
	Public *bool `json:"public,omitempty"`
}

// IsPublic returns whether the package is public (defaults to true).
func (p *Package) IsPublic() bool {
	if p.Public == nil {
		return true
	}
	return *p.Public
}

// Dependencies describes the project's dependencies.
type Dependencies struct {
	// Runtime lists runtime dependencies.
	Runtime []Dependency `json:"runtime,omitempty"`

	// Development lists development dependencies.
	Development []Dependency `json:"development,omitempty"`
}

// Dependency represents a single dependency.
type Dependency struct {
	// Name is the dependency name/identifier.
	Name string `json:"name"`

	// Purpose describes what this dependency is used for.
	Purpose string `json:"purpose,omitempty"`
}

// Testing describes the testing strategy.
type Testing struct {
	// Framework is the testing framework used.
	Framework string `json:"framework,omitempty"`

	// Coverage describes coverage requirements or current coverage.
	Coverage string `json:"coverage,omitempty"`

	// Patterns lists testing patterns and conventions.
	Patterns []string `json:"patterns,omitempty"`
}

// Files describes important files in the project.
type Files struct {
	// EntryPoints lists main entry point files.
	EntryPoints []string `json:"entryPoints,omitempty"`

	// Config lists configuration files.
	Config []string `json:"config,omitempty"`

	// Ignore lists files/patterns to ignore during analysis.
	Ignore []string `json:"ignore,omitempty"`
}

// Note represents an additional note or gotcha.
type Note struct {
	// Title is an optional title for the note.
	Title string `json:"title,omitempty"`

	// Content is the note content.
	Content string `json:"content"`

	// Severity indicates the importance (info, warning, critical).
	Severity string `json:"severity,omitempty"`
}

// GetSeverity returns the severity, defaulting to "info".
func (n *Note) GetSeverity() string {
	if n.Severity == "" {
		return "info"
	}
	return n.Severity
}

// Related represents a related project or resource.
type Related struct {
	// Name is the name of the related item.
	Name string `json:"name"`

	// URL is an optional URL.
	URL string `json:"url,omitempty"`

	// Description is an optional description.
	Description string `json:"description,omitempty"`
}

// NewContext creates a new empty Context.
func NewContext(name string) *Context {
	return &Context{
		Name:     name,
		Commands: make(map[string]string),
	}
}

// ReadFile reads a Context from a JSON file.
func ReadFile(path string) (*Context, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &ParseError{Path: path, Err: err}
	}
	return Parse(data)
}

// Parse parses JSON data into a Context.
func Parse(data []byte) (*Context, error) {
	var ctx Context
	if err := json.Unmarshal(data, &ctx); err != nil {
		return nil, &ParseError{Err: err}
	}
	return &ctx, nil
}

// Marshal converts the Context to JSON.
func (c *Context) Marshal() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

// WriteFile writes the Context to a JSON file using DefaultFileMode.
func (c *Context) WriteFile(path string) error {
	return c.WriteFileWithMode(path, DefaultFileMode)
}

// WriteFileWithMode writes the Context to a JSON file with the specified permission mode.
func (c *Context) WriteFileWithMode(path string, mode os.FileMode) error {
	data, err := c.Marshal()
	if err != nil {
		return &WriteError{Path: path, Err: err}
	}
	return os.WriteFile(path, data, mode)
}

// AddPackage adds a package to the context.
func (c *Context) AddPackage(path, purpose string) {
	c.Packages = append(c.Packages, Package{Path: path, Purpose: purpose})
}

// AddConvention adds a convention to the context.
func (c *Context) AddConvention(convention string) {
	c.Conventions = append(c.Conventions, convention)
}

// AddNote adds a note to the context.
func (c *Context) AddNote(content string) {
	c.Notes = append(c.Notes, Note{Content: content})
}

// AddNoteWithSeverity adds a note with a specific severity.
func (c *Context) AddNoteWithSeverity(title, content, severity string) {
	c.Notes = append(c.Notes, Note{Title: title, Content: content, Severity: severity})
}

// SetCommand sets a command.
func (c *Context) SetCommand(name, command string) {
	if c.Commands == nil {
		c.Commands = make(map[string]string)
	}
	c.Commands[name] = command
}

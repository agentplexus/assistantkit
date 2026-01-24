package core

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// DefaultFileMode is the default permission for generated files.
const DefaultFileMode fs.FileMode = 0600

// DefaultDirMode is the default permission for generated directories.
const DefaultDirMode fs.FileMode = 0700

// Adapter converts between canonical Skill and tool-specific formats.
type Adapter interface {
	// Name returns the adapter identifier (e.g., "claude", "codex").
	Name() string

	// SkillFileName returns the skill definition filename (e.g., "SKILL.md").
	SkillFileName() string

	// DefaultDir returns the default directory name for skills.
	DefaultDir() string

	// Parse converts tool-specific bytes to canonical Skill.
	Parse(data []byte) (*Skill, error)

	// Marshal converts canonical Skill to tool-specific bytes.
	Marshal(skill *Skill) ([]byte, error)

	// ReadFile reads from path and returns canonical Skill.
	ReadFile(path string) (*Skill, error)

	// WriteFile writes canonical Skill to path.
	WriteFile(skill *Skill, path string) error

	// WriteSkillDir writes the complete skill directory structure.
	WriteSkillDir(skill *Skill, baseDir string) error
}

// Registry manages adapter registration and lookup.
type Registry struct {
	mu       sync.RWMutex
	adapters map[string]Adapter
}

// NewRegistry creates a new adapter registry.
func NewRegistry() *Registry {
	return &Registry{
		adapters: make(map[string]Adapter),
	}
}

// Register adds an adapter to the registry.
func (r *Registry) Register(adapter Adapter) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adapters[adapter.Name()] = adapter
}

// GetAdapter returns an adapter by name.
func (r *Registry) GetAdapter(name string) (Adapter, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	adapter, ok := r.adapters[name]
	return adapter, ok
}

// AdapterNames returns all registered adapter names sorted alphabetically.
func (r *Registry) AdapterNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Convert converts skill data from one format to another.
func (r *Registry) Convert(data []byte, from, to string) ([]byte, error) {
	fromAdapter, ok := r.GetAdapter(from)
	if !ok {
		return nil, fmt.Errorf("unknown source adapter: %s", from)
	}

	toAdapter, ok := r.GetAdapter(to)
	if !ok {
		return nil, fmt.Errorf("unknown target adapter: %s", to)
	}

	skill, err := fromAdapter.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", from, err)
	}

	return toAdapter.Marshal(skill)
}

// DefaultRegistry is the global adapter registry.
var DefaultRegistry = NewRegistry()

// Register adds an adapter to the default registry.
func Register(adapter Adapter) {
	DefaultRegistry.Register(adapter)
}

// GetAdapter returns an adapter from the default registry.
func GetAdapter(name string) (Adapter, bool) {
	return DefaultRegistry.GetAdapter(name)
}

// AdapterNames returns adapter names from the default registry.
func AdapterNames() []string {
	return DefaultRegistry.AdapterNames()
}

// Convert converts using the default registry.
func Convert(data []byte, from, to string) ([]byte, error) {
	return DefaultRegistry.Convert(data, from, to)
}

// ReadCanonicalFile reads a canonical skill file (JSON or Markdown with YAML frontmatter).
func ReadCanonicalFile(path string) (*Skill, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &ReadError{Path: path, Err: err}
	}

	// Detect format: if it starts with "---" or has .md extension, parse as markdown
	ext := filepath.Ext(path)
	if ext == ".md" || (len(data) >= 3 && string(data[:3]) == "---") {
		skill, err := ParseSkillMarkdown(data)
		if err != nil {
			return nil, &ParseError{Format: "markdown", Path: path, Err: err}
		}
		// Infer name from filename if not set
		if skill.Name == "" {
			base := filepath.Base(path)
			skill.Name = strings.TrimSuffix(base, filepath.Ext(base))
		}
		return skill, nil
	}

	// Fall back to JSON
	var skill Skill
	if err := json.Unmarshal(data, &skill); err != nil {
		return nil, &ParseError{Format: "canonical", Path: path, Err: err}
	}

	return &skill, nil
}

// WriteCanonicalFile writes a canonical skill.json file.
func WriteCanonicalFile(skill *Skill, path string) error {
	data, err := json.MarshalIndent(skill, "", "  ")
	if err != nil {
		return &MarshalError{Format: "canonical", Err: err}
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, DefaultDirMode); err != nil {
		return &WriteError{Path: path, Err: err}
	}

	if err := os.WriteFile(path, append(data, '\n'), DefaultFileMode); err != nil {
		return &WriteError{Path: path, Err: err}
	}

	return nil
}

// ReadCanonicalDir reads all skill files from a directory.
// Supports both:
// - Subdirectories with skill.json files
// - Direct .md files with YAML frontmatter
func ReadCanonicalDir(dir string) ([]*Skill, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, &ReadError{Path: dir, Err: err}
	}

	var skills []*Skill
	for _, entry := range entries {
		// Handle direct .md files (flat structure)
		if !entry.IsDir() {
			ext := filepath.Ext(entry.Name())
			if ext == ".md" {
				skillPath := filepath.Join(dir, entry.Name())
				skill, err := ReadCanonicalFile(skillPath)
				if err != nil {
					return nil, err
				}
				skills = append(skills, skill)
			}
			continue
		}

		// Handle subdirectories with skill.json
		skillPath := filepath.Join(dir, entry.Name(), "skill.json")
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			continue
		}

		skill, err := ReadCanonicalFile(skillPath)
		if err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

// WriteSkillsToDir writes multiple skills to a directory using the specified adapter.
func WriteSkillsToDir(skills []*Skill, dir string, adapterName string) error {
	adapter, ok := GetAdapter(adapterName)
	if !ok {
		return fmt.Errorf("unknown adapter: %s", adapterName)
	}

	if err := os.MkdirAll(dir, DefaultDirMode); err != nil {
		return &WriteError{Path: dir, Err: err}
	}

	for _, skill := range skills {
		if err := adapter.WriteSkillDir(skill, dir); err != nil {
			return err
		}
	}

	return nil
}

// ParseSkillMarkdown parses a Markdown file with YAML frontmatter into a Skill.
func ParseSkillMarkdown(data []byte) (*Skill, error) {
	content := string(data)

	if !strings.HasPrefix(content, "---") {
		// No frontmatter, treat entire content as instructions
		return &Skill{Instructions: strings.TrimSpace(content)}, nil
	}

	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return &Skill{Instructions: strings.TrimSpace(content)}, nil
	}

	skill := &Skill{}

	// Parse simple YAML key: value pairs from frontmatter
	lines := strings.Split(strings.TrimSpace(parts[1]), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, ":")
		if idx <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		// Remove quotes if present
		value = strings.Trim(value, "\"'")

		switch key {
		case "name":
			skill.Name = value
		case "description":
			skill.Description = value
		case "triggers":
			skill.Triggers = parseList(value)
		case "dependencies":
			skill.Dependencies = parseList(value)
		case "scripts":
			skill.Scripts = parseList(value)
		case "references":
			skill.References = parseList(value)
		case "assets":
			skill.Assets = parseList(value)
		}
	}

	// Body becomes instructions
	skill.Instructions = strings.TrimSpace(parts[2])

	return skill, nil
}

// parseList parses a comma-separated or bracket-enclosed list.
func parseList(s string) []string {
	s = strings.Trim(s, "[]")
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, "\"'")
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

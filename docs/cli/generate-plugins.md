# Generate Command

The `generate` command creates platform-specific plugins from a unified specs directory.

!!! note "v0.9.0 Update"
    As of v0.9.0, the main `generate` command is the recommended way to generate plugins. The `generate plugins`, `generate agents`, `generate all`, and `generate deployment` subcommands are deprecated.

## Synopsis

```bash
assistantkit generate [flags]
```

## Description

This command reads plugin definitions from a unified specs directory and generates complete platform-specific plugins for each deployment target. Each target receives agents, commands, skills, and plugin manifest.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--specs` | `specs` | Path to unified specs directory |
| `--target` | `local` | Deployment target (looks for `specs/deployments/<target>.json`) |
| `--output` | `.` | Output base directory for relative paths |

## Supported Platforms

- **claude-code**: Claude Code plugins (`.claude-plugin/`, commands/, skills/, agents/)
- **kiro-cli**: Kiro IDE Powers (POWER.md + mcp.json) or Kiro Agents (agents/*.json)
- **gemini-cli**: Gemini CLI extensions (gemini-extension.json, commands/, agents/)

## Specs Directory Structure

The unified specs directory should contain:

```
specs/
├── plugin.json          # Plugin metadata
├── agents/              # Agent definitions (*.md with YAML frontmatter)
│   ├── coordinator.md
│   └── writer.md
├── commands/            # Command definitions (*.md or *.json)
│   └── release.md
├── skills/              # Skill definitions (*.md or *.json)
│   └── review.md
├── teams/               # Team workflow definitions (optional)
│   └── my-team.json
└── deployments/         # Deployment configurations
    ├── local.json       # Local development (default)
    └── production.json  # Production deployment
```

### plugin.json

The plugin metadata file defines the plugin name, version, keywords, and MCP server configurations:

```json
{
  "name": "my-plugin",
  "displayName": "My Plugin",
  "version": "1.0.0",
  "description": "A plugin for AI assistants",
  "keywords": ["keyword1", "keyword2"],
  "mcpServers": {
    "my-server": {
      "command": "my-mcp-server",
      "args": []
    }
  }
}
```

### agents/*.md

Agent definitions using multi-agent-spec format with YAML frontmatter:

```markdown
---
name: release-coordinator
description: Orchestrates software releases
model: sonnet
tools: [Read, Write, Bash, Glob, Grep]
skills: [version-analysis, commit-classification]
---

You are a release coordinator agent responsible for...
```

### commands/*.md

Command definitions for slash commands:

```markdown
---
name: release
description: Execute full release workflow
arguments: [version]
dependencies: [version-analysis]
---

# Release Command

When executing a release, follow these steps...
```

### skills/*.md

Skill definitions for reusable capabilities:

```markdown
---
name: code-review
description: Reviews code for best practices
triggers: [review code, check code]
---

# Code Review Skill

When reviewing code, analyze for...
```

### deployments/*.json

Deployment configurations defining output targets:

```json
{
  "team": "my-team",
  "targets": [
    {
      "name": "local-claude",
      "platform": "claude-code",
      "output": "plugins/claude"
    },
    {
      "name": "local-kiro",
      "platform": "kiro-cli",
      "output": "plugins/kiro"
    },
    {
      "name": "local-gemini",
      "platform": "gemini-cli",
      "output": "plugins/gemini"
    }
  ]
}
```

## Generated Output

Each deployment target receives a complete plugin:

### Claude Code (`claude-code`)

```
plugins/claude/
├── .claude-plugin/
│   └── plugin.json       # Claude plugin manifest
├── commands/
│   └── release.md        # Command instructions
├── skills/
│   └── code-review/
│       └── SKILL.md      # Skill instructions
└── agents/
    └── release-coordinator.md  # Agent definition
```

### Kiro CLI (`kiro-cli`)

```
plugins/kiro/
├── POWER.md              # Power description (or agents/*.json)
├── mcp.json              # MCP server configuration
└── steering/
    └── code-review.md    # Steering files from skills
```

### Gemini CLI (`gemini-cli`)

```
plugins/gemini/
├── gemini-extension.json # Extension manifest
├── commands/
│   └── release.toml      # Command in TOML format
└── agents/
    └── release-coordinator.toml  # Agent in TOML format
```

## Examples

Generate plugins using defaults:

```bash
assistantkit generate
```

Use a specific deployment target:

```bash
assistantkit generate --target=production
```

Generate with custom directories:

```bash
assistantkit generate --specs=my-specs --target=local --output=/path/to/output
```

## Deprecated Subcommands

The following subcommands are deprecated and will show warnings when used:

| Deprecated | Replacement |
|------------|-------------|
| `generate plugins` | `generate --specs=... --target=...` |
| `generate agents` | `generate --specs=... --target=...` |
| `generate all` | `generate --specs=... --target=...` |
| `generate deployment` | `generate --specs=... --target=...` |

## See Also

- [Plugin Structure](../plugins/structure.md) - Learn about plugin components
- [Commands](../plugins/commands.md) - Command definition details
- [Skills](../plugins/skills.md) - Skill definition details
- [Agents](../plugins/agents.md) - Agent definition details

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/agentplexus/assistantkit/generate"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate platform-specific plugins from canonical specs",
	Long: `Generate platform-specific plugins from canonical JSON specifications.

Supported platforms:
  - claude: Claude Code plugins (.claude-plugin/)
  - kiro: Kiro IDE Powers (POWER.md + mcp.json)
  - gemini: Gemini CLI extensions (gemini-extension.json)`,
}

var (
	specDir    string
	outputDir  string
	platforms  []string
	configFile string
)

var generatePluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Generate plugins for all configured platforms",
	Long: `Generate platform-specific plugins from a canonical spec directory.

The spec directory should contain:
  - plugin.json: Plugin metadata
  - commands/: Command definitions (*.json)
  - skills/: Skill definitions (*.json)
  - agents/: Agent definitions (*.json)

Example:
  assistantkit generate plugins --spec=plugins/spec --output=plugins --platforms=claude,kiro`,
	RunE: runGeneratePlugins,
}

var (
	deploymentSpecDir string
	deploymentFile    string
)

var generateDeploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Generate deployment artifacts from multi-agent-spec definitions",
	Long: `Generate platform-specific deployment artifacts from multi-agent-spec format.

The specs directory should contain:
  - agents/: Agent definitions (*.md with YAML frontmatter)
  - teams/: Team definitions (*.json)
  - deployments/: Deployment definitions (*.json)

Each target in the deployment file specifies a platform and output directory.

Supported platforms:
  - claude-code: Claude Code agent markdown files
  - kiro-cli: Kiro CLI agent JSON files
  - gemini-cli: Gemini CLI agent TOML files

Example:
  assistantkit generate deployment --specs=specs --deployment=specs/deployments/my-team.json`,
	RunE: runGenerateDeployment,
}

func init() {
	generateCmd.AddCommand(generatePluginsCmd)
	generateCmd.AddCommand(generateDeploymentCmd)

	generatePluginsCmd.Flags().StringVar(&specDir, "spec", "plugins/spec", "Path to canonical spec directory")
	generatePluginsCmd.Flags().StringVar(&outputDir, "output", "plugins", "Output directory for generated plugins")
	generatePluginsCmd.Flags().StringSliceVar(&platforms, "platforms", []string{"claude", "kiro"}, "Platforms to generate (claude,kiro,gemini)")
	generatePluginsCmd.Flags().StringVar(&configFile, "config", "", "Config file (default: assistantkit.yaml if exists)")

	generateDeploymentCmd.Flags().StringVar(&deploymentSpecDir, "specs", "specs", "Path to multi-agent-spec directory")
	generateDeploymentCmd.Flags().StringVar(&deploymentFile, "deployment", "", "Path to deployment definition file (required)")
	_ = generateDeploymentCmd.MarkFlagRequired("deployment")
}

func runGenerateDeployment(cmd *cobra.Command, args []string) error {
	// Resolve paths
	absSpecsDir, err := filepath.Abs(deploymentSpecDir)
	if err != nil {
		return fmt.Errorf("resolving specs dir: %w", err)
	}

	absDeploymentFile, err := filepath.Abs(deploymentFile)
	if err != nil {
		return fmt.Errorf("resolving deployment file: %w", err)
	}

	// Validate paths exist
	if _, err := os.Stat(absSpecsDir); os.IsNotExist(err) {
		return fmt.Errorf("specs directory not found: %s", absSpecsDir)
	}
	if _, err := os.Stat(absDeploymentFile); os.IsNotExist(err) {
		return fmt.Errorf("deployment file not found: %s", absDeploymentFile)
	}

	// Print header
	fmt.Println("=== AssistantKit Deployment Generator ===")
	fmt.Printf("Specs directory: %s\n", absSpecsDir)
	fmt.Printf("Deployment file: %s\n", absDeploymentFile)
	fmt.Println()

	// Generate deployment
	result, err := generate.Deployment(absSpecsDir, absDeploymentFile)
	if err != nil {
		return fmt.Errorf("generating deployment: %w", err)
	}

	// Print results
	fmt.Printf("Team: %s\n", result.TeamName)
	fmt.Printf("Loaded: %d agents\n\n", result.AgentCount)

	fmt.Println("Generated targets:")
	for _, target := range result.TargetsGenerated {
		dir := result.GeneratedDirs[target]
		fmt.Printf("  - %s: %s\n", target, dir)
	}

	fmt.Println("\nDone!")
	return nil
}

func runGeneratePlugins(cmd *cobra.Command, args []string) error {
	// Resolve paths
	absSpecDir, err := filepath.Abs(specDir)
	if err != nil {
		return fmt.Errorf("resolving spec dir: %w", err)
	}

	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("resolving output dir: %w", err)
	}

	// Validate spec directory exists
	if _, err := os.Stat(absSpecDir); os.IsNotExist(err) {
		return fmt.Errorf("spec directory not found: %s", absSpecDir)
	}

	// Print header
	fmt.Println("=== AssistantKit Plugin Generator ===")
	fmt.Printf("Spec directory: %s\n", absSpecDir)
	fmt.Printf("Output directory: %s\n", absOutputDir)
	fmt.Printf("Platforms: %s\n", strings.Join(platforms, ", "))
	fmt.Println()

	// Generate plugins
	result, err := generate.Plugins(absSpecDir, absOutputDir, platforms)
	if err != nil {
		return fmt.Errorf("generating plugins: %w", err)
	}

	// Print results
	fmt.Printf("Loaded: %d commands, %d skills, %d agents\n\n",
		result.CommandCount, result.SkillCount, result.AgentCount)

	for platform, dir := range result.GeneratedDirs {
		fmt.Printf("Generated %s: %s\n", platform, dir)
	}

	fmt.Println("\nDone!")
	return nil
}

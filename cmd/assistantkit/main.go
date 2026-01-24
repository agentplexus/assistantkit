// Command assistantkit provides CLI tools for AI assistant plugin development.
//
// Usage:
//
//	assistantkit generate plugins [flags]
//	assistantkit generate power [flags]
//
// Generate plugins from canonical specs:
//
//	assistantkit generate plugins --spec=plugins/spec --output=plugins
//
// Generate a single power:
//
//	assistantkit generate power --name=mypower --output=~/.kiro/powers/mypower
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "assistantkit",
	Short: "CLI tools for AI assistant plugin development",
	Long: `assistantkit provides tools for creating, validating, and generating
AI assistant plugins across multiple platforms (Claude Code, Kiro IDE, Gemini CLI).

Use canonical JSON specs to define your plugins once, then generate
platform-specific formats automatically.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

package requirements

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Prompter handles human-in-the-loop interactions.
type Prompter interface {
	// Info displays an informational message.
	Info(message string)
	// Warn displays a warning message.
	Warn(message string)
	// Error displays an error message.
	Error(message string)
	// Confirm asks a yes/no question.
	Confirm(message string) (bool, error)
	// Choose presents options and returns the selected index (-1 if cancelled).
	Choose(message string, options []string) (int, error)
}

// CLIPrompter implements Prompter for terminal interaction.
type CLIPrompter struct {
	In  io.Reader
	Out io.Writer
}

// NewCLIPrompter creates a CLIPrompter using stdin/stdout.
func NewCLIPrompter() *CLIPrompter {
	return &CLIPrompter{
		In:  os.Stdin,
		Out: os.Stdout,
	}
}

func (p *CLIPrompter) Info(message string) {
	fmt.Fprintf(p.Out, "ℹ️  %s\n", message)
}

func (p *CLIPrompter) Warn(message string) {
	fmt.Fprintf(p.Out, "⚠️  %s\n", message)
}

func (p *CLIPrompter) Error(message string) {
	fmt.Fprintf(p.Out, "❌ %s\n", message)
}

func (p *CLIPrompter) Confirm(message string) (bool, error) {
	fmt.Fprintf(p.Out, "\n%s [y/N]: ", message)

	reader := bufio.NewReader(p.In)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes", nil
}

func (p *CLIPrompter) Choose(message string, options []string) (int, error) {
	fmt.Fprintf(p.Out, "\n%s\n", message)
	for i, opt := range options {
		fmt.Fprintf(p.Out, "  %d) %s\n", i+1, opt)
	}
	fmt.Fprintf(p.Out, "  0) Cancel\n")
	fmt.Fprintf(p.Out, "Choice: ")

	reader := bufio.NewReader(p.In)
	input, err := reader.ReadString('\n')
	if err != nil {
		return -1, err
	}

	input = strings.TrimSpace(input)
	if input == "" || input == "0" {
		return -1, nil
	}

	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(options) {
		return -1, fmt.Errorf("invalid choice: %s", input)
	}

	return choice - 1, nil
}

// EnsureRequirements checks requirements and prompts for installation of missing ones.
// Returns the final CheckResult after any installations.
func EnsureRequirements(requires []string, prompter Prompter) CheckResult {
	checker := NewChecker()
	return EnsureRequirementsWithChecker(requires, checker, prompter)
}

// EnsureRequirementsWithChecker uses a custom checker.
func EnsureRequirementsWithChecker(requires []string, checker *Checker, prompter Prompter) CheckResult {
	result := checker.Check(requires)

	if result.AllSatisfied() {
		return result
	}

	// Handle unknown tools
	for _, name := range result.Unknown {
		prompter.Error(fmt.Sprintf("Unknown tool '%s' is not installed and not in registry", name))
	}

	// Handle missing tools with available install methods
	for i, missing := range result.Missing {
		installed := promptForInstall(missing, prompter)
		if installed {
			// Move from missing to satisfied
			result.Satisfied = append(result.Satisfied, missing.Requirement.Name)
			// Mark as installed (we'll filter later)
			result.Missing[i].SuggestedMethod = nil
		}
	}

	// Filter out installed items from missing
	var stillMissing []MissingRequirement
	for _, m := range result.Missing {
		if m.SuggestedMethod != nil || len(m.AvailableMethods) == 0 {
			// Still missing (either not attempted or no methods available)
			if !checker.IsInstalled(m.Requirement.Name) {
				stillMissing = append(stillMissing, m)
			}
		}
	}
	result.Missing = stillMissing

	return result
}

// promptForInstall prompts the user to install a missing tool.
// Returns true if installation succeeded.
func promptForInstall(missing MissingRequirement, prompter Prompter) bool {
	req := missing.Requirement

	prompter.Warn(fmt.Sprintf("Required tool '%s' is not installed", req.Name))
	prompter.Info(fmt.Sprintf("Purpose: %s", req.Purpose))

	if len(missing.AvailableMethods) == 0 {
		prompter.Error("No install method available (missing prerequisites)")
		if len(req.InstallMethods) > 0 {
			prompter.Info("Possible install methods (prerequisites not met):")
			for _, m := range req.InstallMethods {
				prereqs := strings.Join(m.Requires, ", ")
				if prereqs == "" {
					prereqs = "none"
				}
				prompter.Info(fmt.Sprintf("  - %s: %s (requires: %s)", m.Name, m.Command, prereqs))
			}
		}
		if req.Homepage != "" {
			prompter.Info(fmt.Sprintf("See: %s", req.Homepage))
		}
		return false
	}

	// If only one method, ask to install directly
	if len(missing.AvailableMethods) == 1 {
		method := missing.AvailableMethods[0]
		prompter.Info(fmt.Sprintf("Install command: %s", method.Command))

		confirmed, err := prompter.Confirm(fmt.Sprintf("Install %s now?", req.Name))
		if err != nil || !confirmed {
			prompter.Info("Skipping installation")
			return false
		}

		return runInstall(method, prompter)
	}

	// Multiple methods - let user choose
	options := make([]string, len(missing.AvailableMethods))
	for i, m := range missing.AvailableMethods {
		options[i] = fmt.Sprintf("%s: %s", m.Name, m.Command)
	}

	choice, err := prompter.Choose("Choose install method:", options)
	if err != nil || choice < 0 {
		prompter.Info("Skipping installation")
		return false
	}

	return runInstall(missing.AvailableMethods[choice], prompter)
}

// runInstall executes an install command.
func runInstall(method InstallMethod, prompter Prompter) bool {
	prompter.Info(fmt.Sprintf("Running: %s", method.Command))

	cmd := exec.Command("sh", "-c", method.Command) //nolint:gosec // G204: intentional command execution for CLI tool
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		prompter.Error(fmt.Sprintf("Installation failed: %v", err))
		return false
	}

	prompter.Info("Installation completed successfully")
	return true
}

// FormatMissingError creates a user-friendly error message for missing requirements.
func FormatMissingError(result CheckResult) string {
	if result.AllSatisfied() {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("Missing required tools:\n")

	for _, m := range result.Missing {
		sb.WriteString(fmt.Sprintf("\n  %s - %s\n", m.Requirement.Name, m.Requirement.Purpose))
		if m.SuggestedMethod != nil {
			sb.WriteString(fmt.Sprintf("    Install: %s\n", m.SuggestedMethod.Command))
		} else if len(m.AvailableMethods) > 0 {
			sb.WriteString(fmt.Sprintf("    Install: %s\n", m.AvailableMethods[0].Command))
		}
		if m.Requirement.Homepage != "" {
			sb.WriteString(fmt.Sprintf("    See: %s\n", m.Requirement.Homepage))
		}
	}

	for _, name := range result.Unknown {
		sb.WriteString(fmt.Sprintf("\n  %s - unknown tool (not in registry)\n", name))
	}

	return sb.String()
}

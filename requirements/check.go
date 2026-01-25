package requirements

import (
	"os/exec"
	"runtime"
	"strings"
)

// Checker validates requirements and finds available install methods.
type Checker struct {
	Registry Registry
}

// NewChecker creates a Checker with the default registry.
func NewChecker() *Checker {
	return &Checker{Registry: DefaultRegistry}
}

// NewCheckerWithRegistry creates a Checker with a custom registry.
func NewCheckerWithRegistry(reg Registry) *Checker {
	return &Checker{Registry: reg}
}

// Check validates a list of required tools and returns the results.
func (c *Checker) Check(requires []string) CheckResult {
	result := CheckResult{}

	for _, name := range requires {
		req := c.Registry.Get(name)
		if req == nil {
			// Unknown tool - still check if it exists
			if c.isInstalled(name) {
				result.Satisfied = append(result.Satisfied, name)
			} else {
				result.Unknown = append(result.Unknown, name)
			}
			continue
		}

		if c.isInstalledByCheck(req.Check) {
			result.Satisfied = append(result.Satisfied, name)
		} else {
			missing := MissingRequirement{
				Requirement:      *req,
				AvailableMethods: c.findAvailableMethods(*req),
			}
			if len(missing.AvailableMethods) > 0 {
				missing.SuggestedMethod = &missing.AvailableMethods[0]
			}
			result.Missing = append(result.Missing, missing)
		}
	}

	return result
}

// IsInstalled checks if a single tool is installed.
func (c *Checker) IsInstalled(name string) bool {
	req := c.Registry.Get(name)
	if req == nil {
		return c.isInstalled(name)
	}
	return c.isInstalledByCheck(req.Check)
}

// isInstalled checks if a command exists in PATH.
func (c *Checker) isInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// isInstalledByCheck runs a check command to verify installation.
func (c *Checker) isInstalledByCheck(check string) bool {
	if check == "" {
		return false
	}

	// Parse the check command
	parts := strings.Fields(check)
	if len(parts) == 0 {
		return false
	}

	cmd := exec.Command(parts[0], parts[1:]...) //nolint:gosec // G204: intentional command execution for CLI tool
	err := cmd.Run()
	return err == nil
}

// findAvailableMethods returns install methods whose prerequisites are met.
func (c *Checker) findAvailableMethods(req Requirement) []InstallMethod {
	var available []InstallMethod
	platform := runtime.GOOS

	for _, method := range req.InstallMethods {
		// Check platform restriction
		if len(method.Platforms) > 0 {
			platformMatch := false
			for _, p := range method.Platforms {
				if p == platform {
					platformMatch = true
					break
				}
			}
			if !platformMatch {
				continue
			}
		}

		// Check prerequisites
		prereqsMet := true
		for _, prereq := range method.Requires {
			if !c.isInstalled(prereq) {
				prereqsMet = false
				break
			}
		}

		if prereqsMet {
			available = append(available, method)
		}
	}

	return available
}

// GetInstallCommand returns the best install command for a tool.
// Returns empty string if no method is available.
func (c *Checker) GetInstallCommand(name string) string {
	req := c.Registry.Get(name)
	if req == nil {
		return ""
	}

	methods := c.findAvailableMethods(*req)
	if len(methods) == 0 {
		return ""
	}

	return methods[0].Command
}

// GetAllInstallCommands returns all available install commands for a tool.
func (c *Checker) GetAllInstallCommands(name string) []InstallMethod {
	req := c.Registry.Get(name)
	if req == nil {
		return nil
	}

	return c.findAvailableMethods(*req)
}

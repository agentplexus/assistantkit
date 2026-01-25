// Package core provides the Publisher interface for marketplace submissions.
package core

import "context"

// Publisher defines the interface for publishing plugins to marketplaces.
type Publisher interface {
	// Name returns the marketplace identifier (e.g., "claude", "gemini").
	Name() string

	// Validate checks if the plugin directory has all required files.
	Validate(pluginDir string) error

	// Publish submits the plugin to the marketplace.
	// Returns the PR URL on success.
	Publish(ctx context.Context, opts PublishOptions) (*PublishResult, error)
}

// PublishOptions configures the publish operation.
type PublishOptions struct {
	// PluginDir is the local directory containing the plugin files.
	PluginDir string

	// PluginName is the name of the plugin (used as directory name in marketplace).
	PluginName string

	// GitHubToken is the GitHub personal access token for API authentication.
	// Required scopes: repo, workflow
	GitHubToken string

	// ForkOwner is the GitHub username/org that owns the fork.
	// If empty, uses the authenticated user.
	ForkOwner string

	// Branch is the name of the branch to create for the PR.
	// If empty, defaults to "add-<plugin-name>".
	Branch string

	// Title is the PR title.
	// If empty, defaults to "Add <plugin-name> plugin".
	Title string

	// Body is the PR description.
	// If empty, a default description is generated.
	Body string

	// DryRun if true, validates and prepares but doesn't create the PR.
	DryRun bool

	// Verbose enables detailed logging.
	Verbose bool
}

// PublishResult contains the result of a publish operation.
type PublishResult struct {
	// PRURL is the URL of the created pull request.
	PRURL string

	// PRNumber is the PR number.
	PRNumber int

	// Branch is the branch name used for the PR.
	Branch string

	// ForkURL is the URL of the fork repository.
	ForkURL string

	// Status is a human-readable status message.
	Status string

	// FilesAdded lists the files that were added/updated.
	FilesAdded []string
}

// MarketplaceConfig defines the target repository for a marketplace.
type MarketplaceConfig struct {
	// Owner is the GitHub org/user that owns the marketplace repo.
	Owner string

	// Repo is the repository name.
	Repo string

	// BaseBranch is the default branch to target (usually "main").
	BaseBranch string

	// PluginPath is the path within the repo where plugins are stored.
	// e.g., "external_plugins" for Claude official marketplace.
	PluginPath string

	// RequiredFiles lists files that must exist in the plugin directory.
	RequiredFiles []string
}

// DefaultFileMode is the default permission for created files.
const DefaultFileMode = 0600

// DefaultDirMode is the default permission for created directories.
const DefaultDirMode = 0755

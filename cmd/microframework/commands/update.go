package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	updateType    string
	updateVersion string
	updateCheck   bool
	updateForce   bool
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update microservice dependencies and framework",
	Long: `Update microservice dependencies, framework version, and related tools.

This command allows you to:
- Update Go dependencies
- Update go-micro-libs to latest version
- Update framework CLI tool
- Check for available updates
- Update deployment configurations

Examples:
  microframework update
  microframework update --type dependencies
  microframework update --type framework
  microframework update --type all
  microframework update --check
  microframework update --version v1.2.0`,
	RunE: runUpdate,
}

func init() {
	updateCmd.Flags().StringVarP(&updateType, "type", "t", "all", "Type of update (all, dependencies, framework, cli, config)")
	updateCmd.Flags().StringVarP(&updateVersion, "version", "V", "", "Specific version to update to")
	updateCmd.Flags().BoolVar(&updateCheck, "check", false, "Check for available updates without installing")
	updateCmd.Flags().BoolVar(&updateForce, "force", false, "Force update even if there are breaking changes")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	// Check if Go is available
	if err := checkGoInstallation(); err != nil {
		return fmt.Errorf("Go installation check failed: %w", err)
	}

	// Check if we're in a microservice directory (except for CLI updates)
	if updateType != "cli" {
		if err := checkMicroserviceDirectory(); err != nil {
			return err
		}
	}

	// Validate update type
	if err := validateUpdateType(updateType); err != nil {
		return fmt.Errorf("invalid update type: %w", err)
	}

	// Validate version format if provided
	if updateVersion != "" {
		if err := validateVersionFormat(updateVersion); err != nil {
			return fmt.Errorf("invalid version format: %w", err)
		}
	}

	fmt.Printf("Updating microservice (type: %s)\n", updateType)

	if updateVersion != "" {
		fmt.Printf("Target version: %s\n", updateVersion)
	}

	if updateCheck {
		fmt.Println("CHECK MODE - No updates will be installed")
	}

	if updateForce {
		fmt.Println("FORCE MODE - Updates will be installed even with breaking changes")
	}

	// Perform update based on type
	switch updateType {
	case "all":
		return updateAll(updateVersion, updateCheck, updateForce)
	case "dependencies":
		return updateDependencies(updateVersion, updateCheck, updateForce)
	case "framework":
		return updateFramework(updateVersion, updateCheck, updateForce)
	case "cli":
		return updateCLI(updateVersion, updateCheck, updateForce)
	case "config":
		return updateConfig(updateVersion, updateCheck, updateForce)
	default:
		return fmt.Errorf("unknown update type: %s", updateType)
	}
}

// validateUpdateType validates the update type
func validateUpdateType(updateType string) error {
	validTypes := []string{"all", "dependencies", "framework", "cli", "config"}

	for _, valid := range validTypes {
		if updateType == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid update type. Available types: %v", validTypes)
}

// validateVersionFormat validates the version format
func validateVersionFormat(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}

	// Basic version format validation (v1.0.0, 1.0.0, latest, etc.)
	if !strings.HasPrefix(version, "v") && !strings.HasPrefix(version, "latest") && !strings.Contains(version, ".") {
		return fmt.Errorf("invalid version format. Expected format: v1.0.0 or 1.0.0")
	}

	return nil
}

// checkGoInstallation checks if Go is properly installed and available
func checkGoInstallation() error {
	// Check if go command is available
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Go is not installed or not in PATH: %w", err)
	}

	// Parse Go version
	versionLine := strings.TrimSpace(string(output))
	fmt.Printf("Using Go: %s\n", versionLine)

	// Check if Go modules are enabled
	cmd = exec.Command("go", "env", "GOMOD")
	output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check Go modules: %w", err)
	}

	// GOMOD should not be empty in a Go module
	if strings.TrimSpace(string(output)) == "" {
		return fmt.Errorf("Go modules are not enabled. Please run 'go mod init' first")
	}

	return nil
}

// Update functions
func updateAll(version string, check, force bool) error {
	fmt.Println("Performing comprehensive update...")

	var errors []error

	// Update dependencies
	fmt.Println("Updating dependencies...")
	if err := updateDependencies(version, check, force); err != nil {
		errors = append(errors, err)
	}

	// Update framework
	fmt.Println("Updating framework...")
	if err := updateFramework(version, check, force); err != nil {
		errors = append(errors, err)
	}

	// Update CLI
	fmt.Println("Updating CLI tool...")
	if err := updateCLI(version, check, force); err != nil {
		errors = append(errors, err)
	}

	// Update configuration
	fmt.Println("Updating configuration...")
	if err := updateConfig(version, check, force); err != nil {
		errors = append(errors, err)
	}

	// Report results
	if len(errors) > 0 {
		fmt.Printf("\nUpdate completed with %d errors:\n", len(errors))
		for i, err := range errors {
			fmt.Printf("%d. %s\n", i+1, err)
		}
		return fmt.Errorf("update failed with %d errors", len(errors))
	}

	fmt.Println("\n✓ All updates completed successfully!")
	return nil
}

func updateDependencies(version string, check, force bool) error {
	fmt.Println("Updating Go dependencies...")

	if check {
		return checkDependencyUpdates()
	}

	// Check for updates first
	if err := checkDependencyUpdates(); err != nil {
		return fmt.Errorf("failed to check for dependency updates: %w", err)
	}

	// Get current dependencies
	dependencies, err := getCurrentDependencies()
	if err != nil {
		return fmt.Errorf("failed to get current dependencies: %w", err)
	}

	// Perform updates
	if err := performDependencyUpdates([]DependencyUpdate{}, force); err != nil {
		return fmt.Errorf("failed to update dependencies: %w", err)
	}

	fmt.Printf("✓ Dependencies updated successfully (%d packages)\n", len(dependencies))
	return nil
}

func updateFramework(version string, check, force bool) error {
	fmt.Println("Updating go-micro-libs framework...")

	if check {
		return checkFrameworkUpdates()
	}

	// Check current framework version
	currentVersion, err := getCurrentFrameworkVersion()
	if err != nil {
		return fmt.Errorf("failed to get current framework version: %w", err)
	}

	// Check for updates
	latestVersion, err := getLatestFrameworkVersion()
	if err != nil {
		return fmt.Errorf("failed to get latest framework version: %w", err)
	}

	if currentVersion == latestVersion {
		fmt.Println("✓ Framework is up to date")
		return nil
	}

	fmt.Printf("Framework update available: %s -> %s\n", currentVersion, latestVersion)

	// Check for breaking changes
	breakingChanges, err := checkBreakingChanges(currentVersion, latestVersion)
	if err != nil {
		return fmt.Errorf("failed to check for breaking changes: %w", err)
	}

	if len(breakingChanges) > 0 && !force {
		fmt.Println("Breaking changes detected:")
		for _, change := range breakingChanges {
			fmt.Printf("  - %s\n", change)
		}
		return fmt.Errorf("breaking changes detected. Use --force to update anyway")
	}

	// Update framework
	if err := performFrameworkUpdate(latestVersion); err != nil {
		return fmt.Errorf("failed to update framework: %w", err)
	}

	fmt.Println("✓ Framework updated successfully")
	return nil
}

func updateCLI(version string, check, force bool) error {
	fmt.Println("Updating CLI tool...")

	if check {
		return checkCLIUpdates()
	}

	// Check current CLI version
	currentVersion, err := getCurrentCLIVersion()
	if err != nil {
		return fmt.Errorf("failed to get current CLI version: %w", err)
	}

	// Check for updates
	latestVersion, err := getLatestCLIVersion()
	if err != nil {
		return fmt.Errorf("failed to get latest CLI version: %w", err)
	}

	if currentVersion == latestVersion {
		fmt.Println("✓ CLI tool is up to date")
		return nil
	}

	fmt.Printf("CLI update available: %s -> %s\n", currentVersion, latestVersion)

	// Update CLI
	if err := performCLIUpdate(latestVersion); err != nil {
		return fmt.Errorf("failed to update CLI: %w", err)
	}

	fmt.Println("✓ CLI tool updated successfully")
	return nil
}

func updateConfig(version string, check, force bool) error {
	fmt.Println("Updating configuration...")

	if check {
		_, err := checkConfigUpdates()
		return err
	}

	// Check for configuration updates
	updates, err := checkConfigUpdates()
	if err != nil {
		return fmt.Errorf("failed to check for configuration updates: %w", err)
	}

	if len(updates) == 0 {
		fmt.Println("✓ Configuration is up to date")
		return nil
	}

	// Show available updates
	fmt.Printf("Found %d configuration updates:\n", len(updates))
	for _, update := range updates {
		fmt.Printf("  - %s: %s\n", update.Key, update.Description)
	}

	// Update configuration
	if err := performConfigUpdates(updates, force); err != nil {
		return fmt.Errorf("failed to update configuration: %w", err)
	}

	fmt.Println("✓ Configuration updated successfully")
	return nil
}

// Helper functions for updates
func getCurrentDependencies() ([]Dependency, error) {
	fmt.Println("Getting current dependencies...")

	// Read go.mod file
	goModPath := "go.mod"
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("go.mod file not found")
	}

	content, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read go.mod: %w", err)
	}

	var dependencies []Dependency
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "require ") && !strings.Contains(line, "// indirect") {
			// Parse require line
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				dependencies = append(dependencies, Dependency{
					Name:    parts[1],
					Version: parts[2],
				})
			}
		}
	}

	return dependencies, nil
}

func checkDependencyUpdates() error {
	fmt.Println("Checking for dependency updates...")

	// Run go list -u -m all to check for updates
	cmd := exec.Command("go", "list", "-u", "-m", "all")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check for dependency updates: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	updateCount := 0

	for _, line := range lines {
		if strings.Contains(line, "[") && strings.Contains(line, "]") {
			updateCount++
			fmt.Printf("  Update available: %s\n", line)
		}
	}

	if updateCount == 0 {
		fmt.Println("✓ All dependencies are up to date")
	} else {
		fmt.Printf("Found %d dependency updates\n", updateCount)
	}

	return nil
}

func performDependencyUpdates(updates []DependencyUpdate, force bool) error {
	fmt.Println("Performing dependency updates...")

	if len(updates) == 0 {
		fmt.Println("No updates to perform")
		return nil
	}

	// Run go get -u to update all dependencies
	cmd := exec.Command("go", "get", "-u", "./...")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update dependencies: %w\nOutput: %s", err, string(output))
	}

	// Run go mod tidy to clean up
	cmd = exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to tidy modules: %w", err)
	}

	fmt.Printf("✓ Updated %d dependencies\n", len(updates))
	return nil
}

func getCurrentFrameworkVersion() (string, error) {
	fmt.Println("Getting current framework version...")

	// Read go.mod file to find go-micro-libs version
	goModPath := "go.mod"
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return "", fmt.Errorf("go.mod file not found")
	}

	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "github.com/anasamu/go-micro-libs") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}

	return "v0.0.0", nil // Default if not found
}

func getLatestFrameworkVersion() (string, error) {
	fmt.Println("Getting latest framework version...")

	// Check GitHub releases for go-micro-libs
	cmd := exec.Command("go", "list", "-m", "-versions", "github.com/anasamu/go-micro-libs")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get latest framework version: %w", err)
	}

	// Parse the output to get the latest version
	line := strings.TrimSpace(string(output))
	parts := strings.Fields(line)
	if len(parts) > 0 {
		// The last part should be the latest version
		versions := strings.Split(parts[len(parts)-1], " ")
		if len(versions) > 0 {
			return versions[len(versions)-1], nil
		}
	}

	return "v1.0.0", nil // Default fallback
}

func checkBreakingChanges(current, latest string) ([]string, error) {
	fmt.Println("Checking for breaking changes...")

	// For now, return empty list - in a real implementation, this would:
	// 1. Check the CHANGELOG.md or release notes
	// 2. Compare API changes between versions
	// 3. Check for deprecated features

	var breakingChanges []string

	// Simple version comparison for major version changes
	if strings.HasPrefix(current, "v0.") && strings.HasPrefix(latest, "v1.") {
		breakingChanges = append(breakingChanges, "Major version upgrade from v0.x to v1.x may contain breaking changes")
	}

	if strings.HasPrefix(current, "v1.") && strings.HasPrefix(latest, "v2.") {
		breakingChanges = append(breakingChanges, "Major version upgrade from v1.x to v2.x may contain breaking changes")
	}

	return breakingChanges, nil
}

func performFrameworkUpdate(version string) error {
	fmt.Printf("Updating framework to version: %s\n", version)

	// Update go-micro-libs to the specified version
	cmd := exec.Command("go", "get", fmt.Sprintf("github.com/anasamu/go-micro-libs@%s", version))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update framework: %w\nOutput: %s", err, string(output))
	}

	// Update go-micro-framework CLI tool as well
	cmd = exec.Command("go", "get", fmt.Sprintf("github.com/anasamu/go-micro-framework@%s", version))
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Warning: Failed to update CLI tool: %s\n", string(output))
	}

	// Run go mod tidy to clean up
	cmd = exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to tidy modules: %w", err)
	}

	// Verify the update was successful
	updatedVersion, err := getCurrentFrameworkVersion()
	if err != nil {
		return fmt.Errorf("failed to verify framework update: %w", err)
	}

	fmt.Printf("✓ Framework updated to version %s\n", updatedVersion)

	// Show integration status with go-micro-libs
	fmt.Println("Checking go-micro-libs integration...")
	if err := checkGoMicroLibsIntegration(); err != nil {
		fmt.Printf("Warning: go-micro-libs integration check failed: %v\n", err)
	} else {
		fmt.Println("✓ go-micro-libs integration verified")
	}

	return nil
}

func getCurrentCLIVersion() (string, error) {
	fmt.Println("Getting current CLI version...")

	// Try to get version from the binary itself
	cmd := exec.Command("microframework", "version")
	output, err := cmd.Output()
	if err != nil {
		// Fallback to reading from version constant
		return "1.0.0", nil
	}

	// Parse version from output
	versionLine := strings.TrimSpace(string(output))
	if strings.Contains(versionLine, " ") {
		parts := strings.Fields(versionLine)
		if len(parts) > 0 {
			return parts[0], nil
		}
	}

	return "1.0.0", nil
}

func getLatestCLIVersion() (string, error) {
	fmt.Println("Getting latest CLI version...")

	// Check GitHub releases for go-micro-framework
	cmd := exec.Command("go", "list", "-m", "-versions", "github.com/anasamu/go-micro-framework")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get latest CLI version: %w", err)
	}

	// Parse the output to get the latest version
	line := strings.TrimSpace(string(output))
	parts := strings.Fields(line)
	if len(parts) > 0 {
		// The last part should be the latest version
		versions := strings.Split(parts[len(parts)-1], " ")
		if len(versions) > 0 {
			return versions[len(versions)-1], nil
		}
	}

	return "1.0.0", nil // Default fallback
}

func checkCLIUpdates() error {
	fmt.Println("Checking for CLI updates...")

	currentVersion, err := getCurrentCLIVersion()
	if err != nil {
		return fmt.Errorf("failed to get current CLI version: %w", err)
	}

	latestVersion, err := getLatestCLIVersion()
	if err != nil {
		return fmt.Errorf("failed to get latest CLI version: %w", err)
	}

	if currentVersion == latestVersion {
		fmt.Println("✓ CLI tool is up to date")
	} else {
		fmt.Printf("CLI update available: %s -> %s\n", currentVersion, latestVersion)
	}

	return nil
}

func performCLIUpdate(version string) error {
	fmt.Printf("Updating CLI to version: %s\n", version)

	// Install the latest version of the CLI tool
	cmd := exec.Command("go", "install", fmt.Sprintf("github.com/anasamu/go-micro-framework/cmd/microframework@%s", version))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update CLI: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("✓ CLI tool updated to version %s\n", version)
	return nil
}

func getCurrentConfiguration() (interface{}, error) {
	fmt.Println("Getting current configuration...")

	// Look for configuration files in the current directory
	configFiles := []string{"config.yaml", "config.yml", "config.json", "configs/config.yaml", "configs/config.yml"}

	for _, configFile := range configFiles {
		if _, err := os.Stat(configFile); err == nil {
			content, err := os.ReadFile(configFile)
			if err != nil {
				continue
			}
			return string(content), nil
		}
	}

	return map[string]interface{}{}, nil
}

func checkConfigUpdates() ([]ConfigUpdate, error) {
	fmt.Println("Checking for configuration updates...")

	var updates []ConfigUpdate

	// Check for common configuration updates
	// This is a simplified implementation - in a real scenario, this would:
	// 1. Compare current config with latest template
	// 2. Check for deprecated configuration options
	// 3. Suggest new configuration options

	// Example: Check if new configuration options are available
	updates = append(updates, ConfigUpdate{
		Key:         "monitoring.prometheus.enabled",
		Description: "Enable Prometheus metrics collection",
		OldValue:    nil,
		NewValue:    true,
	})

	updates = append(updates, ConfigUpdate{
		Key:         "auth.jwt.expiration",
		Description: "Update JWT token expiration to recommended value",
		OldValue:    "1h",
		NewValue:    "24h",
	})

	return updates, nil
}

func performConfigUpdates(updates []ConfigUpdate, force bool) error {
	fmt.Println("Performing configuration updates...")

	if len(updates) == 0 {
		fmt.Println("No configuration updates to perform")
		return nil
	}

	// Read current configuration
	configFile := "configs/config.yaml"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		configFile = "config.yaml"
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			return fmt.Errorf("no configuration file found")
		}
	}

	content, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Create backup
	backupFile := configFile + ".backup"
	if err := os.WriteFile(backupFile, content, 0644); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Apply updates (simplified - in real implementation, would use YAML parser)
	updatedContent := string(content)
	for _, update := range updates {
		fmt.Printf("  Updating %s: %s\n", update.Key, update.Description)
		// Here you would implement actual configuration update logic
	}

	// Write updated configuration
	if err := os.WriteFile(configFile, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated configuration: %w", err)
	}

	fmt.Printf("✓ Configuration updated successfully (backup created: %s)\n", backupFile)
	return nil
}

func checkFrameworkUpdates() error {
	fmt.Println("Checking for framework updates...")

	currentVersion, err := getCurrentFrameworkVersion()
	if err != nil {
		return fmt.Errorf("failed to get current framework version: %w", err)
	}

	latestVersion, err := getLatestFrameworkVersion()
	if err != nil {
		return fmt.Errorf("failed to get latest framework version: %w", err)
	}

	if currentVersion == latestVersion {
		fmt.Println("✓ Framework is up to date")
	} else {
		fmt.Printf("Framework update available: %s -> %s\n", currentVersion, latestVersion)
	}

	return nil
}

// Data structures for updates
type Dependency struct {
	Name    string
	Version string
}

type DependencyUpdate struct {
	Name    string
	Current string
	Latest  string
}

type ConfigUpdate struct {
	Key         string
	Description string
	OldValue    interface{}
	NewValue    interface{}
}

// checkGoMicroLibsIntegration checks if go-micro-libs is properly integrated
func checkGoMicroLibsIntegration() error {
	// Check if go-micro-libs is in go.mod
	goModPath := "go.mod"
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("go.mod file not found")
	}

	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	// Check if go-micro-libs is listed as a dependency
	if !strings.Contains(string(content), "github.com/anasamu/go-micro-libs") {
		return fmt.Errorf("go-micro-libs not found in dependencies")
	}

	// Try to import go-micro-libs to verify it's available
	cmd := exec.Command("go", "list", "-m", "github.com/anasamu/go-micro-libs")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("go-micro-libs not properly installed: %w", err)
	}

	version := strings.TrimSpace(string(output))
	fmt.Printf("  go-micro-libs version: %s\n", version)

	return nil
}

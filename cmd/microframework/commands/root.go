package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version information
	version = "1.0.0"
	commit  = "dev"
	date    = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "microframework",
	Short: "Go Micro Framework - CLI tool for microservices development",
	Long: `Go Micro Framework is a powerful CLI tool that simplifies microservices development
by integrating the microservices-library-go libraries.

This tool allows you to:
- Generate new microservices with predefined structures
- Add features and integrations to existing services
- Manage configurations and deployments
- Integrate with various microservice patterns and providers

Examples:
  microframework new user-service --with-auth --with-database
  microframework add ai --provider openai
  microframework generate handler user
  microframework deploy --env production`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(versionCmd)

	// Global flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.microframework.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("dry-run", "", false, "show what would be done without making changes")
}

// GetRootCmd returns the root command for use in main.go
func GetRootCmd() *cobra.Command {
	return rootCmd
}

package commands

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/anasamu/go-micro-libs/database"
	"github.com/anasamu/go-micro-libs/database/migrations"
	"github.com/anasamu/go-micro-libs/database/providers/cassandra"
	"github.com/anasamu/go-micro-libs/database/providers/cockroachdb"
	"github.com/anasamu/go-micro-libs/database/providers/influxdb"
	"github.com/anasamu/go-micro-libs/database/providers/mongodb"
	"github.com/anasamu/go-micro-libs/database/providers/mysql"
	"github.com/anasamu/go-micro-libs/database/providers/postgresql"
	"github.com/anasamu/go-micro-libs/database/providers/redis"
	"github.com/anasamu/go-micro-libs/database/providers/sqlite"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long: `Database migration commands for managing database schema changes.

This command integrates with the go-micro-libs database migration system to provide:
- Create new migration files
- Apply pending migrations (up)
- Rollback migrations (down)
- Check migration status
- Reset database
- Validate migration files

Examples:
  microframework migrate create add_users_table
  microframework migrate up
  microframework migrate down
  microframework migrate status
  microframework migrate reset
  microframework migrate validate`,
}

var (
	migrateProvider string
	migrateDir      string
	migrateName     string
	migrateConfig   string
	migrateVerbose  bool
	migrateTable    string
)

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Global flags for migrate command
	migrateCmd.PersistentFlags().StringVar(&migrateProvider, "provider", "postgresql", "Database provider (postgresql, mysql, sqlite, cassandra, cockroachdb, influxdb, mongodb, redis)")
	migrateCmd.PersistentFlags().StringVar(&migrateDir, "dir", "./migrations", "Migrations directory")
	migrateCmd.PersistentFlags().StringVar(&migrateConfig, "config", "", "Configuration file path")
	migrateCmd.PersistentFlags().BoolVar(&migrateVerbose, "verbose", false, "Enable verbose logging")
	migrateCmd.PersistentFlags().StringVar(&migrateTable, "table", "schema_migrations", "Migration table name")

	// Subcommands
	migrateCmd.AddCommand(migrateCreateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
	migrateCmd.AddCommand(migrateResetCmd)
	migrateCmd.AddCommand(migrateValidateCmd)
}

// migrateCreateCmd creates a new migration file
var migrateCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration file",
	Long:  `Create a new migration file with the specified name.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := runMigrateCreate(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating migration: %v\n", err)
			os.Exit(1)
		}
	},
}

// migrateUpCmd applies all pending migrations
var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	Long:  `Apply all pending migrations to the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runMigrateUp(); err != nil {
			fmt.Fprintf(os.Stderr, "Error applying migrations: %v\n", err)
			os.Exit(1)
		}
	},
}

// migrateDownCmd rolls back the last migration
var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	Long:  `Rollback the last applied migration.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runMigrateDown(); err != nil {
			fmt.Fprintf(os.Stderr, "Error rolling back migration: %v\n", err)
			os.Exit(1)
		}
	},
}

// migrateStatusCmd shows the status of all migrations
var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Long:  `Show the status of all migrations (applied and pending).`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runMigrateStatus(); err != nil {
			fmt.Fprintf(os.Stderr, "Error getting migration status: %v\n", err)
			os.Exit(1)
		}
	},
}

// migrateResetCmd resets the database and reapplies all migrations
var migrateResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset database and reapply all migrations",
	Long:  `Reset the database by rolling back all migrations and then reapplying them.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runMigrateReset(); err != nil {
			fmt.Fprintf(os.Stderr, "Error resetting database: %v\n", err)
			os.Exit(1)
		}
	},
}

// migrateValidateCmd validates all migration files
var migrateValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate migration files",
	Long:  `Validate all migration files for correctness.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runMigrateValidate(); err != nil {
			fmt.Fprintf(os.Stderr, "Error validating migrations: %v\n", err)
			os.Exit(1)
		}
	},
}

// runMigrateCreate creates a new migration file
func runMigrateCreate(name string) error {
	// Setup logger
	logger := setupLogger()

	// Create migrations directory if it doesn't exist
	if err := os.MkdirAll(migrateDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Create a dummy provider for CLI operations (we only need it for the CLI manager)
	provider, err := createProvider(migrateProvider, logger)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	// Create CLI manager
	cliManager := migrations.NewCLIManager(provider, migrateDir, logger)

	// Set custom migration table name if provided
	if migrateTable != "schema_migrations" {
		if err := cliManager.SetMigrationTableName(migrateTable); err != nil {
			return fmt.Errorf("failed to set migration table name: %w", err)
		}
	}

	// Create migration
	if err := cliManager.CreateMigration(name); err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}

	fmt.Printf("Migration '%s' created successfully in %s\n", name, migrateDir)
	return nil
}

// runMigrateUp applies all pending migrations
func runMigrateUp() error {
	// Setup logger
	logger := setupLogger()

	// Create database manager
	config := database.DefaultManagerConfig()
	databaseManager := database.NewDatabaseManager(config, logger)

	// Create and register provider
	provider, err := createProvider(migrateProvider, logger)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	if err := databaseManager.RegisterProvider(provider); err != nil {
		return fmt.Errorf("failed to register provider: %w", err)
	}

	// Connect to database
	ctx := context.Background()
	if err := databaseManager.Connect(ctx, migrateProvider); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer databaseManager.Close()

	// Get provider and create CLI manager
	dbProvider, err := databaseManager.GetProvider(migrateProvider)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	cliManager := migrations.NewCLIManager(dbProvider, migrateDir, logger)

	// Set custom migration table name if provided
	if migrateTable != "schema_migrations" {
		if err := cliManager.SetMigrationTableName(migrateTable); err != nil {
			return fmt.Errorf("failed to set migration table name: %w", err)
		}
	}

	// Apply migrations
	if err := cliManager.Up(ctx); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully")
	return nil
}

// runMigrateDown rolls back the last migration
func runMigrateDown() error {
	// Setup logger
	logger := setupLogger()

	// Create database manager
	config := database.DefaultManagerConfig()
	databaseManager := database.NewDatabaseManager(config, logger)

	// Create and register provider
	provider, err := createProvider(migrateProvider, logger)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	if err := databaseManager.RegisterProvider(provider); err != nil {
		return fmt.Errorf("failed to register provider: %w", err)
	}

	// Connect to database
	ctx := context.Background()
	if err := databaseManager.Connect(ctx, migrateProvider); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer databaseManager.Close()

	// Get provider and create CLI manager
	dbProvider, err := databaseManager.GetProvider(migrateProvider)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	cliManager := migrations.NewCLIManager(dbProvider, migrateDir, logger)

	// Set custom migration table name if provided
	if migrateTable != "schema_migrations" {
		if err := cliManager.SetMigrationTableName(migrateTable); err != nil {
			return fmt.Errorf("failed to set migration table name: %w", err)
		}
	}

	// Rollback migration
	if err := cliManager.Down(ctx); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	fmt.Println("Migration rolled back successfully")
	return nil
}

// runMigrateStatus shows the status of all migrations
func runMigrateStatus() error {
	// Setup logger
	logger := setupLogger()

	// Create database manager
	config := database.DefaultManagerConfig()
	databaseManager := database.NewDatabaseManager(config, logger)

	// Create and register provider
	provider, err := createProvider(migrateProvider, logger)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	if err := databaseManager.RegisterProvider(provider); err != nil {
		return fmt.Errorf("failed to register provider: %w", err)
	}

	// Connect to database
	ctx := context.Background()
	if err := databaseManager.Connect(ctx, migrateProvider); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer databaseManager.Close()

	// Get provider and create CLI manager
	dbProvider, err := databaseManager.GetProvider(migrateProvider)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	cliManager := migrations.NewCLIManager(dbProvider, migrateDir, logger)

	// Set custom migration table name if provided
	if migrateTable != "schema_migrations" {
		if err := cliManager.SetMigrationTableName(migrateTable); err != nil {
			return fmt.Errorf("failed to set migration table name: %w", err)
		}
	}

	// Show status
	if err := cliManager.Status(ctx); err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	return nil
}

// runMigrateReset resets the database and reapplies all migrations
func runMigrateReset() error {
	// Setup logger
	logger := setupLogger()

	// Create database manager
	config := database.DefaultManagerConfig()
	databaseManager := database.NewDatabaseManager(config, logger)

	// Create and register provider
	provider, err := createProvider(migrateProvider, logger)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	if err := databaseManager.RegisterProvider(provider); err != nil {
		return fmt.Errorf("failed to register provider: %w", err)
	}

	// Connect to database
	ctx := context.Background()
	if err := databaseManager.Connect(ctx, migrateProvider); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer databaseManager.Close()

	// Get provider and create CLI manager
	dbProvider, err := databaseManager.GetProvider(migrateProvider)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	cliManager := migrations.NewCLIManager(dbProvider, migrateDir, logger)

	// Set custom migration table name if provided
	if migrateTable != "schema_migrations" {
		if err := cliManager.SetMigrationTableName(migrateTable); err != nil {
			return fmt.Errorf("failed to set migration table name: %w", err)
		}
	}

	// Reset database
	if err := cliManager.Reset(ctx); err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}

	fmt.Println("Database reset successfully")
	return nil
}

// runMigrateValidate validates all migration files
func runMigrateValidate() error {
	// Setup logger
	logger := setupLogger()

	// Create a dummy provider for CLI operations (we only need it for the CLI manager)
	provider, err := createProvider(migrateProvider, logger)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	// Create CLI manager
	cliManager := migrations.NewCLIManager(provider, migrateDir, logger)

	// Set custom migration table name if provided
	if migrateTable != "schema_migrations" {
		if err := cliManager.SetMigrationTableName(migrateTable); err != nil {
			return fmt.Errorf("failed to set migration table name: %w", err)
		}
	}

	// Validate migrations
	if err := cliManager.Validate(); err != nil {
		return fmt.Errorf("migration validation failed: %w", err)
	}

	fmt.Println("All migrations are valid")
	return nil
}

// setupLogger creates a logger with appropriate level
func setupLogger() *logrus.Logger {
	logger := logrus.New()
	if migrateVerbose {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
	return logger
}

// createProvider creates a database provider based on the provider name
func createProvider(providerName string, logger *logrus.Logger) (database.DatabaseProvider, error) {
	switch providerName {
	case "postgresql":
		provider := postgresql.NewProvider(logger)
		config := map[string]interface{}{
			"host":     getEnv("POSTGRES_HOST", "localhost"),
			"port":     getEnvInt("POSTGRES_PORT", 5432),
			"user":     getEnv("POSTGRES_USER", "postgres"),
			"password": getEnv("POSTGRES_PASSWORD", "password"),
			"database": getEnv("POSTGRES_DATABASE", "testdb"),
			"ssl_mode": getEnv("POSTGRES_SSL_MODE", "disable"),
		}
		if err := provider.Configure(config); err != nil {
			return nil, err
		}
		return provider, nil

	case "mysql":
		provider := mysql.NewProvider(logger)
		config := map[string]interface{}{
			"host":     getEnv("MYSQL_HOST", "localhost"),
			"port":     getEnvInt("MYSQL_PORT", 3306),
			"user":     getEnv("MYSQL_USER", "root"),
			"password": getEnv("MYSQL_PASSWORD", "password"),
			"database": getEnv("MYSQL_DATABASE", "testdb"),
		}
		if err := provider.Configure(config); err != nil {
			return nil, err
		}
		return provider, nil

	case "sqlite":
		provider := sqlite.NewProvider(logger)
		config := map[string]interface{}{
			"file": getEnv("SQLITE_FILE", "./test.db"),
		}
		if err := provider.Configure(config); err != nil {
			return nil, err
		}
		return provider, nil

	case "cassandra":
		provider := cassandra.NewProvider(logger)
		config := map[string]interface{}{
			"hosts":       []string{getEnv("CASSANDRA_HOST", "localhost")},
			"keyspace":    getEnv("CASSANDRA_KEYSPACE", "test_keyspace"),
			"username":    getEnv("CASSANDRA_USERNAME", ""),
			"password":    getEnv("CASSANDRA_PASSWORD", ""),
			"consistency": getEnv("CASSANDRA_CONSISTENCY", "quorum"),
		}
		if err := provider.Configure(config); err != nil {
			return nil, err
		}
		return provider, nil

	case "cockroachdb":
		provider := cockroachdb.NewProvider(logger)
		config := map[string]interface{}{
			"host":     getEnv("COCKROACHDB_HOST", "localhost"),
			"port":     getEnvInt("COCKROACHDB_PORT", 26257),
			"user":     getEnv("COCKROACHDB_USER", "root"),
			"password": getEnv("COCKROACHDB_PASSWORD", ""),
			"database": getEnv("COCKROACHDB_DATABASE", "defaultdb"),
			"ssl_mode": getEnv("COCKROACHDB_SSL_MODE", "require"),
			"cluster":  getEnv("COCKROACHDB_CLUSTER", ""),
		}
		if err := provider.Configure(config); err != nil {
			return nil, err
		}
		return provider, nil

	case "influxdb":
		provider := influxdb.NewProvider(logger)
		config := map[string]interface{}{
			"url":    getEnv("INFLUXDB_URL", "http://localhost:8086"),
			"token":  getEnv("INFLUXDB_TOKEN", ""),
			"org":    getEnv("INFLUXDB_ORG", ""),
			"bucket": getEnv("INFLUXDB_BUCKET", ""),
		}
		if err := provider.Configure(config); err != nil {
			return nil, err
		}
		return provider, nil

	case "mongodb":
		provider := mongodb.NewProvider(logger)
		config := map[string]interface{}{
			"uri":      getEnv("MONGO_URI", "mongodb://localhost:27017"),
			"database": getEnv("MONGO_DATABASE", "testdb"),
		}
		if err := provider.Configure(config); err != nil {
			return nil, err
		}
		return provider, nil

	case "redis":
		provider := redis.NewProvider(logger)
		config := map[string]interface{}{
			"host": getEnv("REDIS_HOST", "localhost"),
			"port": getEnvInt("REDIS_PORT", 6379),
			"db":   getEnvInt("REDIS_DB", 0),
		}
		if err := provider.Configure(config); err != nil {
			return nil, err
		}
		return provider, nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerName)
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an environment variable as integer with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

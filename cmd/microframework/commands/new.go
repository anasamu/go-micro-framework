package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anasamu/go-micro-framework/internal/generator"
	"github.com/spf13/cobra"
)

var (
	serviceType        string
	withAuth           bool
	withDatabase       bool
	withMessaging      bool
	withMonitoring     bool
	withAI             bool
	withStorage        bool
	withCache          bool
	withDiscovery      bool
	withCircuitBreaker bool
	withRateLimit      bool
	withChaos          bool
	withFailover       bool
	withEvent          bool
	withScheduling     bool
	withBackup         bool
	withPayment        bool
	withFileGen        bool
	outputDir          string
	force              bool
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <service-name>",
	Short: "Generate a new microservice",
	Long: `Generate a new microservice with the specified name and features.

This command creates a complete microservice project structure including:
- Go source code with proper architecture
- Configuration files
- Docker and Kubernetes deployment files
- Tests and documentation
- Integration with microservices-library-go

Examples:
  microframework new user-service
  microframework new order-service --with-auth --with-database
  microframework new notification-service --with-messaging --with-ai
  microframework new payment-service --with-payment --with-database --with-monitoring`,
	Args: cobra.ExactArgs(1),
	RunE: runNew,
}

func init() {
	// Service type
	newCmd.Flags().StringVarP(&serviceType, "type", "t", "rest", "Service type (rest, graphql, grpc, websocket, event, scheduled, worker, gateway, proxy)")

	// Core features
	newCmd.Flags().BoolVar(&withAuth, "with-auth", false, "Include authentication (JWT, OAuth)")
	newCmd.Flags().BoolVar(&withDatabase, "with-database", false, "Include database (PostgreSQL, Redis)")
	newCmd.Flags().BoolVar(&withMessaging, "with-messaging", false, "Include messaging (Kafka, RabbitMQ)")
	newCmd.Flags().BoolVar(&withMonitoring, "with-monitoring", false, "Include monitoring (Prometheus, Jaeger, Grafana)")

	// Optional features
	newCmd.Flags().BoolVar(&withAI, "with-ai", false, "Include AI services (OpenAI, Anthropic)")
	newCmd.Flags().BoolVar(&withStorage, "with-storage", false, "Include storage (S3, GCS)")
	newCmd.Flags().BoolVar(&withCache, "with-cache", false, "Include caching (Redis, Memory)")
	newCmd.Flags().BoolVar(&withDiscovery, "with-discovery", false, "Include service discovery (Consul, Kubernetes)")
	newCmd.Flags().BoolVar(&withCircuitBreaker, "with-circuit-breaker", false, "Include circuit breaker patterns")
	newCmd.Flags().BoolVar(&withRateLimit, "with-rate-limit", false, "Include rate limiting")
	newCmd.Flags().BoolVar(&withChaos, "with-chaos", false, "Include chaos engineering")
	newCmd.Flags().BoolVar(&withFailover, "with-failover", false, "Include failover mechanisms")
	newCmd.Flags().BoolVar(&withEvent, "with-event", false, "Include event sourcing")
	newCmd.Flags().BoolVar(&withScheduling, "with-scheduling", false, "Include task scheduling")
	newCmd.Flags().BoolVar(&withBackup, "with-backup", false, "Include backup services")
	newCmd.Flags().BoolVar(&withPayment, "with-payment", false, "Include payment processing")
	newCmd.Flags().BoolVar(&withFileGen, "with-filegen", false, "Include file generation")

	// Output options
	newCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory for the generated service")
	newCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing files")
}

func runNew(cmd *cobra.Command, args []string) error {
	serviceName := args[0]

	// Validate service name
	if err := validateServiceName(serviceName); err != nil {
		return fmt.Errorf("invalid service name: %w", err)
	}

	// Check if output directory exists and is not empty
	fullOutputDir := filepath.Join(outputDir, serviceName)
	if !force {
		if err := checkOutputDirectory(fullOutputDir); err != nil {
			return err
		}
	}

	// Create generator configuration
	config := &generator.GeneratorConfig{
		ServiceName:        serviceName,
		ServiceType:        serviceType,
		WithAuth:           withAuth,
		WithDatabase:       withDatabase,
		WithMessaging:      withMessaging,
		WithMonitoring:     withMonitoring,
		WithAI:             withAI,
		WithStorage:        withStorage,
		WithCache:          withCache,
		WithDiscovery:      withDiscovery,
		WithCircuitBreaker: withCircuitBreaker,
		WithRateLimit:      withRateLimit,
		WithChaos:          withChaos,
		WithFailover:       withFailover,
		WithEvent:          withEvent,
		WithScheduling:     withScheduling,
		WithBackup:         withBackup,
		WithPayment:        withPayment,
		WithFileGen:        withFileGen,
		OutputDir:          outputDir,
	}

	// Create service generator
	generator := generator.NewServiceGenerator(config)

	// Generate the service
	fmt.Printf("Generating microservice: %s\n", serviceName)
	fmt.Printf("Service type: %s\n", serviceType)
	fmt.Printf("Output directory: %s\n", fullOutputDir)

	if withAuth {
		fmt.Println("✓ Authentication enabled")
	}
	if withDatabase {
		fmt.Println("✓ Database integration enabled")
	}
	if withMessaging {
		fmt.Println("✓ Messaging enabled")
	}
	if withMonitoring {
		fmt.Println("✓ Monitoring enabled")
	}
	if withAI {
		fmt.Println("✓ AI services enabled")
	}
	if withStorage {
		fmt.Println("✓ Storage enabled")
	}
	if withCache {
		fmt.Println("✓ Caching enabled")
	}
	if withDiscovery {
		fmt.Println("✓ Service discovery enabled")
	}
	if withCircuitBreaker {
		fmt.Println("✓ Circuit breaker enabled")
	}
	if withRateLimit {
		fmt.Println("✓ Rate limiting enabled")
	}
	if withChaos {
		fmt.Println("✓ Chaos engineering enabled")
	}
	if withFailover {
		fmt.Println("✓ Failover enabled")
	}
	if withEvent {
		fmt.Println("✓ Event sourcing enabled")
	}
	if withScheduling {
		fmt.Println("✓ Task scheduling enabled")
	}
	if withBackup {
		fmt.Println("✓ Backup services enabled")
	}
	if withPayment {
		fmt.Println("✓ Payment processing enabled")
	}
	if withFileGen {
		fmt.Println("✓ File generation enabled")
	}

	fmt.Println("\nGenerating service structure...")

	if err := generator.GenerateService(); err != nil {
		return fmt.Errorf("failed to generate service: %w", err)
	}

	fmt.Printf("\n✓ Service '%s' generated successfully!\n", serviceName)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("1. cd %s\n", fullOutputDir)
	fmt.Printf("2. go mod tidy\n")
	fmt.Printf("3. cp .env.example .env\n")
	fmt.Printf("4. Edit .env with your configuration\n")
	fmt.Printf("5. go run cmd/main.go\n")
	fmt.Printf("\nFor more information, see the README.md file.\n")

	return nil
}

// validateServiceName validates the service name
func validateServiceName(name string) error {
	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	// Check for valid characters (alphanumeric and hyphens)
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-') {
			return fmt.Errorf("service name can only contain alphanumeric characters and hyphens")
		}
	}

	// Check length
	if len(name) < 3 || len(name) > 50 {
		return fmt.Errorf("service name must be between 3 and 50 characters")
	}

	// Check if it starts and ends with alphanumeric
	if name[0] == '-' || name[len(name)-1] == '-' {
		return fmt.Errorf("service name cannot start or end with a hyphen")
	}

	return nil
}

// checkOutputDirectory checks if the output directory exists and is not empty
func checkOutputDirectory(path string) error {
	if _, err := os.Stat(path); err == nil {
		// Directory exists, check if it's empty
		entries, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("failed to read directory: %w", err)
		}

		if len(entries) > 0 {
			return fmt.Errorf("directory '%s' already exists and is not empty. Use --force to overwrite", path)
		}
	}

	return nil
}

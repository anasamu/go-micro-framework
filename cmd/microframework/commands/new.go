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
	withAuth           string
	withDatabase       string
	withMessaging      string
	withMonitoring     string
	withAI             string
	withStorage        string
	withCache          string
	withDiscovery      string
	withCircuitBreaker string
	withRateLimit      string
	withChaos          string
	withFailover       string
	withEvent          string
	withScheduling     string
	withBackup         string
	withPayment        string
	withFileGen        string
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
  microframework new order-service --with-auth=jwt --with-database=postgres
  microframework new notification-service --with-messaging=kafka --with-ai=openai
  microframework new payment-service --with-payment=stripe --with-database=postgres --with-monitoring=prometheus`,
	Args: cobra.ExactArgs(1),
	RunE: runNew,
}

func init() {
	// Service type
	newCmd.Flags().StringVarP(&serviceType, "type", "t", "rest", "Service type (rest, graphql, grpc, websocket, event, scheduled, worker, gateway, proxy)")

	// Core features
	newCmd.Flags().StringVar(&withAuth, "with-auth", "", "Include authentication (jwt, oauth, ldap, saml)")
	newCmd.Flags().StringVar(&withDatabase, "with-database", "", "Include database (postgres, mysql, redis, mongodb)")
	newCmd.Flags().StringVar(&withMessaging, "with-messaging", "", "Include messaging (kafka, rabbitmq, nats)")
	newCmd.Flags().StringVar(&withMonitoring, "with-monitoring", "", "Include monitoring (prometheus, jaeger, grafana)")

	// Optional features
	newCmd.Flags().StringVar(&withAI, "with-ai", "", "Include AI services (openai, anthropic, google)")
	newCmd.Flags().StringVar(&withStorage, "with-storage", "", "Include storage (s3, gcs, azure)")
	newCmd.Flags().StringVar(&withCache, "with-cache", "", "Include caching (redis, memcached, memory)")
	newCmd.Flags().StringVar(&withDiscovery, "with-discovery", "", "Include service discovery (consul, kubernetes)")
	newCmd.Flags().StringVar(&withCircuitBreaker, "with-circuit-breaker", "", "Include circuit breaker patterns")
	newCmd.Flags().StringVar(&withRateLimit, "with-rate-limit", "", "Include rate limiting")
	newCmd.Flags().StringVar(&withChaos, "with-chaos", "", "Include chaos engineering")
	newCmd.Flags().StringVar(&withFailover, "with-failover", "", "Include failover mechanisms")
	newCmd.Flags().StringVar(&withEvent, "with-event", "", "Include event sourcing")
	newCmd.Flags().StringVar(&withScheduling, "with-scheduling", "", "Include task scheduling")
	newCmd.Flags().StringVar(&withBackup, "with-backup", "", "Include backup services")
	newCmd.Flags().StringVar(&withPayment, "with-payment", "", "Include payment processing")
	newCmd.Flags().StringVar(&withFileGen, "with-filegen", "", "Include file generation")

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
		WithAuth:           withAuth != "",
		WithDatabase:       withDatabase != "",
		WithMessaging:      withMessaging != "",
		WithMonitoring:     withMonitoring != "",
		WithAI:             withAI != "",
		WithStorage:        withStorage != "",
		WithCache:          withCache != "",
		WithDiscovery:      withDiscovery != "",
		WithCircuitBreaker: withCircuitBreaker != "",
		WithRateLimit:      withRateLimit != "",
		WithChaos:          withChaos != "",
		WithFailover:       withFailover != "",
		WithEvent:          withEvent != "",
		WithScheduling:     withScheduling != "",
		WithBackup:         withBackup != "",
		WithPayment:        withPayment != "",
		WithFileGen:        withFileGen != "",
		OutputDir:          outputDir,
		// Provider specifications
		AuthProvider:       withAuth,
		DatabaseProvider:   withDatabase,
		MessagingProvider:  withMessaging,
		MonitoringProvider: withMonitoring,
		AIProvider:         withAI,
		StorageProvider:    withStorage,
		CacheProvider:      withCache,
		DiscoveryProvider:  withDiscovery,
		PaymentProvider:    withPayment,
	}

	// Create service generator
	generator := generator.NewServiceGenerator(config)

	// Generate the service
	fmt.Printf("Generating microservice: %s\n", serviceName)
	fmt.Printf("Service type: %s\n", serviceType)
	fmt.Printf("Output directory: %s\n", fullOutputDir)

	if withAuth != "" {
		fmt.Printf("✓ Authentication enabled (%s)\n", withAuth)
	}
	if withDatabase != "" {
		fmt.Printf("✓ Database integration enabled (%s)\n", withDatabase)
	}
	if withMessaging != "" {
		fmt.Printf("✓ Messaging enabled (%s)\n", withMessaging)
	}
	if withMonitoring != "" {
		fmt.Printf("✓ Monitoring enabled (%s)\n", withMonitoring)
	}
	if withAI != "" {
		fmt.Printf("✓ AI services enabled (%s)\n", withAI)
	}
	if withStorage != "" {
		fmt.Printf("✓ Storage enabled (%s)\n", withStorage)
	}
	if withCache != "" {
		fmt.Printf("✓ Caching enabled (%s)\n", withCache)
	}
	if withDiscovery != "" {
		fmt.Printf("✓ Service discovery enabled (%s)\n", withDiscovery)
	}
	if withCircuitBreaker != "" {
		fmt.Printf("✓ Circuit breaker enabled (%s)\n", withCircuitBreaker)
	}
	if withRateLimit != "" {
		fmt.Printf("✓ Rate limiting enabled (%s)\n", withRateLimit)
	}
	if withChaos != "" {
		fmt.Printf("✓ Chaos engineering enabled (%s)\n", withChaos)
	}
	if withFailover != "" {
		fmt.Printf("✓ Failover enabled (%s)\n", withFailover)
	}
	if withEvent != "" {
		fmt.Printf("✓ Event sourcing enabled (%s)\n", withEvent)
	}
	if withScheduling != "" {
		fmt.Printf("✓ Task scheduling enabled (%s)\n", withScheduling)
	}
	if withBackup != "" {
		fmt.Printf("✓ Backup services enabled (%s)\n", withBackup)
	}
	if withPayment != "" {
		fmt.Printf("✓ Payment processing enabled (%s)\n", withPayment)
	}
	if withFileGen != "" {
		fmt.Printf("✓ File generation enabled (%s)\n", withFileGen)
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

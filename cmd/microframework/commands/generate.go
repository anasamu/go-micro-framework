package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/anasamu/go-micro-framework/internal/generator"
	"github.com/spf13/cobra"
)

var (
	generateType         string
	serviceName          string
	outputPath           string
	protobufPackage      string
	graphqlSchema        string
	grpcServices         []string
	graphqlTypes         []string
	graphqlQueries       []string
	graphqlMutations     []string
	graphqlSubscriptions []string
	forceGenerate        bool
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate <type>",
	Short: "Generate protobuf files or GraphQL schemas",
	Long: `Generate protobuf files for gRPC services or GraphQL schemas for GraphQL services.

This command supports:
- protobuf: Generate .proto files for gRPC services
- graphql: Generate GraphQL schema files
- service: Generate both protobuf and GraphQL for a service

Examples:
  microframework generate protobuf --service-name=user-service --grpc-services=UserService,AuthService
  microframework generate graphql --service-name=user-service --graphql-types=User,Profile --graphql-queries=getUser,getUsers
  microframework generate service --service-name=user-service --grpc-services=UserService --graphql-types=User,Profile`,
	Args: cobra.ExactArgs(1),
	RunE: runGenerate,
}

func init() {
	// Generate type
	generateCmd.Flags().StringVarP(&generateType, "type", "t", "", "Type to generate (protobuf, graphql, service)")

	// Service configuration
	generateCmd.Flags().StringVar(&serviceName, "service-name", "", "Name of the service")
	generateCmd.Flags().StringVarP(&outputPath, "output", "o", ".", "Output directory for generated files")
	generateCmd.Flags().StringVar(&protobufPackage, "protobuf-package", "", "Protobuf package name")
	generateCmd.Flags().StringVar(&graphqlSchema, "graphql-schema", "", "GraphQL schema name")

	// gRPC configuration
	generateCmd.Flags().StringSliceVar(&grpcServices, "grpc-services", []string{}, "gRPC service names (comma-separated)")

	// GraphQL configuration
	generateCmd.Flags().StringSliceVar(&graphqlTypes, "graphql-types", []string{}, "GraphQL type names (comma-separated)")
	generateCmd.Flags().StringSliceVar(&graphqlQueries, "graphql-queries", []string{}, "GraphQL query names (comma-separated)")
	generateCmd.Flags().StringSliceVar(&graphqlMutations, "graphql-mutations", []string{}, "GraphQL mutation names (comma-separated)")
	generateCmd.Flags().StringSliceVar(&graphqlSubscriptions, "graphql-subscriptions", []string{}, "GraphQL subscription names (comma-separated)")

	// Options
	generateCmd.Flags().BoolVar(&forceGenerate, "force", false, "Overwrite existing files")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	generateType = args[0]

	// Validate generate type
	if err := validateGenerateType(generateType); err != nil {
		return fmt.Errorf("invalid generate type: %w", err)
	}

	// Validate service name
	if serviceName == "" {
		return fmt.Errorf("service name is required")
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	switch generateType {
	case "protobuf":
		return generateProtobuf()
	case "graphql":
		return generateGraphQL()
	case "service":
		return generateService()
	default:
		return fmt.Errorf("unsupported generate type: %s", generateType)
	}
}

// validateGenerateType validates the generate type
func validateGenerateType(generateType string) error {
	validTypes := []string{"protobuf", "graphql", "service"}
	for _, validType := range validTypes {
		if generateType == validType {
			return nil
		}
	}
	return fmt.Errorf("generate type must be one of: %s", strings.Join(validTypes, ", "))
}

// generateProtobuf generates protobuf files for gRPC services
func generateProtobuf() error {
	fmt.Printf("Generating protobuf files for service: %s\n", serviceName)

	// Set default protobuf package if not provided
	if protobufPackage == "" {
		protobufPackage = strings.ReplaceAll(serviceName, "-", "_")
	}

	// Create protobuf generator configuration
	config := &generator.ProtobufConfig{
		ServiceName:   serviceName,
		PackageName:   protobufPackage,
		GRPCServices:  grpcServices,
		OutputPath:    outputPath,
		ForceGenerate: forceGenerate,
	}

	// Create protobuf generator
	protobufGenerator := generator.NewProtobufGenerator(config)

	// Generate protobuf files
	if err := protobufGenerator.GenerateProtobuf(); err != nil {
		return fmt.Errorf("failed to generate protobuf files: %w", err)
	}

	fmt.Printf("✓ Protobuf files generated successfully!\n")
	fmt.Printf("Generated files:\n")
	for _, service := range grpcServices {
		fmt.Printf("  - %s.proto\n", strings.ToLower(service))
	}

	return nil
}

// generateGraphQL generates GraphQL schema files
func generateGraphQL() error {
	fmt.Printf("Generating GraphQL schema for service: %s\n", serviceName)

	// Set default GraphQL schema name if not provided
	if graphqlSchema == "" {
		graphqlSchema = strings.ReplaceAll(serviceName, "-", "_")
	}

	// Create GraphQL generator configuration
	config := &generator.GraphQLConfig{
		ServiceName:   serviceName,
		SchemaName:    graphqlSchema,
		Types:         graphqlTypes,
		Queries:       graphqlQueries,
		Mutations:     graphqlMutations,
		Subscriptions: graphqlSubscriptions,
		OutputPath:    outputPath,
		ForceGenerate: forceGenerate,
	}

	// Create GraphQL generator
	graphqlGenerator := generator.NewGraphQLGenerator(config)

	// Generate GraphQL schema
	if err := graphqlGenerator.GenerateGraphQL(); err != nil {
		return fmt.Errorf("failed to generate GraphQL schema: %w", err)
	}

	fmt.Printf("✓ GraphQL schema generated successfully!\n")
	fmt.Printf("Generated files:\n")
	fmt.Printf("  - %s.graphql\n", graphqlSchema)
	fmt.Printf("  - %s_schema.go\n", graphqlSchema)

	return nil
}

// generateService generates both protobuf and GraphQL for a service
func generateService() error {
	fmt.Printf("Generating both protobuf and GraphQL for service: %s\n", serviceName)

	// Generate protobuf files
	if len(grpcServices) > 0 {
		if err := generateProtobuf(); err != nil {
			return fmt.Errorf("failed to generate protobuf files: %w", err)
		}
	}

	// Generate GraphQL schema
	if len(graphqlTypes) > 0 || len(graphqlQueries) > 0 || len(graphqlMutations) > 0 || len(graphqlSubscriptions) > 0 {
		if err := generateGraphQL(); err != nil {
			return fmt.Errorf("failed to generate GraphQL schema: %w", err)
		}
	}

	fmt.Printf("✓ Service generation completed successfully!\n")

	return nil
}

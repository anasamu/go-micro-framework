package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// GraphQLConfig holds configuration for GraphQL generation
type GraphQLConfig struct {
	ServiceName   string
	SchemaName    string
	Types         []string
	Queries       []string
	Mutations     []string
	Subscriptions []string
	OutputPath    string
	ForceGenerate bool
}

// GraphQLGenerator handles the generation of GraphQL schema files
type GraphQLGenerator struct {
	config *GraphQLConfig
}

// NewGraphQLGenerator creates a new GraphQL generator
func NewGraphQLGenerator(config *GraphQLConfig) *GraphQLGenerator {
	return &GraphQLGenerator{
		config: config,
	}
}

// GenerateGraphQL generates GraphQL schema files
func (gg *GraphQLGenerator) GenerateGraphQL() error {
	// Create GraphQL directory
	graphqlDir := filepath.Join(gg.config.OutputPath, "graphql")
	if err := os.MkdirAll(graphqlDir, 0755); err != nil {
		return fmt.Errorf("failed to create GraphQL directory: %w", err)
	}

	// Generate GraphQL schema file
	if err := gg.generateGraphQLSchema(graphqlDir); err != nil {
		return fmt.Errorf("failed to generate GraphQL schema: %w", err)
	}

	// Generate Go schema file
	if err := gg.generateGoSchema(graphqlDir); err != nil {
		return fmt.Errorf("failed to generate Go schema: %w", err)
	}

	return nil
}

// generateGraphQLSchema generates the GraphQL schema file
func (gg *GraphQLGenerator) generateGraphQLSchema(graphqlDir string) error {
	// Create GraphQL schema file
	fileName := gg.config.SchemaName + ".graphql"
	filePath := filepath.Join(graphqlDir, fileName)

	// Check if file exists and force is not set
	if !gg.config.ForceGenerate {
		if _, err := os.Stat(filePath); err == nil {
			return fmt.Errorf("file %s already exists, use --force to overwrite", fileName)
		}
	}

	// Parse template
	tmpl, err := template.New("schema.graphql").Parse(graphQLSchemaTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse GraphQL schema template: %w", err)
	}

	// Create template data
	data := map[string]interface{}{
		"ServiceName":   gg.config.ServiceName,
		"SchemaName":    gg.config.SchemaName,
		"Types":         gg.config.Types,
		"Queries":       gg.config.Queries,
		"Mutations":     gg.config.Mutations,
		"Subscriptions": gg.config.Subscriptions,
	}

	// Write file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// generateGoSchema generates the Go schema file
func (gg *GraphQLGenerator) generateGoSchema(graphqlDir string) error {
	// Create Go schema file
	fileName := gg.config.SchemaName + "_schema.go"
	filePath := filepath.Join(graphqlDir, fileName)

	// Check if file exists and force is not set
	if !gg.config.ForceGenerate {
		if _, err := os.Stat(filePath); err == nil {
			return fmt.Errorf("file %s already exists, use --force to overwrite", fileName)
		}
	}

	// Parse template
	tmpl, err := template.New("schema.go").Parse(graphQLGoSchemaTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse GraphQL Go schema template: %w", err)
	}

	// Create template data
	data := map[string]interface{}{
		"ServiceName":   gg.config.ServiceName,
		"SchemaName":    gg.config.SchemaName,
		"Types":         gg.config.Types,
		"Queries":       gg.config.Queries,
		"Mutations":     gg.config.Mutations,
		"Subscriptions": gg.config.Subscriptions,
	}

	// Write file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// GraphQL schema template
const graphQLSchemaTemplate = `# {{.ServiceName}} GraphQL Schema

# Scalar types
scalar Time
scalar JSON

# Input types
input PaginationInput {
  page: Int = 1
  limit: Int = 10
}

input SortInput {
  field: String!
  order: SortOrder = ASC
}

enum SortOrder {
  ASC
  DESC
}

# Response types
type Response {
  success: Boolean!
  message: String
  errors: [String!]
}

type PaginationInfo {
  page: Int!
  limit: Int!
  total: Int!
  totalPages: Int!
}

# Health check
type Health {
  status: String!
  message: String!
  timestamp: Time!
}

# Service info
type ServiceInfo {
  name: String!
  version: String!
  description: String!
  startedAt: Time!
}

{{range .Types}}
# {{.}} type
type {{.}} {
  id: ID!
  name: String!
  description: String
  createdAt: Time!
  updatedAt: Time!
}

# {{.}} input
input {{.}}Input {
  name: String!
  description: String
}

# {{.}} update input
input {{.}}UpdateInput {
  id: ID!
  name: String
  description: String
}

# {{.}} list response
type {{.}}ListResponse {
  {{. | lower}}s: [{{.}}!]!
  pagination: PaginationInfo!
  response: Response!
}
{{end}}

# Query type
type Query {
  # Health check
  health: Health!
  
  # Service info
  serviceInfo: ServiceInfo!
  
{{range .Queries}}
  # {{.}} query
  {{. | lower}}: {{.}}!
{{end}}

{{range .Types}}
  # {{.}} queries
  {{. | lower}}(id: ID!): {{.}}!
  {{. | lower}}s(pagination: PaginationInput, sort: SortInput, search: String): {{.}}ListResponse!
{{end}}
}

# Mutation type
type Mutation {
{{range .Mutations}}
  # {{.}} mutation
  {{. | lower}}: {{.}}!
{{end}}

{{range .Types}}
  # {{.}} mutations
  create{{.}}(input: {{.}}Input!): {{.}}!
  update{{.}}(input: {{.}}UpdateInput!): {{.}}!
  delete{{.}}(id: ID!): Response!
{{end}}
}

# Subscription type
type Subscription {
{{range .Subscriptions}}
  # {{.}} subscription
  {{. | lower}}: {{.}}!
{{end}}

{{range .Types}}
  # {{.}} subscriptions
  {{. | lower}}Created: {{.}}!
  {{. | lower}}Updated: {{.}}!
  {{. | lower}}Deleted: ID!
{{end}}
}
`

// GraphQL Go schema template
const graphQLGoSchemaTemplate = `package graphql

import (
	"context"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/anasamu/go-micro-libs/communication/providers/graphql"
)

// {{.SchemaName}}Schema represents the GraphQL schema for {{.ServiceName}}
type {{.SchemaName}}Schema struct {
	schema *graphql.Schema
	provider *graphql.Provider
}

// New{{.SchemaName}}Schema creates a new GraphQL schema
func New{{.SchemaName}}Schema(provider *graphql.Provider) (*{{.SchemaName}}Schema, error) {
	schema, err := create{{.SchemaName}}Schema()
	if err != nil {
		return nil, err
	}

	return &{{.SchemaName}}Schema{
		schema:   schema,
		provider: provider,
	}, nil
}

// GetSchema returns the GraphQL schema
func (s *{{.SchemaName}}Schema) GetSchema() *graphql.Schema {
	return s.schema
}

// create{{.SchemaName}}Schema creates the GraphQL schema
func create{{.SchemaName}}Schema() (*graphql.Schema, error) {
	// Define scalar types
	timeType := graphql.NewScalar(graphql.ScalarConfig{
		Name:        "Time",
		Description: "Time scalar type",
		Serialize: func(value interface{}) interface{} {
			if t, ok := value.(time.Time); ok {
				return t.Format(time.RFC3339)
			}
			return nil
		},
		ParseValue: func(value interface{}) interface{} {
			if str, ok := value.(string); ok {
				if t, err := time.Parse(time.RFC3339, str); err == nil {
					return t
				}
			}
			return nil
		},
	})

	jsonType := graphql.NewScalar(graphql.ScalarConfig{
		Name:        "JSON",
		Description: "JSON scalar type",
		Serialize: func(value interface{}) interface{} {
			return value
		},
		ParseValue: func(value interface{}) interface{} {
			return value
		},
	})

	// Define enums
	sortOrderEnum := graphql.NewEnum(graphql.EnumConfig{
		Name:        "SortOrder",
		Description: "Sort order enum",
		Values: graphql.EnumValueConfigMap{
			"ASC": &graphql.EnumValueConfig{
				Value: "ASC",
			},
			"DESC": &graphql.EnumValueConfig{
				Value: "DESC",
			},
		},
	})

	// Define input types
	paginationInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "PaginationInput",
		Description: "Pagination input",
		Fields: graphql.InputObjectConfigFieldMap{
			"page": &graphql.InputObjectFieldConfig{
				Type:        graphql.Int,
				Description: "Page number",
			},
			"limit": &graphql.InputObjectFieldConfig{
				Type:        graphql.Int,
				Description: "Items per page",
			},
		},
	})

	sortInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "SortInput",
		Description: "Sort input",
		Fields: graphql.InputObjectConfigFieldMap{
			"field": &graphql.InputObjectFieldConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Field to sort by",
			},
			"order": &graphql.InputObjectFieldConfig{
				Type:        sortOrderEnum,
				Description: "Sort order",
			},
		},
	})

	// Define response types
	responseType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Response",
		Description: "Response type",
		Fields: graphql.Fields{
			"success": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Boolean),
				Description: "Success status",
			},
			"message": &graphql.Field{
				Type:        graphql.String,
				Description: "Response message",
			},
			"errors": &graphql.Field{
				Type:        graphql.NewList(graphql.NewNonNull(graphql.String)),
				Description: "Error messages",
			},
		},
	})

	paginationInfoType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "PaginationInfo",
		Description: "Pagination info",
		Fields: graphql.Fields{
			"page": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Current page",
			},
			"limit": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Items per page",
			},
			"total": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Total items",
			},
			"totalPages": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Total pages",
			},
		},
	})

	// Define health type
	healthType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Health",
		Description: "Health check type",
		Fields: graphql.Fields{
			"status": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Health status",
			},
			"message": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Health message",
			},
			"timestamp": &graphql.Field{
				Type:        graphql.NewNonNull(timeType),
				Description: "Health check timestamp",
			},
		},
	})

	// Define service info type
	serviceInfoType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "ServiceInfo",
		Description: "Service info type",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Service name",
			},
			"version": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Service version",
			},
			"description": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Service description",
			},
			"startedAt": &graphql.Field{
				Type:        graphql.NewNonNull(timeType),
				Description: "Service start time",
			},
		},
	})

{{range .Types}}
	// Define {{.}} type
	{{. | lower}}Type := graphql.NewObject(graphql.ObjectConfig{
		Name:        "{{.}}",
		Description: "{{.}} type",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "{{.}} ID",
			},
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "{{.}} name",
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "{{.}} description",
			},
			"createdAt": &graphql.Field{
				Type:        graphql.NewNonNull(timeType),
				Description: "Creation timestamp",
			},
			"updatedAt": &graphql.Field{
				Type:        graphql.NewNonNull(timeType),
				Description: "Last update timestamp",
			},
		},
	})

	// Define {{.}} input type
	{{. | lower}}InputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "{{.}}Input",
		Description: "{{.}} input",
		Fields: graphql.InputObjectConfigFieldMap{
			"name": &graphql.InputObjectFieldConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "{{.}} name",
			},
			"description": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "{{.}} description",
			},
		},
	})

	// Define {{.}} update input type
	{{. | lower}}UpdateInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "{{.}}UpdateInput",
		Description: "{{.}} update input",
		Fields: graphql.InputObjectConfigFieldMap{
			"id": &graphql.InputObjectFieldConfig{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "{{.}} ID",
			},
			"name": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "{{.}} name",
			},
			"description": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "{{.}} description",
			},
		},
	})

	// Define {{.}} list response type
	{{. | lower}}ListResponseType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "{{.}}ListResponse",
		Description: "{{.}} list response",
		Fields: graphql.Fields{
			"{{. | lower}}s": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull({{. | lower}}Type))),
				Description: "List of {{.}}s",
			},
			"pagination": &graphql.Field{
				Type:        graphql.NewNonNull(paginationInfoType),
				Description: "Pagination info",
			},
			"response": &graphql.Field{
				Type:        graphql.NewNonNull(responseType),
				Description: "Response info",
			},
		},
	})
{{end}}

	// Define query type
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Query",
		Description: "Query type",
		Fields: graphql.Fields{
			"health": &graphql.Field{
				Type:        graphql.NewNonNull(healthType),
				Description: "Health check",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{
						"status":    "healthy",
						"message":   "Service is running",
						"timestamp": time.Now(),
					}, nil
				},
			},
			"serviceInfo": &graphql.Field{
				Type:        graphql.NewNonNull(serviceInfoType),
				Description: "Service information",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{
						"name":        "{{.ServiceName}}",
						"version":     "1.0.0",
						"description": "{{.ServiceName}} microservice",
						"startedAt":   time.Now(),
					}, nil
				},
			},
{{range .Queries}}
			"{{. | lower}}": &graphql.Field{
				Type:        graphql.NewNonNull({{. | lower}}Type),
				Description: "{{.}} query",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Implement {{.}} query logic
					return map[string]interface{}{
						"id":          "1",
						"name":        "{{.}}",
						"description": "{{.}} description",
						"createdAt":   time.Now(),
						"updatedAt":   time.Now(),
					}, nil
				},
			},
{{end}}
{{range .Types}}
			"{{. | lower}}": &graphql.Field{
				Type:        graphql.NewNonNull({{. | lower}}Type),
				Description: "Get {{.}} by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.ID),
						Description: "{{.}} ID",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Implement {{.}} query logic
					return map[string]interface{}{
						"id":          p.Args["id"],
						"name":        "{{.}}",
						"description": "{{.}} description",
						"createdAt":   time.Now(),
						"updatedAt":   time.Now(),
					}, nil
				},
			},
			"{{. | lower}}s": &graphql.Field{
				Type:        graphql.NewNonNull({{. | lower}}ListResponseType),
				Description: "List {{.}}s",
				Args: graphql.FieldConfigArgument{
					"pagination": &graphql.ArgumentConfig{
						Type:        paginationInputType,
						Description: "Pagination input",
					},
					"sort": &graphql.ArgumentConfig{
						Type:        sortInputType,
						Description: "Sort input",
					},
					"search": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Search term",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Implement {{.}} list query logic
					return map[string]interface{}{
						"{{. | lower}}s": []interface{}{
							map[string]interface{}{
								"id":          "1",
								"name":        "{{.}} 1",
								"description": "{{.}} 1 description",
								"createdAt":   time.Now(),
								"updatedAt":   time.Now(),
							},
						},
						"pagination": map[string]interface{}{
							"page":       1,
							"limit":      10,
							"total":      1,
							"totalPages": 1,
						},
						"response": map[string]interface{}{
							"success": true,
							"message": "{{.}}s retrieved successfully",
							"errors":  []string{},
						},
					}, nil
				},
			},
{{end}}
		},
	})

	// Define mutation type
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Mutation",
		Description: "Mutation type",
		Fields: graphql.Fields{
{{range .Mutations}}
			"{{. | lower}}": &graphql.Field{
				Type:        graphql.NewNonNull({{. | lower}}Type),
				Description: "{{.}} mutation",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Implement {{.}} mutation logic
					return map[string]interface{}{
						"id":          "1",
						"name":        "{{.}}",
						"description": "{{.}} description",
						"createdAt":   time.Now(),
						"updatedAt":   time.Now(),
					}, nil
				},
			},
{{end}}
{{range .Types}}
			"create{{.}}": &graphql.Field{
				Type:        graphql.NewNonNull({{. | lower}}Type),
				Description: "Create {{.}}",
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull({{. | lower}}InputType),
						Description: "{{.}} input",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Implement create {{.}} logic
					return map[string]interface{}{
						"id":          "1",
						"name":        "{{.}}",
						"description": "{{.}} description",
						"createdAt":   time.Now(),
						"updatedAt":   time.Now(),
					}, nil
				},
			},
			"update{{.}}": &graphql.Field{
				Type:        graphql.NewNonNull({{. | lower}}Type),
				Description: "Update {{.}}",
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull({{. | lower}}UpdateInputType),
						Description: "{{.}} update input",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Implement update {{.}} logic
					return map[string]interface{}{
						"id":          "1",
						"name":        "{{.}}",
						"description": "{{.}} description",
						"createdAt":   time.Now(),
						"updatedAt":   time.Now(),
					}, nil
				},
			},
			"delete{{.}}": &graphql.Field{
				Type:        graphql.NewNonNull(responseType),
				Description: "Delete {{.}}",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.ID),
						Description: "{{.}} ID",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Implement delete {{.}} logic
					return map[string]interface{}{
						"success": true,
						"message": "{{.}} deleted successfully",
						"errors":  []string{},
					}, nil
				},
			},
{{end}}
		},
	})

	// Define subscription type
	subscriptionType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Subscription",
		Description: "Subscription type",
		Fields: graphql.Fields{
{{range .Subscriptions}}
			"{{. | lower}}": &graphql.Field{
				Type:        graphql.NewNonNull({{. | lower}}Type),
				Description: "{{.}} subscription",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Implement {{.}} subscription logic
					return map[string]interface{}{
						"id":          "1",
						"name":        "{{.}}",
						"description": "{{.}} description",
						"createdAt":   time.Now(),
						"updatedAt":   time.Now(),
					}, nil
				},
			},
{{end}}
{{range .Types}}
			"{{. | lower}}Created": &graphql.Field{
				Type:        graphql.NewNonNull({{. | lower}}Type),
				Description: "{{.}} created subscription",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Implement {{.}} created subscription logic
					return map[string]interface{}{
						"id":          "1",
						"name":        "{{.}}",
						"description": "{{.}} description",
						"createdAt":   time.Now(),
						"updatedAt":   time.Now(),
					}, nil
				},
			},
			"{{. | lower}}Updated": &graphql.Field{
				Type:        graphql.NewNonNull({{. | lower}}Type),
				Description: "{{.}} updated subscription",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Implement {{.}} updated subscription logic
					return map[string]interface{}{
						"id":          "1",
						"name":        "{{.}}",
						"description": "{{.}} description",
						"createdAt":   time.Now(),
						"updatedAt":   time.Now(),
					}, nil
				},
			},
			"{{. | lower}}Deleted": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "{{.}} deleted subscription",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Implement {{.}} deleted subscription logic
					return "1", nil
				},
			},
{{end}}
		},
	})

	// Create schema
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:        queryType,
		Mutation:     mutationType,
		Subscription: subscriptionType,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create GraphQL schema: %w", err)
	}

	return &schema, nil
}
`

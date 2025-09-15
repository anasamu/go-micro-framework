package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// ProtobufConfig holds configuration for protobuf generation
type ProtobufConfig struct {
	ServiceName   string
	PackageName   string
	GRPCServices  []string
	OutputPath    string
	ForceGenerate bool
}

// ProtobufGenerator handles the generation of protobuf files
type ProtobufGenerator struct {
	config *ProtobufConfig
}

// NewProtobufGenerator creates a new protobuf generator
func NewProtobufGenerator(config *ProtobufConfig) *ProtobufGenerator {
	return &ProtobufGenerator{
		config: config,
	}
}

// GenerateProtobuf generates protobuf files for gRPC services
func (pg *ProtobufGenerator) GenerateProtobuf() error {
	// Create protobuf directory
	protobufDir := filepath.Join(pg.config.OutputPath, "protobuf")
	if err := os.MkdirAll(protobufDir, 0755); err != nil {
		return fmt.Errorf("failed to create protobuf directory: %w", err)
	}

	// Generate protobuf files for each service
	for _, serviceName := range pg.config.GRPCServices {
		if err := pg.generateServiceProtobuf(serviceName, protobufDir); err != nil {
			return fmt.Errorf("failed to generate protobuf for service %s: %w", serviceName, err)
		}
	}

	// Generate main protobuf file
	if err := pg.generateMainProtobuf(protobufDir); err != nil {
		return fmt.Errorf("failed to generate main protobuf file: %w", err)
	}

	return nil
}

// generateServiceProtobuf generates a protobuf file for a specific service
func (pg *ProtobufGenerator) generateServiceProtobuf(serviceName, protobufDir string) error {
	// Create service-specific protobuf file
	fileName := strings.ToLower(serviceName) + ".proto"
	filePath := filepath.Join(protobufDir, fileName)

	// Check if file exists and force is not set
	if !pg.config.ForceGenerate {
		if _, err := os.Stat(filePath); err == nil {
			return fmt.Errorf("file %s already exists, use --force to overwrite", fileName)
		}
	}

	// Parse template
	tmpl, err := template.New("service.proto").Parse(serviceProtobufTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse service protobuf template: %w", err)
	}

	// Create template data
	data := map[string]interface{}{
		"ServiceName":      serviceName,
		"PackageName":      pg.config.PackageName,
		"ServiceNameLower": strings.ToLower(serviceName),
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

// generateMainProtobuf generates the main protobuf file
func (pg *ProtobufGenerator) generateMainProtobuf(protobufDir string) error {
	// Create main protobuf file
	fileName := pg.config.ServiceName + ".proto"
	filePath := filepath.Join(protobufDir, fileName)

	// Check if file exists and force is not set
	if !pg.config.ForceGenerate {
		if _, err := os.Stat(filePath); err == nil {
			return fmt.Errorf("file %s already exists, use --force to overwrite", fileName)
		}
	}

	// Parse template
	tmpl, err := template.New("main.proto").Parse(mainProtobufTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse main protobuf template: %w", err)
	}

	// Create template data
	data := map[string]interface{}{
		"ServiceName":  pg.config.ServiceName,
		"PackageName":  pg.config.PackageName,
		"GRPCServices": pg.config.GRPCServices,
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

// Service protobuf template
const serviceProtobufTemplate = `syntax = "proto3";

package {{.PackageName}};

option go_package = "github.com/anasamu/{{.ServiceNameLower}}/protobuf";

import "google/protobuf/empty.proto";
import "google/protobuf/timestamp.proto";

// {{.ServiceName}} service definition
service {{.ServiceName}} {
  // Health check
  rpc HealthCheck(google.protobuf.Empty) returns (HealthResponse);
  
  // Get service info
  rpc GetServiceInfo(google.protobuf.Empty) returns (ServiceInfoResponse);
  
  // Example CRUD operations
  rpc Create{{.ServiceName}}(Create{{.ServiceName}}Request) returns ({{.ServiceName}}Response);
  rpc Get{{.ServiceName}}(Get{{.ServiceName}}Request) returns ({{.ServiceName}}Response);
  rpc Update{{.ServiceName}}(Update{{.ServiceName}}Request) returns ({{.ServiceName}}Response);
  rpc Delete{{.ServiceName}}(Delete{{.ServiceName}}Request) returns (google.protobuf.Empty);
  rpc List{{.ServiceName}}s(List{{.ServiceName}}sRequest) returns (List{{.ServiceName}}sResponse);
}

// Health response
message HealthResponse {
  string status = 1;
  string message = 2;
  google.protobuf.Timestamp timestamp = 3;
}

// Service info response
message ServiceInfoResponse {
  string name = 1;
  string version = 2;
  string description = 3;
  google.protobuf.Timestamp started_at = 4;
}

// {{.ServiceName}} entity
message {{.ServiceName}} {
  string id = 1;
  string name = 2;
  string description = 3;
  google.protobuf.Timestamp created_at = 4;
  google.protobuf.Timestamp updated_at = 5;
}

// Create {{.ServiceName}} request
message Create{{.ServiceName}}Request {
  string name = 1;
  string description = 2;
}

// Get {{.ServiceName}} request
message Get{{.ServiceName}}Request {
  string id = 1;
}

// Update {{.ServiceName}} request
message Update{{.ServiceName}}Request {
  string id = 1;
  string name = 2;
  string description = 3;
}

// Delete {{.ServiceName}} request
message Delete{{.ServiceName}}Request {
  string id = 1;
}

// List {{.ServiceName}}s request
message List{{.ServiceName}}sRequest {
  int32 page = 1;
  int32 limit = 2;
  string search = 3;
}

// {{.ServiceName}} response
message {{.ServiceName}}Response {
  {{.ServiceName}} {{.ServiceNameLower}} = 1;
  string message = 2;
  bool success = 3;
}

// List {{.ServiceName}}s response
message List{{.ServiceName}}sResponse {
  repeated {{.ServiceName}} {{.ServiceNameLower}}s = 1;
  int32 total = 2;
  int32 page = 3;
  int32 limit = 4;
  string message = 5;
  bool success = 6;
}
`

// Main protobuf template
const mainProtobufTemplate = `syntax = "proto3";

package {{.PackageName}};

option go_package = "github.com/anasamu/{{.ServiceName}}/protobuf";

import "google/protobuf/empty.proto";
import "google/protobuf/timestamp.proto";

// Main service definition
service {{.ServiceName}}Service {
  // Health check
  rpc HealthCheck(google.protobuf.Empty) returns (HealthResponse);
  
  // Get service info
  rpc GetServiceInfo(google.protobuf.Empty) returns (ServiceInfoResponse);
}

// Health response
message HealthResponse {
  string status = 1;
  string message = 2;
  google.protobuf.Timestamp timestamp = 3;
}

// Service info response
message ServiceInfoResponse {
  string name = 1;
  string version = 2;
  string description = 3;
  google.protobuf.Timestamp started_at = 4;
}

{{range .GRPCServices}}
// {{.}} service definition
service {{.}} {
  // Health check
  rpc HealthCheck(google.protobuf.Empty) returns (HealthResponse);
  
  // Get service info
  rpc GetServiceInfo(google.protobuf.Empty) returns (ServiceInfoResponse);
  
  // Example CRUD operations
  rpc Create{{.}}(Create{{.}}Request) returns ({{.}}Response);
  rpc Get{{.}}(Get{{.}}Request) returns ({{.}}Response);
  rpc Update{{.}}(Update{{.}}Request) returns ({{.}}Response);
  rpc Delete{{.}}(Delete{{.}}Request) returns (google.protobuf.Empty);
  rpc List{{.}}s(List{{.}}sRequest) returns (List{{.}}sResponse);
}

// {{.}} entity
message {{.}} {
  string id = 1;
  string name = 2;
  string description = 3;
  google.protobuf.Timestamp created_at = 4;
  google.protobuf.Timestamp updated_at = 5;
}

// Create {{.}} request
message Create{{.}}Request {
  string name = 1;
  string description = 2;
}

// Get {{.}} request
message Get{{.}}Request {
  string id = 1;
}

// Update {{.}} request
message Update{{.}}Request {
  string id = 1;
  string name = 2;
  string description = 3;
}

// Delete {{.}} request
message Delete{{.}}Request {
  string id = 1;
}

// List {{.}}s request
message List{{.}}sRequest {
  int32 page = 1;
  int32 limit = 2;
  string search = 3;
}

// {{.}} response
message {{.}}Response {
  {{.}} {{. | lower}} = 1;
  string message = 2;
  bool success = 3;
}

// List {{.}}s response
message List{{.}}sResponse {
  repeated {{.}} {{. | lower}}s = 1;
  int32 total = 2;
  int32 page = 3;
  int32 limit = 4;
  string message = 5;
  bool success = 6;
}
{{end}}
`

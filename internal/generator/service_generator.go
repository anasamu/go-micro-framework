package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/anasamu/go-micro-framework/internal/templates"
)

// ServiceGenerator handles the generation of microservice projects
type ServiceGenerator struct {
	templates map[string]*template.Template
	config    *GeneratorConfig
}

// GeneratorConfig holds configuration for service generation
type GeneratorConfig struct {
	ServiceName        string
	ServiceType        string
	WithAuth           bool
	WithDatabase       bool
	WithMessaging      bool
	WithMonitoring     bool
	WithAI             bool
	WithStorage        bool
	WithCache          bool
	WithDiscovery      bool
	WithCircuitBreaker bool
	WithRateLimit      bool
	WithChaos          bool
	WithFailover       bool
	WithEvent          bool
	WithScheduling     bool
	WithBackup         bool
	WithPayment        bool
	WithFileGen        bool
	OutputDir          string
	// Provider specifications
	AuthProvider       string
	DatabaseProvider   string
	MessagingProvider  string
	MonitoringProvider string
	AIProvider         string
	StorageProvider    string
	CacheProvider      string
	DiscoveryProvider  string
	PaymentProvider    string
}

// NewServiceGenerator creates a new service generator
func NewServiceGenerator(config *GeneratorConfig) *ServiceGenerator {
	return &ServiceGenerator{
		templates: make(map[string]*template.Template),
		config:    config,
	}
}

// GenerateService generates a complete microservice project
func (sg *ServiceGenerator) GenerateService() error {
	// Create project directory structure
	if err := sg.createProjectStructure(); err != nil {
		return fmt.Errorf("failed to create project structure: %w", err)
	}

	// Generate main.go
	if err := sg.generateMain(); err != nil {
		return fmt.Errorf("failed to generate main.go: %w", err)
	}

	// Generate go.mod
	if err := sg.generateGoMod(); err != nil {
		return fmt.Errorf("failed to generate go.mod: %w", err)
	}

	// Generate configuration files
	if err := sg.generateConfig(); err != nil {
		return fmt.Errorf("failed to generate configuration: %w", err)
	}

	// Generate handlers
	if err := sg.generateHandlers(); err != nil {
		return fmt.Errorf("failed to generate handlers: %w", err)
	}

	// Generate models
	if err := sg.generateModels(); err != nil {
		return fmt.Errorf("failed to generate models: %w", err)
	}

	// Generate repositories
	if err := sg.generateRepositories(); err != nil {
		return fmt.Errorf("failed to generate repositories: %w", err)
	}

	// Generate services
	if err := sg.generateServices(); err != nil {
		return fmt.Errorf("failed to generate services: %w", err)
	}

	// Generate middleware
	if err := sg.generateMiddleware(); err != nil {
		return fmt.Errorf("failed to generate middleware: %w", err)
	}

	// Generate Docker files
	if err := sg.generateDocker(); err != nil {
		return fmt.Errorf("failed to generate Docker files: %w", err)
	}

	// Generate Kubernetes manifests
	if err := sg.generateKubernetes(); err != nil {
		return fmt.Errorf("failed to generate Kubernetes manifests: %w", err)
	}

	// Generate tests
	if err := sg.generateTests(); err != nil {
		return fmt.Errorf("failed to generate tests: %w", err)
	}

	// Generate documentation
	if err := sg.generateDocumentation(); err != nil {
		return fmt.Errorf("failed to generate documentation: %w", err)
	}

	return nil
}

// createProjectStructure creates the directory structure for the service
func (sg *ServiceGenerator) createProjectStructure() error {
	baseDir := filepath.Join(sg.config.OutputDir, sg.config.ServiceName)

	dirs := []string{
		"cmd",
		"internal/handlers",
		"internal/models",
		"internal/repositories",
		"internal/services",
		"internal/middleware",
		"pkg/types",
		"configs",
		"deployments/docker",
		"deployments/kubernetes",
		"deployments/helm",
		"tests/unit",
		"tests/integration",
		"tests/e2e",
		"docs",
		"scripts",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(baseDir, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", fullPath, err)
		}
	}

	return nil
}

// generateMain generates the main.go file
func (sg *ServiceGenerator) generateMain() error {
	tmpl, err := template.New("main.go").Parse(templates.MainTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse main template: %w", err)
	}

	outputPath := filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "cmd", "main.go")
	return sg.writeTemplate(tmpl, outputPath, sg.config)
}

// generateGoMod generates the go.mod file
func (sg *ServiceGenerator) generateGoMod() error {
	tmpl, err := template.New("go.mod").Parse(templates.GoModTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse go.mod template: %w", err)
	}

	outputPath := filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "go.mod")
	return sg.writeTemplate(tmpl, outputPath, sg.config)
}

// generateConfig generates configuration files
func (sg *ServiceGenerator) generateConfig() error {
	// Generate config.yaml
	tmpl, err := template.New("config.yaml").Parse(templates.ConfigTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse config template: %w", err)
	}

	outputPath := filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "configs", "config.yaml")
	if err := sg.writeTemplate(tmpl, outputPath, sg.config); err != nil {
		return err
	}

	// Generate config.dev.yaml
	tmpl, err = template.New("config.dev.yaml").Parse(templates.ConfigDevTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse config.dev template: %w", err)
	}

	outputPath = filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "configs", "config.dev.yaml")
	return sg.writeTemplate(tmpl, outputPath, sg.config)
}

// generateHandlers generates HTTP handlers
func (sg *ServiceGenerator) generateHandlers() error {
	tmpl, err := template.New("handlers.go").Parse(templates.HandlersTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse handlers template: %w", err)
	}

	outputPath := filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "internal", "handlers", "handlers.go")
	return sg.writeTemplate(tmpl, outputPath, sg.config)
}

// generateModels generates data models
func (sg *ServiceGenerator) generateModels() error {
	tmpl, err := template.New("models.go").Parse(templates.ModelsTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse models template: %w", err)
	}

	outputPath := filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "internal", "models", "models.go")
	return sg.writeTemplate(tmpl, outputPath, sg.config)
}

// generateRepositories generates data repositories
func (sg *ServiceGenerator) generateRepositories() error {
	tmpl, err := template.New("repositories.go").Parse(templates.RepositoriesTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse repositories template: %w", err)
	}

	outputPath := filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "internal", "repositories", "repositories.go")
	return sg.writeTemplate(tmpl, outputPath, sg.config)
}

// generateServices generates business logic services
func (sg *ServiceGenerator) generateServices() error {
	tmpl, err := template.New("services.go").Parse(templates.ServicesTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse services template: %w", err)
	}

	outputPath := filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "internal", "services", "services.go")
	return sg.writeTemplate(tmpl, outputPath, sg.config)
}

// generateMiddleware generates middleware components
func (sg *ServiceGenerator) generateMiddleware() error {
	tmpl, err := template.New("middleware.go").Parse(templates.MiddlewareTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse middleware template: %w", err)
	}

	outputPath := filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "internal", "middleware", "middleware.go")
	return sg.writeTemplate(tmpl, outputPath, sg.config)
}

// generateDocker generates Docker-related files
func (sg *ServiceGenerator) generateDocker() error {
	// Generate Dockerfile
	tmpl, err := template.New("Dockerfile").Parse(templates.DockerfileTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse Dockerfile template: %w", err)
	}

	outputPath := filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "deployments", "docker", "Dockerfile")
	if err := sg.writeTemplate(tmpl, outputPath, sg.config); err != nil {
		return err
	}

	// Generate docker-compose.yml
	tmpl, err = template.New("docker-compose.yml").Parse(templates.DockerComposeTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse docker-compose template: %w", err)
	}

	outputPath = filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "deployments", "docker", "docker-compose.yml")
	return sg.writeTemplate(tmpl, outputPath, sg.config)
}

// generateKubernetes generates Kubernetes manifests
func (sg *ServiceGenerator) generateKubernetes() error {
	// Generate deployment.yaml
	tmpl, err := template.New("deployment.yaml").Parse(templates.KubernetesDeploymentTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse deployment template: %w", err)
	}

	outputPath := filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "deployments", "kubernetes", "deployment.yaml")
	if err := sg.writeTemplate(tmpl, outputPath, sg.config); err != nil {
		return err
	}

	// Generate service.yaml
	tmpl, err = template.New("service.yaml").Parse(templates.KubernetesServiceTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse service template: %w", err)
	}

	outputPath = filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "deployments", "kubernetes", "service.yaml")
	if err := sg.writeTemplate(tmpl, outputPath, sg.config); err != nil {
		return err
	}

	// Generate configmap.yaml
	tmpl, err = template.New("configmap.yaml").Parse(templates.KubernetesConfigMapTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse configmap template: %w", err)
	}

	outputPath = filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "deployments", "kubernetes", "configmap.yaml")
	return sg.writeTemplate(tmpl, outputPath, sg.config)
}

// generateTests generates test files
func (sg *ServiceGenerator) generateTests() error {
	// Generate unit tests
	tmpl, err := template.New("unit_test.go").Parse(templates.UnitTestTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse unit test template: %w", err)
	}

	outputPath := filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "tests", "unit", "service_test.go")
	if err := sg.writeTemplate(tmpl, outputPath, sg.config); err != nil {
		return err
	}

	// Generate integration tests
	tmpl, err = template.New("integration_test.go").Parse(templates.IntegrationTestTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse integration test template: %w", err)
	}

	outputPath = filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "tests", "integration", "integration_test.go")
	return sg.writeTemplate(tmpl, outputPath, sg.config)
}

// generateDocumentation generates documentation files
func (sg *ServiceGenerator) generateDocumentation() error {
	// Generate README.md
	tmpl, err := template.New("README.md").Parse(templates.ReadmeTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse README template: %w", err)
	}

	outputPath := filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "README.md")
	if err := sg.writeTemplate(tmpl, outputPath, sg.config); err != nil {
		return err
	}

	// Generate API documentation
	tmpl, err = template.New("API.md").Parse(templates.APITemplate)
	if err != nil {
		return fmt.Errorf("failed to parse API template: %w", err)
	}

	outputPath = filepath.Join(sg.config.OutputDir, sg.config.ServiceName, "docs", "API.md")
	return sg.writeTemplate(tmpl, outputPath, sg.config)
}

// writeTemplate writes a template to a file
func (sg *ServiceGenerator) writeTemplate(tmpl *template.Template, outputPath string, data interface{}) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

package templates

// Template constants for service generation
const (
	MainTemplate = `package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"
	
	// Use go-micro-libs library
	microservices "github.com/anasamu/go-micro-libs"
)

func main() {
	ctx := context.Background()
	
	// Initialize using go-micro-libs library
	configManager := microservices.NewConfigManager()
	loggingManager := microservices.NewLoggingManager(
		microservices.DefaultLoggingManagerConfig(),
		logger,
	)
	monitoringManager := microservices.NewMonitoringManager(
		microservices.DefaultMonitoringManagerConfig(),
		logger,
	)
	databaseManager := microservices.NewDatabaseManager(
		microservices.DefaultDatabaseManagerConfig(),
		logger,
	)
	authManager := microservices.NewAuthManager(
		microservices.DefaultAuthManagerConfig(),
		logger,
	)
	middlewareManager := microservices.NewMiddlewareManager(
		microservices.DefaultMiddlewareManagerConfig(),
		logger,
	)
	communicationManager := microservices.NewCommunicationManager(
		microservices.DefaultCommunicationManagerConfig(),
		logger,
	)
	
	// Bootstrap service using go-micro-libs library
	if err := bootstrapService(ctx, configManager, loggingManager, monitoringManager, 
		databaseManager, authManager, middlewareManager, communicationManager); err != nil {
		log.Fatal("Failed to bootstrap service:", err)
	}
	
	log.Println("Service started successfully")
}

func bootstrapService(ctx context.Context, 
	configManager *config_gateway.ConfigManager,
	loggingManager *logging_gateway.LoggingManager,
	monitoringManager *monitoring_gateway.MonitoringManager,
	databaseManager *database_gateway.DatabaseManager,
	authManager *auth_gateway.AuthManager,
	middlewareManager *middleware_gateway.MiddlewareManager,
	communicationManager *communication_gateway.CommunicationManager) error {
	
	// Load configuration using existing library
	if err := configManager.Load(); err != nil {
		return err
	}
	
	// Initialize logging using existing library
	if err := loggingManager.Initialize(); err != nil {
		return err
	}
	
	// Start monitoring using existing library
	if err := monitoringManager.Start(); err != nil {
		return err
	}
	
	// Connect to database using existing library
	if err := databaseManager.Connect(ctx); err != nil {
		return err
	}
	
	// Initialize authentication using existing library
	if err := authManager.Initialize(); err != nil {
		return err
	}
	
	// Setup middleware using existing library
	if err := middlewareManager.SetupChain(); err != nil {
		return err
	}
	
	// Start communication server using existing library
	if err := communicationManager.Start(); err != nil {
		return err
	}
	
	return nil
}
`

	GoModTemplate = `module {{.ServiceName}}

go 1.21

require (
	// Use go-micro-libs library
	github.com/anasamu/go-micro-libs v1.0.0
	
	// Core dependencies
	github.com/gin-gonic/gin v1.9.1
	github.com/spf13/cobra v1.7.0
	github.com/spf13/viper v1.16.0
	
	// Database dependencies
	gorm.io/gorm v1.25.4
	gorm.io/driver/postgres v1.5.2
	gorm.io/driver/mysql v1.5.1
	gorm.io/driver/sqlite v1.5.2
	
	// Monitoring dependencies
	github.com/prometheus/client_golang v1.16.0
	go.opentelemetry.io/otel v1.16.0
	go.opentelemetry.io/otel/trace v1.16.0
	
	// Logging dependencies
	github.com/sirupsen/logrus v1.9.3
	github.com/rs/zerolog v1.30.0
	
	// Testing dependencies
	github.com/stretchr/testify v1.8.4
	github.com/golang/mock v1.6.0
)
`

	ConfigTemplate = `# Configuration for {{.ServiceName}}
service:
  name: "{{.ServiceName}}"
  version: "1.0.0"
  port: 8080
  environment: "development"

# Core configurations using existing libraries
config:
  providers:
    file:
      path: "./configs"
      format: "yaml"
    env:
      prefix: "{{.ServiceName}}_"

logging:
  providers:
    console:
      level: "info"
      format: "json"
    file:
      path: "/var/log/{{.ServiceName}}.log"
      level: "debug"

monitoring:
  providers:
    prometheus:
      endpoint: "http://localhost:9090"
      port: 9090
    jaeger:
      endpoint: "http://localhost:14268"
      service_name: "{{.ServiceName}}"

{{if .WithDatabase}}
database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100
      max_idle_connections: 10
    redis:
      url: "${REDIS_URL}"
      db: 0
{{end}}

{{if .WithAuth}}
auth:
  providers:
    jwt:
      secret: "${JWT_SECRET}"
      expiration: "24h"
      issuer: "{{.ServiceName}}"
    oauth:
      client_id: "${OAUTH_CLIENT_ID}"
      client_secret: "${OAUTH_CLIENT_SECRET}"
      redirect_url: "${OAUTH_REDIRECT_URL}"
{{end}}

middleware:
  auth:
    enabled: {{.WithAuth}}
    provider: "jwt"
  rate_limit:
    enabled: true
    requests_per_minute: 100
  circuit_breaker:
    enabled: true
    failure_threshold: 5
    timeout: 30s

communication:
  providers:
    rest:
      port: 8080
      timeout: 30s
    grpc:
      port: 9090
      timeout: 30s
`

	ConfigDevTemplate = `# Development configuration for {{.ServiceName}}
service:
  name: "{{.ServiceName}}"
  version: "1.0.0-dev"
  port: 8080
  environment: "development"

logging:
  providers:
    console:
      level: "debug"
      format: "text"

monitoring:
  providers:
    prometheus:
      endpoint: "http://localhost:9090"
      port: 9090
    jaeger:
      endpoint: "http://localhost:14268"
      service_name: "{{.ServiceName}}-dev"

{{if .WithDatabase}}
database:
  providers:
    postgresql:
      url: "postgres://localhost:5432/{{.ServiceName}}_dev?sslmode=disable"
      max_connections: 10
      max_idle_connections: 5
    redis:
      url: "redis://localhost:6379"
      db: 0
{{end}}

{{if .WithAuth}}
auth:
  providers:
    jwt:
      secret: "dev-secret-key"
      expiration: "24h"
      issuer: "{{.ServiceName}}-dev"
{{end}}

middleware:
  auth:
    enabled: false
  rate_limit:
    enabled: false
  circuit_breaker:
    enabled: false

communication:
  providers:
    rest:
      port: 8080
      timeout: 30s
`

	HandlersTemplate = `package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// ServiceHandler handles HTTP requests
type ServiceHandler struct {
	// Add service dependencies here
}

// NewServiceHandler creates a new handler
func NewServiceHandler() *ServiceHandler {
	return &ServiceHandler{}
}

// HealthCheck returns the health status of the service
func (h *ServiceHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "{{.ServiceName}}",
	})
}

// GetService returns a sample response
func (h *ServiceHandler) GetService(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello from {{.ServiceName}}",
		"data": "sample data",
	})
}

// CreateService creates a new resource
func (h *ServiceHandler) CreateService(c *gin.Context) {
	var request struct {
		Name string ` + "`json:\"name\"`" + `
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Created successfully",
		"data": request,
	})
}
`

	ModelsTemplate = `package models

import (
	"time"
	"gorm.io/gorm"
)

// ServiceModel represents the main entity
type ServiceModel struct {
	ID        uint           ` + "`json:\"id\" gorm:\"primaryKey\"`" + `
	Name      string         ` + "`json:\"name\" gorm:\"not null\"`" + `
	Email     string         ` + "`json:\"email\" gorm:\"uniqueIndex\"`" + `
	CreatedAt time.Time      ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time      ` + "`json:\"updated_at\"`" + `
	DeletedAt gorm.DeletedAt ` + "`json:\"deleted_at\" gorm:\"index\"`" + `
}

// TableName returns the table name for the model
func (ServiceModel) TableName() string {
	return "{{.ServiceName}}s"
}

// CreateServiceRequest represents the request to create a service
type CreateServiceRequest struct {
	Name  string ` + "`json:\"name\" binding:\"required\"`" + `
	Email string ` + "`json:\"email\" binding:\"required,email\"`" + `
}

// UpdateServiceRequest represents the request to update a service
type UpdateServiceRequest struct {
	Name  *string ` + "`json:\"name,omitempty\"`" + `
	Email *string ` + "`json:\"email,omitempty\"`" + `
}

// ServiceResponse represents the response for service operations
type ServiceResponse struct {
	ID        uint      ` + "`json:\"id\"`" + `
	Name      string    ` + "`json:\"name\"`" + `
	Email     string    ` + "`json:\"email\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}
`

	RepositoriesTemplate = `package repositories

import (
	"context"
	"{{.ServiceName}}/internal/models"
	"gorm.io/gorm"
)

// ServiceRepository handles data access
type ServiceRepository struct {
	db *gorm.DB
}

// NewServiceRepository creates a new repository
func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

// Create creates a new service
func (r *ServiceRepository) Create(ctx context.Context, service *models.ServiceModel) error {
	return r.db.WithContext(ctx).Create(service).Error
}

// GetByID retrieves a service by ID
func (r *ServiceRepository) GetByID(ctx context.Context, id uint) (*models.ServiceModel, error) {
	var service models.ServiceModel
	err := r.db.WithContext(ctx).First(&service, id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

// GetByEmail retrieves a service by email
func (r *ServiceRepository) GetByEmail(ctx context.Context, email string) (*models.ServiceModel, error) {
	var service models.ServiceModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

// Update updates a service
func (r *ServiceRepository) Update(ctx context.Context, service *models.ServiceModel) error {
	return r.db.WithContext(ctx).Save(service).Error
}

// Delete soft deletes a service
func (r *ServiceRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ServiceModel{}, id).Error
}

// List retrieves all services with pagination
func (r *ServiceRepository) List(ctx context.Context, offset, limit int) ([]*models.ServiceModel, error) {
	var services []*models.ServiceModel
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&services).Error
	return services, err
}

// Count returns the total number of services
func (r *ServiceRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.ServiceModel{}).Count(&count).Error
	return count, err
}
`

	ServicesTemplate = `package services

import (
	"context"
	"errors"
	"{{.ServiceName}}/internal/models"
	"{{.ServiceName}}/internal/repositories"
)

// ServiceService handles business logic
type ServiceService struct {
	repo *repositories.ServiceRepository
}

// NewServiceService creates a new service
func NewServiceService(repo *repositories.ServiceRepository) *ServiceService {
	return &ServiceService{
		repo: repo,
	}
}

// CreateService creates a new service
func (s *ServiceService) CreateService(ctx context.Context, req *models.CreateServiceRequest) (*models.ServiceResponse, error) {
	// Check if email already exists
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, errors.New("email already exists")
	}
	
	service := &models.ServiceModel{
		Name:  req.Name,
		Email: req.Email,
	}
	
	if err := s.repo.Create(ctx, service); err != nil {
		return nil, err
	}
	
	return &models.ServiceResponse{
		ID:        service.ID,
		Name:      service.Name,
		Email:     service.Email,
		CreatedAt: service.CreatedAt,
		UpdatedAt: service.UpdatedAt,
	}, nil
}

// GetService retrieves a service by ID
func (s *ServiceService) GetService(ctx context.Context, id uint) (*models.ServiceResponse, error) {
	service, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return &models.ServiceResponse{
		ID:        service.ID,
		Name:      service.Name,
		Email:     service.Email,
		CreatedAt: service.CreatedAt,
		UpdatedAt: service.UpdatedAt,
	}, nil
}

// UpdateService updates a service
func (s *ServiceService) UpdateService(ctx context.Context, id uint, req *models.UpdateServiceRequest) (*models.ServiceResponse, error) {
	service, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Update fields if provided
	if req.Name != nil {
		service.Name = *req.Name
	}
	if req.Email != nil {
		// Check if new email already exists
		existing, err := s.repo.GetByEmail(ctx, *req.Email)
		if err == nil && existing != nil && existing.ID != id {
			return nil, errors.New("email already exists")
		}
		service.Email = *req.Email
	}
	
	if err := s.repo.Update(ctx, service); err != nil {
		return nil, err
	}
	
	return &models.ServiceResponse{
		ID:        service.ID,
		Name:      service.Name,
		Email:     service.Email,
		CreatedAt: service.CreatedAt,
		UpdatedAt: service.UpdatedAt,
	}, nil
}

// DeleteService deletes a service
func (s *ServiceService) DeleteService(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// ListServices retrieves services with pagination
func (s *ServiceService) ListServices(ctx context.Context, offset, limit int) ([]*models.ServiceResponse, int64, error) {
	services, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	
	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	
	responses := make([]*models.ServiceResponse, len(services))
	for i, service := range services {
		responses[i] = &models.ServiceResponse{
			ID:        service.ID,
			Name:      service.Name,
			Email:     service.Email,
			CreatedAt: service.CreatedAt,
			UpdatedAt: service.UpdatedAt,
		}
	}
	
	return responses, count, nil
}
`

	MiddlewareTemplate = `package middleware

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware provides request logging
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// RecoveryMiddleware provides panic recovery
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}

// CORSMiddleware provides CORS support
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// RequestIDMiddleware adds a request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// TimeoutMiddleware adds a timeout to requests
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
`

	DockerfileTemplate = `# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create app user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Copy configuration files
COPY --from=builder /app/configs ./configs

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
`

	DockerComposeTemplate = `version: '3.8'

services:
  {{.ServiceName}}:
    build:
      context: .
      dockerfile: deployments/docker/Dockerfile
    ports:
      - "8080:8080"
    environment:
      - ENV=development
      - DATABASE_URL=postgres://postgres:password@postgres:5432/{{.ServiceName}}_dev?sslmode=disable
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=dev-secret-key
    depends_on:
      - postgres
      - redis
    networks:
      - {{.ServiceName}}-network

{{if .WithDatabase}}
  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB={{.ServiceName}}_dev
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - {{.ServiceName}}-network

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - {{.ServiceName}}-network
{{end}}

volumes:
  postgres_data:
  redis_data:

networks:
  {{.ServiceName}}-network:
    driver: bridge
`

	KubernetesDeploymentTemplate = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ServiceName}}
  labels:
    app: {{.ServiceName}}
spec:
  replicas: 3
  selector:
    matchLabels:
      app: {{.ServiceName}}
  template:
    metadata:
      labels:
        app: {{.ServiceName}}
    spec:
      containers:
      - name: {{.ServiceName}}
        image: {{.ServiceName}}:latest
        ports:
        - containerPort: 8080
        env:
        - name: ENV
          value: "production"
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: {{.ServiceName}}-secrets
              key: database-url
        - name: REDIS_URL
          valueFrom:
            secretKeyRef:
              name: {{.ServiceName}}-secrets
              key: redis-url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: {{.ServiceName}}-secrets
              key: jwt-secret
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
`

	KubernetesServiceTemplate = `apiVersion: v1
kind: Service
metadata:
  name: {{.ServiceName}}-service
  labels:
    app: {{.ServiceName}}
spec:
  selector:
    app: {{.ServiceName}}
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP
`

	KubernetesConfigMapTemplate = `apiVersion: v1
kind: ConfigMap
metadata:
  name: {{.ServiceName}}-config
data:
  config.yaml: |
    service:
      name: "{{.ServiceName}}"
      version: "1.0.0"
      port: 8080
      environment: "production"
    
    logging:
      providers:
        console:
          level: "info"
          format: "json"
    
    monitoring:
      providers:
        prometheus:
          endpoint: "http://prometheus:9090"
          port: 9090
        jaeger:
          endpoint: "http://jaeger:14268"
          service_name: "{{.ServiceName}}"
    
    middleware:
      auth:
        enabled: true
        provider: "jwt"
      rate_limit:
        enabled: true
        requests_per_minute: 100
      circuit_breaker:
        enabled: true
        failure_threshold: 5
        timeout: 30s
    
    communication:
      providers:
        rest:
          port: 8080
          timeout: 30s
`

	UnitTestTemplate = "package handlers\n\n" +
		"import (\n" +
		"	\"net/http\"\n" +
		"	\"net/http/httptest\"\n" +
		"	\"testing\"\n" +
		"	\"github.com/gin-gonic/gin\"\n" +
		"	\"github.com/stretchr/testify/assert\"\n" +
		")\n\n" +
		"func TestServiceHandler_HealthCheck(t *testing.T) {\n" +
		"	gin.SetMode(gin.TestMode)\n" +
		"	\n" +
		"	handler := NewServiceHandler()\n" +
		"	router := gin.New()\n" +
		"	router.GET(\"/health\", handler.HealthCheck)\n" +
		"	\n" +
		"	req, _ := http.NewRequest(\"GET\", \"/health\", nil)\n" +
		"	w := httptest.NewRecorder()\n" +
		"	\n" +
		"	router.ServeHTTP(w, req)\n" +
		"	\n" +
		"	assert.Equal(t, http.StatusOK, w.Code)\n" +
		"	assert.Contains(t, w.Body.String(), \"healthy\")\n" +
		"	assert.Contains(t, w.Body.String(), \"{{.ServiceName}}\")\n" +
		"}\n\n" +
		"func TestServiceHandler_GetService(t *testing.T) {\n" +
		"	gin.SetMode(gin.TestMode)\n" +
		"	\n" +
		"	handler := NewServiceHandler()\n" +
		"	router := gin.New()\n" +
		"	router.GET(\"/service\", handler.GetService)\n" +
		"	\n" +
		"	req, _ := http.NewRequest(\"GET\", \"/service\", nil)\n" +
		"	w := httptest.NewRecorder()\n" +
		"	\n" +
		"	router.ServeHTTP(w, req)\n" +
		"	\n" +
		"	assert.Equal(t, http.StatusOK, w.Code)\n" +
		"	assert.Contains(t, w.Body.String(), \"Hello from {{.ServiceName}}\")\n" +
		"}\n\n" +
		"func TestServiceHandler_CreateService(t *testing.T) {\n" +
		"	gin.SetMode(gin.TestMode)\n" +
		"	\n" +
		"	handler := NewServiceHandler()\n" +
		"	router := gin.New()\n" +
		"	router.POST(\"/service\", handler.CreateService)\n" +
		"	\n" +
		"	// Test valid request\n" +
		"	req, _ := http.NewRequest(\"POST\", \"/service\", strings.NewReader(\"{\\\"name\\\":\\\"Test\\\",\\\"email\\\":\\\"test@example.com\\\"}\"))\n" +
		"	req.Header.Set(\"Content-Type\", \"application/json\")\n" +
		"	w := httptest.NewRecorder()\n" +
		"	\n" +
		"	router.ServeHTTP(w, req)\n" +
		"	\n" +
		"	assert.Equal(t, http.StatusCreated, w.Code)\n" +
		"	assert.Contains(t, w.Body.String(), \"Created successfully\")\n" +
		"	\n" +
		"	// Test invalid request\n" +
		"	req, _ = http.NewRequest(\"POST\", \"/service\", strings.NewReader(\"{\\\"invalid\\\":\\\"json\\\"}\"))\n" +
		"	req.Header.Set(\"Content-Type\", \"application/json\")\n" +
		"	w = httptest.NewRecorder()\n" +
		"	\n" +
		"	router.ServeHTTP(w, req)\n" +
		"	\n" +
		"	assert.Equal(t, http.StatusBadRequest, w.Code)\n" +
		"}"

	IntegrationTestTemplate = `package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"{{.ServiceName}}/internal/models"
	"{{.ServiceName}}/internal/handlers"
	"{{.ServiceName}}/internal/services"
	"{{.ServiceName}}/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ServiceIntegrationTestSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
	handler *handlers.ServiceHandler
}

func (suite *ServiceIntegrationTestSuite) SetupSuite() {
	// Setup test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)
	suite.db = db
	
	// Auto migrate
	err = db.AutoMigrate(&models.ServiceModel{})
	suite.Require().NoError(err)
	
	// Setup dependencies
	repo := repositories.NewServiceRepository(db)
	service := services.NewServiceService(repo)
	handler := handlers.NewServiceHandler()
	
	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health", handler.HealthCheck)
	router.GET("/service", handler.GetService)
	router.POST("/service", handler.CreateService)
	
	suite.router = router
	suite.handler = handler
}

func (suite *ServiceIntegrationTestSuite) TearDownSuite() {
	// Cleanup
}

func (suite *ServiceIntegrationTestSuite) TestHealthCheck() {
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	
	suite.router.ServeHTTP(w, req)
	
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *ServiceIntegrationTestSuite) TestCreateAndGetService() {
	// Create service
	createReq := models.CreateServiceRequest{
		Name:  "Test Service",
		Email: "test@example.com",
	}
	
	jsonData, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/service", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	suite.router.ServeHTTP(w, req)
	
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	
	// Verify creation
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	
	assert.Contains(suite.T(), response, "data")
}

func TestServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceIntegrationTestSuite))
}
`

	ReadmeTemplate = "# {{.ServiceName}} Service\n\n" +
		"A microservice built with Go Micro Framework using go-micro-libs library.\n\n" +
		"## Features\n\n" +
		"- REST API with Gin framework\n" +
		"- Database integration with GORM\n" +
		"- Authentication and authorization\n" +
		"- Monitoring and logging\n" +
		"- Docker containerization\n" +
		"- Kubernetes deployment\n" +
		"- Comprehensive testing\n\n" +
		"## Quick Start\n\n" +
		"### Prerequisites\n\n" +
		"- Go 1.21+\n" +
		"- Docker and Docker Compose\n" +
		"- Make (optional)\n\n" +
		"### Development\n\n" +
		"1. Clone the repository:\n" +
		"```bash\n" +
		"git clone <repository-url>\n" +
		"cd {{.ServiceName}}\n" +
		"```\n\n" +
		"2. Install dependencies:\n" +
		"```bash\n" +
		"go mod tidy\n" +
		"```\n\n" +
		"3. Set up environment variables:\n" +
		"```bash\n" +
		"export DATABASE_URL=\"postgres://localhost:5432/{{.ServiceName}}_dev?sslmode=disable\"\n" +
		"export REDIS_URL=\"redis://localhost:6379\"\n" +
		"export JWT_SECRET=\"your-secret-key\"\n" +
		"```\n\n" +
		"4. Run the service:\n" +
		"```bash\n" +
		"go run cmd/main.go\n" +
		"```\n\n" +
		"### Docker\n\n" +
		"1. Build and run with Docker Compose:\n" +
		"```bash\n" +
		"docker-compose -f deployments/docker/docker-compose.yml up --build\n" +
		"```\n\n" +
		"### Kubernetes\n\n" +
		"1. Apply Kubernetes manifests:\n" +
		"```bash\n" +
		"kubectl apply -f deployments/kubernetes/\n" +
		"```\n\n" +
		"## API Endpoints\n\n" +
		"- `GET /health` - Health check\n" +
		"- `GET /service` - Get sample data\n" +
		"- `POST /service` - Create new resource\n\n" +
		"## Configuration\n\n" +
		"The service uses YAML configuration files located in the `configs/` directory:\n\n" +
		"- `config.yaml` - Production configuration\n" +
		"- `config.dev.yaml` - Development configuration\n\n" +
		"## Testing\n\n" +
		"Run tests:\n" +
		"```bash\n" +
		"# Unit tests\n" +
		"go test ./tests/unit/...\n\n" +
		"# Integration tests\n" +
		"go test ./tests/integration/...\n\n" +
		"# All tests\n" +
		"go test ./...\n" +
		"```\n\n" +
		"## Monitoring\n\n" +
		"The service includes built-in monitoring with:\n\n" +
		"- Prometheus metrics\n" +
		"- Jaeger tracing\n" +
		"- Structured logging\n\n" +
		"Access metrics at: `http://localhost:9090/metrics`\n\n" +
		"## Contributing\n\n" +
		"1. Fork the repository\n" +
		"2. Create a feature branch\n" +
		"3. Make your changes\n" +
		"4. Add tests\n" +
		"5. Submit a pull request\n\n" +
		"## License\n\n" +
		"MIT License"

	APITemplate = "# {{.ServiceName}} API Documentation\n\n" +
		"## Overview\n\n" +
		"The {{.ServiceName}} service provides REST API endpoints for managing service resources.\n\n" +
		"## Base URL\n\n" +
		"`http://localhost:8080`\n\n" +
		"## Authentication\n\n" +
		"The API uses JWT-based authentication. Include the token in the Authorization header:\n\n" +
		"```\n" +
		"Authorization: Bearer <your-jwt-token>\n" +
		"```\n\n" +
		"## Endpoints\n\n" +
		"### Health Check\n\n" +
		"Check the health status of the service.\n\n" +
		"```http\n" +
		"GET /health\n" +
		"```\n\n" +
		"**Response:**\n" +
		"```json\n" +
		"{\n" +
		"  \"status\": \"healthy\",\n" +
		"  \"service\": \"{{.ServiceName}}\"\n" +
		"}\n" +
		"```\n\n" +
		"### Get Service\n\n" +
		"Retrieve sample data from the service.\n\n" +
		"```http\n" +
		"GET /service\n" +
		"```\n\n" +
		"**Response:**\n" +
		"```json\n" +
		"{\n" +
		"  \"message\": \"Hello from {{.ServiceName}}\",\n" +
		"  \"data\": \"sample data\"\n" +
		"}\n" +
		"```\n\n" +
		"### Create Service\n\n" +
		"Create a new service resource.\n\n" +
		"```http\n" +
		"POST /service\n" +
		"Content-Type: application/json\n" +
		"```\n\n" +
		"**Request Body:**\n" +
		"```json\n" +
		"{\n" +
		"  \"name\": \"string\",\n" +
		"  \"email\": \"string\"\n" +
		"}\n" +
		"```\n\n" +
		"**Response:**\n" +
		"```json\n" +
		"{\n" +
		"  \"message\": \"Created successfully\",\n" +
		"  \"data\": {\n" +
		"    \"name\": \"string\",\n" +
		"    \"email\": \"string\"\n" +
		"  }\n" +
		"}\n" +
		"```\n\n" +
		"## Error Responses\n\n" +
		"All error responses follow this format:\n\n" +
		"```json\n" +
		"{\n" +
		"  \"error\": \"error message\"\n" +
		"}\n" +
		"```\n\n" +
		"## Status Codes\n\n" +
		"- `200 OK` - Request successful\n" +
		"- `201 Created` - Resource created successfully\n" +
		"- `400 Bad Request` - Invalid request data\n" +
		"- `401 Unauthorized` - Authentication required\n" +
		"- `404 Not Found` - Resource not found\n" +
		"- `500 Internal Server Error` - Server error\n\n" +
		"## Rate Limiting\n\n" +
		"The API implements rate limiting:\n" +
		"- 100 requests per minute per IP address\n" +
		"- Rate limit headers are included in responses\n\n" +
		"## Monitoring\n\n" +
		"The service exposes Prometheus metrics at `/metrics` endpoint."

	MigrationExampleTemplate = `{
  "version": "{{.Timestamp}}",
  "description": "{{.Description}}",
  "up_sql": "-- Add your up migration SQL here\n-- Example:\n-- CREATE TABLE {{.ServiceName | lower}}_users (\n--     id SERIAL PRIMARY KEY,\n--     name VARCHAR(255) NOT NULL,\n--     email VARCHAR(255) UNIQUE NOT NULL,\n--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,\n--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP\n-- );",
  "down_sql": "-- Add your down migration SQL here\n-- Example:\n-- DROP TABLE IF EXISTS {{.ServiceName | lower}}_users;",
  "created_at": "{{.CreatedAt}}",
  "checksum": ""
}`

	UtilsTemplate = `package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// UtilsManager provides utility functions for the service
type UtilsManager struct {
	serviceName string
	serviceID   string
}

// NewUtilsManager creates a new utils manager
func NewUtilsManager(serviceName string) *UtilsManager {
	return &UtilsManager{
		serviceName: serviceName,
		serviceID:   uuid.New().String(),
	}
}

// GetServiceID returns the service instance ID
func (u *UtilsManager) GetServiceID() string {
	return u.serviceID
}

// GetServiceName returns the service name
func (u *UtilsManager) GetServiceName() string {
	return u.serviceName
}

// GenerateUUID generates a new UUID
func (u *UtilsManager) GenerateUUID() string {
	return uuid.New().String()
}

// GenerateUUIDWithNamespace generates a UUID with namespace
func (u *UtilsManager) GenerateUUIDWithNamespace(namespace string) string {
	namespaceUUID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(namespace))
	return uuid.NewSHA1(namespaceUUID, []byte(u.serviceName)).String()
}

// GenerateRandomString generates a random string of specified length
func (u *UtilsManager) GenerateRandomString(length int) string {
	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// LoadEnvironment loads environment variables from .env file
func (u *UtilsManager) LoadEnvironment(envFile string) error {
	if envFile == "" {
		envFile = ".env"
	}
	
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return fmt.Errorf("environment file %s not found", envFile)
	}
	
	return godotenv.Load(envFile)
}

// GetEnv gets environment variable with default value
func (u *UtilsManager) GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvAsInt gets environment variable as integer with default value
func (u *UtilsManager) GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvAsBool gets environment variable as boolean with default value
func (u *UtilsManager) GetEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// ValidateRequiredEnv validates that required environment variables are set
func (u *UtilsManager) ValidateRequiredEnv(requiredVars []string) error {
	var missingVars []string
	
	for _, varName := range requiredVars {
		if os.Getenv(varName) == "" {
			missingVars = append(missingVars, varName)
		}
	}
	
	if len(missingVars) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missingVars, ", "))
	}
	
	return nil
}

// FormatTimestamp formats timestamp to RFC3339 format
func (u *UtilsManager) FormatTimestamp(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseTimestamp parses RFC3339 timestamp
func (u *UtilsManager) ParseTimestamp(timestamp string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestamp)
}

// GetCurrentTimestamp returns current timestamp in RFC3339 format
func (u *UtilsManager) GetCurrentTimestamp() string {
	return u.FormatTimestamp(time.Now())
}

// SanitizeString removes special characters from string
func (u *UtilsManager) SanitizeString(input string) string {
	// Remove special characters except alphanumeric, spaces, hyphens, and underscores
	var result strings.Builder
	for _, char := range input {
		if (char >= 'a' && char <= 'z') || 
		   (char >= 'A' && char <= 'Z') || 
		   (char >= '0' && char <= '9') || 
		   char == ' ' || char == '-' || char == '_' {
			result.WriteRune(char)
		}
	}
	return strings.TrimSpace(result.String())
}

// TruncateString truncates string to specified length
func (u *UtilsManager) TruncateString(input string, maxLength int) string {
	if len(input) <= maxLength {
		return input
	}
	return input[:maxLength] + "..."
}

// IsEmpty checks if string is empty or contains only whitespace
func (u *UtilsManager) IsEmpty(input string) bool {
	return strings.TrimSpace(input) == ""
}

// Contains checks if slice contains the specified item
func (u *UtilsManager) Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveDuplicates removes duplicate strings from slice
func (u *UtilsManager) RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string
	
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// MergeMaps merges two maps, with second map taking precedence
func (u *UtilsManager) MergeMaps(map1, map2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	
	// Add all items from first map
	for k, v := range map1 {
		result[k] = v
	}
	
	// Add/override with items from second map
	for k, v := range map2 {
		result[k] = v
	}
	
	return result
}

// GetMapValue gets value from map with default
func (u *UtilsManager) GetMapValue(m map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if value, exists := m[key]; exists {
		return value
	}
	return defaultValue
}

// ConvertToString converts interface{} to string
func (u *UtilsManager) ConvertToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ConvertToInt converts interface{} to int
func (u *UtilsManager) ConvertToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", value)
	}
}

// ConvertToBool converts interface{} to bool
func (u *UtilsManager) ConvertToBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	case int:
		return v != 0, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", value)
	}
	}
`

	EnvExampleTemplate = `# Environment variables for {{.ServiceName}} service

# Service Configuration
{{.ServiceName | upper}}_SERVICE_NAME={{.ServiceName}}
{{.ServiceName | upper}}_SERVICE_VERSION=1.0.0
{{.ServiceName | upper}}_SERVICE_PORT=8080
{{.ServiceName | upper}}_SERVICE_ENVIRONMENT=development

# Core Configuration
{{.ServiceName | upper}}_CONFIG_PATH=./configs
{{.ServiceName | upper}}_CONFIG_FORMAT=yaml

# Core Logging
{{.ServiceName | upper}}_LOG_LEVEL=info
{{.ServiceName | upper}}_LOG_FORMAT=json
{{.ServiceName | upper}}_LOG_FILE_PATH=/var/log/{{.ServiceName}}.log

# Core Monitoring
{{.ServiceName | upper}}_PROMETHEUS_ENDPOINT=http://localhost:9090
{{.ServiceName | upper}}_JAEGER_ENDPOINT=http://localhost:14268
{{.ServiceName | upper}}_GRAFANA_ENDPOINT=http://localhost:3000

# Core Middleware
{{.ServiceName | upper}}_RATE_LIMIT_ENABLED=true
{{.ServiceName | upper}}_RATE_LIMIT_REQUESTS_PER_MINUTE=100
{{.ServiceName | upper}}_CIRCUIT_BREAKER_ENABLED=true
{{.ServiceName | upper}}_CIRCUIT_BREAKER_FAILURE_THRESHOLD=5

# Core Communication
{{.ServiceName | upper}}_REST_PORT=8080
{{.ServiceName | upper}}_GRPC_PORT=9090
{{.ServiceName | upper}}_REQUEST_TIMEOUT=30s

# Core Utils
{{.ServiceName | upper}}_UUID_VERSION=4
{{.ServiceName | upper}}_UUID_NAMESPACE={{.ServiceName}}
{{.ServiceName | upper}}_VALIDATION_ENABLED=true
{{.ServiceName | upper}}_VALIDATION_STRICT_MODE=false

{{- if .WithDatabase}}
# Database Configuration
{{.ServiceName | upper}}_DATABASE_URL=postgres://localhost:5432/{{.ServiceName}}_dev?sslmode=disable
{{.ServiceName | upper}}_DATABASE_MAX_CONNECTIONS=100
{{.ServiceName | upper}}_DATABASE_MAX_IDLE_CONNECTIONS=10
{{.ServiceName | upper}}_REDIS_URL=redis://localhost:6379
{{.ServiceName | upper}}_REDIS_DB=0
{{.ServiceName | upper}}_REDIS_POOL_SIZE=10
{{- end}}

{{- if .WithAuth}}
# Authentication Configuration
{{.ServiceName | upper}}_JWT_SECRET=your-jwt-secret-key-here
{{.ServiceName | upper}}_JWT_EXPIRATION=24h
{{.ServiceName | upper}}_JWT_ISSUER={{.ServiceName}}
{{.ServiceName | upper}}_OAUTH_CLIENT_ID=your-oauth-client-id
{{.ServiceName | upper}}_OAUTH_CLIENT_SECRET=your-oauth-client-secret
{{.ServiceName | upper}}_OAUTH_REDIRECT_URL=http://localhost:8080/auth/callback
{{- end}}

{{- if .WithMessaging}}
# Messaging Configuration
{{.ServiceName | upper}}_KAFKA_BROKERS=localhost:9092
{{.ServiceName | upper}}_KAFKA_GROUP_ID={{.ServiceName}}
{{.ServiceName | upper}}_RABBITMQ_URL=amqp://localhost:5672
{{.ServiceName | upper}}_RABBITMQ_EXCHANGE={{.ServiceName}}-exchange
{{.ServiceName | upper}}_RABBITMQ_QUEUE={{.ServiceName}}-queue
{{- end}}

{{- if .WithAI}}
# AI Services Configuration
{{.ServiceName | upper}}_OPENAI_API_KEY=your-openai-api-key
{{.ServiceName | upper}}_OPENAI_DEFAULT_MODEL=gpt-4
{{.ServiceName | upper}}_ANTHROPIC_API_KEY=your-anthropic-api-key
{{.ServiceName | upper}}_ANTHROPIC_DEFAULT_MODEL=claude-3-sonnet
{{- end}}

{{- if .WithStorage}}
# Storage Configuration
{{.ServiceName | upper}}_AWS_ACCESS_KEY_ID=your-aws-access-key
{{.ServiceName | upper}}_AWS_SECRET_ACCESS_KEY=your-aws-secret-key
{{.ServiceName | upper}}_AWS_REGION=us-east-1
{{.ServiceName | upper}}_S3_BUCKET=your-s3-bucket
{{.ServiceName | upper}}_GCS_CREDENTIALS_FILE=path/to/gcs-credentials.json
{{.ServiceName | upper}}_GCS_BUCKET=your-gcs-bucket
{{- end}}

{{- if .WithCache}}
# Cache Configuration
{{.ServiceName | upper}}_CACHE_REDIS_URL=redis://localhost:6379
{{.ServiceName | upper}}_CACHE_REDIS_DB=2
{{.ServiceName | upper}}_CACHE_TTL=1h
{{.ServiceName | upper}}_CACHE_MEMORY_MAX_SIZE=1000
{{.ServiceName | upper}}_CACHE_MEMORY_TTL=30m
{{- end}}

{{- if .WithDiscovery}}
# Service Discovery Configuration
{{.ServiceName | upper}}_CONSUL_ADDRESS=localhost:8500
{{.ServiceName | upper}}_CONSUL_TOKEN=your-consul-token
{{.ServiceName | upper}}_KUBERNETES_CONFIG=path/to/kubeconfig
{{- end}}

{{- if .WithPayment}}
# Payment Configuration
{{.ServiceName | upper}}_STRIPE_API_KEY=your-stripe-api-key
{{.ServiceName | upper}}_STRIPE_WEBHOOK_SECRET=your-stripe-webhook-secret
{{.ServiceName | upper}}_PAYPAL_CLIENT_ID=your-paypal-client-id
{{.ServiceName | upper}}_PAYPAL_CLIENT_SECRET=your-paypal-client-secret
{{- end}}

{{- if .WithAPI}}
# API Integration Configuration
{{.ServiceName | upper}}_API_HTTP_TIMEOUT=30s
{{.ServiceName | upper}}_API_HTTP_RETRY_ATTEMPTS=3
{{.ServiceName | upper}}_API_HTTP_RETRY_DELAY=5s
{{.ServiceName | upper}}_API_GRAPHQL_ENDPOINT=http://localhost:4000/graphql
{{.ServiceName | upper}}_API_WEBSOCKET_TIMEOUT=30s
{{- end}}

{{- if .WithEmail}}
# Email Configuration
{{.ServiceName | upper}}_SMTP_HOST=smtp.gmail.com
{{.ServiceName | upper}}_SMTP_PORT=587
{{.ServiceName | upper}}_SMTP_USERNAME=your-email@gmail.com
{{.ServiceName | upper}}_SMTP_PASSWORD=your-email-password
{{.ServiceName | upper}}_SENDGRID_API_KEY=your-sendgrid-api-key
{{.ServiceName | upper}}_MAILGUN_API_KEY=your-mailgun-api-key
{{.ServiceName | upper}}_MAILGUN_DOMAIN=your-mailgun-domain
{{- end}}

# External Services
{{.ServiceName | upper}}_ELASTICSEARCH_ENDPOINT=http://localhost:9200
{{.ServiceName | upper}}_CONSUL_ADDRESS=localhost:8500
{{.ServiceName | upper}}_CONSUL_TOKEN=your-consul-token
`
)

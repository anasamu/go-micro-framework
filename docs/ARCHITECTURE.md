# Architecture Documentation

## Overview

Go Micro Framework adalah framework yang mengintegrasikan library `go-micro-libs` menjadi platform yang mudah digunakan untuk pengembangan microservices. Framework ini menggunakan pola Gateway dan Manager untuk mengorchestrasi berbagai provider dari library yang sudah ada.

## Core Architecture

### 1. Gateway Pattern

Framework menggunakan **Gateway Pattern** sebagai pola arsitektur utama:

```
┌─────────────────────────────────────────────────────────────┐
│                    Go Micro Framework                      │
├─────────────────────────────────────────────────────────────┤
│  CLI Tool (cobra)                                          │
│  ├── new command                                           │
│  ├── add command                                           │
│  ├── generate command                                      │
│  ├── deploy command                                        │
│  └── ... other commands                                    │
├─────────────────────────────────────────────────────────────┤
│  Core Components                                           │
│  ├── Bootstrap Engine                                      │
│  ├── Service Generator                                     │
│  ├── Template System                                       │
│  └── Configuration Manager                                 │
├─────────────────────────────────────────────────────────────┤
│  Library Integration Layer (Gateway)                       │
│  ├── AI Gateway                                            │
│  ├── Auth Gateway                                          │
│  ├── Database Gateway                                      │
│  ├── Storage Gateway                                       │
│  ├── Messaging Gateway                                     │
│  ├── Monitoring Gateway                                    │
│  └── ... other gateways                                    │
├─────────────────────────────────────────────────────────────┤
│  go-micro-libs                                             │
│  ├── ai/                                                   │
│  ├── auth/                                                 │
│  ├── database/                                             │
│  ├── storage/                                              │
│  ├── messaging/                                            │
│  ├── monitoring/                                           │
│  └── ... other libraries                                   │
└─────────────────────────────────────────────────────────────┘
```

### 2. Manager Pattern

Setiap library menggunakan **Manager Pattern** untuk mengelola provider:

```go
type Manager struct {
    providers map[string]Provider
    config    *Config
    logger    *logrus.Logger
}

func (m *Manager) RegisterProvider(name string, provider Provider) error
func (m *Manager) GetProvider(name string) (Provider, error)
func (m *Manager) ListProviders() []string
func (m *Manager) HealthCheck(ctx context.Context) error
```

### 3. Provider Interface

Setiap library memiliki interface provider yang konsisten:

```go
type Provider interface {
    Initialize(ctx context.Context, config interface{}) error
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    HealthCheck(ctx context.Context) error
    GetConfig() interface{}
}
```

## Framework Components

### 1. CLI Tool

CLI tool menggunakan `cobra` untuk command management:

```
microframework/
├── cmd/
│   └── microframework/
│       ├── main.go              # Entry point
│       └── commands/
│           ├── root.go          # Root command
│           ├── new.go           # New service command
│           ├── add.go           # Add feature command
│           ├── generate.go      # Generate component command
│           ├── deploy.go        # Deploy command
│           ├── config.go        # Config management
│           ├── validate.go      # Validation command
│           ├── logs.go          # Logs command
│           ├── health.go        # Health check command
│           ├── update.go        # Update command
│           └── version.go       # Version command
```

### 2. Core Components

#### Bootstrap Engine

Bootstrap engine menginisialisasi dan mengorchestrasi semua komponen:

```go
type Bootstrap struct {
    configManager    *config.Manager
    loggingManager   *logging.Manager
    monitoringManager *monitoring.Manager
    databaseManager  *database.Manager
    authManager      *auth.Manager
    middlewareManager *middleware.Manager
    // ... other managers
}

func (b *Bootstrap) Initialize(ctx context.Context) error {
    // Initialize all managers in order
    if err := b.configManager.Load(); err != nil {
        return err
    }
    
    if err := b.loggingManager.Initialize(); err != nil {
        return err
    }
    
    // ... initialize other managers
    
    return nil
}
```

#### Service Generator

Service generator menggunakan template system untuk menghasilkan kode:

```go
type ServiceGenerator struct {
    templates map[string]*template.Template
    logger    *logrus.Logger
}

func (sg *ServiceGenerator) GenerateService(config *ServiceConfig) error {
    // Load templates
    // Process templates with config
    // Generate files
    // Create directory structure
}
```

#### Template System

Template system menggunakan Go templates untuk code generation:

```
templates/
├── main.go.tmpl                 # Main application template
├── handlers.go.tmpl             # HTTP handlers template
├── services.go.tmpl             # Business services template
├── models.go.tmpl               # Data models template
├── repositories.go.tmpl         # Data access template
├── middleware.go.tmpl           # Middleware template
├── dockerfile.tmpl              # Dockerfile template
├── docker-compose.yml.tmpl      # Docker Compose template
├── kubernetes-deployment.yaml.tmpl  # K8s deployment template
├── kubernetes-service.yaml.tmpl     # K8s service template
├── kubernetes-configmap.yaml.tmpl   # K8s configmap template
├── unit-test.go.tmpl            # Unit test template
├── integration-test.go.tmpl     # Integration test template
├── README.md.tmpl               # README template
└── API.md.tmpl                  # API documentation template
```

### 3. Library Integration

Framework mengintegrasikan semua library dari `go-micro-libs`:

#### AI Services
```go
type AIManager struct {
    providers map[string]AIProvider
    config    *AIConfig
}

type AIProvider interface {
    Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
    GenerateText(ctx context.Context, req *TextRequest) (*TextResponse, error)
    CreateEmbedding(ctx context.Context, text string) (*Embedding, error)
}
```

#### Authentication
```go
type AuthManager struct {
    providers map[string]AuthProvider
    config    *AuthConfig
}

type AuthProvider interface {
    Authenticate(ctx context.Context, req *AuthRequest) (*AuthResponse, error)
    ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
    Authorize(ctx context.Context, user *User, resource string, action string) error
}
```

#### Database
```go
type DatabaseManager struct {
    providers map[string]DatabaseProvider
    config    *DatabaseConfig
}

type DatabaseProvider interface {
    Connect(ctx context.Context) error
    Disconnect(ctx context.Context) error
    Query(ctx context.Context, query string, args ...interface{}) (*Rows, error)
    Exec(ctx context.Context, query string, args ...interface{}) (*Result, error)
    Transaction(ctx context.Context, fn func(*Tx) error) error
}
```

## Service Architecture

### Generated Service Structure

Framework menghasilkan service dengan struktur yang konsisten:

```
user-service/
├── cmd/
│   └── main.go                 # Bootstrap code
├── internal/
│   ├── handlers/               # HTTP handlers
│   │   └── user_handler.go
│   ├── services/               # Business services
│   │   └── user_service.go
│   ├── repositories/           # Data access layer
│   │   └── user_repository.go
│   ├── models/                 # Data models
│   │   └── user.go
│   ├── middleware/             # Custom middleware
│   │   └── auth_middleware.go
│   └── config/                 # Configuration
│       └── config.go
├── pkg/
│   └── types/                  # Public types
│       └── user.go
├── configs/
│   ├── config.yaml             # Default configuration
│   ├── config.dev.yaml         # Development config
│   └── config.prod.yaml        # Production config
├── deployments/
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── kubernetes/
│       ├── deployment.yaml
│       ├── service.yaml
│       └── configmap.yaml
├── tests/
│   ├── unit/
│   │   └── user_test.go
│   ├── integration/
│   │   └── user_integration_test.go
│   └── e2e/
│       └── user_e2e_test.go
├── go.mod                      # Dependencies
├── go.sum                      # Checksums
├── Makefile                    # Build automation
└── README.md                   # Service documentation
```

### Service Bootstrap

Generated service menggunakan bootstrap engine:

```go
func main() {
    ctx := context.Background()
    
    // Initialize bootstrap
    bootstrap := core.NewBootstrap()
    
    // Load configuration
    if err := bootstrap.LoadConfig("configs/config.yaml"); err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // Initialize all components
    if err := bootstrap.Initialize(ctx); err != nil {
        log.Fatal("Failed to initialize:", err)
    }
    
    // Start service
    if err := bootstrap.Start(ctx); err != nil {
        log.Fatal("Failed to start:", err)
    }
    
    // Wait for shutdown signal
    bootstrap.WaitForShutdown()
}
```

## Configuration Management

### Configuration Structure

Framework menggunakan konfigurasi yang sesuai dengan library yang ada:

```yaml
# config.yaml
service:
  name: "user-service"
  version: "1.0.0"
  port: 8080

# Core configurations
config:
  providers:
    file:
      path: "./configs"
      format: "yaml"
    env:
      prefix: "USER_SERVICE_"

logging:
  providers:
    console:
      level: "info"
      format: "json"
    file:
      path: "/var/log/user-service.log"
      level: "debug"

monitoring:
  providers:
    prometheus:
      endpoint: "${PROMETHEUS_ENDPOINT}"
      port: 9090
    jaeger:
      endpoint: "${JAEGER_ENDPOINT}"
      service_name: "user-service"

database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100
    redis:
      url: "${REDIS_URL}"
      db: 0

auth:
  providers:
    jwt:
      secret: "${JWT_SECRET}"
      expiration: "24h"
    oauth:
      client_id: "${OAUTH_CLIENT_ID}"
      client_secret: "${OAUTH_CLIENT_SECRET}"

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

communication:
  providers:
    rest:
      port: 8080
      timeout: 30s
    grpc:
      port: 9090
      timeout: 30s

# Optional configurations
ai:
  providers:
    openai:
      api_key: "${OPENAI_API_KEY}"
      default_model: "gpt-4"

storage:
  providers:
    s3:
      access_key: "${AWS_ACCESS_KEY_ID}"
      secret_key: "${AWS_SECRET_ACCESS_KEY}"
      region: "${AWS_REGION}"
      bucket: "${S3_BUCKET}"

messaging:
  providers:
    kafka:
      brokers: "${KAFKA_BROKERS}"
      group_id: "user-service"
      topics: ["user-events"]

payment:
  providers:
    stripe:
      secret_key: "${STRIPE_SECRET_KEY}"
      webhook_secret: "${STRIPE_WEBHOOK_SECRET}"
```

### Environment Variables

Framework mendukung environment variables untuk konfigurasi:

```bash
# Service Configuration
SERVICE_NAME=user-service
SERVICE_VERSION=1.0.0
SERVICE_PORT=8080

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/mydb
REDIS_URL=redis://localhost:6379/0

# Authentication
JWT_SECRET=your-jwt-secret
OAUTH_CLIENT_ID=your-oauth-client-id
OAUTH_CLIENT_SECRET=your-oauth-client-secret

# Monitoring
PROMETHEUS_ENDPOINT=http://localhost:9090
JAEGER_ENDPOINT=http://localhost:14268

# AI Services
OPENAI_API_KEY=your-openai-api-key
ANTHROPIC_API_KEY=your-anthropic-api-key

# Storage
AWS_ACCESS_KEY_ID=your-aws-access-key
AWS_SECRET_ACCESS_KEY=your-aws-secret-key
S3_BUCKET=your-s3-bucket

# Messaging
KAFKA_BROKERS=localhost:9092

# Payment
STRIPE_SECRET_KEY=your-stripe-secret-key
STRIPE_WEBHOOK_SECRET=your-stripe-webhook-secret
```

## Deployment Architecture

### Docker Deployment

Framework menghasilkan Dockerfile yang optimal:

```dockerfile
# Multi-stage build
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary and configs
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

# Change ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run application
CMD ["./main"]
```

### Kubernetes Deployment

Framework menghasilkan manifest Kubernetes yang production-ready:

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  labels:
    app: user-service
    version: v1.0.0
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
        version: v1.0.0
    spec:
      serviceAccountName: user-service
      securityContext:
        runAsNonRoot: true
        runAsUser: 1001
        fsGroup: 1001
      containers:
      - name: user-service
        image: user-service:latest
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: metrics
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: user-service-secrets
              key: database-url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: user-service-secrets
              key: jwt-secret
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
            - ALL
```

## Security Architecture

### Security Layers

Framework mengimplementasikan multiple layers of security:

1. **Authentication Layer**: JWT, OAuth2, LDAP, SAML
2. **Authorization Layer**: RBAC, ABAC, ACL
3. **Input Validation Layer**: Request/response validation
4. **Security Headers Layer**: CORS, CSRF, XSS protection
5. **Encryption Layer**: TLS, data encryption
6. **Audit Layer**: Comprehensive audit logging

### Security Middleware

```go
type SecurityMiddleware struct {
    authManager    *auth.Manager
    rateLimitManager *ratelimit.Manager
    validator      *validator.Validate
}

func (sm *SecurityMiddleware) SecurityChain() []gin.HandlerFunc {
    return []gin.HandlerFunc{
        sm.CORSMiddleware(),
        sm.SecurityHeadersMiddleware(),
        sm.RateLimitMiddleware(),
        sm.AuthMiddleware(),
        sm.ValidationMiddleware(),
        sm.AuditMiddleware(),
    }
}
```

## Monitoring & Observability

### Three Pillars of Observability

1. **Metrics**: Prometheus, Grafana
2. **Logs**: Structured logging with correlation IDs
3. **Traces**: Jaeger, OpenTelemetry

### Health Checks

```go
type HealthChecker struct {
    managers map[string]Manager
}

func (hc *HealthChecker) CheckHealth(ctx context.Context) *HealthStatus {
    status := &HealthStatus{
        Status: "healthy",
        Checks: make(map[string]CheckResult),
    }
    
    for name, manager := range hc.managers {
        if err := manager.HealthCheck(ctx); err != nil {
            status.Status = "unhealthy"
            status.Checks[name] = CheckResult{
                Status: "unhealthy",
                Error:  err.Error(),
            }
        } else {
            status.Checks[name] = CheckResult{
                Status: "healthy",
            }
        }
    }
    
    return status
}
```

## Performance Architecture

### Connection Pooling

```go
type ConnectionPool struct {
    db     *sql.DB
    redis  *redis.Client
    config *PoolConfig
}

func (cp *ConnectionPool) Initialize() error {
    // Database connection pool
    cp.db.SetMaxOpenConns(cp.config.MaxOpenConns)
    cp.db.SetMaxIdleConns(cp.config.MaxIdleConns)
    cp.db.SetConnMaxLifetime(cp.config.ConnMaxLifetime)
    
    // Redis connection pool
    cp.redis.Options().PoolSize = cp.config.RedisPoolSize
    cp.redis.Options().MinIdleConns = cp.config.RedisMinIdleConns
    
    return nil
}
```

### Caching Strategy

```go
type CacheManager struct {
    providers map[string]CacheProvider
    strategy  CacheStrategy
}

type CacheStrategy interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}, ttl time.Duration) error
    Delete(key string) error
    InvalidatePattern(pattern string) error
}
```

## Error Handling Architecture

### Error Types

```go
type ErrorType int

const (
    ErrorTypeValidation ErrorType = iota
    ErrorTypeAuthentication
    ErrorTypeAuthorization
    ErrorTypeNotFound
    ErrorTypeConflict
    ErrorTypeInternal
    ErrorTypeExternal
)

type FrameworkError struct {
    Type    ErrorType
    Code    string
    Message string
    Details map[string]interface{}
    Cause   error
}
```

### Error Handling Middleware

```go
func ErrorHandlingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            switch e := err.Err.(type) {
            case *FrameworkError:
                c.JSON(e.GetHTTPStatus(), e)
            case *validator.ValidationErrors:
                c.JSON(http.StatusBadRequest, gin.H{
                    "error": "Validation failed",
                    "details": e,
                })
            default:
                c.JSON(http.StatusInternalServerError, gin.H{
                    "error": "Internal server error",
                })
            }
        }
    }
}
```

## Testing Architecture

### Test Types

1. **Unit Tests**: Test individual components
2. **Integration Tests**: Test component interactions
3. **End-to-End Tests**: Test complete workflows
4. **Performance Tests**: Test performance and scalability
5. **Chaos Tests**: Test resilience and fault tolerance

### Test Structure

```go
type TestSuite struct {
    service    Service
    config     *Config
    mocks      map[string]interface{}
    testDB     *sql.DB
    testRedis  *redis.Client
}

func (ts *TestSuite) SetupTest() error {
    // Setup test environment
    // Initialize mocks
    // Setup test database
    // Setup test Redis
}

func (ts *TestSuite) TearDownTest() error {
    // Cleanup test environment
    // Close connections
    // Clean test data
}
```

## Conclusion

Go Micro Framework menggunakan arsitektur yang modular, extensible, dan production-ready. Framework ini mengintegrasikan semua library dari `go-micro-libs` dengan pola Gateway dan Manager yang konsisten, memungkinkan developer untuk fokus pada business logic sambil mendapatkan semua fitur infrastruktur yang diperlukan.

Arsitektur ini memastikan:
- **Scalability**: Mudah menambah provider dan fitur baru
- **Maintainability**: Kode yang clean dan terorganisir
- **Testability**: Mudah untuk testing dan mocking
- **Security**: Built-in security features
- **Observability**: Comprehensive monitoring dan logging
- **Performance**: Optimized untuk production use

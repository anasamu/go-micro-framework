# Go Micro Framework

[![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)](https://github.com/anasamu/go-micro-framework/releases/tag/v1.0.0)
[![Go Version](https://img.shields.io/badge/go-1.24.0-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/anasamu/go-micro-framework)](https://goreportcard.com/report/github.com/anasamu/go-micro-framework)

Framework Go yang powerful dan user-friendly untuk pengembangan microservices dengan integrasi seamless dari 20+ library yang sudah ada di [go-micro-libs](https://github.com/anasamu/go-micro-libs/). Framework ini dirancang untuk memberikan developer experience yang optimal dengan fokus pada business logic, bukan infrastructure.

## 🎉 Version 1.0.0 Released!

This is the first stable release of Go Micro Framework. This release includes a complete framework for microservices development with seamless integration of all 20+ libraries from go-micro-libs.

## 🎯 Vision

Mengembangkan framework Go yang paling powerful dan user-friendly untuk microservices development, dengan integrasi seamless dari semua library yang sudah ada di [go-micro-libs](https://github.com/anasamu/go-micro-libs/).

## 🎯 Mission

- **Zero-Configuration Setup**: Framework menangani semua konfigurasi default
- **Business Logic Focus**: Developer fokus pada business logic, bukan infrastructure
- **Production Ready**: Built-in monitoring, logging, security, dan resilience
- **Extensible**: Mudah menambah fitur baru dan custom providers
- **Library Integration**: Menggunakan semua library yang sudah ada di go-micro-libs

## 📚 Documentation Navigation

### Core Documentation
- **[Framework Overview](docs/FRAMEWORK_OVERVIEW.md)** - Comprehensive overview of the framework
- **[Architecture](docs/ARCHITECTURE.md)** - Detailed architecture documentation
- **[Library Integration](docs/LIBRARY_INTEGRATION.md)** - Complete library integration guide
- **[Library Combination](docs/LIBRARY_COMBINATION.md)** - How to combine multiple libraries

### CLI & Commands
- **[CLI Commands](docs/CLI_COMMANDS.md)** - Complete CLI commands reference
- **[Service Configuration](docs/SERVICE_CONFIGURATION.md)** - Configuration management guide

### Authentication & Security
- **[Authentication](docs/AUTHENTICATION.md)** - Authentication implementation guide
- **[Authorization](docs/AUTHORIZATION.md)** - Authorization and access control
- **[Middleware](docs/MIDDLEWARE.md)** - Middleware implementation and usage

### Communication & Services
- **[Service Communication](docs/SERVICE_COMMUNICATION.md)** - Inter-service communication
- **[gRPC Communication](docs/GRPC_COMMUNICATION.md)** - gRPC implementation guide
- **[Message Broker Communication](docs/MESSAGE_BROKER_COMMUNICATION.md)** - Message broker integration
- **[Service Discovery](docs/SERVICE_DISCOVERY.md)** - Service discovery patterns

### Advanced Features
- **[Circuit Breaker](docs/CIRCUIT_BREAKER.md)** - Circuit breaker implementation
- **[Failover](docs/FAILOVER.md)** - Failover mechanisms and patterns
- **[Deployment](docs/DEPLOYMENT.md)** - Deployment strategies and guides

### Development & Maintenance
- **[API Reference](docs/API.md)** - Complete API documentation
- **[Common Problems](docs/COMMON_PROBLEMS.md)** - Troubleshooting guide
- **[Adding Libraries](docs/ADDING_LIBRARIES.md)** - Guide for adding new libraries

## 🔧 Core Libraries (Always Integrated)

Framework ini secara otomatis mengintegrasikan library berikut di setiap service yang dibuat:

| Library | Description | Providers |
|---------|-------------|-----------|
| **Config** | Configuration management | File, Env, Consul, Vault |
| **Logging** | Structured logging | Console, File, Elasticsearch |
| **Monitoring** | Metrics and tracing | Prometheus, Jaeger, Grafana |
| **Middleware** | HTTP middleware | Auth, Rate Limit, Circuit Breaker |
| **Communication** | Communication protocols | HTTP, gRPC, WebSocket, GraphQL |
| **Utils** | Utility functions | UUID, Environment, Validation |

### Core Library Usage

```go
// Config management
configManager := microservices.NewConfigManager()
configManager.Load()

// Logging
loggingManager := microservices.NewLoggingManager(config, logger)
loggingManager.Initialize()

// Monitoring
monitoringManager := microservices.NewMonitoringManager(config, logger)
monitoringManager.Start()

// Utils
utilsManager := utils.NewUtilsManager("service-name")
serviceID := utilsManager.GetServiceID()
```

## 🏗️ Arsitektur Framework

### Core Components

```
go-micro-framework/
├── cmd/
│   └── microframework/          # CLI binary
├── internal/
│   ├── core/                    # Core framework logic
│   ├── generators/              # Code generation engine
│   ├── templates/               # Go templates for code gen
│   └── validators/              # Configuration validation
├── pkg/                         # Using go-micro-libs package
├── templates/                   # Code generation templates
├── examples/                    # Generated examples
├── docs/                        # Documentation
├── scripts/                     # Build and deployment scripts
└── tests/                       # Integration tests
```

### Framework Features

#### 1. CLI Tool
```bash
microframework new <service-name> [flags]     # Generate new service
microframework add <feature> [flags]          # Add feature to existing service
microframework update [flags]                 # Update framework dependencies
microframework validate [flags]               # Validate service configuration
microframework deploy [flags]                 # Deploy service (Docker/K8s)
microframework logs [flags]                   # View service logs
microframework health [flags]                 # Check service health
microframework config [flags]                 # Manage configuration
```

#### 2. Service Generation
- **REST API**: Standard HTTP REST service
- **GraphQL**: GraphQL API service
- **gRPC**: gRPC service
- **WebSocket**: Real-time WebSocket service
- **Event-Driven**: Event sourcing service
- **Scheduled**: Cron/scheduled task service
- **Worker**: Background job processing service
- **Gateway**: API Gateway service
- **Proxy**: Reverse proxy service

#### 3. Feature Flags
```bash
--type=rest|graphql|grpc|websocket|event|scheduled|worker|gateway|proxy
--with-auth=jwt|oauth|ldap|saml|2fa
--with-db=postgres|mysql|mongodb|redis|elasticsearch|none
--with-messaging=kafka|rabbitmq|nats|sqs|none
--with-cache=redis|memcached|memory|none
--with-storage=s3|gcs|azure|minio|local|none
--with-monitoring=prometheus|jaeger|elasticsearch|datadog|none
--with-discovery=consul|etcd|kubernetes|static|none
--optional=ai,backup,chaos,circuitbreaker,failover,filegen,payment,ratelimit,scheduling
--middleware=auth,logging,monitoring,ratelimit,circuitbreaker,caching
--deployment=docker|kubernetes|helm|terraform
--testing=unit|integration|e2e|benchmark
```

## 📚 Complete Library Integration

Framework ini mengintegrasikan semua library yang sudah ada di [go-micro-libs](https://github.com/anasamu/go-micro-libs/):

### Optional Libraries

| Library | Description | Providers |
|---------|-------------|-----------|
| **AI** | AI services | OpenAI, Anthropic, Google, DeepSeek, X.AI |
| **Auth** | Authentication | JWT, OAuth2, LDAP, SAML, 2FA |
| **Database** | Database abstraction | PostgreSQL, MySQL, MongoDB, Redis, Elasticsearch |
| **Cache** | Caching system | Redis, Memcached, Memory |
| **Storage** | Object storage | S3, GCS, Azure, MinIO |
| **Messaging** | Message queues | Kafka, RabbitMQ, NATS, SQS |
| **Discovery** | Service discovery | Consul, Kubernetes, etcd |
| **Event** | Event sourcing | PostgreSQL, Kafka, NATS |
| **Payment** | Payment processing | Stripe, PayPal, Midtrans, Xendit |
| **Backup** | Backup services | S3, GCS, Local |
| **Chaos** | Chaos engineering | Kubernetes, HTTP, Messaging |
| **Circuit Breaker** | Resilience patterns | Custom, GoBreaker |
| **Failover** | Failover mechanisms | Consul, Kubernetes |
| **FileGen** | File generation | DOCX, Excel, CSV, PDF |
| **Rate Limit** | Rate limiting | Token bucket, Sliding window |
| **Scheduling** | Task scheduling | Cron, Redis-based |
| **API** | Third-party APIs | HTTP, gRPC, GraphQL, WebSocket |
| **Email** | Email services | SMTP, SendGrid, Mailgun |

### Library Usage Examples

#### AI Services
```go
import "github.com/anasamu/go-micro-libs/ai"

// Initialize AI manager
aiManager := ai.NewManager()

// Add OpenAI provider
openaiProvider := openai.NewProvider("your-api-key")
aiManager.AddProvider(openaiProvider)

// Chat with AI
response, err := aiManager.Chat(ctx, "openai", chatReq)
```

#### Database Integration
```go
import "github.com/anasamu/go-micro-libs/database"

// Initialize database manager
dbManager := database.NewManager()

// Add PostgreSQL provider
postgresProvider := postgresql.NewProvider()
dbManager.RegisterProvider(postgresProvider)

// Connect and query
err = dbManager.Connect(ctx, "postgresql")
result, err := dbManager.Query(ctx, "postgresql", "SELECT * FROM users")
```

#### Authentication
```go
import "github.com/anasamu/go-micro-libs/auth"

// Initialize auth manager
authManager := auth.NewManager()

// Add JWT provider
jwtProvider := jwt.NewProvider()
authManager.RegisterProvider(jwtProvider)

// Generate and validate tokens
token, err := authManager.GenerateToken("user123", claims)
claims, err := authManager.ValidateToken(token)
```

## 🚀 Quick Start

### 1. Install CLI Tool
```bash
# Install latest stable version (1.0.0)
go install github.com/anasamu/go-micro-framework/cmd/microframework@v1.0.0

# Or install latest version
go install github.com/anasamu/go-micro-framework/cmd/microframework@latest
```

### 2. Generate New Service
```bash
# Generate basic service (with core libraries automatically integrated)
microframework new user-service

# Generate service with specific features
microframework new order-service \
  --type=rest \
  --with-auth=jwt \
  --with-db=postgres \
  --with-cache=redis \
  --with-monitoring=prometheus \
  --with-messaging=kafka \
  --optional=ai,payment \
  --deployment=docker
```

### 3. Core Libraries Usage

Setiap service yang dibuat sudah memiliki core libraries terintegrasi:

```go
// Config management
configManager := microservices.NewConfigManager()
configManager.Load()

// Logging
loggingManager := microservices.NewLoggingManager(config, logger)
loggingManager.Initialize()

// Monitoring
monitoringManager := microservices.NewMonitoringManager(config, logger)
monitoringManager.Start()

// Utils
utilsManager := utils.NewUtilsManager("service-name")
serviceID := utilsManager.GetServiceID()
```

### 4. Run Service
```bash
cd user-service
go run cmd/main.go
```

## 📝 Generated Service Structure

```
user-service/
├── cmd/
│   └── main.go                 # Bootstrap code
├── internal/
│   ├── handlers/               # HTTP handlers
│   ├── models/                 # Data models
│   ├── repositories/           # Data access layer
│   ├── services/               # Business services
│   └── middleware/             # Custom middleware
├── pkg/
│   └── types/                  # Public types
├── configs/
│   ├── config.yaml             # Configuration
│   ├── config.dev.yaml         # Development config
│   └── config.prod.yaml        # Production config
├── deployments/
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── kubernetes/
│       ├── deployment.yaml
│       └── service.yaml
├── tests/
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── go.mod
├── go.sum
└── README.md
```

## 🔧 Configuration

### Environment Variables
```bash
# Service Configuration
SERVICE_NAME=user-service
SERVICE_VERSION=1.0.0
SERVICE_PORT=8080

# Core Libraries (Always Available)
USER_SERVICE_CONFIG_PATH=./configs
USER_SERVICE_LOG_LEVEL=info
USER_SERVICE_PROMETHEUS_ENDPOINT=http://localhost:9090
USER_SERVICE_JAEGER_ENDPOINT=http://localhost:14268

# Database
DATABASE_URL=postgres://user:pass@localhost/db

# Redis
REDIS_URL=redis://localhost:6379

# Monitoring
PROMETHEUS_ENDPOINT=http://localhost:9090
JAEGER_ENDPOINT=http://localhost:14268

# Authentication
JWT_SECRET=your-jwt-secret

# AI Services
OPENAI_API_KEY=your-openai-key

# Payment
STRIPE_SECRET_KEY=your-stripe-key
```

### Configuration File
```yaml
service:
  name: "user-service"
  version: "1.0.0"
  port: 8080

# Core Libraries Configuration (Always Enabled)
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
    jaeger:
      endpoint: "${JAEGER_ENDPOINT}"

middleware:
  auth:
    enabled: false
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
      timeout: "30s"
    grpc:
      port: 9090
      timeout: "30s"

utils:
  uuid:
    version: 4
    namespace: "user-service"
  environment:
    load_env_file: true
    env_file_path: ".env"

# Optional Libraries (Only if enabled)
database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100

cache:
  providers:
    redis:
      url: "${REDIS_URL}"

auth:
  providers:
    jwt:
      secret: "${JWT_SECRET}"
      expiration: "24h"

monitoring:
  providers:
    prometheus:
      endpoint: "${PROMETHEUS_ENDPOINT}"
    jaeger:
      endpoint: "${JAEGER_ENDPOINT}"

ai:
  providers:
    openai:
      api_key: "${OPENAI_API_KEY}"

payment:
  providers:
    stripe:
      secret_key: "${STRIPE_SECRET_KEY}"
```

## 🧪 Testing

### Unit Tests
```bash
make test-unit
```

### Integration Tests
```bash
make test-integration
```

### End-to-End Tests
```bash
make test-e2e
```

## 🚀 Deployment

### Docker
```bash
make docker-build
make docker-run
```

### Kubernetes
```bash
make k8s-deploy
```

### Helm
```bash
make helm-install
```

## 📊 Monitoring

### Health Check
```bash
curl http://localhost:8080/health
```

### Metrics
```bash
curl http://localhost:8080/metrics
```

### Logs
```bash
microframework logs --service=user-service
```

## 🛠️ CLI Commands Reference

### Command Overview

| Command | Description | Usage |
|---------|-------------|-------|
| `new` | Generate new microservice | `microframework new <service-name> [flags]` |
| `add` | Add features to existing service | `microframework add <feature> [flags]` |
| `generate` | Generate specific components | `microframework generate <type> [flags]` |
| `config` | Manage configuration | `microframework config <subcommand> [flags]` |
| `deploy` | Deploy service | `microframework deploy [flags]` |
| `validate` | Validate service | `microframework validate [flags]` |
| `logs` | View service logs | `microframework logs [flags]` |
| `health` | Check service health | `microframework health [flags]` |
| `update` | Update framework | `microframework update [flags]` |
| `version` | Show version information | `microframework version [flags]` |

### Core Commands

#### `microframework new` - Generate New Service

Generate a new microservice with specified features and configuration.

**Basic Usage:**
```bash
# Generate basic service (with core libraries automatically integrated)
microframework new user-service

# Generate service with specific type
microframework new order-service --type=rest

# Generate service with multiple features
microframework new payment-service \
  --type=rest \
  --with-auth=jwt \
  --with-database=postgres \
  --with-cache=redis \
  --with-monitoring=prometheus
```

**Available Flags:**
| Flag | Description | Options | Default |
|------|-------------|---------|---------|
| `--type`, `-t` | Service type | `rest`, `graphql`, `grpc`, `websocket`, `event`, `scheduled`, `worker`, `gateway`, `proxy` | `rest` |
| `--with-auth` | Include authentication | `jwt`, `oauth`, `ldap`, `saml` | - |
| `--with-database` | Include database | `postgres`, `mysql`, `redis`, `mongodb` | - |
| `--with-messaging` | Include messaging | `kafka`, `rabbitmq`, `nats` | - |
| `--with-monitoring` | Include monitoring | `prometheus`, `jaeger`, `grafana` | - |
| `--with-ai` | Include AI services | `openai`, `anthropic`, `google` | - |
| `--with-storage` | Include storage | `s3`, `gcs`, `azure` | - |
| `--with-cache` | Include caching | `redis`, `memcached`, `memory` | - |
| `--with-discovery` | Include service discovery | `consul`, `kubernetes` | - |
| `--output`, `-o` | Output directory | Path | `.` |
| `--force` | Overwrite existing files | - | `false` |

**Examples:**
```bash
# REST API with authentication and database
microframework new user-service \
  --type=rest \
  --with-auth=jwt \
  --with-database=postgres \
  --with-cache=redis

# GraphQL service with AI integration
microframework new chat-service \
  --type=graphql \
  --with-ai=openai \
  --with-database=postgres \
  --with-cache=redis

# gRPC service with monitoring
microframework new order-service \
  --type=grpc \
  --with-database=postgres \
  --with-monitoring=prometheus \
  --with-messaging=kafka

# Event-driven service
microframework new event-service \
  --type=event \
  --with-database=postgres \
  --with-messaging=kafka \
  --with-event=postgresql

# Worker service with scheduling
microframework new worker-service \
  --type=worker \
  --with-messaging=kafka \
  --with-scheduling=cron \
  --with-monitoring=prometheus

# API Gateway
microframework new api-gateway \
  --type=gateway \
  --with-auth=jwt \
  --with-discovery=consul \
  --with-monitoring=prometheus

# Payment service
microframework new payment-service \
  --type=rest \
  --with-payment=stripe \
  --with-database=postgres \
  --with-monitoring=prometheus

# AI-powered service
microframework new ai-service \
  --type=rest \
  --with-ai=openai \
  --with-database=postgres \
  --with-cache=redis \
  --with-storage=s3
```

#### `microframework add` - Add Features
```bash
# Add AI capabilities
microframework add ai --provider=openai --config=openai.yaml

# Add authentication
microframework add auth --provider=jwt --config=auth.yaml

# Add database
microframework add database --provider=postgresql --config=db.yaml

# Add monitoring
microframework add monitoring --provider=prometheus --config=monitoring.yaml

# Add payment processing
microframework add payment --provider=stripe --config=payment.yaml
```

#### `microframework generate` - Generate Components
```bash
# Generate handler
microframework generate handler user --methods=GET,POST,PUT,DELETE

# Generate service
microframework generate service user --methods=Create,Read,Update,Delete

# Generate repository
microframework generate repository user --database=postgresql

# Generate model
microframework generate model user --fields=id,name,email,created_at

# Generate middleware
microframework generate middleware auth --type=jwt
```

#### `microframework deploy` - Deploy Service
```bash
# Deploy to Docker
microframework deploy --type=docker --config=docker-compose.yml

# Deploy to Kubernetes
microframework deploy --type=kubernetes --namespace=production

# Deploy using Helm
microframework deploy --type=helm --chart=./charts/user-service

# Deploy to cloud
microframework deploy --type=aws --cluster=my-cluster --service=user-service
```

#### `microframework config` - Manage Configuration
```bash
# Show current configuration
microframework config show

# Get specific config value
microframework config get database.url

# Set configuration value
microframework config set database.url "postgres://localhost:5432/mydb"

# Validate configuration
microframework config validate

# Generate configuration template
microframework config generate --template=production
```

#### `microframework validate` - Validate Service
```bash
# Validate all components
microframework validate --type=all

# Validate configuration
microframework validate --type=config

# Validate dependencies
microframework validate --type=dependencies

# Validate with auto-fix
microframework validate --type=all --fix
```

#### `microframework logs` - View Logs
```bash
# View service logs
microframework logs --service=user-service

# Follow logs in real-time
microframework logs --service=user-service --follow

# Filter logs by level
microframework logs --service=user-service --level=error

# View logs from specific time
microframework logs --service=user-service --since=1h
```

#### `microframework health` - Health Checks
```bash
# Check service health
microframework health --service=user-service

# Check all services
microframework health --all

# Detailed health report
microframework health --service=user-service --detailed

# Health check with timeout
microframework health --service=user-service --timeout=30s
```

#### `microframework update` - Update Framework
```bash
# Update framework to latest version
microframework update

# Update specific library
microframework update --library=go-micro-libs

# Update all dependencies
microframework update --all

# Check for updates
microframework update --check
```

#### `microframework version` - Version Information
```bash
# Show framework version
microframework version

# Show detailed version info
microframework version --detailed

# Show library versions
microframework version --libraries
```

## 📚 Advanced Examples

### Complete E-commerce Service
```bash
# Generate complete e-commerce service
microframework new ecommerce-service \
  --type=rest \
  --with-auth=jwt \
  --with-database=postgresql \
  --with-cache=redis \
  --with-messaging=kafka \
  --with-storage=s3 \
  --with-payment=stripe \
  --with-monitoring=prometheus \
  --with-ai=openai \
  --optional=backup,chaos,circuitbreaker \
  --middleware=auth,logging,monitoring,ratelimit \
  --deployment=kubernetes \
  --testing=unit,integration,e2e
```

### Microservices Architecture
```bash
# User Service
microframework new user-service \
  --type=rest \
  --with-auth=jwt \
  --with-database=postgresql \
  --with-cache=redis \
  --with-monitoring=prometheus

# Order Service
microframework new order-service \
  --type=rest \
  --with-database=postgresql \
  --with-messaging=kafka \
  --with-monitoring=prometheus

# Payment Service
microframework new payment-service \
  --type=rest \
  --with-payment=stripe \
  --with-database=postgresql \
  --with-monitoring=prometheus

# Notification Service
microframework new notification-service \
  --type=worker \
  --with-messaging=kafka \
  --with-ai=openai \
  --with-monitoring=prometheus

# API Gateway
microframework new api-gateway \
  --type=gateway \
  --with-auth=jwt \
  --with-monitoring=prometheus \
  --with-discovery=consul
```

### Event-Driven Architecture
```bash
# Event Store
microframework new event-store \
  --type=event \
  --with-database=postgresql \
  --with-event=postgresql \
  --with-monitoring=prometheus

# Command Service
microframework new command-service \
  --type=rest \
  --with-database=postgresql \
  --with-messaging=kafka \
  --with-event=postgresql \
  --with-monitoring=prometheus

# Query Service
microframework new query-service \
  --type=rest \
  --with-database=postgresql \
  --with-cache=redis \
  --with-monitoring=prometheus

# Event Handler
microframework new event-handler \
  --type=worker \
  --with-messaging=kafka \
  --with-database=postgresql \
  --with-monitoring=prometheus
```

## 🔧 Development Workflow

### 1. Project Initialization
```bash
# Create new project
microframework new my-project \
  --type=rest \
  --with-auth=jwt \
  --with-database=postgresql \
  --with-cache=redis \
  --with-monitoring=prometheus

# Navigate to project
cd my-project

# Install dependencies
go mod tidy

# Setup environment
cp .env.example .env
# Edit .env with your configuration
```

### 2. Development
```bash
# Run in development mode
go run cmd/main.go

# Run with hot reload
make dev

# Run tests
make test

# Check code quality
make lint
```

### 3. Testing
```bash
# Unit tests
make test-unit

# Integration tests
make test-integration

# End-to-end tests
make test-e2e

# Coverage report
make test-coverage
```

### 4. Deployment
```bash
# Build for production
make build

# Build Docker image
make docker-build

# Deploy to staging
make deploy-staging

# Deploy to production
make deploy-production
```

## 🏗️ Architecture Patterns

### Clean Architecture
```
user-service/
├── cmd/main.go                 # Application entry point
├── internal/
│   ├── handlers/               # HTTP handlers (presentation layer)
│   │   └── user_handler.go
│   ├── services/               # Business logic (application layer)
│   │   └── user_service.go
│   ├── repositories/           # Data access (infrastructure layer)
│   │   └── user_repository.go
│   ├── models/                 # Domain models (domain layer)
│   │   └── user.go
│   └── middleware/             # Cross-cutting concerns
│       └── auth_middleware.go
├── pkg/types/                  # Public types
└── configs/                    # Configuration files
```

### Domain-Driven Design
```
user-service/
├── cmd/main.go
├── internal/
│   ├── user/                   # User domain
│   │   ├── handlers/
│   │   ├── services/
│   │   ├── repositories/
│   │   └── models/
│   ├── auth/                   # Authentication domain
│   │   ├── handlers/
│   │   ├── services/
│   │   └── models/
│   └── shared/                 # Shared components
│       ├── middleware/
│       └── utils/
```

## 🚀 Performance Optimization

### Connection Pooling
```yaml
# config.yaml
database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100
      max_idle_connections: 10
      connection_max_lifetime: "1h"
```

### Caching Strategy
```yaml
# config.yaml
cache:
  providers:
    redis:
      url: "${REDIS_URL}"
      db: 0
      pool_size: 10
      ttl: "1h"
```

### Monitoring Configuration
```yaml
# config.yaml
monitoring:
  providers:
    prometheus:
      endpoint: "${PROMETHEUS_ENDPOINT}"
      port: 9090
      metrics_path: "/metrics"
    jaeger:
      endpoint: "${JAEGER_ENDPOINT}"
      service_name: "user-service"
      sampling_rate: 0.1
```

## 🔒 Security Best Practices

### Authentication
```yaml
# config.yaml
auth:
  providers:
    jwt:
      secret: "${JWT_SECRET}"
      expiration: "24h"
      issuer: "user-service"
      audience: "api"
    oauth:
      client_id: "${OAUTH_CLIENT_ID}"
      client_secret: "${OAUTH_CLIENT_SECRET}"
      redirect_url: "${OAUTH_REDIRECT_URL}"
      scopes: ["read", "write"]
```

### Security Headers
```yaml
# config.yaml
middleware:
  security:
    enabled: true
    headers:
      - "X-Content-Type-Options: nosniff"
      - "X-Frame-Options: DENY"
      - "X-XSS-Protection: 1; mode=block"
      - "Strict-Transport-Security: max-age=31536000"
```

## 📊 Monitoring & Observability

### Health Checks
```bash
# Basic health check
curl http://localhost:8080/health

# Detailed health check
curl http://localhost:8080/health/detailed

# Readiness check
curl http://localhost:8080/ready

# Liveness check
curl http://localhost:8080/live
```

### Metrics
```bash
# Prometheus metrics
curl http://localhost:8080/metrics

# Custom metrics
curl http://localhost:8080/metrics/custom
```

### Logging
```bash
# View logs
microframework logs --service=user-service

# Filter by level
microframework logs --service=user-service --level=error

# Follow logs
microframework logs --service=user-service --follow
```

## 🧪 Testing Strategies

### Unit Testing
```go
// internal/services/user_service_test.go
func TestUserService_CreateUser(t *testing.T) {
    // Test implementation
}
```

### Integration Testing
```go
// tests/integration/user_test.go
func TestUserIntegration(t *testing.T) {
    // Integration test implementation
}
```

### End-to-End Testing
```go
// tests/e2e/user_e2e_test.go
func TestUserE2E(t *testing.T) {
    // E2E test implementation
}
```

## 🚀 Deployment Strategies

### Docker
```dockerfile
# Generated Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs
CMD ["./main"]
```

### Kubernetes
```yaml
# Generated deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: user-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: user-service-secrets
              key: database-url
```

### Helm
```yaml
# Generated values.yaml
replicaCount: 3
image:
  repository: user-service
  tag: latest
  pullPolicy: IfNotPresent
service:
  type: ClusterIP
  port: 8080
ingress:
  enabled: true
  annotations:
    kubernetes.io/ingress.class: nginx
  hosts:
    - host: user-service.example.com
      paths:
        - path: /
          pathType: Prefix
```

## 🤝 Contributing

### Development Setup
```bash
# Fork and clone repository
git clone https://github.com/anasamu/go-micro-framework.git
cd go-micro-framework

# Install dependencies
go mod tidy

# Install development tools
make tools

# Run tests
make test

# Build framework
make build
```

### Code Style
- Follow Go best practices
- Use `gofmt` for formatting
- Use `golangci-lint` for linting
- Write comprehensive tests
- Document public APIs

### Pull Request Process
1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Create Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **go-micro-libs**: The foundation of this framework
- **Go Community**: For the amazing ecosystem
- **Contributors**: All the amazing people who contribute to this project

## 🚀 Quick Reference

### Common Tasks
- **[Generate New Service](docs/CLI_COMMANDS.md#microframework-new---generate-new-service)** - Create a new microservice
- **[Add Features](docs/CLI_COMMANDS.md#microframework-add---add-features)** - Add features to existing service
- **[Configure Service](docs/SERVICE_CONFIGURATION.md)** - Configure your service
- **[Deploy Service](docs/DEPLOYMENT.md)** - Deploy to various platforms
- **[Troubleshoot Issues](docs/COMMON_PROBLEMS.md)** - Common problems and solutions

### Service Types
- **[REST API](docs/CLI_COMMANDS.md#examples)** - Standard HTTP REST service
- **[GraphQL](docs/CLI_COMMANDS.md#examples)** - GraphQL API service
- **[gRPC](docs/GRPC_COMMUNICATION.md)** - High-performance gRPC service
- **[WebSocket](docs/CLI_COMMANDS.md#examples)** - Real-time WebSocket service
- **[Event-Driven](docs/CLI_COMMANDS.md#examples)** - Event sourcing service
- **[Worker](docs/CLI_COMMANDS.md#examples)** - Background job processing
- **[Gateway](docs/CLI_COMMANDS.md#examples)** - API Gateway service

### Library Integration
- **[Authentication](docs/AUTHENTICATION.md)** - JWT, OAuth2, LDAP, SAML
- **[Database](docs/LIBRARY_INTEGRATION.md#2-database-gomicro-libsdatabase)** - PostgreSQL, MySQL, MongoDB, Redis
- **[AI Services](docs/LIBRARY_INTEGRATION.md#1-ai-services-gomicro-libsai)** - OpenAI, Anthropic, Google, DeepSeek
- **[Messaging](docs/LIBRARY_INTEGRATION.md#5-messaging-gomicro-libsmessaging)** - Kafka, RabbitMQ, NATS, SQS
- **[Storage](docs/LIBRARY_INTEGRATION.md#4-storage-gomicro-libsstorage)** - S3, GCS, Azure, MinIO
- **[Monitoring](docs/LIBRARY_INTEGRATION.md#3-monitoring-gomicro-libsmonitoring)** - Prometheus, Jaeger, Grafana

### Architecture Patterns
- **[Clean Architecture](docs/ARCHITECTURE.md#service-architecture)** - Clean architecture implementation
- **[Domain-Driven Design](docs/ARCHITECTURE.md#service-architecture)** - DDD patterns
- **[Event Sourcing](docs/LIBRARY_INTEGRATION.md#3-event-sourcing-gomicro-libsevent)** - Event-driven architecture
- **[CQRS](docs/LIBRARY_INTEGRATION.md#3-event-sourcing-gomicro-libsevent)** - Command Query Responsibility Segregation
- **[Microservices](docs/ARCHITECTURE.md)** - Microservices patterns

### Security & Middleware
- **[Authentication](docs/AUTHENTICATION.md)** - Authentication implementation
- **[Authorization](docs/AUTHORIZATION.md)** - Authorization and access control
- **[Middleware](docs/MIDDLEWARE.md)** - Middleware implementation
- **[Rate Limiting](docs/MIDDLEWARE.md)** - Rate limiting strategies
- **[Circuit Breaker](docs/CIRCUIT_BREAKER.md)** - Circuit breaker patterns

### Deployment & DevOps
- **[Docker](docs/DEPLOYMENT.md#1-docker-deployment)** - Docker deployment
- **[Kubernetes](docs/DEPLOYMENT.md#2-kubernetes-deployment)** - Kubernetes deployment
- **[Helm](docs/DEPLOYMENT.md#3-helm-deployment)** - Helm charts
- **[Cloud Deployment](docs/DEPLOYMENT.md#4-cloud-deployment)** - AWS, GCP, Azure
- **[CI/CD](docs/DEPLOYMENT.md#5-cicd-integration)** - Continuous integration/deployment

### Development Workflow
- **[Getting Started](docs/FRAMEWORK_OVERVIEW.md#quick-start)** - Quick start guide
- **[Development Setup](docs/FRAMEWORK_OVERVIEW.md#development-workflow)** - Development environment
- **[Testing](docs/FRAMEWORK_OVERVIEW.md#testing-support)** - Testing strategies
- **[Code Generation](docs/CLI_COMMANDS.md#microframework-generate---generate-components)** - Code generation
- **[Configuration Management](docs/SERVICE_CONFIGURATION.md)** - Configuration best practices

---

**Built with ❤️ by the Go Micro Framework team**

**Happy Coding! 🚀**

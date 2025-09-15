# CLI Commands Reference

## 🎯 Overview

GoMicroFramework menyediakan CLI tool yang powerful untuk mengelola microservices. CLI tool ini memungkinkan Anda untuk generate, configure, deploy, dan monitor services dengan mudah.

## 🚀 Installation

### Install CLI Tool

```bash
# Install from source
go install github.com/anasamu/go-micro-framework/cmd/microframework@latest

# Or build from source
git clone https://github.com/anasamu/go-micro-framework.git
cd go-micro-framework
go build -o microframework cmd/microframework/main.go
```

### Verify Installation

```bash
microframework version
```

## 📋 Command Overview

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

## 🔧 Core Commands

### 1. `microframework new` - Generate New Service

Generate a new microservice with specified features and configuration.

#### Basic Usage

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

#### Flags

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
| `--with-circuit-breaker` | Include circuit breaker | - | - |
| `--with-rate-limit` | Include rate limiting | - | - |
| `--with-chaos` | Include chaos engineering | - | - |
| `--with-failover` | Include failover | - | - |
| `--with-event` | Include event sourcing | - | - |
| `--with-scheduling` | Include task scheduling | - | - |
| `--with-backup` | Include backup services | - | - |
| `--with-payment` | Include payment processing | - | - |
| `--with-filegen` | Include file generation | - | - |
| `--with-api` | Include API integration | `http`, `grpc`, `graphql`, `websocket` | - |
| `--with-email` | Include email services | `smtp`, `sendgrid`, `mailgun` | - |
| `--output`, `-o` | Output directory | Path | `.` |
| `--force` | Overwrite existing files | - | `false` |

#### Examples

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

### 2. `microframework add` - Add Features

Add new features to an existing service.

#### Basic Usage

```bash
# Add authentication to existing service
microframework add auth --provider=jwt --config=auth.yaml

# Add database to existing service
microframework add database --provider=postgres --config=db.yaml

# Add AI capabilities
microframework add ai --provider=openai --config=ai.yaml

# Add monitoring
microframework add monitoring --provider=prometheus --config=monitoring.yaml
```

#### Flags

| Flag | Description | Options | Required |
|------|-------------|---------|----------|
| `--provider` | Provider type | Varies by feature | Yes |
| `--config` | Configuration file | Path to config file | No |
| `--force` | Overwrite existing files | - | No |

#### Examples

```bash
# Add JWT authentication
microframework add auth --provider=jwt --config=auth.yaml

# Add PostgreSQL database
microframework add database --provider=postgres --config=db.yaml

# Add Redis cache
microframework add cache --provider=redis --config=cache.yaml

# Add Kafka messaging
microframework add messaging --provider=kafka --config=messaging.yaml

# Add S3 storage
microframework add storage --provider=s3 --config=storage.yaml

# Add OpenAI AI service
microframework add ai --provider=openai --config=ai.yaml

# Add Prometheus monitoring
microframework add monitoring --provider=prometheus --config=monitoring.yaml

# Add Stripe payment
microframework add payment --provider=stripe --config=payment.yaml
```

### 3. `microframework generate` - Generate Components

Generate specific components for a service.

#### Basic Usage

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

#### Component Types

| Type | Description | Flags |
|------|-------------|-------|
| `handler` | HTTP/gRPC handlers | `--methods`, `--type` |
| `service` | Business logic services | `--methods` |
| `repository` | Data access layer | `--database` |
| `model` | Data models | `--fields` |
| `middleware` | HTTP middleware | `--type` |
| `config` | Configuration files | `--template` |
| `test` | Test files | `--type` |

#### Examples

```bash
# Generate REST handler
microframework generate handler user \
  --methods=GET,POST,PUT,DELETE \
  --type=rest

# Generate gRPC handler
microframework generate handler user \
  --methods=Create,Read,Update,Delete \
  --type=grpc

# Generate service with CRUD operations
microframework generate service user \
  --methods=Create,Read,Update,Delete

# Generate repository for PostgreSQL
microframework generate repository user \
  --database=postgresql

# Generate model with fields
microframework generate model user \
  --fields=id,name,email,created_at,updated_at

# Generate JWT middleware
microframework generate middleware auth \
  --type=jwt

# Generate rate limiting middleware
microframework generate middleware ratelimit \
  --type=tokenbucket

# Generate configuration template
microframework generate config \
  --template=production

# Generate unit tests
microframework generate test \
  --type=unit

# Generate integration tests
microframework generate test \
  --type=integration
```

### 4. `microframework config` - Manage Configuration

Manage service configuration.

#### Subcommands

| Subcommand | Description | Usage |
|------------|-------------|-------|
| `show` | Show current configuration | `microframework config show` |
| `get` | Get specific config value | `microframework config get <key>` |
| `set` | Set configuration value | `microframework config set <key> <value>` |
| `validate` | Validate configuration | `microframework config validate` |
| `generate` | Generate configuration template | `microframework config generate` |
| `merge` | Merge configuration files | `microframework config merge <files>` |

#### Examples

```bash
# Show current configuration
microframework config show

# Get specific config value
microframework config get database.url
microframework config get auth.jwt.secret

# Set configuration value
microframework config set database.url "postgres://localhost:5432/mydb"
microframework config set auth.jwt.secret "new-secret"

# Validate configuration
microframework config validate

# Generate configuration template
microframework config generate --template=production

# Merge configuration files
microframework config merge config.yaml config.prod.yaml
```

### 5. `microframework deploy` - Deploy Service

Deploy service to various platforms.

#### Basic Usage

```bash
# Deploy to Docker
microframework deploy --type=docker --config=docker-compose.yml

# Deploy to Kubernetes
microframework deploy --type=kubernetes --namespace=production

# Deploy using Helm
microframework deploy --type=helm --chart=./charts/service

# Deploy to cloud
microframework deploy --type=aws --cluster=my-cluster --service=my-service
```

#### Flags

| Flag | Description | Options | Required |
|------|-------------|---------|----------|
| `--type` | Deployment type | `docker`, `kubernetes`, `helm`, `aws`, `gcp`, `azure` | Yes |
| `--config` | Configuration file | Path to config file | No |
| `--namespace` | Kubernetes namespace | Namespace name | No |
| `--chart` | Helm chart path | Path to chart | No |
| `--cluster` | Cluster name | Cluster name | No |
| `--service` | Service name | Service name | No |
| `--environment` | Environment | `dev`, `staging`, `prod` | No |

#### Examples

```bash
# Deploy to Docker
microframework deploy --type=docker --config=docker-compose.yml

# Deploy to Kubernetes
microframework deploy --type=kubernetes --namespace=production

# Deploy using Helm
microframework deploy --type=helm --chart=./charts/user-service

# Deploy to AWS EKS
microframework deploy --type=aws --cluster=my-cluster --service=user-service

# Deploy to Google GKE
microframework deploy --type=gcp --cluster=my-cluster --service=user-service

# Deploy to Azure AKS
microframework deploy --type=azure --cluster=my-cluster --service=user-service
```

### 6. `microframework validate` - Validate Service

Validate service configuration and dependencies.

#### Basic Usage

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

#### Flags

| Flag | Description | Options | Default |
|------|-------------|---------|---------|
| `--type` | Validation type | `all`, `config`, `dependencies`, `code` | `all` |
| `--fix` | Auto-fix issues | - | `false` |
| `--strict` | Strict validation | - | `false` |

#### Examples

```bash
# Validate everything
microframework validate --type=all

# Validate configuration only
microframework validate --type=config

# Validate dependencies
microframework validate --type=dependencies

# Validate code quality
microframework validate --type=code

# Auto-fix issues
microframework validate --type=all --fix

# Strict validation
microframework validate --type=all --strict
```

### 7. `microframework logs` - View Logs

View and manage service logs.

#### Basic Usage

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

#### Flags

| Flag | Description | Options | Default |
|------|-------------|---------|---------|
| `--service` | Service name | Service name | - |
| `--follow`, `-f` | Follow logs | - | `false` |
| `--level` | Log level | `debug`, `info`, `warn`, `error` | - |
| `--since` | Time since | `1h`, `30m`, `1d` | - |
| `--tail` | Number of lines | Number | `100` |

#### Examples

```bash
# View logs for specific service
microframework logs --service=user-service

# Follow logs in real-time
microframework logs --service=user-service --follow

# View error logs only
microframework logs --service=user-service --level=error

# View logs from last hour
microframework logs --service=user-service --since=1h

# View last 50 lines
microframework logs --service=user-service --tail=50

# Follow error logs
microframework logs --service=user-service --follow --level=error
```

### 8. `microframework health` - Health Checks

Check service health and status.

#### Basic Usage

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

#### Flags

| Flag | Description | Options | Default |
|------|-------------|---------|---------|
| `--service` | Service name | Service name | - |
| `--all` | Check all services | - | `false` |
| `--detailed` | Detailed report | - | `false` |
| `--timeout` | Timeout duration | Duration | `10s` |

#### Examples

```bash
# Check specific service
microframework health --service=user-service

# Check all services
microframework health --all

# Detailed health report
microframework health --service=user-service --detailed

# Health check with timeout
microframework health --service=user-service --timeout=30s

# Check all services with detailed report
microframework health --all --detailed
```

### 9. `microframework update` - Update Framework

Update framework and dependencies.

#### Basic Usage

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

#### Flags

| Flag | Description | Options | Default |
|------|-------------|---------|---------|
| `--library` | Specific library | Library name | - |
| `--all` | Update all dependencies | - | `false` |
| `--check` | Check for updates only | - | `false` |
| `--force` | Force update | - | `false` |

#### Examples

```bash
# Update framework
microframework update

# Update go-micro-libs
microframework update --library=go-micro-libs

# Update all dependencies
microframework update --all

# Check for updates
microframework update --check

# Force update
microframework update --force
```

### 10. `microframework version` - Version Information

Show version information.

#### Basic Usage

```bash
# Show framework version
microframework version

# Show detailed version info
microframework version --detailed

# Show library versions
microframework version --libraries
```

#### Flags

| Flag | Description | Options | Default |
|------|-------------|---------|---------|
| `--detailed` | Detailed information | - | `false` |
| `--libraries` | Show library versions | - | `false` |

#### Examples

```bash
# Show version
microframework version

# Show detailed version
microframework version --detailed

# Show library versions
microframework version --libraries

# Show everything
microframework version --detailed --libraries
```

## 🔧 Advanced Usage

### 1. Service Generation with Multiple Features

```bash
# Complete e-commerce service
microframework new ecommerce-service \
  --type=rest \
  --with-auth=jwt \
  --with-database=postgres \
  --with-cache=redis \
  --with-messaging=kafka \
  --with-storage=s3 \
  --with-payment=stripe \
  --with-monitoring=prometheus \
  --with-ai=openai \
  --with-discovery=consul \
  --with-circuit-breaker \
  --with-rate-limit \
  --with-backup \
  --with-chaos \
  --with-failover \
  --with-event \
  --with-scheduling \
  --with-filegen \
  --with-email \
  --deployment=kubernetes \
  --testing=unit,integration,e2e
```

### 2. Microservices Architecture

```bash
# User Service
microframework new user-service \
  --type=rest \
  --with-auth=jwt \
  --with-database=postgres \
  --with-cache=redis \
  --with-monitoring=prometheus

# Order Service
microframework new order-service \
  --type=rest \
  --with-database=postgres \
  --with-messaging=kafka \
  --with-monitoring=prometheus

# Payment Service
microframework new payment-service \
  --type=rest \
  --with-payment=stripe \
  --with-database=postgres \
  --with-monitoring=prometheus

# Notification Service
microframework new notification-service \
  --type=worker \
  --with-messaging=kafka \
  --with-ai=openai \
  --with-email \
  --with-monitoring=prometheus

# API Gateway
microframework new api-gateway \
  --type=gateway \
  --with-auth=jwt \
  --with-monitoring=prometheus \
  --with-discovery=consul
```

### 3. Event-Driven Architecture

```bash
# Event Store
microframework new event-store \
  --type=event \
  --with-database=postgres \
  --with-event=postgresql \
  --with-monitoring=prometheus

# Command Service
microframework new command-service \
  --type=rest \
  --with-database=postgres \
  --with-messaging=kafka \
  --with-event=postgresql \
  --with-monitoring=prometheus

# Query Service
microframework new query-service \
  --type=rest \
  --with-database=postgres \
  --with-cache=redis \
  --with-monitoring=prometheus

# Event Handler
microframework new event-handler \
  --type=worker \
  --with-messaging=kafka \
  --with-database=postgres \
  --with-monitoring=prometheus
```

## 🔧 Configuration Examples

### 1. Development Environment

```bash
# Generate service for development
microframework new user-service \
  --type=rest \
  --with-database=postgres \
  --with-cache=redis \
  --output=./dev-services
```

### 2. Production Environment

```bash
# Generate service for production
microframework new user-service \
  --type=rest \
  --with-auth=jwt \
  --with-database=postgres \
  --with-cache=redis \
  --with-monitoring=prometheus \
  --with-discovery=consul \
  --with-circuit-breaker \
  --with-rate-limit \
  --with-backup \
  --with-failover \
  --deployment=kubernetes
```

### 3. Testing Environment

```bash
# Generate service for testing
microframework new user-service \
  --type=rest \
  --with-database=postgres \
  --with-cache=redis \
  --testing=unit,integration,e2e
```

## 🔧 Troubleshooting

### Common Issues

1. **Service Name Validation**
   ```bash
   # Valid service names
   microframework new user-service
   microframework new order-service
   microframework new payment-service
   
   # Invalid service names
   microframework new user_service  # Underscore not allowed
   microframework new 123service    # Cannot start with number
   microframework new user-         # Cannot end with hyphen
   ```

2. **Output Directory Issues**
   ```bash
   # Use --force to overwrite existing directory
   microframework new user-service --force
   
   # Specify different output directory
   microframework new user-service --output=./services
   ```

3. **Configuration Issues**
   ```bash
   # Validate configuration
   microframework config validate
   
   # Check configuration
   microframework config show
   ```

### Debug Mode

```bash
# Enable debug mode
export MICROFRAMEWORK_DEBUG=true

# Run command with debug
microframework new user-service --type=rest
```

### Help and Support

```bash
# Get help for specific command
microframework new --help
microframework add --help
microframework generate --help

# Get help for specific subcommand
microframework config --help
microframework deploy --help
```

---

**CLI Commands - Powerful command-line interface for microservices management! 🚀**

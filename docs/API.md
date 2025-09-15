# API Documentation

## Overview

Go Micro Framework menyediakan comprehensive CLI API untuk mengelola microservices development. API ini dirancang untuk memudahkan pengembangan, deployment, dan maintenance microservices menggunakan library `go-micro-libs`.

## CLI Commands

### Core Commands

#### `microframework new` - Generate New Service

Membuat service microservices baru dengan konfigurasi yang ditentukan.

**Syntax:**
```bash
microframework new <service-name> [flags]
```

**Arguments:**
- `service-name` (required): Nama service yang akan dibuat

**Flags:**
```bash
# Service Type
--type=rest|graphql|grpc|websocket|event|scheduled|worker|gateway|proxy

# Core Features
--with-auth=jwt|oauth|ldap|saml|2fa
--with-database=postgres|mysql|mongodb|redis|elasticsearch|none
--with-cache=redis|memcached|memory|none
--with-storage=s3|gcs|azure|minio|local|none
--with-monitoring=prometheus|jaeger|elasticsearch|datadog|none
--with-discovery=consul|etcd|kubernetes|static|none
--with-messaging=kafka|rabbitmq|nats|sqs|none

# Optional Features
--optional=ai,backup,chaos,circuitbreaker,failover,filegen,payment,ratelimit,scheduling

# Middleware
--middleware=auth,logging,monitoring,ratelimit,circuitbreaker,caching

# Deployment
--deployment=docker|kubernetes|helm|terraform

# Testing
--testing=unit|integration|e2e|benchmark

# Output
--output=./output-directory
--force
```

**Examples:**
```bash
# Basic REST service
microframework new user-service

# REST service with authentication and database
microframework new user-service \
  --type=rest \
  --with-auth=jwt \
  --with-database=postgresql \
  --with-cache=redis

# Event-driven service
microframework new order-service \
  --type=event \
  --with-messaging=kafka \
  --with-database=postgresql \
  --with-event=postgresql

# AI-powered service
microframework new chat-service \
  --type=rest \
  --with-ai=openai \
  --with-database=postgresql \
  --with-cache=redis

# Complete e-commerce service
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

**Output:**
```
Creating service: user-service
✓ Generated service structure
✓ Generated main.go
✓ Generated handlers
✓ Generated services
✓ Generated models
✓ Generated repositories
✓ Generated middleware
✓ Generated configuration files
✓ Generated Docker files
✓ Generated Kubernetes manifests
✓ Generated tests
✓ Generated documentation
✓ Generated Makefile
✓ Generated go.mod with dependencies

Service created successfully!
Next steps:
1. cd user-service
2. go mod tidy
3. cp .env.example .env
4. Edit .env with your configuration
5. go run cmd/main.go
```

#### `microframework add` - Add Features

Menambahkan fitur ke service yang sudah ada.

**Syntax:**
```bash
microframework add <feature> [flags]
```

**Arguments:**
- `feature` (required): Fitur yang akan ditambahkan (ai, auth, database, storage, messaging, monitoring, payment, etc.)

**Flags:**
```bash
--provider=<provider-name>
--config=<config-file>
--output=./service-directory
--force
```

**Examples:**
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

Menghasilkan komponen spesifik untuk service.

**Syntax:**
```bash
microframework generate <component-type> <name> [flags]
```

**Arguments:**
- `component-type` (required): Jenis komponen (handler, service, repository, model, middleware)
- `name` (required): Nama komponen

**Flags:**
```bash
# For handlers
--methods=GET,POST,PUT,DELETE
--path=/api/v1/users

# For services
--methods=Create,Read,Update,Delete
--database=postgresql

# For repositories
--database=postgresql
--table=users

# For models
--fields=id,name,email,created_at,updated_at
--database=postgresql

# For middleware
--type=jwt|oauth|rate-limit|circuit-breaker

# Output
--output=./service-directory
--force
```

**Examples:**
```bash
# Generate handler
microframework generate handler user --methods=GET,POST,PUT,DELETE --path=/api/v1/users

# Generate service
microframework generate service user --methods=Create,Read,Update,Delete

# Generate repository
microframework generate repository user --database=postgresql --table=users

# Generate model
microframework generate model user --fields=id,name,email,created_at,updated_at

# Generate middleware
microframework generate middleware auth --type=jwt
```

#### `microframework deploy` - Deploy Service

Deploy service ke environment yang ditentukan.

**Syntax:**
```bash
microframework deploy [flags]
```

**Flags:**
```bash
# Deployment type
--type=docker|kubernetes|helm|aws|gcp|azure

# Environment
--env=development|staging|production

# Configuration
--config=<config-file>
--namespace=<kubernetes-namespace>
--chart=<helm-chart-path>

# Cloud specific
--cluster=<cluster-name>
--service=<service-name>
--region=<region>

# Options
--dry-run
--wait
--timeout=30s
```

**Examples:**
```bash
# Deploy to Docker
microframework deploy --type=docker --config=docker-compose.yml

# Deploy to Kubernetes
microframework deploy --type=kubernetes --namespace=production --env=production

# Deploy using Helm
microframework deploy --type=helm --chart=./charts/user-service --namespace=production

# Deploy to AWS ECS
microframework deploy --type=aws --cluster=my-cluster --service=user-service --region=us-east-1

# Deploy to Google Cloud Run
microframework deploy --type=gcp --service=user-service --region=us-central1

# Deploy to Azure Container Instances
microframework deploy --type=azure --resource-group=my-rg --service=user-service
```

#### `microframework config` - Manage Configuration

Mengelola konfigurasi service.

**Subcommands:**
- `show` - Menampilkan konfigurasi saat ini
- `get <key>` - Mendapatkan nilai konfigurasi
- `set <key> <value>` - Mengatur nilai konfigurasi
- `validate` - Memvalidasi konfigurasi
- `generate` - Generate template konfigurasi

**Syntax:**
```bash
microframework config <subcommand> [flags]
```

**Examples:**
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

Memvalidasi service configuration dan dependencies.

**Syntax:**
```bash
microframework validate [flags]
```

**Flags:**
```bash
--type=all|config|dependencies|templates|code
--fix
--output=<report-file>
```

**Examples:**
```bash
# Validate all components
microframework validate --type=all

# Validate configuration only
microframework validate --type=config

# Validate with auto-fix
microframework validate --type=all --fix

# Generate validation report
microframework validate --type=all --output=validation-report.json
```

#### `microframework logs` - View Logs

Melihat logs dari service.

**Syntax:**
```bash
microframework logs [flags]
```

**Flags:**
```bash
--service=<service-name>
--follow
--level=debug|info|warn|error
--since=<duration>
--until=<duration>
--lines=<number>
```

**Examples:**
```bash
# View service logs
microframework logs --service=user-service

# Follow logs in real-time
microframework logs --service=user-service --follow

# Filter logs by level
microframework logs --service=user-service --level=error

# View logs from specific time
microframework logs --service=user-service --since=1h

# View last 100 lines
microframework logs --service=user-service --lines=100
```

#### `microframework health` - Health Checks

Melakukan health check pada service.

**Syntax:**
```bash
microframework health [flags]
```

**Flags:**
```bash
--service=<service-name>
--all
--detailed
--timeout=<duration>
--interval=<duration>
--count=<number>
```

**Examples:**
```bash
# Check service health
microframework health --service=user-service

# Check all services
microframework health --all

# Detailed health report
microframework health --service=user-service --detailed

# Health check with timeout
microframework health --service=user-service --timeout=30s

# Continuous health check
microframework health --service=user-service --interval=10s --count=5
```

#### `microframework update` - Update Framework

Update framework dan dependencies.

**Syntax:**
```bash
microframework update [flags]
```

**Flags:**
```bash
--library=<library-name>
--all
--check
--force
```

**Examples:**
```bash
# Update framework to latest version
microframework update

# Update specific library
microframework update --library=go-micro-libs

# Update all dependencies
microframework update --all

# Check for updates
microframework update --check

# Force update
microframework update --force
```

#### `microframework version` - Version Information

Menampilkan informasi versi framework.

**Syntax:**
```bash
microframework version [flags]
```

**Flags:**
```bash
--detailed
--libraries
--json
```

**Examples:**
```bash
# Show framework version
microframework version

# Show detailed version info
microframework version --detailed

# Show library versions
microframework version --libraries

# Show version in JSON format
microframework version --json
```

## Configuration API

### Configuration Structure

Framework menggunakan konfigurasi YAML yang sesuai dengan library yang ada:

```yaml
# config.yaml
service:
  name: "user-service"
  version: "1.0.0"
  port: 8080
  environment: "production"

# Core configurations
config:
  providers:
    file:
      path: "./configs"
      format: "yaml"
    env:
      prefix: "USER_SERVICE_"
    consul:
      address: "${CONSUL_ADDRESS}"
      token: "${CONSUL_TOKEN}"

logging:
  providers:
    console:
      level: "info"
      format: "json"
    file:
      path: "/var/log/user-service.log"
      level: "debug"
      max_size: 100
      max_backups: 3
      max_age: 28
    elasticsearch:
      endpoint: "${ELASTICSEARCH_ENDPOINT}"
      index: "user-service-logs"

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
    grafana:
      endpoint: "${GRAFANA_ENDPOINT}"

database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100
      max_idle_connections: 10
      connection_max_lifetime: "1h"
    redis:
      url: "${REDIS_URL}"
      db: 0
      pool_size: 10

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

middleware:
  auth:
    enabled: true
    provider: "jwt"
  rate_limit:
    enabled: true
    requests_per_minute: 100
    burst: 10
  circuit_breaker:
    enabled: true
    failure_threshold: 5
    timeout: 30s
    max_requests: 3

communication:
  providers:
    rest:
      port: 8080
      timeout: 30s
      read_timeout: 10s
      write_timeout: 10s
    grpc:
      port: 9090
      timeout: 30s
      max_recv_msg_size: 4194304
      max_send_msg_size: 4194304

# Optional configurations
ai:
  providers:
    openai:
      api_key: "${OPENAI_API_KEY}"
      default_model: "gpt-4"
      timeout: 30s
    anthropic:
      api_key: "${ANTHROPIC_API_KEY}"
      default_model: "claude-3-sonnet"
      timeout: 30s

storage:
  providers:
    s3:
      access_key: "${AWS_ACCESS_KEY_ID}"
      secret_key: "${AWS_SECRET_ACCESS_KEY}"
      region: "${AWS_REGION}"
      bucket: "${S3_BUCKET}"
    gcs:
      credentials_file: "${GCS_CREDENTIALS_FILE}"
      bucket: "${GCS_BUCKET}"

messaging:
  providers:
    kafka:
      brokers: "${KAFKA_BROKERS}"
      group_id: "user-service"
      topics: ["user-events", "user-commands"]
    rabbitmq:
      url: "${RABBITMQ_URL}"
      exchange: "user-exchange"
      queue: "user-queue"

scheduling:
  providers:
    cron:
      timezone: "UTC"
    redis:
      url: "${REDIS_URL}"
      db: 1

backup:
  providers:
    s3:
      bucket: "${BACKUP_S3_BUCKET}"
      region: "${BACKUP_S3_REGION}"
    gcs:
      bucket: "${BACKUP_GCS_BUCKET}"

chaos:
  providers:
    chaos_monkey:
      enabled: false
      failure_rate: 0.1
      latency: "100ms"

failover:
  providers:
    consul:
      address: "${CONSUL_ADDRESS}"
      service_name: "user-service"
      health_check_interval: "10s"

event:
  providers:
    postgresql:
      url: "${EVENT_POSTGRES_URL}"
      table: "events"
    kafka:
      brokers: "${EVENT_KAFKA_BROKERS}"
      topic: "events"

discovery:
  providers:
    consul:
      address: "${CONSUL_ADDRESS}"
      token: "${CONSUL_TOKEN}"
    kubernetes:
      config_path: "${KUBERNETES_CONFIG}"

cache:
  providers:
    redis:
      url: "${CACHE_REDIS_URL}"
      db: 2
      ttl: "1h"
    memory:
      max_size: 1000
      ttl: "30m"

ratelimit:
  providers:
    redis:
      url: "${RATELIMIT_REDIS_URL}"
      db: 3
    memory:
      requests_per_minute: 100

circuitbreaker:
  providers:
    memory:
      failure_threshold: 5
      timeout: 30s
      max_requests: 3

payment:
  providers:
    stripe:
      secret_key: "${STRIPE_SECRET_KEY}"
      webhook_secret: "${STRIPE_WEBHOOK_SECRET}"
    paypal:
      client_id: "${PAYPAL_CLIENT_ID}"
      client_secret: "${PAYPAL_CLIENT_SECRET}"
      sandbox: true
```

## Generated Service API

### REST API Endpoints

Generated service menyediakan endpoint REST yang konsisten:

```bash
# Health Check
GET /health
GET /health/live
GET /health/ready
GET /health/detailed

# Metrics
GET /metrics

# API Endpoints
GET /api/v1/users
POST /api/v1/users
GET /api/v1/users/{id}
PUT /api/v1/users/{id}
DELETE /api/v1/users/{id}

# Authentication
POST /api/v1/auth/login
POST /api/v1/auth/register
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
```

### gRPC API

Generated service juga menyediakan gRPC API:

```protobuf
service UserService {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}
```

### GraphQL API

Untuk service GraphQL:

```graphql
type Query {
  user(id: ID!): User
  users(first: Int, after: String): UserConnection
}

type Mutation {
  createUser(input: CreateUserInput!): User
  updateUser(id: ID!, input: UpdateUserInput!): User
  deleteUser(id: ID!): Boolean
}

type User {
  id: ID!
  name: String!
  email: String!
  createdAt: DateTime!
  updatedAt: DateTime!
}
```

## Error Handling API

### Error Response Format

```json
{
  "error": {
    "type": "validation_error",
    "code": "INVALID_INPUT",
    "message": "Validation failed",
    "details": {
      "field": "email",
      "reason": "invalid_email_format"
    },
    "timestamp": "2024-01-01T00:00:00Z",
    "request_id": "req_123456789"
  }
}
```

### HTTP Status Codes

- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid request
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Access denied
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict
- `422 Unprocessable Entity` - Validation error
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - Service unavailable

## Monitoring API

### Health Check Response

```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z",
  "service": {
    "name": "user-service",
    "version": "1.0.0"
  },
  "checks": {
    "database": {
      "status": "healthy",
      "response_time": "5ms"
    },
    "redis": {
      "status": "healthy",
      "response_time": "2ms"
    },
    "external_api": {
      "status": "healthy",
      "response_time": "50ms"
    }
  }
}
```

### Metrics Format

Prometheus metrics format:

```
# HTTP requests
http_requests_total{method="GET",path="/api/v1/users",status="200"} 100
http_request_duration_seconds{method="GET",path="/api/v1/users",quantile="0.5"} 0.1
http_request_duration_seconds{method="GET",path="/api/v1/users",quantile="0.9"} 0.2
http_request_duration_seconds{method="GET",path="/api/v1/users",quantile="0.99"} 0.5

# Database connections
database_connections_active{provider="postgresql"} 5
database_connections_idle{provider="postgresql"} 10

# Cache operations
cache_operations_total{provider="redis",operation="get",status="hit"} 1000
cache_operations_total{provider="redis",operation="get",status="miss"} 100
```

## Authentication API

### JWT Token Format

```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "user123",
    "iss": "user-service",
    "aud": "api",
    "exp": 1640995200,
    "iat": 1640908800,
    "roles": ["user", "admin"]
  }
}
```

### OAuth2 Flow

```bash
# Authorization Code Flow
GET /oauth/authorize?client_id=xxx&redirect_uri=xxx&response_type=code&scope=read,write

# Token Exchange
POST /oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code&code=xxx&client_id=xxx&client_secret=xxx&redirect_uri=xxx
```

## Webhook API

### Webhook Payload Format

```json
{
  "event": "user.created",
  "timestamp": "2024-01-01T00:00:00Z",
  "data": {
    "user": {
      "id": "user123",
      "name": "John Doe",
      "email": "john@example.com"
    }
  },
  "metadata": {
    "source": "user-service",
    "version": "1.0.0"
  }
}
```

## Rate Limiting API

### Rate Limit Headers

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
X-RateLimit-Retry-After: 60
```

## Circuit Breaker API

### Circuit Breaker Status

```json
{
  "circuit_breaker": {
    "name": "external_api",
    "state": "closed",
    "failure_count": 0,
    "success_count": 100,
    "failure_threshold": 5,
    "timeout": 30
  }
}
```

## Conclusion

Go Micro Framework menyediakan comprehensive API untuk mengelola seluruh lifecycle microservices development. API ini dirancang untuk memudahkan developer dalam membuat, mengelola, dan deploy microservices dengan integrasi seamless dari semua library yang ada di `go-micro-libs`.

Dengan API ini, developer dapat:
- Generate service dengan konfigurasi yang tepat
- Mengelola konfigurasi dengan mudah
- Deploy ke berbagai platform
- Monitor dan maintain service
- Integrate dengan semua library yang tersedia

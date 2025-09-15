# Service Configuration Guide

## 🎯 Overview

GoMicroFramework menyediakan sistem konfigurasi yang powerful dan fleksibel untuk mengelola konfigurasi service. Framework mendukung multiple configuration sources, hot reloading, dan validation.

## 🔧 Configuration Sources

### 1. Environment Variables
```bash
# Service Configuration
SERVICE_NAME=user-service
SERVICE_VERSION=1.0.0
SERVICE_PORT=8080

# Database
DATABASE_URL=postgres://user:pass@localhost/db
REDIS_URL=redis://localhost:6379

# Monitoring
PROMETHEUS_ENDPOINT=http://localhost:9090
JAEGER_ENDPOINT=http://localhost:14268

# Authentication
JWT_SECRET=your-jwt-secret
```

### 2. Configuration Files
```yaml
# config.yaml
service:
  name: "user-service"
  version: "1.0.0"
  port: 8080

# Core libraries (always enabled)
config:
  providers:
    file:
      path: "./configs"
      format: "yaml"
    env:
      prefix: "SERVICE_"

logging:
  providers:
    console:
      level: "info"
      format: "json"
    file:
      path: "/var/log/service.log"
      level: "debug"

monitoring:
  providers:
    prometheus:
      endpoint: ":9090"
    jaeger:
      endpoint: "http://localhost:14268"
      service_name: "user-service"
```

### 3. External Sources
- Consul
- Vault
- etcd
- Kubernetes ConfigMaps

## 📁 Configuration Structure

### Generated Configuration Files

```
service-name/
├── configs/
│   ├── config.yaml             # Default configuration
│   ├── config.dev.yaml         # Development config
│   ├── config.staging.yaml     # Staging config
│   ├── config.prod.yaml        # Production config
│   └── schema.yaml             # Configuration schema
├── .env.example                # Environment variables template
└── .env                        # Local environment variables
```

### Configuration Hierarchy

```
Command Line Flags > Environment Variables > Config Files > Defaults
```

## 🔧 Core Configuration

### Service Configuration

```yaml
# config.yaml
service:
  name: "user-service"
  version: "1.0.0"
  port: 8080
  environment: "development"
  region: "us-east-1"
  zone: "us-east-1a"
  
  # Health check configuration
  health:
    enabled: true
    path: "/health"
    port: 8081
    
  # Graceful shutdown
  shutdown:
    timeout: "30s"
    signals: ["SIGTERM", "SIGINT"]
```

### Config Management

```yaml
# config.yaml
config:
  providers:
    file:
      path: "./configs"
      format: "yaml"
      watch: true
    env:
      prefix: "SERVICE_"
      required: true
    consul:
      address: "localhost:8500"
      token: "${CONSUL_TOKEN}"
      key_prefix: "services/user-service"
    vault:
      address: "http://localhost:8200"
      token: "${VAULT_TOKEN}"
      path: "secret/services/user-service"
      
  # Configuration validation
  validation:
    enabled: true
    schema_path: "./configs/schema.yaml"
    
  # Hot reloading
  hot_reload:
    enabled: true
    interval: "5s"
```

## 🔧 Library Configurations

### Logging Configuration

```yaml
# config.yaml
logging:
  providers:
    console:
      enabled: true
      level: "info"
      format: "json"
      timestamp: true
      caller: true
      
    file:
      enabled: true
      path: "/var/log/user-service.log"
      level: "debug"
      max_size: 100
      max_backups: 3
      max_age: 28
      compress: true
      
    elasticsearch:
      enabled: false
      endpoint: "http://localhost:9200"
      index: "user-service-logs"
      level: "info"
      
  # Correlation ID
  correlation:
    enabled: true
    header: "X-Correlation-ID"
    generate: true
```

### Monitoring Configuration

```yaml
# config.yaml
monitoring:
  providers:
    prometheus:
      enabled: true
      endpoint: ":9090"
      path: "/metrics"
      namespace: "user_service"
      
    jaeger:
      enabled: true
      endpoint: "http://localhost:14268"
      service_name: "user-service"
      sampling_rate: 0.1
      tags:
        environment: "development"
        version: "1.0.0"
        
    grafana:
      enabled: false
      endpoint: "http://localhost:3000"
      dashboard_id: "user-service"
      
  # Health checks
  health:
    enabled: true
    path: "/health"
    port: 8081
    checks:
      - name: "database"
        type: "database"
        timeout: "5s"
      - name: "cache"
        type: "cache"
        timeout: "3s"
      - name: "external-api"
        type: "http"
        url: "http://external-api/health"
        timeout: "10s"
```

### Middleware Configuration

```yaml
# config.yaml
middleware:
  auth:
    enabled: true
    provider: "jwt"
    secret: "${JWT_SECRET}"
    expiration: "24h"
    issuer: "user-service"
    audience: "api"
    
  rate_limit:
    enabled: true
    provider: "redis"
    requests_per_minute: 100
    burst: 10
    key_func: "ip"
    
  circuit_breaker:
    enabled: true
    provider: "gobreaker"
    failure_threshold: 5
    timeout: "30s"
    max_requests: 3
    
  caching:
    enabled: true
    provider: "redis"
    ttl: "1h"
    key_prefix: "user-service:"
    
  security:
    enabled: true
    headers:
      - "X-Content-Type-Options: nosniff"
      - "X-Frame-Options: DENY"
      - "X-XSS-Protection: 1; mode=block"
      - "Strict-Transport-Security: max-age=31536000"
    cors:
      enabled: true
      origins: ["http://localhost:3000"]
      methods: ["GET", "POST", "PUT", "DELETE"]
      headers: ["Content-Type", "Authorization"]
```

### Communication Configuration

```yaml
# config.yaml
communication:
  providers:
    http:
      enabled: true
      port: 8080
      timeout: "30s"
      read_timeout: "10s"
      write_timeout: "10s"
      idle_timeout: "120s"
      max_header_bytes: 1048576
      
    grpc:
      enabled: true
      port: 9090
      timeout: "30s"
      max_recv_msg_size: 4194304
      max_send_msg_size: 4194304
      keepalive:
        time: "30s"
        timeout: "5s"
        permit_without_stream: true
        
    websocket:
      enabled: false
      port: 8082
      path: "/ws"
      read_buffer_size: 1024
      write_buffer_size: 1024
      check_origin: true
      
    graphql:
      enabled: false
      port: 8083
      path: "/graphql"
      playground: true
      introspection: true
```

## 🔧 Optional Library Configurations

### Database Configuration

```yaml
# config.yaml
database:
  providers:
    postgresql:
      enabled: true
      url: "${DATABASE_URL}"
      max_connections: 100
      max_idle_connections: 10
      connection_max_lifetime: "1h"
      connection_max_idle_time: "30m"
      
    redis:
      enabled: true
      url: "${REDIS_URL}"
      db: 0
      pool_size: 10
      min_idle_conns: 5
      max_conn_age: "1h"
      pool_timeout: "4s"
      idle_timeout: "5m"
      
    mongodb:
      enabled: false
      url: "${MONGODB_URL}"
      database: "user-service"
      max_pool_size: 100
      min_pool_size: 10
      max_idle_time: "30m"
```

### Cache Configuration

```yaml
# config.yaml
cache:
  providers:
    redis:
      enabled: true
      url: "${REDIS_URL}"
      db: 1
      pool_size: 10
      ttl: "1h"
      key_prefix: "cache:"
      
    memory:
      enabled: false
      max_size: 1000
      ttl: "30m"
      cleanup_interval: "5m"
      
    memcached:
      enabled: false
      servers: ["localhost:11211"]
      ttl: "1h"
      max_idle_conns: 2
      timeout: "100ms"
```

### Storage Configuration

```yaml
# config.yaml
storage:
  providers:
    s3:
      enabled: true
      region: "us-east-1"
      access_key_id: "${AWS_ACCESS_KEY_ID}"
      secret_access_key: "${AWS_SECRET_ACCESS_KEY}"
      bucket: "${S3_BUCKET}"
      endpoint: ""
      force_path_style: false
      
    gcs:
      enabled: false
      credentials_file: "${GCS_CREDENTIALS_FILE}"
      bucket: "${GCS_BUCKET}"
      project_id: "${GCP_PROJECT_ID}"
      
    azure:
      enabled: false
      account_name: "${AZURE_ACCOUNT_NAME}"
      account_key: "${AZURE_ACCOUNT_KEY}"
      container: "${AZURE_CONTAINER}"
```

### Messaging Configuration

```yaml
# config.yaml
messaging:
  providers:
    kafka:
      enabled: true
      brokers: ["localhost:9092"]
      group_id: "user-service"
      topics: ["users", "orders"]
      consumer:
        auto_offset_reset: "latest"
        enable_auto_commit: true
        auto_commit_interval: "1s"
      producer:
        required_acks: 1
        timeout: "10s"
        retry_max: 3
        
    rabbitmq:
      enabled: false
      url: "amqp://guest:guest@localhost:5672/"
      exchange: "user-exchange"
      queue: "user-queue"
      routing_key: "user.created"
      
    nats:
      enabled: false
      url: "nats://localhost:4222"
      cluster_id: "user-service"
      client_id: "user-service-1"
```

### Authentication Configuration

```yaml
# config.yaml
auth:
  providers:
    jwt:
      enabled: true
      secret: "${JWT_SECRET}"
      expiration: "24h"
      issuer: "user-service"
      audience: "api"
      algorithm: "HS256"
      
    oauth:
      enabled: false
      client_id: "${OAUTH_CLIENT_ID}"
      client_secret: "${OAUTH_CLIENT_SECRET}"
      redirect_url: "${OAUTH_REDIRECT_URL}"
      scopes: ["read", "write"]
      provider: "google"
      
    ldap:
      enabled: false
      server: "ldap://localhost:389"
      base_dn: "dc=example,dc=com"
      bind_dn: "cn=admin,dc=example,dc=com"
      bind_password: "${LDAP_PASSWORD}"
```

### AI Configuration

```yaml
# config.yaml
ai:
  providers:
    openai:
      enabled: true
      api_key: "${OPENAI_API_KEY}"
      base_url: "https://api.openai.com/v1"
      default_model: "gpt-4"
      timeout: "30s"
      max_tokens: 1000
      
    anthropic:
      enabled: false
      api_key: "${ANTHROPIC_API_KEY}"
      default_model: "claude-3-sonnet"
      timeout: "30s"
      max_tokens: 1000
      
    google:
      enabled: false
      api_key: "${GOOGLE_API_KEY}"
      default_model: "gemini-pro"
      timeout: "30s"
```

### Payment Configuration

```yaml
# config.yaml
payment:
  providers:
    stripe:
      enabled: true
      secret_key: "${STRIPE_SECRET_KEY}"
      publishable_key: "${STRIPE_PUBLISHABLE_KEY}"
      webhook_secret: "${STRIPE_WEBHOOK_SECRET}"
      currency: "usd"
      
    paypal:
      enabled: false
      client_id: "${PAYPAL_CLIENT_ID}"
      client_secret: "${PAYPAL_CLIENT_SECRET}"
      environment: "sandbox"
      currency: "USD"
```

## 🔧 Environment-Specific Configurations

### Development Configuration

```yaml
# config.dev.yaml
service:
  environment: "development"
  port: 8080

logging:
  providers:
    console:
      level: "debug"
      format: "text"
    file:
      enabled: false

monitoring:
  providers:
    prometheus:
      enabled: true
    jaeger:
      enabled: true
      sampling_rate: 1.0

database:
  providers:
    postgresql:
      url: "postgres://localhost:5432/user_service_dev"
    redis:
      url: "redis://localhost:6379/0"

middleware:
  auth:
    enabled: false
  rate_limit:
    enabled: false
  circuit_breaker:
    enabled: false
```

### Staging Configuration

```yaml
# config.staging.yaml
service:
  environment: "staging"
  port: 8080

logging:
  providers:
    console:
      level: "info"
      format: "json"
    file:
      enabled: true
      level: "info"

monitoring:
  providers:
    prometheus:
      enabled: true
    jaeger:
      enabled: true
      sampling_rate: 0.1

database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
    redis:
      url: "${REDIS_URL}"

middleware:
  auth:
    enabled: true
  rate_limit:
    enabled: true
    requests_per_minute: 50
  circuit_breaker:
    enabled: true
```

### Production Configuration

```yaml
# config.prod.yaml
service:
  environment: "production"
  port: 8080

logging:
  providers:
    console:
      level: "warn"
      format: "json"
    file:
      enabled: true
      level: "info"
    elasticsearch:
      enabled: true

monitoring:
  providers:
    prometheus:
      enabled: true
    jaeger:
      enabled: true
      sampling_rate: 0.01
    grafana:
      enabled: true

database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 200
      max_idle_connections: 20
    redis:
      url: "${REDIS_URL}"
      pool_size: 20

middleware:
  auth:
    enabled: true
  rate_limit:
    enabled: true
    requests_per_minute: 1000
  circuit_breaker:
    enabled: true
  security:
    enabled: true
```

## 🔧 Configuration Validation

### Schema Definition

```yaml
# schema.yaml
type: object
properties:
  service:
    type: object
    properties:
      name:
        type: string
        minLength: 3
        maxLength: 50
      version:
        type: string
        pattern: "^[0-9]+\\.[0-9]+\\.[0-9]+$"
      port:
        type: integer
        minimum: 1
        maximum: 65535
    required: ["name", "version", "port"]
    
  logging:
    type: object
    properties:
      providers:
        type: object
        properties:
          console:
            type: object
            properties:
              level:
                type: string
                enum: ["debug", "info", "warn", "error"]
              format:
                type: string
                enum: ["text", "json"]
                
  database:
    type: object
    properties:
      providers:
        type: object
        properties:
          postgresql:
            type: object
            properties:
              url:
                type: string
                format: uri
              max_connections:
                type: integer
                minimum: 1
                maximum: 1000
                
required: ["service", "logging"]
```

### Validation Commands

```bash
# Validate configuration
microframework config validate

# Validate with schema
microframework config validate --schema=./configs/schema.yaml

# Validate specific environment
microframework config validate --env=production
```

## 🔧 Hot Reloading

### Configuration Hot Reload

```yaml
# config.yaml
config:
  hot_reload:
    enabled: true
    interval: "5s"
    watch_files: true
    watch_directories: ["./configs"]
    
  providers:
    file:
      watch: true
      debounce: "1s"
```

### Implementation

```go
// Hot reload implementation
func (c *ConfigManager) WatchConfig() {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()
    
    go func() {
        for {
            select {
            case event := <-watcher.Events:
                if event.Op&fsnotify.Write == fsnotify.Write {
                    c.ReloadConfig()
                }
            case err := <-watcher.Errors:
                log.Error("Config watcher error:", err)
            }
        }
    }()
    
    err = watcher.Add("./configs")
    if err != nil {
        log.Fatal(err)
    }
}
```

## 🔧 Configuration Management

### CLI Commands

```bash
# Show current configuration
microframework config show

# Get specific value
microframework config get database.url

# Set configuration value
microframework config set database.url "postgres://localhost:5432/mydb"

# Merge configurations
microframework config merge config.yaml config.prod.yaml

# Generate configuration template
microframework config generate --template=production
```

### Programmatic Access

```go
// Access configuration programmatically
config := configManager.GetConfig()

// Get specific values
dbURL := config.GetString("database.url")
port := config.GetInt("service.port")
enabled := config.GetBool("middleware.auth.enabled")

// Get nested values
loggingLevel := config.GetString("logging.providers.console.level")

// Set values
config.Set("database.url", "postgres://localhost:5432/mydb")
config.Set("service.port", 8080)
```

## 🔧 Best Practices

### 1. Configuration Organization

```yaml
# Organize configuration by environment
configs/
├── config.yaml          # Base configuration
├── config.dev.yaml      # Development overrides
├── config.staging.yaml  # Staging overrides
├── config.prod.yaml     # Production overrides
└── schema.yaml          # Validation schema
```

### 2. Secret Management

```yaml
# Use environment variables for secrets
database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"  # Never hardcode secrets
      
auth:
  providers:
    jwt:
      secret: "${JWT_SECRET}"  # Use environment variables
```

### 3. Configuration Validation

```yaml
# Always validate configuration
config:
  validation:
    enabled: true
    schema_path: "./configs/schema.yaml"
    strict: true
```

### 4. Environment-Specific Settings

```yaml
# Use environment-specific configurations
service:
  environment: "${ENVIRONMENT:-development}"
  
logging:
  providers:
    console:
      level: "${LOG_LEVEL:-info}"
      
monitoring:
  providers:
    jaeger:
      sampling_rate: "${JAEGER_SAMPLING_RATE:-0.1}"
```

---

**Service Configuration - Flexible and powerful configuration management! 🚀**

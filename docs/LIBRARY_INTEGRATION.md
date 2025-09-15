# Library Integration Guide

## 🎯 Overview

GoMicroFramework mengintegrasikan semua library yang sudah ada di [go-micro-libs](https://github.com/anasamu/go-micro-libs/) untuk memberikan pengalaman development yang seamless. Setiap service yang dibuat secara otomatis menggunakan library yang sesuai dengan konfigurasi yang dipilih.

## 📚 Library Architecture

### Manager-Provider Pattern

Setiap library mengikuti pola arsitektur yang konsisten:

```
Manager (Interface)
    ↓
Provider (Implementation)
    ↓
External Service (Database, Cache, etc.)
```

### Core Components

1. **Manager**: Mengelola multiple provider dan menyediakan interface terpadu
2. **Provider**: Implementasi spesifik untuk setiap layanan eksternal
3. **Types**: Definisi tipe data dan interface yang digunakan
4. **Configuration**: Konfigurasi untuk setiap provider

## 🔧 Core Libraries (Always Integrated)

### 1. Config Management (`go-micro-libs/config`)

**Purpose**: Manajemen konfigurasi multi-source dengan hot reloading

**Providers**:
- File (YAML, JSON, TOML)
- Environment Variables
- Consul
- Vault

**Usage**:
```go
import (
    "github.com/anasamu/go-micro-libs/config"
    "github.com/anasamu/go-micro-libs/config/providers/file"
    "github.com/anasamu/go-micro-libs/config/providers/env"
)

// Initialize config manager
configManager := config.NewManager()

// Add providers
fileProvider := file.NewProvider("./configs")
envProvider := env.NewProvider("SERVICE_")

configManager.AddProvider(fileProvider)
configManager.AddProvider(envProvider)

// Load configuration
err := configManager.Load()
if err != nil {
    log.Fatal(err)
}

// Get configuration value
dbURL := configManager.GetString("database.url")
```

**Configuration**:
```yaml
# config.yaml
config:
  providers:
    file:
      path: "./configs"
      format: "yaml"
    env:
      prefix: "SERVICE_"
    consul:
      address: "localhost:8500"
      token: "your-token"
```

### 2. Logging (`go-micro-libs/logging`)

**Purpose**: Structured logging dengan multiple providers

**Providers**:
- Console
- File
- Elasticsearch

**Usage**:
```go
import (
    "github.com/anasamu/go-micro-libs/logging"
    "github.com/anasamu/go-micro-libs/logging/providers/console"
    "github.com/anasamu/go-micro-libs/logging/providers/file"
)

// Initialize logging manager
loggingManager := logging.NewManager()

// Add providers
consoleProvider := console.NewProvider()
fileProvider := file.NewProvider("/var/log/service.log")

loggingManager.AddProvider(consoleProvider)
loggingManager.AddProvider(fileProvider)

// Initialize logging
err := loggingManager.Initialize()
if err != nil {
    log.Fatal(err)
}

// Use logger
logger := loggingManager.GetLogger()
logger.Info("Service started", "port", 8080)
```

**Configuration**:
```yaml
# config.yaml
logging:
  providers:
    console:
      level: "info"
      format: "json"
    file:
      path: "/var/log/service.log"
      level: "debug"
      max_size: 100
      max_backups: 3
      max_age: 28
```

### 3. Monitoring (`go-micro-libs/monitoring`)

**Purpose**: Metrics, tracing, dan health checks

**Providers**:
- Prometheus
- Jaeger
- Elasticsearch

**Usage**:
```go
import (
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/monitoring/providers/prometheus"
    "github.com/anasamu/go-micro-libs/monitoring/providers/jaeger"
)

// Initialize monitoring manager
monitoringManager := monitoring.NewManager()

// Add providers
prometheusProvider := prometheus.NewProvider(":9090")
jaegerProvider := jaeger.NewProvider("http://localhost:14268")

monitoringManager.AddProvider(prometheusProvider)
monitoringManager.AddProvider(jaegerProvider)

// Start monitoring
err := monitoringManager.Start()
if err != nil {
    log.Fatal(err)
}

// Create metrics
counter := monitoringManager.CreateCounter("requests_total", "Total requests")
histogram := monitoringManager.CreateHistogram("request_duration", "Request duration")

// Use metrics
counter.Inc()
histogram.Observe(0.5)
```

**Configuration**:
```yaml
# config.yaml
monitoring:
  providers:
    prometheus:
      endpoint: ":9090"
      path: "/metrics"
    jaeger:
      endpoint: "http://localhost:14268"
      service_name: "my-service"
      sampling_rate: 0.1
```

### 4. Middleware (`go-micro-libs/middleware`)

**Purpose**: HTTP middleware untuk authentication, rate limiting, dll

**Providers**:
- Authentication
- Rate Limiting
- Circuit Breaker
- Caching
- Logging
- Monitoring
- Security

**Usage**:
```go
import (
    "github.com/anasamu/go-micro-libs/middleware"
    "github.com/anasamu/go-micro-libs/middleware/providers/auth"
    "github.com/anasamu/go-micro-libs/middleware/providers/ratelimit"
)

// Initialize middleware manager
middlewareManager := middleware.NewManager()

// Add middleware
authMiddleware := auth.NewMiddleware("jwt")
rateLimitMiddleware := ratelimit.NewMiddleware(100, time.Minute)

middlewareManager.AddMiddleware(authMiddleware)
middlewareManager.AddMiddleware(rateLimitMiddleware)

// Apply middleware
http.Handle("/api/", middlewareManager.Apply(handler))
```

**Configuration**:
```yaml
# config.yaml
middleware:
  auth:
    enabled: true
    provider: "jwt"
    secret: "${JWT_SECRET}"
  rate_limit:
    enabled: true
    requests_per_minute: 100
    burst: 10
  circuit_breaker:
    enabled: true
    failure_threshold: 5
    timeout: 30s
```

### 5. Communication (`go-micro-libs/communication`)

**Purpose**: Protokol komunikasi (HTTP, gRPC, WebSocket, GraphQL)

**Providers**:
- HTTP/REST
- gRPC
- WebSocket
- GraphQL
- SSE (Server-Sent Events)
- QUIC

**Usage**:
```go
import (
    "github.com/anasamu/go-micro-libs/communication"
    "github.com/anasamu/go-micro-libs/communication/providers/http"
    "github.com/anasamu/go-micro-libs/communication/providers/grpc"
)

// Initialize communication manager
commManager := communication.NewManager()

// Add providers
httpProvider := http.NewProvider(":8080")
grpcProvider := grpc.NewProvider(":9090")

commManager.AddProvider(httpProvider)
commManager.AddProvider(grpcProvider)

// Start communication
err := commManager.Start()
if err != nil {
    log.Fatal(err)
}
```

**Configuration**:
```yaml
# config.yaml
communication:
  providers:
    http:
      port: 8080
      timeout: 30s
      read_timeout: 10s
      write_timeout: 10s
    grpc:
      port: 9090
      timeout: 30s
      max_recv_msg_size: 4194304
      max_send_msg_size: 4194304
```

## 🔧 Optional Libraries

### 1. AI Services (`go-micro-libs/ai`)

**Purpose**: Integrasi dengan berbagai AI services

**Providers**:
- OpenAI
- Anthropic
- Google
- DeepSeek
- X.AI

**Usage**:
```go
import (
    "github.com/anasamu/go-micro-libs/ai"
    "github.com/anasamu/go-micro-libs/ai/providers/openai"
    "github.com/anasamu/go-micro-libs/ai/types"
)

// Initialize AI manager
aiManager := ai.NewManager()

// Add OpenAI provider
openaiProvider := openai.NewProvider("your-api-key")
aiManager.AddProvider(openaiProvider)

// Chat with AI
ctx := context.Background()
chatReq := &types.ChatRequest{
    Messages: []types.Message{
        {Role: "user", Content: "Hello, how are you?"},
    },
    Model: "gpt-4",
}

response, err := aiManager.Chat(ctx, "openai", chatReq)
if err != nil {
    log.Fatal(err)
}

fmt.Println("AI Response:", response.Choices[0].Message.Content)
```

**Configuration**:
```yaml
# config.yaml
ai:
  providers:
    openai:
      api_key: "${OPENAI_API_KEY}"
      base_url: "https://api.openai.com/v1"
      default_model: "gpt-4"
      timeout: 30s
    anthropic:
      api_key: "${ANTHROPIC_API_KEY}"
      default_model: "claude-3-sonnet"
      timeout: 30s
```

### 2. Database (`go-micro-libs/database`)

**Purpose**: Database abstraction dengan multiple providers

**Providers**:
- PostgreSQL
- MySQL
- MongoDB
- Redis
- SQLite
- Cassandra
- CockroachDB
- Elasticsearch
- InfluxDB

**Usage**:
```go
import (
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/database/providers/postgresql"
)

// Initialize database manager
dbManager := database.NewManager()

// Add PostgreSQL provider
postgresProvider := postgresql.NewProvider()
config := map[string]interface{}{
    "host":     "localhost",
    "port":     5432,
    "user":     "postgres",
    "password": "password",
    "database": "mydb",
}

err := postgresProvider.Configure(config)
if err != nil {
    log.Fatal(err)
}

dbManager.RegisterProvider(postgresProvider)

// Connect to database
ctx := context.Background()
err = dbManager.Connect(ctx, "postgresql")
if err != nil {
    log.Fatal(err)
}

// Execute query
result, err := dbManager.Query(ctx, "postgresql", "SELECT * FROM users LIMIT 10")
if err != nil {
    log.Fatal(err)
}
```

**Configuration**:
```yaml
# config.yaml
database:
  providers:
    postgresql:
      host: "localhost"
      port: 5432
      user: "postgres"
      password: "${DB_PASSWORD}"
      database: "mydb"
      max_connections: 100
      max_idle_connections: 10
      connection_max_lifetime: "1h"
    redis:
      host: "localhost"
      port: 6379
      password: "${REDIS_PASSWORD}"
      db: 0
      pool_size: 10
```

### 3. Cache (`go-micro-libs/cache`)

**Purpose**: Caching system dengan fallback

**Providers**:
- Redis
- Memcached
- Memory

**Usage**:
```go
import (
    "github.com/anasamu/go-micro-libs/cache"
    "github.com/anasamu/go-micro-libs/cache/providers/redis"
)

// Initialize cache manager
cacheManager := cache.NewManager()

// Add Redis provider
redisProvider := redis.NewProvider()
config := map[string]interface{}{
    "host": "localhost",
    "port": 6379,
    "db":   0,
}

err := redisProvider.Configure(config)
if err != nil {
    log.Fatal(err)
}

cacheManager.RegisterProvider(redisProvider)

// Connect to cache
ctx := context.Background()
err = cacheManager.Connect(ctx, "redis")
if err != nil {
    log.Fatal(err)
}

// Set cache
err = cacheManager.Set(ctx, "key", "value", 10*time.Minute)
if err != nil {
    log.Fatal(err)
}

// Get cache
var value string
err = cacheManager.Get(ctx, "key", &value)
if err != nil {
    log.Fatal(err)
}
```

**Configuration**:
```yaml
# config.yaml
cache:
  providers:
    redis:
      host: "localhost"
      port: 6379
      password: "${REDIS_PASSWORD}"
      db: 0
      pool_size: 10
      ttl: "1h"
    memory:
      max_size: 1000
      ttl: "30m"
```

### 4. Storage (`go-micro-libs/storage`)

**Purpose**: Object storage abstraction

**Providers**:
- AWS S3
- Google Cloud Storage
- Azure Blob Storage
- MinIO

**Usage**:
```go
import (
    "github.com/anasamu/go-micro-libs/storage"
    "github.com/anasamu/go-micro-libs/storage/providers/s3"
)

// Initialize storage manager
storageManager := storage.NewManager()

// Add S3 provider
s3Provider := s3.NewProvider()
config := map[string]interface{}{
    "region":            "us-east-1",
    "access_key_id":     "your-access-key",
    "secret_access_key": "your-secret-key",
    "bucket":            "my-bucket",
}

err := s3Provider.Configure(config)
if err != nil {
    log.Fatal(err)
}

storageManager.RegisterProvider(s3Provider)

// Upload file
ctx := context.Background()
content := strings.NewReader("Hello, World!")

putReq := &storage.PutObjectRequest{
    Bucket:      "my-bucket",
    Key:         "test.txt",
    Content:     content,
    Size:        13,
    ContentType: "text/plain",
}

response, err := storageManager.PutObject(ctx, "s3", putReq)
if err != nil {
    log.Fatal(err)
}
```

**Configuration**:
```yaml
# config.yaml
storage:
  providers:
    s3:
      region: "us-east-1"
      access_key_id: "${AWS_ACCESS_KEY_ID}"
      secret_access_key: "${AWS_SECRET_ACCESS_KEY}"
      bucket: "my-bucket"
    gcs:
      credentials_file: "${GCS_CREDENTIALS_FILE}"
      bucket: "my-bucket"
```

### 5. Messaging (`go-micro-libs/messaging`)

**Purpose**: Message queue integration

**Providers**:
- Kafka
- RabbitMQ
- NATS
- AWS SQS

**Usage**:
```go
import (
    "github.com/anasamu/go-micro-libs/messaging"
    "github.com/anasamu/go-micro-libs/messaging/providers/kafka"
)

// Initialize messaging manager
msgManager := messaging.NewManager()

// Add Kafka provider
kafkaProvider := kafka.NewProvider()
config := map[string]interface{}{
    "brokers": []string{"localhost:9092"},
}

err := kafkaProvider.Configure(config)
if err != nil {
    log.Fatal(err)
}

msgManager.RegisterProvider(kafkaProvider)

// Connect to messaging
ctx := context.Background()
err = msgManager.Connect(ctx, "kafka")
if err != nil {
    log.Fatal(err)
}

// Publish message
message := messaging.CreateMessage("user.created", map[string]interface{}{
    "user_id": "123",
    "email":   "user@example.com",
})

publishReq := &messaging.PublishRequest{
    Topic:   "users",
    Message: message,
}

response, err := msgManager.PublishMessage(ctx, "kafka", publishReq)
if err != nil {
    log.Fatal(err)
}
```

**Configuration**:
```yaml
# config.yaml
messaging:
  providers:
    kafka:
      brokers: ["localhost:9092"]
      group_id: "my-service"
      topics: ["users", "orders"]
    rabbitmq:
      url: "amqp://guest:guest@localhost:5672/"
      exchange: "my-exchange"
      queue: "my-queue"
```

## 🔧 Advanced Libraries

### 1. Authentication (`go-micro-libs/auth`)

**Purpose**: Authentication dan authorization

**Providers**:
- JWT
- OAuth2
- LDAP
- SAML
- 2FA

**Usage**:
```go
import (
    "github.com/anasamu/go-micro-libs/auth"
    "github.com/anasamu/go-micro-libs/auth/providers/jwt"
)

// Initialize auth manager
authManager := auth.NewManager()

// Add JWT provider
jwtProvider := jwt.NewProvider()
config := map[string]interface{}{
    "secret":     "your-secret",
    "expiration": "24h",
    "issuer":     "my-service",
}

err := jwtProvider.Configure(config)
if err != nil {
    log.Fatal(err)
}

authManager.RegisterProvider(jwtProvider)

// Generate token
token, err := authManager.GenerateToken("user123", map[string]interface{}{
    "role": "admin",
})
if err != nil {
    log.Fatal(err)
}

// Validate token
claims, err := authManager.ValidateToken(token)
if err != nil {
    log.Fatal(err)
}
```

**Configuration**:
```yaml
# config.yaml
auth:
  providers:
    jwt:
      secret: "${JWT_SECRET}"
      expiration: "24h"
      issuer: "my-service"
      audience: "api"
    oauth:
      client_id: "${OAUTH_CLIENT_ID}"
      client_secret: "${OAUTH_CLIENT_SECRET}"
      redirect_url: "${OAUTH_REDIRECT_URL}"
      scopes: ["read", "write"]
```

### 2. Circuit Breaker (`go-micro-libs/circuitbreaker`)

**Purpose**: Resilience patterns

**Providers**:
- Custom
- GoBreaker

**Usage**:
```go
import (
    "github.com/anasamu/go-micro-libs/circuitbreaker"
    "github.com/anasamu/go-micro-libs/circuitbreaker/providers/gobreaker"
)

// Initialize circuit breaker manager
cbManager := circuitbreaker.NewManager()

// Add GoBreaker provider
gobreakerProvider := gobreaker.NewProvider()
config := map[string]interface{}{
    "max_requests": 3,
    "interval":     "10s",
    "timeout":      "30s",
}

err := gobreakerProvider.Configure(config)
if err != nil {
    log.Fatal(err)
}

cbManager.RegisterProvider(gobreakerProvider)

// Execute with circuit breaker
result, err := cbManager.Execute("my-service", func() (interface{}, error) {
    // Your service call here
    return externalService.Call()
})
if err != nil {
    log.Fatal(err)
}
```

**Configuration**:
```yaml
# config.yaml
circuitbreaker:
  providers:
    gobreaker:
      max_requests: 3
      interval: "10s"
      timeout: "30s"
      ready_to_trip: "failure_ratio"
```

### 3. Rate Limiting (`go-micro-libs/ratelimit`)

**Purpose**: Rate limiting dengan berbagai algoritma

**Providers**:
- Token Bucket
- Sliding Window
- Leaky Bucket

**Usage**:
```go
import (
    "github.com/anasamu/go-micro-libs/ratelimit"
    "github.com/anasamu/go-micro-libs/ratelimit/providers/tokenbucket"
)

// Initialize rate limit manager
rlManager := ratelimit.NewManager()

// Add token bucket provider
tbProvider := tokenbucket.NewProvider()
config := map[string]interface{}{
    "rate":  100,
    "burst": 10,
}

err := tbProvider.Configure(config)
if err != nil {
    log.Fatal(err)
}

rlManager.RegisterProvider(tbProvider)

// Check rate limit
allowed, err := rlManager.Allow("user123")
if err != nil {
    log.Fatal(err)
}

if !allowed {
    // Rate limit exceeded
    return errors.New("rate limit exceeded")
}
```

**Configuration**:
```yaml
# config.yaml
ratelimit:
  providers:
    tokenbucket:
      rate: 100
      burst: 10
    slidingwindow:
      window: "1m"
      limit: 100
```

## 🔧 Library Combination Examples

### Example 1: E-commerce Service

```go
// Combine multiple libraries for e-commerce service
func setupEcommerceService() {
    // Core libraries (always available)
    configManager := config.NewManager()
    loggingManager := logging.NewManager()
    monitoringManager := monitoring.NewManager()
    middlewareManager := middleware.NewManager()
    commManager := communication.NewManager()
    
    // Optional libraries
    dbManager := database.NewManager()
    cacheManager := cache.NewManager()
    storageManager := storage.NewManager()
    msgManager := messaging.NewManager()
    authManager := auth.NewManager()
    paymentManager := payment.NewManager()
    
    // Initialize all managers
    // ... initialization code ...
}
```

### Example 2: AI-Powered Chat Service

```go
// Combine AI, database, cache, and messaging
func setupChatService() {
    // Core libraries
    configManager := config.NewManager()
    loggingManager := logging.NewManager()
    monitoringManager := monitoring.NewManager()
    
    // AI and data libraries
    aiManager := ai.NewManager()
    dbManager := database.NewManager()
    cacheManager := cache.NewManager()
    msgManager := messaging.NewManager()
    
    // Initialize all managers
    // ... initialization code ...
}
```

### Example 3: Event-Driven Service

```go
// Combine event sourcing, messaging, and database
func setupEventService() {
    // Core libraries
    configManager := config.NewManager()
    loggingManager := logging.NewManager()
    monitoringManager := monitoring.NewManager()
    
    // Event and data libraries
    eventManager := event.NewManager()
    msgManager := messaging.NewManager()
    dbManager := database.NewManager()
    cacheManager := cache.NewManager()
    
    // Initialize all managers
    // ... initialization code ...
}
```

## 🔧 Configuration Management

### Environment Variables

```bash
# Core libraries
SERVICE_NAME=my-service
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

# AI Services
OPENAI_API_KEY=your-openai-key

# Storage
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
S3_BUCKET=my-bucket

# Messaging
KAFKA_BROKERS=localhost:9092
```

### Configuration File

```yaml
# config.yaml
service:
  name: "my-service"
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
      service_name: "my-service"

middleware:
  auth:
    enabled: true
    provider: "jwt"
  rate_limit:
    enabled: true
    requests_per_minute: 100

communication:
  providers:
    http:
      port: 8080
      timeout: 30s
    grpc:
      port: 9090
      timeout: 30s

# Optional libraries (only if enabled)
database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100

cache:
  providers:
    redis:
      url: "${REDIS_URL}"

storage:
  providers:
    s3:
      region: "us-east-1"
      access_key_id: "${AWS_ACCESS_KEY_ID}"
      secret_access_key: "${AWS_SECRET_ACCESS_KEY}"
      bucket: "${S3_BUCKET}"

messaging:
  providers:
    kafka:
      brokers: ["${KAFKA_BROKERS}"]

auth:
  providers:
    jwt:
      secret: "${JWT_SECRET}"
      expiration: "24h"

ai:
  providers:
    openai:
      api_key: "${OPENAI_API_KEY}"
```

## 🔧 Best Practices

### 1. Library Initialization

```go
// Initialize libraries in correct order
func initializeLibraries() error {
    // 1. Config first
    configManager := config.NewManager()
    if err := configManager.Load(); err != nil {
        return err
    }
    
    // 2. Logging second
    loggingManager := logging.NewManager()
    if err := loggingManager.Initialize(); err != nil {
        return err
    }
    
    // 3. Monitoring third
    monitoringManager := monitoring.NewManager()
    if err := monitoringManager.Start(); err != nil {
        return err
    }
    
    // 4. Other libraries
    // ... initialize other libraries ...
    
    return nil
}
```

### 2. Error Handling

```go
// Proper error handling for library operations
func handleLibraryError(err error, libraryName string) {
    if err != nil {
        logger.Error("Library error", 
            "library", libraryName,
            "error", err.Error(),
        )
        // Handle error appropriately
    }
}
```

### 3. Resource Cleanup

```go
// Cleanup resources on shutdown
func cleanupLibraries() {
    // Stop monitoring
    monitoringManager.Stop()
    
    // Close database connections
    dbManager.Close()
    
    // Close cache connections
    cacheManager.Close()
    
    // Close messaging connections
    msgManager.Close()
}
```

## 🔧 Troubleshooting

### Common Issues

1. **Library Not Found**
   - Check if library is properly imported
   - Verify go.mod includes the library
   - Run `go mod tidy`

2. **Configuration Errors**
   - Check configuration file format
   - Verify environment variables
   - Check provider configuration

3. **Connection Issues**
   - Verify external service is running
   - Check network connectivity
   - Verify credentials and configuration

4. **Performance Issues**
   - Check connection pooling settings
   - Monitor resource usage
   - Optimize configuration

### Debug Mode

```go
// Enable debug mode for troubleshooting
func enableDebugMode() {
    // Set debug level for logging
    loggingManager.SetLevel("debug")
    
    // Enable debug mode for monitoring
    monitoringManager.EnableDebug()
    
    // Enable debug mode for other libraries
    // ... enable debug for other libraries ...
}
```

---

**Library Integration - Seamlessly integrate 20+ microservices libraries! 🚀**

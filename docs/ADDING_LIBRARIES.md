# Adding Libraries to Services

## 🎯 Overview

GoMicroFramework memungkinkan Anda untuk menambahkan library tambahan ke service yang sudah ada. Setiap service secara otomatis memiliki core libraries terintegrasi, dan Anda dapat menambahkan library optional sesuai kebutuhan.

## 🔧 Core Libraries (Always Available)

Setiap service yang dibuat secara otomatis memiliki library berikut:

| Library | Description | Usage |
|---------|-------------|-------|
| **Config** | Configuration management | `config.NewManager()` |
| **Logging** | Structured logging | `logging.NewManager()` |
| **Monitoring** | Metrics and tracing | `monitoring.NewManager()` |
| **Middleware** | HTTP middleware | `middleware.NewManager()` |
| **Communication** | Communication protocols | `communication.NewManager()` |
| **Utils** | Utility functions | `utils.NewManager()` |

## 🔧 Adding Optional Libraries

### 1. Using CLI Command

#### Add Single Library

```bash
# Add authentication
microframework add auth --provider=jwt --config=auth.yaml

# Add database
microframework add database --provider=postgres --config=db.yaml

# Add cache
microframework add cache --provider=redis --config=cache.yaml

# Add messaging
microframework add messaging --provider=kafka --config=messaging.yaml

# Add storage
microframework add storage --provider=s3 --config=storage.yaml

# Add AI services
microframework add ai --provider=openai --config=ai.yaml

# Add payment processing
microframework add payment --provider=stripe --config=payment.yaml
```

#### Add Multiple Libraries

```bash
# Add multiple libraries at once
microframework add auth --provider=jwt
microframework add database --provider=postgres
microframework add cache --provider=redis
microframework add messaging --provider=kafka
```

### 2. Manual Integration

#### Step 1: Update go.mod

```go
// go.mod
module user-service

go 1.21

require (
    // Core libraries (always included)
    github.com/anasamu/go-micro-libs v0.1.0
    
    // Optional libraries (add as needed)
    // Database
    gorm.io/gorm v1.25.4
    gorm.io/driver/postgres v1.5.2
    
    // Cache
    github.com/go-redis/redis/v8 v8.11.5
    
    // Messaging
    github.com/Shopify/sarama v1.38.1
    
    // Storage
    github.com/aws/aws-sdk-go v1.44.327
    
    // AI
    github.com/sashabaranov/go-openai v1.14.2
    
    // Payment
    github.com/stripe/stripe-go/v72 v72.122.0
)
```

#### Step 2: Update Configuration

```yaml
# config.yaml
# Add configuration for new libraries

# Database configuration
database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100
      max_idle_connections: 10

# Cache configuration
cache:
  providers:
    redis:
      url: "${REDIS_URL}"
      db: 0
      pool_size: 10

# Messaging configuration
messaging:
  providers:
    kafka:
      brokers: ["localhost:9092"]
      group_id: "user-service"
      topics: ["users"]

# Storage configuration
storage:
  providers:
    s3:
      region: "us-east-1"
      access_key_id: "${AWS_ACCESS_KEY_ID}"
      secret_access_key: "${AWS_SECRET_ACCESS_KEY}"
      bucket: "${S3_BUCKET}"

# AI configuration
ai:
  providers:
    openai:
      api_key: "${OPENAI_API_KEY}"
      default_model: "gpt-4"

# Payment configuration
payment:
  providers:
    stripe:
      secret_key: "${STRIPE_SECRET_KEY}"
      publishable_key: "${STRIPE_PUBLISHABLE_KEY}"
```

#### Step 3: Update Main.go

```go
// cmd/main.go
package main

import (
    "context"
    "log"
    
    // Core libraries (always available)
    "github.com/anasamu/go-micro-libs/config"
    "github.com/anasamu/go-micro-libs/logging"
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/middleware"
    "github.com/anasamu/go-micro-libs/communication"
    
    // Optional libraries (add as needed)
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/cache"
    "github.com/anasamu/go-micro-libs/messaging"
    "github.com/anasamu/go-micro-libs/storage"
    "github.com/anasamu/go-micro-libs/ai"
    "github.com/anasamu/go-micro-libs/payment"
)

func main() {
    ctx := context.Background()
    
    // Initialize core managers (always available)
    configManager := config.NewManager()
    loggingManager := logging.NewManager()
    monitoringManager := monitoring.NewManager()
    middlewareManager := middleware.NewManager()
    communicationManager := communication.NewManager()
    
    // Initialize optional managers (add as needed)
    dbManager := database.NewManager()
    cacheManager := cache.NewManager()
    msgManager := messaging.NewManager()
    storageManager := storage.NewManager()
    aiManager := ai.NewManager()
    paymentManager := payment.NewManager()
    
    // Bootstrap service
    if err := bootstrapService(ctx, configManager, loggingManager, monitoringManager, 
        middlewareManager, communicationManager, dbManager, cacheManager, 
        msgManager, storageManager, aiManager, paymentManager); err != nil {
        log.Fatal("Failed to bootstrap service:", err)
    }
    
    log.Println("Service started successfully")
}

func bootstrapService(ctx context.Context, 
    configManager *config.ConfigManager,
    loggingManager *logging.LoggingManager,
    monitoringManager *monitoring.MonitoringManager,
    middlewareManager *middleware.MiddlewareManager,
    communicationManager *communication.CommunicationManager,
    dbManager *database.DatabaseManager,
    cacheManager *cache.CacheManager,
    msgManager *messaging.MessagingManager,
    storageManager *storage.StorageManager,
    aiManager *ai.AIManager,
    paymentManager *payment.PaymentManager) error {
    
    // Load configuration
    if err := configManager.Load(); err != nil {
        return err
    }
    
    // Initialize logging
    if err := loggingManager.Initialize(); err != nil {
        return err
    }
    
    // Start monitoring
    if err := monitoringManager.Start(); err != nil {
        return err
    }
    
    // Setup middleware
    if err := middlewareManager.SetupChain(); err != nil {
        return err
    }
    
    // Start communication
    if err := communicationManager.Start(); err != nil {
        return err
    }
    
    // Initialize optional libraries (add as needed)
    
    // Database
    if err := dbManager.Connect(ctx); err != nil {
        return err
    }
    
    // Cache
    if err := cacheManager.Connect(ctx); err != nil {
        return err
    }
    
    // Messaging
    if err := msgManager.Connect(ctx); err != nil {
        return err
    }
    
    // Storage
    if err := storageManager.Initialize(); err != nil {
        return err
    }
    
    // AI
    if err := aiManager.Initialize(); err != nil {
        return err
    }
    
    // Payment
    if err := paymentManager.Initialize(); err != nil {
        return err
    }
    
    return nil
}
```

## 🔧 Library-Specific Integration

### 1. Database Integration

#### Add PostgreSQL

```bash
# Using CLI
microframework add database --provider=postgres --config=db.yaml
```

```go
// Manual integration
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
```

#### Add Redis

```bash
# Using CLI
microframework add cache --provider=redis --config=cache.yaml
```

```go
// Manual integration
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
```

### 2. Messaging Integration

#### Add Kafka

```bash
# Using CLI
microframework add messaging --provider=kafka --config=messaging.yaml
```

```go
// Manual integration
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
```

### 3. Storage Integration

#### Add S3

```bash
# Using CLI
microframework add storage --provider=s3 --config=storage.yaml
```

```go
// Manual integration
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

// Initialize storage
err = storageManager.Initialize()
if err != nil {
    log.Fatal(err)
}
```

### 4. AI Integration

#### Add OpenAI

```bash
# Using CLI
microframework add ai --provider=openai --config=ai.yaml
```

```go
// Manual integration
import (
    "github.com/anasamu/go-micro-libs/ai"
    "github.com/anasamu/go-micro-libs/ai/providers/openai"
)

// Initialize AI manager
aiManager := ai.NewManager()

// Add OpenAI provider
openaiProvider := openai.NewProvider()
config := map[string]interface{}{
    "api_key": "your-api-key",
    "base_url": "https://api.openai.com/v1",
}

err := openaiProvider.Configure(config)
if err != nil {
    log.Fatal(err)
}

aiManager.RegisterProvider(openaiProvider)

// Initialize AI
err = aiManager.Initialize()
if err != nil {
    log.Fatal(err)
}
```

### 5. Payment Integration

#### Add Stripe

```bash
# Using CLI
microframework add payment --provider=stripe --config=payment.yaml
```

```go
// Manual integration
import (
    "github.com/anasamu/go-micro-libs/payment"
    "github.com/anasamu/go-micro-libs/payment/providers/stripe"
)

// Initialize payment manager
paymentManager := payment.NewManager()

// Add Stripe provider
stripeProvider := stripe.NewProvider()
config := map[string]interface{}{
    "secret_key": "your-secret-key",
    "publishable_key": "your-publishable-key",
}

err := stripeProvider.Configure(config)
if err != nil {
    log.Fatal(err)
}

paymentManager.RegisterProvider(stripeProvider)

// Initialize payment
err = paymentManager.Initialize()
if err != nil {
    log.Fatal(err)
}
```

## 🔧 Advanced Library Integration

### 1. Custom Provider Integration

```go
// Create custom provider
type CustomProvider struct {
    config map[string]interface{}
}

func (p *CustomProvider) Configure(config map[string]interface{}) error {
    p.config = config
    return nil
}

func (p *CustomProvider) Connect(ctx context.Context) error {
    // Custom connection logic
    return nil
}

// Register custom provider
customProvider := &CustomProvider{}
err := manager.RegisterProvider(customProvider)
if err != nil {
    log.Fatal(err)
}
```

### 2. Multiple Provider Integration

```go
// Add multiple providers of the same type
dbManager := database.NewManager()

// Add PostgreSQL provider
postgresProvider := postgresql.NewProvider()
postgresConfig := map[string]interface{}{
    "host": "postgres-primary",
    "port": 5432,
    "user": "postgres",
    "password": "password",
    "database": "mydb",
}
postgresProvider.Configure(postgresConfig)
dbManager.RegisterProvider(postgresProvider)

// Add MySQL provider
mysqlProvider := mysql.NewProvider()
mysqlConfig := map[string]interface{}{
    "host": "mysql-secondary",
    "port": 3306,
    "user": "mysql",
    "password": "password",
    "database": "mydb",
}
mysqlProvider.Configure(mysqlConfig)
dbManager.RegisterProvider(mysqlProvider)

// Use specific provider
err := dbManager.Connect(ctx, "postgresql")
if err != nil {
    // Fallback to MySQL
    err = dbManager.Connect(ctx, "mysql")
    if err != nil {
        log.Fatal(err)
    }
}
```

### 3. Conditional Library Integration

```go
// Load configuration first
configManager := config.NewManager()
configManager.Load()

// Check if library should be enabled
if configManager.GetBool("database.enabled") {
    dbManager := database.NewManager()
    // Initialize database
}

if configManager.GetBool("cache.enabled") {
    cacheManager := cache.NewManager()
    // Initialize cache
}

if configManager.GetBool("messaging.enabled") {
    msgManager := messaging.NewManager()
    // Initialize messaging
}
```

## 🔧 Configuration Templates

### 1. Database Configuration Template

```yaml
# configs/database.yaml
database:
  providers:
    postgresql:
      enabled: true
      url: "${DATABASE_URL}"
      max_connections: 100
      max_idle_connections: 10
      connection_max_lifetime: "1h"
      
    mysql:
      enabled: false
      url: "${MYSQL_URL}"
      max_connections: 100
      max_idle_connections: 10
      
    mongodb:
      enabled: false
      url: "${MONGODB_URL}"
      database: "user-service"
      max_pool_size: 100
```

### 2. Cache Configuration Template

```yaml
# configs/cache.yaml
cache:
  providers:
    redis:
      enabled: true
      url: "${REDIS_URL}"
      db: 0
      pool_size: 10
      ttl: "1h"
      
    memory:
      enabled: false
      max_size: 1000
      ttl: "30m"
      
    memcached:
      enabled: false
      servers: ["localhost:11211"]
      ttl: "1h"
```

### 3. Messaging Configuration Template

```yaml
# configs/messaging.yaml
messaging:
  providers:
    kafka:
      enabled: true
      brokers: ["localhost:9092"]
      group_id: "user-service"
      topics: ["users", "orders"]
      
    rabbitmq:
      enabled: false
      url: "amqp://guest:guest@localhost:5672/"
      exchange: "user-exchange"
      queue: "user-queue"
      
    nats:
      enabled: false
      url: "nats://localhost:4222"
      cluster_id: "user-service"
```

## 🔧 Best Practices

### 1. Library Initialization Order

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
    
    // 4. Database
    if configManager.GetBool("database.enabled") {
        dbManager := database.NewManager()
        if err := dbManager.Connect(ctx); err != nil {
            return err
        }
    }
    
    // 5. Cache
    if configManager.GetBool("cache.enabled") {
        cacheManager := cache.NewManager()
        if err := cacheManager.Connect(ctx); err != nil {
            return err
        }
    }
    
    // 6. Other libraries...
    
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
    if dbManager != nil {
        dbManager.Close()
    }
    
    // Close cache connections
    if cacheManager != nil {
        cacheManager.Close()
    }
    
    // Close messaging connections
    if msgManager != nil {
        msgManager.Close()
    }
}
```

### 4. Health Checks

```go
// Register health checks for libraries
func registerHealthChecks() {
    // Database health check
    if dbManager != nil {
        monitoringManager.RegisterHealthCheck("database", func() error {
            return dbManager.HealthCheck(ctx)
        })
    }
    
    // Cache health check
    if cacheManager != nil {
        monitoringManager.RegisterHealthCheck("cache", func() error {
            return cacheManager.HealthCheck(ctx)
        })
    }
    
    // Messaging health check
    if msgManager != nil {
        monitoringManager.RegisterHealthCheck("messaging", func() error {
            return msgManager.HealthCheck(ctx)
        })
    }
}
```

## 🔧 Troubleshooting

### Common Issues

1. **Library Not Found**
   ```bash
   # Check if library is properly imported
   go mod tidy
   go mod download
   ```

2. **Configuration Errors**
   ```bash
   # Validate configuration
   microframework config validate
   ```

3. **Connection Issues**
   ```bash
   # Check external service is running
   # Verify network connectivity
   # Check credentials and configuration
   ```

4. **Performance Issues**
   ```bash
   # Check connection pooling settings
   # Monitor resource usage
   # Optimize configuration
   ```

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

**Adding Libraries - Seamlessly integrate additional libraries to your services! 🚀**

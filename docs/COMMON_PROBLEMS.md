# Common Microservices Problems and Solutions

## 🎯 Overview

GoMicroFramework menyediakan solusi untuk berbagai masalah umum yang terjadi dalam pengembangan microservices. Dokumen ini menjelaskan masalah-masalah yang sering terjadi dan bagaimana framework membantu mengatasi masalah tersebut.

## 🔧 Common Problems and Solutions

### 1. Service Communication Issues

#### Problem: Network Timeouts and Failures
**Symptoms:**
- Services fail to communicate with each other
- Timeout errors in service calls
- Network connectivity issues

**Solutions with GoMicroFramework:**
```go
// Circuit breaker for handling failures
circuitBreaker := circuitbreaker.NewCircuitBreaker(
    5,              // failure threshold
    30*time.Second, // timeout
    3,              // max requests in half-open state
)

// Retry mechanism with exponential backoff
retryConfig := &retry.Config{
    MaxRetries: 3,
    Backoff:    retry.ExponentialBackoff,
    Interval:   1 * time.Second,
}

// Timeout configuration
httpClient := &http.Client{
    Timeout: 30 * time.Second,
}
```

#### Problem: Service Discovery Failures
**Symptoms:**
- Services cannot find each other
- DNS resolution issues
- Service registry unavailability

**Solutions with GoMicroFramework:**
```go
// Multiple service discovery providers
discoveryManager := discovery.NewManager(config, logger, metrics)

// Register multiple providers
discoveryManager.RegisterProvider("consul", consulProvider)
discoveryManager.RegisterProvider("etcd", etcdProvider)
discoveryManager.RegisterProvider("kubernetes", k8sProvider)

// Fallback mechanism
services, err := discoveryManager.DiscoverServices(ctx, "consul", "user-service")
if err != nil {
    // Fallback to etcd
    services, err = discoveryManager.DiscoverServices(ctx, "etcd", "user-service")
}
```

### 2. Data Consistency Issues

#### Problem: Distributed Transactions
**Symptoms:**
- Data inconsistency across services
- Partial transaction failures
- Rollback complexity

**Solutions with GoMicroFramework:**
```go
// Saga pattern implementation
sagaManager := saga.NewManager(config, logger, metrics)

// Define saga steps
sagaSteps := []saga.Step{
    {
        Name: "create-user",
        Action: func(ctx context.Context) error {
            return userService.CreateUser(ctx, user)
        },
        Compensation: func(ctx context.Context) error {
            return userService.DeleteUser(ctx, user.ID)
        },
    },
    {
        Name: "create-profile",
        Action: func(ctx context.Context) error {
            return profileService.CreateProfile(ctx, profile)
        },
        Compensation: func(ctx context.Context) error {
            return profileService.DeleteProfile(ctx, profile.ID)
        },
    },
}

// Execute saga
err := sagaManager.ExecuteSaga(ctx, sagaSteps)
```

#### Problem: Eventual Consistency
**Symptoms:**
- Data synchronization delays
- Inconsistent reads
- Event ordering issues

**Solutions with GoMicroFramework:**
```go
// Event sourcing with event store
eventStore := event.NewEventStore(dbManager, msgManager)

// Append events
err := eventStore.AppendEvent(ctx, "user-stream", "user-created", userCreatedEvent)
if err != nil {
    return err
}

// Publish event to message broker
message := messaging.CreateMessage("user.created", userCreatedEvent)
publishReq := &messaging.PublishRequest{
    Topic:   "user-events",
    Message: message,
}

_, err = msgManager.PublishMessage(ctx, "kafka", publishReq)
```

### 3. Performance Issues

#### Problem: Slow Database Queries
**Symptoms:**
- High database response times
- Database connection pool exhaustion
- Slow query performance

**Solutions with GoMicroFramework:**
```go
// Database connection pooling
dbConfig := &database.Config{
    MaxConnections:     100,
    MaxIdleConnections: 10,
    ConnectionMaxLifetime: 1 * time.Hour,
}

// Query caching
cacheManager := cache.NewManager(config, logger, metrics)

// Cache query results
cacheKey := fmt.Sprintf("user:%s", userID)
var user User
err := cacheManager.Get(ctx, cacheKey, &user)
if err != nil {
    // Cache miss - query database
    user, err = dbManager.Query(ctx, "postgresql", query, userID)
    if err == nil {
        // Cache the result
        cacheManager.Set(ctx, cacheKey, user, 1*time.Hour)
    }
}
```

#### Problem: Memory Leaks
**Symptoms:**
- Increasing memory usage
- Service crashes due to OOM
- Performance degradation

**Solutions with GoMicroFramework:**
```go
// Resource cleanup
func (s *Service) cleanup() {
    // Close database connections
    if s.dbManager != nil {
        s.dbManager.Close()
    }
    
    // Close cache connections
    if s.cacheManager != nil {
        s.cacheManager.Close()
    }
    
    // Close message broker connections
    if s.msgManager != nil {
        s.msgManager.Close()
    }
}

// Memory monitoring
func (s *Service) monitorMemory() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        
        if m.Alloc > 100*1024*1024 { // 100MB threshold
            s.logger.Warn("High memory usage detected",
                "alloc_mb", m.Alloc/1024/1024,
                "sys_mb", m.Sys/1024/1024,
            )
        }
    }
}
```

### 4. Security Issues

#### Problem: Authentication Failures
**Symptoms:**
- Unauthorized access
- Token validation errors
- Session management issues

**Solutions with GoMicroFramework:**
```go
// JWT authentication with validation
jwtService := auth.NewJWTService(secret, expiration)

// Validate token
claims, err := jwtService.ValidateToken(ctx, tokenString)
if err != nil {
    return nil, fmt.Errorf("invalid token: %w", err)
}

// Check token blacklist
if jwtService.IsTokenBlacklisted(ctx, tokenString) {
    return nil, fmt.Errorf("token is blacklisted")
}

// Rate limiting for authentication endpoints
rateLimiter := ratelimit.NewRateLimiter(100, 10) // 100 requests per minute, burst 10

if !rateLimiter.Allow(clientIP) {
    return nil, fmt.Errorf("rate limit exceeded")
}
```

#### Problem: Authorization Issues
**Symptoms:**
- Insufficient permissions
- Role-based access control failures
- Policy enforcement issues

**Solutions with GoMicroFramework:**
```go
// RBAC authorization
rbacService := auth.NewRBACService(dbManager, cacheManager)

// Check user permissions
hasPermission, err := rbacService.HasPermission(ctx, userID, "user:read")
if err != nil {
    return nil, err
}

if !hasPermission {
    return nil, fmt.Errorf("insufficient permissions")
}

// ABAC authorization
abacService := auth.NewABACService(dbManager, cacheManager)

// Evaluate policy
context := &auth.Context{
    User: map[string]interface{}{
        "id":   userID,
        "role": "user",
    },
    Resource: map[string]interface{}{
        "name": "user-profile",
        "owner": userID,
    },
    Environment: map[string]interface{}{
        "time": time.Now().Format("15:04"),
    },
}

allowed, err := abacService.EvaluateAccess(ctx, userID, "user-profile", "read", context)
if err != nil {
    return nil, err
}

if !allowed {
    return nil, fmt.Errorf("access denied")
}
```

### 5. Monitoring and Observability Issues

#### Problem: Lack of Visibility
**Symptoms:**
- Difficult to debug issues
- No performance metrics
- Poor error tracking

**Solutions with GoMicroFramework:**
```go
// Comprehensive monitoring
monitoringManager := monitoring.NewManager(config, logger)

// Metrics collection
monitoringManager.IncrementCounter("requests_total", map[string]string{
    "method": "GET",
    "path":   "/users",
    "status": "200",
})

monitoringManager.RecordHistogram("request_duration_seconds", 
    duration.Seconds(), map[string]string{
        "method": "GET",
        "path":   "/users",
    })

// Distributed tracing
span := monitoringManager.StartSpan("user-service", "get-user")
defer span.Finish()

span.SetTag("user.id", userID)
span.SetTag("operation", "get-user")

// Health checks
monitoringManager.RegisterHealthCheck("database", func() error {
    return dbManager.HealthCheck(ctx)
})

monitoringManager.RegisterHealthCheck("cache", func() error {
    return cacheManager.HealthCheck(ctx)
})
```

#### Problem: Log Management
**Symptoms:**
- Inconsistent log formats
- Difficult log aggregation
- Poor log correlation

**Solutions with GoMicroFramework:**
```go
// Structured logging
loggingManager := logging.NewManager(config)

// Correlation ID tracking
correlationID := generateCorrelationID()
ctx = context.WithValue(ctx, "correlation_id", correlationID)

// Structured logging with context
loggingManager.Info("User created",
    "user_id", user.ID,
    "email", user.Email,
    "correlation_id", correlationID,
    "duration", time.Since(start),
)

// Log aggregation
loggingManager.ConfigureElasticsearch(config.Elasticsearch)
loggingManager.ConfigureFluentd(config.Fluentd)
```

### 6. Deployment and Configuration Issues

#### Problem: Configuration Management
**Symptoms:**
- Configuration drift
- Environment-specific issues
- Secret management problems

**Solutions with GoMicroFramework:**
```go
// Multi-source configuration
configManager := config.NewManager()

// Load from multiple sources
configManager.LoadFromFile("config.yaml")
configManager.LoadFromEnv("SERVICE_")
configManager.LoadFromConsul("service/config")

// Configuration validation
schema := &config.Schema{
    Required: []string{"database.url", "redis.url"},
    Types: map[string]string{
        "database.url": "string",
        "redis.url":    "string",
        "port":         "int",
    },
}

err := configManager.Validate(schema)
if err != nil {
    return fmt.Errorf("configuration validation failed: %w", err)
}

// Hot reloading
configManager.Watch(func(newConfig *config.Config) {
    logger.Info("Configuration reloaded")
    // Update service configuration
})
```

#### Problem: Service Deployment
**Symptoms:**
- Deployment failures
- Service startup issues
- Health check failures

**Solutions with GoMicroFramework:**
```go
// Graceful shutdown
func (s *Service) gracefulShutdown() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    <-sigChan
    logger.Info("Shutting down service...")
    
    // Stop accepting new requests
    s.server.Shutdown(ctx)
    
    // Wait for existing requests to complete
    s.wg.Wait()
    
    // Cleanup resources
    s.cleanup()
    
    logger.Info("Service shutdown complete")
}

// Health check endpoint
func (s *Service) healthCheck(c *gin.Context) {
    health := map[string]string{
        "status": "healthy",
        "timestamp": time.Now().Format(time.RFC3339),
    }
    
    // Check dependencies
    if err := s.dbManager.HealthCheck(ctx); err != nil {
        health["database"] = "unhealthy"
        health["status"] = "unhealthy"
    } else {
        health["database"] = "healthy"
    }
    
    if err := s.cacheManager.HealthCheck(ctx); err != nil {
        health["cache"] = "unhealthy"
        health["status"] = "unhealthy"
    } else {
        health["cache"] = "healthy"
    }
    
    if health["status"] == "healthy" {
        c.JSON(http.StatusOK, health)
    } else {
        c.JSON(http.StatusServiceUnavailable, health)
    }
}
```

### 7. Error Handling and Recovery

#### Problem: Error Propagation
**Symptoms:**
- Errors not properly handled
- Poor error messages
- Difficult error debugging

**Solutions with GoMicroFramework:**
```go
// Structured error handling
type ServiceError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
    Cause   error  `json:"-"`
}

func (e *ServiceError) Error() string {
    return e.Message
}

// Error middleware
func errorMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            // Log error with context
            logger.Error("Request error",
                "error", err.Error(),
                "path", c.Request.URL.Path,
                "method", c.Request.Method,
                "correlation_id", c.GetString("correlation_id"),
            )
            
            // Return structured error response
            serviceErr := &ServiceError{
                Code:    "INTERNAL_ERROR",
                Message: "An internal error occurred",
                Details: err.Error(),
            }
            
            c.JSON(http.StatusInternalServerError, serviceErr)
        }
    }
}

// Retry mechanism
func retryOperation(operation func() error, maxRetries int) error {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        err := operation()
        if err == nil {
            return nil
        }
        
        lastErr = err
        
        // Exponential backoff
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    
    return fmt.Errorf("operation failed after %d retries: %w", maxRetries, lastErr)
}
```

### 8. Testing Issues

#### Problem: Integration Testing
**Symptoms:**
- Difficult to test service interactions
- Flaky tests
- Poor test coverage

**Solutions with GoMicroFramework:**
```go
// Integration test framework
type IntegrationTest struct {
    services map[string]*Service
    config   *config.Config
    logger   *logging.Logger
}

func (it *IntegrationTest) SetupTest() error {
    // Start test services
    for name, service := range it.services {
        err := service.Start()
        if err != nil {
            return fmt.Errorf("failed to start service %s: %w", name, err)
        }
    }
    
    // Wait for services to be healthy
    return it.waitForHealthy()
}

func (it *IntegrationTest) TearDownTest() error {
    // Stop all services
    for name, service := range it.services {
        err := service.Stop()
        if err != nil {
            logger.Warn("Failed to stop service", "service", name, "error", err)
        }
    }
    
    return nil
}

// Mock services
func (it *IntegrationTest) MockService(name string, mock *MockService) {
    it.services[name] = mock
}

// Test service communication
func TestServiceCommunication(t *testing.T) {
    test := &IntegrationTest{
        services: make(map[string]*Service),
    }
    
    // Setup test
    err := test.SetupTest()
    assert.NoError(t, err)
    defer test.TearDownTest()
    
    // Test service communication
    client := NewServiceClient("user-service")
    user, err := client.GetUser("123")
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

## 🔧 Best Practices

### 1. Error Handling
```go
// Always handle errors properly
result, err := service.DoSomething()
if err != nil {
    // Log error with context
    logger.Error("Operation failed", "error", err, "context", context)
    
    // Return appropriate error
    return nil, fmt.Errorf("operation failed: %w", err)
}

// Use structured errors
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field %s: %s", e.Field, e.Message)
}
```

### 2. Resource Management
```go
// Always cleanup resources
func (s *Service) Start() error {
    // Start services
    if err := s.startServices(); err != nil {
        s.cleanup()
        return err
    }
    
    // Setup graceful shutdown
    go s.gracefulShutdown()
    
    return nil
}

func (s *Service) cleanup() {
    // Close all connections
    if s.dbManager != nil {
        s.dbManager.Close()
    }
    
    if s.cacheManager != nil {
        s.cacheManager.Close()
    }
    
    if s.msgManager != nil {
        s.msgManager.Close()
    }
}
```

### 3. Monitoring and Observability
```go
// Always include monitoring
func (s *Service) handleRequest(c *gin.Context) {
    start := time.Now()
    
    // Record metrics
    defer func() {
        duration := time.Since(start)
        metrics.RecordHistogram("request_duration_seconds", duration.Seconds())
        metrics.IncrementCounter("requests_total")
    }()
    
    // Process request
    c.Next()
}

// Health checks
func (s *Service) setupHealthChecks() {
    healthManager.RegisterHealthCheck("database", s.dbManager.HealthCheck)
    healthManager.RegisterHealthCheck("cache", s.cacheManager.HealthCheck)
    healthManager.RegisterHealthCheck("messaging", s.msgManager.HealthCheck)
}
```

### 4. Configuration Management
```go
// Use environment-specific configuration
func loadConfig() *config.Config {
    config := &config.Config{}
    
    // Load base configuration
    config.LoadFromFile("config.yaml")
    
    // Override with environment-specific config
    env := os.Getenv("ENVIRONMENT")
    if env != "" {
        config.LoadFromFile(fmt.Sprintf("config.%s.yaml", env))
    }
    
    // Override with environment variables
    config.LoadFromEnv("SERVICE_")
    
    return config
}
```

### 5. Security
```go
// Always validate input
func (s *Service) validateInput(input interface{}) error {
    // Use validation library
    validator := validator.New()
    
    if err := validator.Struct(input); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    return nil
}

// Use secure defaults
func (s *Service) setupSecurity() {
    // Enable HTTPS
    if s.config.TLS.Enabled {
        s.server.TLSConfig = &tls.Config{
            MinVersion: tls.VersionTLS12,
        }
    }
    
    // Set security headers
    s.server.Use(helmet.New())
    
    // Enable CORS
    s.server.Use(cors.New(cors.Config{
        AllowOrigins: s.config.CORS.AllowedOrigins,
        AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders: []string{"Content-Type", "Authorization"},
    }))
}
```

## 🔧 Troubleshooting Guide

### 1. Service Startup Issues
```bash
# Check service logs
docker logs <service-name>

# Check service health
curl http://localhost:8080/health

# Check service configuration
microframework config validate

# Check service dependencies
microframework health check
```

### 2. Communication Issues
```bash
# Check service discovery
microframework discovery list

# Check network connectivity
ping <service-address>

# Check service endpoints
curl http://<service-address>:<port>/health

# Check circuit breaker status
microframework circuit-breaker status
```

### 3. Performance Issues
```bash
# Check service metrics
curl http://localhost:9090/metrics

# Check service logs for errors
grep "ERROR" /var/log/service.log

# Check resource usage
docker stats <service-name>

# Check database performance
microframework database stats
```

### 4. Security Issues
```bash
# Check authentication
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/users

# Check authorization
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/admin

# Check rate limiting
curl -v http://localhost:8080/api/users

# Check security headers
curl -I http://localhost:8080/api/users
```

---

**Common Problems - Comprehensive solutions for microservices challenges! 🚀**

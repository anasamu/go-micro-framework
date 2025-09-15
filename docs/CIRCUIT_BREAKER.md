# Circuit Breaker Implementation

## 🎯 Overview

GoMicroFramework menyediakan implementasi circuit breaker yang komprehensif untuk melindungi services dari cascading failures dan meningkatkan resilience. Circuit breaker pattern membantu mencegah service dari terus-menerus mencoba operasi yang gagal, memberikan waktu untuk service yang bermasalah untuk pulih.

## 🔧 Circuit Breaker States

### 1. Closed State
- Normal operation
- Requests are allowed through
- Failure count is tracked
- When failure threshold is reached, transitions to Open state

### 2. Open State
- Requests are immediately rejected
- No calls to the failing service
- After timeout period, transitions to Half-Open state

### 3. Half-Open State
- Limited number of test requests are allowed
- If requests succeed, transitions to Closed state
- If requests fail, transitions back to Open state

## 🔧 Circuit Breaker Configuration

### 1. Generate Service with Circuit Breaker

```bash
# Generate service with circuit breaker
microframework new user-service --with-circuit-breaker --with-database=postgres --with-monitoring=prometheus
```

### 2. Configuration

```yaml
# config.yaml
service:
  name: "user-service"
  version: "1.0.0"
  port: 8080

# Core libraries
config:
  providers:
    file:
      path: "./configs"
      format: "yaml"

logging:
  providers:
    console:
      level: "info"
      format: "json"

monitoring:
  providers:
    prometheus:
      endpoint: ":9090"
    jaeger:
      endpoint: "http://localhost:14268"
      service_name: "user-service"

# Circuit breaker configuration
circuit_breaker:
  providers:
    memory:
      enabled: true
      failure_threshold: 5
      timeout: 30s
      max_requests: 3
      success_threshold: 2
      failure_ratio: 0.5
      min_requests: 10
      
    redis:
      enabled: false
      url: "${REDIS_URL}"
      db: 1
      key_prefix: "circuit_breaker:"
      ttl: "1h"
      
    consul:
      enabled: false
      address: "localhost:8500"
      token: "${CONSUL_TOKEN}"
      service_name: "user-service"
      check_interval: "10s"

# Database for circuit breaker data
database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100

# Cache for circuit breaker state
cache:
  providers:
    redis:
      url: "${REDIS_URL}"
      db: 0
      pool_size: 10
```

## 🔧 Circuit Breaker Implementation

### 1. Circuit Breaker Manager

```go
// internal/circuitbreaker/manager.go
package circuitbreaker

import (
    "context"
    "fmt"
    "sync"
    "time"
    
    "github.com/anasamu/go-micro-libs/circuitbreaker"
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/logging"
)

type Manager struct {
    breakers    map[string]*CircuitBreaker
    mutex       sync.RWMutex
    config      *Config
    logger      *logging.Logger
    metrics     *monitoring.Metrics
}

type Config struct {
    Provider        string        `yaml:"provider"`
    FailureThreshold int          `yaml:"failure_threshold"`
    Timeout         time.Duration `yaml:"timeout"`
    MaxRequests     int           `yaml:"max_requests"`
    SuccessThreshold int          `yaml:"success_threshold"`
    FailureRatio    float64       `yaml:"failure_ratio"`
    MinRequests     int           `yaml:"min_requests"`
}

type CircuitBreaker struct {
    name            string
    state           State
    failureCount    int
    successCount    int
    requestCount    int
    lastFailureTime time.Time
    timeout         time.Duration
    failureThreshold int
    successThreshold int
    maxRequests     int
    failureRatio    float64
    minRequests     int
    mutex           sync.RWMutex
    logger          *logging.Logger
    metrics         *monitoring.Metrics
}

type State int

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

func (s State) String() string {
    switch s {
    case StateClosed:
        return "closed"
    case StateOpen:
        return "open"
    case StateHalfOpen:
        return "half-open"
    default:
        return "unknown"
    }
}

func NewManager(config *Config, logger *logging.Logger, metrics *monitoring.Metrics) *Manager {
    return &Manager{
        breakers: make(map[string]*CircuitBreaker),
        config:   config,
        logger:   logger,
        metrics:  metrics,
    }
}

func (m *Manager) GetBreaker(name string) *CircuitBreaker {
    m.mutex.RLock()
    breaker, exists := m.breakers[name]
    m.mutex.RUnlock()
    
    if !exists {
        m.mutex.Lock()
        defer m.mutex.Unlock()
        
        // Double-check after acquiring write lock
        breaker, exists = m.breakers[name]
        if !exists {
            breaker = NewCircuitBreaker(name, m.config, m.logger, m.metrics)
            m.breakers[name] = breaker
        }
    }
    
    return breaker
}

func (m *Manager) Execute(ctx context.Context, name string, operation func() (interface{}, error)) (interface{}, error) {
    breaker := m.GetBreaker(name)
    return breaker.Execute(ctx, operation)
}

func (m *Manager) GetState(name string) State {
    breaker := m.GetBreaker(name)
    return breaker.GetState()
}

func (m *Manager) Reset(name string) {
    breaker := m.GetBreaker(name)
    breaker.Reset()
}

func (m *Manager) GetStats(name string) *Stats {
    breaker := m.GetBreaker(name)
    return breaker.GetStats()
}
```

### 2. Circuit Breaker Implementation

```go
// internal/circuitbreaker/breaker.go
package circuitbreaker

import (
    "context"
    "fmt"
    "sync"
    "time"
    
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/logging"
)

func NewCircuitBreaker(name string, config *Config, logger *logging.Logger, metrics *monitoring.Metrics) *CircuitBreaker {
    return &CircuitBreaker{
        name:            name,
        state:           StateClosed,
        failureCount:    0,
        successCount:    0,
        requestCount:    0,
        timeout:         config.Timeout,
        failureThreshold: config.FailureThreshold,
        successThreshold: config.SuccessThreshold,
        maxRequests:     config.MaxRequests,
        failureRatio:    config.FailureRatio,
        minRequests:     config.MinRequests,
        logger:          logger,
        metrics:         metrics,
    }
}

func (cb *CircuitBreaker) Execute(ctx context.Context, operation func() (interface{}, error)) (interface{}, error) {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    // Check if circuit breaker should allow the request
    if !cb.shouldAllowRequest() {
        cb.recordRejectedRequest()
        return nil, fmt.Errorf("circuit breaker %s is %s", cb.name, cb.state.String())
    }
    
    // Execute the operation
    result, err := operation()
    
    // Record the result
    if err != nil {
        cb.recordFailure()
    } else {
        cb.recordSuccess()
    }
    
    // Update state based on results
    cb.updateState()
    
    return result, err
}

func (cb *CircuitBreaker) shouldAllowRequest() bool {
    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        // Check if timeout has passed
        if time.Since(cb.lastFailureTime) >= cb.timeout {
            cb.state = StateHalfOpen
            cb.requestCount = 0
            cb.logger.Info("Circuit breaker transitioning to half-open", "name", cb.name)
            return true
        }
        return false
    case StateHalfOpen:
        // Allow limited number of requests
        return cb.requestCount < cb.maxRequests
    default:
        return false
    }
}

func (cb *CircuitBreaker) recordFailure() {
    cb.failureCount++
    cb.requestCount++
    cb.lastFailureTime = time.Now()
    
    cb.logger.Debug("Circuit breaker recorded failure", 
        "name", cb.name, 
        "failure_count", cb.failureCount,
        "request_count", cb.requestCount)
    
    // Record metrics
    cb.metrics.IncrementCounter("circuit_breaker_failures_total", map[string]string{
        "name": cb.name,
    })
}

func (cb *CircuitBreaker) recordSuccess() {
    cb.successCount++
    cb.requestCount++
    
    cb.logger.Debug("Circuit breaker recorded success", 
        "name", cb.name, 
        "success_count", cb.successCount,
        "request_count", cb.requestCount)
    
    // Record metrics
    cb.metrics.IncrementCounter("circuit_breaker_successes_total", map[string]string{
        "name": cb.name,
    })
}

func (cb *CircuitBreaker) recordRejectedRequest() {
    cb.logger.Debug("Circuit breaker rejected request", "name", cb.name)
    
    // Record metrics
    cb.metrics.IncrementCounter("circuit_breaker_rejected_requests_total", map[string]string{
        "name": cb.name,
    })
}

func (cb *CircuitBreaker) updateState() {
    switch cb.state {
    case StateClosed:
        // Check if we should open the circuit
        if cb.shouldOpen() {
            cb.state = StateOpen
            cb.logger.Warn("Circuit breaker opened", "name", cb.name)
            
            // Record metrics
            cb.metrics.IncrementCounter("circuit_breaker_opened_total", map[string]string{
                "name": cb.name,
            })
        }
    case StateHalfOpen:
        // Check if we should close the circuit
        if cb.shouldClose() {
            cb.state = StateClosed
            cb.failureCount = 0
            cb.successCount = 0
            cb.requestCount = 0
            cb.logger.Info("Circuit breaker closed", "name", cb.name)
            
            // Record metrics
            cb.metrics.IncrementCounter("circuit_breaker_closed_total", map[string]string{
                "name": cb.name,
            })
        } else if cb.shouldReopen() {
            cb.state = StateOpen
            cb.logger.Warn("Circuit breaker reopened", "name", cb.name)
            
            // Record metrics
            cb.metrics.IncrementCounter("circuit_breaker_reopened_total", map[string]string{
                "name": cb.name,
            })
        }
    }
}

func (cb *CircuitBreaker) shouldOpen() bool {
    // Check if we have enough requests to make a decision
    if cb.requestCount < cb.minRequests {
        return false
    }
    
    // Check failure threshold
    if cb.failureCount >= cb.failureThreshold {
        return true
    }
    
    // Check failure ratio
    if cb.requestCount > 0 {
        failureRatio := float64(cb.failureCount) / float64(cb.requestCount)
        if failureRatio >= cb.failureRatio {
            return true
        }
    }
    
    return false
}

func (cb *CircuitBreaker) shouldClose() bool {
    // Check if we have enough successful requests
    return cb.successCount >= cb.successThreshold
}

func (cb *CircuitBreaker) shouldReopen() bool {
    // Check if we have too many failures in half-open state
    return cb.failureCount > 0
}

func (cb *CircuitBreaker) GetState() State {
    cb.mutex.RLock()
    defer cb.mutex.RUnlock()
    return cb.state
}

func (cb *CircuitBreaker) Reset() {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    cb.state = StateClosed
    cb.failureCount = 0
    cb.successCount = 0
    cb.requestCount = 0
    cb.lastFailureTime = time.Time{}
    
    cb.logger.Info("Circuit breaker reset", "name", cb.name)
    
    // Record metrics
    cb.metrics.IncrementCounter("circuit_breaker_reset_total", map[string]string{
        "name": cb.name,
    })
}

func (cb *CircuitBreaker) GetStats() *Stats {
    cb.mutex.RLock()
    defer cb.mutex.RUnlock()
    
    return &Stats{
        Name:            cb.name,
        State:           cb.state.String(),
        FailureCount:    cb.failureCount,
        SuccessCount:    cb.successCount,
        RequestCount:    cb.requestCount,
        LastFailureTime: cb.lastFailureTime,
        FailureRatio:    cb.getFailureRatio(),
    }
}

func (cb *CircuitBreaker) getFailureRatio() float64 {
    if cb.requestCount == 0 {
        return 0
    }
    return float64(cb.failureCount) / float64(cb.requestCount)
}

type Stats struct {
    Name            string    `json:"name"`
    State           string    `json:"state"`
    FailureCount    int       `json:"failure_count"`
    SuccessCount    int       `json:"success_count"`
    RequestCount    int       `json:"request_count"`
    LastFailureTime time.Time `json:"last_failure_time"`
    FailureRatio    float64   `json:"failure_ratio"`
}
```

### 3. Circuit Breaker Middleware

```go
// internal/middleware/circuit_breaker.go
package middleware

import (
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/anasamu/go-micro-libs/circuitbreaker"
    "github.com/anasamu/go-micro-libs/monitoring"
)

type CircuitBreakerMiddleware struct {
    manager *circuitbreaker.Manager
    config  *Config
}

type Config struct {
    Enabled          bool          `yaml:"enabled"`
    FailureThreshold int           `yaml:"failure_threshold"`
    Timeout          time.Duration `yaml:"timeout"`
    MaxRequests      int           `yaml:"max_requests"`
    SuccessThreshold int           `yaml:"success_threshold"`
    FailureRatio     float64       `yaml:"failure_ratio"`
    MinRequests      int           `yaml:"min_requests"`
}

func NewCircuitBreakerMiddleware(manager *circuitbreaker.Manager, config *Config) *CircuitBreakerMiddleware {
    return &CircuitBreakerMiddleware{
        manager: manager,
        config:  config,
    }
}

func (m *CircuitBreakerMiddleware) CircuitBreaker(name string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Execute request through circuit breaker
        result, err := m.manager.Execute(c.Request.Context(), name, func() (interface{}, error) {
            // Process request
            c.Next()
            
            // Check if request was successful
            if c.Writer.Status() >= 500 {
                return nil, fmt.Errorf("server error: %d", c.Writer.Status())
            }
            
            return nil, nil
        })
        
        if err != nil {
            // Circuit breaker rejected the request
            m.handleCircuitBreakerError(c, err)
            return
        }
        
        // Request was processed successfully
        c.Next()
    }
}

func (m *CircuitBreakerMiddleware) handleCircuitBreakerError(c *gin.Context, err error) {
    // Record metrics
    monitoring.IncrementCounter("circuit_breaker_rejected_requests_total", map[string]string{
        "endpoint": c.Request.URL.Path,
    })
    
    // Set appropriate status code
    c.JSON(http.StatusServiceUnavailable, gin.H{
        "error": "Service temporarily unavailable",
        "message": "Circuit breaker is open",
    })
    c.Abort()
}
```

### 4. Circuit Breaker for External Services

```go
// internal/clients/circuit_breaker_client.go
package clients

import (
    "context"
    "fmt"
    "time"
    
    "github.com/anasamu/go-micro-libs/circuitbreaker"
    "github.com/anasamu/go-micro-libs/monitoring"
)

type CircuitBreakerClient struct {
    httpClient *http.Client
    manager    *circuitbreaker.Manager
    baseURL    string
    name       string
}

func NewCircuitBreakerClient(httpClient *http.Client, manager *circuitbreaker.Manager, baseURL, name string) *CircuitBreakerClient {
    return &CircuitBreakerClient{
        httpClient: httpClient,
        manager:    manager,
        baseURL:    baseURL,
        name:       name,
    }
}

func (c *CircuitBreakerClient) Get(ctx context.Context, path string) (*http.Response, error) {
    var response *http.Response
    
    result, err := c.manager.Execute(ctx, c.name, func() (interface{}, error) {
        // Make HTTP request
        url := c.baseURL + path
        req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
        if err != nil {
            return nil, err
        }
        
        resp, err := c.httpClient.Do(req)
        if err != nil {
            return nil, err
        }
        
        // Check if response indicates failure
        if resp.StatusCode >= 500 {
            resp.Body.Close()
            return nil, fmt.Errorf("server error: %d", resp.StatusCode)
        }
        
        return resp, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    response = result.(*http.Response)
    return response, nil
}

func (c *CircuitBreakerClient) Post(ctx context.Context, path string, body []byte) (*http.Response, error) {
    var response *http.Response
    
    result, err := c.manager.Execute(ctx, c.name, func() (interface{}, error) {
        // Make HTTP request
        url := c.baseURL + path
        req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
        if err != nil {
            return nil, err
        }
        
        req.Header.Set("Content-Type", "application/json")
        
        resp, err := c.httpClient.Do(req)
        if err != nil {
            return nil, err
        }
        
        // Check if response indicates failure
        if resp.StatusCode >= 500 {
            resp.Body.Close()
            return nil, fmt.Errorf("server error: %d", resp.StatusCode)
        }
        
        return resp, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    response = result.(*http.Response)
    return response, nil
}
```

### 5. Circuit Breaker for Database Operations

```go
// internal/database/circuit_breaker_db.go
package database

import (
    "context"
    "fmt"
    
    "github.com/anasamu/go-micro-libs/circuitbreaker"
    "github.com/anasamu/go-micro-libs/database"
)

type CircuitBreakerDB struct {
    db      *database.DatabaseManager
    manager *circuitbreaker.Manager
    name    string
}

func NewCircuitBreakerDB(db *database.DatabaseManager, manager *circuitbreaker.Manager, name string) *CircuitBreakerDB {
    return &CircuitBreakerDB{
        db:      db,
        manager: manager,
        name:    name,
    }
}

func (cbd *CircuitBreakerDB) Query(ctx context.Context, provider, query string, args ...interface{}) (*database.Result, error) {
    var result *database.Result
    
    res, err := cbd.manager.Execute(ctx, cbd.name, func() (interface{}, error) {
        // Execute database query
        res, err := cbd.db.Query(ctx, provider, query, args...)
        if err != nil {
            return nil, err
        }
        
        return res, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    result = res.(*database.Result)
    return result, nil
}

func (cbd *CircuitBreakerDB) Exec(ctx context.Context, provider, query string, args ...interface{}) (*database.Result, error) {
    var result *database.Result
    
    res, err := cbd.manager.Execute(ctx, cbd.name, func() (interface{}, error) {
        // Execute database command
        res, err := cbd.db.Exec(ctx, provider, query, args...)
        if err != nil {
            return nil, err
        }
        
        return res, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    result = res.(*database.Result)
    return result, nil
}
```

### 6. Circuit Breaker for Message Broker

```go
// internal/messaging/circuit_breaker_messaging.go
package messaging

import (
    "context"
    
    "github.com/anasamu/go-micro-libs/circuitbreaker"
    "github.com/anasamu/go-micro-libs/messaging"
)

type CircuitBreakerMessaging struct {
    msgManager *messaging.MessagingManager
    manager    *circuitbreaker.Manager
    name       string
}

func NewCircuitBreakerMessaging(msgManager *messaging.MessagingManager, manager *circuitbreaker.Manager, name string) *CircuitBreakerMessaging {
    return &CircuitBreakerMessaging{
        msgManager: msgManager,
        manager:    manager,
        name:       name,
    }
}

func (cbm *CircuitBreakerMessaging) PublishMessage(ctx context.Context, provider string, req *messaging.PublishRequest) (*messaging.PublishResponse, error) {
    var response *messaging.PublishResponse
    
    result, err := cbm.manager.Execute(ctx, cbm.name, func() (interface{}, error) {
        // Publish message
        resp, err := cbm.msgManager.PublishMessage(ctx, provider, req)
        if err != nil {
            return nil, err
        }
        
        return resp, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    response = result.(*messaging.PublishResponse)
    return response, nil
}

func (cbm *CircuitBreakerMessaging) Subscribe(ctx context.Context, provider string, req *messaging.SubscribeRequest) (<-chan *messaging.Message, error) {
    var messageChan <-chan *messaging.Message
    
    result, err := cbm.manager.Execute(ctx, cbm.name, func() (interface{}, error) {
        // Subscribe to messages
        msgChan, err := cbm.msgManager.Subscribe(ctx, provider, req)
        if err != nil {
            return nil, err
        }
        
        return msgChan, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    messageChan = result.(<-chan *messaging.Message)
    return messageChan, nil
}
```

## 🔧 Circuit Breaker Monitoring

### 1. Health Check Endpoint

```go
// internal/handlers/health_handler.go
package handlers

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/anasamu/go-micro-libs/circuitbreaker"
)

type HealthHandler struct {
    manager *circuitbreaker.Manager
}

func NewHealthHandler(manager *circuitbreaker.Manager) *HealthHandler {
    return &HealthHandler{
        manager: manager,
    }
}

func (h *HealthHandler) GetHealth(c *gin.Context) {
    // Get circuit breaker states
    states := make(map[string]string)
    
    // This would need to be implemented to get all circuit breaker names
    // For now, we'll use a simple example
    names := []string{"database", "external-api", "message-broker"}
    
    for _, name := range names {
        state := h.manager.GetState(name)
        states[name] = state.String()
    }
    
    // Check if any circuit breakers are open
    hasOpenBreakers := false
    for _, state := range states {
        if state == "open" {
            hasOpenBreakers = true
            break
        }
    }
    
    if hasOpenBreakers {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "unhealthy",
            "circuit_breakers": states,
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "status": "healthy",
        "circuit_breakers": states,
    })
}

func (h *HealthHandler) GetCircuitBreakerStats(c *gin.Context) {
    name := c.Param("name")
    
    stats := h.manager.GetStats(name)
    if stats == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Circuit breaker not found"})
        return
    }
    
    c.JSON(http.StatusOK, stats)
}

func (h *HealthHandler) ResetCircuitBreaker(c *gin.Context) {
    name := c.Param("name")
    
    h.manager.Reset(name)
    
    c.JSON(http.StatusOK, gin.H{
        "message": fmt.Sprintf("Circuit breaker %s reset", name),
    })
}
```

### 2. Metrics Collection

```go
// internal/monitoring/circuit_breaker_metrics.go
package monitoring

import (
    "time"
    
    "github.com/anasamu/go-micro-libs/circuitbreaker"
    "github.com/anasamu/go-micro-libs/monitoring"
)

type CircuitBreakerMetrics struct {
    manager *circuitbreaker.Manager
    metrics *monitoring.Metrics
}

func NewCircuitBreakerMetrics(manager *circuitbreaker.Manager, metrics *monitoring.Metrics) *CircuitBreakerMetrics {
    return &CircuitBreakerMetrics{
        manager: manager,
        metrics: metrics,
    }
}

func (cbm *CircuitBreakerMetrics) CollectMetrics() {
    // This would typically be called periodically
    // For now, we'll show how to collect metrics for a specific circuit breaker
    
    names := []string{"database", "external-api", "message-broker"}
    
    for _, name := range names {
        stats := cbm.manager.GetStats(name)
        if stats != nil {
            // Record circuit breaker state
            cbm.metrics.SetGauge("circuit_breaker_state", float64(stats.State), map[string]string{
                "name": name,
            })
            
            // Record failure count
            cbm.metrics.SetGauge("circuit_breaker_failure_count", float64(stats.FailureCount), map[string]string{
                "name": name,
            })
            
            // Record success count
            cbm.metrics.SetGauge("circuit_breaker_success_count", float64(stats.SuccessCount), map[string]string{
                "name": name,
            })
            
            // Record request count
            cbm.metrics.SetGauge("circuit_breaker_request_count", float64(stats.RequestCount), map[string]string{
                "name": name,
            })
            
            // Record failure ratio
            cbm.metrics.SetGauge("circuit_breaker_failure_ratio", stats.FailureRatio, map[string]string{
                "name": name,
            })
        }
    }
}
```

## 🔧 Best Practices

### 1. Circuit Breaker Configuration

```go
// Optimal circuit breaker configuration
func getOptimalConfig() *Config {
    return &Config{
        FailureThreshold: 5,        // Open after 5 failures
        Timeout:         30 * time.Second, // Wait 30s before trying again
        MaxRequests:     3,         // Allow 3 test requests in half-open state
        SuccessThreshold: 2,        // Close after 2 successful requests
        FailureRatio:    0.5,       // Open if 50% of requests fail
        MinRequests:     10,        // Need at least 10 requests to make a decision
    }
}
```

### 2. Fallback Mechanisms

```go
// Implement fallback mechanisms
func (c *CircuitBreakerClient) GetWithFallback(ctx context.Context, path string) (*http.Response, error) {
    response, err := c.Get(ctx, path)
    if err != nil {
        // Circuit breaker is open or request failed
        // Implement fallback logic
        return c.getFallbackResponse(ctx, path)
    }
    
    return response, nil
}

func (c *CircuitBreakerClient) getFallbackResponse(ctx context.Context, path string) (*http.Response, error) {
    // Return cached response or default response
    // This is a simplified example
    return &http.Response{
        StatusCode: http.StatusOK,
        Body:       io.NopCloser(strings.NewReader(`{"message": "Service temporarily unavailable"}`)),
    }, nil
}
```

### 3. Circuit Breaker Testing

```go
// Test circuit breaker behavior
func TestCircuitBreaker(t *testing.T) {
    config := &Config{
        FailureThreshold: 3,
        Timeout:         1 * time.Second,
        MaxRequests:     2,
        SuccessThreshold: 1,
        FailureRatio:    0.5,
        MinRequests:     5,
    }
    
    manager := NewManager(config, logger, metrics)
    
    // Test normal operation
    result, err := manager.Execute(ctx, "test", func() (interface{}, error) {
        return "success", nil
    })
    
    assert.NoError(t, err)
    assert.Equal(t, "success", result)
    
    // Test failure threshold
    for i := 0; i < 3; i++ {
        _, err := manager.Execute(ctx, "test", func() (interface{}, error) {
            return nil, fmt.Errorf("test error")
        })
        assert.Error(t, err)
    }
    
    // Circuit breaker should be open now
    state := manager.GetState("test")
    assert.Equal(t, StateOpen, state)
    
    // Requests should be rejected
    _, err = manager.Execute(ctx, "test", func() (interface{}, error) {
        return "success", nil
    })
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "circuit breaker is open")
}
```

### 4. Circuit Breaker Monitoring

```go
// Monitor circuit breaker health
func (m *Manager) MonitorHealth() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        m.mutex.RLock()
        for name, breaker := range m.breakers {
            stats := breaker.GetStats()
            
            // Log circuit breaker state
            m.logger.Info("Circuit breaker stats",
                "name", stats.Name,
                "state", stats.State,
                "failure_count", stats.FailureCount,
                "success_count", stats.SuccessCount,
                "request_count", stats.RequestCount,
                "failure_ratio", stats.FailureRatio,
            )
            
            // Alert if circuit breaker is open for too long
            if stats.State == "open" && time.Since(stats.LastFailureTime) > 5*time.Minute {
                m.logger.Warn("Circuit breaker has been open for too long",
                    "name", name,
                    "duration", time.Since(stats.LastFailureTime),
                )
            }
        }
        m.mutex.RUnlock()
    }
}
```

---

**Circuit Breaker - Resilient and fault-tolerant microservices! 🚀**

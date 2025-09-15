# Failover Implementation

## 🎯 Overview

GoMicroFramework menyediakan sistem failover yang komprehensif untuk memastikan high availability dan reliability dari microservices. Failover mechanism memungkinkan automatic switching ke backup services atau alternative endpoints ketika primary service mengalami failure.

## 🔧 Failover Strategies

### 1. Active-Passive Failover
- Primary service handles all requests
- Secondary service remains idle until failover
- Automatic switching when primary fails
- Manual or automatic failback

### 2. Active-Active Failover
- Multiple services handle requests simultaneously
- Load balancing across all services
- Automatic removal of failed services
- Continuous operation during failures

### 3. Database Failover
- Primary and secondary database instances
- Automatic replication and synchronization
- Read/write splitting
- Connection pooling and health checks

### 4. Service Discovery Failover
- Multiple service instances registration
- Health check monitoring
- Automatic service removal/addition
- Load balancing with failover

## 🔧 Failover Configuration

### 1. Generate Service with Failover

```bash
# Generate service with failover
microframework new user-service --with-failover --with-database=postgres --with-monitoring=prometheus
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

# Failover configuration
failover:
  providers:
    consul:
      enabled: true
      address: "localhost:8500"
      token: "${CONSUL_TOKEN}"
      service_name: "user-service"
      health_check:
        enabled: true
        path: "/health"
        interval: "10s"
        timeout: "3s"
        deregister_after: "30s"
      failover:
        enabled: true
        strategy: "active-passive"
        max_retries: 3
        retry_interval: "5s"
        timeout: "30s"
        
    kubernetes:
      enabled: false
      config_path: "${KUBERNETES_CONFIG}"
      namespace: "default"
      service_name: "user-service"
      health_check:
        enabled: true
        path: "/health"
        interval: "10s"
        timeout: "3s"
      failover:
        enabled: true
        strategy: "active-active"
        max_retries: 3
        retry_interval: "5s"
        timeout: "30s"
        
    static:
      enabled: false
      endpoints:
        - "http://user-service-1:8080"
        - "http://user-service-2:8080"
        - "http://user-service-3:8080"
      health_check:
        enabled: true
        path: "/health"
        interval: "10s"
        timeout: "3s"
      failover:
        enabled: true
        strategy: "round-robin"
        max_retries: 3
        retry_interval: "5s"
        timeout: "30s"

# Database failover
database:
  providers:
    postgresql:
      primary:
        url: "${PRIMARY_DATABASE_URL}"
        max_connections: 100
        max_idle_connections: 10
      secondary:
        url: "${SECONDARY_DATABASE_URL}"
        max_connections: 100
        max_idle_connections: 10
      failover:
        enabled: true
        strategy: "active-passive"
        health_check_interval: "10s"
        timeout: "30s"
        max_retries: 3

# Cache failover
cache:
  providers:
    redis:
      primary:
        url: "${PRIMARY_REDIS_URL}"
        db: 0
        pool_size: 10
      secondary:
        url: "${SECONDARY_REDIS_URL}"
        db: 0
        pool_size: 10
      failover:
        enabled: true
        strategy: "active-passive"
        health_check_interval: "10s"
        timeout: "30s"
        max_retries: 3
```

## 🔧 Failover Implementation

### 1. Failover Manager

```go
// internal/failover/manager.go
package failover

import (
    "context"
    "fmt"
    "sync"
    "time"
    
    "github.com/anasamu/go-micro-libs/failover"
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/logging"
)

type Manager struct {
    providers    map[string]FailoverProvider
    config       *Config
    logger       *logging.Logger
    metrics      *monitoring.Metrics
    healthChecks map[string]*HealthCheck
    mutex        sync.RWMutex
}

type Config struct {
    Provider        string        `yaml:"provider"`
    Strategy        string        `yaml:"strategy"`
    MaxRetries      int           `yaml:"max_retries"`
    RetryInterval   time.Duration `yaml:"retry_interval"`
    Timeout         time.Duration `yaml:"timeout"`
    HealthCheck     HealthCheckConfig `yaml:"health_check"`
}

type HealthCheckConfig struct {
    Enabled          bool          `yaml:"enabled"`
    Path             string        `yaml:"path"`
    Interval         time.Duration `yaml:"interval"`
    Timeout          time.Duration `yaml:"timeout"`
    DeregisterAfter  time.Duration `yaml:"deregister_after"`
}

type FailoverProvider interface {
    GetEndpoints(ctx context.Context) ([]*Endpoint, error)
    RegisterEndpoint(ctx context.Context, endpoint *Endpoint) error
    DeregisterEndpoint(ctx context.Context, endpoint *Endpoint) error
    HealthCheck(ctx context.Context, endpoint *Endpoint) error
}

type Endpoint struct {
    ID       string            `json:"id"`
    Address  string            `json:"address"`
    Port     int               `json:"port"`
    Metadata map[string]string `json:"metadata"`
    Status   EndpointStatus    `json:"status"`
    LastSeen time.Time         `json:"last_seen"`
}

type EndpointStatus int

const (
    StatusHealthy EndpointStatus = iota
    StatusUnhealthy
    StatusUnknown
)

func (s EndpointStatus) String() string {
    switch s {
    case StatusHealthy:
        return "healthy"
    case StatusUnhealthy:
        return "unhealthy"
    case StatusUnknown:
        return "unknown"
    default:
        return "unknown"
    }
}

type HealthCheck struct {
    endpoint   *Endpoint
    interval   time.Duration
    timeout    time.Duration
    provider   FailoverProvider
    stopChan   chan struct{}
    logger     *logging.Logger
    metrics    *monitoring.Metrics
}

func NewManager(config *Config, logger *logging.Logger, metrics *monitoring.Metrics) *Manager {
    return &Manager{
        providers:    make(map[string]FailoverProvider),
        config:       config,
        logger:       logger,
        metrics:      metrics,
        healthChecks: make(map[string]*HealthCheck),
    }
}

func (m *Manager) RegisterProvider(name string, provider FailoverProvider) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    m.providers[name] = provider
    m.logger.Info("Failover provider registered", "name", name)
}

func (m *Manager) GetEndpoints(ctx context.Context, providerName string) ([]*Endpoint, error) {
    m.mutex.RLock()
    provider, exists := m.providers[providerName]
    m.mutex.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("provider %s not found", providerName)
    }
    
    endpoints, err := provider.GetEndpoints(ctx)
    if err != nil {
        m.logger.Error("Failed to get endpoints", "provider", providerName, "error", err)
        return nil, err
    }
    
    // Filter healthy endpoints
    healthyEndpoints := make([]*Endpoint, 0)
    for _, endpoint := range endpoints {
        if endpoint.Status == StatusHealthy {
            healthyEndpoints = append(healthyEndpoints, endpoint)
        }
    }
    
    return healthyEndpoints, nil
}

func (m *Manager) StartHealthChecks(ctx context.Context, providerName string) error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    provider, exists := m.providers[providerName]
    if !exists {
        return fmt.Errorf("provider %s not found", providerName)
    }
    
    // Get endpoints
    endpoints, err := provider.GetEndpoints(ctx)
    if err != nil {
        return err
    }
    
    // Start health checks for each endpoint
    for _, endpoint := range endpoints {
        healthCheck := &HealthCheck{
            endpoint: endpoint,
            interval: m.config.HealthCheck.Interval,
            timeout:  m.config.HealthCheck.Timeout,
            provider: provider,
            stopChan: make(chan struct{}),
            logger:   m.logger,
            metrics:  m.metrics,
        }
        
        m.healthChecks[endpoint.ID] = healthCheck
        go healthCheck.Start(ctx)
    }
    
    return nil
}

func (m *Manager) StopHealthChecks() {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    for _, healthCheck := range m.healthChecks {
        healthCheck.Stop()
    }
}

func (m *Manager) GetHealthyEndpoints(ctx context.Context, providerName string) ([]*Endpoint, error) {
    endpoints, err := m.GetEndpoints(ctx, providerName)
    if err != nil {
        return nil, err
    }
    
    // Filter only healthy endpoints
    healthyEndpoints := make([]*Endpoint, 0)
    for _, endpoint := range endpoints {
        if endpoint.Status == StatusHealthy {
            healthyEndpoints = append(healthyEndpoints, endpoint)
        }
    }
    
    return healthyEndpoints, nil
}
```

### 2. Health Check Implementation

```go
// internal/failover/health_check.go
package failover

import (
    "context"
    "net/http"
    "time"
    
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/logging"
)

func (hc *HealthCheck) Start(ctx context.Context) {
    ticker := time.NewTicker(hc.interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            hc.performHealthCheck(ctx)
        case <-hc.stopChan:
            return
        case <-ctx.Done():
            return
        }
    }
}

func (hc *HealthCheck) Stop() {
    close(hc.stopChan)
}

func (hc *HealthCheck) performHealthCheck(ctx context.Context) {
    start := time.Now()
    
    // Perform health check
    err := hc.provider.HealthCheck(ctx, hc.endpoint)
    
    duration := time.Since(start)
    
    if err != nil {
        // Health check failed
        if hc.endpoint.Status == StatusHealthy {
            hc.endpoint.Status = StatusUnhealthy
            hc.logger.Warn("Endpoint became unhealthy",
                "endpoint", hc.endpoint.Address,
                "error", err,
                "duration", duration,
            )
            
            // Record metrics
            hc.metrics.IncrementCounter("failover_endpoint_unhealthy_total", map[string]string{
                "endpoint": hc.endpoint.Address,
            })
        }
    } else {
        // Health check succeeded
        if hc.endpoint.Status == StatusUnhealthy {
            hc.endpoint.Status = StatusHealthy
            hc.logger.Info("Endpoint became healthy",
                "endpoint", hc.endpoint.Address,
                "duration", duration,
            )
            
            // Record metrics
            hc.metrics.IncrementCounter("failover_endpoint_healthy_total", map[string]string{
                "endpoint": hc.endpoint.Address,
            })
        }
    }
    
    // Update last seen time
    hc.endpoint.LastSeen = time.Now()
    
    // Record health check duration
    hc.metrics.RecordHistogram("failover_health_check_duration_seconds",
        duration.Seconds(), map[string]string{
            "endpoint": hc.endpoint.Address,
            "status":   hc.endpoint.Status.String(),
        })
}
```

### 3. Consul Failover Provider

```go
// internal/failover/consul_provider.go
package failover

import (
    "context"
    "fmt"
    "net/http"
    "time"
    
    "github.com/hashicorp/consul/api"
    "github.com/anasamu/go-micro-libs/logging"
)

type ConsulProvider struct {
    client      *api.Client
    serviceName string
    logger      *logging.Logger
}

func NewConsulProvider(address, token, serviceName string, logger *logging.Logger) (*ConsulProvider, error) {
    config := api.DefaultConfig()
    config.Address = address
    config.Token = token
    
    client, err := api.NewClient(config)
    if err != nil {
        return nil, err
    }
    
    return &ConsulProvider{
        client:      client,
        serviceName: serviceName,
        logger:      logger,
    }, nil
}

func (cp *ConsulProvider) GetEndpoints(ctx context.Context) ([]*Endpoint, error) {
    services, _, err := cp.client.Health().Service(cp.serviceName, "", true, nil)
    if err != nil {
        return nil, err
    }
    
    endpoints := make([]*Endpoint, len(services))
    for i, service := range services {
        endpoint := &Endpoint{
            ID:      service.Service.ID,
            Address: service.Service.Address,
            Port:    service.Service.Port,
            Metadata: service.Service.Meta,
            Status:  StatusHealthy,
            LastSeen: time.Now(),
        }
        
        // Check if service is healthy
        if service.Checks.AggregatedStatus() != api.HealthPassing {
            endpoint.Status = StatusUnhealthy
        }
        
        endpoints[i] = endpoint
    }
    
    return endpoints, nil
}

func (cp *ConsulProvider) RegisterEndpoint(ctx context.Context, endpoint *Endpoint) error {
    registration := &api.AgentServiceRegistration{
        ID:      endpoint.ID,
        Name:    cp.serviceName,
        Port:    endpoint.Port,
        Address: endpoint.Address,
        Meta:    endpoint.Metadata,
        Check: &api.AgentServiceCheck{
            HTTP:                           fmt.Sprintf("http://%s:%d/health", endpoint.Address, endpoint.Port),
            Timeout:                        "3s",
            Interval:                       "10s",
            DeregisterCriticalServiceAfter: "30s",
        },
    }
    
    err := cp.client.Agent().ServiceRegister(registration)
    if err != nil {
        return err
    }
    
    cp.logger.Info("Endpoint registered with Consul",
        "endpoint", endpoint.Address,
        "port", endpoint.Port,
    )
    
    return nil
}

func (cp *ConsulProvider) DeregisterEndpoint(ctx context.Context, endpoint *Endpoint) error {
    err := cp.client.Agent().ServiceDeregister(endpoint.ID)
    if err != nil {
        return err
    }
    
    cp.logger.Info("Endpoint deregistered from Consul",
        "endpoint", endpoint.Address,
        "port", endpoint.Port,
    )
    
    return nil
}

func (cp *ConsulProvider) HealthCheck(ctx context.Context, endpoint *Endpoint) error {
    url := fmt.Sprintf("http://%s:%d/health", endpoint.Address, endpoint.Port)
    
    client := &http.Client{
        Timeout: 3 * time.Second,
    }
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return err
    }
    
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
    }
    
    return nil
}
```

### 4. Failover Client

```go
// internal/clients/failover_client.go
package clients

import (
    "context"
    "fmt"
    "net/http"
    "time"
    
    "github.com/anasamu/go-micro-libs/failover"
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/logging"
)

type FailoverClient struct {
    manager     *failover.Manager
    httpClient  *http.Client
    providerName string
    strategy    string
    maxRetries  int
    retryInterval time.Duration
    logger      *logging.Logger
    metrics     *monitoring.Metrics
}

func NewFailoverClient(manager *failover.Manager, httpClient *http.Client, providerName, strategy string, maxRetries int, retryInterval time.Duration, logger *logging.Logger, metrics *monitoring.Metrics) *FailoverClient {
    return &FailoverClient{
        manager:      manager,
        httpClient:   httpClient,
        providerName: providerName,
        strategy:     strategy,
        maxRetries:   maxRetries,
        retryInterval: retryInterval,
        logger:       logger,
        metrics:      metrics,
    }
}

func (fc *FailoverClient) Get(ctx context.Context, path string) (*http.Response, error) {
    return fc.executeWithFailover(ctx, "GET", path, nil)
}

func (fc *FailoverClient) Post(ctx context.Context, path string, body []byte) (*http.Response, error) {
    return fc.executeWithFailover(ctx, "POST", path, body)
}

func (fc *FailoverClient) Put(ctx context.Context, path string, body []byte) (*http.Response, error) {
    return fc.executeWithFailover(ctx, "PUT", path, body)
}

func (fc *FailoverClient) Delete(ctx context.Context, path string) (*http.Response, error) {
    return fc.executeWithFailover(ctx, "DELETE", path, nil)
}

func (fc *FailoverClient) executeWithFailover(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
    var lastErr error
    
    for attempt := 0; attempt < fc.maxRetries; attempt++ {
        // Get healthy endpoints
        endpoints, err := fc.manager.GetHealthyEndpoints(ctx, fc.providerName)
        if err != nil {
            fc.logger.Error("Failed to get healthy endpoints", "error", err)
            return nil, err
        }
        
        if len(endpoints) == 0 {
            return nil, fmt.Errorf("no healthy endpoints available")
        }
        
        // Select endpoint based on strategy
        endpoint := fc.selectEndpoint(endpoints, attempt)
        
        // Execute request
        response, err := fc.executeRequest(ctx, method, path, body, endpoint)
        if err == nil {
            // Request successful
            fc.metrics.IncrementCounter("failover_requests_success_total", map[string]string{
                "endpoint": endpoint.Address,
                "method":   method,
            })
            return response, nil
        }
        
        // Request failed
        lastErr = err
        fc.logger.Warn("Request failed, trying next endpoint",
            "endpoint", endpoint.Address,
            "attempt", attempt+1,
            "error", err,
        )
        
        fc.metrics.IncrementCounter("failover_requests_failed_total", map[string]string{
            "endpoint": endpoint.Address,
            "method":   method,
        })
        
        // Wait before retry
        if attempt < fc.maxRetries-1 {
            time.Sleep(fc.retryInterval)
        }
    }
    
    return nil, fmt.Errorf("all endpoints failed: %w", lastErr)
}

func (fc *FailoverClient) selectEndpoint(endpoints []*failover.Endpoint, attempt int) *failover.Endpoint {
    switch fc.strategy {
    case "round-robin":
        return endpoints[attempt%len(endpoints)]
    case "random":
        return endpoints[time.Now().UnixNano()%int64(len(endpoints))]
    case "first":
        return endpoints[0]
    default:
        return endpoints[attempt%len(endpoints)]
    }
}

func (fc *FailoverClient) executeRequest(ctx context.Context, method, path string, body []byte, endpoint *failover.Endpoint) (*http.Response, error) {
    url := fmt.Sprintf("http://%s:%d%s", endpoint.Address, endpoint.Port, path)
    
    var req *http.Request
    var err error
    
    if body != nil {
        req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
        if err != nil {
            return nil, err
        }
        req.Header.Set("Content-Type", "application/json")
    } else {
        req, err = http.NewRequestWithContext(ctx, method, url, nil)
        if err != nil {
            return nil, err
        }
    }
    
    resp, err := fc.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    
    // Check if response indicates failure
    if resp.StatusCode >= 500 {
        resp.Body.Close()
        return nil, fmt.Errorf("server error: %d", resp.StatusCode)
    }
    
    return resp, nil
}
```

### 5. Database Failover

```go
// internal/database/failover_db.go
package database

import (
    "context"
    "fmt"
    "sync"
    "time"
    
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/logging"
)

type FailoverDB struct {
    primary    *database.DatabaseManager
    secondary  *database.DatabaseManager
    strategy   string
    mutex      sync.RWMutex
    logger     *logging.Logger
    metrics    *monitoring.Metrics
}

func NewFailoverDB(primary, secondary *database.DatabaseManager, strategy string, logger *logging.Logger, metrics *monitoring.Metrics) *FailoverDB {
    return &FailoverDB{
        primary:   primary,
        secondary: secondary,
        strategy:  strategy,
        logger:    logger,
        metrics:   metrics,
    }
}

func (fdb *FailoverDB) Query(ctx context.Context, provider, query string, args ...interface{}) (*database.Result, error) {
    // Try primary first
    result, err := fdb.primary.Query(ctx, provider, query, args...)
    if err == nil {
        fdb.metrics.IncrementCounter("database_queries_success_total", map[string]string{
            "database": "primary",
        })
        return result, nil
    }
    
    // Primary failed, try secondary
    fdb.logger.Warn("Primary database query failed, trying secondary", "error", err)
    
    result, err = fdb.secondary.Query(ctx, provider, query, args...)
    if err != nil {
        fdb.metrics.IncrementCounter("database_queries_failed_total", map[string]string{
            "database": "secondary",
        })
        return nil, fmt.Errorf("both primary and secondary databases failed: %w", err)
    }
    
    fdb.metrics.IncrementCounter("database_queries_success_total", map[string]string{
        "database": "secondary",
    })
    
    return result, nil
}

func (fdb *FailoverDB) Exec(ctx context.Context, provider, query string, args ...interface{}) (*database.Result, error) {
    // For write operations, try primary first
    result, err := fdb.primary.Exec(ctx, provider, query, args...)
    if err == nil {
        fdb.metrics.IncrementCounter("database_executions_success_total", map[string]string{
            "database": "primary",
        })
        return result, nil
    }
    
    // Primary failed, try secondary
    fdb.logger.Warn("Primary database execution failed, trying secondary", "error", err)
    
    result, err = fdb.secondary.Exec(ctx, provider, query, args...)
    if err != nil {
        fdb.metrics.IncrementCounter("database_executions_failed_total", map[string]string{
            "database": "secondary",
        })
        return nil, fmt.Errorf("both primary and secondary databases failed: %w", err)
    }
    
    fdb.metrics.IncrementCounter("database_executions_success_total", map[string]string{
        "database": "secondary",
    })
    
    return result, nil
}

func (fdb *FailoverDB) HealthCheck(ctx context.Context) error {
    // Check primary database
    err := fdb.primary.HealthCheck(ctx)
    if err != nil {
        fdb.logger.Warn("Primary database health check failed", "error", err)
        
        // Check secondary database
        err = fdb.secondary.HealthCheck(ctx)
        if err != nil {
            fdb.logger.Error("Secondary database health check failed", "error", err)
            return fmt.Errorf("both primary and secondary databases are unhealthy")
        }
    }
    
    return nil
}
```

### 6. Cache Failover

```go
// internal/cache/failover_cache.go
package cache

import (
    "context"
    "fmt"
    "sync"
    "time"
    
    "github.com/anasamu/go-micro-libs/cache"
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/logging"
)

type FailoverCache struct {
    primary    *cache.CacheManager
    secondary  *cache.CacheManager
    strategy   string
    mutex      sync.RWMutex
    logger     *logging.Logger
    metrics    *monitoring.Metrics
}

func NewFailoverCache(primary, secondary *cache.CacheManager, strategy string, logger *logging.Logger, metrics *monitoring.Metrics) *FailoverCache {
    return &FailoverCache{
        primary:   primary,
        secondary: secondary,
        strategy:  strategy,
        logger:    logger,
        metrics:   metrics,
    }
}

func (fc *FailoverCache) Get(ctx context.Context, key string) (interface{}, error) {
    // Try primary first
    value, err := fc.primary.Get(ctx, key)
    if err == nil {
        fc.metrics.IncrementCounter("cache_gets_success_total", map[string]string{
            "cache": "primary",
        })
        return value, nil
    }
    
    // Primary failed, try secondary
    fc.logger.Warn("Primary cache get failed, trying secondary", "key", key, "error", err)
    
    value, err = fc.secondary.Get(ctx, key)
    if err != nil {
        fc.metrics.IncrementCounter("cache_gets_failed_total", map[string]string{
            "cache": "secondary",
        })
        return nil, fmt.Errorf("both primary and secondary caches failed: %w", err)
    }
    
    fc.metrics.IncrementCounter("cache_gets_success_total", map[string]string{
        "cache": "secondary",
    })
    
    return value, nil
}

func (fc *FailoverCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    // Set in both primary and secondary
    var primaryErr, secondaryErr error
    
    // Set in primary
    primaryErr = fc.primary.Set(ctx, key, value, ttl)
    if primaryErr != nil {
        fc.logger.Warn("Primary cache set failed", "key", key, "error", primaryErr)
    }
    
    // Set in secondary
    secondaryErr = fc.secondary.Set(ctx, key, value, ttl)
    if secondaryErr != nil {
        fc.logger.Warn("Secondary cache set failed", "key", key, "error", secondaryErr)
    }
    
    // If both failed, return error
    if primaryErr != nil && secondaryErr != nil {
        fc.metrics.IncrementCounter("cache_sets_failed_total", map[string]string{
            "cache": "both",
        })
        return fmt.Errorf("both primary and secondary caches failed: primary=%v, secondary=%v", primaryErr, secondaryErr)
    }
    
    // Record success
    if primaryErr == nil {
        fc.metrics.IncrementCounter("cache_sets_success_total", map[string]string{
            "cache": "primary",
        })
    }
    if secondaryErr == nil {
        fc.metrics.IncrementCounter("cache_sets_success_total", map[string]string{
            "cache": "secondary",
        })
    }
    
    return nil
}

func (fc *FailoverCache) Delete(ctx context.Context, key string) error {
    // Delete from both primary and secondary
    var primaryErr, secondaryErr error
    
    // Delete from primary
    primaryErr = fc.primary.Delete(ctx, key)
    if primaryErr != nil {
        fc.logger.Warn("Primary cache delete failed", "key", key, "error", primaryErr)
    }
    
    // Delete from secondary
    secondaryErr = fc.secondary.Delete(ctx, key)
    if secondaryErr != nil {
        fc.logger.Warn("Secondary cache delete failed", "key", key, "error", secondaryErr)
    }
    
    // If both failed, return error
    if primaryErr != nil && secondaryErr != nil {
        fc.metrics.IncrementCounter("cache_deletes_failed_total", map[string]string{
            "cache": "both",
        })
        return fmt.Errorf("both primary and secondary caches failed: primary=%v, secondary=%v", primaryErr, secondaryErr)
    }
    
    // Record success
    if primaryErr == nil {
        fc.metrics.IncrementCounter("cache_deletes_success_total", map[string]string{
            "cache": "primary",
        })
    }
    if secondaryErr == nil {
        fc.metrics.IncrementCounter("cache_deletes_success_total", map[string]string{
            "cache": "secondary",
        })
    }
    
    return nil
}

func (fc *FailoverCache) HealthCheck(ctx context.Context) error {
    // Check primary cache
    err := fc.primary.HealthCheck(ctx)
    if err != nil {
        fc.logger.Warn("Primary cache health check failed", "error", err)
        
        // Check secondary cache
        err = fc.secondary.HealthCheck(ctx)
        if err != nil {
            fc.logger.Error("Secondary cache health check failed", "error", err)
            return fmt.Errorf("both primary and secondary caches are unhealthy")
        }
    }
    
    return nil
}
```

## 🔧 Best Practices

### 1. Failover Configuration

```go
// Optimal failover configuration
func getOptimalFailoverConfig() *Config {
    return &Config{
        Strategy:      "active-passive",  // Use active-passive for critical services
        MaxRetries:    3,                // Retry up to 3 times
        RetryInterval: 5 * time.Second,  // Wait 5s between retries
        Timeout:       30 * time.Second, // 30s timeout for operations
        HealthCheck: HealthCheckConfig{
            Enabled:         true,
            Path:           "/health",
            Interval:       10 * time.Second,
            Timeout:        3 * time.Second,
            DeregisterAfter: 30 * time.Second,
        },
    }
}
```

### 2. Monitoring and Alerting

```go
// Monitor failover events
func (m *Manager) MonitorFailover() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        m.mutex.RLock()
        for name, healthCheck := range m.healthChecks {
            if healthCheck.endpoint.Status == StatusUnhealthy {
                // Alert if endpoint has been unhealthy for too long
                if time.Since(healthCheck.endpoint.LastSeen) > 5*time.Minute {
                    m.logger.Warn("Endpoint has been unhealthy for too long",
                        "endpoint", name,
                        "duration", time.Since(healthCheck.endpoint.LastSeen),
                    )
                    
                    // Send alert
                    m.sendAlert(name, "endpoint_unhealthy")
                }
            }
        }
        m.mutex.RUnlock()
    }
}

func (m *Manager) sendAlert(endpoint, alertType string) {
    // Implement alerting logic
    m.logger.Info("Sending alert", "endpoint", endpoint, "type", alertType)
}
```

### 3. Testing Failover

```go
// Test failover behavior
func TestFailover(t *testing.T) {
    config := &Config{
        Strategy:      "active-passive",
        MaxRetries:    3,
        RetryInterval: 1 * time.Second,
        Timeout:       5 * time.Second,
    }
    
    manager := NewManager(config, logger, metrics)
    
    // Test normal operation
    endpoints, err := manager.GetHealthyEndpoints(ctx, "test-provider")
    assert.NoError(t, err)
    assert.NotEmpty(t, endpoints)
    
    // Test failover
    // Simulate endpoint failure
    // Verify failover behavior
}
```

### 4. Graceful Degradation

```go
// Implement graceful degradation
func (fc *FailoverClient) GetWithFallback(ctx context.Context, path string) (*http.Response, error) {
    response, err := fc.Get(ctx, path)
    if err != nil {
        // All endpoints failed, return fallback response
        return fc.getFallbackResponse(ctx, path)
    }
    
    return response, nil
}

func (fc *FailoverClient) getFallbackResponse(ctx context.Context, path string) (*http.Response, error) {
    // Return cached response or default response
    return &http.Response{
        StatusCode: http.StatusOK,
        Body:       io.NopCloser(strings.NewReader(`{"message": "Service temporarily unavailable"}`)),
    }, nil
}
```

---

**Failover - High availability and reliability for microservices! 🚀**

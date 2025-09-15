# Middleware Implementation

## 🎯 Overview

GoMicroFramework menyediakan sistem middleware yang komprehensif untuk menangani berbagai aspek dari request processing seperti authentication, authorization, logging, monitoring, rate limiting, circuit breaking, dan caching. Middleware dapat dikombinasikan untuk membangun pipeline yang powerful dan fleksibel.

## 🔧 Supported Middleware Types

### 1. Authentication Middleware
- JWT token validation
- OAuth2 token validation
- Session-based authentication
- API key validation

### 2. Authorization Middleware
- Role-based access control
- Permission-based access control
- Resource-based access control
- Policy-based access control

### 3. Logging Middleware
- Request/response logging
- Structured logging
- Correlation ID tracking
- Performance logging

### 4. Monitoring Middleware
- Metrics collection
- Distributed tracing
- Health checks
- Performance monitoring

### 5. Rate Limiting Middleware
- Token bucket algorithm
- Sliding window algorithm
- IP-based rate limiting
- User-based rate limiting

### 6. Circuit Breaker Middleware
- Failure detection
- Automatic recovery
- Fallback mechanisms
- Health monitoring

### 7. Caching Middleware
- Response caching
- Query result caching
- Session caching
- Distributed caching

## 🔧 Middleware Setup

### 1. Generate Service with Middleware

```bash
# Generate service with middleware
microframework new user-service --with-auth=jwt --with-database=postgres --with-monitoring=prometheus
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

# Middleware configuration
middleware:
  auth:
    enabled: true
    provider: "jwt"
    secret: "${JWT_SECRET}"
    expiration: "24h"
    
  authorization:
    enabled: true
    provider: "rbac"
    default_role: "user"
    
  logging:
    enabled: true
    level: "info"
    format: "json"
    include_body: false
    include_headers: true
    
  monitoring:
    enabled: true
    metrics: true
    tracing: true
    health_checks: true
    
  rate_limit:
    enabled: true
    provider: "redis"
    requests_per_minute: 100
    burst: 10
    key_strategy: "ip"
    
  circuit_breaker:
    enabled: true
    failure_threshold: 5
    timeout: 30s
    max_requests: 3
    
  caching:
    enabled: true
    provider: "redis"
    ttl: "1h"
    key_strategy: "path"
    
  compression:
    enabled: true
    algorithm: "gzip"
    level: 6
    
  security:
    enabled: true
    cors:
      enabled: true
      origins: ["*"]
      methods: ["GET", "POST", "PUT", "DELETE"]
      headers: ["Content-Type", "Authorization"]
    csrf:
      enabled: false
    xss:
      enabled: true
    helmet:
      enabled: true

# Database for middleware data
database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100

# Cache for middleware
cache:
  providers:
    redis:
      url: "${REDIS_URL}"
      db: 0
      pool_size: 10
```

## 🔧 Middleware Implementation

### 1. Middleware Manager

```go
// internal/middleware/manager.go
package middleware

import (
    "context"
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/anasamu/go-micro-libs/middleware"
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/logging"
)

type MiddlewareManager struct {
    middlewares []gin.HandlerFunc
    config      *Config
    logger      *logging.Logger
    metrics     *monitoring.Metrics
}

type Config struct {
    Auth           AuthConfig           `yaml:"auth"`
    Authorization  AuthorizationConfig  `yaml:"authorization"`
    Logging        LoggingConfig        `yaml:"logging"`
    Monitoring     MonitoringConfig     `yaml:"monitoring"`
    RateLimit      RateLimitConfig      `yaml:"rate_limit"`
    CircuitBreaker CircuitBreakerConfig `yaml:"circuit_breaker"`
    Caching        CachingConfig        `yaml:"caching"`
    Compression    CompressionConfig    `yaml:"compression"`
    Security       SecurityConfig       `yaml:"security"`
}

func NewMiddlewareManager(config *Config, logger *logging.Logger, metrics *monitoring.Metrics) *MiddlewareManager {
    return &MiddlewareManager{
        middlewares: make([]gin.HandlerFunc, 0),
        config:      config,
        logger:      logger,
        metrics:     metrics,
    }
}

func (m *MiddlewareManager) SetupMiddleware() []gin.HandlerFunc {
    // Setup middleware in order
    if m.config.Security.Enabled {
        m.addSecurityMiddleware()
    }
    
    if m.config.Logging.Enabled {
        m.addLoggingMiddleware()
    }
    
    if m.config.Monitoring.Enabled {
        m.addMonitoringMiddleware()
    }
    
    if m.config.RateLimit.Enabled {
        m.addRateLimitMiddleware()
    }
    
    if m.config.CircuitBreaker.Enabled {
        m.addCircuitBreakerMiddleware()
    }
    
    if m.config.Caching.Enabled {
        m.addCachingMiddleware()
    }
    
    if m.config.Compression.Enabled {
        m.addCompressionMiddleware()
    }
    
    if m.config.Auth.Enabled {
        m.addAuthMiddleware()
    }
    
    if m.config.Authorization.Enabled {
        m.addAuthorizationMiddleware()
    }
    
    return m.middlewares
}

func (m *MiddlewareManager) addSecurityMiddleware() {
    // CORS middleware
    if m.config.Security.CORS.Enabled {
        m.middlewares = append(m.middlewares, m.corsMiddleware())
    }
    
    // Helmet middleware
    if m.config.Security.Helmet.Enabled {
        m.middlewares = append(m.middlewares, m.helmetMiddleware())
    }
    
    // XSS protection
    if m.config.Security.XSS.Enabled {
        m.middlewares = append(m.middlewares, m.xssMiddleware())
    }
}

func (m *MiddlewareManager) addLoggingMiddleware() {
    m.middlewares = append(m.middlewares, m.loggingMiddleware())
}

func (m *MiddlewareManager) addMonitoringMiddleware() {
    m.middlewares = append(m.middlewares, m.monitoringMiddleware())
}

func (m *MiddlewareManager) addRateLimitMiddleware() {
    m.middlewares = append(m.middlewares, m.rateLimitMiddleware())
}

func (m *MiddlewareManager) addCircuitBreakerMiddleware() {
    m.middlewares = append(m.middlewares, m.circuitBreakerMiddleware())
}

func (m *MiddlewareManager) addCachingMiddleware() {
    m.middlewares = append(m.middlewares, m.cachingMiddleware())
}

func (m *MiddlewareManager) addCompressionMiddleware() {
    m.middlewares = append(m.middlewares, m.compressionMiddleware())
}

func (m *MiddlewareManager) addAuthMiddleware() {
    m.middlewares = append(m.middlewares, m.authMiddleware())
}

func (m *MiddlewareManager) addAuthorizationMiddleware() {
    m.middlewares = append(m.middlewares, m.authorizationMiddleware())
}
```

### 2. Logging Middleware

```go
// internal/middleware/logging.go
package middleware

import (
    "bytes"
    "io"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/anasamu/go-micro-libs/logging"
    "github.com/anasamu/go-micro-libs/monitoring"
)

type responseWriter struct {
    gin.ResponseWriter
    body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
    w.body.Write(b)
    return w.ResponseWriter.Write(b)
}

func (m *MiddlewareManager) loggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        // Generate correlation ID
        correlationID := generateCorrelationID()
        c.Set("correlation_id", correlationID)
        
        // Set correlation ID in response header
        c.Header("X-Correlation-ID", correlationID)
        
        // Log request
        m.logger.Info("Request started",
            "method", c.Request.Method,
            "path", c.Request.URL.Path,
            "query", c.Request.URL.RawQuery,
            "ip", c.ClientIP(),
            "user_agent", c.Request.UserAgent(),
            "correlation_id", correlationID,
        )
        
        // Capture response body
        var responseBody bytes.Buffer
        writer := &responseWriter{
            ResponseWriter: c.Writer,
            body:          &responseBody,
        }
        c.Writer = writer
        
        // Process request
        c.Next()
        
        // Calculate duration
        duration := time.Since(start)
        
        // Log response
        m.logger.Info("Request completed",
            "method", c.Request.Method,
            "path", c.Request.URL.Path,
            "status", c.Writer.Status(),
            "duration", duration,
            "size", c.Writer.Size(),
            "correlation_id", correlationID,
        )
        
        // Record metrics
        m.metrics.RecordHistogram("http_request_duration_seconds", 
            duration.Seconds(), map[string]string{
                "method": c.Request.Method,
                "path":   c.Request.URL.Path,
                "status": string(rune(c.Writer.Status())),
            })
        
        m.metrics.IncrementCounter("http_requests_total", map[string]string{
            "method": c.Request.Method,
            "path":   c.Request.URL.Path,
            "status": string(rune(c.Writer.Status())),
        })
    }
}
```

### 3. Monitoring Middleware

```go
// internal/middleware/monitoring.go
package middleware

import (
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/anasamu/go-micro-libs/monitoring"
)

func (m *MiddlewareManager) monitoringMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        // Start span for distributed tracing
        span := m.metrics.StartSpan(c.Request.Context(), "http_request",
            map[string]string{
                "method": c.Request.Method,
                "path":   c.Request.URL.Path,
            })
        defer span.Finish()
        
        // Set span in context
        c.Set("span", span)
        
        // Process request
        c.Next()
        
        // Record metrics
        duration := time.Since(start)
        
        m.metrics.RecordHistogram("http_request_duration_seconds", 
            duration.Seconds(), map[string]string{
                "method": c.Request.Method,
                "path":   c.Request.URL.Path,
                "status": string(rune(c.Writer.Status())),
            })
        
        m.metrics.IncrementCounter("http_requests_total", map[string]string{
            "method": c.Request.Method,
            "path":   c.Request.URL.Path,
            "status": string(rune(c.Writer.Status())),
        })
        
        // Set span tags
        span.SetTag("http.status_code", c.Writer.Status())
        span.SetTag("http.method", c.Request.Method)
        span.SetTag("http.url", c.Request.URL.String())
    }
}
```

### 4. Rate Limiting Middleware

```go
// internal/middleware/rate_limit.go
package middleware

import (
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/anasamu/go-micro-libs/ratelimit"
    "github.com/anasamu/go-micro-libs/monitoring"
)

func (m *MiddlewareManager) rateLimitMiddleware() gin.HandlerFunc {
    rateLimiter := ratelimit.NewRateLimiter(
        m.config.RateLimit.RequestsPerMinute,
        m.config.RateLimit.Burst,
        m.config.RateLimit.KeyStrategy,
    )
    
    return func(c *gin.Context) {
        // Get rate limit key
        key := m.getRateLimitKey(c)
        
        // Check rate limit
        allowed, err := rateLimiter.Allow(c.Request.Context(), key)
        if err != nil {
            m.logger.Error("Rate limit check failed", "error", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit check failed"})
            c.Abort()
            return
        }
        
        if !allowed {
            // Record rate limit exceeded
            m.metrics.IncrementCounter("rate_limit_exceeded_total", map[string]string{
                "key": key,
            })
            
            c.Header("X-RateLimit-Limit", string(rune(m.config.RateLimit.RequestsPerMinute)))
            c.Header("X-RateLimit-Remaining", "0")
            c.Header("X-RateLimit-Reset", string(rune(time.Now().Add(time.Minute).Unix())))
            
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded",
                "retry_after": 60,
            })
            c.Abort()
            return
        }
        
        // Set rate limit headers
        c.Header("X-RateLimit-Limit", string(rune(m.config.RateLimit.RequestsPerMinute)))
        c.Header("X-RateLimit-Remaining", string(rune(rateLimiter.GetRemaining(key))))
        c.Header("X-RateLimit-Reset", string(rune(time.Now().Add(time.Minute).Unix())))
        
        c.Next()
    }
}

func (m *MiddlewareManager) getRateLimitKey(c *gin.Context) string {
    switch m.config.RateLimit.KeyStrategy {
    case "ip":
        return c.ClientIP()
    case "user":
        userID, exists := c.Get("user_id")
        if exists {
            return userID.(string)
        }
        return c.ClientIP()
    case "path":
        return c.Request.URL.Path
    default:
        return c.ClientIP()
    }
}
```

### 5. Circuit Breaker Middleware

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

func (m *MiddlewareManager) circuitBreakerMiddleware() gin.HandlerFunc {
    circuitBreaker := circuitbreaker.NewCircuitBreaker(
        m.config.CircuitBreaker.FailureThreshold,
        m.config.CircuitBreaker.Timeout,
        m.config.CircuitBreaker.MaxRequests,
    )
    
    return func(c *gin.Context) {
        // Get circuit breaker key
        key := m.getCircuitBreakerKey(c)
        
        // Check circuit breaker state
        state := circuitBreaker.GetState(key)
        if state == circuitbreaker.StateOpen {
            m.metrics.IncrementCounter("circuit_breaker_open_total", map[string]string{
                "key": key,
            })
            
            c.JSON(http.StatusServiceUnavailable, gin.H{
                "error": "Service temporarily unavailable",
            })
            c.Abort()
            return
        }
        
        // Execute with circuit breaker
        result, err := circuitBreaker.Execute(c.Request.Context(), key, func() (interface{}, error) {
            // Process request
            c.Next()
            
            // Check if request was successful
            if c.Writer.Status() >= 500 {
                return nil, fmt.Errorf("server error: %d", c.Writer.Status())
            }
            
            return nil, nil
        })
        
        if err != nil {
            m.logger.Error("Circuit breaker execution failed", "error", err)
            
            // Record circuit breaker failure
            m.metrics.IncrementCounter("circuit_breaker_failure_total", map[string]string{
                "key": key,
            })
            
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Request failed",
            })
            c.Abort()
            return
        }
        
        // Record circuit breaker success
        m.metrics.IncrementCounter("circuit_breaker_success_total", map[string]string{
            "key": key,
        })
    }
}

func (m *MiddlewareManager) getCircuitBreakerKey(c *gin.Context) string {
    return c.Request.URL.Path
}
```

### 6. Caching Middleware

```go
// internal/middleware/caching.go
package middleware

import (
    "crypto/md5"
    "fmt"
    "io"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/anasamu/go-micro-libs/cache"
    "github.com/anasamu/go-micro-libs/monitoring"
)

func (m *MiddlewareManager) cachingMiddleware() gin.HandlerFunc {
    cacheManager := cache.NewCacheManager(m.config.Caching.Provider, m.config.Caching.TTL)
    
    return func(c *gin.Context) {
        // Only cache GET requests
        if c.Request.Method != "GET" {
            c.Next()
            return
        }
        
        // Generate cache key
        cacheKey := m.generateCacheKey(c)
        
        // Check cache
        cachedResponse, err := cacheManager.Get(c.Request.Context(), cacheKey)
        if err == nil && cachedResponse != nil {
            // Cache hit
            m.metrics.IncrementCounter("cache_hits_total", map[string]string{
                "key": cacheKey,
            })
            
            // Set cached response
            c.Header("X-Cache", "HIT")
            c.Header("X-Cache-Key", cacheKey)
            c.Data(200, "application/json", cachedResponse.([]byte))
            c.Abort()
            return
        }
        
        // Cache miss
        m.metrics.IncrementCounter("cache_misses_total", map[string]string{
            "key": cacheKey,
        })
        
        // Capture response
        var responseBody bytes.Buffer
        writer := &responseWriter{
            ResponseWriter: c.Writer,
            body:          &responseBody,
        }
        c.Writer = writer
        
        // Process request
        c.Next()
        
        // Cache response if successful
        if c.Writer.Status() == 200 {
            responseData := responseBody.Bytes()
            
            // Cache the response
            err = cacheManager.Set(c.Request.Context(), cacheKey, responseData)
            if err != nil {
                m.logger.Error("Failed to cache response", "error", err)
            }
            
            c.Header("X-Cache", "MISS")
            c.Header("X-Cache-Key", cacheKey)
        }
    }
}

func (m *MiddlewareManager) generateCacheKey(c *gin.Context) string {
    // Create hash of request
    h := md5.New()
    io.WriteString(h, c.Request.Method)
    io.WriteString(h, c.Request.URL.String())
    io.WriteString(h, c.GetHeader("Authorization"))
    
    return fmt.Sprintf("cache:%x", h.Sum(nil))
}
```

### 7. Compression Middleware

```go
// internal/middleware/compression.go
package middleware

import (
    "compress/gzip"
    "net/http"
    "strings"
    
    "github.com/gin-gonic/gin"
)

func (m *MiddlewareManager) compressionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check if client supports gzip
        if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
            c.Next()
            return
        }
        
        // Create gzip writer
        gz := gzip.NewWriter(c.Writer)
        defer gz.Close()
        
        // Set headers
        c.Header("Content-Encoding", "gzip")
        c.Header("Vary", "Accept-Encoding")
        
        // Wrap response writer
        c.Writer = &gzipResponseWriter{
            ResponseWriter: c.Writer,
            gz:            gz,
        }
        
        c.Next()
    }
}

type gzipResponseWriter struct {
    gin.ResponseWriter
    gz *gzip.Writer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
    return w.gz.Write(data)
}

func (w *gzipResponseWriter) WriteString(s string) (int, error) {
    return w.gz.Write([]byte(s))
}
```

### 8. Security Middleware

```go
// internal/middleware/security.go
package middleware

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
)

func (m *MiddlewareManager) corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Set CORS headers
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Header("Access-Control-Allow-Credentials", "true")
        
        // Handle preflight requests
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        
        c.Next()
    }
}

func (m *MiddlewareManager) helmetMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Set security headers
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        
        c.Next()
    }
}

func (m *MiddlewareManager) xssMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // XSS protection is handled by helmet middleware
        c.Next()
    }
}
```

### 9. Authentication Middleware

```go
// internal/middleware/auth.go
package middleware

import (
    "net/http"
    "strings"
    
    "github.com/gin-gonic/gin"
    "github.com/anasamu/go-micro-libs/auth"
    "github.com/anasamu/go-micro-libs/monitoring"
)

func (m *MiddlewareManager) authMiddleware() gin.HandlerFunc {
    jwtService := auth.NewJWTService(m.config.Auth.Secret, m.config.Auth.Expiration)
    
    return func(c *gin.Context) {
        // Get token from header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            m.metrics.IncrementCounter("auth_errors_total", map[string]string{
                "error": "missing_authorization_header",
            })
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        // Extract token
        tokenString := m.extractToken(authHeader)
        if tokenString == "" {
            m.metrics.IncrementCounter("auth_errors_total", map[string]string{
                "error": "invalid_authorization_header",
            })
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
            c.Abort()
            return
        }
        
        // Validate token
        claims, err := jwtService.ValidateToken(c.Request.Context(), tokenString)
        if err != nil {
            m.metrics.IncrementCounter("auth_errors_total", map[string]string{
                "error": "invalid_token",
            })
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        // Set user context
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("email", claims.Email)
        c.Set("roles", claims.Roles)
        
        // Record success
        m.metrics.IncrementCounter("auth_success_total", map[string]string{
            "user_id": claims.UserID,
        })
        
        c.Next()
    }
}

func (m *MiddlewareManager) extractToken(authHeader string) string {
    if !strings.HasPrefix(authHeader, "Bearer ") {
        return ""
    }
    return strings.TrimPrefix(authHeader, "Bearer ")
}
```

### 10. Authorization Middleware

```go
// internal/middleware/authorization.go
package middleware

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/anasamu/go-micro-libs/auth"
    "github.com/anasamu/go-micro-libs/monitoring"
)

func (m *MiddlewareManager) authorizationMiddleware() gin.HandlerFunc {
    rbacService := auth.NewRBACService(m.config.Authorization.Provider)
    
    return func(c *gin.Context) {
        // Get user ID from context
        userID, exists := c.Get("user_id")
        if !exists {
            m.metrics.IncrementCounter("authorization_errors_total", map[string]string{
                "error": "user_not_authenticated",
            })
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
            c.Abort()
            return
        }
        
        // Get required permission from route
        permission := c.GetString("required_permission")
        if permission == "" {
            // No permission required
            c.Next()
            return
        }
        
        // Check permission
        hasPermission, err := rbacService.HasPermission(c.Request.Context(), userID.(string), permission)
        if err != nil {
            m.metrics.IncrementCounter("authorization_errors_total", map[string]string{
                "error": "permission_check_failed",
            })
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permission"})
            c.Abort()
            return
        }
        
        if !hasPermission {
            m.metrics.IncrementCounter("authorization_errors_total", map[string]string{
                "error": "insufficient_permissions",
                "user_id": userID.(string),
            })
            c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }
        
        // Record success
        m.metrics.IncrementCounter("authorization_success_total", map[string]string{
            "permission": permission,
            "user_id":    userID.(string),
        })
        
        c.Next()
    }
}
```

## 🔧 Middleware Configuration

### 1. Configuration Types

```go
// internal/middleware/config.go
package middleware

type AuthConfig struct {
    Enabled    bool   `yaml:"enabled"`
    Provider   string `yaml:"provider"`
    Secret     string `yaml:"secret"`
    Expiration string `yaml:"expiration"`
}

type AuthorizationConfig struct {
    Enabled     bool   `yaml:"enabled"`
    Provider    string `yaml:"provider"`
    DefaultRole string `yaml:"default_role"`
}

type LoggingConfig struct {
    Enabled       bool   `yaml:"enabled"`
    Level         string `yaml:"level"`
    Format        string `yaml:"format"`
    IncludeBody   bool   `yaml:"include_body"`
    IncludeHeaders bool  `yaml:"include_headers"`
}

type MonitoringConfig struct {
    Enabled      bool `yaml:"enabled"`
    Metrics      bool `yaml:"metrics"`
    Tracing      bool `yaml:"tracing"`
    HealthChecks bool `yaml:"health_checks"`
}

type RateLimitConfig struct {
    Enabled            bool   `yaml:"enabled"`
    Provider           string `yaml:"provider"`
    RequestsPerMinute  int    `yaml:"requests_per_minute"`
    Burst              int    `yaml:"burst"`
    KeyStrategy        string `yaml:"key_strategy"`
}

type CircuitBreakerConfig struct {
    Enabled          bool          `yaml:"enabled"`
    FailureThreshold int           `yaml:"failure_threshold"`
    Timeout          time.Duration `yaml:"timeout"`
    MaxRequests      int           `yaml:"max_requests"`
}

type CachingConfig struct {
    Enabled     bool          `yaml:"enabled"`
    Provider    string        `yaml:"provider"`
    TTL         time.Duration `yaml:"ttl"`
    KeyStrategy string        `yaml:"key_strategy"`
}

type CompressionConfig struct {
    Enabled    bool   `yaml:"enabled"`
    Algorithm  string `yaml:"algorithm"`
    Level      int    `yaml:"level"`
}

type SecurityConfig struct {
    Enabled bool        `yaml:"enabled"`
    CORS    CORSConfig  `yaml:"cors"`
    CSRF    CSRFConfig  `yaml:"csrf"`
    XSS     XSSConfig   `yaml:"xss"`
    Helmet  HelmetConfig `yaml:"helmet"`
}

type CORSConfig struct {
    Enabled bool     `yaml:"enabled"`
    Origins []string `yaml:"origins"`
    Methods []string `yaml:"methods"`
    Headers []string `yaml:"headers"`
}

type CSRFConfig struct {
    Enabled bool `yaml:"enabled"`
}

type XSSConfig struct {
    Enabled bool `yaml:"enabled"`
}

type HelmetConfig struct {
    Enabled bool `yaml:"enabled"`
}
```

## 🔧 Best Practices

### 1. Middleware Order

```go
// Correct middleware order
func (m *MiddlewareManager) SetupMiddleware() []gin.HandlerFunc {
    // 1. Security first (CORS, Helmet, XSS)
    if m.config.Security.Enabled {
        m.addSecurityMiddleware()
    }
    
    // 2. Logging second
    if m.config.Logging.Enabled {
        m.addLoggingMiddleware()
    }
    
    // 3. Monitoring third
    if m.config.Monitoring.Enabled {
        m.addMonitoringMiddleware()
    }
    
    // 4. Rate limiting fourth
    if m.config.RateLimit.Enabled {
        m.addRateLimitMiddleware()
    }
    
    // 5. Circuit breaker fifth
    if m.config.CircuitBreaker.Enabled {
        m.addCircuitBreakerMiddleware()
    }
    
    // 6. Caching sixth
    if m.config.Caching.Enabled {
        m.addCachingMiddleware()
    }
    
    // 7. Compression seventh
    if m.config.Compression.Enabled {
        m.addCompressionMiddleware()
    }
    
    // 8. Authentication eighth
    if m.config.Auth.Enabled {
        m.addAuthMiddleware()
    }
    
    // 9. Authorization last
    if m.config.Authorization.Enabled {
        m.addAuthorizationMiddleware()
    }
    
    return m.middlewares
}
```

### 2. Error Handling

```go
// Proper error handling in middleware
func (m *MiddlewareManager) handleMiddlewareError(c *gin.Context, err error, message string) {
    m.logger.Error("Middleware error", "error", err, "message", message)
    
    m.metrics.IncrementCounter("middleware_errors_total", map[string]string{
        "error": err.Error(),
    })
    
    c.JSON(http.StatusInternalServerError, gin.H{
        "error": message,
    })
    c.Abort()
}
```

### 3. Performance Optimization

```go
// Optimize middleware performance
func (m *MiddlewareManager) optimizeMiddleware() {
    // Use connection pooling for database operations
    // Implement caching for frequently accessed data
    // Use async operations where possible
    // Minimize memory allocations
    // Use efficient data structures
}
```

### 4. Testing

```go
// Test middleware
func TestLoggingMiddleware(t *testing.T) {
    // Setup test
    gin.SetMode(gin.TestMode)
    router := gin.New()
    
    config := &Config{
        Logging: LoggingConfig{
            Enabled: true,
            Level:   "info",
            Format:  "json",
        },
    }
    
    manager := NewMiddlewareManager(config, logger, metrics)
    router.Use(manager.loggingMiddleware())
    
    // Test request
    req, _ := http.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    
    router.ServeHTTP(w, req)
    
    // Assertions
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Header().Get("X-Correlation-ID"), "")
}
```

---

**Middleware - Powerful and flexible middleware system for microservices! 🚀**

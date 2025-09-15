# Authorization Implementation

## 🎯 Overview

GoMicroFramework menyediakan sistem authorization yang komprehensif dengan dukungan untuk berbagai model seperti RBAC (Role-Based Access Control), ABAC (Attribute-Based Access Control), dan ACL (Access Control Lists). Framework ini memungkinkan implementasi authorization yang granular dan scalable untuk microservices.

## 🔧 Supported Authorization Models

### 1. RBAC (Role-Based Access Control)
- Role-based permissions
- Hierarchical roles
- Dynamic role assignment

### 2. ABAC (Attribute-Based Access Control)
- Attribute-based policies
- Context-aware decisions
- Fine-grained control

### 3. ACL (Access Control Lists)
- Resource-based permissions
- User-resource mappings
- Simple permission model

### 4. Policy-Based Authorization
- JSON-based policies
- Rule-based decisions
- External policy engines

## 🔧 Authorization Setup

### 1. Generate Service with Authorization

```bash
# Generate service with RBAC authorization
microframework new user-service --with-auth=jwt --with-database=postgres

# Generate service with ABAC authorization
microframework new resource-service --with-auth=jwt --with-database=postgres
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

# Authentication configuration
auth:
  providers:
    jwt:
      enabled: true
      secret: "${JWT_SECRET}"
      expiration: "24h"
      issuer: "user-service"
      audience: "api"

# Authorization configuration
authorization:
  providers:
    rbac:
      enabled: true
      hierarchical: true
      default_role: "user"
      
    abac:
      enabled: false
      policy_engine: "opa"
      policy_path: "./policies"
      
    acl:
      enabled: false
      cache_ttl: "1h"
      
    policy:
      enabled: false
      engine: "json"
      policies_path: "./policies"

# Database for authorization data
database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100
```

## 🔧 RBAC Implementation

### 1. RBAC Models

```go
// internal/models/rbac.go
package models

import (
    "time"
)

type Role struct {
    ID          string    `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Description string    `json:"description" db:"description"`
    Permissions []string  `json:"permissions" db:"permissions"`
    ParentID    *string   `json:"parent_id" db:"parent_id"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Permission struct {
    ID          string    `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Description string    `json:"description" db:"description"`
    Resource    string    `json:"resource" db:"resource"`
    Action      string    `json:"action" db:"action"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type UserRole struct {
    ID        string    `json:"id" db:"id"`
    UserID    string    `json:"user_id" db:"user_id"`
    RoleID    string    `json:"role_id" db:"role_id"`
    AssignedBy string   `json:"assigned_by" db:"assigned_by"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    ExpiresAt *time.Time `json:"expires_at" db:"expires_at"`
}

type RolePermission struct {
    ID           string    `json:"id" db:"id"`
    RoleID       string    `json:"role_id" db:"role_id"`
    PermissionID string    `json:"permission_id" db:"permission_id"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
```

### 2. RBAC Service

```go
// internal/auth/rbac_service.go
package auth

import (
    "context"
    "fmt"
    "time"
    
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/cache"
)

type RBACService struct {
    dbManager    *database.DatabaseManager
    cacheManager *cache.CacheManager
}

func NewRBACService(dbManager *database.DatabaseManager, cacheManager *cache.CacheManager) *RBACService {
    return &RBACService{
        dbManager:    dbManager,
        cacheManager: cacheManager,
    }
}

func (s *RBACService) CreateRole(ctx context.Context, role *Role) error {
    query := `INSERT INTO roles (id, name, description, permissions, parent_id, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
    
    now := time.Now()
    _, err := s.dbManager.Exec(ctx, "postgresql", query, 
        role.ID, role.Name, role.Description, role.Permissions, role.ParentID, now, now)
    
    if err != nil {
        return err
    }
    
    // Invalidate cache
    s.cacheManager.Delete(ctx, "roles")
    
    return nil
}

func (s *RBACService) GetRole(ctx context.Context, roleID string) (*Role, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("role:%s", roleID)
    var role Role
    err := s.cacheManager.Get(ctx, cacheKey, &role)
    if err == nil {
        return &role, nil
    }
    
    // Get from database
    query := `SELECT id, name, description, permissions, parent_id, created_at, updated_at 
              FROM roles WHERE id = $1`
    
    result, err := s.dbManager.Query(ctx, "postgresql", query, roleID)
    if err != nil {
        return nil, err
    }
    
    if result.Next() {
        err = result.Scan(&role.ID, &role.Name, &role.Description, 
            &role.Permissions, &role.ParentID, &role.CreatedAt, &role.UpdatedAt)
        if err != nil {
            return nil, err
        }
    } else {
        return nil, fmt.Errorf("role not found")
    }
    
    // Cache the result
    s.cacheManager.Set(ctx, cacheKey, role, 1*time.Hour)
    
    return &role, nil
}

func (s *RBACService) AssignRoleToUser(ctx context.Context, userID, roleID, assignedBy string, expiresAt *time.Time) error {
    query := `INSERT INTO user_roles (id, user_id, role_id, assigned_by, created_at, expires_at) 
              VALUES ($1, $2, $3, $4, $5, $6)`
    
    now := time.Now()
    _, err := s.dbManager.Exec(ctx, "postgresql", query, 
        generateID(), userID, roleID, assignedBy, now, expiresAt)
    
    if err != nil {
        return err
    }
    
    // Invalidate user cache
    s.cacheManager.Delete(ctx, fmt.Sprintf("user_roles:%s", userID))
    
    return nil
}

func (s *RBACService) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
    query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
    
    _, err := s.dbManager.Exec(ctx, "postgresql", query, userID, roleID)
    if err != nil {
        return err
    }
    
    // Invalidate user cache
    s.cacheManager.Delete(ctx, fmt.Sprintf("user_roles:%s", userID))
    
    return nil
}

func (s *RBACService) GetUserRoles(ctx context.Context, userID string) ([]*Role, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("user_roles:%s", userID)
    var roles []*Role
    err := s.cacheManager.Get(ctx, cacheKey, &roles)
    if err == nil {
        return roles, nil
    }
    
    // Get from database
    query := `SELECT r.id, r.name, r.description, r.permissions, r.parent_id, r.created_at, r.updated_at
              FROM roles r
              JOIN user_roles ur ON r.id = ur.role_id
              WHERE ur.user_id = $1 AND (ur.expires_at IS NULL OR ur.expires_at > NOW())`
    
    result, err := s.dbManager.Query(ctx, "postgresql", query, userID)
    if err != nil {
        return nil, err
    }
    
    for result.Next() {
        role := &Role{}
        err = result.Scan(&role.ID, &role.Name, &role.Description, 
            &role.Permissions, &role.ParentID, &role.CreatedAt, &role.UpdatedAt)
        if err != nil {
            return nil, err
        }
        roles = append(roles, role)
    }
    
    // Cache the result
    s.cacheManager.Set(ctx, cacheKey, roles, 1*time.Hour)
    
    return roles, nil
}

func (s *RBACService) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("user_permissions:%s", userID)
    var permissions []string
    err := s.cacheManager.Get(ctx, cacheKey, &permissions)
    if err == nil {
        return permissions, nil
    }
    
    // Get user roles
    roles, err := s.GetUserRoles(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // Collect permissions from roles
    permissionSet := make(map[string]bool)
    for _, role := range roles {
        for _, permission := range role.Permissions {
            permissionSet[permission] = true
        }
        
        // Get inherited permissions from parent roles
        inheritedPermissions, err := s.getInheritedPermissions(ctx, role)
        if err != nil {
            return nil, err
        }
        
        for _, permission := range inheritedPermissions {
            permissionSet[permission] = true
        }
    }
    
    // Convert to slice
    for permission := range permissionSet {
        permissions = append(permissions, permission)
    }
    
    // Cache the result
    s.cacheManager.Set(ctx, cacheKey, permissions, 1*time.Hour)
    
    return permissions, nil
}

func (s *RBACService) HasPermission(ctx context.Context, userID, permission string) (bool, error) {
    permissions, err := s.GetUserPermissions(ctx, userID)
    if err != nil {
        return false, err
    }
    
    for _, p := range permissions {
        if p == permission {
            return true, nil
        }
    }
    
    return false, nil
}

func (s *RBACService) HasRole(ctx context.Context, userID, roleName string) (bool, error) {
    roles, err := s.GetUserRoles(ctx, userID)
    if err != nil {
        return false, err
    }
    
    for _, role := range roles {
        if role.Name == roleName {
            return true, nil
        }
    }
    
    return false, nil
}

func (s *RBACService) getInheritedPermissions(ctx context.Context, role *Role) ([]string, error) {
    if role.ParentID == nil {
        return []string{}, nil
    }
    
    parentRole, err := s.GetRole(ctx, *role.ParentID)
    if err != nil {
        return nil, err
    }
    
    permissions := make([]string, len(parentRole.Permissions))
    copy(permissions, parentRole.Permissions)
    
    // Get permissions from grandparent roles
    inheritedPermissions, err := s.getInheritedPermissions(ctx, parentRole)
    if err != nil {
        return nil, err
    }
    
    permissions = append(permissions, inheritedPermissions...)
    
    return permissions, nil
}
```

### 3. RBAC Middleware

```go
// internal/middleware/rbac_middleware.go
package middleware

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/anasamu/go-micro-libs/auth"
    "github.com/anasamu/go-micro-libs/monitoring"
)

type RBACMiddleware struct {
    rbacService *auth.RBACService
}

func NewRBACMiddleware(rbacService *auth.RBACService) *RBACMiddleware {
    return &RBACMiddleware{
        rbacService: rbacService,
    }
}

func (m *RBACMiddleware) RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Record metrics
        monitoring.IncrementCounter("authorization_requests_total", map[string]string{
            "permission": permission,
            "endpoint":   c.Request.URL.Path,
        })
        
        // Get user ID from context
        userID, exists := c.Get("user_id")
        if !exists {
            monitoring.IncrementCounter("authorization_errors_total", map[string]string{
                "error": "user_not_authenticated",
            })
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
            c.Abort()
            return
        }
        
        // Check permission
        hasPermission, err := m.rbacService.HasPermission(c.Request.Context(), userID.(string), permission)
        if err != nil {
            monitoring.IncrementCounter("authorization_errors_total", map[string]string{
                "error": "permission_check_failed",
            })
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permission"})
            c.Abort()
            return
        }
        
        if !hasPermission {
            monitoring.IncrementCounter("authorization_errors_total", map[string]string{
                "error": "insufficient_permissions",
                "user_id": userID.(string),
            })
            c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }
        
        // Record success
        monitoring.IncrementCounter("authorization_success_total", map[string]string{
            "permission": permission,
            "user_id":    userID.(string),
        })
        
        c.Next()
    }
}

func (m *RBACMiddleware) RequireRole(roleName string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Record metrics
        monitoring.IncrementCounter("authorization_requests_total", map[string]string{
            "role":     roleName,
            "endpoint": c.Request.URL.Path,
        })
        
        // Get user ID from context
        userID, exists := c.Get("user_id")
        if !exists {
            monitoring.IncrementCounter("authorization_errors_total", map[string]string{
                "error": "user_not_authenticated",
            })
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
            c.Abort()
            return
        }
        
        // Check role
        hasRole, err := m.rbacService.HasRole(c.Request.Context(), userID.(string), roleName)
        if err != nil {
            monitoring.IncrementCounter("authorization_errors_total", map[string]string{
                "error": "role_check_failed",
            })
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check role"})
            c.Abort()
            return
        }
        
        if !hasRole {
            monitoring.IncrementCounter("authorization_errors_total", map[string]string{
                "error": "insufficient_role",
                "user_id": userID.(string),
            })
            c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient role"})
            c.Abort()
            return
        }
        
        // Record success
        monitoring.IncrementCounter("authorization_success_total", map[string]string{
            "role":    roleName,
            "user_id": userID.(string),
        })
        
        c.Next()
    }
}

func (m *RBACMiddleware) RequireAnyRole(roleNames ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get user ID from context
        userID, exists := c.Get("user_id")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
            c.Abort()
            return
        }
        
        // Check if user has any of the required roles
        for _, roleName := range roleNames {
            hasRole, err := m.rbacService.HasRole(c.Request.Context(), userID.(string), roleName)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check role"})
                c.Abort()
                return
            }
            
            if hasRole {
                c.Next()
                return
            }
        }
        
        c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient role"})
        c.Abort()
    }
}
```

## 🔧 ABAC Implementation

### 1. ABAC Models

```go
// internal/models/abac.go
package models

import (
    "time"
)

type Policy struct {
    ID          string    `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Description string    `json:"description" db:"description"`
    Rules       []Rule    `json:"rules" db:"rules"`
    Enabled     bool      `json:"enabled" db:"enabled"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Rule struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Effect      string                 `json:"effect"` // "allow" or "deny"
    Conditions  map[string]interface{} `json:"conditions"`
    Actions     []string               `json:"actions"`
    Resources   []string               `json:"resources"`
}

type Attribute struct {
    Name  string      `json:"name"`
    Value interface{} `json:"value"`
    Type  string      `json:"type"` // "string", "number", "boolean", "array"
}

type Context struct {
    User        map[string]interface{} `json:"user"`
    Resource    map[string]interface{} `json:"resource"`
    Environment map[string]interface{} `json:"environment"`
    Request     map[string]interface{} `json:"request"`
}
```

### 2. ABAC Service

```go
// internal/auth/abac_service.go
package auth

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/cache"
)

type ABACService struct {
    dbManager    *database.DatabaseManager
    cacheManager *cache.CacheManager
}

func NewABACService(dbManager *database.DatabaseManager, cacheManager *cache.CacheManager) *ABACService {
    return &ABACService{
        dbManager:    dbManager,
        cacheManager: cacheManager,
    }
}

func (s *ABACService) CreatePolicy(ctx context.Context, policy *Policy) error {
    // Serialize rules
    rulesJSON, err := json.Marshal(policy.Rules)
    if err != nil {
        return err
    }
    
    query := `INSERT INTO policies (id, name, description, rules, enabled, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
    
    now := time.Now()
    _, err = s.dbManager.Exec(ctx, "postgresql", query, 
        policy.ID, policy.Name, policy.Description, rulesJSON, policy.Enabled, now, now)
    
    if err != nil {
        return err
    }
    
    // Invalidate cache
    s.cacheManager.Delete(ctx, "policies")
    
    return nil
}

func (s *ABACService) GetPolicy(ctx context.Context, policyID string) (*Policy, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("policy:%s", policyID)
    var policy Policy
    err := s.cacheManager.Get(ctx, cacheKey, &policy)
    if err == nil {
        return &policy, nil
    }
    
    // Get from database
    query := `SELECT id, name, description, rules, enabled, created_at, updated_at 
              FROM policies WHERE id = $1`
    
    result, err := s.dbManager.Query(ctx, "postgresql", query, policyID)
    if err != nil {
        return nil, err
    }
    
    var rulesJSON []byte
    if result.Next() {
        err = result.Scan(&policy.ID, &policy.Name, &policy.Description, 
            &rulesJSON, &policy.Enabled, &policy.CreatedAt, &policy.UpdatedAt)
        if err != nil {
            return nil, err
        }
    } else {
        return nil, fmt.Errorf("policy not found")
    }
    
    // Deserialize rules
    err = json.Unmarshal(rulesJSON, &policy.Rules)
    if err != nil {
        return nil, err
    }
    
    // Cache the result
    s.cacheManager.Set(ctx, cacheKey, policy, 1*time.Hour)
    
    return &policy, nil
}

func (s *ABACService) EvaluatePolicy(ctx context.Context, policyID string, context *Context) (bool, error) {
    // Get policy
    policy, err := s.GetPolicy(ctx, policyID)
    if err != nil {
        return false, err
    }
    
    if !policy.Enabled {
        return false, nil
    }
    
    // Evaluate rules
    for _, rule := range policy.Rules {
        if s.evaluateRule(rule, context) {
            return rule.Effect == "allow", nil
        }
    }
    
    return false, nil
}

func (s *ABACService) EvaluateAccess(ctx context.Context, userID, resource, action string, context *Context) (bool, error) {
    // Get all enabled policies
    policies, err := s.GetEnabledPolicies(ctx)
    if err != nil {
        return false, err
    }
    
    // Set default context
    if context == nil {
        context = &Context{}
    }
    
    // Set user context
    if context.User == nil {
        context.User = make(map[string]interface{})
    }
    context.User["id"] = userID
    
    // Set resource context
    if context.Resource == nil {
        context.Resource = make(map[string]interface{})
    }
    context.Resource["name"] = resource
    
    // Set request context
    if context.Request == nil {
        context.Request = make(map[string]interface{})
    }
    context.Request["action"] = action
    
    // Evaluate policies
    for _, policy := range policies {
        allowed, err := s.EvaluatePolicy(ctx, policy.ID, context)
        if err != nil {
            return false, err
        }
        
        if allowed {
            return true, nil
        }
    }
    
    return false, nil
}

func (s *ABACService) GetEnabledPolicies(ctx context.Context) ([]*Policy, error) {
    // Check cache first
    var policies []*Policy
    err := s.cacheManager.Get(ctx, "enabled_policies", &policies)
    if err == nil {
        return policies, nil
    }
    
    // Get from database
    query := `SELECT id, name, description, rules, enabled, created_at, updated_at 
              FROM policies WHERE enabled = true`
    
    result, err := s.dbManager.Query(ctx, "postgresql", query)
    if err != nil {
        return nil, err
    }
    
    for result.Next() {
        policy := &Policy{}
        var rulesJSON []byte
        err = result.Scan(&policy.ID, &policy.Name, &policy.Description, 
            &rulesJSON, &policy.Enabled, &policy.CreatedAt, &policy.UpdatedAt)
        if err != nil {
            return nil, err
        }
        
        // Deserialize rules
        err = json.Unmarshal(rulesJSON, &policy.Rules)
        if err != nil {
            return nil, err
        }
        
        policies = append(policies, policy)
    }
    
    // Cache the result
    s.cacheManager.Set(ctx, "enabled_policies", policies, 1*time.Hour)
    
    return policies, nil
}

func (s *ABACService) evaluateRule(rule Rule, context *Context) bool {
    // Check if rule applies to the request
    if !s.ruleAppliesToRequest(rule, context) {
        return false
    }
    
    // Evaluate conditions
    for condition, expectedValue := range rule.Conditions {
        if !s.evaluateCondition(condition, expectedValue, context) {
            return false
        }
    }
    
    return true
}

func (s *ABACService) ruleAppliesToRequest(rule Rule, context *Context) bool {
    // Check actions
    if len(rule.Actions) > 0 {
        requestAction, exists := context.Request["action"]
        if !exists {
            return false
        }
        
        actionMatched := false
        for _, action := range rule.Actions {
            if action == requestAction {
                actionMatched = true
                break
            }
        }
        
        if !actionMatched {
            return false
        }
    }
    
    // Check resources
    if len(rule.Resources) > 0 {
        resourceName, exists := context.Resource["name"]
        if !exists {
            return false
        }
        
        resourceMatched := false
        for _, resource := range rule.Resources {
            if resource == resourceName {
                resourceMatched = true
                break
            }
        }
        
        if !resourceMatched {
            return false
        }
    }
    
    return true
}

func (s *ABACService) evaluateCondition(condition string, expectedValue interface{}, context *Context) bool {
    // Parse condition (e.g., "user.role", "resource.owner", "environment.time")
    parts := strings.Split(condition, ".")
    if len(parts) != 2 {
        return false
    }
    
    contextType := parts[0]
    attribute := parts[1]
    
    var actualValue interface{}
    var exists bool
    
    switch contextType {
    case "user":
        actualValue, exists = context.User[attribute]
    case "resource":
        actualValue, exists = context.Resource[attribute]
    case "environment":
        actualValue, exists = context.Environment[attribute]
    case "request":
        actualValue, exists = context.Request[attribute]
    default:
        return false
    }
    
    if !exists {
        return false
    }
    
    // Compare values
    return s.compareValues(actualValue, expectedValue)
}

func (s *ABACService) compareValues(actual, expected interface{}) bool {
    // Simple equality comparison
    // In a real implementation, you might want to support more complex comparisons
    return actual == expected
}
```

### 3. ABAC Middleware

```go
// internal/middleware/abac_middleware.go
package middleware

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/anasamu/go-micro-libs/auth"
    "github.com/anasamu/go-micro-libs/monitoring"
)

type ABACMiddleware struct {
    abacService *auth.ABACService
}

func NewABACMiddleware(abacService *auth.ABACService) *ABACMiddleware {
    return &ABACMiddleware{
        abacService: abacService,
    }
}

func (m *ABACMiddleware) RequireAccess(resource, action string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Record metrics
        monitoring.IncrementCounter("abac_requests_total", map[string]string{
            "resource": resource,
            "action":   action,
            "endpoint": c.Request.URL.Path,
        })
        
        // Get user ID from context
        userID, exists := c.Get("user_id")
        if !exists {
            monitoring.IncrementCounter("abac_errors_total", map[string]string{
                "error": "user_not_authenticated",
            })
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
            c.Abort()
            return
        }
        
        // Create context for ABAC evaluation
        context := &auth.Context{
            User: map[string]interface{}{
                "id":   userID,
                "ip":   c.ClientIP(),
                "role": c.GetString("role"),
            },
            Resource: map[string]interface{}{
                "name": resource,
                "id":   c.Param("id"),
            },
            Environment: map[string]interface{}{
                "time": time.Now().Format("15:04"),
                "date": time.Now().Format("2006-01-02"),
            },
            Request: map[string]interface{}{
                "action": action,
                "method": c.Request.Method,
                "path":   c.Request.URL.Path,
            },
        }
        
        // Evaluate access
        allowed, err := m.abacService.EvaluateAccess(c.Request.Context(), userID.(string), resource, action, context)
        if err != nil {
            monitoring.IncrementCounter("abac_errors_total", map[string]string{
                "error": "evaluation_failed",
            })
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to evaluate access"})
            c.Abort()
            return
        }
        
        if !allowed {
            monitoring.IncrementCounter("abac_errors_total", map[string]string{
                "error": "access_denied",
                "user_id": userID.(string),
            })
            c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
            c.Abort()
            return
        }
        
        // Record success
        monitoring.IncrementCounter("abac_success_total", map[string]string{
            "resource": resource,
            "action":   action,
            "user_id":  userID.(string),
        })
        
        c.Next()
    }
}
```

## 🔧 ACL Implementation

### 1. ACL Models

```go
// internal/models/acl.go
package models

import (
    "time"
)

type ACL struct {
    ID        string    `json:"id" db:"id"`
    UserID    string    `json:"user_id" db:"user_id"`
    Resource  string    `json:"resource" db:"resource"`
    Action    string    `json:"action" db:"action"`
    Granted   bool      `json:"granted" db:"granted"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ACLGroup struct {
    ID        string    `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Resources []string  `json:"resources" db:"resources"`
    Actions   []string  `json:"actions" db:"actions"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserACLGroup struct {
    ID        string    `json:"id" db:"id"`
    UserID    string    `json:"user_id" db:"user_id"`
    GroupID   string    `json:"group_id" db:"group_id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```

### 2. ACL Service

```go
// internal/auth/acl_service.go
package auth

import (
    "context"
    "fmt"
    "time"
    
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/cache"
)

type ACLService struct {
    dbManager    *database.DatabaseManager
    cacheManager *cache.CacheManager
}

func NewACLService(dbManager *database.DatabaseManager, cacheManager *cache.CacheManager) *ACLService {
    return &ACLService{
        dbManager:    dbManager,
        cacheManager: cacheManager,
    }
}

func (s *ACLService) GrantAccess(ctx context.Context, userID, resource, action string) error {
    query := `INSERT INTO acls (id, user_id, resource, action, granted, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) 
              ON CONFLICT (user_id, resource, action) 
              DO UPDATE SET granted = $5, updated_at = $7`
    
    now := time.Now()
    _, err := s.dbManager.Exec(ctx, "postgresql", query, 
        generateID(), userID, resource, action, true, now, now)
    
    if err != nil {
        return err
    }
    
    // Invalidate cache
    s.cacheManager.Delete(ctx, fmt.Sprintf("user_acls:%s", userID))
    
    return nil
}

func (s *ACLService) RevokeAccess(ctx context.Context, userID, resource, action string) error {
    query := `INSERT INTO acls (id, user_id, resource, action, granted, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) 
              ON CONFLICT (user_id, resource, action) 
              DO UPDATE SET granted = $5, updated_at = $7`
    
    now := time.Now()
    _, err := s.dbManager.Exec(ctx, "postgresql", query, 
        generateID(), userID, resource, action, false, now, now)
    
    if err != nil {
        return err
    }
    
    // Invalidate cache
    s.cacheManager.Delete(ctx, fmt.Sprintf("user_acls:%s", userID))
    
    return nil
}

func (s *ACLService) HasAccess(ctx context.Context, userID, resource, action string) (bool, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("user_acls:%s", userID)
    var acls []ACL
    err := s.cacheManager.Get(ctx, cacheKey, &acls)
    if err != nil {
        // Get from database
        acls, err = s.getUserACLs(ctx, userID)
        if err != nil {
            return false, err
        }
        
        // Cache the result
        s.cacheManager.Set(ctx, cacheKey, acls, 1*time.Hour)
    }
    
    // Check for specific permission
    for _, acl := range acls {
        if acl.Resource == resource && acl.Action == action {
            return acl.Granted, nil
        }
    }
    
    // Check for wildcard permissions
    for _, acl := range acls {
        if acl.Resource == "*" && acl.Action == action {
            return acl.Granted, nil
        }
        if acl.Resource == resource && acl.Action == "*" {
            return acl.Granted, nil
        }
        if acl.Resource == "*" && acl.Action == "*" {
            return acl.Granted, nil
        }
    }
    
    return false, nil
}

func (s *ACLService) getUserACLs(ctx context.Context, userID string) ([]ACL, error) {
    query := `SELECT id, user_id, resource, action, granted, created_at, updated_at 
              FROM acls WHERE user_id = $1`
    
    result, err := s.dbManager.Query(ctx, "postgresql", query, userID)
    if err != nil {
        return nil, err
    }
    
    var acls []ACL
    for result.Next() {
        acl := ACL{}
        err = result.Scan(&acl.ID, &acl.UserID, &acl.Resource, &acl.Action, 
            &acl.Granted, &acl.CreatedAt, &acl.UpdatedAt)
        if err != nil {
            return nil, err
        }
        acls = append(acls, acl)
    }
    
    return acls, nil
}
```

## 🔧 Best Practices

### 1. Permission Caching

```go
// Cache permissions for better performance
func (s *RBACService) GetUserPermissionsWithCache(ctx context.Context, userID string) ([]string, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("user_permissions:%s", userID)
    var permissions []string
    err := s.cacheManager.Get(ctx, cacheKey, &permissions)
    if err == nil {
        return permissions, nil
    }
    
    // Get from database
    permissions, err = s.GetUserPermissions(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // Cache with TTL
    s.cacheManager.Set(ctx, cacheKey, permissions, 1*time.Hour)
    
    return permissions, nil
}
```

### 2. Permission Validation

```go
// Validate permissions before granting
func (s *RBACService) ValidatePermission(permission string) error {
    // Check if permission exists
    query := `SELECT COUNT(*) FROM permissions WHERE name = $1`
    
    result, err := s.dbManager.Query(ctx, "postgresql", query, permission)
    if err != nil {
        return err
    }
    
    var count int
    if result.Next() {
        result.Scan(&count)
    }
    
    if count == 0 {
        return fmt.Errorf("permission %s does not exist", permission)
    }
    
    return nil
}
```

### 3. Audit Logging

```go
// Audit logging for authorization events
func (s *RBACService) AssignRoleWithAudit(ctx context.Context, userID, roleID, assignedBy string) error {
    err := s.AssignRoleToUser(ctx, userID, roleID, assignedBy, nil)
    if err != nil {
        return err
    }
    
    // Log audit event
    s.auditLogger.LogAuthEvent(ctx, &AuditEvent{
        UserID:    userID,
        Action:    "role_assigned",
        Details:   fmt.Sprintf("Role %s assigned by %s", roleID, assignedBy),
        Timestamp: time.Now(),
    })
    
    return nil
}
```

### 4. Permission Inheritance

```go
// Implement permission inheritance
func (s *RBACService) GetInheritedPermissions(ctx context.Context, roleID string) ([]string, error) {
    role, err := s.GetRole(ctx, roleID)
    if err != nil {
        return nil, err
    }
    
    permissions := make([]string, len(role.Permissions))
    copy(permissions, role.Permissions)
    
    // Get permissions from parent roles
    if role.ParentID != nil {
        inheritedPermissions, err := s.GetInheritedPermissions(ctx, *role.ParentID)
        if err != nil {
            return nil, err
        }
        
        permissions = append(permissions, inheritedPermissions...)
    }
    
    return permissions, nil
}
```

---

**Authorization - Comprehensive and flexible authorization for microservices! 🚀**

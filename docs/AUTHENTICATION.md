# Authentication Implementation

## 🎯 Overview

GoMicroFramework menyediakan sistem authentication yang komprehensif dengan dukungan untuk berbagai provider seperti JWT, OAuth2, LDAP, SAML, dan 2FA. Framework ini memungkinkan implementasi authentication yang aman dan scalable untuk microservices.

## 🔧 Supported Authentication Providers

### 1. JWT (JSON Web Tokens)
- Stateless authentication
- Self-contained tokens
- Ideal for microservices

### 2. OAuth2
- Industry standard for authorization
- Support for multiple providers (Google, GitHub, etc.)
- Secure token exchange

### 3. LDAP
- Enterprise directory integration
- Active Directory support
- Group-based authentication

### 4. SAML
- Enterprise SSO
- Identity provider integration
- XML-based assertions

### 5. 2FA (Two-Factor Authentication)
- TOTP (Time-based One-Time Password)
- SMS-based verification
- Email-based verification

## 🔧 Authentication Setup

### 1. Generate Service with Authentication

```bash
# Generate service with JWT authentication
microframework new user-service --with-auth=jwt --with-database=postgres

# Generate service with OAuth2 authentication
microframework new auth-service --with-auth=oauth --with-database=postgres

# Generate service with LDAP authentication
microframework new enterprise-service --with-auth=ldap --with-database=postgres
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
      algorithm: "HS256"
      refresh_token_expiration: "7d"
      
    oauth:
      enabled: false
      client_id: "${OAUTH_CLIENT_ID}"
      client_secret: "${OAUTH_CLIENT_SECRET}"
      redirect_url: "${OAUTH_REDIRECT_URL}"
      scopes: ["read", "write"]
      provider: "google"
      auth_url: "https://accounts.google.com/o/oauth2/auth"
      token_url: "https://oauth2.googleapis.com/token"
      user_info_url: "https://www.googleapis.com/oauth2/v2/userinfo"
      
    ldap:
      enabled: false
      server: "ldap://localhost:389"
      base_dn: "dc=example,dc=com"
      bind_dn: "cn=admin,dc=example,dc=com"
      bind_password: "${LDAP_PASSWORD}"
      user_search_base: "ou=users,dc=example,dc=com"
      user_search_filter: "(uid=%s)"
      group_search_base: "ou=groups,dc=example,dc=com"
      group_search_filter: "(member=%s)"
      
    saml:
      enabled: false
      entity_id: "user-service"
      sso_url: "https://sso.example.com/saml/sso"
      slo_url: "https://sso.example.com/saml/slo"
      certificate: "${SAML_CERTIFICATE}"
      private_key: "${SAML_PRIVATE_KEY}"
      
    two_factor:
      enabled: false
      provider: "totp"
      issuer: "user-service"
      algorithm: "SHA1"
      digits: 6
      period: 30
      sms:
        provider: "twilio"
        account_sid: "${TWILIO_ACCOUNT_SID}"
        auth_token: "${TWILIO_AUTH_TOKEN}"
        from_number: "${TWILIO_FROM_NUMBER}"
      email:
        provider: "smtp"
        smtp_host: "${SMTP_HOST}"
        smtp_port: 587
        smtp_username: "${SMTP_USERNAME}"
        smtp_password: "${SMTP_PASSWORD}"
        from_email: "${FROM_EMAIL}"

# Database for user storage
database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100
```

## 🔧 JWT Authentication Implementation

### 1. JWT Service

```go
// internal/auth/jwt_service.go
package auth

import (
    "context"
    "errors"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "github.com/anasamu/go-micro-libs/auth"
    "github.com/anasamu/go-micro-libs/auth/providers/jwt"
    "github.com/anasamu/go-micro-libs/database"
)

type JWTService struct {
    jwtProvider *jwt.Provider
    dbManager   *database.DatabaseManager
    secret      string
    expiration  time.Duration
}

type Claims struct {
    UserID    string   `json:"user_id"`
    Username  string   `json:"username"`
    Email     string   `json:"email"`
    Roles     []string `json:"roles"`
    jwt.RegisteredClaims
}

func NewJWTService(jwtProvider *jwt.Provider, dbManager *database.DatabaseManager, secret string, expiration time.Duration) *JWTService {
    return &JWTService{
        jwtProvider: jwtProvider,
        dbManager:   dbManager,
        secret:      secret,
        expiration:  expiration,
    }
}

func (s *JWTService) GenerateToken(ctx context.Context, user *User) (string, error) {
    // Create claims
    claims := &Claims{
        UserID:   user.ID,
        Username: user.Username,
        Email:    user.Email,
        Roles:    user.Roles,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "user-service",
            Audience:  []string{"api"},
            Subject:   user.ID,
        },
    }
    
    // Create token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    // Sign token
    tokenString, err := token.SignedString([]byte(s.secret))
    if err != nil {
        return "", err
    }
    
    // Store token in database (optional)
    err = s.storeToken(ctx, user.ID, tokenString)
    if err != nil {
        log.Printf("Failed to store token: %v", err)
    }
    
    return tokenString, nil
}

func (s *JWTService) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
    // Parse token
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(s.secret), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    // Validate token
    if !token.Valid {
        return nil, errors.New("invalid token")
    }
    
    // Extract claims
    claims, ok := token.Claims.(*Claims)
    if !ok {
        return nil, errors.New("invalid claims")
    }
    
    // Check if token is blacklisted
    if s.isTokenBlacklisted(ctx, tokenString) {
        return nil, errors.New("token is blacklisted")
    }
    
    return claims, nil
}

func (s *JWTService) RefreshToken(ctx context.Context, tokenString string) (string, error) {
    // Validate current token
    claims, err := s.ValidateToken(ctx, tokenString)
    if err != nil {
        return "", err
    }
    
    // Get user from database
    user, err := s.getUserByID(ctx, claims.UserID)
    if err != nil {
        return "", err
    }
    
    // Generate new token
    return s.GenerateToken(ctx, user)
}

func (s *JWTService) RevokeToken(ctx context.Context, tokenString string) error {
    // Add token to blacklist
    return s.blacklistToken(ctx, tokenString)
}

func (s *JWTService) storeToken(ctx context.Context, userID, tokenString string) error {
    query := `INSERT INTO user_tokens (user_id, token, created_at, expires_at) 
              VALUES ($1, $2, $3, $4)`
    
    now := time.Now()
    expiresAt := now.Add(s.expiration)
    
    _, err := s.dbManager.Exec(ctx, "postgresql", query, userID, tokenString, now, expiresAt)
    return err
}

func (s *JWTService) isTokenBlacklisted(ctx context.Context, tokenString string) bool {
    query := `SELECT COUNT(*) FROM blacklisted_tokens WHERE token = $1`
    
    result, err := s.dbManager.Query(ctx, "postgresql", query, tokenString)
    if err != nil {
        return false
    }
    
    var count int
    if result.Next() {
        result.Scan(&count)
    }
    
    return count > 0
}

func (s *JWTService) blacklistToken(ctx context.Context, tokenString string) error {
    query := `INSERT INTO blacklisted_tokens (token, created_at) VALUES ($1, $2)`
    
    _, err := s.dbManager.Exec(ctx, "postgresql", query, tokenString, time.Now())
    return err
}

func (s *JWTService) getUserByID(ctx context.Context, userID string) (*User, error) {
    query := `SELECT id, username, email, roles FROM users WHERE id = $1`
    
    result, err := s.dbManager.Query(ctx, "postgresql", query, userID)
    if err != nil {
        return nil, err
    }
    
    var user User
    if result.Next() {
        err = result.Scan(&user.ID, &user.Username, &user.Email, &user.Roles)
        if err != nil {
            return nil, err
        }
    }
    
    return &user, nil
}
```

### 2. JWT Middleware

```go
// internal/middleware/jwt_middleware.go
package middleware

import (
    "context"
    "net/http"
    "strings"
    
    "github.com/gin-gonic/gin"
    "github.com/anasamu/go-micro-libs/auth"
    "github.com/anasamu/go-micro-libs/monitoring"
)

type JWTMiddleware struct {
    jwtService *auth.JWTService
}

func NewJWTMiddleware(jwtService *auth.JWTService) *JWTMiddleware {
    return &JWTMiddleware{
        jwtService: jwtService,
    }
}

func (m *JWTMiddleware) RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Record metrics
        monitoring.IncrementCounter("auth_requests_total", map[string]string{
            "endpoint": c.Request.URL.Path,
        })
        
        // Get token from header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            monitoring.IncrementCounter("auth_errors_total", map[string]string{
                "error": "missing_authorization_header",
            })
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        // Extract token
        tokenString := m.extractToken(authHeader)
        if tokenString == "" {
            monitoring.IncrementCounter("auth_errors_total", map[string]string{
                "error": "invalid_authorization_header",
            })
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
            c.Abort()
            return
        }
        
        // Validate token
        claims, err := m.jwtService.ValidateToken(c.Request.Context(), tokenString)
        if err != nil {
            monitoring.IncrementCounter("auth_errors_total", map[string]string{
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
        monitoring.IncrementCounter("auth_success_total", map[string]string{
            "user_id": claims.UserID,
        })
        
        c.Next()
    }
}

func (m *JWTMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check if user is authenticated
        userID, exists := c.Get("user_id")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
            c.Abort()
            return
        }
        
        // Get user roles
        roles, exists := c.Get("roles")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{"error": "No roles found"})
            c.Abort()
            return
        }
        
        // Check if user has required role
        userRoles := roles.([]string)
        hasRole := false
        for _, role := range userRoles {
            if role == requiredRole {
                hasRole = true
                break
            }
        }
        
        if !hasRole {
            monitoring.IncrementCounter("auth_errors_total", map[string]string{
                "error": "insufficient_permissions",
                "user_id": userID.(string),
            })
            c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}

func (m *JWTMiddleware) extractToken(authHeader string) string {
    // Check if header starts with "Bearer "
    if !strings.HasPrefix(authHeader, "Bearer ") {
        return ""
    }
    
    // Extract token
    return strings.TrimPrefix(authHeader, "Bearer ")
}
```

## 🔧 OAuth2 Authentication Implementation

### 1. OAuth2 Service

```go
// internal/auth/oauth2_service.go
package auth

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    
    "github.com/anasamu/go-micro-libs/auth"
    "github.com/anasamu/go-micro-libs/auth/providers/oauth2"
    "github.com/anasamu/go-micro-libs/database"
)

type OAuth2Service struct {
    oauth2Provider *oauth2.Provider
    dbManager      *database.DatabaseManager
    clientID       string
    clientSecret   string
    redirectURL    string
    authURL        string
    tokenURL       string
    userInfoURL    string
}

type OAuth2User struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Name     string `json:"name"`
    Picture  string `json:"picture"`
}

func NewOAuth2Service(oauth2Provider *oauth2.Provider, dbManager *database.DatabaseManager, 
    clientID, clientSecret, redirectURL, authURL, tokenURL, userInfoURL string) *OAuth2Service {
    return &OAuth2Service{
        oauth2Provider: oauth2Provider,
        dbManager:      dbManager,
        clientID:       clientID,
        clientSecret:   clientSecret,
        redirectURL:    redirectURL,
        authURL:        authURL,
        tokenURL:       tokenURL,
        userInfoURL:    userInfoURL,
    }
}

func (s *OAuth2Service) GetAuthURL(ctx context.Context, state string) string {
    params := url.Values{}
    params.Add("client_id", s.clientID)
    params.Add("redirect_uri", s.redirectURL)
    params.Add("response_type", "code")
    params.Add("scope", "openid email profile")
    params.Add("state", state)
    
    return fmt.Sprintf("%s?%s", s.authURL, params.Encode())
}

func (s *OAuth2Service) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
    // Exchange authorization code for access token
    tokenReq := &oauth2.TokenRequest{
        ClientID:     s.clientID,
        ClientSecret: s.clientSecret,
        Code:         code,
        RedirectURI:  s.redirectURL,
        GrantType:    "authorization_code",
    }
    
    token, err := s.oauth2Provider.ExchangeToken(ctx, tokenReq)
    if err != nil {
        return nil, err
    }
    
    return token, nil
}

func (s *OAuth2Service) GetUserInfo(ctx context.Context, token *oauth2.Token) (*OAuth2User, error) {
    // Get user info from OAuth2 provider
    userInfoReq := &oauth2.UserInfoRequest{
        AccessToken: token.AccessToken,
        UserInfoURL: s.userInfoURL,
    }
    
    userInfo, err := s.oauth2Provider.GetUserInfo(ctx, userInfoReq)
    if err != nil {
        return nil, err
    }
    
    // Parse user info
    var oauth2User OAuth2User
    err = json.Unmarshal(userInfo, &oauth2User)
    if err != nil {
        return nil, err
    }
    
    return &oauth2User, nil
}

func (s *OAuth2Service) CreateOrUpdateUser(ctx context.Context, oauth2User *OAuth2User) (*User, error) {
    // Check if user exists
    user, err := s.getUserByOAuth2ID(ctx, oauth2User.ID)
    if err != nil {
        // User doesn't exist, create new user
        user = &User{
            OAuth2ID:  oauth2User.ID,
            Username:  oauth2User.Username,
            Email:     oauth2User.Email,
            Name:      oauth2User.Name,
            Picture:   oauth2User.Picture,
            Roles:     []string{"user"},
        }
        
        err = s.createUser(ctx, user)
        if err != nil {
            return nil, err
        }
    } else {
        // Update existing user
        user.Username = oauth2User.Username
        user.Email = oauth2User.Email
        user.Name = oauth2User.Name
        user.Picture = oauth2User.Picture
        
        err = s.updateUser(ctx, user)
        if err != nil {
            return nil, err
        }
    }
    
    return user, nil
}

func (s *OAuth2Service) getUserByOAuth2ID(ctx context.Context, oauth2ID string) (*User, error) {
    query := `SELECT id, oauth2_id, username, email, name, picture, roles 
              FROM users WHERE oauth2_id = $1`
    
    result, err := s.dbManager.Query(ctx, "postgresql", query, oauth2ID)
    if err != nil {
        return nil, err
    }
    
    var user User
    if result.Next() {
        err = result.Scan(&user.ID, &user.OAuth2ID, &user.Username, 
            &user.Email, &user.Name, &user.Picture, &user.Roles)
        if err != nil {
            return nil, err
        }
    } else {
        return nil, fmt.Errorf("user not found")
    }
    
    return &user, nil
}

func (s *OAuth2Service) createUser(ctx context.Context, user *User) error {
    query := `INSERT INTO users (oauth2_id, username, email, name, picture, roles, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
    
    now := time.Now()
    _, err := s.dbManager.Exec(ctx, "postgresql", query, 
        user.OAuth2ID, user.Username, user.Email, user.Name, user.Picture, user.Roles, now, now)
    
    return err
}

func (s *OAuth2Service) updateUser(ctx context.Context, user *User) error {
    query := `UPDATE users SET username = $1, email = $2, name = $3, picture = $4, updated_at = $5 
              WHERE oauth2_id = $6`
    
    _, err := s.dbManager.Exec(ctx, "postgresql", query, 
        user.Username, user.Email, user.Name, user.Picture, time.Now(), user.OAuth2ID)
    
    return err
}
```

### 2. OAuth2 Handlers

```go
// internal/handlers/oauth2_handler.go
package handlers

import (
    "context"
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/anasamu/go-micro-libs/auth"
    "github.com/anasamu/go-micro-libs/monitoring"
)

type OAuth2Handler struct {
    oauth2Service *auth.OAuth2Service
    jwtService    *auth.JWTService
}

func NewOAuth2Handler(oauth2Service *auth.OAuth2Service, jwtService *auth.JWTService) *OAuth2Handler {
    return &OAuth2Handler{
        oauth2Service: oauth2Service,
        jwtService:    jwtService,
    }
}

func (h *OAuth2Handler) Login(c *gin.Context) {
    // Generate state parameter
    state := generateState()
    
    // Store state in session or cache
    // ... store state logic ...
    
    // Get authorization URL
    authURL := h.oauth2Service.GetAuthURL(c.Request.Context(), state)
    
    // Record metrics
    monitoring.IncrementCounter("oauth2_login_requests_total", map[string]string{
        "provider": "google",
    })
    
    c.JSON(http.StatusOK, gin.H{
        "auth_url": authURL,
        "state":    state,
    })
}

func (h *OAuth2Handler) Callback(c *gin.Context) {
    // Get authorization code and state
    code := c.Query("code")
    state := c.Query("state")
    
    if code == "" || state == "" {
        monitoring.IncrementCounter("oauth2_errors_total", map[string]string{
            "error": "missing_code_or_state",
        })
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or state parameter"})
        return
    }
    
    // Validate state parameter
    // ... validate state logic ...
    
    // Exchange code for token
    token, err := h.oauth2Service.ExchangeCodeForToken(c.Request.Context(), code)
    if err != nil {
        monitoring.IncrementCounter("oauth2_errors_total", map[string]string{
            "error": "token_exchange_failed",
        })
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
        return
    }
    
    // Get user info
    oauth2User, err := h.oauth2Service.GetUserInfo(c.Request.Context(), token)
    if err != nil {
        monitoring.IncrementCounter("oauth2_errors_total", map[string]string{
            "error": "user_info_failed",
        })
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
        return
    }
    
    // Create or update user
    user, err := h.oauth2Service.CreateOrUpdateUser(c.Request.Context(), oauth2User)
    if err != nil {
        monitoring.IncrementCounter("oauth2_errors_total", map[string]string{
            "error": "user_creation_failed",
        })
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create/update user"})
        return
    }
    
    // Generate JWT token
    jwtToken, err := h.jwtService.GenerateToken(c.Request.Context(), user)
    if err != nil {
        monitoring.IncrementCounter("oauth2_errors_total", map[string]string{
            "error": "jwt_generation_failed",
        })
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
        return
    }
    
    // Record success
    monitoring.IncrementCounter("oauth2_success_total", map[string]string{
        "user_id": user.ID,
    })
    
    c.JSON(http.StatusOK, gin.H{
        "access_token": jwtToken,
        "token_type":   "Bearer",
        "expires_in":   86400, // 24 hours
        "user":         user,
    })
}

func generateState() string {
    // Generate random state parameter
    return "random-state-string"
}
```

## 🔧 2FA Authentication Implementation

### 1. TOTP Service

```go
// internal/auth/totp_service.go
package auth

import (
    "context"
    "crypto/rand"
    "encoding/base32"
    "fmt"
    "time"
    
    "github.com/pquerna/otp"
    "github.com/pquerna/otp/totp"
    "github.com/anasamu/go-micro-libs/auth"
    "github.com/anasamu/go-micro-libs/database"
)

type TOTPService struct {
    dbManager *database.DatabaseManager
    issuer    string
    algorithm otp.Algorithm
    digits    otp.Digits
    period    uint
}

func NewTOTPService(dbManager *database.DatabaseManager, issuer string) *TOTPService {
    return &TOTPService{
        dbManager: dbManager,
        issuer:    issuer,
        algorithm: otp.AlgorithmSHA1,
        digits:    otp.DigitsSix,
        period:    30,
    }
}

func (s *TOTPService) GenerateSecret(ctx context.Context, userID string) (*otp.Key, error) {
    // Generate TOTP secret
    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      s.issuer,
        AccountName: userID,
        Algorithm:   s.algorithm,
        Digits:      s.digits,
        Period:      s.period,
    })
    if err != nil {
        return nil, err
    }
    
    // Store secret in database
    err = s.storeSecret(ctx, userID, key.Secret())
    if err != nil {
        return nil, err
    }
    
    return key, nil
}

func (s *TOTPService) ValidateCode(ctx context.Context, userID, code string) (bool, error) {
    // Get user's secret
    secret, err := s.getSecret(ctx, userID)
    if err != nil {
        return false, err
    }
    
    // Validate TOTP code
    valid := totp.Validate(code, secret)
    
    // Record validation attempt
    s.recordValidationAttempt(ctx, userID, valid)
    
    return valid, nil
}

func (s *TOTPService) GenerateRecoveryCodes(ctx context.Context, userID string) ([]string, error) {
    // Generate recovery codes
    codes := make([]string, 10)
    for i := 0; i < 10; i++ {
        code := s.generateRecoveryCode()
        codes[i] = code
    }
    
    // Store recovery codes
    err := s.storeRecoveryCodes(ctx, userID, codes)
    if err != nil {
        return nil, err
    }
    
    return codes, nil
}

func (s *TOTPService) ValidateRecoveryCode(ctx context.Context, userID, code string) (bool, error) {
    // Get user's recovery codes
    codes, err := s.getRecoveryCodes(ctx, userID)
    if err != nil {
        return false, err
    }
    
    // Check if code exists
    for i, recoveryCode := range codes {
        if recoveryCode == code {
            // Remove used code
            codes = append(codes[:i], codes[i+1:]...)
            s.storeRecoveryCodes(ctx, userID, codes)
            return true, nil
        }
    }
    
    return false, nil
}

func (s *TOTPService) storeSecret(ctx context.Context, userID, secret string) error {
    query := `INSERT INTO user_totp_secrets (user_id, secret, created_at) 
              VALUES ($1, $2, $3) 
              ON CONFLICT (user_id) DO UPDATE SET secret = $2, updated_at = $3`
    
    _, err := s.dbManager.Exec(ctx, "postgresql", query, userID, secret, time.Now())
    return err
}

func (s *TOTPService) getSecret(ctx context.Context, userID string) (string, error) {
    query := `SELECT secret FROM user_totp_secrets WHERE user_id = $1`
    
    result, err := s.dbManager.Query(ctx, "postgresql", query, userID)
    if err != nil {
        return "", err
    }
    
    var secret string
    if result.Next() {
        err = result.Scan(&secret)
        if err != nil {
            return "", err
        }
    } else {
        return "", fmt.Errorf("TOTP secret not found")
    }
    
    return secret, nil
}

func (s *TOTPService) recordValidationAttempt(ctx context.Context, userID string, valid bool) {
    query := `INSERT INTO totp_validation_attempts (user_id, valid, created_at) 
              VALUES ($1, $2, $3)`
    
    _, err := s.dbManager.Exec(ctx, "postgresql", query, userID, valid, time.Now())
    if err != nil {
        log.Printf("Failed to record validation attempt: %v", err)
    }
}

func (s *TOTPService) storeRecoveryCodes(ctx context.Context, userID string, codes []string) error {
    query := `INSERT INTO user_recovery_codes (user_id, codes, created_at) 
              VALUES ($1, $2, $3) 
              ON CONFLICT (user_id) DO UPDATE SET codes = $2, updated_at = $3`
    
    codesJSON, err := json.Marshal(codes)
    if err != nil {
        return err
    }
    
    _, err = s.dbManager.Exec(ctx, "postgresql", query, userID, codesJSON, time.Now())
    return err
}

func (s *TOTPService) getRecoveryCodes(ctx context.Context, userID string) ([]string, error) {
    query := `SELECT codes FROM user_recovery_codes WHERE user_id = $1`
    
    result, err := s.dbManager.Query(ctx, "postgresql", query, userID)
    if err != nil {
        return nil, err
    }
    
    var codesJSON []byte
    if result.Next() {
        err = result.Scan(&codesJSON)
        if err != nil {
            return nil, err
        }
    } else {
        return nil, fmt.Errorf("recovery codes not found")
    }
    
    var codes []string
    err = json.Unmarshal(codesJSON, &codes)
    if err != nil {
        return nil, err
    }
    
    return codes, nil
}

func (s *TOTPService) generateRecoveryCode() string {
    // Generate random recovery code
    bytes := make([]byte, 8)
    rand.Read(bytes)
    return base32.StdEncoding.EncodeToString(bytes)
}
```

### 2. 2FA Handlers

```go
// internal/handlers/2fa_handler.go
package handlers

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/anasamu/go-micro-libs/auth"
    "github.com/anasamu/go-micro-libs/monitoring"
)

type TwoFAHandler struct {
    totpService *auth.TOTPService
    jwtService  *auth.JWTService
}

func NewTwoFAHandler(totpService *auth.TOTPService, jwtService *auth.JWTService) *TwoFAHandler {
    return &TwoFAHandler{
        totpService: totpService,
        jwtService:  jwtService,
    }
}

func (h *TwoFAHandler) SetupTOTP(c *gin.Context) {
    userID := c.GetString("user_id")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
        return
    }
    
    // Generate TOTP secret
    key, err := h.totpService.GenerateSecret(c.Request.Context(), userID)
    if err != nil {
        monitoring.IncrementCounter("2fa_errors_total", map[string]string{
            "error": "secret_generation_failed",
        })
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate TOTP secret"})
        return
    }
    
    // Generate recovery codes
    recoveryCodes, err := h.totpService.GenerateRecoveryCodes(c.Request.Context(), userID)
    if err != nil {
        monitoring.IncrementCounter("2fa_errors_total", map[string]string{
            "error": "recovery_codes_generation_failed",
        })
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate recovery codes"})
        return
    }
    
    monitoring.IncrementCounter("2fa_setup_total", map[string]string{
        "user_id": userID,
    })
    
    c.JSON(http.StatusOK, gin.H{
        "secret":         key.Secret(),
        "qr_code_url":    key.URL(),
        "recovery_codes": recoveryCodes,
    })
}

func (h *TwoFAHandler) VerifyTOTP(c *gin.Context) {
    userID := c.GetString("user_id")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
        return
    }
    
    var req struct {
        Code string `json:"code" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Validate TOTP code
    valid, err := h.totpService.ValidateCode(c.Request.Context(), userID, req.Code)
    if err != nil {
        monitoring.IncrementCounter("2fa_errors_total", map[string]string{
            "error": "validation_failed",
        })
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate TOTP code"})
        return
    }
    
    if !valid {
        monitoring.IncrementCounter("2fa_errors_total", map[string]string{
            "error": "invalid_code",
        })
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TOTP code"})
        return
    }
    
    monitoring.IncrementCounter("2fa_verification_success_total", map[string]string{
        "user_id": userID,
    })
    
    c.JSON(http.StatusOK, gin.H{"message": "TOTP code verified successfully"})
}

func (h *TwoFAHandler) VerifyRecoveryCode(c *gin.Context) {
    userID := c.GetString("user_id")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
        return
    }
    
    var req struct {
        Code string `json:"code" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Validate recovery code
    valid, err := h.totpService.ValidateRecoveryCode(c.Request.Context(), userID, req.Code)
    if err != nil {
        monitoring.IncrementCounter("2fa_errors_total", map[string]string{
            "error": "recovery_validation_failed",
        })
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate recovery code"})
        return
    }
    
    if !valid {
        monitoring.IncrementCounter("2fa_errors_total", map[string]string{
            "error": "invalid_recovery_code",
        })
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recovery code"})
        return
    }
    
    monitoring.IncrementCounter("2fa_recovery_success_total", map[string]string{
        "user_id": userID,
    })
    
    c.JSON(http.StatusOK, gin.H{"message": "Recovery code verified successfully"})
}
```

## 🔧 Best Practices

### 1. Token Security

```go
// Secure token generation
func (s *JWTService) GenerateSecureToken(ctx context.Context, user *User) (string, error) {
    // Use strong secret
    if len(s.secret) < 32 {
        return "", errors.New("secret too short")
    }
    
    // Set short expiration
    claims := &Claims{
        UserID: user.ID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Short expiration
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.secret))
}
```

### 2. Rate Limiting

```go
// Rate limiting for authentication endpoints
func (h *AuthHandler) LoginWithRateLimit(c *gin.Context) {
    clientIP := c.ClientIP()
    
    // Check rate limit
    if h.rateLimiter.IsLimited(clientIP) {
        monitoring.IncrementCounter("auth_rate_limited_total", map[string]string{
            "ip": clientIP,
        })
        c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
        return
    }
    
    // Process login
    h.Login(c)
}
```

### 3. Audit Logging

```go
// Audit logging for authentication events
func (s *JWTService) GenerateTokenWithAudit(ctx context.Context, user *User) (string, error) {
    token, err := s.GenerateToken(ctx, user)
    if err != nil {
        return "", err
    }
    
    // Log authentication event
    s.auditLogger.LogAuthEvent(ctx, &AuditEvent{
        UserID:    user.ID,
        Action:    "token_generated",
        Timestamp: time.Now(),
        IP:        getClientIP(ctx),
        UserAgent: getUserAgent(ctx),
    })
    
    return token, nil
}
```

### 4. Session Management

```go
// Session management for authentication
func (s *JWTService) RevokeAllUserTokens(ctx context.Context, userID string) error {
    // Blacklist all user tokens
    query := `UPDATE user_tokens SET revoked = true WHERE user_id = $1`
    
    _, err := s.dbManager.Exec(ctx, "postgresql", query, userID)
    if err != nil {
        return err
    }
    
    // Log revocation event
    s.auditLogger.LogAuthEvent(ctx, &AuditEvent{
        UserID:    userID,
        Action:    "all_tokens_revoked",
        Timestamp: time.Now(),
    })
    
    return nil
}
```

---

**Authentication - Secure and comprehensive authentication for microservices! 🚀**

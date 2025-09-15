# Library Combination Guide

## 🎯 Overview

GoMicroFramework memungkinkan Anda untuk menggabungkan 2-5 library dalam satu service untuk membangun aplikasi yang powerful dan feature-rich. Panduan ini menunjukkan cara menggabungkan library secara efektif dengan contoh praktis.

## 🔧 Common Library Combinations

### 1. Database + Cache + Monitoring (3 Libraries)

#### Use Case: High-Performance Data Service

```go
// cmd/main.go
package main

import (
    "context"
    "log"
    
    "github.com/anasamu/go-micro-libs/config"
    "github.com/anasamu/go-micro-libs/logging"
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/cache"
)

func main() {
    ctx := context.Background()
    
    // Initialize core managers
    configManager := config.NewManager()
    loggingManager := logging.NewManager()
    monitoringManager := monitoring.NewManager()
    
    // Initialize combination libraries
    dbManager := database.NewManager()
    cacheManager := cache.NewManager()
    
    // Bootstrap service
    if err := bootstrapService(ctx, configManager, loggingManager, 
        monitoringManager, dbManager, cacheManager); err != nil {
        log.Fatal("Failed to bootstrap service:", err)
    }
    
    log.Println("Data service started successfully")
}

func bootstrapService(ctx context.Context, 
    configManager *config.ConfigManager,
    loggingManager *logging.LoggingManager,
    monitoringManager *monitoring.MonitoringManager,
    dbManager *database.DatabaseManager,
    cacheManager *cache.CacheManager) error {
    
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
    
    // Connect to database
    if err := dbManager.Connect(ctx); err != nil {
        return err
    }
    
    // Connect to cache
    if err := cacheManager.Connect(ctx); err != nil {
        return err
    }
    
    // Register health checks
    monitoringManager.RegisterHealthCheck("database", func() error {
        return dbManager.HealthCheck(ctx)
    })
    
    monitoringManager.RegisterHealthCheck("cache", func() error {
        return cacheManager.HealthCheck(ctx)
    })
    
    return nil
}
```

#### Service Implementation

```go
// internal/services/user_service.go
package services

import (
    "context"
    "time"
    
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/cache"
    "github.com/anasamu/go-micro-libs/monitoring"
)

type UserService struct {
    dbManager       *database.DatabaseManager
    cacheManager    *cache.CacheManager
    monitoringManager *monitoring.MonitoringManager
}

func NewUserService(dbManager *database.DatabaseManager, 
    cacheManager *cache.CacheManager,
    monitoringManager *monitoring.MonitoringManager) *UserService {
    return &UserService{
        dbManager:       dbManager,
        cacheManager:    cacheManager,
        monitoringManager: monitoringManager,
    }
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("user:%s", userID)
    var user User
    
    err := s.cacheManager.Get(ctx, cacheKey, &user)
    if err == nil {
        // Cache hit - record metric
        s.monitoringManager.IncrementCounter("cache_hits_total", map[string]string{
            "operation": "get_user",
        })
        return &user, nil
    }
    
    // Cache miss - record metric
    s.monitoringManager.IncrementCounter("cache_misses_total", map[string]string{
        "operation": "get_user",
    })
    
    // Get from database
    start := time.Now()
    user, err = s.getUserFromDB(ctx, userID)
    if err != nil {
        s.monitoringManager.IncrementCounter("database_errors_total", map[string]string{
            "operation": "get_user",
        })
        return nil, err
    }
    
    // Record database latency
    s.monitoringManager.RecordHistogram("database_latency_seconds", 
        time.Since(start).Seconds(), map[string]string{
            "operation": "get_user",
        })
    
    // Cache the result
    s.cacheManager.Set(ctx, cacheKey, user, 1*time.Hour)
    
    return &user, nil
}

func (s *UserService) getUserFromDB(ctx context.Context, userID string) (*User, error) {
    query := "SELECT id, name, email, created_at FROM users WHERE id = ?"
    result, err := s.dbManager.Query(ctx, "postgresql", query, userID)
    if err != nil {
        return nil, err
    }
    
    var user User
    if result.Next() {
        err = result.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
        if err != nil {
            return nil, err
        }
    }
    
    return &user, nil
}
```

#### Configuration

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
    env:
      prefix: "USER_SERVICE_"

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

# Combination libraries
database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100
      max_idle_connections: 10

cache:
  providers:
    redis:
      url: "${REDIS_URL}"
      db: 0
      pool_size: 10
      ttl: "1h"
```

### 2. Database + Messaging + AI (3 Libraries)

#### Use Case: AI-Powered Content Service

```go
// cmd/main.go
package main

import (
    "context"
    "log"
    
    "github.com/anasamu/go-micro-libs/config"
    "github.com/anasamu/go-micro-libs/logging"
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/messaging"
    "github.com/anasamu/go-micro-libs/ai"
)

func main() {
    ctx := context.Background()
    
    // Initialize core managers
    configManager := config.NewManager()
    loggingManager := logging.NewManager()
    monitoringManager := monitoring.NewManager()
    
    // Initialize combination libraries
    dbManager := database.NewManager()
    msgManager := messaging.NewManager()
    aiManager := ai.NewManager()
    
    // Bootstrap service
    if err := bootstrapService(ctx, configManager, loggingManager, 
        monitoringManager, dbManager, msgManager, aiManager); err != nil {
        log.Fatal("Failed to bootstrap service:", err)
    }
    
    log.Println("AI content service started successfully")
}

func bootstrapService(ctx context.Context, 
    configManager *config.ConfigManager,
    loggingManager *logging.LoggingManager,
    monitoringManager *monitoring.MonitoringManager,
    dbManager *database.DatabaseManager,
    msgManager *messaging.MessagingManager,
    aiManager *ai.AIManager) error {
    
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
    
    // Connect to database
    if err := dbManager.Connect(ctx); err != nil {
        return err
    }
    
    // Connect to messaging
    if err := msgManager.Connect(ctx); err != nil {
        return err
    }
    
    // Initialize AI
    if err := aiManager.Initialize(); err != nil {
        return err
    }
    
    return nil
}
```

#### Service Implementation

```go
// internal/services/content_service.go
package services

import (
    "context"
    "encoding/json"
    
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/messaging"
    "github.com/anasamu/go-micro-libs/ai"
    "github.com/anasamu/go-micro-libs/ai/types"
)

type ContentService struct {
    dbManager    *database.DatabaseManager
    msgManager   *messaging.MessagingManager
    aiManager    *ai.AIManager
}

func NewContentService(dbManager *database.DatabaseManager, 
    msgManager *messaging.MessagingManager,
    aiManager *ai.AIManager) *ContentService {
    return &ContentService{
        dbManager:  dbManager,
        msgManager: msgManager,
        aiManager:  aiManager,
    }
}

func (s *ContentService) GenerateContent(ctx context.Context, request *ContentRequest) (*Content, error) {
    // Generate content using AI
    chatReq := &types.ChatRequest{
        Messages: []types.Message{
            {Role: "user", Content: request.Prompt},
        },
        Model: "gpt-4",
    }
    
    response, err := s.aiManager.Chat(ctx, "openai", chatReq)
    if err != nil {
        return nil, err
    }
    
    // Create content object
    content := &Content{
        ID:          generateID(),
        Title:       request.Title,
        Content:     response.Choices[0].Message.Content,
        AuthorID:    request.AuthorID,
        Status:      "draft",
        CreatedAt:   time.Now(),
    }
    
    // Save to database
    err = s.saveContent(ctx, content)
    if err != nil {
        return nil, err
    }
    
    // Publish event
    event := &ContentCreatedEvent{
        ContentID: content.ID,
        AuthorID:  content.AuthorID,
        Title:     content.Title,
        CreatedAt: content.CreatedAt,
    }
    
    eventData, _ := json.Marshal(event)
    message := messaging.CreateMessage("content.created", eventData)
    
    publishReq := &messaging.PublishRequest{
        Topic:   "content-events",
        Message: message,
    }
    
    _, err = s.msgManager.PublishMessage(ctx, "kafka", publishReq)
    if err != nil {
        // Log error but don't fail the operation
        log.Printf("Failed to publish event: %v", err)
    }
    
    return content, nil
}

func (s *ContentService) saveContent(ctx context.Context, content *Content) error {
    query := `INSERT INTO contents (id, title, content, author_id, status, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6)`
    
    _, err := s.dbManager.Exec(ctx, "postgresql", query, 
        content.ID, content.Title, content.Content, content.AuthorID, 
        content.Status, content.CreatedAt)
    
    return err
}
```

#### Configuration

```yaml
# config.yaml
service:
  name: "content-service"
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

# Combination libraries
database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100

messaging:
  providers:
    kafka:
      brokers: ["localhost:9092"]
      group_id: "content-service"
      topics: ["content-events"]

ai:
  providers:
    openai:
      api_key: "${OPENAI_API_KEY}"
      default_model: "gpt-4"
```

### 3. Database + Cache + Storage + Payment (4 Libraries)

#### Use Case: E-commerce Order Service

```go
// cmd/main.go
package main

import (
    "context"
    "log"
    
    "github.com/anasamu/go-micro-libs/config"
    "github.com/anasamu/go-micro-libs/logging"
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/cache"
    "github.com/anasamu/go-micro-libs/storage"
    "github.com/anasamu/go-micro-libs/payment"
)

func main() {
    ctx := context.Background()
    
    // Initialize core managers
    configManager := config.NewManager()
    loggingManager := logging.NewManager()
    monitoringManager := monitoring.NewManager()
    
    // Initialize combination libraries
    dbManager := database.NewManager()
    cacheManager := cache.NewManager()
    storageManager := storage.NewManager()
    paymentManager := payment.NewManager()
    
    // Bootstrap service
    if err := bootstrapService(ctx, configManager, loggingManager, 
        monitoringManager, dbManager, cacheManager, storageManager, paymentManager); err != nil {
        log.Fatal("Failed to bootstrap service:", err)
    }
    
    log.Println("E-commerce order service started successfully")
}

func bootstrapService(ctx context.Context, 
    configManager *config.ConfigManager,
    loggingManager *logging.LoggingManager,
    monitoringManager *monitoring.MonitoringManager,
    dbManager *database.DatabaseManager,
    cacheManager *cache.CacheManager,
    storageManager *storage.StorageManager,
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
    
    // Connect to database
    if err := dbManager.Connect(ctx); err != nil {
        return err
    }
    
    // Connect to cache
    if err := cacheManager.Connect(ctx); err != nil {
        return err
    }
    
    // Initialize storage
    if err := storageManager.Initialize(); err != nil {
        return err
    }
    
    // Initialize payment
    if err := paymentManager.Initialize(); err != nil {
        return err
    }
    
    return nil
}
```

#### Service Implementation

```go
// internal/services/order_service.go
package services

import (
    "context"
    "fmt"
    "time"
    
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/cache"
    "github.com/anasamu/go-micro-libs/storage"
    "github.com/anasamu/go-micro-libs/payment"
)

type OrderService struct {
    dbManager       *database.DatabaseManager
    cacheManager    *cache.CacheManager
    storageManager  *storage.StorageManager
    paymentManager  *payment.PaymentManager
}

func NewOrderService(dbManager *database.DatabaseManager, 
    cacheManager *cache.CacheManager,
    storageManager *storage.StorageManager,
    paymentManager *payment.PaymentManager) *OrderService {
    return &OrderService{
        dbManager:      dbManager,
        cacheManager:   cacheManager,
        storageManager: storageManager,
        paymentManager: paymentManager,
    }
}

func (s *OrderService) CreateOrder(ctx context.Context, request *CreateOrderRequest) (*Order, error) {
    // Check product availability in cache
    cacheKey := fmt.Sprintf("product:%s:stock", request.ProductID)
    var stock int
    err := s.cacheManager.Get(ctx, cacheKey, &stock)
    if err != nil {
        // Cache miss - get from database
        stock, err = s.getProductStock(ctx, request.ProductID)
        if err != nil {
            return nil, err
        }
        
        // Cache the result
        s.cacheManager.Set(ctx, cacheKey, stock, 5*time.Minute)
    }
    
    if stock < request.Quantity {
        return nil, fmt.Errorf("insufficient stock")
    }
    
    // Process payment
    paymentReq := &payment.PaymentRequest{
        Amount:   request.TotalAmount,
        Currency: "USD",
        CustomerID: request.CustomerID,
        PaymentMethodID: request.PaymentMethodID,
    }
    
    paymentResp, err := s.paymentManager.ProcessPayment(ctx, "stripe", paymentReq)
    if err != nil {
        return nil, fmt.Errorf("payment failed: %w", err)
    }
    
    // Create order
    order := &Order{
        ID:            generateID(),
        CustomerID:    request.CustomerID,
        ProductID:     request.ProductID,
        Quantity:      request.Quantity,
        TotalAmount:   request.TotalAmount,
        PaymentID:     paymentResp.PaymentID,
        Status:        "confirmed",
        CreatedAt:     time.Now(),
    }
    
    // Save order to database
    err = s.saveOrder(ctx, order)
    if err != nil {
        // Refund payment if order creation fails
        s.paymentManager.RefundPayment(ctx, "stripe", &payment.RefundRequest{
            PaymentID: paymentResp.PaymentID,
            Amount:    request.TotalAmount,
        })
        return nil, err
    }
    
    // Update product stock
    err = s.updateProductStock(ctx, request.ProductID, stock-request.Quantity)
    if err != nil {
        return nil, err
    }
    
    // Invalidate cache
    s.cacheManager.Delete(ctx, cacheKey)
    
    // Generate invoice and store in storage
    invoiceData := s.generateInvoice(order)
    invoiceKey := fmt.Sprintf("invoices/%s.pdf", order.ID)
    
    putReq := &storage.PutObjectRequest{
        Bucket:      "ecommerce-invoices",
        Key:         invoiceKey,
        Content:     strings.NewReader(invoiceData),
        Size:        int64(len(invoiceData)),
        ContentType: "application/pdf",
    }
    
    _, err = s.storageManager.PutObject(ctx, "s3", putReq)
    if err != nil {
        // Log error but don't fail the operation
        log.Printf("Failed to store invoice: %v", err)
    }
    
    return order, nil
}

func (s *OrderService) getProductStock(ctx context.Context, productID string) (int, error) {
    query := "SELECT stock FROM products WHERE id = ?"
    result, err := s.dbManager.Query(ctx, "postgresql", query, productID)
    if err != nil {
        return 0, err
    }
    
    var stock int
    if result.Next() {
        err = result.Scan(&stock)
        if err != nil {
            return 0, err
        }
    }
    
    return stock, nil
}

func (s *OrderService) saveOrder(ctx context.Context, order *Order) error {
    query := `INSERT INTO orders (id, customer_id, product_id, quantity, total_amount, 
              payment_id, status, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
    
    _, err := s.dbManager.Exec(ctx, "postgresql", query, 
        order.ID, order.CustomerID, order.ProductID, order.Quantity, 
        order.TotalAmount, order.PaymentID, order.Status, order.CreatedAt)
    
    return err
}

func (s *OrderService) updateProductStock(ctx context.Context, productID string, newStock int) error {
    query := "UPDATE products SET stock = $1 WHERE id = $2"
    _, err := s.dbManager.Exec(ctx, "postgresql", query, newStock, productID)
    return err
}

func (s *OrderService) generateInvoice(order *Order) string {
    // Generate invoice content (simplified)
    return fmt.Sprintf("Invoice for Order %s\nAmount: $%.2f", order.ID, order.TotalAmount)
}
```

#### Configuration

```yaml
# config.yaml
service:
  name: "order-service"
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

# Combination libraries
database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100

cache:
  providers:
    redis:
      url: "${REDIS_URL}"
      db: 0
      pool_size: 10

storage:
  providers:
    s3:
      region: "us-east-1"
      access_key_id: "${AWS_ACCESS_KEY_ID}"
      secret_access_key: "${AWS_SECRET_ACCESS_KEY}"
      bucket: "ecommerce-invoices"

payment:
  providers:
    stripe:
      secret_key: "${STRIPE_SECRET_KEY}"
      publishable_key: "${STRIPE_PUBLISHABLE_KEY}"
```

### 4. Database + Cache + Messaging + AI + Storage (5 Libraries)

#### Use Case: AI-Powered Document Processing Service

```go
// cmd/main.go
package main

import (
    "context"
    "log"
    
    "github.com/anasamu/go-micro-libs/config"
    "github.com/anasamu/go-micro-libs/logging"
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/cache"
    "github.com/anasamu/go-micro-libs/messaging"
    "github.com/anasamu/go-micro-libs/ai"
    "github.com/anasamu/go-micro-libs/storage"
)

func main() {
    ctx := context.Background()
    
    // Initialize core managers
    configManager := config.NewManager()
    loggingManager := logging.NewManager()
    monitoringManager := monitoring.NewManager()
    
    // Initialize combination libraries
    dbManager := database.NewManager()
    cacheManager := cache.NewManager()
    msgManager := messaging.NewManager()
    aiManager := ai.NewManager()
    storageManager := storage.NewManager()
    
    // Bootstrap service
    if err := bootstrapService(ctx, configManager, loggingManager, 
        monitoringManager, dbManager, cacheManager, msgManager, aiManager, storageManager); err != nil {
        log.Fatal("Failed to bootstrap service:", err)
    }
    
    log.Println("AI document processing service started successfully")
}

func bootstrapService(ctx context.Context, 
    configManager *config.ConfigManager,
    loggingManager *logging.LoggingManager,
    monitoringManager *monitoring.MonitoringManager,
    dbManager *database.DatabaseManager,
    cacheManager *cache.CacheManager,
    msgManager *messaging.MessagingManager,
    aiManager *ai.AIManager,
    storageManager *storage.StorageManager) error {
    
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
    
    // Connect to database
    if err := dbManager.Connect(ctx); err != nil {
        return err
    }
    
    // Connect to cache
    if err := cacheManager.Connect(ctx); err != nil {
        return err
    }
    
    // Connect to messaging
    if err := msgManager.Connect(ctx); err != nil {
        return err
    }
    
    // Initialize AI
    if err := aiManager.Initialize(); err != nil {
        return err
    }
    
    // Initialize storage
    if err := storageManager.Initialize(); err != nil {
        return err
    }
    
    return nil
}
```

#### Service Implementation

```go
// internal/services/document_service.go
package services

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/cache"
    "github.com/anasamu/go-micro-libs/messaging"
    "github.com/anasamu/go-micro-libs/ai"
    "github.com/anasamu/go-micro-libs/ai/types"
    "github.com/anasamu/go-micro-libs/storage"
)

type DocumentService struct {
    dbManager      *database.DatabaseManager
    cacheManager   *cache.CacheManager
    msgManager     *messaging.MessagingManager
    aiManager      *ai.AIManager
    storageManager *storage.StorageManager
}

func NewDocumentService(dbManager *database.DatabaseManager, 
    cacheManager *cache.CacheManager,
    msgManager *messaging.MessagingManager,
    aiManager *ai.AIManager,
    storageManager *storage.StorageManager) *DocumentService {
    return &DocumentService{
        dbManager:      dbManager,
        cacheManager:   cacheManager,
        msgManager:     msgManager,
        aiManager:      aiManager,
        storageManager: storageManager,
    }
}

func (s *DocumentService) ProcessDocument(ctx context.Context, request *ProcessDocumentRequest) (*Document, error) {
    // Check if document is already processed (cache)
    cacheKey := fmt.Sprintf("document:%s:processed", request.DocumentID)
    var processedDoc Document
    err := s.cacheManager.Get(ctx, cacheKey, &processedDoc)
    if err == nil {
        return &processedDoc, nil
    }
    
    // Get document from storage
    getReq := &storage.GetObjectRequest{
        Bucket: "documents",
        Key:    request.DocumentID,
    }
    
    getResp, err := s.storageManager.GetObject(ctx, "s3", getReq)
    if err != nil {
        return nil, fmt.Errorf("failed to get document: %w", err)
    }
    
    // Process document with AI
    chatReq := &types.ChatRequest{
        Messages: []types.Message{
            {Role: "user", Content: fmt.Sprintf("Process this document: %s", string(getResp.Content))},
        },
        Model: "gpt-4",
    }
    
    response, err := s.aiManager.Chat(ctx, "openai", chatReq)
    if err != nil {
        return nil, fmt.Errorf("AI processing failed: %w", err)
    }
    
    // Create processed document
    processedDoc = Document{
        ID:            request.DocumentID,
        OriginalPath:  request.DocumentID,
        ProcessedContent: response.Choices[0].Message.Content,
        Status:        "processed",
        ProcessedAt:   time.Now(),
    }
    
    // Save to database
    err = s.saveDocument(ctx, &processedDoc)
    if err != nil {
        return nil, err
    }
    
    // Cache the result
    s.cacheManager.Set(ctx, cacheKey, processedDoc, 1*time.Hour)
    
    // Publish processing complete event
    event := &DocumentProcessedEvent{
        DocumentID: processedDoc.ID,
        Status:     processedDoc.Status,
        ProcessedAt: processedDoc.ProcessedAt,
    }
    
    eventData, _ := json.Marshal(event)
    message := messaging.CreateMessage("document.processed", eventData)
    
    publishReq := &messaging.PublishRequest{
        Topic:   "document-events",
        Message: message,
    }
    
    _, err = s.msgManager.PublishMessage(ctx, "kafka", publishReq)
    if err != nil {
        log.Printf("Failed to publish event: %v", err)
    }
    
    return &processedDoc, nil
}

func (s *DocumentService) saveDocument(ctx context.Context, doc *Document) error {
    query := `INSERT INTO documents (id, original_path, processed_content, status, processed_at) 
              VALUES ($1, $2, $3, $4, $5)`
    
    _, err := s.dbManager.Exec(ctx, "postgresql", query, 
        doc.ID, doc.OriginalPath, doc.ProcessedContent, doc.Status, doc.ProcessedAt)
    
    return err
}
```

#### Configuration

```yaml
# config.yaml
service:
  name: "document-service"
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

# Combination libraries
database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100

cache:
  providers:
    redis:
      url: "${REDIS_URL}"
      db: 0
      pool_size: 10

messaging:
  providers:
    kafka:
      brokers: ["localhost:9092"]
      group_id: "document-service"
      topics: ["document-events"]

ai:
  providers:
    openai:
      api_key: "${OPENAI_API_KEY}"
      default_model: "gpt-4"

storage:
  providers:
    s3:
      region: "us-east-1"
      access_key_id: "${AWS_ACCESS_KEY_ID}"
      secret_access_key: "${AWS_SECRET_ACCESS_KEY}"
      bucket: "documents"
```

## 🔧 Best Practices for Library Combination

### 1. Initialization Order

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
    dbManager := database.NewManager()
    if err := dbManager.Connect(ctx); err != nil {
        return err
    }
    
    // 5. Cache
    cacheManager := cache.NewManager()
    if err := cacheManager.Connect(ctx); err != nil {
        return err
    }
    
    // 6. Other libraries...
    
    return nil
}
```

### 2. Error Handling

```go
// Proper error handling for library combinations
func (s *UserService) GetUserWithFallback(ctx context.Context, userID string) (*User, error) {
    // Try cache first
    user, err := s.getUserFromCache(ctx, userID)
    if err == nil {
        return user, nil
    }
    
    // Fallback to database
    user, err = s.getUserFromDB(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // Cache the result for next time
    s.cacheUser(ctx, user)
    
    return user, nil
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
    
    // Close storage connections
    storageManager.Close()
}
```

### 4. Health Checks

```go
// Register health checks for all libraries
func registerHealthChecks() {
    // Database health check
    monitoringManager.RegisterHealthCheck("database", func() error {
        return dbManager.HealthCheck(ctx)
    })
    
    // Cache health check
    monitoringManager.RegisterHealthCheck("cache", func() error {
        return cacheManager.HealthCheck(ctx)
    })
    
    // Messaging health check
    monitoringManager.RegisterHealthCheck("messaging", func() error {
        return msgManager.HealthCheck(ctx)
    })
    
    // AI health check
    monitoringManager.RegisterHealthCheck("ai", func() error {
        return aiManager.HealthCheck(ctx)
    })
    
    // Storage health check
    monitoringManager.RegisterHealthCheck("storage", func() error {
        return storageManager.HealthCheck(ctx)
    })
}
```

## 🔧 Performance Optimization

### 1. Connection Pooling

```yaml
# Optimize connection pools for multiple libraries
database:
  providers:
    postgresql:
      max_connections: 100
      max_idle_connections: 10
      connection_max_lifetime: "1h"

cache:
  providers:
    redis:
      pool_size: 20
      min_idle_conns: 5
      max_conn_age: "1h"
```

### 2. Caching Strategy

```go
// Implement multi-level caching
func (s *UserService) GetUserWithMultiLevelCache(ctx context.Context, userID string) (*User, error) {
    // L1: Memory cache
    user, err := s.getUserFromMemoryCache(userID)
    if err == nil {
        return user, nil
    }
    
    // L2: Redis cache
    user, err = s.getUserFromRedisCache(ctx, userID)
    if err == nil {
        s.setMemoryCache(userID, user)
        return user, nil
    }
    
    // L3: Database
    user, err = s.getUserFromDB(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // Update all caches
    s.setMemoryCache(userID, user)
    s.setRedisCache(ctx, userID, user)
    
    return user, nil
}
```

### 3. Batch Operations

```go
// Batch operations for better performance
func (s *UserService) GetUsersBatch(ctx context.Context, userIDs []string) ([]*User, error) {
    // Check cache for all users
    cachedUsers := make(map[string]*User)
    missingIDs := make([]string, 0)
    
    for _, userID := range userIDs {
        var user User
        err := s.cacheManager.Get(ctx, fmt.Sprintf("user:%s", userID), &user)
        if err == nil {
            cachedUsers[userID] = &user
        } else {
            missingIDs = append(missingIDs, userID)
        }
    }
    
    // Get missing users from database in batch
    if len(missingIDs) > 0 {
        dbUsers, err := s.getUsersFromDBBatch(ctx, missingIDs)
        if err != nil {
            return nil, err
        }
        
        // Cache the results
        for _, user := range dbUsers {
            s.cacheManager.Set(ctx, fmt.Sprintf("user:%s", user.ID), user, 1*time.Hour)
            cachedUsers[user.ID] = user
        }
    }
    
    // Return users in the same order as requested
    result := make([]*User, len(userIDs))
    for i, userID := range userIDs {
        result[i] = cachedUsers[userID]
    }
    
    return result, nil
}
```

---

**Library Combination - Build powerful services by combining multiple libraries! 🚀**

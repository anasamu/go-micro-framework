# Message Broker Communication

## 🎯 Overview

GoMicroFramework menyediakan dukungan lengkap untuk message broker communication antar services. Framework mendukung berbagai message brokers seperti Kafka, RabbitMQ, NATS, dan AWS SQS untuk asynchronous communication.

## 🔧 Supported Message Brokers

### 1. Apache Kafka
- High-throughput, distributed streaming platform
- Ideal for event streaming and real-time data processing
- Supports both pub/sub and queue patterns

### 2. RabbitMQ
- Reliable message broker with advanced routing
- Supports multiple messaging patterns
- Built-in clustering and high availability

### 3. NATS
- Lightweight, high-performance messaging system
- Simple pub/sub and request-reply patterns
- Cloud-native and scalable

### 4. AWS SQS
- Managed message queue service
- Fully managed and serverless
- Integrates with other AWS services

## 🔧 Message Broker Setup

### 1. Generate Service with Messaging

```bash
# Generate service with Kafka messaging
microframework new user-service --with-messaging=kafka --with-database=postgres

# Generate service with RabbitMQ messaging
microframework new order-service --with-messaging=rabbitmq --with-database=postgres

# Generate service with NATS messaging
microframework new notification-service --with-messaging=nats --with-ai=openai
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

# Message broker configuration
messaging:
  providers:
    kafka:
      enabled: true
      brokers: ["localhost:9092"]
      group_id: "user-service"
      topics: ["user-events", "user-commands"]
      consumer:
        auto_offset_reset: "latest"
        enable_auto_commit: true
        auto_commit_interval: "1s"
        session_timeout: "30s"
        heartbeat_interval: "3s"
      producer:
        required_acks: 1
        timeout: "10s"
        retry_max: 3
        compression: "gzip"
        
    rabbitmq:
      enabled: false
      url: "amqp://guest:guest@localhost:5672/"
      exchange: "user-exchange"
      queue: "user-queue"
      routing_key: "user.created"
      durable: true
      auto_delete: false
      exclusive: false
      no_wait: false
      
    nats:
      enabled: false
      url: "nats://localhost:4222"
      cluster_id: "user-service"
      client_id: "user-service-1"
      max_reconnect: 5
      reconnect_wait: "2s"
      timeout: "10s"
      
    sqs:
      enabled: false
      region: "us-east-1"
      access_key_id: "${AWS_ACCESS_KEY_ID}"
      secret_access_key: "${AWS_SECRET_ACCESS_KEY}"
      queue_url: "${SQS_QUEUE_URL}"
      visibility_timeout: "30s"
      message_retention_period: "14d"
```

## 🔧 Kafka Implementation

### 1. Kafka Producer

```go
// internal/messaging/kafka_producer.go
package messaging

import (
    "context"
    "encoding/json"
    "log"
    "time"
    
    "github.com/anasamu/go-micro-libs/messaging"
    "github.com/anasamu/go-micro-libs/messaging/providers/kafka"
)

type KafkaProducer struct {
    msgManager *messaging.MessagingManager
    topic      string
}

func NewKafkaProducer(msgManager *messaging.MessagingManager, topic string) *KafkaProducer {
    return &KafkaProducer{
        msgManager: msgManager,
        topic:      topic,
    }
}

func (p *KafkaProducer) PublishUserCreated(ctx context.Context, user *User) error {
    event := &UserCreatedEvent{
        UserID:    user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: time.Now(),
    }
    
    return p.publishEvent(ctx, "user.created", event)
}

func (p *KafkaProducer) PublishUserUpdated(ctx context.Context, user *User) error {
    event := &UserUpdatedEvent{
        UserID:    user.ID,
        Name:      user.Name,
        Email:     user.Email,
        UpdatedAt: time.Now(),
    }
    
    return p.publishEvent(ctx, "user.updated", event)
}

func (p *KafkaProducer) PublishUserDeleted(ctx context.Context, userID string) error {
    event := &UserDeletedEvent{
        UserID:    userID,
        DeletedAt: time.Now(),
    }
    
    return p.publishEvent(ctx, "user.deleted", event)
}

func (p *KafkaProducer) publishEvent(ctx context.Context, eventType string, event interface{}) error {
    // Serialize event
    eventData, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    // Create message
    message := messaging.CreateMessage(eventType, eventData)
    
    // Publish message
    publishReq := &messaging.PublishRequest{
        Topic:   p.topic,
        Message: message,
    }
    
    response, err := p.msgManager.PublishMessage(ctx, "kafka", publishReq)
    if err != nil {
        log.Printf("Failed to publish event %s: %v", eventType, err)
        return err
    }
    
    log.Printf("Published event %s with message ID: %s", eventType, response.MessageID)
    return nil
}
```

### 2. Kafka Consumer

```go
// internal/messaging/kafka_consumer.go
package messaging

import (
    "context"
    "encoding/json"
    "log"
    "sync"
    
    "github.com/anasamu/go-micro-libs/messaging"
    "github.com/anasamu/go-micro-libs/messaging/providers/kafka"
)

type KafkaConsumer struct {
    msgManager   *messaging.MessagingManager
    topic        string
    groupID      string
    handlers     map[string]MessageHandler
    stopChan     chan struct{}
    wg           sync.WaitGroup
}

type MessageHandler func(ctx context.Context, message *messaging.Message) error

func NewKafkaConsumer(msgManager *messaging.MessagingManager, topic, groupID string) *KafkaConsumer {
    return &KafkaConsumer{
        msgManager: msgManager,
        topic:      topic,
        groupID:    groupID,
        handlers:   make(map[string]MessageHandler),
        stopChan:   make(chan struct{}),
    }
}

func (c *KafkaConsumer) RegisterHandler(eventType string, handler MessageHandler) {
    c.handlers[eventType] = handler
}

func (c *KafkaConsumer) Start(ctx context.Context) error {
    // Subscribe to topic
    subscribeReq := &messaging.SubscribeRequest{
        Topic:   c.topic,
        GroupID: c.groupID,
    }
    
    messageChan, err := c.msgManager.Subscribe(ctx, "kafka", subscribeReq)
    if err != nil {
        return err
    }
    
    // Start consumer loop
    c.wg.Add(1)
    go c.consumeMessages(ctx, messageChan)
    
    return nil
}

func (c *KafkaConsumer) Stop() {
    close(c.stopChan)
    c.wg.Wait()
}

func (c *KafkaConsumer) consumeMessages(ctx context.Context, messageChan <-chan *messaging.Message) {
    defer c.wg.Done()
    
    for {
        select {
        case message := <-messageChan:
            if err := c.handleMessage(ctx, message); err != nil {
                log.Printf("Failed to handle message: %v", err)
            }
            
        case <-c.stopChan:
            return
            
        case <-ctx.Done():
            return
        }
    }
}

func (c *KafkaConsumer) handleMessage(ctx context.Context, message *messaging.Message) error {
    // Get handler for event type
    handler, exists := c.handlers[message.Type]
    if !exists {
        log.Printf("No handler found for event type: %s", message.Type)
        return nil
    }
    
    // Call handler
    return handler(ctx, message)
}
```

### 3. Event Handlers

```go
// internal/handlers/user_event_handlers.go
package handlers

import (
    "context"
    "encoding/json"
    "log"
    
    "github.com/anasamu/go-micro-libs/messaging"
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/cache"
)

type UserEventHandlers struct {
    dbManager    *database.DatabaseManager
    cacheManager *cache.CacheManager
}

func NewUserEventHandlers(dbManager *database.DatabaseManager, cacheManager *cache.CacheManager) *UserEventHandlers {
    return &UserEventHandlers{
        dbManager:    dbManager,
        cacheManager: cacheManager,
    }
}

func (h *UserEventHandlers) HandleUserCreated(ctx context.Context, message *messaging.Message) error {
    var event UserCreatedEvent
    if err := json.Unmarshal(message.Data, &event); err != nil {
        return err
    }
    
    log.Printf("Handling user created event: %+v", event)
    
    // Update cache
    cacheKey := fmt.Sprintf("user:%s", event.UserID)
    user := &User{
        ID:        event.UserID,
        Name:      event.Name,
        Email:     event.Email,
        CreatedAt: event.CreatedAt,
    }
    
    if err := h.cacheManager.Set(ctx, cacheKey, user, 1*time.Hour); err != nil {
        log.Printf("Failed to cache user: %v", err)
    }
    
    // Send welcome email (async)
    go h.sendWelcomeEmail(ctx, user)
    
    return nil
}

func (h *UserEventHandlers) HandleUserUpdated(ctx context.Context, message *messaging.Message) error {
    var event UserUpdatedEvent
    if err := json.Unmarshal(message.Data, &event); err != nil {
        return err
    }
    
    log.Printf("Handling user updated event: %+v", event)
    
    // Invalidate cache
    cacheKey := fmt.Sprintf("user:%s", event.UserID)
    if err := h.cacheManager.Delete(ctx, cacheKey); err != nil {
        log.Printf("Failed to invalidate cache: %v", err)
    }
    
    return nil
}

func (h *UserEventHandlers) HandleUserDeleted(ctx context.Context, message *messaging.Message) error {
    var event UserDeletedEvent
    if err := json.Unmarshal(message.Data, &event); err != nil {
        return err
    }
    
    log.Printf("Handling user deleted event: %+v", event)
    
    // Remove from cache
    cacheKey := fmt.Sprintf("user:%s", event.UserID)
    if err := h.cacheManager.Delete(ctx, cacheKey); err != nil {
        log.Printf("Failed to remove from cache: %v", err)
    }
    
    // Cleanup related data
    go h.cleanupUserData(ctx, event.UserID)
    
    return nil
}

func (h *UserEventHandlers) sendWelcomeEmail(ctx context.Context, user *User) {
    // Send welcome email logic
    log.Printf("Sending welcome email to user: %s", user.Email)
}

func (h *UserEventHandlers) cleanupUserData(ctx context.Context, userID string) {
    // Cleanup related data logic
    log.Printf("Cleaning up data for user: %s", userID)
}
```

## 🔧 RabbitMQ Implementation

### 1. RabbitMQ Producer

```go
// internal/messaging/rabbitmq_producer.go
package messaging

import (
    "context"
    "encoding/json"
    "log"
    "time"
    
    "github.com/anasamu/go-micro-libs/messaging"
    "github.com/anasamu/go-micro-libs/messaging/providers/rabbitmq"
)

type RabbitMQProducer struct {
    msgManager *messaging.MessagingManager
    exchange   string
    routingKey string
}

func NewRabbitMQProducer(msgManager *messaging.MessagingManager, exchange, routingKey string) *RabbitMQProducer {
    return &RabbitMQProducer{
        msgManager: msgManager,
        exchange:   exchange,
        routingKey: routingKey,
    }
}

func (p *RabbitMQProducer) PublishUserCreated(ctx context.Context, user *User) error {
    event := &UserCreatedEvent{
        UserID:    user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: time.Now(),
    }
    
    return p.publishEvent(ctx, "user.created", event)
}

func (p *RabbitMQProducer) publishEvent(ctx context.Context, eventType string, event interface{}) error {
    // Serialize event
    eventData, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    // Create message
    message := messaging.CreateMessage(eventType, eventData)
    
    // Publish message
    publishReq := &messaging.PublishRequest{
        Exchange:   p.exchange,
        RoutingKey: p.routingKey,
        Message:    message,
    }
    
    response, err := p.msgManager.PublishMessage(ctx, "rabbitmq", publishReq)
    if err != nil {
        log.Printf("Failed to publish event %s: %v", eventType, err)
        return err
    }
    
    log.Printf("Published event %s with message ID: %s", eventType, response.MessageID)
    return nil
}
```

### 2. RabbitMQ Consumer

```go
// internal/messaging/rabbitmq_consumer.go
package messaging

import (
    "context"
    "log"
    "sync"
    
    "github.com/anasamu/go-micro-libs/messaging"
)

type RabbitMQConsumer struct {
    msgManager *messaging.MessagingManager
    exchange   string
    queue      string
    handlers   map[string]MessageHandler
    stopChan   chan struct{}
    wg         sync.WaitGroup
}

func NewRabbitMQConsumer(msgManager *messaging.MessagingManager, exchange, queue string) *RabbitMQConsumer {
    return &RabbitMQConsumer{
        msgManager: msgManager,
        exchange:   exchange,
        queue:      queue,
        handlers:   make(map[string]MessageHandler),
        stopChan:   make(chan struct{}),
    }
}

func (c *RabbitMQConsumer) RegisterHandler(eventType string, handler MessageHandler) {
    c.handlers[eventType] = handler
}

func (c *RabbitMQConsumer) Start(ctx context.Context) error {
    // Subscribe to queue
    subscribeReq := &messaging.SubscribeRequest{
        Exchange: c.exchange,
        Queue:    c.queue,
    }
    
    messageChan, err := c.msgManager.Subscribe(ctx, "rabbitmq", subscribeReq)
    if err != nil {
        return err
    }
    
    // Start consumer loop
    c.wg.Add(1)
    go c.consumeMessages(ctx, messageChan)
    
    return nil
}

func (c *RabbitMQConsumer) Stop() {
    close(c.stopChan)
    c.wg.Wait()
}

func (c *RabbitMQConsumer) consumeMessages(ctx context.Context, messageChan <-chan *messaging.Message) {
    defer c.wg.Done()
    
    for {
        select {
        case message := <-messageChan:
            if err := c.handleMessage(ctx, message); err != nil {
                log.Printf("Failed to handle message: %v", err)
            }
            
        case <-c.stopChan:
            return
            
        case <-ctx.Done():
            return
        }
    }
}

func (c *RabbitMQConsumer) handleMessage(ctx context.Context, message *messaging.Message) error {
    // Get handler for event type
    handler, exists := c.handlers[message.Type]
    if !exists {
        log.Printf("No handler found for event type: %s", message.Type)
        return nil
    }
    
    // Call handler
    return handler(ctx, message)
}
```

## 🔧 NATS Implementation

### 1. NATS Producer

```go
// internal/messaging/nats_producer.go
package messaging

import (
    "context"
    "encoding/json"
    "log"
    "time"
    
    "github.com/anasamu/go-micro-libs/messaging"
)

type NATSProducer struct {
    msgManager *messaging.MessagingManager
    subject    string
}

func NewNATSProducer(msgManager *messaging.MessagingManager, subject string) *NATSProducer {
    return &NATSProducer{
        msgManager: msgManager,
        subject:    subject,
    }
}

func (p *NATSProducer) PublishUserCreated(ctx context.Context, user *User) error {
    event := &UserCreatedEvent{
        UserID:    user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: time.Now(),
    }
    
    return p.publishEvent(ctx, "user.created", event)
}

func (p *NATSProducer) publishEvent(ctx context.Context, eventType string, event interface{}) error {
    // Serialize event
    eventData, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    // Create message
    message := messaging.CreateMessage(eventType, eventData)
    
    // Publish message
    publishReq := &messaging.PublishRequest{
        Subject: p.subject,
        Message: message,
    }
    
    response, err := p.msgManager.PublishMessage(ctx, "nats", publishReq)
    if err != nil {
        log.Printf("Failed to publish event %s: %v", eventType, err)
        return err
    }
    
    log.Printf("Published event %s with message ID: %s", eventType, response.MessageID)
    return nil
}
```

### 2. NATS Consumer

```go
// internal/messaging/nats_consumer.go
package messaging

import (
    "context"
    "log"
    "sync"
    
    "github.com/anasamu/go-micro-libs/messaging"
)

type NATSConsumer struct {
    msgManager *messaging.MessagingManager
    subject    string
    handlers   map[string]MessageHandler
    stopChan   chan struct{}
    wg         sync.WaitGroup
}

func NewNATSConsumer(msgManager *messaging.MessagingManager, subject string) *NATSConsumer {
    return &NATSConsumer{
        msgManager: msgManager,
        subject:    subject,
        handlers:   make(map[string]MessageHandler),
        stopChan:   make(chan struct{}),
    }
}

func (c *NATSConsumer) RegisterHandler(eventType string, handler MessageHandler) {
    c.handlers[eventType] = handler
}

func (c *NATSConsumer) Start(ctx context.Context) error {
    // Subscribe to subject
    subscribeReq := &messaging.SubscribeRequest{
        Subject: c.subject,
    }
    
    messageChan, err := c.msgManager.Subscribe(ctx, "nats", subscribeReq)
    if err != nil {
        return err
    }
    
    // Start consumer loop
    c.wg.Add(1)
    go c.consumeMessages(ctx, messageChan)
    
    return nil
}

func (c *NATSConsumer) Stop() {
    close(c.stopChan)
    c.wg.Wait()
}

func (c *NATSConsumer) consumeMessages(ctx context.Context, messageChan <-chan *messaging.Message) {
    defer c.wg.Done()
    
    for {
        select {
        case message := <-messageChan:
            if err := c.handleMessage(ctx, message); err != nil {
                log.Printf("Failed to handle message: %v", err)
            }
            
        case <-c.stopChan:
            return
            
        case <-ctx.Done():
            return
        }
    }
}

func (c *NATSConsumer) handleMessage(ctx context.Context, message *messaging.Message) error {
    // Get handler for event type
    handler, exists := c.handlers[message.Type]
    if !exists {
        log.Printf("No handler found for event type: %s", message.Type)
        return nil
    }
    
    // Call handler
    return handler(ctx, message)
}
```

## 🔧 Service Integration

### 1. Service with Message Broker

```go
// cmd/main.go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/anasamu/go-micro-libs/config"
    "github.com/anasamu/go-micro-libs/logging"
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/messaging"
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/cache"
    
    "user-service/internal/messaging"
    "user-service/internal/handlers"
    "user-service/internal/services"
)

func main() {
    ctx := context.Background()
    
    // Initialize core managers
    configManager := config.NewManager()
    loggingManager := logging.NewManager()
    monitoringManager := monitoring.NewManager()
    msgManager := messaging.NewManager()
    dbManager := database.NewManager()
    cacheManager := cache.NewManager()
    
    // Bootstrap service
    if err := bootstrapService(ctx, configManager, loggingManager, 
        monitoringManager, msgManager, dbManager, cacheManager); err != nil {
        log.Fatal("Failed to bootstrap service:", err)
    }
    
    // Start message consumers
    if err := startMessageConsumers(ctx, msgManager, dbManager, cacheManager); err != nil {
        log.Fatal("Failed to start message consumers:", err)
    }
    
    // Wait for shutdown signal
    waitForShutdown()
}

func startMessageConsumers(ctx context.Context, msgManager *messaging.MessagingManager, 
    dbManager *database.DatabaseManager, cacheManager *cache.CacheManager) error {
    
    // Create event handlers
    eventHandlers := handlers.NewUserEventHandlers(dbManager, cacheManager)
    
    // Create Kafka consumer
    kafkaConsumer := messaging.NewKafkaConsumer(msgManager, "user-events", "user-service")
    kafkaConsumer.RegisterHandler("user.created", eventHandlers.HandleUserCreated)
    kafkaConsumer.RegisterHandler("user.updated", eventHandlers.HandleUserUpdated)
    kafkaConsumer.RegisterHandler("user.deleted", eventHandlers.HandleUserDeleted)
    
    // Start consumer
    if err := kafkaConsumer.Start(ctx); err != nil {
        return err
    }
    
    return nil
}

func waitForShutdown() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    <-sigChan
    log.Println("Shutting down service...")
}
```

### 2. Service with Producer

```go
// internal/services/user_service.go
package services

import (
    "context"
    "time"
    
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/messaging"
)

type UserService struct {
    dbManager    *database.DatabaseManager
    msgManager   *messaging.MessagingManager
    producer     *messaging.KafkaProducer
}

func NewUserService(dbManager *database.DatabaseManager, msgManager *messaging.MessagingManager) *UserService {
    producer := messaging.NewKafkaProducer(msgManager, "user-events")
    
    return &UserService{
        dbManager:  dbManager,
        msgManager: msgManager,
        producer:   producer,
    }
}

func (s *UserService) CreateUser(ctx context.Context, user *User) (*User, error) {
    // Save user to database
    user.ID = generateID()
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    
    err := s.saveUser(ctx, user)
    if err != nil {
        return nil, err
    }
    
    // Publish user created event
    if err := s.producer.PublishUserCreated(ctx, user); err != nil {
        log.Printf("Failed to publish user created event: %v", err)
        // Don't fail the operation if event publishing fails
    }
    
    return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *User) (*User, error) {
    // Update user in database
    user.UpdatedAt = time.Now()
    
    err := s.updateUser(ctx, user)
    if err != nil {
        return nil, err
    }
    
    // Publish user updated event
    if err := s.producer.PublishUserUpdated(ctx, user); err != nil {
        log.Printf("Failed to publish user updated event: %v", err)
    }
    
    return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
    // Delete user from database
    err := s.deleteUser(ctx, userID)
    if err != nil {
        return err
    }
    
    // Publish user deleted event
    if err := s.producer.PublishUserDeleted(ctx, userID); err != nil {
        log.Printf("Failed to publish user deleted event: %v", err)
    }
    
    return nil
}

func (s *UserService) saveUser(ctx context.Context, user *User) error {
    query := `INSERT INTO users (id, name, email, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5)`
    
    _, err := s.dbManager.Exec(ctx, "postgresql", query, 
        user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
    
    return err
}

func (s *UserService) updateUser(ctx context.Context, user *User) error {
    query := `UPDATE users SET name = $1, email = $2, updated_at = $3 WHERE id = $4`
    
    _, err := s.dbManager.Exec(ctx, "postgresql", query, 
        user.Name, user.Email, user.UpdatedAt, user.ID)
    
    return err
}

func (s *UserService) deleteUser(ctx context.Context, userID string) error {
    query := `DELETE FROM users WHERE id = $1`
    
    _, err := s.dbManager.Exec(ctx, "postgresql", query, userID)
    
    return err
}
```

## 🔧 Event Sourcing

### 1. Event Store

```go
// internal/events/event_store.go
package events

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/anasamu/go-micro-libs/database"
    "github.com/anasamu/go-micro-libs/messaging"
)

type EventStore struct {
    dbManager  *database.DatabaseManager
    msgManager *messaging.MessagingManager
}

type Event struct {
    ID        string    `json:"id"`
    StreamID  string    `json:"stream_id"`
    Type      string    `json:"type"`
    Data      []byte    `json:"data"`
    Version   int       `json:"version"`
    CreatedAt time.Time `json:"created_at"`
}

func NewEventStore(dbManager *database.DatabaseManager, msgManager *messaging.MessagingManager) *EventStore {
    return &EventStore{
        dbManager:  dbManager,
        msgManager: msgManager,
    }
}

func (es *EventStore) AppendEvent(ctx context.Context, streamID string, eventType string, data interface{}) error {
    // Serialize event data
    eventData, err := json.Marshal(data)
    if err != nil {
        return err
    }
    
    // Get next version
    version, err := es.getNextVersion(ctx, streamID)
    if err != nil {
        return err
    }
    
    // Create event
    event := &Event{
        ID:        generateID(),
        StreamID:  streamID,
        Type:      eventType,
        Data:      eventData,
        Version:   version,
        CreatedAt: time.Now(),
    }
    
    // Save event to database
    err = es.saveEvent(ctx, event)
    if err != nil {
        return err
    }
    
    // Publish event to message broker
    message := messaging.CreateMessage(eventType, eventData)
    publishReq := &messaging.PublishRequest{
        Topic:   "events",
        Message: message,
    }
    
    _, err = es.msgManager.PublishMessage(ctx, "kafka", publishReq)
    if err != nil {
        log.Printf("Failed to publish event: %v", err)
    }
    
    return nil
}

func (es *EventStore) GetEvents(ctx context.Context, streamID string, fromVersion int) ([]*Event, error) {
    query := `SELECT id, stream_id, type, data, version, created_at 
              FROM events 
              WHERE stream_id = $1 AND version >= $2 
              ORDER BY version ASC`
    
    result, err := es.dbManager.Query(ctx, "postgresql", query, streamID, fromVersion)
    if err != nil {
        return nil, err
    }
    
    var events []*Event
    for result.Next() {
        event := &Event{}
        err = result.Scan(&event.ID, &event.StreamID, &event.Type, 
            &event.Data, &event.Version, &event.CreatedAt)
        if err != nil {
            return nil, err
        }
        events = append(events, event)
    }
    
    return events, nil
}

func (es *EventStore) getNextVersion(ctx context.Context, streamID string) (int, error) {
    query := `SELECT COALESCE(MAX(version), 0) + 1 FROM events WHERE stream_id = $1`
    
    result, err := es.dbManager.Query(ctx, "postgresql", query, streamID)
    if err != nil {
        return 0, err
    }
    
    var version int
    if result.Next() {
        err = result.Scan(&version)
        if err != nil {
            return 0, err
        }
    }
    
    return version, nil
}

func (es *EventStore) saveEvent(ctx context.Context, event *Event) error {
    query := `INSERT INTO events (id, stream_id, type, data, version, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6)`
    
    _, err := es.dbManager.Exec(ctx, "postgresql", query, 
        event.ID, event.StreamID, event.Type, event.Data, event.Version, event.CreatedAt)
    
    return err
}
```

## 🔧 Best Practices

### 1. Message Serialization

```go
// Use consistent message format
type Message struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    Data      map[string]interface{} `json:"data"`
    Metadata  map[string]string      `json:"metadata"`
    Timestamp time.Time              `json:"timestamp"`
    Version   string                 `json:"version"`
}

func CreateMessage(eventType string, data interface{}) *Message {
    return &Message{
        ID:        generateID(),
        Type:      eventType,
        Data:      data,
        Metadata:  make(map[string]string),
        Timestamp: time.Now(),
        Version:   "1.0",
    }
}
```

### 2. Error Handling

```go
// Proper error handling for message processing
func (h *UserEventHandlers) HandleUserCreated(ctx context.Context, message *messaging.Message) error {
    var event UserCreatedEvent
    if err := json.Unmarshal(message.Data, &event); err != nil {
        log.Printf("Failed to unmarshal event: %v", err)
        return err
    }
    
    // Process event with retry logic
    return h.processWithRetry(ctx, func() error {
        return h.processUserCreated(ctx, &event)
    })
}

func (h *UserEventHandlers) processWithRetry(ctx context.Context, fn func() error) error {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        if err := fn(); err != nil {
            if i == maxRetries-1 {
                return err
            }
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }
        return nil
    }
    return nil
}
```

### 3. Dead Letter Queue

```go
// Dead letter queue for failed messages
func (c *KafkaConsumer) handleMessageWithDLQ(ctx context.Context, message *messaging.Message) error {
    handler, exists := c.handlers[message.Type]
    if !exists {
        return nil
    }
    
    err := handler(ctx, message)
    if err != nil {
        // Send to dead letter queue
        dlqMessage := &messaging.Message{
            ID:        message.ID,
            Type:      "dlq." + message.Type,
            Data:      message.Data,
            Metadata:  message.Metadata,
            Timestamp: time.Now(),
        }
        
        publishReq := &messaging.PublishRequest{
            Topic:   "dlq",
            Message: dlqMessage,
        }
        
        _, dlqErr := c.msgManager.PublishMessage(ctx, "kafka", publishReq)
        if dlqErr != nil {
            log.Printf("Failed to send to DLQ: %v", dlqErr)
        }
    }
    
    return err
}
```

### 4. Message Ordering

```go
// Ensure message ordering for critical events
func (p *KafkaProducer) PublishOrderedEvent(ctx context.Context, eventType string, event interface{}, key string) error {
    eventData, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    message := messaging.CreateMessage(eventType, eventData)
    
    publishReq := &messaging.PublishRequest{
        Topic:   p.topic,
        Message: message,
        Key:     key, // Use key for partitioning
    }
    
    _, err = p.msgManager.PublishMessage(ctx, "kafka", publishReq)
    return err
}
```

---

**Message Broker Communication - Reliable asynchronous communication between services! 🚀**

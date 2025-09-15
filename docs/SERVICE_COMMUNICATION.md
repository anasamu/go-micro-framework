# Service Communication Guide

## 🎯 Overview

GoMicroFramework menyediakan berbagai protokol komunikasi untuk interaksi antar services. Framework mendukung HTTP/REST, gRPC, WebSocket, GraphQL, dan protokol komunikasi lainnya dengan konfigurasi yang mudah.

## 🔧 Communication Protocols

### 1. HTTP/REST Communication

#### Basic HTTP Service

```go
// internal/handlers/user_handler.go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/anasamu/go-micro-libs/communication"
)

type UserHandler struct {
    userService UserService
    commManager *communication.CommunicationManager
}

func NewUserHandler(userService UserService, commManager *communication.CommunicationManager) *UserHandler {
    return &UserHandler{
        userService: userService,
        commManager: commManager,
    }
}

func (h *UserHandler) GetUser(c *gin.Context) {
    userID := c.Param("id")
    
    user, err := h.userService.GetUser(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    createdUser, err := h.userService.CreateUser(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, createdUser)
}
```

#### HTTP Client

```go
// internal/clients/user_client.go
package clients

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    
    "github.com/anasamu/go-micro-libs/communication"
)

type UserClient struct {
    baseURL    string
    httpClient *http.Client
    commManager *communication.CommunicationManager
}

func NewUserClient(baseURL string, commManager *communication.CommunicationManager) *UserClient {
    return &UserClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        commManager: commManager,
    }
}

func (c *UserClient) GetUser(ctx context.Context, userID string) (*User, error) {
    url := fmt.Sprintf("%s/users/%s", c.baseURL, userID)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
    }
    
    var user User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, err
    }
    
    return &user, nil
}

func (c *UserClient) CreateUser(ctx context.Context, user *User) (*User, error) {
    url := fmt.Sprintf("%s/users", c.baseURL)
    
    jsonData, err := json.Marshal(user)
    if err != nil {
        return nil, err
    }
    
    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusCreated {
        return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
    }
    
    var createdUser User
    if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
        return nil, err
    }
    
    return &createdUser, nil
}
```

### 2. gRPC Communication

#### gRPC Service

```go
// internal/grpc/user_service.go
package grpc

import (
    "context"
    
    "google.golang.org/grpc"
    "github.com/anasamu/go-micro-libs/communication"
    pb "user-service/pkg/proto"
)

type UserGRPCService struct {
    pb.UnimplementedUserServiceServer
    userService UserService
    commManager *communication.CommunicationManager
}

func NewUserGRPCService(userService UserService, commManager *communication.CommunicationManager) *UserGRPCService {
    return &UserGRPCService{
        userService: userService,
        commManager: commManager,
    }
}

func (s *UserGRPCService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    user, err := s.userService.GetUser(req.UserId)
    if err != nil {
        return nil, err
    }
    
    return &pb.GetUserResponse{
        User: &pb.User{
            Id:    user.ID,
            Name:  user.Name,
            Email: user.Email,
        },
    }, nil
}

func (s *UserGRPCService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    user := &User{
        Name:  req.Name,
        Email: req.Email,
    }
    
    createdUser, err := s.userService.CreateUser(*user)
    if err != nil {
        return nil, err
    }
    
    return &pb.CreateUserResponse{
        User: &pb.User{
            Id:    createdUser.ID,
            Name:  createdUser.Name,
            Email: createdUser.Email,
        },
    }, nil
}
```

#### gRPC Client

```go
// internal/clients/user_grpc_client.go
package clients

import (
    "context"
    
    "google.golang.org/grpc"
    "github.com/anasamu/go-micro-libs/communication"
    pb "user-service/pkg/proto"
)

type UserGRPCClient struct {
    client      pb.UserServiceClient
    conn        *grpc.ClientConn
    commManager *communication.CommunicationManager
}

func NewUserGRPCClient(address string, commManager *communication.CommunicationManager) (*UserGRPCClient, error) {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    
    client := pb.NewUserServiceClient(conn)
    
    return &UserGRPCClient{
        client:      client,
        conn:        conn,
        commManager: commManager,
    }, nil
}

func (c *UserGRPCClient) GetUser(ctx context.Context, userID string) (*User, error) {
    req := &pb.GetUserRequest{
        UserId: userID,
    }
    
    resp, err := c.client.GetUser(ctx, req)
    if err != nil {
        return nil, err
    }
    
    return &User{
        ID:    resp.User.Id,
        Name:  resp.User.Name,
        Email: resp.User.Email,
    }, nil
}

func (c *UserGRPCClient) CreateUser(ctx context.Context, user *User) (*User, error) {
    req := &pb.CreateUserRequest{
        Name:  user.Name,
        Email: user.Email,
    }
    
    resp, err := c.client.CreateUser(ctx, req)
    if err != nil {
        return nil, err
    }
    
    return &User{
        ID:    resp.User.Id,
        Name:  resp.User.Name,
        Email: resp.User.Email,
    }, nil
}

func (c *UserGRPCClient) Close() error {
    return c.conn.Close()
}
```

### 3. WebSocket Communication

#### WebSocket Service

```go
// internal/websocket/user_websocket.go
package websocket

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    
    "github.com/gorilla/websocket"
    "github.com/anasamu/go-micro-libs/communication"
)

type UserWebSocketService struct {
    upgrader    websocket.Upgrader
    userService UserService
    commManager *communication.CommunicationManager
    clients     map[*websocket.Conn]bool
    broadcast   chan []byte
}

func NewUserWebSocketService(userService UserService, commManager *communication.CommunicationManager) *UserWebSocketService {
    return &UserWebSocketService{
        upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool {
                return true
            },
        },
        userService: userService,
        commManager: commManager,
        clients:     make(map[*websocket.Conn]bool),
        broadcast:   make(chan []byte),
    }
}

func (s *UserWebSocketService) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := s.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("WebSocket upgrade error:", err)
        return
    }
    defer conn.Close()
    
    s.clients[conn] = true
    
    for {
        var msg map[string]interface{}
        err := conn.ReadJSON(&msg)
        if err != nil {
            log.Println("WebSocket read error:", err)
            delete(s.clients, conn)
            break
        }
        
        // Handle message
        s.handleMessage(conn, msg)
    }
}

func (s *UserWebSocketService) handleMessage(conn *websocket.Conn, msg map[string]interface{}) {
    msgType, ok := msg["type"].(string)
    if !ok {
        return
    }
    
    switch msgType {
    case "get_user":
        userID, ok := msg["user_id"].(string)
        if !ok {
            return
        }
        
        user, err := s.userService.GetUser(userID)
        if err != nil {
            s.sendError(conn, err.Error())
            return
        }
        
        response := map[string]interface{}{
            "type": "user_response",
            "data": user,
        }
        
        s.sendMessage(conn, response)
        
    case "create_user":
        userData, ok := msg["data"].(map[string]interface{})
        if !ok {
            return
        }
        
        user := &User{
            Name:  userData["name"].(string),
            Email: userData["email"].(string),
        }
        
        createdUser, err := s.userService.CreateUser(*user)
        if err != nil {
            s.sendError(conn, err.Error())
            return
        }
        
        response := map[string]interface{}{
            "type": "user_created",
            "data": createdUser,
        }
        
        s.sendMessage(conn, response)
    }
}

func (s *UserWebSocketService) sendMessage(conn *websocket.Conn, msg map[string]interface{}) {
    err := conn.WriteJSON(msg)
    if err != nil {
        log.Println("WebSocket write error:", err)
        delete(s.clients, conn)
    }
}

func (s *UserWebSocketService) sendError(conn *websocket.Conn, errorMsg string) {
    response := map[string]interface{}{
        "type":  "error",
        "error": errorMsg,
    }
    
    s.sendMessage(conn, response)
}

func (s *UserWebSocketService) Broadcast(message []byte) {
    for client := range s.clients {
        err := client.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            log.Println("WebSocket broadcast error:", err)
            client.Close()
            delete(s.clients, client)
        }
    }
}
```

### 4. GraphQL Communication

#### GraphQL Service

```go
// internal/graphql/user_resolver.go
package graphql

import (
    "context"
    
    "github.com/graphql-go/graphql"
    "github.com/anasamu/go-micro-libs/communication"
)

type UserResolver struct {
    userService UserService
    commManager *communication.CommunicationManager
}

func NewUserResolver(userService UserService, commManager *communication.CommunicationManager) *UserResolver {
    return &UserResolver{
        userService: userService,
        commManager: commManager,
    }
}

func (r *UserResolver) GetUser(p graphql.ResolveParams) (interface{}, error) {
    userID, ok := p.Args["id"].(string)
    if !ok {
        return nil, nil
    }
    
    user, err := r.userService.GetUser(userID)
    if err != nil {
        return nil, err
    }
    
    return user, nil
}

func (r *UserResolver) CreateUser(p graphql.ResolveParams) (interface{}, error) {
    input, ok := p.Args["input"].(map[string]interface{})
    if !ok {
        return nil, nil
    }
    
    user := &User{
        Name:  input["name"].(string),
        Email: input["email"].(string),
    }
    
    createdUser, err := r.userService.CreateUser(*user)
    if err != nil {
        return nil, err
    }
    
    return createdUser, nil
}

func (r *UserResolver) GetUsers(p graphql.ResolveParams) (interface{}, error) {
    users, err := r.userService.GetUsers()
    if err != nil {
        return nil, err
    }
    
    return users, nil
}
```

#### GraphQL Schema

```go
// internal/graphql/schema.go
package graphql

import (
    "github.com/graphql-go/graphql"
)

func NewUserSchema(userResolver *UserResolver) *graphql.Schema {
    userType := graphql.NewObject(graphql.ObjectConfig{
        Name: "User",
        Fields: graphql.Fields{
            "id": &graphql.Field{
                Type: graphql.String,
            },
            "name": &graphql.Field{
                Type: graphql.String,
            },
            "email": &graphql.Field{
                Type: graphql.String,
            },
            "created_at": &graphql.Field{
                Type: graphql.DateTime,
            },
        },
    })
    
    createUserInputType := graphql.NewInputObject(graphql.InputObjectConfig{
        Name: "CreateUserInput",
        Fields: graphql.InputObjectConfigFieldMap{
            "name": &graphql.InputObjectFieldConfig{
                Type: graphql.NewNonNull(graphql.String),
            },
            "email": &graphql.InputObjectFieldConfig{
                Type: graphql.NewNonNull(graphql.String),
            },
        },
    })
    
    queryType := graphql.NewObject(graphql.ObjectConfig{
        Name: "Query",
        Fields: graphql.Fields{
            "user": &graphql.Field{
                Type:        userType,
                Description: "Get user by ID",
                Args: graphql.FieldConfigArgument{
                    "id": &graphql.ArgumentConfig{
                        Type: graphql.NewNonNull(graphql.String),
                    },
                },
                Resolve: userResolver.GetUser,
            },
            "users": &graphql.Field{
                Type:        graphql.NewList(userType),
                Description: "Get all users",
                Resolve:     userResolver.GetUsers,
            },
        },
    })
    
    mutationType := graphql.NewObject(graphql.ObjectConfig{
        Name: "Mutation",
        Fields: graphql.Fields{
            "createUser": &graphql.Field{
                Type:        userType,
                Description: "Create new user",
                Args: graphql.FieldConfigArgument{
                    "input": &graphql.ArgumentConfig{
                        Type: graphql.NewNonNull(createUserInputType),
                    },
                },
                Resolve: userResolver.CreateUser,
            },
        },
    })
    
    schema, _ := graphql.NewSchema(graphql.SchemaConfig{
        Query:    queryType,
        Mutation: mutationType,
    })
    
    return &schema
}
```

## 🔧 Service Discovery Integration

### 1. Consul Service Discovery

```go
// internal/discovery/consul_discovery.go
package discovery

import (
    "context"
    "fmt"
    "time"
    
    "github.com/hashicorp/consul/api"
    "github.com/anasamu/go-micro-libs/discovery"
)

type ConsulDiscovery struct {
    client      *api.Client
    serviceName string
    serviceID   string
    servicePort int
    commManager *communication.CommunicationManager
}

func NewConsulDiscovery(serviceName string, servicePort int, commManager *communication.CommunicationManager) (*ConsulDiscovery, error) {
    client, err := api.NewClient(api.DefaultConfig())
    if err != nil {
        return nil, err
    }
    
    return &ConsulDiscovery{
        client:      client,
        serviceName: serviceName,
        servicePort: servicePort,
        commManager: commManager,
    }, nil
}

func (d *ConsulDiscovery) RegisterService(ctx context.Context) error {
    serviceID := fmt.Sprintf("%s-%d", d.serviceName, time.Now().Unix())
    
    registration := &api.AgentServiceRegistration{
        ID:      serviceID,
        Name:    d.serviceName,
        Port:    d.servicePort,
        Address: "localhost",
        Check: &api.AgentServiceCheck{
            HTTP:                           fmt.Sprintf("http://localhost:%d/health", d.servicePort),
            Timeout:                        "3s",
            Interval:                       "10s",
            DeregisterCriticalServiceAfter: "30s",
        },
    }
    
    err := d.client.Agent().ServiceRegister(registration)
    if err != nil {
        return err
    }
    
    d.serviceID = serviceID
    return nil
}

func (d *ConsulDiscovery) DiscoverServices(ctx context.Context, serviceName string) ([]*discovery.ServiceInfo, error) {
    services, _, err := d.client.Health().Service(serviceName, "", true, nil)
    if err != nil {
        return nil, err
    }
    
    var serviceInfos []*discovery.ServiceInfo
    for _, service := range services {
        serviceInfo := &discovery.ServiceInfo{
            ID:      service.Service.ID,
            Name:    service.Service.Service,
            Address: service.Service.Address,
            Port:    service.Service.Port,
            Tags:    service.Service.Tags,
        }
        serviceInfos = append(serviceInfos, serviceInfo)
    }
    
    return serviceInfos, nil
}

func (d *ConsulDiscovery) DeregisterService(ctx context.Context) error {
    if d.serviceID == "" {
        return nil
    }
    
    return d.client.Agent().ServiceDeregister(d.serviceID)
}
```

### 2. Service Client with Discovery

```go
// internal/clients/discovery_client.go
package clients

import (
    "context"
    "fmt"
    "math/rand"
    "time"
    
    "github.com/anasamu/go-micro-libs/discovery"
    "github.com/anasamu/go-micro-libs/communication"
)

type DiscoveryClient struct {
    discoveryManager *discovery.DiscoveryManager
    commManager      *communication.CommunicationManager
    serviceName      string
}

func NewDiscoveryClient(discoveryManager *discovery.DiscoveryManager, commManager *communication.CommunicationManager, serviceName string) *DiscoveryClient {
    return &DiscoveryClient{
        discoveryManager: discoveryManager,
        commManager:      commManager,
        serviceName:      serviceName,
    }
}

func (c *DiscoveryClient) GetServiceURL(ctx context.Context) (string, error) {
    services, err := c.discoveryManager.DiscoverServices(ctx, c.serviceName)
    if err != nil {
        return "", err
    }
    
    if len(services) == 0 {
        return "", fmt.Errorf("no services found for %s", c.serviceName)
    }
    
    // Simple load balancing - random selection
    service := services[rand.Intn(len(services))]
    return fmt.Sprintf("http://%s:%d", service.Address, service.Port), nil
}

func (c *DiscoveryClient) GetUser(ctx context.Context, userID string) (*User, error) {
    serviceURL, err := c.GetServiceURL(ctx)
    if err != nil {
        return nil, err
    }
    
    // Create HTTP client for the discovered service
    client := NewUserClient(serviceURL, c.commManager)
    return client.GetUser(ctx, userID)
}
```

## 🔧 Load Balancing

### 1. Round Robin Load Balancer

```go
// internal/loadbalancer/round_robin.go
package loadbalancer

import (
    "sync"
    "github.com/anasamu/go-micro-libs/discovery"
)

type RoundRobinLoadBalancer struct {
    services []*discovery.ServiceInfo
    current  int
    mutex    sync.Mutex
}

func NewRoundRobinLoadBalancer() *RoundRobinLoadBalancer {
    return &RoundRobinLoadBalancer{
        services: make([]*discovery.ServiceInfo, 0),
        current:  0,
    }
}

func (lb *RoundRobinLoadBalancer) SetServices(services []*discovery.ServiceInfo) {
    lb.mutex.Lock()
    defer lb.mutex.Unlock()
    
    lb.services = services
    lb.current = 0
}

func (lb *RoundRobinLoadBalancer) GetNextService() *discovery.ServiceInfo {
    lb.mutex.Lock()
    defer lb.mutex.Unlock()
    
    if len(lb.services) == 0 {
        return nil
    }
    
    service := lb.services[lb.current]
    lb.current = (lb.current + 1) % len(lb.services)
    
    return service
}
```

### 2. Weighted Load Balancer

```go
// internal/loadbalancer/weighted.go
package loadbalancer

import (
    "math/rand"
    "sync"
    "time"
    
    "github.com/anasamu/go-micro-libs/discovery"
)

type WeightedService struct {
    Service *discovery.ServiceInfo
    Weight  int
}

type WeightedLoadBalancer struct {
    services []WeightedService
    totalWeight int
    mutex    sync.RWMutex
}

func NewWeightedLoadBalancer() *WeightedLoadBalancer {
    rand.Seed(time.Now().UnixNano())
    return &WeightedLoadBalancer{
        services: make([]WeightedService, 0),
    }
}

func (lb *WeightedLoadBalancer) SetServices(services []WeightedService) {
    lb.mutex.Lock()
    defer lb.mutex.Unlock()
    
    lb.services = services
    lb.totalWeight = 0
    
    for _, service := range services {
        lb.totalWeight += service.Weight
    }
}

func (lb *WeightedLoadBalancer) GetNextService() *discovery.ServiceInfo {
    lb.mutex.RLock()
    defer lb.mutex.RUnlock()
    
    if len(lb.services) == 0 {
        return nil
    }
    
    random := rand.Intn(lb.totalWeight)
    current := 0
    
    for _, service := range lb.services {
        current += service.Weight
        if random < current {
            return service.Service
        }
    }
    
    return lb.services[0].Service
}
```

## 🔧 Configuration

### Communication Configuration

```yaml
# config.yaml
communication:
  providers:
    http:
      enabled: true
      port: 8080
      timeout: "30s"
      read_timeout: "10s"
      write_timeout: "10s"
      idle_timeout: "120s"
      max_header_bytes: 1048576
      
    grpc:
      enabled: true
      port: 9090
      timeout: "30s"
      max_recv_msg_size: 4194304
      max_send_msg_size: 4194304
      keepalive:
        time: "30s"
        timeout: "5s"
        permit_without_stream: true
        
    websocket:
      enabled: false
      port: 8082
      path: "/ws"
      read_buffer_size: 1024
      write_buffer_size: 1024
      check_origin: true
      
    graphql:
      enabled: false
      port: 8083
      path: "/graphql"
      playground: true
      introspection: true

# Service discovery
discovery:
  providers:
    consul:
      enabled: true
      address: "localhost:8500"
      token: "${CONSUL_TOKEN}"
      service_name: "user-service"
      service_port: 8080
      health_check:
        enabled: true
        path: "/health"
        interval: "10s"
        timeout: "3s"
        deregister_after: "30s"
```

## 🔧 Best Practices

### 1. Error Handling

```go
// Proper error handling for service communication
func (c *UserClient) GetUserWithRetry(ctx context.Context, userID string, maxRetries int) (*User, error) {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        user, err := c.GetUser(ctx, userID)
        if err == nil {
            return user, nil
        }
        
        lastErr = err
        
        // Exponential backoff
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    
    return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
```

### 2. Circuit Breaker Integration

```go
// Circuit breaker for service communication
func (c *UserClient) GetUserWithCircuitBreaker(ctx context.Context, userID string) (*User, error) {
    result, err := c.circuitBreaker.Execute(func() (interface{}, error) {
        return c.GetUser(ctx, userID)
    })
    
    if err != nil {
        return nil, err
    }
    
    return result.(*User), nil
}
```

### 3. Timeout Management

```go
// Timeout management for service communication
func (c *UserClient) GetUserWithTimeout(ctx context.Context, userID string, timeout time.Duration) (*User, error) {
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()
    
    return c.GetUser(ctx, userID)
}
```

### 4. Health Checks

```go
// Health check for service communication
func (c *UserClient) HealthCheck(ctx context.Context) error {
    url := fmt.Sprintf("%s/health", c.baseURL)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return err
    }
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("health check failed: %d", resp.StatusCode)
    }
    
    return nil
}
```

---

**Service Communication - Powerful and flexible inter-service communication! 🚀**

# gRPC Inter-Service Communication

## 🎯 Overview

GoMicroFramework menyediakan dukungan lengkap untuk gRPC communication antar services. gRPC adalah protokol RPC yang high-performance, language-agnostic, dan ideal untuk microservices communication.

## 🔧 gRPC Service Setup

### 1. Generate gRPC Service

```bash
# Generate gRPC service
microframework new user-service --type=grpc --with-database=postgres --with-monitoring=prometheus
```

### 2. Protocol Buffer Definition

```protobuf
// pkg/proto/user.proto
syntax = "proto3";

package user;

option go_package = "github.com/user-service/pkg/proto";

// User service definition
service UserService {
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
    rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
    rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
    rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
    rpc StreamUsers(StreamUsersRequest) returns (stream User);
}

// User message
message User {
    string id = 1;
    string name = 2;
    string email = 3;
    int64 created_at = 4;
    int64 updated_at = 5;
}

// Request/Response messages
message GetUserRequest {
    string user_id = 1;
}

message GetUserResponse {
    User user = 1;
}

message CreateUserRequest {
    string name = 1;
    string email = 2;
}

message CreateUserResponse {
    User user = 1;
}

message UpdateUserRequest {
    string user_id = 1;
    string name = 2;
    string email = 3;
}

message UpdateUserResponse {
    User user = 1;
}

message DeleteUserRequest {
    string user_id = 1;
}

message DeleteUserResponse {
    bool success = 1;
}

message ListUsersRequest {
    int32 page = 1;
    int32 page_size = 2;
    string filter = 3;
}

message ListUsersResponse {
    repeated User users = 1;
    int32 total = 2;
    int32 page = 3;
    int32 page_size = 4;
}

message StreamUsersRequest {
    string filter = 1;
}
```

### 3. Generate Go Code

```bash
# Generate Go code from proto files
protoc --go_out=. --go-grpc_out=. pkg/proto/user.proto
```

## 🔧 gRPC Server Implementation

### 1. Server Setup

```go
// cmd/main.go
package main

import (
    "context"
    "log"
    "net"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    
    "github.com/anasamu/go-micro-libs/config"
    "github.com/anasamu/go-micro-libs/logging"
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/communication"
    "github.com/anasamu/go-micro-libs/database"
    
    pb "user-service/pkg/proto"
    "user-service/internal/grpc"
    "user-service/internal/services"
)

func main() {
    ctx := context.Background()
    
    // Initialize core managers
    configManager := config.NewManager()
    loggingManager := logging.NewManager()
    monitoringManager := monitoring.NewManager()
    commManager := communication.NewManager()
    dbManager := database.NewManager()
    
    // Bootstrap service
    if err := bootstrapService(ctx, configManager, loggingManager, 
        monitoringManager, commManager, dbManager); err != nil {
        log.Fatal("Failed to bootstrap service:", err)
    }
    
    // Start gRPC server
    if err := startGRPCServer(ctx, commManager, dbManager); err != nil {
        log.Fatal("Failed to start gRPC server:", err)
    }
}

func startGRPCServer(ctx context.Context, commManager *communication.CommunicationManager, dbManager *database.DatabaseManager) error {
    // Create gRPC server
    server := grpc.NewServer(
        grpc.UnaryInterceptor(monitoring.UnaryServerInterceptor()),
        grpc.StreamInterceptor(monitoring.StreamServerInterceptor()),
    )
    
    // Register service
    userService := services.NewUserService(dbManager)
    userGRPCService := grpc.NewUserGRPCService(userService, commManager)
    pb.RegisterUserServiceServer(server, userGRPCService)
    
    // Enable reflection for debugging
    reflection.Register(server)
    
    // Start server
    listener, err := net.Listen("tcp", ":9090")
    if err != nil {
        return err
    }
    
    log.Println("gRPC server started on :9090")
    return server.Serve(listener)
}
```

### 2. Service Implementation

```go
// internal/grpc/user_service.go
package grpc

import (
    "context"
    "log"
    
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    
    "github.com/anasamu/go-micro-libs/communication"
    "github.com/anasamu/go-micro-libs/monitoring"
    
    pb "user-service/pkg/proto"
    "user-service/internal/services"
)

type UserGRPCService struct {
    pb.UnimplementedUserServiceServer
    userService services.UserService
    commManager *communication.CommunicationManager
}

func NewUserGRPCService(userService services.UserService, commManager *communication.CommunicationManager) *UserGRPCService {
    return &UserGRPCService{
        userService: userService,
        commManager: commManager,
    }
}

func (s *UserGRPCService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    // Record metrics
    monitoring.IncrementCounter("grpc_requests_total", map[string]string{
        "method": "GetUser",
        "service": "UserService",
    })
    
    // Validate request
    if req.UserId == "" {
        return nil, status.Error(codes.InvalidArgument, "user_id is required")
    }
    
    // Call business logic
    user, err := s.userService.GetUser(ctx, req.UserId)
    if err != nil {
        log.Printf("GetUser error: %v", err)
        return nil, status.Error(codes.NotFound, "user not found")
    }
    
    // Convert to protobuf
    pbUser := &pb.User{
        Id:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt.Unix(),
        UpdatedAt: user.UpdatedAt.Unix(),
    }
    
    return &pb.GetUserResponse{
        User: pbUser,
    }, nil
}

func (s *UserGRPCService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    // Record metrics
    monitoring.IncrementCounter("grpc_requests_total", map[string]string{
        "method": "CreateUser",
        "service": "UserService",
    })
    
    // Validate request
    if req.Name == "" || req.Email == "" {
        return nil, status.Error(codes.InvalidArgument, "name and email are required")
    }
    
    // Create user
    user := &services.User{
        Name:  req.Name,
        Email: req.Email,
    }
    
    createdUser, err := s.userService.CreateUser(ctx, user)
    if err != nil {
        log.Printf("CreateUser error: %v", err)
        return nil, status.Error(codes.Internal, "failed to create user")
    }
    
    // Convert to protobuf
    pbUser := &pb.User{
        Id:        createdUser.ID,
        Name:      createdUser.Name,
        Email:     createdUser.Email,
        CreatedAt: createdUser.CreatedAt.Unix(),
        UpdatedAt: createdUser.UpdatedAt.Unix(),
    }
    
    return &pb.CreateUserResponse{
        User: pbUser,
    }, nil
}

func (s *UserGRPCService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
    // Record metrics
    monitoring.IncrementCounter("grpc_requests_total", map[string]string{
        "method": "UpdateUser",
        "service": "UserService",
    })
    
    // Validate request
    if req.UserId == "" {
        return nil, status.Error(codes.InvalidArgument, "user_id is required")
    }
    
    // Update user
    user := &services.User{
        ID:    req.UserId,
        Name:  req.Name,
        Email: req.Email,
    }
    
    updatedUser, err := s.userService.UpdateUser(ctx, user)
    if err != nil {
        log.Printf("UpdateUser error: %v", err)
        return nil, status.Error(codes.NotFound, "user not found")
    }
    
    // Convert to protobuf
    pbUser := &pb.User{
        Id:        updatedUser.ID,
        Name:      updatedUser.Name,
        Email:     updatedUser.Email,
        CreatedAt: updatedUser.CreatedAt.Unix(),
        UpdatedAt: updatedUser.UpdatedAt.Unix(),
    }
    
    return &pb.UpdateUserResponse{
        User: pbUser,
    }, nil
}

func (s *UserGRPCService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
    // Record metrics
    monitoring.IncrementCounter("grpc_requests_total", map[string]string{
        "method": "DeleteUser",
        "service": "UserService",
    })
    
    // Validate request
    if req.UserId == "" {
        return nil, status.Error(codes.InvalidArgument, "user_id is required")
    }
    
    // Delete user
    err := s.userService.DeleteUser(ctx, req.UserId)
    if err != nil {
        log.Printf("DeleteUser error: %v", err)
        return nil, status.Error(codes.NotFound, "user not found")
    }
    
    return &pb.DeleteUserResponse{
        Success: true,
    }, nil
}

func (s *UserGRPCService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
    // Record metrics
    monitoring.IncrementCounter("grpc_requests_total", map[string]string{
        "method": "ListUsers",
        "service": "UserService",
    })
    
    // Validate request
    if req.Page <= 0 {
        req.Page = 1
    }
    if req.PageSize <= 0 {
        req.PageSize = 10
    }
    
    // Get users
    users, total, err := s.userService.ListUsers(ctx, req.Page, req.PageSize, req.Filter)
    if err != nil {
        log.Printf("ListUsers error: %v", err)
        return nil, status.Error(codes.Internal, "failed to list users")
    }
    
    // Convert to protobuf
    pbUsers := make([]*pb.User, len(users))
    for i, user := range users {
        pbUsers[i] = &pb.User{
            Id:        user.ID,
            Name:      user.Name,
            Email:     user.Email,
            CreatedAt: user.CreatedAt.Unix(),
            UpdatedAt: user.UpdatedAt.Unix(),
        }
    }
    
    return &pb.ListUsersResponse{
        Users:    pbUsers,
        Total:    int32(total),
        Page:     req.Page,
        PageSize: req.PageSize,
    }, nil
}

func (s *UserGRPCService) StreamUsers(req *pb.StreamUsersRequest, stream pb.UserService_StreamUsersServer) error {
    // Record metrics
    monitoring.IncrementCounter("grpc_requests_total", map[string]string{
        "method": "StreamUsers",
        "service": "UserService",
    })
    
    // Stream users
    userChan := make(chan *services.User, 100)
    errChan := make(chan error, 1)
    
    go func() {
        defer close(userChan)
        err := s.userService.StreamUsers(stream.Context(), req.Filter, userChan)
        if err != nil {
            errChan <- err
        }
    }()
    
    for {
        select {
        case user, ok := <-userChan:
            if !ok {
                return nil
            }
            
            // Convert to protobuf
            pbUser := &pb.User{
                Id:        user.ID,
                Name:      user.Name,
                Email:     user.Email,
                CreatedAt: user.CreatedAt.Unix(),
                UpdatedAt: user.UpdatedAt.Unix(),
            }
            
            if err := stream.Send(pbUser); err != nil {
                return err
            }
            
        case err := <-errChan:
            return err
            
        case <-stream.Context().Done():
            return stream.Context().Err()
        }
    }
}
```

## 🔧 gRPC Client Implementation

### 1. Client Setup

```go
// internal/clients/user_grpc_client.go
package clients

import (
    "context"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/keepalive"
    
    "github.com/anasamu/go-micro-libs/communication"
    "github.com/anasamu/go-micro-libs/monitoring"
    
    pb "user-service/pkg/proto"
)

type UserGRPCClient struct {
    client      pb.UserServiceClient
    conn        *grpc.ClientConn
    commManager *communication.CommunicationManager
}

func NewUserGRPCClient(address string, commManager *communication.CommunicationManager) (*UserGRPCClient, error) {
    // Create connection with keepalive
    conn, err := grpc.Dial(address,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithKeepaliveParams(keepalive.ClientParameters{
            Time:                10 * time.Second,
            Timeout:             3 * time.Second,
            PermitWithoutStream: true,
        }),
        grpc.WithUnaryInterceptor(monitoring.UnaryClientInterceptor()),
        grpc.WithStreamInterceptor(monitoring.StreamClientInterceptor()),
    )
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
    // Record metrics
    monitoring.IncrementCounter("grpc_client_requests_total", map[string]string{
        "method": "GetUser",
        "service": "UserService",
    })
    
    // Create request
    req := &pb.GetUserRequest{
        UserId: userID,
    }
    
    // Call gRPC service
    resp, err := c.client.GetUser(ctx, req)
    if err != nil {
        monitoring.IncrementCounter("grpc_client_errors_total", map[string]string{
            "method": "GetUser",
            "service": "UserService",
        })
        return nil, err
    }
    
    // Convert response
    user := &User{
        ID:        resp.User.Id,
        Name:      resp.User.Name,
        Email:     resp.User.Email,
        CreatedAt: time.Unix(resp.User.CreatedAt, 0),
        UpdatedAt: time.Unix(resp.User.UpdatedAt, 0),
    }
    
    return user, nil
}

func (c *UserGRPCClient) CreateUser(ctx context.Context, user *User) (*User, error) {
    // Record metrics
    monitoring.IncrementCounter("grpc_client_requests_total", map[string]string{
        "method": "CreateUser",
        "service": "UserService",
    })
    
    // Create request
    req := &pb.CreateUserRequest{
        Name:  user.Name,
        Email: user.Email,
    }
    
    // Call gRPC service
    resp, err := c.client.CreateUser(ctx, req)
    if err != nil {
        monitoring.IncrementCounter("grpc_client_errors_total", map[string]string{
            "method": "CreateUser",
            "service": "UserService",
        })
        return nil, err
    }
    
    // Convert response
    createdUser := &User{
        ID:        resp.User.Id,
        Name:      resp.User.Name,
        Email:     resp.User.Email,
        CreatedAt: time.Unix(resp.User.CreatedAt, 0),
        UpdatedAt: time.Unix(resp.User.UpdatedAt, 0),
    }
    
    return createdUser, nil
}

func (c *UserGRPCClient) UpdateUser(ctx context.Context, user *User) (*User, error) {
    // Record metrics
    monitoring.IncrementCounter("grpc_client_requests_total", map[string]string{
        "method": "UpdateUser",
        "service": "UserService",
    })
    
    // Create request
    req := &pb.UpdateUserRequest{
        UserId: user.ID,
        Name:   user.Name,
        Email:  user.Email,
    }
    
    // Call gRPC service
    resp, err := c.client.UpdateUser(ctx, req)
    if err != nil {
        monitoring.IncrementCounter("grpc_client_errors_total", map[string]string{
            "method": "UpdateUser",
            "service": "UserService",
        })
        return nil, err
    }
    
    // Convert response
    updatedUser := &User{
        ID:        resp.User.Id,
        Name:      resp.User.Name,
        Email:     resp.User.Email,
        CreatedAt: time.Unix(resp.User.CreatedAt, 0),
        UpdatedAt: time.Unix(resp.User.UpdatedAt, 0),
    }
    
    return updatedUser, nil
}

func (c *UserGRPCClient) DeleteUser(ctx context.Context, userID string) error {
    // Record metrics
    monitoring.IncrementCounter("grpc_client_requests_total", map[string]string{
        "method": "DeleteUser",
        "service": "UserService",
    })
    
    // Create request
    req := &pb.DeleteUserRequest{
        UserId: userID,
    }
    
    // Call gRPC service
    _, err := c.client.DeleteUser(ctx, req)
    if err != nil {
        monitoring.IncrementCounter("grpc_client_errors_total", map[string]string{
            "method": "DeleteUser",
            "service": "UserService",
        })
        return err
    }
    
    return nil
}

func (c *UserGRPCClient) ListUsers(ctx context.Context, page, pageSize int32, filter string) ([]*User, int32, error) {
    // Record metrics
    monitoring.IncrementCounter("grpc_client_requests_total", map[string]string{
        "method": "ListUsers",
        "service": "UserService",
    })
    
    // Create request
    req := &pb.ListUsersRequest{
        Page:     page,
        PageSize: pageSize,
        Filter:   filter,
    }
    
    // Call gRPC service
    resp, err := c.client.ListUsers(ctx, req)
    if err != nil {
        monitoring.IncrementCounter("grpc_client_errors_total", map[string]string{
            "method": "ListUsers",
            "service": "UserService",
        })
        return nil, 0, err
    }
    
    // Convert response
    users := make([]*User, len(resp.Users))
    for i, pbUser := range resp.Users {
        users[i] = &User{
            ID:        pbUser.Id,
            Name:      pbUser.Name,
            Email:     pbUser.Email,
            CreatedAt: time.Unix(pbUser.CreatedAt, 0),
            UpdatedAt: time.Unix(pbUser.UpdatedAt, 0),
        }
    }
    
    return users, resp.Total, nil
}

func (c *UserGRPCClient) StreamUsers(ctx context.Context, filter string) (<-chan *User, error) {
    // Record metrics
    monitoring.IncrementCounter("grpc_client_requests_total", map[string]string{
        "method": "StreamUsers",
        "service": "UserService",
    })
    
    // Create request
    req := &pb.StreamUsersRequest{
        Filter: filter,
    }
    
    // Create stream
    stream, err := c.client.StreamUsers(ctx, req)
    if err != nil {
        monitoring.IncrementCounter("grpc_client_errors_total", map[string]string{
            "method": "StreamUsers",
            "service": "UserService",
        })
        return nil, err
    }
    
    // Create channel for users
    userChan := make(chan *User, 100)
    
    // Start goroutine to receive users
    go func() {
        defer close(userChan)
        
        for {
            pbUser, err := stream.Recv()
            if err != nil {
                return
            }
            
            // Convert to User
            user := &User{
                ID:        pbUser.Id,
                Name:      pbUser.Name,
                Email:     pbUser.Email,
                CreatedAt: time.Unix(pbUser.CreatedAt, 0),
                UpdatedAt: time.Unix(pbUser.UpdatedAt, 0),
            }
            
            select {
            case userChan <- user:
            case <-ctx.Done():
                return
            }
        }
    }()
    
    return userChan, nil
}

func (c *UserGRPCClient) Close() error {
    return c.conn.Close()
}
```

## 🔧 Service Discovery Integration

### 1. gRPC with Service Discovery

```go
// internal/clients/discovery_grpc_client.go
package clients

import (
    "context"
    "fmt"
    "math/rand"
    "time"
    
    "github.com/anasamu/go-micro-libs/discovery"
    "github.com/anasamu/go-micro-libs/communication"
)

type DiscoveryGRPCClient struct {
    discoveryManager *discovery.DiscoveryManager
    commManager      *communication.CommunicationManager
    serviceName      string
    clients          map[string]*UserGRPCClient
}

func NewDiscoveryGRPCClient(discoveryManager *discovery.DiscoveryManager, 
    commManager *communication.CommunicationManager, serviceName string) *DiscoveryGRPCClient {
    return &DiscoveryGRPCClient{
        discoveryManager: discoveryManager,
        commManager:      commManager,
        serviceName:      serviceName,
        clients:          make(map[string]*UserGRPCClient),
    }
}

func (c *DiscoveryGRPCClient) GetUser(ctx context.Context, userID string) (*User, error) {
    client, err := c.getClient(ctx)
    if err != nil {
        return nil, err
    }
    
    return client.GetUser(ctx, userID)
}

func (c *DiscoveryGRPCClient) getClient(ctx context.Context) (*UserGRPCClient, error) {
    // Discover services
    services, err := c.discoveryManager.DiscoverServices(ctx, c.serviceName)
    if err != nil {
        return nil, err
    }
    
    if len(services) == 0 {
        return nil, fmt.Errorf("no services found for %s", c.serviceName)
    }
    
    // Select service (simple load balancing)
    service := services[rand.Intn(len(services))]
    address := fmt.Sprintf("%s:%d", service.Address, service.Port)
    
    // Check if client already exists
    if client, exists := c.clients[address]; exists {
        return client, nil
    }
    
    // Create new client
    client, err := NewUserGRPCClient(address, c.commManager)
    if err != nil {
        return nil, err
    }
    
    // Cache client
    c.clients[address] = client
    
    return client, nil
}
```

## 🔧 Load Balancing

### 1. gRPC Load Balancer

```go
// internal/loadbalancer/grpc_load_balancer.go
package loadbalancer

import (
    "context"
    "fmt"
    "sync"
    
    "github.com/anasamu/go-micro-libs/discovery"
    "github.com/anasamu/go-micro-libs/communication"
)

type GRPCLoadBalancer struct {
    discoveryManager *discovery.DiscoveryManager
    commManager      *communication.CommunicationManager
    serviceName      string
    clients          map[string]*UserGRPCClient
    mutex            sync.RWMutex
}

func NewGRPCLoadBalancer(discoveryManager *discovery.DiscoveryManager, 
    commManager *communication.CommunicationManager, serviceName string) *GRPCLoadBalancer {
    return &GRPCLoadBalancer{
        discoveryManager: discoveryManager,
        commManager:      commManager,
        serviceName:      serviceName,
        clients:          make(map[string]*UserGRPCClient),
    }
}

func (lb *GRPCLoadBalancer) GetUser(ctx context.Context, userID string) (*User, error) {
    client, err := lb.getClient(ctx)
    if err != nil {
        return nil, err
    }
    
    return client.GetUser(ctx, userID)
}

func (lb *GRPCLoadBalancer) getClient(ctx context.Context) (*UserGRPCClient, error) {
    // Discover services
    services, err := lb.discoveryManager.DiscoverServices(ctx, lb.serviceName)
    if err != nil {
        return nil, err
    }
    
    if len(services) == 0 {
        return nil, fmt.Errorf("no services found for %s", lb.serviceName)
    }
    
    // Select service using round-robin
    service := lb.selectService(services)
    address := fmt.Sprintf("%s:%d", service.Address, service.Port)
    
    lb.mutex.RLock()
    client, exists := lb.clients[address]
    lb.mutex.RUnlock()
    
    if exists {
        return client, nil
    }
    
    // Create new client
    client, err = NewUserGRPCClient(address, lb.commManager)
    if err != nil {
        return nil, err
    }
    
    // Cache client
    lb.mutex.Lock()
    lb.clients[address] = client
    lb.mutex.Unlock()
    
    return client, nil
}

func (lb *GRPCLoadBalancer) selectService(services []*discovery.ServiceInfo) *discovery.ServiceInfo {
    // Simple round-robin selection
    return services[0] // In real implementation, use proper round-robin
}
```

## 🔧 Configuration

### gRPC Configuration

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

# gRPC communication
communication:
  providers:
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
      compression:
        enabled: true
        algorithm: "gzip"
      tls:
        enabled: false
        cert_file: ""
        key_file: ""

# Service discovery
discovery:
  providers:
    consul:
      enabled: true
      address: "localhost:8500"
      token: "${CONSUL_TOKEN}"
      service_name: "user-service"
      service_port: 9090
      health_check:
        enabled: true
        path: "/health"
        interval: "10s"
        timeout: "3s"
        deregister_after: "30s"

# Database
database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100
```

## 🔧 Best Practices

### 1. Error Handling

```go
// Proper error handling for gRPC
func (s *UserGRPCService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    // Validate request
    if req.UserId == "" {
        return nil, status.Error(codes.InvalidArgument, "user_id is required")
    }
    
    // Call business logic
    user, err := s.userService.GetUser(ctx, req.UserId)
    if err != nil {
        // Log error
        log.Printf("GetUser error: %v", err)
        
        // Return appropriate gRPC error
        if errors.Is(err, services.ErrUserNotFound) {
            return nil, status.Error(codes.NotFound, "user not found")
        }
        
        return nil, status.Error(codes.Internal, "internal server error")
    }
    
    // Convert and return response
    return &pb.GetUserResponse{
        User: convertToPBUser(user),
    }, nil
}
```

### 2. Timeout Management

```go
// Timeout management for gRPC
func (c *UserGRPCClient) GetUserWithTimeout(ctx context.Context, userID string, timeout time.Duration) (*User, error) {
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()
    
    return c.GetUser(ctx, userID)
}
```

### 3. Circuit Breaker Integration

```go
// Circuit breaker for gRPC
func (c *UserGRPCClient) GetUserWithCircuitBreaker(ctx context.Context, userID string) (*User, error) {
    result, err := c.circuitBreaker.Execute(func() (interface{}, error) {
        return c.GetUser(ctx, userID)
    })
    
    if err != nil {
        return nil, err
    }
    
    return result.(*User), nil
}
```

### 4. Health Checks

```go
// Health check for gRPC service
func (s *UserGRPCService) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
    // Check database connection
    if err := s.dbManager.HealthCheck(ctx); err != nil {
        return &pb.HealthCheckResponse{
            Status: pb.HealthCheckResponse_NOT_SERVING,
        }, nil
    }
    
    return &pb.HealthCheckResponse{
        Status: pb.HealthCheckResponse_SERVING,
    }, nil
}
```

### 5. Metrics and Monitoring

```go
// Metrics for gRPC
func (s *UserGRPCService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    start := time.Now()
    
    // Record request metrics
    monitoring.IncrementCounter("grpc_requests_total", map[string]string{
        "method": "GetUser",
        "service": "UserService",
    })
    
    // Call business logic
    resp, err := s.getUser(ctx, req)
    
    // Record latency
    monitoring.RecordHistogram("grpc_request_duration_seconds", 
        time.Since(start).Seconds(), map[string]string{
            "method": "GetUser",
            "service": "UserService",
        })
    
    // Record error metrics
    if err != nil {
        monitoring.IncrementCounter("grpc_errors_total", map[string]string{
            "method": "GetUser",
            "service": "UserService",
            "error": err.Error(),
        })
    }
    
    return resp, err
}
```

---

**gRPC Communication - High-performance inter-service communication! 🚀**

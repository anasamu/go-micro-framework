# Service Discovery Implementation

## 🎯 Overview

GoMicroFramework menyediakan sistem service discovery yang komprehensif untuk memungkinkan services menemukan dan berkomunikasi dengan services lain secara dinamis. Framework mendukung berbagai service discovery providers seperti Consul, etcd, Kubernetes, dan static configuration.

## 🔧 Supported Service Discovery Providers

### 1. Consul
- Distributed service mesh
- Health checking
- Key-value storage
- Service segmentation

### 2. etcd
- Distributed key-value store
- High availability
- Strong consistency
- Watch functionality

### 3. Kubernetes
- Native Kubernetes service discovery
- DNS-based discovery
- Service mesh integration
- Pod and service monitoring

### 4. Static
- Static service configuration
- Manual service registration
- Simple service discovery
- Development and testing

## 🔧 Service Discovery Configuration

### 1. Generate Service with Service Discovery

```bash
# Generate service with Consul service discovery
microframework new user-service --with-discovery=consul --with-database=postgres --with-monitoring=prometheus

# Generate service with etcd service discovery
microframework new order-service --with-discovery=etcd --with-database=postgres

# Generate service with Kubernetes service discovery
microframework new payment-service --with-discovery=kubernetes --with-database=postgres
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

# Service discovery configuration
discovery:
  providers:
    consul:
      enabled: true
      address: "localhost:8500"
      token: "${CONSUL_TOKEN}"
      datacenter: "dc1"
      service_name: "user-service"
      service_port: 8080
      service_address: "localhost"
      service_tags: ["api", "user", "v1"]
      service_meta:
        version: "1.0.0"
        environment: "production"
      health_check:
        enabled: true
        path: "/health"
        interval: "10s"
        timeout: "3s"
        deregister_after: "30s"
        tcp_check:
          enabled: false
          port: 8080
        http_check:
          enabled: true
          path: "/health"
          method: "GET"
          headers:
            "Content-Type": "application/json"
      watch:
        enabled: true
        services: ["user-service", "order-service", "payment-service"]
        
    etcd:
      enabled: false
      endpoints: ["http://localhost:2379"]
      username: "${ETCD_USERNAME}"
      password: "${ETCD_PASSWORD}"
      tls:
        enabled: false
        cert_file: ""
        key_file: ""
        ca_file: ""
      service_prefix: "/services/"
      service_name: "user-service"
      service_port: 8080
      service_address: "localhost"
      ttl: "30s"
      health_check:
        enabled: true
        path: "/health"
        interval: "10s"
        timeout: "3s"
        
    kubernetes:
      enabled: false
      config_path: "${KUBERNETES_CONFIG}"
      namespace: "default"
      service_name: "user-service"
      service_port: 8080
      service_address: "localhost"
      labels:
        app: "user-service"
        version: "1.0.0"
      annotations:
        description: "User service for user management"
      health_check:
        enabled: true
        path: "/health"
        interval: "10s"
        timeout: "3s"
        
    static:
      enabled: false
      services:
        - name: "user-service"
          address: "localhost"
          port: 8080
          tags: ["api", "user"]
          meta:
            version: "1.0.0"
        - name: "order-service"
          address: "localhost"
          port: 8081
          tags: ["api", "order"]
          meta:
            version: "1.0.0"

# Database for service discovery data
database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100
```

## 🔧 Service Discovery Implementation

### 1. Service Discovery Manager

```go
// internal/discovery/manager.go
package discovery

import (
    "context"
    "fmt"
    "sync"
    "time"
    
    "github.com/anasamu/go-micro-libs/discovery"
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/logging"
)

type Manager struct {
    providers    map[string]DiscoveryProvider
    config       *Config
    logger       *logging.Logger
    metrics      *monitoring.Metrics
    watchers     map[string]*ServiceWatcher
    mutex        sync.RWMutex
}

type Config struct {
    Provider        string        `yaml:"provider"`
    ServiceName     string        `yaml:"service_name"`
    ServicePort     int           `yaml:"service_port"`
    ServiceAddress  string        `yaml:"service_address"`
    ServiceTags     []string      `yaml:"service_tags"`
    ServiceMeta     map[string]string `yaml:"service_meta"`
    HealthCheck     HealthCheckConfig `yaml:"health_check"`
    Watch           WatchConfig   `yaml:"watch"`
}

type HealthCheckConfig struct {
    Enabled         bool          `yaml:"enabled"`
    Path            string        `yaml:"path"`
    Interval        time.Duration `yaml:"interval"`
    Timeout         time.Duration `yaml:"timeout"`
    DeregisterAfter time.Duration `yaml:"deregister_after"`
    TCPCheck        TCPCheckConfig `yaml:"tcp_check"`
    HTTPCheck       HTTPCheckConfig `yaml:"http_check"`
}

type TCPCheckConfig struct {
    Enabled bool `yaml:"enabled"`
    Port    int  `yaml:"port"`
}

type HTTPCheckConfig struct {
    Enabled bool              `yaml:"enabled"`
    Path    string            `yaml:"path"`
    Method  string            `yaml:"method"`
    Headers map[string]string `yaml:"headers"`
}

type WatchConfig struct {
    Enabled  bool     `yaml:"enabled"`
    Services []string `yaml:"services"`
}

type DiscoveryProvider interface {
    RegisterService(ctx context.Context, service *ServiceInfo) error
    DeregisterService(ctx context.Context, serviceID string) error
    DiscoverServices(ctx context.Context, serviceName string) ([]*ServiceInfo, error)
    WatchServices(ctx context.Context, serviceName string, callback func([]*ServiceInfo)) error
    HealthCheck(ctx context.Context, serviceID string) error
}

type ServiceInfo struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Address     string            `json:"address"`
    Port        int               `json:"port"`
    Tags        []string          `json:"tags"`
    Meta        map[string]string `json:"meta"`
    Status      ServiceStatus     `json:"status"`
    LastSeen    time.Time         `json:"last_seen"`
}

type ServiceStatus int

const (
    StatusHealthy ServiceStatus = iota
    StatusUnhealthy
    StatusUnknown
)

func (s ServiceStatus) String() string {
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

type ServiceWatcher struct {
    serviceName string
    provider    DiscoveryProvider
    callback    func([]*ServiceInfo)
    stopChan    chan struct{}
    logger      *logging.Logger
    metrics     *monitoring.Metrics
}

func NewManager(config *Config, logger *logging.Logger, metrics *monitoring.Metrics) *Manager {
    return &Manager{
        providers: make(map[string]DiscoveryProvider),
        config:    config,
        logger:    logger,
        metrics:   metrics,
        watchers:  make(map[string]*ServiceWatcher),
    }
}

func (m *Manager) RegisterProvider(name string, provider DiscoveryProvider) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    m.providers[name] = provider
    m.logger.Info("Discovery provider registered", "name", name)
}

func (m *Manager) RegisterService(ctx context.Context, providerName string) error {
    m.mutex.RLock()
    provider, exists := m.providers[providerName]
    m.mutex.RUnlock()
    
    if !exists {
        return fmt.Errorf("provider %s not found", providerName)
    }
    
    service := &ServiceInfo{
        ID:         generateServiceID(),
        Name:       m.config.ServiceName,
        Address:    m.config.ServiceAddress,
        Port:       m.config.ServicePort,
        Tags:       m.config.ServiceTags,
        Meta:       m.config.ServiceMeta,
        Status:     StatusHealthy,
        LastSeen:   time.Now(),
    }
    
    err := provider.RegisterService(ctx, service)
    if err != nil {
        m.logger.Error("Failed to register service", "provider", providerName, "error", err)
        return err
    }
    
    m.logger.Info("Service registered", "provider", providerName, "service", service.Name)
    
    // Record metrics
    m.metrics.IncrementCounter("service_registered_total", map[string]string{
        "provider": providerName,
        "service":  service.Name,
    })
    
    return nil
}

func (m *Manager) DeregisterService(ctx context.Context, providerName, serviceID string) error {
    m.mutex.RLock()
    provider, exists := m.providers[providerName]
    m.mutex.RUnlock()
    
    if !exists {
        return fmt.Errorf("provider %s not found", providerName)
    }
    
    err := provider.DeregisterService(ctx, serviceID)
    if err != nil {
        m.logger.Error("Failed to deregister service", "provider", providerName, "error", err)
        return err
    }
    
    m.logger.Info("Service deregistered", "provider", providerName, "service_id", serviceID)
    
    // Record metrics
    m.metrics.IncrementCounter("service_deregistered_total", map[string]string{
        "provider": providerName,
    })
    
    return nil
}

func (m *Manager) DiscoverServices(ctx context.Context, providerName, serviceName string) ([]*ServiceInfo, error) {
    m.mutex.RLock()
    provider, exists := m.providers[providerName]
    m.mutex.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("provider %s not found", providerName)
    }
    
    services, err := provider.DiscoverServices(ctx, serviceName)
    if err != nil {
        m.logger.Error("Failed to discover services", "provider", providerName, "service", serviceName, "error", err)
        return nil, err
    }
    
    // Record metrics
    m.metrics.IncrementCounter("service_discovery_requests_total", map[string]string{
        "provider": providerName,
        "service":  serviceName,
    })
    
    return services, nil
}

func (m *Manager) WatchServices(ctx context.Context, providerName, serviceName string, callback func([]*ServiceInfo)) error {
    m.mutex.RLock()
    provider, exists := m.providers[providerName]
    m.mutex.RUnlock()
    
    if !exists {
        return fmt.Errorf("provider %s not found", providerName)
    }
    
    watcher := &ServiceWatcher{
        serviceName: serviceName,
        provider:    provider,
        callback:    callback,
        stopChan:    make(chan struct{}),
        logger:      m.logger,
        metrics:     m.metrics,
    }
    
    m.mutex.Lock()
    m.watchers[serviceName] = watcher
    m.mutex.Unlock()
    
    go watcher.Start(ctx)
    
    m.logger.Info("Started watching services", "provider", providerName, "service", serviceName)
    
    return nil
}

func (m *Manager) StopWatching(serviceName string) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    watcher, exists := m.watchers[serviceName]
    if exists {
        watcher.Stop()
        delete(m.watchers, serviceName)
        m.logger.Info("Stopped watching services", "service", serviceName)
    }
}

func (m *Manager) StartWatching(ctx context.Context, providerName string) error {
    if !m.config.Watch.Enabled {
        return nil
    }
    
    for _, serviceName := range m.config.Watch.Services {
        err := m.WatchServices(ctx, providerName, serviceName, func(services []*ServiceInfo) {
            m.logger.Info("Service list updated", "service", serviceName, "count", len(services))
            
            // Record metrics
            m.metrics.SetGauge("discovered_services_count", float64(len(services)), map[string]string{
                "service": serviceName,
            })
        })
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

### 2. Service Watcher Implementation

```go
// internal/discovery/watcher.go
package discovery

import (
    "context"
    "time"
    
    "github.com/anasamu/go-micro-libs/monitoring"
    "github.com/anasamu/go-micro-libs/logging"
)

func (sw *ServiceWatcher) Start(ctx context.Context) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            sw.checkServices(ctx)
        case <-sw.stopChan:
            return
        case <-ctx.Done():
            return
        }
    }
}

func (sw *ServiceWatcher) Stop() {
    close(sw.stopChan)
}

func (sw *ServiceWatcher) checkServices(ctx context.Context) {
    services, err := sw.provider.DiscoverServices(ctx, sw.serviceName)
    if err != nil {
        sw.logger.Error("Failed to discover services", "service", sw.serviceName, "error", err)
        return
    }
    
    // Call callback with updated services
    sw.callback(services)
    
    // Record metrics
    sw.metrics.SetGauge("discovered_services_count", float64(len(services)), map[string]string{
        "service": sw.serviceName,
    })
}
```

### 3. Consul Provider Implementation

```go
// internal/discovery/consul_provider.go
package discovery

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
    datacenter  string
    logger      *logging.Logger
}

func NewConsulProvider(address, token, datacenter string, logger *logging.Logger) (*ConsulProvider, error) {
    config := api.DefaultConfig()
    config.Address = address
    config.Token = token
    config.Datacenter = datacenter
    
    client, err := api.NewClient(config)
    if err != nil {
        return nil, err
    }
    
    return &ConsulProvider{
        client:     client,
        datacenter: datacenter,
        logger:     logger,
    }, nil
}

func (cp *ConsulProvider) RegisterService(ctx context.Context, service *ServiceInfo) error {
    registration := &api.AgentServiceRegistration{
        ID:      service.ID,
        Name:    service.Name,
        Port:    service.Port,
        Address: service.Address,
        Tags:    service.Tags,
        Meta:    service.Meta,
        Check: &api.AgentServiceCheck{
            HTTP:                           fmt.Sprintf("http://%s:%d/health", service.Address, service.Port),
            Timeout:                        "3s",
            Interval:                       "10s",
            DeregisterCriticalServiceAfter: "30s",
        },
    }
    
    err := cp.client.Agent().ServiceRegister(registration)
    if err != nil {
        return err
    }
    
    cp.logger.Info("Service registered with Consul",
        "service", service.Name,
        "address", service.Address,
        "port", service.Port,
    )
    
    return nil
}

func (cp *ConsulProvider) DeregisterService(ctx context.Context, serviceID string) error {
    err := cp.client.Agent().ServiceDeregister(serviceID)
    if err != nil {
        return err
    }
    
    cp.logger.Info("Service deregistered from Consul", "service_id", serviceID)
    
    return nil
}

func (cp *ConsulProvider) DiscoverServices(ctx context.Context, serviceName string) ([]*ServiceInfo, error) {
    services, _, err := cp.client.Health().Service(serviceName, "", true, nil)
    if err != nil {
        return nil, err
    }
    
    serviceInfos := make([]*ServiceInfo, len(services))
    for i, service := range services {
        serviceInfo := &ServiceInfo{
            ID:       service.Service.ID,
            Name:     service.Service.Service,
            Address:  service.Service.Address,
            Port:     service.Service.Port,
            Tags:     service.Service.Tags,
            Meta:     service.Service.Meta,
            Status:   StatusHealthy,
            LastSeen: time.Now(),
        }
        
        // Check service health
        if service.Checks.AggregatedStatus() != api.HealthPassing {
            serviceInfo.Status = StatusUnhealthy
        }
        
        serviceInfos[i] = serviceInfo
    }
    
    return serviceInfos, nil
}

func (cp *ConsulProvider) WatchServices(ctx context.Context, serviceName string, callback func([]*ServiceInfo)) error {
    // Use Consul's watch functionality
    watchConfig := &api.QueryOptions{
        WaitTime: 10 * time.Second,
    }
    
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                services, meta, err := cp.client.Health().Service(serviceName, "", true, watchConfig)
                if err != nil {
                    cp.logger.Error("Failed to watch services", "service", serviceName, "error", err)
                    time.Sleep(5 * time.Second)
                    continue
                }
                
                // Convert to ServiceInfo
                serviceInfos := make([]*ServiceInfo, len(services))
                for i, service := range services {
                    serviceInfo := &ServiceInfo{
                        ID:       service.Service.ID,
                        Name:     service.Service.Service,
                        Address:  service.Service.Address,
                        Port:     service.Service.Port,
                        Tags:     service.Service.Tags,
                        Meta:     service.Service.Meta,
                        Status:   StatusHealthy,
                        LastSeen: time.Now(),
                    }
                    
                    if service.Checks.AggregatedStatus() != api.HealthPassing {
                        serviceInfo.Status = StatusUnhealthy
                    }
                    
                    serviceInfos[i] = serviceInfo
                }
                
                // Call callback
                callback(serviceInfos)
                
                // Update watch index
                watchConfig.WaitIndex = meta.LastIndex
            }
        }
    }()
    
    return nil
}

func (cp *ConsulProvider) HealthCheck(ctx context.Context, serviceID string) error {
    // Check if service is healthy in Consul
    services, _, err := cp.client.Health().Service(serviceID, "", true, nil)
    if err != nil {
        return err
    }
    
    if len(services) == 0 {
        return fmt.Errorf("service %s not found", serviceID)
    }
    
    if services[0].Checks.AggregatedStatus() != api.HealthPassing {
        return fmt.Errorf("service %s is unhealthy", serviceID)
    }
    
    return nil
}
```

### 4. etcd Provider Implementation

```go
// internal/discovery/etcd_provider.go
package discovery

import (
    "context"
    "encoding/json"
    "fmt"
    "path"
    "time"
    
    "go.etcd.io/etcd/clientv3"
    "github.com/anasamu/go-micro-libs/logging"
)

type EtcdProvider struct {
    client       *clientv3.Client
    servicePrefix string
    logger       *logging.Logger
}

func NewEtcdProvider(endpoints []string, username, password string, servicePrefix string, logger *logging.Logger) (*EtcdProvider, error) {
    config := clientv3.Config{
        Endpoints:   endpoints,
        DialTimeout: 5 * time.Second,
    }
    
    if username != "" && password != "" {
        config.Username = username
        config.Password = password
    }
    
    client, err := clientv3.New(config)
    if err != nil {
        return nil, err
    }
    
    return &EtcdProvider{
        client:       client,
        servicePrefix: servicePrefix,
        logger:       logger,
    }, nil
}

func (ep *EtcdProvider) RegisterService(ctx context.Context, service *ServiceInfo) error {
    key := path.Join(ep.servicePrefix, service.Name, service.ID)
    
    serviceData, err := json.Marshal(service)
    if err != nil {
        return err
    }
    
    // Register service with TTL
    lease, err := ep.client.Grant(ctx, 30) // 30 second TTL
    if err != nil {
        return err
    }
    
    _, err = ep.client.Put(ctx, key, string(serviceData), clientv3.WithLease(lease.ID))
    if err != nil {
        return err
    }
    
    // Keep alive the lease
    ch, err := ep.client.KeepAlive(ctx, lease.ID)
    if err != nil {
        return err
    }
    
    // Start keep alive goroutine
    go func() {
        for range ch {
            // Keep alive response received
        }
    }()
    
    ep.logger.Info("Service registered with etcd",
        "service", service.Name,
        "address", service.Address,
        "port", service.Port,
    )
    
    return nil
}

func (ep *EtcdProvider) DeregisterService(ctx context.Context, serviceID string) error {
    // Find and delete service
    resp, err := ep.client.Get(ctx, ep.servicePrefix, clientv3.WithPrefix())
    if err != nil {
        return err
    }
    
    for _, kv := range resp.Kvs {
        var service ServiceInfo
        err = json.Unmarshal(kv.Value, &service)
        if err != nil {
            continue
        }
        
        if service.ID == serviceID {
            _, err = ep.client.Delete(ctx, string(kv.Key))
            if err != nil {
                return err
            }
            
            ep.logger.Info("Service deregistered from etcd", "service_id", serviceID)
            return nil
        }
    }
    
    return fmt.Errorf("service %s not found", serviceID)
}

func (ep *EtcdProvider) DiscoverServices(ctx context.Context, serviceName string) ([]*ServiceInfo, error) {
    key := path.Join(ep.servicePrefix, serviceName)
    
    resp, err := ep.client.Get(ctx, key, clientv3.WithPrefix())
    if err != nil {
        return nil, err
    }
    
    services := make([]*ServiceInfo, 0, len(resp.Kvs))
    for _, kv := range resp.Kvs {
        var service ServiceInfo
        err = json.Unmarshal(kv.Value, &service)
        if err != nil {
            continue
        }
        
        services = append(services, &service)
    }
    
    return services, nil
}

func (ep *EtcdProvider) WatchServices(ctx context.Context, serviceName string, callback func([]*ServiceInfo)) error {
    key := path.Join(ep.servicePrefix, serviceName)
    
    go func() {
        watchChan := ep.client.Watch(ctx, key, clientv3.WithPrefix())
        
        for watchResp := range watchChan {
            if watchResp.Err() != nil {
                ep.logger.Error("etcd watch error", "error", watchResp.Err())
                continue
            }
            
            // Get current services
            services, err := ep.DiscoverServices(ctx, serviceName)
            if err != nil {
                ep.logger.Error("Failed to discover services", "error", err)
                continue
            }
            
            // Call callback
            callback(services)
        }
    }()
    
    return nil
}

func (ep *EtcdProvider) HealthCheck(ctx context.Context, serviceID string) error {
    // Check if service exists in etcd
    resp, err := ep.client.Get(ctx, ep.servicePrefix, clientv3.WithPrefix())
    if err != nil {
        return err
    }
    
    for _, kv := range resp.Kvs {
        var service ServiceInfo
        err = json.Unmarshal(kv.Value, &service)
        if err != nil {
            continue
        }
        
        if service.ID == serviceID {
            return nil
        }
    }
    
    return fmt.Errorf("service %s not found", serviceID)
}
```

### 5. Kubernetes Provider Implementation

```go
// internal/discovery/kubernetes_provider.go
package discovery

import (
    "context"
    "fmt"
    "time"
    
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/labels"
    "github.com/anasamu/go-micro-libs/logging"
)

type KubernetesProvider struct {
    client      *kubernetes.Clientset
    namespace   string
    logger      *logging.Logger
}

func NewKubernetesProvider(configPath, namespace string, logger *logging.Logger) (*KubernetesProvider, error) {
    var config *rest.Config
    var err error
    
    if configPath != "" {
        config, err = clientcmd.BuildConfigFromFlags("", configPath)
    } else {
        config, err = rest.InClusterConfig()
    }
    
    if err != nil {
        return nil, err
    }
    
    client, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, err
    }
    
    return &KubernetesProvider{
        client:    client,
        namespace: namespace,
        logger:    logger,
    }, nil
}

func (kp *KubernetesProvider) RegisterService(ctx context.Context, service *ServiceInfo) error {
    // In Kubernetes, services are typically registered via Service and Endpoints resources
    // This is usually handled by the Kubernetes cluster itself
    // We'll just log the registration
    
    kp.logger.Info("Service registered with Kubernetes",
        "service", service.Name,
        "address", service.Address,
        "port", service.Port,
    )
    
    return nil
}

func (kp *KubernetesProvider) DeregisterService(ctx context.Context, serviceID string) error {
    // In Kubernetes, services are typically deregistered via Service and Endpoints resources
    // This is usually handled by the Kubernetes cluster itself
    // We'll just log the deregistration
    
    kp.logger.Info("Service deregistered from Kubernetes", "service_id", serviceID)
    
    return nil
}

func (kp *KubernetesProvider) DiscoverServices(ctx context.Context, serviceName string) ([]*ServiceInfo, error) {
    // Get services from Kubernetes
    services, err := kp.client.CoreV1().Services(kp.namespace).List(ctx, metav1.ListOptions{
        LabelSelector: labels.Set{"app": serviceName}.AsSelector().String(),
    })
    if err != nil {
        return nil, err
    }
    
    serviceInfos := make([]*ServiceInfo, 0, len(services.Items))
    for _, service := range services.Items {
        // Get endpoints for the service
        endpoints, err := kp.client.CoreV1().Endpoints(kp.namespace).Get(ctx, service.Name, metav1.GetOptions{})
        if err != nil {
            continue
        }
        
        for _, subset := range endpoints.Subsets {
            for _, address := range subset.Addresses {
                for _, port := range subset.Ports {
                    serviceInfo := &ServiceInfo{
                        ID:       fmt.Sprintf("%s-%s", service.Name, address.IP),
                        Name:     service.Name,
                        Address:  address.IP,
                        Port:     int(port.Port),
                        Tags:     []string{"kubernetes"},
                        Meta:     service.Annotations,
                        Status:   StatusHealthy,
                        LastSeen: time.Now(),
                    }
                    
                    serviceInfos = append(serviceInfos, serviceInfo)
                }
            }
        }
    }
    
    return serviceInfos, nil
}

func (kp *KubernetesProvider) WatchServices(ctx context.Context, serviceName string, callback func([]*ServiceInfo)) error {
    // Use Kubernetes watch functionality
    go func() {
        watch, err := kp.client.CoreV1().Services(kp.namespace).Watch(ctx, metav1.ListOptions{
            LabelSelector: labels.Set{"app": serviceName}.AsSelector().String(),
        })
        if err != nil {
            kp.logger.Error("Failed to watch services", "service", serviceName, "error", err)
            return
        }
        
        for event := range watch.ResultChan() {
            // Get current services
            services, err := kp.DiscoverServices(ctx, serviceName)
            if err != nil {
                kp.logger.Error("Failed to discover services", "error", err)
                continue
            }
            
            // Call callback
            callback(services)
        }
    }()
    
    return nil
}

func (kp *KubernetesProvider) HealthCheck(ctx context.Context, serviceID string) error {
    // Check if service exists in Kubernetes
    services, err := kp.client.CoreV1().Services(kp.namespace).List(ctx, metav1.ListOptions{})
    if err != nil {
        return err
    }
    
    for _, service := range services.Items {
        if service.Name == serviceID {
            return nil
        }
    }
    
    return fmt.Errorf("service %s not found", serviceID)
}
```

### 6. Static Provider Implementation

```go
// internal/discovery/static_provider.go
package discovery

import (
    "context"
    "fmt"
    "time"
    
    "github.com/anasamu/go-micro-libs/logging"
)

type StaticProvider struct {
    services map[string][]*ServiceInfo
    logger   *logging.Logger
}

func NewStaticProvider(services map[string][]*ServiceInfo, logger *logging.Logger) *StaticProvider {
    return &StaticProvider{
        services: services,
        logger:   logger,
    }
}

func (sp *StaticProvider) RegisterService(ctx context.Context, service *ServiceInfo) error {
    // In static provider, services are pre-configured
    // We'll just log the registration
    
    sp.logger.Info("Service registered with static provider",
        "service", service.Name,
        "address", service.Address,
        "port", service.Port,
    )
    
    return nil
}

func (sp *StaticProvider) DeregisterService(ctx context.Context, serviceID string) error {
    // In static provider, services are pre-configured
    // We'll just log the deregistration
    
    sp.logger.Info("Service deregistered from static provider", "service_id", serviceID)
    
    return nil
}

func (sp *StaticProvider) DiscoverServices(ctx context.Context, serviceName string) ([]*ServiceInfo, error) {
    services, exists := sp.services[serviceName]
    if !exists {
        return nil, fmt.Errorf("service %s not found", serviceName)
    }
    
    return services, nil
}

func (sp *StaticProvider) WatchServices(ctx context.Context, serviceName string, callback func([]*ServiceInfo)) error {
    // Static provider doesn't support watching
    // We'll just call the callback once with current services
    
    services, err := sp.DiscoverServices(ctx, serviceName)
    if err != nil {
        return err
    }
    
    callback(services)
    
    return nil
}

func (sp *StaticProvider) HealthCheck(ctx context.Context, serviceID string) error {
    // In static provider, all services are considered healthy
    return nil
}
```

## 🔧 Best Practices

### 1. Service Discovery Configuration

```go
// Optimal service discovery configuration
func getOptimalDiscoveryConfig() *Config {
    return &Config{
        ServiceName:    "user-service",
        ServicePort:    8080,
        ServiceAddress: "localhost",
        ServiceTags:    []string{"api", "user", "v1"},
        ServiceMeta: map[string]string{
            "version":     "1.0.0",
            "environment": "production",
        },
        HealthCheck: HealthCheckConfig{
            Enabled:         true,
            Path:           "/health",
            Interval:       10 * time.Second,
            Timeout:        3 * time.Second,
            DeregisterAfter: 30 * time.Second,
            HTTPCheck: HTTPCheckConfig{
                Enabled: true,
                Path:    "/health",
                Method:  "GET",
                Headers: map[string]string{
                    "Content-Type": "application/json",
                },
            },
        },
        Watch: WatchConfig{
            Enabled:  true,
            Services: []string{"user-service", "order-service", "payment-service"},
        },
    }
}
```

### 2. Service Discovery Monitoring

```go
// Monitor service discovery events
func (m *Manager) MonitorDiscovery() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        m.mutex.RLock()
        for name, watcher := range m.watchers {
            // Log service discovery status
            m.logger.Info("Service discovery status",
                "service", name,
                "watcher_active", true,
            )
        }
        m.mutex.RUnlock()
    }
}
```

### 3. Service Discovery Testing

```go
// Test service discovery behavior
func TestServiceDiscovery(t *testing.T) {
    config := &Config{
        ServiceName:    "test-service",
        ServicePort:    8080,
        ServiceAddress: "localhost",
        ServiceTags:    []string{"test"},
    }
    
    manager := NewManager(config, logger, metrics)
    
    // Test service registration
    err := manager.RegisterService(ctx, "test-provider")
    assert.NoError(t, err)
    
    // Test service discovery
    services, err := manager.DiscoverServices(ctx, "test-provider", "test-service")
    assert.NoError(t, err)
    assert.NotEmpty(t, services)
    
    // Test service deregistration
    err = manager.DeregisterService(ctx, "test-provider", "test-service-id")
    assert.NoError(t, err)
}
```

### 4. Service Discovery Health Checks

```go
// Implement health checks for service discovery
func (m *Manager) HealthCheck(ctx context.Context) error {
    // Check if all providers are healthy
    for name, provider := range m.providers {
        err := provider.HealthCheck(ctx, "health-check")
        if err != nil {
            m.logger.Error("Provider health check failed", "provider", name, "error", err)
            return err
        }
    }
    
    return nil
}
```

---

**Service Discovery - Dynamic and reliable service discovery for microservices! 🚀**

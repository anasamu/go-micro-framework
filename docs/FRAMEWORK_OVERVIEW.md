# GoMicroFramework - Framework Overview

## 🎯 Introduction

GoMicroFramework adalah framework Go yang powerful dan user-friendly untuk pengembangan microservices dengan integrasi seamless dari 20+ library yang sudah ada di [go-micro-libs](https://github.com/anasamu/go-micro-libs/). Framework ini dirancang untuk memberikan developer experience yang optimal dengan fokus pada business logic, bukan infrastructure.

## 🏗️ Architecture Overview

### Core Philosophy

1. **Zero-Configuration Setup**: Framework menangani semua konfigurasi default
2. **Business Logic Focus**: Developer fokus pada business logic, bukan infrastructure
3. **Production Ready**: Built-in monitoring, logging, security, dan resilience
4. **Extensible**: Mudah menambah fitur baru dan custom providers
5. **Library Integration**: Menggunakan semua library yang sudah ada di go-micro-libs

### Framework Structure

```
go-micro-framework/
├── cmd/
│   └── microframework/          # CLI binary
├── internal/
│   ├── core/                    # Core framework logic
│   ├── generator/               # Code generation engine
│   └── templates/               # Go templates for code gen
├── pkg/                         # Framework packages
├── docs/                        # Documentation
├── examples/                    # Generated examples
└── scripts/                     # Build and deployment scripts
```

## 🔧 Core Components

### 1. CLI Tool (`microframework`)

CLI tool yang powerful untuk:
- Generate new microservices
- Add features to existing services
- Manage configuration
- Deploy services
- Monitor and debug services

### 2. Service Generator

Engine yang menghasilkan:
- Complete service structure
- Configuration files
- Docker and Kubernetes manifests
- Tests and documentation
- Integration with go-micro-libs

### 3. Bootstrap Engine

Framework yang menangani:
- Service initialization
- Library integration
- Configuration management
- Health checks
- Graceful shutdown

## 📚 Library Integration

### Core Libraries (Always Integrated)

Framework ini secara otomatis mengintegrasikan library berikut di setiap service:

| Library | Description | Provider |
|---------|-------------|----------|
| **Config** | Configuration management | File, Env, Consul, Vault |
| **Logging** | Structured logging | Console, File, Elasticsearch |
| **Monitoring** | Metrics and tracing | Prometheus, Jaeger, Grafana |
| **Middleware** | HTTP middleware | Auth, Rate Limit, Circuit Breaker |
| **Communication** | Communication protocols | HTTP, gRPC, WebSocket, GraphQL |
| **Utils** | Utility functions | UUID, Environment, Validation |

### Optional Libraries

| Library | Description | Provider |
|---------|-------------|----------|
| **AI** | AI services | OpenAI, Anthropic, Google, DeepSeek, X.AI |
| **Auth** | Authentication | JWT, OAuth2, LDAP, SAML, 2FA |
| **Database** | Database abstraction | PostgreSQL, MySQL, MongoDB, Redis |
| **Cache** | Caching system | Redis, Memcached, Memory |
| **Storage** | Object storage | S3, GCS, Azure, MinIO |
| **Messaging** | Message queues | Kafka, RabbitMQ, NATS, SQS |
| **Discovery** | Service discovery | Consul, Kubernetes, etcd |
| **Event** | Event sourcing | PostgreSQL, Kafka, NATS |
| **Payment** | Payment processing | Stripe, PayPal, Midtrans, Xendit |
| **Backup** | Backup services | S3, GCS, Local |
| **Chaos** | Chaos engineering | Kubernetes, HTTP, Messaging |
| **Circuit Breaker** | Resilience patterns | Custom, GoBreaker |
| **Failover** | Failover mechanisms | Consul, Kubernetes |
| **FileGen** | File generation | DOCX, Excel, CSV, PDF |
| **Rate Limit** | Rate limiting | Token bucket, Sliding window |
| **Scheduling** | Task scheduling | Cron, Redis-based |
| **API** | Third-party APIs | HTTP, gRPC, GraphQL, WebSocket |
| **Email** | Email services | SMTP, SendGrid, Mailgun |

## 🚀 Service Types

Framework mendukung berbagai jenis service:

### 1. REST API Service
- Standard HTTP REST service
- JSON API endpoints
- Built-in middleware support
- OpenAPI documentation

### 2. GraphQL Service
- GraphQL API service
- Schema generation
- Resolver implementation
- Subscription support

### 3. gRPC Service
- High-performance gRPC service
- Protocol buffer support
- Streaming support
- Service reflection

### 4. WebSocket Service
- Real-time WebSocket service
- Bidirectional communication
- Connection management
- Message broadcasting

### 5. Event-Driven Service
- Event sourcing service
- Event store integration
- Command/Query separation
- Event replay

### 6. Scheduled Service
- Cron/scheduled task service
- Task scheduling
- Job management
- Retry mechanisms

### 7. Worker Service
- Background job processing
- Queue integration
- Job distribution
- Error handling

### 8. Gateway Service
- API Gateway service
- Request routing
- Load balancing
- Rate limiting

### 9. Proxy Service
- Reverse proxy service
- Request forwarding
- Response caching
- Load balancing

## 🔧 Configuration Management

### Multi-Source Configuration

Framework mendukung konfigurasi dari berbagai sumber:

1. **Environment Variables**
2. **YAML Files**
3. **Consul**
4. **Vault**
5. **Command Line Flags**

### Configuration Hierarchy

```
Command Line Flags > Environment Variables > Config Files > Defaults
```

### Hot Reloading

Framework mendukung hot reloading untuk:
- Configuration changes
- Service discovery updates
- Feature flags
- Environment variables

## 🏗️ Generated Service Structure

Setiap service yang dibuat memiliki struktur yang konsisten:

```
service-name/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── handlers/               # HTTP/gRPC handlers
│   ├── services/               # Business logic
│   ├── repositories/           # Data access layer
│   ├── models/                 # Data models
│   └── middleware/             # Custom middleware
├── pkg/
│   └── types/                  # Public types
├── configs/
│   ├── config.yaml             # Default configuration
│   ├── config.dev.yaml         # Development config
│   └── config.prod.yaml        # Production config
├── deployments/
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── kubernetes/
│       ├── deployment.yaml
│       └── service.yaml
├── tests/
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── docs/
│   ├── api.md
│   └── deployment.md
├── scripts/
│   ├── build.sh
│   ├── test.sh
│   └── deploy.sh
├── go.mod                      # Dependencies
├── go.sum
├── Makefile
└── README.md
```

## 🔒 Security Features

### Built-in Security

1. **Authentication**
   - JWT token validation
   - OAuth2 integration
   - Multi-factor authentication
   - Session management

2. **Authorization**
   - Role-based access control (RBAC)
   - Attribute-based access control (ABAC)
   - Access control lists (ACL)
   - Permission management

3. **Security Headers**
   - CORS configuration
   - CSRF protection
   - XSS protection
   - Content Security Policy

4. **Rate Limiting**
   - Request rate limiting
   - IP-based limiting
   - User-based limiting
   - API key limiting

## 📊 Monitoring & Observability

### Built-in Monitoring

1. **Metrics**
   - Prometheus metrics
   - Custom metrics
   - Business metrics
   - Performance metrics

2. **Tracing**
   - Distributed tracing
   - Request correlation
   - Performance profiling
   - Error tracking

3. **Logging**
   - Structured logging
   - Correlation IDs
   - Log levels
   - Log aggregation

4. **Health Checks**
   - Liveness probes
   - Readiness probes
   - Dependency checks
   - Custom health checks

## 🧪 Testing Support

### Test Types

1. **Unit Tests**
   - Business logic testing
   - Mock integration
   - Coverage reporting
   - Fast execution

2. **Integration Tests**
   - Database testing
   - External service testing
   - End-to-end workflows
   - Performance testing

3. **End-to-End Tests**
   - Full system testing
   - User journey testing
   - API testing
   - UI testing

### Test Utilities

- Test data generation
- Mock services
- Test containers
- Performance benchmarks

## 🚀 Deployment Support

### Containerization

1. **Docker**
   - Multi-stage builds
   - Optimized images
   - Security scanning
   - Image signing

2. **Kubernetes**
   - Deployment manifests
   - Service definitions
   - ConfigMaps and Secrets
   - Ingress configuration

3. **Helm**
   - Chart generation
   - Value templates
   - Dependency management
   - Release management

### Cloud Deployment

- AWS EKS
- Google GKE
- Azure AKS
- DigitalOcean Kubernetes

## 🔄 Development Workflow

### 1. Service Generation
```bash
microframework new user-service --with-auth=jwt --with-database=postgres
```

### 2. Development
```bash
cd user-service
go run cmd/main.go
```

### 3. Testing
```bash
make test
make test-integration
make test-e2e
```

### 4. Deployment
```bash
make docker-build
make k8s-deploy
```

## 🎯 Best Practices

### 1. Service Design
- Single responsibility principle
- Domain-driven design
- Clean architecture
- Microservice patterns

### 2. Configuration
- Environment-specific configs
- Secret management
- Configuration validation
- Hot reloading

### 3. Monitoring
- Comprehensive metrics
- Distributed tracing
- Structured logging
- Health checks

### 4. Security
- Authentication and authorization
- Input validation
- Rate limiting
- Security headers

### 5. Testing
- Test-driven development
- Comprehensive test coverage
- Integration testing
- Performance testing

## 🔮 Future Roadmap

### Planned Features

1. **Multi-language Support**
   - Python wrappers
   - Node.js wrappers
   - Java wrappers

2. **Cloud Native**
   - Native cloud provider integration
   - Serverless support
   - Edge computing

3. **Advanced Features**
   - Blockchain integration
   - IoT device support
   - Machine learning pipelines

4. **Enterprise Features**
   - Enterprise SSO
   - Advanced compliance
   - Policy governance
   - Multi-tenancy

## 🤝 Community & Support

### Getting Help

1. **Documentation**: Comprehensive guides and examples
2. **GitHub Issues**: Bug reports and feature requests
3. **Discussions**: Community discussions and Q&A
4. **Discord/Slack**: Real-time community chat

### Contributing

1. **Code Contributions**: Bug fixes and new features
2. **Documentation**: Improving guides and examples
3. **Testing**: Adding test cases and improving coverage
4. **Community**: Helping other developers

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**GoMicroFramework - Building the future of Go microservices development! 🚀**

# Examples

This directory contains examples of services generated using Go Micro Framework.

## Available Examples

### 1. Basic REST Service

A simple REST API service with basic CRUD operations.

**Generated with:**
```bash
microframework new basic-rest-service \
  --type=rest \
  --with-auth=jwt \
  --with-database=postgresql
```

**Features:**
- REST API with CRUD operations
- JWT authentication
- PostgreSQL database
- Health checks
- Basic monitoring

### 2. Event-Driven Service

An event-driven service using Kafka for messaging.

**Generated with:**
```bash
microframework new event-driven-service \
  --type=event \
  --with-messaging=kafka \
  --with-database=postgresql \
  --with-event=postgresql
```

**Features:**
- Event sourcing with PostgreSQL
- Kafka messaging
- Event handlers
- CQRS pattern
- Event store

### 3. AI-Powered Service

A service that integrates AI capabilities for intelligent processing.

**Generated with:**
```bash
microframework new ai-powered-service \
  --type=rest \
  --with-ai=openai \
  --with-database=postgresql \
  --with-cache=redis
```

**Features:**
- OpenAI integration
- Intelligent content processing
- Caching with Redis
- AI-powered recommendations

### 4. Payment Service

A complete payment processing service.

**Generated with:**
```bash
microframework new payment-service \
  --type=rest \
  --with-payment=stripe \
  --with-database=postgresql \
  --with-monitoring=prometheus
```

**Features:**
- Stripe payment integration
- Payment processing
- Webhook handling
- Transaction management
- Comprehensive monitoring

### 5. E-commerce Service

A complete e-commerce microservice with all features.

**Generated with:**
```bash
microframework new ecommerce-service \
  --type=rest \
  --with-auth=jwt \
  --with-database=postgresql \
  --with-cache=redis \
  --with-messaging=kafka \
  --with-storage=s3 \
  --with-payment=stripe \
  --with-monitoring=prometheus \
  --with-ai=openai \
  --optional=backup,chaos,circuitbreaker \
  --middleware=auth,logging,monitoring,ratelimit \
  --deployment=kubernetes \
  --testing=unit,integration,e2e
```

**Features:**
- Complete e-commerce functionality
- All library integrations
- Production-ready deployment
- Comprehensive testing
- Advanced monitoring

### 6. Microservices Architecture

A complete microservices architecture example.

**Services:**
- User Service
- Order Service
- Payment Service
- Notification Service
- API Gateway

**Generated with:**
```bash
# User Service
microframework new user-service \
  --type=rest \
  --with-auth=jwt \
  --with-database=postgresql \
  --with-cache=redis \
  --with-monitoring=prometheus

# Order Service
microframework new order-service \
  --type=rest \
  --with-database=postgresql \
  --with-messaging=kafka \
  --with-monitoring=prometheus

# Payment Service
microframework new payment-service \
  --type=rest \
  --with-payment=stripe \
  --with-database=postgresql \
  --with-monitoring=prometheus

# Notification Service
microframework new notification-service \
  --type=worker \
  --with-messaging=kafka \
  --with-ai=openai \
  --with-monitoring=prometheus

# API Gateway
microframework new api-gateway \
  --type=gateway \
  --with-auth=jwt \
  --with-monitoring=prometheus \
  --with-discovery=consul
```

## Running Examples

### Prerequisites

1. **Install Go Micro Framework:**
   ```bash
   go install github.com/anasamu/go-micro-framework/cmd/microframework@latest
   ```

2. **Install Dependencies:**
   - Docker and Docker Compose
   - Kubernetes (optional)
   - Required environment variables

### Basic Setup

1. **Clone the examples:**
   ```bash
   git clone https://github.com/anasamu/go-micro-framework.git
   cd go-micro-framework/examples
   ```

2. **Choose an example:**
   ```bash
   cd basic-rest-service
   ```

3. **Install dependencies:**
   ```bash
   go mod tidy
   ```

4. **Setup environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

5. **Run the service:**
   ```bash
   go run cmd/main.go
   ```

### Docker Setup

1. **Build and run with Docker:**
   ```bash
   docker-compose up -d
   ```

2. **Check service health:**
   ```bash
   curl http://localhost:8080/health
   ```

### Kubernetes Setup

1. **Deploy to Kubernetes:**
   ```bash
   kubectl apply -f deployments/kubernetes/
   ```

2. **Check deployment:**
   ```bash
   kubectl get pods
   kubectl get services
   ```

## Example Configurations

### Environment Variables

```bash
# Service Configuration
SERVICE_NAME=basic-rest-service
SERVICE_VERSION=1.0.0
SERVICE_PORT=8080

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/mydb

# Authentication
JWT_SECRET=your-jwt-secret

# Monitoring
PROMETHEUS_ENDPOINT=http://localhost:9090
JAEGER_ENDPOINT=http://localhost:14268

# AI Services (for AI-powered examples)
OPENAI_API_KEY=your-openai-api-key

# Payment (for payment examples)
STRIPE_SECRET_KEY=your-stripe-secret-key
STRIPE_WEBHOOK_SECRET=your-stripe-webhook-secret

# Storage (for e-commerce examples)
AWS_ACCESS_KEY_ID=your-aws-access-key
AWS_SECRET_ACCESS_KEY=your-aws-secret-key
S3_BUCKET=your-s3-bucket

# Messaging (for event-driven examples)
KAFKA_BROKERS=localhost:9092
```

### Docker Compose

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/mydb
      - JWT_SECRET=your-jwt-secret
    depends_on:
      - db
      - redis

  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=mydb
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  jaeger:
    image: jaegertracing/all-in-one
    ports:
      - "16686:16686"
      - "14268:14268"

volumes:
  postgres_data:
```

## Testing Examples

### Unit Tests

```bash
# Run unit tests
make test-unit

# Run with coverage
make test-coverage
```

### Integration Tests

```bash
# Run integration tests
make test-integration

# Run with Docker
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

### End-to-End Tests

```bash
# Run E2E tests
make test-e2e

# Run with specific environment
./scripts/test.sh --pattern E2E --env staging
```

## API Testing

### Health Check

```bash
curl http://localhost:8080/health
```

### Authentication

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password"}'

# Use token
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/users
```

### CRUD Operations

```bash
# Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name": "John Doe", "email": "john@example.com"}'

# Get user
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/users/1

# Update user
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name": "Jane Doe", "email": "jane@example.com"}'

# Delete user
curl -X DELETE -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/users/1
```

## Monitoring Examples

### Prometheus Metrics

```bash
curl http://localhost:8080/metrics
```

### Jaeger Tracing

Visit http://localhost:16686 to view traces.

### Grafana Dashboard

Visit http://localhost:3000 to view dashboards.

## Troubleshooting

### Common Issues

1. **Database Connection Issues:**
   - Check DATABASE_URL format
   - Ensure database is running
   - Verify network connectivity

2. **Authentication Issues:**
   - Check JWT_SECRET is set
   - Verify token format
   - Check token expiration

3. **Monitoring Issues:**
   - Ensure Prometheus is running
   - Check endpoint accessibility
   - Verify configuration

### Debug Mode

```bash
# Run with debug logging
LOG_LEVEL=debug go run cmd/main.go

# Run with verbose output
go run cmd/main.go --verbose
```

### Logs

```bash
# View logs
microframework logs --service=basic-rest-service

# Follow logs
microframework logs --service=basic-rest-service --follow

# Filter by level
microframework logs --service=basic-rest-service --level=error
```

## Contributing Examples

We welcome contributions of new examples! Please follow these guidelines:

1. **Create a new directory** for your example
2. **Include a README.md** with setup instructions
3. **Provide configuration files** (Docker, Kubernetes, etc.)
4. **Add tests** for your example
5. **Document the features** and use cases

### Example Structure

```
examples/
├── your-example/
│   ├── README.md
│   ├── cmd/
│   ├── internal/
│   ├── configs/
│   ├── deployments/
│   ├── tests/
│   ├── docker-compose.yml
│   └── Makefile
```

## Support

For questions about examples:
- Check the documentation
- Open an issue on GitHub
- Join our Discord community

Happy coding! 🚀

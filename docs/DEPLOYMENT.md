# Deployment Guide

## Overview

This guide covers deploying services generated with Go Micro Framework to various platforms and environments.

## Deployment Options

### 1. Docker Deployment

#### Basic Docker Deployment

```bash
# Build Docker image
docker build -t user-service:latest .

# Run container
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@host:5432/db" \
  -e JWT_SECRET="your-secret" \
  user-service:latest
```

#### Docker Compose

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
      - PROMETHEUS_ENDPOINT=http://prometheus:9090
    depends_on:
      - db
      - redis
      - prometheus

  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=mydb
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

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

#### Multi-stage Dockerfile

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary and configs
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

RUN chown -R appuser:appgroup /root

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]
```

### 2. Kubernetes Deployment

#### Basic Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  labels:
    app: user-service
    version: v1.0.0
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
        version: v1.0.0
    spec:
      serviceAccountName: user-service
      securityContext:
        runAsNonRoot: true
        runAsUser: 1001
        fsGroup: 1001
      containers:
      - name: user-service
        image: user-service:latest
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: metrics
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: user-service-secrets
              key: database-url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: user-service-secrets
              key: jwt-secret
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
            - ALL
```

#### Kubernetes Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: user-service
  labels:
    app: user-service
spec:
  selector:
    app: user-service
  ports:
  - name: http
    port: 80
    targetPort: 8080
  - name: metrics
    port: 9090
    targetPort: 9090
  type: ClusterIP
```

#### Kubernetes ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: user-service-config
data:
  config.yaml: |
    service:
      name: "user-service"
      version: "1.0.0"
      port: 8080
    
    logging:
      providers:
        console:
          level: "info"
          format: "json"
    
    monitoring:
      providers:
        prometheus:
          endpoint: "http://prometheus:9090"
```

#### Kubernetes Secrets

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: user-service-secrets
type: Opaque
data:
  database-url: cG9zdGdyZXM6Ly91c2VyOnBhc3NAaG9zdDo1NDMyL2Ri
  jwt-secret: eW91ci1qd3Qtc2VjcmV0
```

#### Kubernetes Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: user-service-ingress
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
  - hosts:
    - user-service.example.com
    secretName: user-service-tls
  rules:
  - host: user-service.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: user-service
            port:
              number: 80
```

### 3. Helm Deployment

#### Helm Chart Structure

```
charts/user-service/
├── Chart.yaml
├── values.yaml
├── templates/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── configmap.yaml
│   ├── secret.yaml
│   ├── ingress.yaml
│   └── _helpers.tpl
└── README.md
```

#### Chart.yaml

```yaml
apiVersion: v2
name: user-service
description: User Service Helm Chart
type: application
version: 1.0.0
appVersion: "1.0.0"
keywords:
  - microservice
  - user
  - api
home: https://github.com/anasamu/go-micro-framework
sources:
  - https://github.com/anasamu/go-micro-framework
maintainers:
  - name: Go Micro Framework Team
    email: team@example.com
```

#### values.yaml

```yaml
replicaCount: 3

image:
  repository: user-service
  tag: latest
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80
  targetPort: 8080

ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: user-service.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: user-service-tls
      hosts:
        - user-service.example.com

resources:
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 100m
    memory: 128Mi

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
  targetMemoryUtilizationPercentage: 80

nodeSelector: {}

tolerations: []

affinity: {}

config:
  service:
    name: "user-service"
    version: "1.0.0"
    port: 8080
  
  logging:
    providers:
      console:
        level: "info"
        format: "json"
  
  monitoring:
    providers:
      prometheus:
        endpoint: "http://prometheus:9090"

secrets:
  database-url: "postgres://user:pass@host:5432/db"
  jwt-secret: "your-jwt-secret"
```

#### Deploy with Helm

```bash
# Add Helm repository
helm repo add go-micro-framework https://charts.example.com
helm repo update

# Install chart
helm install user-service go-micro-framework/user-service \
  --namespace production \
  --create-namespace \
  --values values.yaml

# Upgrade chart
helm upgrade user-service go-micro-framework/user-service \
  --namespace production \
  --values values.yaml

# Uninstall chart
helm uninstall user-service --namespace production
```

### 4. Cloud Provider Deployment

#### AWS ECS Deployment

```yaml
# task-definition.json
{
  "family": "user-service",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "512",
  "memory": "1024",
  "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::account:role/ecsTaskRole",
  "containerDefinitions": [
    {
      "name": "user-service",
      "image": "account.dkr.ecr.region.amazonaws.com/user-service:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "SERVICE_NAME",
          "value": "user-service"
        }
      ],
      "secrets": [
        {
          "name": "DATABASE_URL",
          "valueFrom": "arn:aws:ssm:region:account:parameter/user-service/database-url"
        },
        {
          "name": "JWT_SECRET",
          "valueFrom": "arn:aws:ssm:region:account:parameter/user-service/jwt-secret"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/user-service",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      },
      "healthCheck": {
        "command": [
          "CMD-SHELL",
          "wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1"
        ],
        "interval": 30,
        "timeout": 5,
        "retries": 3,
        "startPeriod": 60
      }
    }
  ]
}
```

#### Google Cloud Run

```yaml
# cloud-run.yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: user-service
  annotations:
    run.googleapis.com/ingress: all
    run.googleapis.com/execution-environment: gen2
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/maxScale: "10"
        run.googleapis.com/cpu-throttling: "false"
        run.googleapis.com/execution-environment: gen2
    spec:
      containerConcurrency: 100
      timeoutSeconds: 300
      containers:
      - image: gcr.io/project/user-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: user-service-secrets
              key: database-url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: user-service-secrets
              key: jwt-secret
        resources:
          limits:
            cpu: "2"
            memory: "2Gi"
          requests:
            cpu: "1"
            memory: "1Gi"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

#### Azure Container Instances

```yaml
# azure-container-instance.yaml
apiVersion: 2019-12-01
location: eastus
name: user-service
properties:
  containers:
  - name: user-service
    properties:
      image: user-service.azurecr.io/user-service:latest
      ports:
      - port: 8080
        protocol: TCP
      environmentVariables:
      - name: DATABASE_URL
        secureValue: postgres://user:pass@host:5432/db
      - name: JWT_SECRET
        secureValue: your-jwt-secret
      resources:
        requests:
          cpu: 1
          memoryInGb: 1
        limits:
          cpu: 2
          memoryInGb: 2
      livenessProbe:
        httpGet:
          path: /health/live
          port: 8080
        initialDelaySeconds: 30
        periodSeconds: 10
      readinessProbe:
        httpGet:
          path: /health/ready
          port: 8080
        initialDelaySeconds: 5
        periodSeconds: 5
  osType: Linux
  ipAddress:
    type: Public
    ports:
    - protocol: TCP
      port: 8080
  restartPolicy: Always
type: Microsoft.ContainerInstance/containerGroups
```

## Environment-Specific Configurations

### Development Environment

```yaml
# config.dev.yaml
service:
  name: "user-service"
  version: "1.0.0"
  port: 8080
  environment: "development"

logging:
  providers:
    console:
      level: "debug"
      format: "text"

database:
  providers:
    postgresql:
      url: "postgres://user:pass@localhost:5432/mydb_dev"
      max_connections: 10

monitoring:
  providers:
    prometheus:
      endpoint: "http://localhost:9090"
```

### Staging Environment

```yaml
# config.staging.yaml
service:
  name: "user-service"
  version: "1.0.0"
  port: 8080
  environment: "staging"

logging:
  providers:
    console:
      level: "info"
      format: "json"
    file:
      path: "/var/log/user-service.log"
      level: "info"

database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 50

monitoring:
  providers:
    prometheus:
      endpoint: "http://prometheus-staging:9090"
    jaeger:
      endpoint: "http://jaeger-staging:14268"
```

### Production Environment

```yaml
# config.prod.yaml
service:
  name: "user-service"
  version: "1.0.0"
  port: 8080
  environment: "production"

logging:
  providers:
    console:
      level: "warn"
      format: "json"
    file:
      path: "/var/log/user-service.log"
      level: "info"
      max_size: 100
      max_backups: 3
      max_age: 28
    elasticsearch:
      endpoint: "${ELASTICSEARCH_ENDPOINT}"
      index: "user-service-logs"

database:
  providers:
    postgresql:
      url: "${DATABASE_URL}"
      max_connections: 100
      max_idle_connections: 10
      connection_max_lifetime: "1h"

monitoring:
  providers:
    prometheus:
      endpoint: "${PROMETHEUS_ENDPOINT}"
    jaeger:
      endpoint: "${JAEGER_ENDPOINT}"
    grafana:
      endpoint: "${GRAFANA_ENDPOINT}"

middleware:
  rate_limit:
    enabled: true
    requests_per_minute: 1000
  circuit_breaker:
    enabled: true
    failure_threshold: 5
    timeout: 30s
```

## CI/CD Pipeline

### GitHub Actions

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Run tests
      run: make test
    
    - name: Run integration tests
      run: make test-integration
    
    - name: Run security scan
      run: make security

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2
    
    - name: Log in to Container Registry
      uses: docker/login-action@v2
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}
    
    - name: Extract metadata
      id: meta
      uses: docker/metadata-action@v4
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
        tags: |
          type=ref,event=branch
          type=ref,event=pr
          type=semver,pattern={{version}}
          type=semver,pattern={{major}}.{{minor}}
    
    - name: Build and push Docker image
      uses: docker/build-push-action@v4
      with:
        context: .
        push: true
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}

  deploy-staging:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
    - uses: actions/checkout@v3
    
    - name: Deploy to staging
      run: |
        helm upgrade --install user-service-staging ./charts/user-service \
          --namespace staging \
          --create-namespace \
          --values charts/user-service/values.staging.yaml \
          --set image.tag=${{ github.sha }}

  deploy-production:
    needs: deploy-staging
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment: production
    steps:
    - uses: actions/checkout@v3
    
    - name: Deploy to production
      run: |
        helm upgrade --install user-service ./charts/user-service \
          --namespace production \
          --create-namespace \
          --values charts/user-service/values.prod.yaml \
          --set image.tag=${{ github.sha }}
```

### GitLab CI

```yaml
stages:
  - test
  - build
  - deploy-staging
  - deploy-production

variables:
  DOCKER_DRIVER: overlay2
  DOCKER_TLS_CERTDIR: "/certs"

test:
  stage: test
  image: golang:1.21
  script:
    - make test
    - make test-integration
    - make security
  coverage: '/coverage: \d+\.\d+%/'

build:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
  only:
    - main

deploy-staging:
  stage: deploy-staging
  image: alpine/helm:latest
  script:
    - helm upgrade --install user-service-staging ./charts/user-service \
        --namespace staging \
        --create-namespace \
        --values charts/user-service/values.staging.yaml \
        --set image.tag=$CI_COMMIT_SHA
  only:
    - main

deploy-production:
  stage: deploy-production
  image: alpine/helm:latest
  script:
    - helm upgrade --install user-service ./charts/user-service \
        --namespace production \
        --create-namespace \
        --values charts/user-service/values.prod.yaml \
        --set image.tag=$CI_COMMIT_SHA
  only:
    - main
  when: manual
```

## Monitoring and Observability

### Prometheus Configuration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "rules/*.yml"

scrape_configs:
  - job_name: 'user-service'
    static_configs:
      - targets: ['user-service:9090']
    scrape_interval: 5s
    metrics_path: /metrics

  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093
```

### Grafana Dashboard

```json
{
  "dashboard": {
    "title": "User Service Dashboard",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{path}}"
          }
        ]
      },
      {
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      },
      {
        "title": "Error Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total{status=~\"5..\"}[5m])",
            "legendFormat": "5xx errors"
          }
        ]
      }
    ]
  }
}
```

## Security Considerations

### Container Security

```dockerfile
# Use non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

USER appuser

# Use read-only root filesystem
securityContext:
  readOnlyRootFilesystem: true

# Drop all capabilities
securityContext:
  capabilities:
    drop:
    - ALL
```

### Network Security

```yaml
# NetworkPolicy
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: user-service-netpol
spec:
  podSelector:
    matchLabels:
      app: user-service
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - namespaceSelector:
        matchLabels:
          name: database
    ports:
    - protocol: TCP
      port: 5432
```

### Secret Management

```yaml
# External Secrets Operator
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: aws-secrets-manager
spec:
  provider:
    aws:
      service: SecretsManager
      region: us-east-1
      auth:
        secretRef:
          accessKeyID:
            name: aws-credentials
            key: access-key-id
          secretAccessKey:
            name: aws-credentials
            key: secret-access-key

---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: user-service-secrets
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: aws-secrets-manager
    kind: SecretStore
  target:
    name: user-service-secrets
    creationPolicy: Owner
  data:
  - secretKey: database-url
    remoteRef:
      key: user-service/database-url
  - secretKey: jwt-secret
    remoteRef:
      key: user-service/jwt-secret
```

## Troubleshooting

### Common Issues

1. **Service won't start:**
   - Check logs: `kubectl logs -f deployment/user-service`
   - Verify configuration: `kubectl describe configmap user-service-config`
   - Check secrets: `kubectl describe secret user-service-secrets`

2. **Database connection issues:**
   - Verify DATABASE_URL format
   - Check network connectivity
   - Verify credentials

3. **Health check failures:**
   - Check health endpoint: `curl http://localhost:8080/health`
   - Verify dependencies are running
   - Check resource limits

### Debug Commands

```bash
# Check pod status
kubectl get pods -l app=user-service

# Check service endpoints
kubectl get endpoints user-service

# Check ingress
kubectl describe ingress user-service-ingress

# Check logs
kubectl logs -f deployment/user-service

# Port forward for local testing
kubectl port-forward svc/user-service 8080:80

# Check resource usage
kubectl top pods -l app=user-service

# Check events
kubectl get events --sort-by=.metadata.creationTimestamp
```

## Best Practices

1. **Use specific image tags** instead of `latest`
2. **Set resource limits** to prevent resource exhaustion
3. **Use health checks** for better reliability
4. **Implement proper logging** for debugging
5. **Use secrets management** for sensitive data
6. **Enable monitoring** for observability
7. **Use network policies** for security
8. **Implement proper backup** strategies
9. **Use blue-green deployments** for zero downtime
10. **Test in staging** before production deployment

## Conclusion

This deployment guide provides comprehensive instructions for deploying Go Micro Framework services to various platforms. Choose the deployment method that best fits your infrastructure and requirements.

For more information, refer to the [Architecture Documentation](ARCHITECTURE.md) and [API Documentation](API.md).

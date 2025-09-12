#!/bin/bash

# Go Micro Framework Deploy Script
# This script deploys the framework to various environments

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variables
BINARY_NAME="microframework"
BUILD_DIR="build"
DIST_DIR="dist"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
DOCKER_IMAGE="microframework"
DOCKER_TAG="${VERSION}"

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Build the binary
build_binary() {
    log_info "Building binary..."
    go build -o ${BUILD_DIR}/${BINARY_NAME} cmd/microframework/main.go
    log_success "Binary built successfully"
}

# Build Docker image
build_docker_image() {
    log_info "Building Docker image..."
    docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
    docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
    log_success "Docker image built successfully"
}

# Push Docker image
push_docker_image() {
    log_info "Pushing Docker image..."
    docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
    docker push ${DOCKER_IMAGE}:latest
    log_success "Docker image pushed successfully"
}

# Deploy to local Docker
deploy_local_docker() {
    log_info "Deploying to local Docker..."
    
    # Stop existing container if running
    if docker ps -q -f name=${BINARY_NAME} | grep -q .; then
        log_info "Stopping existing container..."
        docker stop ${BINARY_NAME}
        docker rm ${BINARY_NAME}
    fi
    
    # Run new container
    docker run -d --name ${BINARY_NAME} -p 8080:8080 ${DOCKER_IMAGE}:latest
    log_success "Deployed to local Docker successfully"
}

# Deploy to Kubernetes
deploy_kubernetes() {
    local namespace=$1
    log_info "Deploying to Kubernetes namespace: ${namespace}"
    
    # Create namespace if it doesn't exist
    kubectl create namespace ${namespace} --dry-run=client -o yaml | kubectl apply -f -
    
    # Apply Kubernetes manifests
    if [ -d "deployments/kubernetes" ]; then
        kubectl apply -f deployments/kubernetes/ -n ${namespace}
        log_success "Kubernetes deployment completed"
    else
        log_error "Kubernetes deployment files not found"
        exit 1
    fi
}

# Deploy using Helm
deploy_helm() {
    local namespace=$1
    local chart_path=$2
    log_info "Deploying using Helm to namespace: ${namespace}"
    
    if [ -z "$chart_path" ]; then
        chart_path="deployments/helm"
    fi
    
    if [ -d "$chart_path" ]; then
        helm upgrade --install ${BINARY_NAME} ${chart_path} --namespace ${namespace} --create-namespace
        log_success "Helm deployment completed"
    else
        log_error "Helm chart not found at: ${chart_path}"
        exit 1
    fi
}

# Deploy to AWS ECS
deploy_aws_ecs() {
    local cluster_name=$1
    local service_name=$2
    log_info "Deploying to AWS ECS cluster: ${cluster_name}, service: ${service_name}"
    
    # Update ECS service
    aws ecs update-service --cluster ${cluster_name} --service ${service_name} --force-new-deployment
    log_success "AWS ECS deployment initiated"
}

# Deploy to Google Cloud Run
deploy_gcp_cloud_run() {
    local service_name=$1
    local region=$2
    log_info "Deploying to Google Cloud Run service: ${service_name}, region: ${region}"
    
    # Deploy to Cloud Run
    gcloud run deploy ${service_name} --source . --region ${region} --platform managed
    log_success "Google Cloud Run deployment completed"
}

# Deploy to Azure Container Instances
deploy_azure_container_instances() {
    local resource_group=$1
    local container_name=$2
    log_info "Deploying to Azure Container Instances: ${container_name}"
    
    # Deploy to Azure Container Instances
    az container create --resource-group ${resource_group} --name ${container_name} --image ${DOCKER_IMAGE}:latest
    log_success "Azure Container Instances deployment completed"
}

# Health check
health_check() {
    local endpoint=$1
    local timeout=${2:-30}
    
    log_info "Performing health check on: ${endpoint}"
    
    for i in $(seq 1 ${timeout}); do
        if curl -f -s ${endpoint}/health > /dev/null 2>&1; then
            log_success "Health check passed"
            return 0
        fi
        log_info "Health check attempt ${i}/${timeout} failed, retrying in 5 seconds..."
        sleep 5
    done
    
    log_error "Health check failed after ${timeout} attempts"
    return 1
}

# Rollback deployment
rollback_deployment() {
    local deployment_type=$1
    local namespace=$2
    
    log_info "Rolling back ${deployment_type} deployment..."
    
    case $deployment_type in
        kubernetes)
            kubectl rollout undo deployment/${BINARY_NAME} -n ${namespace}
            ;;
        helm)
            helm rollback ${BINARY_NAME} -n ${namespace}
            ;;
        docker)
            docker stop ${BINARY_NAME}
            docker rm ${BINARY_NAME}
            docker run -d --name ${BINARY_NAME} -p 8080:8080 ${DOCKER_IMAGE}:previous
            ;;
        *)
            log_error "Unsupported deployment type for rollback: ${deployment_type}"
            exit 1
            ;;
    esac
    
    log_success "Rollback completed"
}

# Clean deployment artifacts
clean_deployment() {
    log_info "Cleaning deployment artifacts..."
    
    # Stop and remove Docker containers
    if docker ps -q -f name=${BINARY_NAME} | grep -q .; then
        docker stop ${BINARY_NAME}
        docker rm ${BINARY_NAME}
    fi
    
    # Remove Docker images
    docker rmi ${DOCKER_IMAGE}:${DOCKER_TAG} 2>/dev/null || true
    docker rmi ${DOCKER_IMAGE}:latest 2>/dev/null || true
    
    # Clean build artifacts
    rm -rf ${BUILD_DIR}
    rm -rf ${DIST_DIR}
    
    log_success "Deployment artifacts cleaned"
}

# Show deployment help
show_help() {
    echo "Go Micro Framework Deploy Script"
    echo "================================"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -t, --type TYPE              Deployment type (docker, kubernetes, helm, aws, gcp, azure)"
    echo "  -n, --namespace NAMESPACE    Kubernetes/Helm namespace (default: default)"
    echo "  -c, --cluster CLUSTER        AWS ECS cluster name"
    echo "  -s, --service SERVICE        Service name"
    echo "  -r, --region REGION          Cloud region"
    echo "  -g, --resource-group GROUP   Azure resource group"
    echo "  -p, --chart-path PATH        Helm chart path"
    echo "  -e, --endpoint ENDPOINT      Health check endpoint"
    echo "  --timeout TIMEOUT            Health check timeout (default: 30)"
    echo "  --rollback                   Rollback deployment"
    echo "  --clean                      Clean deployment artifacts"
    echo "  -h, --help                   Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 --type docker                    # Deploy to local Docker"
    echo "  $0 --type kubernetes --namespace prod # Deploy to Kubernetes"
    echo "  $0 --type helm --namespace prod     # Deploy using Helm"
    echo "  $0 --type aws --cluster my-cluster  # Deploy to AWS ECS"
    echo "  $0 --type gcp --service my-service  # Deploy to Google Cloud Run"
    echo "  $0 --type azure --resource-group my-rg # Deploy to Azure"
    echo "  $0 --rollback --type kubernetes     # Rollback Kubernetes deployment"
    echo "  $0 --clean                          # Clean deployment artifacts"
}

# Main function
main() {
    local deployment_type=""
    local namespace="default"
    local cluster_name=""
    local service_name=""
    local region=""
    local resource_group=""
    local chart_path=""
    local endpoint=""
    local timeout=30
    local rollback=false
    local clean_artifacts=false
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -t|--type)
                deployment_type="$2"
                shift 2
                ;;
            -n|--namespace)
                namespace="$2"
                shift 2
                ;;
            -c|--cluster)
                cluster_name="$2"
                shift 2
                ;;
            -s|--service)
                service_name="$2"
                shift 2
                ;;
            -r|--region)
                region="$2"
                shift 2
                ;;
            -g|--resource-group)
                resource_group="$2"
                shift 2
                ;;
            -p|--chart-path)
                chart_path="$2"
                shift 2
                ;;
            -e|--endpoint)
                endpoint="$2"
                shift 2
                ;;
            --timeout)
                timeout="$2"
                shift 2
                ;;
            --rollback)
                rollback=true
                shift
                ;;
            --clean)
                clean_artifacts=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # Execute actions
    if [ "$clean_artifacts" = true ]; then
        clean_deployment
        exit 0
    fi
    
    if [ "$rollback" = true ]; then
        if [ -z "$deployment_type" ]; then
            log_error "Deployment type is required for rollback"
            exit 1
        fi
        rollback_deployment "$deployment_type" "$namespace"
        exit 0
    fi
    
    if [ -z "$deployment_type" ]; then
        log_error "Deployment type is required"
        show_help
        exit 1
    fi
    
    # Build binary if needed
    if [ "$deployment_type" = "docker" ] || [ "$deployment_type" = "kubernetes" ] || [ "$deployment_type" = "helm" ]; then
        build_binary
    fi
    
    # Build Docker image if needed
    if [ "$deployment_type" = "docker" ] || [ "$deployment_type" = "kubernetes" ] || [ "$deployment_type" = "helm" ]; then
        build_docker_image
    fi
    
    # Deploy based on type
    case $deployment_type in
        docker)
            deploy_local_docker
            if [ -n "$endpoint" ]; then
                health_check "$endpoint" "$timeout"
            fi
            ;;
        kubernetes)
            deploy_kubernetes "$namespace"
            if [ -n "$endpoint" ]; then
                health_check "$endpoint" "$timeout"
            fi
            ;;
        helm)
            deploy_helm "$namespace" "$chart_path"
            if [ -n "$endpoint" ]; then
                health_check "$endpoint" "$timeout"
            fi
            ;;
        aws)
            if [ -z "$cluster_name" ] || [ -z "$service_name" ]; then
                log_error "Cluster name and service name are required for AWS deployment"
                exit 1
            fi
            deploy_aws_ecs "$cluster_name" "$service_name"
            ;;
        gcp)
            if [ -z "$service_name" ] || [ -z "$region" ]; then
                log_error "Service name and region are required for GCP deployment"
                exit 1
            fi
            deploy_gcp_cloud_run "$service_name" "$region"
            ;;
        azure)
            if [ -z "$resource_group" ] || [ -z "$service_name" ]; then
                log_error "Resource group and service name are required for Azure deployment"
                exit 1
            fi
            deploy_azure_container_instances "$resource_group" "$service_name"
            ;;
        *)
            log_error "Unsupported deployment type: ${deployment_type}"
            show_help
            exit 1
            ;;
    esac
    
    log_success "Deployment script completed"
}

# Run main function
main "$@"

# Go Micro Framework Makefile
# Provides common build, test, and development tasks

# Variables
BINARY_NAME=microframework
BINARY_PATH=cmd/microframework
BUILD_DIR=build
DIST_DIR=dist
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOLINT=golangci-lint

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: all build clean test deps fmt vet lint security install uninstall help

# Default target
all: clean deps fmt vet test build

# Help target
help: ## Show this help message
	@echo "$(BLUE)Go Micro Framework$(NC)"
	@echo "$(BLUE)===================$(NC)"
	@echo ""
	@echo "$(GREEN)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the binary
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(BINARY_PATH)
	@echo "$(GREEN)✓ Build completed$(NC)"

build-linux: ## Build for Linux
	@echo "$(BLUE)Building for Linux...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(BINARY_PATH)
	@echo "$(GREEN)✓ Linux build completed$(NC)"

build-darwin: ## Build for macOS
	@echo "$(BLUE)Building for macOS...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(BINARY_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(BINARY_PATH)
	@echo "$(GREEN)✓ macOS build completed$(NC)"

build-windows: ## Build for Windows
	@echo "$(BLUE)Building for Windows...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(BINARY_PATH)
	@echo "$(GREEN)✓ Windows build completed$(NC)"

build-all: ## Build for all platforms
	@echo "$(BLUE)Building for all platforms...$(NC)"
	$(MAKE) build-linux
	$(MAKE) build-darwin
	$(MAKE) build-windows
	@echo "$(GREEN)✓ All platform builds completed$(NC)"

# Development targets
dev: ## Run in development mode
	@echo "$(BLUE)Running in development mode...$(NC)"
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(BINARY_PATH)
	./$(BUILD_DIR)/$(BINARY_NAME)

run: ## Run the application
	@echo "$(BLUE)Running $(BINARY_NAME)...$(NC)"
	$(GOCMD) run $(BINARY_PATH)/main.go

# Test targets
test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"

test-unit: ## Run unit tests only
	@echo "$(BLUE)Running unit tests...$(NC)"
	$(GOTEST) -v -short ./...

test-integration: ## Run integration tests
	@echo "$(BLUE)Running integration tests...$(NC)"
	$(GOTEST) -v -run=Integration ./...

test-benchmark: ## Run benchmark tests
	@echo "$(BLUE)Running benchmark tests...$(NC)"
	$(GOTEST) -v -bench=. ./...

# Code quality targets
fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GOFMT) ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	$(GOVET) ./...
	@echo "$(GREEN)✓ Go vet completed$(NC)"

lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run ./...; \
	else \
		echo "$(YELLOW)Warning: golangci-lint not installed. Installing...$(NC)"; \
		$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		$(GOLINT) run ./...; \
	fi
	@echo "$(GREEN)✓ Linting completed$(NC)"

security: ## Run security checks
	@echo "$(BLUE)Running security checks...$(NC)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "$(YELLOW)Warning: gosec not installed. Installing...$(NC)"; \
		$(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi
	@echo "$(GREEN)✓ Security checks completed$(NC)"

# Dependency targets
deps: ## Download dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

deps-update: ## Update dependencies
	@echo "$(BLUE)Updating dependencies...$(NC)"
	$(GOMOD) tidy
	$(GOGET) -u ./...
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

deps-vendor: ## Vendor dependencies
	@echo "$(BLUE)Vendoring dependencies...$(NC)"
	$(GOMOD) vendor
	@echo "$(GREEN)✓ Dependencies vendored$(NC)"

# Installation targets
install: ## Install the binary
	@echo "$(BLUE)Installing $(BINARY_NAME)...$(NC)"
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(BINARY_PATH)
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "$(GREEN)✓ $(BINARY_NAME) installed to /usr/local/bin/$(NC)"

uninstall: ## Uninstall the binary
	@echo "$(BLUE)Uninstalling $(BINARY_NAME)...$(NC)"
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)✓ $(BINARY_NAME) uninstalled$(NC)"

# Docker targets
docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	docker build -t $(BINARY_NAME):$(VERSION) .
	docker build -t $(BINARY_NAME):latest .
	@echo "$(GREEN)✓ Docker image built$(NC)"

docker-run: ## Run Docker container
	@echo "$(BLUE)Running Docker container...$(NC)"
	docker run --rm -it $(BINARY_NAME):latest

docker-push: ## Push Docker image
	@echo "$(BLUE)Pushing Docker image...$(NC)"
	docker push $(BINARY_NAME):$(VERSION)
	docker push $(BINARY_NAME):latest
	@echo "$(GREEN)✓ Docker image pushed$(NC)"

# Cleanup targets
clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	rm -f coverage.out coverage.html
	@echo "$(GREEN)✓ Cleanup completed$(NC)"

clean-deps: ## Clean dependency cache
	@echo "$(BLUE)Cleaning dependency cache...$(NC)"
	$(GOCLEAN) -cache -modcache
	@echo "$(GREEN)✓ Dependency cache cleaned$(NC)"

# Release targets
release: ## Create release build
	@echo "$(BLUE)Creating release build...$(NC)"
	@mkdir -p $(DIST_DIR)
	$(MAKE) build-all
	@cp $(BUILD_DIR)/* $(DIST_DIR)/
	@echo "$(GREEN)✓ Release build created in $(DIST_DIR)/$(NC)"

release-tag: ## Create and push git tag
	@echo "$(BLUE)Creating git tag...$(NC)"
	@if [ -z "$(TAG)" ]; then \
		echo "$(RED)Error: TAG variable is required$(NC)"; \
		echo "Usage: make release-tag TAG=v1.0.0"; \
		exit 1; \
	fi
	git tag -a $(TAG) -m "Release $(TAG)"
	git push origin $(TAG)
	@echo "$(GREEN)✓ Tag $(TAG) created and pushed$(NC)"

# Documentation targets
docs: ## Generate documentation
	@echo "$(BLUE)Generating documentation...$(NC)"
	@if command -v godoc >/dev/null 2>&1; then \
		godoc -http=:6060; \
	else \
		echo "$(YELLOW)Warning: godoc not installed. Installing...$(NC)"; \
		$(GOGET) golang.org/x/tools/cmd/godoc@latest; \
		godoc -http=:6060; \
	fi

# CI/CD targets
ci: ## Run CI pipeline
	@echo "$(BLUE)Running CI pipeline...$(NC)"
	$(MAKE) deps
	$(MAKE) fmt
	$(MAKE) vet
	$(MAKE) lint
	$(MAKE) security
	$(MAKE) test
	$(MAKE) build
	@echo "$(GREEN)✓ CI pipeline completed$(NC)"

ci-test: ## Run CI tests
	@echo "$(BLUE)Running CI tests...$(NC)"
	$(MAKE) test-coverage
	@echo "$(GREEN)✓ CI tests completed$(NC)"

# Development tools
tools: ## Install development tools
	@echo "$(BLUE)Installing development tools...$(NC)"
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	$(GOGET) golang.org/x/tools/cmd/godoc@latest
	$(GOGET) github.com/golang/mock/mockgen@latest
	@echo "$(GREEN)✓ Development tools installed$(NC)"

# Version info
version: ## Show version information
	@echo "$(BLUE)Version Information$(NC)"
	@echo "$(GREEN)Version:$(NC) $(VERSION)"
	@echo "$(GREEN)Commit:$(NC) $(COMMIT)"
	@echo "$(GREEN)Date:$(NC) $(DATE)"
	@echo "$(GREEN)Go Version:$(NC) $$($(GOCMD) version)"

# Quick development cycle
quick: ## Quick development cycle (fmt, vet, test, build)
	@echo "$(BLUE)Running quick development cycle...$(NC)"
	$(MAKE) fmt
	$(MAKE) vet
	$(MAKE) test
	$(MAKE) build
	@echo "$(GREEN)✓ Quick development cycle completed$(NC)"

# Watch for changes and rebuild
watch: ## Watch for changes and rebuild
	@echo "$(BLUE)Watching for changes...$(NC)"
	@if command -v fswatch >/dev/null 2>&1; then \
		fswatch -o . | while read f; do \
			echo "$(YELLOW)Change detected, rebuilding...$(NC)"; \
			$(MAKE) quick; \
		done; \
	else \
		echo "$(RED)Error: fswatch not installed$(NC)"; \
		echo "Install with: brew install fswatch (macOS) or apt-get install fswatch (Ubuntu)"; \
	fi

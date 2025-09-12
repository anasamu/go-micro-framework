#!/bin/bash

# Go Micro Framework Build Script
# This script builds the framework for multiple platforms

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variables
BINARY_NAME="microframework"
BINARY_PATH="cmd/microframework"
BUILD_DIR="build"
DIST_DIR="dist"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS="-ldflags \"-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}\""

# Platform configurations
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

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

# Clean function
clean() {
    log_info "Cleaning build artifacts..."
    rm -rf ${BUILD_DIR}
    rm -rf ${DIST_DIR}
    go clean
    log_success "Cleanup completed"
}

# Build function
build() {
    local platform=$1
    local os=$(echo $platform | cut -d'/' -f1)
    local arch=$(echo $platform | cut -d'/' -f2)
    
    log_info "Building for ${os}/${arch}..."
    
    # Create build directory
    mkdir -p ${BUILD_DIR}
    
    # Set output filename
    local output_name="${BINARY_NAME}"
    if [ "$os" = "windows" ]; then
        output_name="${BINARY_NAME}.exe"
    fi
    
    # Build
    GOOS=${os} GOARCH=${arch} go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-${os}-${arch}${output_name:${#BINARY_NAME}} ${BINARY_PATH}
    
    log_success "Build completed for ${os}/${arch}"
}

# Build all platforms
build_all() {
    log_info "Building for all platforms..."
    
    for platform in "${PLATFORMS[@]}"; do
        build $platform
    done
    
    log_success "All platform builds completed"
}

# Build specific platform
build_platform() {
    local platform=$1
    
    if [[ " ${PLATFORMS[@]} " =~ " ${platform} " ]]; then
        build $platform
    else
        log_error "Unsupported platform: ${platform}"
        log_info "Supported platforms: ${PLATFORMS[*]}"
        exit 1
    fi
}

# Create distribution package
create_dist() {
    log_info "Creating distribution package..."
    
    mkdir -p ${DIST_DIR}
    
    # Copy binaries
    cp ${BUILD_DIR}/* ${DIST_DIR}/
    
    # Create archive for each platform
    for platform in "${PLATFORMS[@]}"; do
        local os=$(echo $platform | cut -d'/' -f1)
        local arch=$(echo $platform | cut -d'/' -f2)
        local binary_name="${BINARY_NAME}"
        
        if [ "$os" = "windows" ]; then
            binary_name="${BINARY_NAME}.exe"
        fi
        
        if [ -f "${DIST_DIR}/${BINARY_NAME}-${os}-${arch}${binary_name:${#BINARY_NAME}}" ]; then
            log_info "Creating archive for ${os}/${arch}..."
            cd ${DIST_DIR}
            tar -czf "${BINARY_NAME}-${VERSION}-${os}-${arch}.tar.gz" "${BINARY_NAME}-${os}-${arch}${binary_name:${#BINARY_NAME}}"
            cd ..
            log_success "Archive created: ${BINARY_NAME}-${VERSION}-${os}-${arch}.tar.gz"
        fi
    done
    
    log_success "Distribution package created"
}

# Show help
show_help() {
    echo "Go Micro Framework Build Script"
    echo "==============================="
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -p, --platform PLATFORM    Build for specific platform (e.g., linux/amd64)"
    echo "  -a, --all                  Build for all platforms"
    echo "  -c, --clean                Clean build artifacts"
    echo "  -d, --dist                 Create distribution package"
    echo "  -h, --help                 Show this help message"
    echo ""
    echo "Supported platforms:"
    for platform in "${PLATFORMS[@]}"; do
        echo "  - ${platform}"
    done
    echo ""
    echo "Examples:"
    echo "  $0 --all                   # Build for all platforms"
    echo "  $0 --platform linux/amd64  # Build for Linux AMD64"
    echo "  $0 --clean                 # Clean build artifacts"
    echo "  $0 --dist                  # Create distribution package"
}

# Main function
main() {
    local platform=""
    local build_all=false
    local clean_build=false
    local create_dist_package=false
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -p|--platform)
                platform="$2"
                shift 2
                ;;
            -a|--all)
                build_all=true
                shift
                ;;
            -c|--clean)
                clean_build=true
                shift
                ;;
            -d|--dist)
                create_dist_package=true
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
    if [ "$clean_build" = true ]; then
        clean
    fi
    
    if [ "$build_all" = true ]; then
        build_all
    elif [ -n "$platform" ]; then
        build_platform "$platform"
    else
        # Default: build for current platform
        log_info "Building for current platform..."
        go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ${BINARY_PATH}
        log_success "Build completed for current platform"
    fi
    
    if [ "$create_dist_package" = true ]; then
        create_dist
    fi
    
    log_success "Build script completed"
}

# Run main function
main "$@"

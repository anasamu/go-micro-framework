#!/bin/bash

# Go Micro Framework Test Script
# This script runs various types of tests

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variables
COVERAGE_FILE="coverage.out"
COVERAGE_HTML="coverage.html"
TEST_TIMEOUT="10m"

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

# Run unit tests
run_unit_tests() {
    log_info "Running unit tests..."
    go test -v -short -timeout=${TEST_TIMEOUT} ./...
    log_success "Unit tests completed"
}

# Run integration tests
run_integration_tests() {
    log_info "Running integration tests..."
    go test -v -run=Integration -timeout=${TEST_TIMEOUT} ./...
    log_success "Integration tests completed"
}

# Run all tests
run_all_tests() {
    log_info "Running all tests..."
    go test -v -timeout=${TEST_TIMEOUT} ./...
    log_success "All tests completed"
}

# Run tests with coverage
run_tests_with_coverage() {
    log_info "Running tests with coverage..."
    go test -v -coverprofile=${COVERAGE_FILE} -timeout=${TEST_TIMEOUT} ./...
    
    if [ -f "${COVERAGE_FILE}" ]; then
        log_info "Generating coverage report..."
        go tool cover -html=${COVERAGE_FILE} -o ${COVERAGE_HTML}
        log_success "Coverage report generated: ${COVERAGE_HTML}"
        
        # Show coverage percentage
        coverage_percent=$(go tool cover -func=${COVERAGE_FILE} | grep total | awk '{print $3}')
        log_info "Total coverage: ${coverage_percent}"
    else
        log_warning "Coverage file not generated"
    fi
}

# Run benchmark tests
run_benchmark_tests() {
    log_info "Running benchmark tests..."
    go test -v -bench=. -benchmem -timeout=${TEST_TIMEOUT} ./...
    log_success "Benchmark tests completed"
}

# Run race detection tests
run_race_tests() {
    log_info "Running race detection tests..."
    go test -v -race -timeout=${TEST_TIMEOUT} ./...
    log_success "Race detection tests completed"
}

# Run tests with verbose output
run_verbose_tests() {
    log_info "Running tests with verbose output..."
    go test -v -timeout=${TEST_TIMEOUT} ./...
    log_success "Verbose tests completed"
}

# Run tests for specific package
run_package_tests() {
    local package=$1
    log_info "Running tests for package: ${package}"
    go test -v -timeout=${TEST_TIMEOUT} ./${package}
    log_success "Package tests completed for: ${package}"
}

# Run tests with specific pattern
run_pattern_tests() {
    local pattern=$1
    log_info "Running tests with pattern: ${pattern}"
    go test -v -run=${pattern} -timeout=${TEST_TIMEOUT} ./...
    log_success "Pattern tests completed for: ${pattern}"
}

# Clean test artifacts
clean_test_artifacts() {
    log_info "Cleaning test artifacts..."
    rm -f ${COVERAGE_FILE}
    rm -f ${COVERAGE_HTML}
    go clean -testcache
    log_success "Test artifacts cleaned"
}

# Show test help
show_help() {
    echo "Go Micro Framework Test Script"
    echo "=============================="
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -u, --unit                    Run unit tests only"
    echo "  -i, --integration             Run integration tests only"
    echo "  -a, --all                     Run all tests"
    echo "  -c, --coverage                Run tests with coverage"
    echo "  -b, --benchmark               Run benchmark tests"
    echo "  -r, --race                    Run race detection tests"
    echo "  -v, --verbose                 Run tests with verbose output"
    echo "  -p, --package PACKAGE         Run tests for specific package"
    echo "  -t, --pattern PATTERN         Run tests with specific pattern"
    echo "  --clean                       Clean test artifacts"
    echo "  -h, --help                    Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 --unit                     # Run unit tests"
    echo "  $0 --integration              # Run integration tests"
    echo "  $0 --coverage                 # Run tests with coverage"
    echo "  $0 --package ./internal/core  # Run tests for specific package"
    echo "  $0 --pattern TestUser         # Run tests matching pattern"
    echo "  $0 --clean                    # Clean test artifacts"
}

# Main function
main() {
    local run_unit=false
    local run_integration=false
    local run_all=false
    local run_coverage=false
    local run_benchmark=false
    local run_race=false
    local run_verbose=false
    local package=""
    local pattern=""
    local clean_artifacts=false
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -u|--unit)
                run_unit=true
                shift
                ;;
            -i|--integration)
                run_integration=true
                shift
                ;;
            -a|--all)
                run_all=true
                shift
                ;;
            -c|--coverage)
                run_coverage=true
                shift
                ;;
            -b|--benchmark)
                run_benchmark=true
                shift
                ;;
            -r|--race)
                run_race=true
                shift
                ;;
            -v|--verbose)
                run_verbose=true
                shift
                ;;
            -p|--package)
                package="$2"
                shift 2
                ;;
            -t|--pattern)
                pattern="$2"
                shift 2
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
        clean_test_artifacts
    fi
    
    if [ "$run_unit" = true ]; then
        run_unit_tests
    elif [ "$run_integration" = true ]; then
        run_integration_tests
    elif [ "$run_all" = true ]; then
        run_all_tests
    elif [ "$run_coverage" = true ]; then
        run_tests_with_coverage
    elif [ "$run_benchmark" = true ]; then
        run_benchmark_tests
    elif [ "$run_race" = true ]; then
        run_race_tests
    elif [ "$run_verbose" = true ]; then
        run_verbose_tests
    elif [ -n "$package" ]; then
        run_package_tests "$package"
    elif [ -n "$pattern" ]; then
        run_pattern_tests "$pattern"
    else
        # Default: run all tests
        run_all_tests
    fi
    
    log_success "Test script completed"
}

# Run main function
main "$@"

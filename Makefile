# Makefile for RSS Reader

# Variables
BINARY_SERVICE := rss-service
BINARY_UI := rss-ui
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(BUILD_TIME)

# Build directories
BUILD_DIR := build
DIST_DIR := dist

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

.PHONY: all build build-service build-ui clean test deps fmt install run-service run-ui release help

# Default target
all: clean build

# Build both binaries
build: build-service build-ui

# Build the service binary
build-service:
	@echo "Building service..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_SERVICE) ./cmd/service

# Build the UI binary
build-ui:
	@echo "Building UI..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_UI) ./cmd/ui

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	$(GOCLEAN) ./...

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Update dependencies
deps:
	@echo "Updating dependencies..."
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	gofmt -s -w .

# Install the application
install: build
	@echo "Installing..."
	@mkdir -p $(HOME)/.rss-reader/bin
	@cp $(BUILD_DIR)/$(BINARY_SERVICE) $(HOME)/.rss-reader/bin/
	@cp $(BUILD_DIR)/$(BINARY_UI) $(HOME)/.rss-reader/bin/
	@chmod +x scripts/install.sh
	@scripts/install.sh

# Run the service
run-service: build-service
	@echo "Running service..."
	@$(BUILD_DIR)/$(BINARY_SERVICE)

# Run the UI
run-ui: build-ui
	@echo "Running UI..."
	@$(BUILD_DIR)/$(BINARY_UI)

# Create a release build
release:
	@echo "Creating release build..."
	goreleaser --snapshot --rm-dist

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	@go install github.com/goreleaser/goreleaser@latest

# Help
help:
	@echo "Make targets:"
	@echo "  build         - Build both service and UI binaries"
	@echo "  build-service - Build only the service binary"
	@echo "  build-ui      - Build only the UI binary"
	@echo "  clean         - Remove build artifacts"
	@echo "  test          - Run tests"
	@echo "  deps          - Update dependencies"
	@echo "  fmt           - Format code"
	@echo "  install       - Install the application"
	@echo "  run-service   - Run the service"
	@echo "  run-ui        - Run the UI"
	@echo "  release       - Create a release build"
	@echo "  dev-setup     - Setup development environment"
	@echo "  help          - Show this help"

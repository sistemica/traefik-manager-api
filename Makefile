# Traefik Manager Project Makefile

# Project variables
PROJECT_NAME := traefik-manager
VERSION := 1.0.0
GO_VERSION := 1.23.5
MAIN_PATH := cmd/server/main.go
BUILD_DIR := bin
DOCKER_IMAGE := ${PROJECT_NAME}
DOCKER_TAG := latest

# Ensure all targets are phony
.PHONY: all build run test clean help docker-build docker-run docker-compose-up docker-compose-down lint fmt

# Build, test, and run targets
all: clean build test

# Build the application
build:
	@echo "Building ${PROJECT_NAME}..."
	@mkdir -p ${BUILD_DIR}
	go build -o ${BUILD_DIR}/${PROJECT_NAME} ${MAIN_PATH}
	@echo "✓ Build complete: ${BUILD_DIR}/${PROJECT_NAME}"

# Run the application
run:
	@echo "Running ${PROJECT_NAME}..."
	@mkdir -p data
	STORAGE_FILE_PATH="./data/traefik-manager.json" go run ${MAIN_PATH}

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	@rm -rf ${BUILD_DIR}
	@go clean
	@echo "✓ Clean complete"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Format complete"

# Lint code
lint:
	@echo "Linting code..."
	@if command -v golint >/dev/null 2>&1; then \
		golint ./...; \
	else \
		echo "Warning: golint not installed. Run: go install golang.org/x/lint/golint@latest"; \
	fi
	@echo "✓ Lint complete"

# Docker commands
docker-build:
	@echo "Building Docker image: ${DOCKER_IMAGE}:${DOCKER_TAG}..."
	docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
	@echo "✓ Docker build complete"

docker-run:
	@echo "Running Docker container: ${DOCKER_IMAGE}:${DOCKER_TAG}..."
	docker run -p 9000:9000 -d --name ${PROJECT_NAME} ${DOCKER_IMAGE}:${DOCKER_TAG}
	@echo "✓ Docker container started"

docker-compose-up:
	@echo "Starting Docker Compose services..."
	docker-compose up -d
	@echo "✓ Services started. Traefik dashboard: http://localhost:8080, Traefik Manager API: http://localhost:9000/api/v1"

docker-compose-down:
	@echo "Stopping Docker Compose services..."
	docker-compose down
	@echo "✓ Services stopped"

docker-logs:
	@echo "Showing Docker Compose logs..."
	docker-compose logs -f

# Quick start guide
help:
	@echo "Traefik Manager - Make Targets:"
	@echo ""
	@echo "Build Commands:"
	@echo "  make build                - Build the application"
	@echo "  make clean                - Remove generated files"
	@echo "  make all                  - Clean, build, and test"
	@echo ""
	@echo "Development Commands:"
	@echo "  make run                  - Run the application locally"
	@echo "  make test                 - Run tests"
	@echo "  make fmt                  - Format code"
	@echo "  make lint                 - Lint code"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make docker-build         - Build Docker image"
	@echo "  make docker-run           - Run Docker container"
	@echo "  make docker-compose-up    - Start all services with Docker Compose"
	@echo "  make docker-compose-down  - Stop all services"
	@echo "  make docker-logs          - Show logs from Docker Compose services"
	@echo ""
	@echo "Prerequisites:"
	@echo "  - Go ${GO_VERSION} or later"
	@echo "  - Docker and Docker Compose for container-related commands"
.DEFAULT_GOAL := help

help:
	@echo "KV-STORE - MAKEFILE"
	@echo ""
	@echo "Available commands:"
	@echo " run              Run the application (Redis storage - default)"
	@echo " run-memory       Run the application (in-memory storage)"
	@echo " build            Build the application binary"
	@echo " test             Run unit tests"
	@echo " test-integration Run integration tests (requires Docker)"
	@echo " test-all         Run all tests (unit + integration)"
	@echo " lint             Run linter"
	@echo " docker-up        Start Redis using Docker Compose"
	@echo " docker-down      Stop Redis"
	@echo ""

run:
	@echo "Starting application with Redis storage..."
	@docker-compose up -d
	@sleep 2
	STORAGE_TYPE=redis go run cmd/api/main.go

run-memory:
	@echo "Starting application with in-memory storage..."
	STORAGE_TYPE=memory go run cmd/api/main.go

build:
	@echo "Building application..."
	go build -o bin/server cmd/api/main.go
	@echo "Binary created at bin/server"

test:
	@echo "Running unit tests..."
	go test ./...

test-integration:
	@echo "Running integration tests (requires Docker)..."
	go test -tags=integration ./...

test-all:
	@echo "Running all tests (unit + integration)..."
	go test ./...
	go test -tags=integration ./...

lint:
	@echo "Running linter..."
	golangci-lint run ./...

docker-up:
	@echo "Starting Redis..."
	docker-compose up -d
	@echo "Redis is running on localhost:6379"

docker-down:
	@echo "Stopping Redis..."
	docker-compose down

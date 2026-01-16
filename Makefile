.DEFAULT_GOAL := help

help:
	@echo "KV-Store"
	@echo ""
	@echo "Available commands:"
	@echo "  run                Run application with Redis"
	@echo "  run-memory         Run application with in-memory storage"
	@echo "  test               Run all tests (unit + integration)"
	@echo "  test-unit          Run unit tests only"
	@echo "  test-integration   Run integration tests only"
	@echo "  lint               Run linter"
	@echo ""

run:
	@echo "Starting application..."
	@echo "Stopping any process on port 8080..."
	@lsof -ti:8080 | xargs kill -9 2>/dev/null || true
	@docker-compose down 2>/dev/null || true
	@docker-compose up -d
	@sleep 2
	STORAGE_TYPE=redis go run cmd/api/main.go

run-memory:
	@echo "Starting application (in-memory)..."
	STORAGE_TYPE=memory go run cmd/api/main.go

test:
	@echo "Running all tests..."
	@go test ./... -short
	@go test -tags=integration ./... -v

test-unit:
	@echo "Running unit tests..."
	@go test ./... -short

test-integration:
	@echo "Running integration tests..."
	@go test -tags=integration ./... -v

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

.PHONY: help run run-memory test test-unit test-integration lint

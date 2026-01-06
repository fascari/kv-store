# KV Store

Key-value store with HTTP API using Go and Chi router. Supports both in-memory and Redis storage backends.

## Features

- ✅ Clean Architecture (handlers, use cases, stores)
- ✅ Pluggable storage backends (memory or Redis)
- ✅ RESTful HTTP API
- ✅ Docker Compose for local Redis
- ✅ Environment-based configuration
- ✅ Structured logging with Zap


## Configuration

Copy the example environment file and adjust as needed:

```bash
cp .env.example .env
```

### Environment Variables

```bash
# Storage type: "memory" or "redis" (default: redis)
STORAGE_TYPE=redis

# Redis configuration (when STORAGE_TYPE=redis)
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Server configuration
SERVER_PORT=8080
```

## Running

### With Redis Storage (Default)

```bash
# Start Redis and run the application
make run

# Or manually
docker-compose up -d
go run cmd/api/main.go
```

### With In-Memory Storage

```bash
# Run with in-memory storage
make run-memory

# Or manually
STORAGE_TYPE=memory go run cmd/api/main.go
```

### Stopping Services

```bash
# Stop Redis
make docker-down
```

## Testing

### Unit Tests

Run fast unit tests (in-memory storage only):

```bash
make test
# or
go test ./...
```

### Integration Tests

Run integration tests with Redis using testcontainers (requires Docker):

```bash
make test-integration
# or
go test -tags=integration ./...
```

### All Tests

Run both unit and integration tests:

```bash
make test-all
```

**Note:** Integration tests automatically spin up Redis containers using testcontainers, so you don't need to manually start Redis.

### Other Commands

```bash
# Run linter
make lint

# Format code
make fmt
```

Server starts on `http://localhost:8080` (or the port specified in `SERVER_PORT`)

## Endpoints

### Health Check

**GET /health**

```bash
curl http://localhost:8080/health
```

Response:
```json
{"status": "ok"}
```

### Store a Value

**POST /api/keys**

```bash
curl -X POST http://localhost:8080/api/keys \
  -H "Content-Type: application/json" \
  -d '{"key": "name", "value": "John Doe"}'
```

Response:
```json
{"key": "name", "value": "John Doe"}
```

Store complex values:
```bash
curl -X POST http://localhost:8080/api/keys \
  -H "Content-Type: application/json" \
  -d '{"key": "config", "value": {"theme": "dark", "language": "en"}}'
```

### Retrieve a Value

**GET /api/keys/{key}**

```bash
curl http://localhost:8080/api/keys/name
```

Response:
```json
{"key": "name", "value": "John Doe"}
```

### Delete a Key

**DELETE /api/keys/{key}**

```bash
curl -X DELETE http://localhost:8080/api/keys/name
```

Response: `204 No Content`

## Design Principles

- Separation of concerns - each package has a single responsibility
- Early returns - no nested if/else statements
- Type grouping - related types grouped with `type ()`
- Idiomatic Go - following [Google Go Style Guide](https://google.github.io/styleguide/go/)

## Dependencies

- [Chi Router](https://github.com/go-chi/chi) - HTTP router
- [Zap Logger](https://github.com/uber-go/zap) - Structured logger
- [golangci-lint](https://golangci-lint.run/) - Go linter (dev dependency)

### Why Chi?

Chi was chosen because:
- **Lightweight and fast** - minimal overhead, pure stdlib compatible
- **Idiomatic Go** - follows standard `http.Handler` interface
- **Flexible routing** - supports route groups, middleware chaining, and URL parameters
- **No reflection** - compile-time route definition
- **Battle-tested** - widely used in production systems

### Why Zap Logger?

Zap was chosen because:
- **Performance** - fastest structured logger in Go ecosystem
- **Structured logging** - JSON output for easy parsing and analysis
- **Type-safe** - compile-time type checking for log fields
- **Zero allocations** - minimal memory overhead
- **Development mode** - human-readable output for local development


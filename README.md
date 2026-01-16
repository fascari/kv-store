# KV Store

A simple key-value store with distributed locking for multi-server environments.

## Features

- ✅ RESTful HTTP API
- ✅ Distributed Locking (Redis)
- ✅ In-Memory or Redis Storage
- ✅ Clean Architecture
- ✅ Comprehensive Tests

## Quick Start

```bash
# Start with Redis (distributed)
make run

# Or in-memory (single server)
make run-memory
```

Server starts on `http://localhost:8080`

## API Usage

### Save a key
```bash
curl -X POST http://localhost:8080/api/keys \
  -H "Content-Type: application/json" \
  -d '{"key": "user:1", "value": {"name": "Alice"}}'
```

### Retrieve a key
```bash
curl http://localhost:8080/api/keys/user:1
```

### Delete a key
```bash
curl -X DELETE http://localhost:8080/api/keys/user:1
```

### Health check
```bash
curl http://localhost:8080/health
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `STORAGE_TYPE` | `redis` | Storage backend: `memory` or `redis` |
| `SERVER_PORT` | `8080` | HTTP server port |
| `REDIS_ADDR` | `localhost:6379` | Redis address |
| `REDIS_PASSWORD` | - | Redis password |
| `REDIS_DB` | `0` | Redis database |

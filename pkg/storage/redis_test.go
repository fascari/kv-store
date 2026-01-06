//go:build integration

package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/felipeascari/kv-store/pkg/storage"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestRedis_Save(t *testing.T) {
	ctx := context.Background()
	container, store := setupRedisContainer(t, ctx)
	defer terminateContainer(t, container)
	defer closeRedis(t, store)

	runSaveTests(t, store, saveTests)
}

func TestRedis_Retrieve(t *testing.T) {
	ctx := context.Background()
	container, store := setupRedisContainer(t, ctx)
	defer terminateContainer(t, container)
	defer closeRedis(t, store)

	runRetrieveTests(t, store, retrieveTests())
}

func TestRedis_Delete(t *testing.T) {
	ctx := context.Background()
	container, store := setupRedisContainer(t, ctx)
	defer terminateContainer(t, container)
	defer closeRedis(t, store)

	runDeleteTests(t, store, deleteTests)
}

func TestRedis_Ping(t *testing.T) {
	ctx := context.Background()
	container, store := setupRedisContainer(t, ctx)
	defer terminateContainer(t, container)
	defer closeRedis(t, store)

	err := store.Ping(5 * time.Second)
	require.NoError(t, err)
}

func setupRedisContainer(t *testing.T, ctx context.Context) (*redis.RedisContainer, *storage.Redis) {
	t.Helper()

	redisContainer, err := redis.Run(ctx, "redis:7-alpine")
	require.NoError(t, err, "failed to start Redis container")

	host, err := redisContainer.Host(ctx)
	require.NoError(t, err, "failed to get Redis host")

	port, err := redisContainer.MappedPort(ctx, "6379/tcp")
	require.NoError(t, err, "failed to get Redis port")

	addr := host + ":" + port.Port()

	store, err := storage.NewRedis(addr, "", 0)
	require.NoError(t, err, "failed to create Redis store")

	return redisContainer, store
}

func terminateContainer(t *testing.T, container *redis.RedisContainer) {
	t.Helper()
	if container != nil {
		_ = testcontainers.TerminateContainer(container)
	}
}

func closeRedis(t *testing.T, store *storage.Redis) {
	t.Helper()
	_ = store.Close()
}

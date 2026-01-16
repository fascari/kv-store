//go:build integration

package lock_test

import (
	"context"
	"testing"
	"time"

	"github.com/felipeascari/kv-store/pkg/lock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestRedisLock(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	client, cleanup := setupRedis(t, ctx)
	defer cleanup()

	t.Run("should acquire and release lock", func(t *testing.T) {
		redisLock := lock.NewRedisLock(client, 5*time.Second)

		token, err := redisLock.Acquire(ctx, "test-key")
		require.NoError(t, err)
		require.Greater(t, token, int64(0))

		lockKey := "lock:test-key"
		exists := client.Exists(ctx, lockKey).Val()
		require.Equal(t, int64(1), exists)

		err = redisLock.Release(ctx, "test-key", token)
		require.NoError(t, err)

		exists = client.Exists(ctx, lockKey).Val()
		require.Equal(t, int64(0), exists)
	})

	t.Run("should prevent double acquisition", func(t *testing.T) {
		client.FlushDB(ctx)
		redisLock := lock.NewRedisLock(client, 5*time.Second)

		token1, err := redisLock.Acquire(ctx, "double-key")
		require.NoError(t, err)

		ctx2, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		_, err = redisLock.Acquire(ctx2, "double-key")
		require.Error(t, err)
		require.Contains(t, err.Error(), "lock already held")

		err = redisLock.Release(ctx, "double-key", token1)
		require.NoError(t, err)

		token3, err := redisLock.Acquire(ctx, "double-key")
		require.NoError(t, err)
		require.Greater(t, token3, token1)

		_ = redisLock.Release(ctx, "double-key", token3)
	})

	t.Run("should validate tokens correctly", func(t *testing.T) {
		client.FlushDB(ctx)
		redisLock := lock.NewRedisLock(client, 5*time.Second)

		token, err := redisLock.Acquire(ctx, "validate-key")
		require.NoError(t, err)

		valid, err := redisLock.ValidateToken(ctx, "validate-key", token)
		require.NoError(t, err)
		require.True(t, valid)

		valid, err = redisLock.ValidateToken(ctx, "validate-key", token+1)
		require.NoError(t, err)
		require.False(t, valid)

		valid, err = redisLock.ValidateToken(ctx, "nonexistent-key", token)
		require.NoError(t, err)
		require.False(t, valid)

		_ = redisLock.Release(ctx, "validate-key", token)
	})

	t.Run("should expire lock to prevent zombie processes", func(t *testing.T) {
		client.FlushDB(ctx)
		redisLock := lock.NewRedisLock(client, 500*time.Millisecond)

		token, err := redisLock.Acquire(ctx, "expiring-key")
		require.NoError(t, err)

		valid, err := redisLock.ValidateToken(ctx, "expiring-key", token)
		require.NoError(t, err)
		require.True(t, valid)

		time.Sleep(1 * time.Second)

		valid, err = redisLock.ValidateToken(ctx, "expiring-key", token)
		require.NoError(t, err)
		require.False(t, valid)

		token2, err := redisLock.Acquire(ctx, "expiring-key")
		require.NoError(t, err)
		require.Greater(t, token2, token)

		_ = redisLock.Release(ctx, "expiring-key", token2)
	})

	t.Run("should reject release with wrong token", func(t *testing.T) {
		client.FlushDB(ctx)
		redisLock := lock.NewRedisLock(client, 5*time.Second)

		token, err := redisLock.Acquire(ctx, "wrong-token-key")
		require.NoError(t, err)

		err = redisLock.Release(ctx, "wrong-token-key", token+1)
		require.Error(t, err)

		valid, err := redisLock.ValidateToken(ctx, "wrong-token-key", token)
		require.NoError(t, err)
		require.True(t, valid)

		err = redisLock.Release(ctx, "wrong-token-key", token)
		require.NoError(t, err)
	})
}

func setupRedis(t *testing.T, ctx context.Context) (*redis.Client, func()) {
	t.Helper()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:7-alpine",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections"),
		},
		Started: true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "6379")
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: host + ":" + port.Port(),
	})

	err = client.Ping(ctx).Err()
	require.NoError(t, err)

	cleanup := func() {
		_ = client.Close()
		_ = container.Terminate(ctx)
	}

	return client, cleanup
}

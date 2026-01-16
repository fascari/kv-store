//go:build integration

package storage_test

import (
	"context"
	"testing"

	"github.com/felipeascari/kv-store/pkg/storage"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestRedis(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	store, cleanup := setupRedis(t, ctx)
	defer cleanup()

	t.Run("should save and retrieve values", func(t *testing.T) {
		tests := []struct {
			name  string
			key   string
			value any
			want  any
		}{
			{
				name:  "string value",
				key:   "name",
				value: "Alice",
				want:  "Alice",
			},
			{
				name:  "integer value",
				key:   "age",
				value: 9,
				want:  float64(9),
			},
			{
				name:  "map value",
				key:   "user",
				value: map[string]any{"id": 1},
				want:  map[string]any{"id": float64(1)},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				require.NoError(t, store.Save(tt.key, tt.value))

				value, err := store.Retrieve(tt.key)
				require.NoError(t, err)
				require.Equal(t, tt.want, value)
			})
		}
	})

	t.Run("should return error when retrieving non-existent key", func(t *testing.T) {
		_, err := store.Retrieve("nonexistent")
		require.ErrorIs(t, err, storage.ErrKeyNotFound)
	})

	t.Run("should delete existing key", func(t *testing.T) {
		require.NoError(t, store.Save("temp", "value"))
		require.NoError(t, store.Delete("temp"))

		_, err := store.Retrieve("temp")
		require.ErrorIs(t, err, storage.ErrKeyNotFound)
	})

	t.Run("should return error when deleting non-existent key", func(t *testing.T) {
		err := store.Delete("nonexistent")
		require.ErrorIs(t, err, storage.ErrKeyNotFound)
	})
}

func setupRedis(t *testing.T, ctx context.Context) (*storage.Redis, func()) {
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

	store, err := storage.NewRedis(host+":"+port.Port(), "", 0)
	require.NoError(t, err)

	cleanup := func() {
		_ = store.Close()
		_ = container.Terminate(ctx)
	}

	return store, cleanup
}

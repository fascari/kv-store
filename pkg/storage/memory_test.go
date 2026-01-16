package storage_test

import (
	"testing"

	"github.com/felipeascari/kv-store/pkg/storage"
	"github.com/stretchr/testify/require"
)

func TestMemorySave(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    any
		expected any
	}{
		{
			name:     "should save string value",
			key:      "name",
			value:    "Alice",
			expected: "Alice",
		},
		{
			name:     "should save integer value",
			key:      "age",
			value:    30,
			expected: 30,
		},
		{
			name:     "should save map value",
			key:      "user",
			value:    map[string]any{"id": 1, "name": "Bob"},
			expected: map[string]any{"id": 1, "name": "Bob"},
		},
		{
			name:     "should overwrite existing value",
			key:      "counter",
			value:    42,
			expected: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.NewMemory()

			err := store.Save(tt.key, tt.value)
			require.NoError(t, err)

			value, err := store.Retrieve(tt.key)
			require.NoError(t, err)
			require.Equal(t, tt.expected, value)
		})
	}
}

func TestMemoryRetrieve(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*storage.Memory)
		key         string
		expectError bool
		expectValue any
	}{
		{
			name: "should retrieve existing key",
			setup: func(m *storage.Memory) {
				_ = m.Save("key1", "value1")
			},
			key:         "key1",
			expectError: false,
			expectValue: "value1",
		},
		{
			name: "should return error when retrieving non-existent key",
			setup: func(_ *storage.Memory) {
				// no setup
			},
			key:         "nonexistent",
			expectError: true,
		},
		{
			name: "should retrieve after multiple saves",
			setup: func(m *storage.Memory) {
				_ = m.Save("a", "1")
				_ = m.Save("b", "2")
				_ = m.Save("c", "3")
			},
			key:         "b",
			expectError: false,
			expectValue: "2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.NewMemory()
			tt.setup(store)

			value, err := store.Retrieve(tt.key)

			if tt.expectError {
				require.Error(t, err)
				require.ErrorIs(t, err, storage.ErrKeyNotFound)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectValue, value)
		})
	}
}

func TestMemoryDelete(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*storage.Memory)
		key         string
		expectError bool
		verifyGone  bool
	}{
		{
			name: "should delete existing key",
			setup: func(m *storage.Memory) {
				_ = m.Save("key1", "value1")
			},
			key:         "key1",
			expectError: false,
			verifyGone:  true,
		},
		{
			name: "should return error when deleting non-existent key",
			setup: func(_ *storage.Memory) {
				// no setup
			},
			key:         "nonexistent",
			expectError: true,
			verifyGone:  false,
		},
		{
			name: "should delete one of multiple keys",
			setup: func(m *storage.Memory) {
				_ = m.Save("a", "1")
				_ = m.Save("b", "2")
				_ = m.Save("c", "3")
			},
			key:         "b",
			expectError: false,
			verifyGone:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.NewMemory()
			tt.setup(store)

			err := store.Delete(tt.key)

			if tt.expectError {
				require.Error(t, err)
				require.ErrorIs(t, err, storage.ErrKeyNotFound)
				return
			}

			require.NoError(t, err)

			if tt.verifyGone {
				_, err := store.Retrieve(tt.key)
				require.Error(t, err)
				require.ErrorIs(t, err, storage.ErrKeyNotFound)
			}
		})
	}
}

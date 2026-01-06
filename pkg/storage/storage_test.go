package storage_test

import (
	"testing"

	"github.com/felipeascari/kv-store/pkg/storage"
	"github.com/stretchr/testify/require"
)

var (
	saveTests = []saveTest{
		{
			name:  "string value",
			key:   "key1",
			value: "value1",
		},
		{
			name:  "numeric value",
			key:   "key2",
			value: 42,
		},
		{
			name:  "boolean value",
			key:   "key3",
			value: true,
		},
		{
			name: "complex value",
			key:  "key4",
			value: map[string]any{
				"name": "John",
				"age":  30.0,
			},
		},
	}

	deleteTests = []deleteTest{
		{
			name:      "existing key",
			setupKey:  "to_delete",
			setupVal:  "value",
			deleteKey: "to_delete",
			wantError: nil,
		},
		{
			name:      "non-existing key",
			setupKey:  "",
			deleteKey: "nonexistent",
			wantError: storage.ErrKeyNotFound,
		},
	}
)

type (
	saveTest struct {
		name  string
		key   string
		value any
	}

	retrieveTest struct {
		name      string
		setupKey  string
		setupVal  any
		key       string
		wantValue any
		wantError error
		validate  func(t *testing.T, value any)
	}

	deleteTest struct {
		name      string
		setupKey  string
		setupVal  any
		deleteKey string
		wantError error
	}
)

func retrieveTests() []retrieveTest {
	return []retrieveTest{
		{
			name:      "existing key",
			setupKey:  "existing",
			setupVal:  "test_value",
			key:       "existing",
			wantValue: "test_value",
			wantError: nil,
		},
		{
			name:      "non-existing key",
			setupKey:  "",
			key:       "nonexistent",
			wantValue: nil,
			wantError: storage.ErrKeyNotFound,
		},
		{
			name:     "complex value",
			setupKey: "complex",
			setupVal: map[string]any{
				"name": "John",
				"age":  30.0,
			},
			key: "complex",
			validate: func(t *testing.T, value any) {
				m, ok := value.(map[string]any)
				require.True(t, ok, "expected map[string]any, got %T", value)
				require.Equal(t, "John", m["name"])
				require.Equal(t, 30.0, m["age"])
			},
			wantError: nil,
		},
	}
}

func runSaveTests(t *testing.T, store storage.Store, tests []saveTest) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Save(tt.key, tt.value)
			require.NoError(t, err)
		})
	}
}

func runRetrieveTests(t *testing.T, store storage.Store, tests []retrieveTest) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupKey != "" {
				err := store.Save(tt.setupKey, tt.setupVal)
				require.NoError(t, err)
			}

			value, err := store.Retrieve(tt.key)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				return
			}

			require.NoError(t, err)

			if tt.validate != nil {
				tt.validate(t, value)
			} else if tt.wantValue != nil {
				require.Equal(t, tt.wantValue, value)
			}
		})
	}
}

func runDeleteTests(t *testing.T, store storage.Store, tests []deleteTest) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupKey != "" {
				err := store.Save(tt.setupKey, tt.setupVal)
				require.NoError(t, err)
			}

			err := store.Delete(tt.deleteKey)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				return
			}

			require.NoError(t, err)

			_, err = store.Retrieve(tt.deleteKey)
			require.ErrorIs(t, err, storage.ErrKeyNotFound)
		})
	}
}

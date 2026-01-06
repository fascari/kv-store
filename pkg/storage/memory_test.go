package storage_test

import (
	"testing"

	"github.com/felipeascari/kv-store/pkg/storage"
)

func TestMemory_Save(t *testing.T) {
	store := storage.NewMemory()
	runSaveTests(t, store, saveTests)
}

func TestMemory_Retrieve(t *testing.T) {
	store := storage.NewMemory()
	runRetrieveTests(t, store, retrieveTests())
}

func TestMemory_Delete(t *testing.T) {
	store := storage.NewMemory()
	runDeleteTests(t, store, deleteTests)
}

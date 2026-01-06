package storage

import (
	"errors"
	"sync"
)

var ErrKeyNotFound = errors.New("key not found")

type Memory struct {
	mu      sync.RWMutex
	storage MemoryStorage
}

func NewMemory() *Memory {
	return &Memory{
		storage: NewMemoryStorage(),
	}
}

func (s *Memory) Save(key string, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.storage.Set(key, value)
	return nil
}

func (s *Memory) Retrieve(key string) (any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.storage.Get(key)
	if !exists {
		return nil, ErrKeyNotFound
	}

	return value, nil
}

func (s *Memory) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.storage.Get(key); !exists {
		return ErrKeyNotFound
	}

	s.storage.Delete(key)
	return nil
}

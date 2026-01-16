package storage

import "sync"

type Memory struct {
	mu    sync.RWMutex
	store map[string]any
}

func NewMemory() *Memory {
	return &Memory{
		store: make(map[string]any),
	}
}

func (m *Memory) Save(key string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[key] = value
	return nil
}

func (m *Memory) Retrieve(key string) (any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, exists := m.store[key]
	if !exists {
		return nil, ErrKeyNotFound
	}

	return value, nil
}

func (m *Memory) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.store[key]; !exists {
		return ErrKeyNotFound
	}

	delete(m.store, key)
	return nil
}

package storage

type MemoryStorage map[string]any

func NewMemoryStorage() MemoryStorage {
	return make(MemoryStorage)
}

func (m MemoryStorage) Set(key string, value any) {
	m[key] = value
}

func (m MemoryStorage) Get(key string) (any, bool) {
	value, exists := m[key]
	return value, exists
}

func (m MemoryStorage) Delete(key string) bool {
	if _, exists := m[key]; !exists {
		return false
	}
	delete(m, key)
	return true
}

func (m MemoryStorage) Has(key string) bool {
	_, exists := m[key]
	return exists
}

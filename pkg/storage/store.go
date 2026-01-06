package storage

// Store defines the interface for key-value storage operations
type Store interface {
	Save(key string, value any) error
	Retrieve(key string) (any, error)
	Delete(key string) error
}

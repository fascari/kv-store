package storage

type Store interface {
	Save(key string, value any) error
	Retrieve(key string) (any, error)
	Delete(key string) error
}

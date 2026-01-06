package delete

import "github.com/felipeascari/kv-store/pkg/storage"

type UseCase struct {
	store storage.Store
}

func NewUseCase(s storage.Store) UseCase {
	return UseCase{store: s}
}

func (u UseCase) Execute(key string) error {
	return u.store.Delete(key)
}

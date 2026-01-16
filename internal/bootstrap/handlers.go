package bootstrap

import (
	"github.com/felipeascari/kv-store/internal/handler/delete"
	"github.com/felipeascari/kv-store/internal/handler/retrieve"
	"github.com/felipeascari/kv-store/internal/handler/save"
	deleteUseCase "github.com/felipeascari/kv-store/internal/usecase/delete"
	retrieveUseCase "github.com/felipeascari/kv-store/internal/usecase/retrieve"
	saveUseCase "github.com/felipeascari/kv-store/internal/usecase/save"
	"github.com/felipeascari/kv-store/pkg/storage"
)

type Handlers struct {
	Save     *save.Handler
	Retrieve *retrieve.Handler
	Delete   *delete.Handler
}

func NewHandlers(store storage.Store) *Handlers {
	saveUC := saveUseCase.NewUseCase(store)
	retrieveUC := retrieveUseCase.NewUseCase(store)
	deleteUC := deleteUseCase.NewUseCase(store)

	return &Handlers{
		Save:     save.New(saveUC),
		Retrieve: retrieve.New(retrieveUC),
		Delete:   delete.New(deleteUC),
	}
}

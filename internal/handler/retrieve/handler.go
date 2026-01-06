package retrieve

import (
	"errors"
	"net/http"

	"github.com/felipeascari/kv-store/internal/usecase/retrieve"
	pkghttp "github.com/felipeascari/kv-store/pkg/http"
	"github.com/felipeascari/kv-store/pkg/storage"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	useCase retrieve.UseCase
}

func NewHandler(useCase retrieve.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "" {
		pkghttp.BadRequest(w, "key is required")
		return
	}

	value, err := h.useCase.Execute(key)
	if err != nil {
		if errors.Is(err, storage.ErrKeyNotFound) {
			pkghttp.NotFound(w, "key not found")
			return
		}
		pkghttp.InternalServerError(w, "internal server error")
		return
	}

	pkghttp.JSON(w, http.StatusOK, Response{
		Key:   key,
		Value: value,
	})
}

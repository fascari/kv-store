package save

import (
	"encoding/json"
	"net/http"

	"github.com/felipeascari/kv-store/internal/usecase/save"
	pkghttp "github.com/felipeascari/kv-store/pkg/http"
)

type Handler struct {
	useCase save.UseCase
}

func New(useCase save.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var req Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.BadRequest(w, "invalid request body")
		return
	}

	if req.Key == "" {
		pkghttp.BadRequest(w, "key is required")
		return
	}

	if err := h.useCase.Execute(req.Key, req.Value); err != nil {
		pkghttp.InternalServerError(w, "failed to save key")
		return
	}

	pkghttp.JSON(w, http.StatusCreated, Response{
		Key:   req.Key,
		Value: req.Value,
	})
}

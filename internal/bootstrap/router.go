package bootstrap

import (
	"net/http"

	pkghttp "github.com/felipeascari/kv-store/pkg/http"
	"github.com/felipeascari/kv-store/pkg/middleware"
	"github.com/felipeascari/kv-store/pkg/storage"
	"github.com/go-chi/chi/v5"
)

func setupRouter(kvStore storage.Store) *chi.Mux {
	r := chi.NewRouter()

	middleware.Setup(r)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		pkghttp.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	handlers := NewHandlers(kvStore)

	r.Route("/api", func(r chi.Router) {
		r.Post("/keys", handlers.Save.Handle)
		r.Get("/keys/{key}", handlers.Retrieve.Handle)
		r.Delete("/keys/{key}", handlers.Delete.Handle)
	})

	return r
}

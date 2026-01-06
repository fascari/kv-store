package bootstrap

import (
	"fmt"
	"net/http"

	"github.com/felipeascari/kv-store/pkg/config"
	"github.com/felipeascari/kv-store/pkg/logger"
	"github.com/felipeascari/kv-store/pkg/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	router  *chi.Mux
	addr    string
	logger  *zap.Logger
	storage storage.Store
}

func NewServer() (*Server, error) {
	if err := logger.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	kvStorage, err := createStorage(cfg.Storage)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	router := setupRouter(kvStorage)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)

	logger.Logger().Info("initialized storage", zap.String("type", cfg.Storage.Type.String()))

	return &Server{
		router:  router,
		addr:    addr,
		logger:  logger.Logger(),
		storage: kvStorage,
	}, nil
}

func createStorage(cfg config.StorageConfig) (storage.Store, error) {
	switch cfg.Type {
	case storage.TypeRedis:
		return storage.NewRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	case storage.TypeMemory:
		return storage.NewMemory(), nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Type)
	}
}

func (s *Server) Start() error {
	s.logger.Info("starting server", zap.String("address", s.addr))
	return http.ListenAndServe(s.addr, s.router)
}

func (s *Server) Address() string {
	return fmt.Sprintf("http://localhost%s", s.addr)
}

func (s *Server) Shutdown() {
	s.logger.Info("shutting down server")

	if redis, ok := s.storage.(*storage.Redis); ok {
		if err := redis.Close(); err != nil {
			s.logger.Error("failed to close Redis connection", zap.Error(err))
		}
	}

	logger.Sync()
}

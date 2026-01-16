package bootstrap

import (
	"fmt"
	"net/http"
	"time"

	"github.com/felipeascari/kv-store/pkg/config"
	"github.com/felipeascari/kv-store/pkg/lock"
	"github.com/felipeascari/kv-store/pkg/logger"
	"github.com/felipeascari/kv-store/pkg/storage"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Server struct {
	router      *chi.Mux
	addr        string
	logger      *zap.Logger
	store       storage.Store
	redisClient *redis.Client
}

func NewServer() (*Server, error) {
	if err := logger.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	kvStore, redisClient, err := createStorage(cfg.Storage)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	router := setupRouter(kvStore)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)

	logger.Logger().Info("initialized storage", zap.String("type", cfg.Storage.Type.String()))

	return &Server{
		router:      router,
		addr:        addr,
		logger:      logger.Logger(),
		store:       kvStore,
		redisClient: redisClient,
	}, nil
}

func createStorage(cfg config.StorageConfig) (storage.Store, *redis.Client, error) {
	switch cfg.Type {
	case storage.TypeRedis:
		redisStore, err := storage.NewRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
		if err != nil {
			return nil, nil, err
		}

		// Wrap Redis store with distributed locking (fencing tokens)
		// This prevents zombie processes and ensures consistency
		lockMgr := lock.NewManager(
			lock.NewRedisLock(redisStore.Client(), 5*time.Second),
		)
		lockedStore := storage.NewLockedStore(redisStore, lockMgr)

		return lockedStore, redisStore.Client(), nil

	case storage.TypeMemory:
		return storage.NewMemory(), nil, nil

	default:
		return nil, nil, fmt.Errorf("unknown storage type: %s", cfg.Type)
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

	if s.redisClient != nil {
		if err := s.redisClient.Close(); err != nil {
			s.logger.Error("failed to close Redis connection", zap.Error(err))
		}
	}

	logger.Sync()
}

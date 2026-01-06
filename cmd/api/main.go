package main

import (
	"log"

	"github.com/felipeascari/kv-store/internal/bootstrap"
	"github.com/felipeascari/kv-store/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	server, err := bootstrap.NewServer()
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	defer server.Shutdown()

	if err := server.Start(); err != nil {
		logger.Logger().Fatal("server failed to start", zap.Error(err))
	}
}

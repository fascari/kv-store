package http

import (
	"encoding/json"
	"net/http"

	"github.com/felipeascari/kv-store/pkg/logger"
	"go.uber.org/zap"
)

func JSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Logger().Error("failed to encode JSON response", zap.Error(err))
	}
}

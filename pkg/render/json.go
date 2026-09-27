package render

import (
	"encoding/json"
	"net/http"
	"project/pkg/logger"

	"go.uber.org/zap"
)

func JSON(w http.ResponseWriter, body any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, "json encode error", http.StatusInternalServerError)
		logger.Log.Error("json encode error", zap.Error(err))
	}
}

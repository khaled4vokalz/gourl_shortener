package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/khaled4vokalz/gourl_shortener/internal/cache"
	"github.com/khaled4vokalz/gourl_shortener/internal/db"
)

type HealthCheckResponse struct {
	Database bool `json:"database_is_alive"`
	Cache    bool `json:"cache_is_alive"`
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request, storage db.Storage, cache cache.Cache) {
	storage_is_alive := storage.IsAlive()
	cache_is_alive := cache.IsAlive()

	response := HealthCheckResponse{
		Database: storage_is_alive,
		Cache:    cache_is_alive,
	}

	statusCode := http.StatusOK
	if !storage_is_alive || !cache_is_alive {
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

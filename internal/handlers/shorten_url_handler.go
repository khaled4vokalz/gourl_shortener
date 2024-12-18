package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/khaled4vokalz/gourl_shortener/internal/cache"
	"github.com/khaled4vokalz/gourl_shortener/internal/config"
	"github.com/khaled4vokalz/gourl_shortener/internal/db"
	logger "github.com/khaled4vokalz/gourl_shortener/internal/logging"
	"github.com/khaled4vokalz/gourl_shortener/internal/middlewares"
	"github.com/khaled4vokalz/gourl_shortener/internal/service"
	"github.com/khaled4vokalz/gourl_shortener/internal/utils"
)

type Request struct {
	URL       string     `json:"url"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

func ShortenUrlHandler(w http.ResponseWriter, r *http.Request, storage db.Storage, cache cache.Cache, settings config.ShortenerSettings) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestId = middlewares.GetRequestID(r.Context())
	var request Request

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logger.GetLogger().Infow("-----------------------", "r.expires_at", request)

	if !utils.IsValidURL(request.URL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	baseURL := "http://localhost:8080" // Fallback in case Origin is not available
	if origin := r.Header.Get("Origin"); origin != "" {
		baseURL = origin
	} else if host := r.Host; host != "" {
		baseURL = fmt.Sprintf("http://%s", host)
	}
	length := settings.Length
	shortened := service.GenerateShortenedURL(request.URL, settings.Length)
	url, err := storage.Get(shortened)
	var attempt_count int8 = 1
	for url != "" {
		if attempt_count > settings.MaxAttempt {
			// bail out, we can not try more than allowed max attempt
			http.Error(w, fmt.Sprintf("Failed to generate a unique URL, attempted %d times :(", settings.MaxAttempt), http.StatusInternalServerError)
			return
		}
		logger.GetLogger().Debugw(fmt.Sprintf("Key '%s' is not unique, generating new one.", shortened), "request-id", requestId, "url", request.URL)
		length++
		shortened = service.GenerateShortenedURL(request.URL, length)
		url, err = storage.Get(shortened)
		if err != nil {
			logger.GetLogger().Errorw("Failed to get URL info from database", "error", err)
			http.Error(w, fmt.Sprintf("Failed to generate a shortened URL"), http.StatusInternalServerError)
			return
		}
		attempt_count++
	}

	err = storage.Save(shortened, request.URL, *request.ExpiresAt)
	if err != nil {
		logger.GetLogger().Errorw(fmt.Sprintf("Failed to store URL '%s' in database.", request.URL), "request-id", requestId, "error", err)
		http.Error(w, "Failed to store URL", http.StatusInternalServerError)
		return
	}

	// Cache the URL (keep it for 10 minutes now. TODO: make this configurable)
	err = cache.Set(shortened, request.URL, 10*time.Minute)
	if err != nil {
		logger.GetLogger().Errorw(fmt.Sprintf("Failed to cache URL '%s'", request.URL), "request-id", requestId, "error", err)
		http.Error(w, "Failed to cache URL", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"shortened_url": fmt.Sprintf("%s/%s", baseURL, shortened),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

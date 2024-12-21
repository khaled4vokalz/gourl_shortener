package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/khaled4vokalz/gourl_shortener/internal/cache"
	errors "github.com/khaled4vokalz/gourl_shortener/internal/common"
	"github.com/khaled4vokalz/gourl_shortener/internal/config"
	"github.com/khaled4vokalz/gourl_shortener/internal/db"
	logger "github.com/khaled4vokalz/gourl_shortener/internal/logging"
	"github.com/khaled4vokalz/gourl_shortener/internal/middlewares"
	"github.com/khaled4vokalz/gourl_shortener/internal/service"
	"github.com/khaled4vokalz/gourl_shortener/internal/utils"
)

type Request struct {
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
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

	if !utils.IsValidURL(request.URL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	baseURL := "http://localhost:8080" // Fallback in case Host is not available
	if host := r.Host; host != "" {
		baseURL = fmt.Sprintf("http://%s", host)
	}
	length := settings.Length
	var attempt_count int8 = 1
	var shortened string
	for true {
		shortened = service.GenerateShortenedURL(request.URL, length)
		_, err := storage.Get(shortened)
		if err == errors.NotFound || err == errors.Expired {
			// so we don't have duplicates in the DB
			break
		} else if err != nil {
			logger.GetLogger().Errorw(fmt.Sprintf("Failed to query database for key '%s'", shortened), "request-id", requestId, "error", err)
			http.Error(w, "Error fetching from database", http.StatusInternalServerError)
			return
		} else if attempt_count > settings.MaxAttempt {
			// bail out, we can not try more than allowed max attempt
			// TODO: this doesn't look like a good solution, we should figure out a diffirent way
			// so that we must be able to create a unique key
			http.Error(w, fmt.Sprintf("Failed to generate a unique URL, attempted %d times :(", settings.MaxAttempt), http.StatusInternalServerError)
			return
		} else {
			logger.GetLogger().Debugw(fmt.Sprintf("Key '%s' is not unique, generating new one.", shortened), "request-id", requestId, "url", request.URL)
			attempt_count++
			length++
		}
	}

	err := storage.Save(shortened, request.URL, request.ExpiresAt)
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

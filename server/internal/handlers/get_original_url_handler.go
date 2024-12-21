package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/khaled4vokalz/gourl_shortener/internal/cache"
	errors "github.com/khaled4vokalz/gourl_shortener/internal/common"
	"github.com/khaled4vokalz/gourl_shortener/internal/db"
	logger "github.com/khaled4vokalz/gourl_shortener/internal/logging"
	"github.com/khaled4vokalz/gourl_shortener/internal/middlewares"
)

func GetOriginalUrlHandler(w http.ResponseWriter, r *http.Request, storage db.Storage, cache cache.Cache) {
	log := logger.GetLogger()
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	key := parts[len(parts)-1]

	// is it in the cache?
	originalURL, err := cache.Get(key)

	var requestId = middlewares.GetRequestID(r.Context())

	if err == errors.NotFound {
		// Cache miss, fallback to database
		originalURL, err = storage.Get(key)
		if err == errors.NotFound {
			http.Error(w, fmt.Sprintf("URL not found for key: %s", key), http.StatusNotFound)
			return
		} else if err == errors.Expired {
			http.Error(w, fmt.Sprintf("URL expired key: %s", key), http.StatusGone)
			return
		} else if err != nil {
			logger.GetLogger().Errorw(fmt.Sprintf("Failed to query database for key '%s'", key), "request-id", requestId, "error", err)
			http.Error(w, "Error fetching from database", http.StatusInternalServerError)
			return
		}
		w.Header().Set("X-Cache-Status", "Miss")

		// Cache the URL (keep it for 10 minutes now. TODO: make this configurable)
		cache.Set(key, originalURL, 10*time.Minute)
	} else if err != nil {
		logger.GetLogger().Errorw(fmt.Sprintf("Failed to query cache for key '%s'", key), "request-id", requestId, "error", err)
		http.Error(w, "Error fetching from cache", http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("X-Cache-Status", "Hit")
	}

	log.Debugw(fmt.Sprintf("Redirect url is %s", originalURL), "request-id", requestId)
	http.Redirect(w, r, originalURL, http.StatusPermanentRedirect)
}

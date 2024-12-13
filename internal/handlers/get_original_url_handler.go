package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/khaled4vokalz/gourl_shortener/internal/cache"
	errors "github.com/khaled4vokalz/gourl_shortener/internal/common"
	"github.com/khaled4vokalz/gourl_shortener/internal/db"
)

func GetOriginalUrlHandler(w http.ResponseWriter, r *http.Request, storage db.Storage, cache cache.Cache) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	key := parts[len(parts)-1]

	// is it in the cache?
	originalURL, err := cache.Get(key)

	if err == errors.NotFound {
		// Cache miss, fallback to database
		originalURL, exists := storage.Get(key)
		if exists == false {
			http.Error(w, fmt.Sprintf("URL not found for key: %s", key), http.StatusNotFound)
			return
		}
		w.Header().Set("X-Cache-Status", "Miss")

		// Cache the URL (keep it for 10 minutes now. TODO: make this configurable)
		cache.Set(key, originalURL, 10*time.Minute)
	} else if err != nil {
		http.Error(w, "Error fetching from cache", http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("X-Cache-Status", "Hit")
	}

	http.Redirect(w, r, originalURL, http.StatusPermanentRedirect)
}

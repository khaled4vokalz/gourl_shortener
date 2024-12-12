package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/khaled4vokalz/gourl_shortener/internal/db"
)

func GetOriginalUrlHandler(w http.ResponseWriter, r *http.Request, storage db.Storage) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	key := parts[len(parts)-1]

	originalURL, exists := storage.Get(key)
	if exists == false {
		http.Error(w, fmt.Sprintf("URL not found for key: %s", key), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

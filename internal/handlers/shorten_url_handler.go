package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/khaled4vokalz/gourl_shortener/internal/db"
	"github.com/khaled4vokalz/gourl_shortener/internal/service"
	"github.com/khaled4vokalz/gourl_shortener/internal/utils"
)

type Request struct {
	URL string `json:"url"` // I like this tag thing in go :+1:
}

func ShortenUrlHandler(w http.ResponseWriter, r *http.Request, storage db.Storage) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

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
	baseURL := "http://localhost:8080" // Fallback in case Origin is not available
	if origin := r.Header.Get("Origin"); origin != "" {
		baseURL = origin
	} else if host := r.Host; host != "" {
		baseURL = fmt.Sprintf("http://%s", host)
	}
	shortened := service.GenerateShortenedURL(request.URL)

	err := storage.Save(shortened, request.URL)
	if err != nil {
		http.Error(w, "Failed to store URL", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"shortened_url": fmt.Sprintf("%s/%s", baseURL, shortened),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

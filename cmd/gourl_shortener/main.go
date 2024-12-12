package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/khaled4vokalz/gourl_shortener/internal/db"
	"github.com/khaled4vokalz/gourl_shortener/internal/handlers"
)

func main() {
	port := "8080"
	// DI :D
	storage := db.NewInMemoryDb()

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShortenUrlHandler(w, r, storage)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/") {
			handlers.GetOriginalUrlHandler(w, r, storage)
		}
	})
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Couldn't start the server because of: %s", err)
	}
}

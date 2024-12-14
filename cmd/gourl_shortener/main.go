package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/khaled4vokalz/gourl_shortener/internal/cache"
	config_loader "github.com/khaled4vokalz/gourl_shortener/internal/config"
	"github.com/khaled4vokalz/gourl_shortener/internal/db"
	"github.com/khaled4vokalz/gourl_shortener/internal/handlers"
)

func main() {
	LoadEnv()
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "DEV"
	}

	loaded_config, _ := config_loader.LoadConfig(fmt.Sprintf("configuration/%s.yaml", strings.ToLower(env)))

	storage, err := db.GetDb(loaded_config.Storage)
	if err != nil {
		log.Fatalf("failed to initialize storage: %s", err)
	}
	cache := cache.GetCache(loaded_config.Cache)

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShortenUrlHandler(w, r, storage, cache, loaded_config.ShortenerProps)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/") {
			handlers.GetOriginalUrlHandler(w, r, storage, cache)
		}
	})

	port := loaded_config.Server.Port // TODO: this should have a default if not set
	host := loaded_config.Server.Host // TODO: this should have a default if not set

	if err := http.ListenAndServe(host+":"+port, nil); err != nil {
		log.Fatalf("Couldn't start the server because of: %s", err)
	}
}

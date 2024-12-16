package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/khaled4vokalz/gourl_shortener/internal/cache"
	config_loader "github.com/khaled4vokalz/gourl_shortener/internal/config"
	"github.com/khaled4vokalz/gourl_shortener/internal/db"
	"github.com/khaled4vokalz/gourl_shortener/internal/handlers"
	logger "github.com/khaled4vokalz/gourl_shortener/internal/logging"
	"github.com/khaled4vokalz/gourl_shortener/internal/middlewares"
)

func main() {
	loaded_config, _ := config_loader.LoadConfig()
	log := logger.GetLogger()
	defer log.Sync()

	storage, err := db.GetDb(loaded_config.Storage)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to initialize storage: %s", err))
	}
	cache, err := cache.GetCache(loaded_config.Cache)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to initialize cache: %s", err))
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShortenUrlHandler(w, r, storage, cache, loaded_config.ShortenerProps)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/") {
			handlers.GetOriginalUrlHandler(w, r, storage, cache)
		}
	})

	port := loaded_config.Server.Port // TODO: this should have a default if not set
	host := loaded_config.Server.Host // TODO: this should have a default if not set

	if err := http.ListenAndServe(host+":"+port, middlewares.RequestIDMiddleware(mux)); err != nil {
		log.Fatal(fmt.Sprintf("Couldn't start the server because of: %s", err))
	}
}

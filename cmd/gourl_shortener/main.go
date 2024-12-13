package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

	config, _ := config_loader.LoadConfig(fmt.Sprintf("configuration/%s.yaml", strings.ToLower(env)))

	// DI :D
	var storage db.Storage
	var error error
	if config.Db_Conn_String != "" {
		storage, error = db.NewPostgresDb(config.Db_Conn_String)
		if error != nil {
			log.Fatalf("Failed to initialize postgresDB :: %s", error)
		}
	} else {
		storage = db.NewInMemoryDb()
	}

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShortenUrlHandler(w, r, storage)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/") {
			handlers.GetOriginalUrlHandler(w, r, storage)
		}
	})

	port := config.Server.Port // TODO: this should have a default if not set
	host := config.Server.Host // TODO: this should have a default if not set

	if err := http.ListenAndServe(host+":"+port, nil); err != nil {
		log.Fatalf("Couldn't start the server because of: %s", err)
	}
}

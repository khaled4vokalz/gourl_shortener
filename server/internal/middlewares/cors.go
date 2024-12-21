package middlewares

import (
	"net/http"

	"github.com/rs/cors"
)

func AddCORSMiddleware(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(next)
}

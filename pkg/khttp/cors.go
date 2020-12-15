package khttp

import (
	"github.com/rs/cors"
	"net/http"
)

func AddCorsMiddleware() func(handler http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "OPTIONS", "HEAD", "DELETE", "PATCH"},
	}).Handler
}

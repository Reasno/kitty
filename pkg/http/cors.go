package http

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
	}).Handler
}

package http // AddDocMiddleware returns a documentation path at /doc/

import (
	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
	"net/http"
)

func AddHealthCheck() func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		router := mux.NewRouter()
		router.PathPrefix("/live").Handler(healthcheck.NewHandler())
		router.PathPrefix("/ready").Handler(healthcheck.NewHandler())
		router.PathPrefix("/").Handler(handler)
		return router
	}
}

func RegisterHealthCheck(httpProviders *[]func(router *mux.Router), _ interface{}) {
	*httpProviders = append(*httpProviders, func(router *mux.Router) {
		router.PathPrefix("/live").Handler(healthcheck.NewHandler())
		router.PathPrefix("/ready").Handler(healthcheck.NewHandler())
	})
}

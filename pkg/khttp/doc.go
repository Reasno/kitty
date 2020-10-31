package khttp

import (
	"github.com/gorilla/mux"
	"net/http"
)

// AddDocMiddleware returns a documentation path at /doc/
func AddDocMiddleware() func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		router := mux.NewRouter()
		router.PathPrefix("/doc/").Handler(getOpenAPIHandler())
		router.PathPrefix("/doc").Handler(http.RedirectHandler("/doc/", 302))
		router.PathPrefix("/").Handler(handler)
		return router
	}
}

func Doc(router *mux.Router) {
	router.PathPrefix("/doc/").Handler(getOpenAPIHandler())
	router.PathPrefix("/doc").Handler(http.RedirectHandler("/doc/", 302))
}

// getOpenAPIHandler serves an OpenAPI UI.
func getOpenAPIHandler() http.Handler {
	return http.StripPrefix("/doc", http.FileServer(http.Dir("./doc")))
}

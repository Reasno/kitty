package khttp // AddDocMiddleware returns a documentation path at /doc/
import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func AddMetricMiddleware() func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		router := mux.NewRouter()
		router.PathPrefix("/metrics").Handler(promhttp.Handler())
		router.PathPrefix("/").Handler(handler)
		return router
	}
}

func Metrics(router *mux.Router) {
	router.PathPrefix("/metrics").Handler(promhttp.Handler())
}

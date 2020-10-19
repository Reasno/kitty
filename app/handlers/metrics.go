package handlers

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

func provideHistogramMetrics() metrics.Histogram {
	var his metrics.Histogram = prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: viper.GetString("app_name"),
		Subsystem: viper.GetString("app_env"),
		Name:     "request_duration_seconds",
		Help:     "Total time spent serving requests.",
	}, []string{"service", "method"})
	return his
}

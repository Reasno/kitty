package handlers

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	logging "github.com/Reasno/kitty/pkg/log"
	jaegermetric "github.com/uber/jaeger-lib/metrics"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	"github.com/go-kit/kit/metrics"
)

func ProvideLogger() log.Logger {
	return log.With(logging.NewLogger(viper.GetString("app_env")), "service", "app")
}

func provideHistogramMetrics() metrics.Histogram {
	var his metrics.Histogram = prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: viper.GetString("app_name"),
		Subsystem: viper.GetString("app_env"),
		Name:     "request_duration_seconds",
		Help:     "Total time spent serving requests.",
	}, []string{"service", "method"})
	return his
}

type logAdapter struct {
	logging log.Logger
}

func (l logAdapter) Infof(msg string, args... interface{})  {
	level.Info(l.logging).Log("msg", fmt.Sprintf(msg, args...))
}

func (l logAdapter) Error(msg string) {
	level.Error(l.logging).Log("msg", msg)
}
func provideJaegerLogAdatper(logging log.Logger) jaeger.Logger {
	return logAdapter{logging: logging}
}

var tracer opentracing.Tracer
func provideOpentracing(log jaeger.Logger) opentracing.Tracer {
	if tracer != nil {
		return tracer
	}
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  viper.GetString("jaeger.sampler.type"),
			Param: viper.GetFloat64("jaeger.sampler.param"),
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: viper.GetBool("jaeger.log.enable"),
		},
	}
	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := log
	jMetricsFactory := jaegermetric.NullFactory

	// Initialize tracer with a logger and a metrics factory
	var err error
	jaegerCloser, err = cfg.InitGlobalTracer(
		viper.GetString("app_name"),
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Error(fmt.Sprintf("Could not initialize jaeger tracer: %s", err.Error()))
	}
	tracer = opentracing.GlobalTracer()
	return tracer
}

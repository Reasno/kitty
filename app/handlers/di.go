package handlers

import (
	"fmt"
	logging "github.com/Reasno/kitty/pkg/log"
	"github.com/Reasno/kitty/pkg/middleware"
	"github.com/Reasno/kitty/pkg/sms"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegermetric "github.com/uber/jaeger-lib/metrics"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
)

func ProvideLogger() log.Logger {
	return log.With(logging.NewLogger(viper.GetString("app_env")), "service", "app")
}

func provideHistogramMetrics() metrics.Histogram {
	var his metrics.Histogram = prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: viper.GetString("app_name"),
		Subsystem: viper.GetString("app_env"),
		Name:      "request_duration_seconds",
		Help:      "Total time spent serving requests.",
	}, []string{"service", "method"})
	return his
}

func provideSeurityConfig() *middleware.SecurityConfig {
	return &middleware.SecurityConfig{
		Enable: viper.GetBool("security.enable"),
		JwtKey: viper.GetString("security.key"),
		JwtId:  viper.GetString("security.kid"),
	}
}

func provideSmsConfig() *sms.SenderConfig {
	return &sms.SenderConfig{
		Tag:        viper.GetString("sms.tag"),
		SendUrl:    viper.GetString("sms.send_url"),
		BalanceUrl: viper.GetString("sms.balance_url"),
		UserName:   viper.GetString("sms.username"),
		Password:   viper.GetString("sms.password"),
	}
}

func provideRedis() *redis.Client {
	return redis.NewClient(&redis.Options{})
}

func provideDialector() gorm.Dialector {
	return sqlite.Open("test.db")
}

func provideGormConfig(l log.Logger) *gorm.Config {
	return &gorm.Config{
		Logger: &logging.GormLogAdapter{l},
	}
}

func provideJaegerLogAdatper(l log.Logger) jaeger.Logger {
	return &logging.JaegerLogAdapter{Logging: l}
}

func provideOpentracing(log jaeger.Logger) opentracing.Tracer {
	if opentracing.IsGlobalTracerRegistered() {
		return opentracing.GlobalTracer()
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
	var (
		tracer   opentracing.Tracer
		canceler io.Closer
		err      error
	)
	canceler, err = cfg.InitGlobalTracer(
		viper.GetString("app_name"),
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Error(fmt.Sprintf("Could not initialize jaeger tracer: %s", err.Error()))
		panic(err)
	}
	tracer = opentracing.GlobalTracer()
	destruct.Add(func() {
		if err := canceler.Close(); err != nil {
			log.Error(err.Error())
		}
	})
	return tracer
}

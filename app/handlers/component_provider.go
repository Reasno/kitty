package handlers

import (
	"fmt"
	"github.com/Reasno/kitty/pkg/config"
	kittyhttp "github.com/Reasno/kitty/pkg/http"
	logging "github.com/Reasno/kitty/pkg/log"
	"github.com/Reasno/kitty/pkg/middleware"
	"github.com/Reasno/kitty/pkg/otgorm"
	"github.com/Reasno/kitty/pkg/otredis"
	"github.com/Reasno/kitty/pkg/sms"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
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
	"sync"
)

func provideFoo() func() string {
	return func() string {
		return "Foo"
	}
}

func provideConfig() (*viper.Viper, error) {
	return config.ProvideChildConfig("app", "global")
}

func provideLogger(conf *viper.Viper) log.Logger {
	return log.With(logging.NewLogger(conf.GetString("env")), "service", "app")
}

func provideHistogramMetrics(conf *viper.Viper) metrics.Histogram {
	var his metrics.Histogram = prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: conf.GetString("name"),
		Subsystem: conf.GetString("env"),
		Name:      "request_duration_seconds",
		Help:      "Total time spent serving requests.",
	}, []string{"service", "method"})
	return his
}

func provideHttpClient(tracer opentracing.Tracer) *kittyhttp.Client {
	return kittyhttp.NewClient(tracer)
}

func provideSecurityConfig(conf *viper.Viper) *middleware.SecurityConfig {
	return &middleware.SecurityConfig{
		Enable: conf.GetBool("security.enable"),
		JwtKey: conf.GetString("security.key"),
		JwtId:  conf.GetString("security.kid"),
	}
}

func provideSmsConfig(doer kittyhttp.Doer, conf *viper.Viper) *sms.TransportConfig {
	return &sms.TransportConfig{
		Tag:        conf.GetString("sms.tag"),
		SendUrl:    conf.GetString("sms.send_url"),
		BalanceUrl: conf.GetString("sms.balance_url"),
		UserName:   conf.GetString("sms.username"),
		Password:   conf.GetString("sms.password"),
		Client:     doer,
	}
}

func provideRedis(logging log.Logger, conf *viper.Viper) (redis.UniversalClient, func()) {
	client := redis.NewUniversalClient(
		&redis.UniversalOptions{
			Addrs: conf.GetStringSlice("redis.addrs"),
			DB:    conf.GetInt("redis.database"),
		})
	client.AddHook(
		otredis.NewHook(conf.GetStringSlice("redis.addrs"),
			conf.GetInt("redis.database")))
	return client, func() {
		if err := client.Close(); err != nil {
			level.Error(logging).Log("err", err.Error())
		}
	}
}

func provideDialector() gorm.Dialector {
	return sqlite.Open("test.db")
}

func provideGormConfig(l log.Logger) *gorm.Config {
	return &gorm.Config{
		Logger: &logging.GormLogAdapter{l},
	}
}

func provideGormDB(dialector gorm.Dialector, config *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, err
	}
	otgorm.AddGormCallbacks(db)
	return db, nil
}

func provideJaegerLogAdatper(l log.Logger) jaeger.Logger {
	return &logging.JaegerLogAdapter{Logging: l}
}

var (
	closer     func()
	tracer     opentracing.Tracer
	onceTracer sync.Once
)
func provideOpentracing(log jaeger.Logger, conf *viper.Viper) (opentracing.Tracer, func()) {
	onceTracer.Do(func() {
		cfg := jaegercfg.Configuration{
			ServiceName: conf.GetString("name"),
			Sampler: &jaegercfg.SamplerConfig{
				Type:  conf.GetString("jaeger.sampler.type"),
				Param: conf.GetFloat64("jaeger.sampler.param"),
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans: conf.GetBool("jaeger.log.enable"),
			},
		}
		// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
		// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
		// frameworks.
		jLogger := log
		jMetricsFactory := jaegermetric.NullFactory

		// Initialize tracer with a logger and a metrics factory
		var (
			canceler io.Closer
			err      error
		)
		tracer, canceler, err = cfg.NewTracer(jaegercfg.Logger(jLogger), jaegercfg.Metrics(jMetricsFactory))
		if err != nil {
			log.Error(fmt.Sprintf("Could not initialize jaeger tracer: %s", err.Error()))
			panic(err)
		}
		closer = func() {
			if err := canceler.Close(); err != nil {
				log.Error(err.Error())
			}
		}
	})
	return tracer, closer
}

package ots3

import (
	"fmt"
	"github.com/Reasno/kitty/pkg/contract"
	kittyhttp "github.com/Reasno/kitty/pkg/khttp"
	logging "github.com/Reasno/kitty/pkg/klog"
	"github.com/go-kit/kit/log"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegermetric "github.com/uber/jaeger-lib/metrics"
	"io"
)

func provideHttpClient(tracer opentracing.Tracer) *kittyhttp.Client {
	return kittyhttp.NewClient(tracer)
}

func provideJaegerLogAdapter(l log.Logger) jaeger.Logger {
	return &logging.JaegerLogAdapter{Logging: l}
}

func provideOpentracing(log jaeger.Logger, conf contract.ConfigReader) (opentracing.Tracer, func(), error) {
	cfg := jaegercfg.Configuration{
		ServiceName: conf.String("name"),
		Sampler: &jaegercfg.SamplerConfig{
			Type:  conf.String("jaeger.sampler.type"),
			Param: conf.Float64("jaeger.sampler.param"),
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: conf.Bool("jaeger.log.enable"),
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
	tracer, canceler, err := cfg.NewTracer(jaegercfg.Logger(jLogger), jaegercfg.Metrics(jMetricsFactory))
	if err != nil {
		log.Error(fmt.Sprintf("Could not initialize jaeger tracer: %s", err.Error()))
		return nil, nil, err
	}
	closer := func() {
		if err := canceler.Close(); err != nil {
			log.Error(err.Error())
		}
	}

	return tracer, closer, nil
}

func provideUploadManager(conf contract.ConfigReader) *Manager {
	return NewManager(
		conf.String("s3.accessKey"),
		conf.String("s3.accessSecret"),
		conf.String("s3.region"),
		conf.String("s3.endpoint"),
		conf.String("s3.bucket"),
		WithLocationFunc(func(location string) (url string) {
			return fmt.Sprintf(conf.String("s3.cdnUrl"), location)
		}),
	)
}

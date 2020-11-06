package handlers

import (
	"fmt"
	"io"

	"github.com/Reasno/kitty/app/svc"
	"github.com/Reasno/kitty/pkg/contract"
	kittyhttp "github.com/Reasno/kitty/pkg/khttp"
	"github.com/Reasno/kitty/pkg/kkafka"
	logging "github.com/Reasno/kitty/pkg/klog"
	"github.com/Reasno/kitty/pkg/kmiddleware"
	"github.com/Reasno/kitty/pkg/otgorm"
	"github.com/Reasno/kitty/pkg/otredis"
	"github.com/Reasno/kitty/pkg/ots3"
	"github.com/Reasno/kitty/pkg/sms"
	"github.com/Reasno/kitty/pkg/wechat"
	kitty "github.com/Reasno/kitty/proto"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegermetric "github.com/uber/jaeger-lib/metrics"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func provideHistogramMetrics(appName contract.AppName, env contract.Env) metrics.Histogram {
	var his metrics.Histogram = prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: appName.String(),
		Subsystem: env.String(),
		Name:      "request_duration_seconds",
		Help:      "Total time spent serving requests.",
	}, []string{"module", "method"})
	return his
}

func provideKeyManager(appName contract.AppName, env contract.Env) otredis.KeyManager {
	return otredis.NewKeyManager(":", appName.String(), env.String())
}

func provideHttpClient(tracer opentracing.Tracer) *kittyhttp.Client {
	return kittyhttp.NewClient(tracer)
}

type userBus struct {
	kkafka.DataStore
}

func provideUserBus(factory *kkafka.KafkaProducerFactory, conf contract.ConfigReader) *userBus {
	return &userBus{kkafka.DataStore{
		Factory: factory,
		Topic:   conf.String("kafka.userBus"),
	}}
}

type eventBus struct {
	kkafka.EventStore
}

func provideEventBus(factory *kkafka.KafkaProducerFactory, conf contract.ConfigReader) *eventBus {
	return &eventBus{kkafka.EventStore{
		Factory: factory,
		Topic:   conf.String("kafka.eventBus"),
	}}
}

func provideKafkaProducerFactory(conf contract.ConfigReader, logger log.Logger, tracer opentracing.Tracer) (*kkafka.KafkaProducerFactory, func()) {
	factory := kkafka.NewKafkaProducerFactoryWithTracer(conf.Strings("kafka.brokers"), logger, tracer)
	return factory, func() {
		_ = factory.Close()
	}
}

func provideUploadManager(tracer opentracing.Tracer, conf contract.ConfigReader, client contract.HttpDoer) *ots3.Manager {
	return ots3.NewManager(
		conf.String("s3.accessKey"),
		conf.String("s3.accessSecret"),
		conf.String("s3.region"),
		conf.String("s3.endpoint"),
		conf.String("s3.bucket"),
		ots3.WithTracer(tracer),
		ots3.WithHttpClient(client),
		ots3.WithLocationFunc(func(location string) (url string) {
			return fmt.Sprintf(conf.String("s3.cdnUrl"), location)
		}),
	)
}

func provideSecurityConfig(conf contract.ConfigReader) *kmiddleware.SecurityConfig {
	return &kmiddleware.SecurityConfig{
		Enable: conf.Bool("security.enable"),
		JwtKey: conf.String("security.key"),
		JwtId:  conf.String("security.kid"),
	}
}

func provideWechatConfig(conf contract.ConfigReader, client contract.HttpDoer) *wechat.WechatConfig {
	return &wechat.WechatConfig{
		WechatAccessTokenUrl: conf.String("wechat.wechatAccessTokenUrl"),
		WeChatGetUserInfoUrl: conf.String("wechat.weChatGetUserInfoUrl"),
		AppId:                conf.String("wechat.appId"),
		AppSecret:            conf.String("wechat.appSecret"),
		Client:               client,
	}
}

func provideSmsConfig(doer contract.HttpDoer, conf contract.ConfigReader) *sms.TransportConfig {
	return &sms.TransportConfig{
		Tag:        conf.String("sms.tag"),
		SendUrl:    conf.String("sms.sendUrl"),
		BalanceUrl: conf.String("sms.balanceUrl"),
		UserName:   conf.String("sms.username"),
		Password:   conf.String("sms.password"),
		Client:     doer,
	}
}

func provideRedis(logging log.Logger, conf contract.ConfigReader, tracer opentracing.Tracer) (redis.UniversalClient, func()) {
	client := redis.NewUniversalClient(
		&redis.UniversalOptions{
			Addrs:    conf.Strings("redis.addrs"),
			DB:       conf.Int("redis.database"),
			Password: conf.String("redis.password"),
		})
	client.AddHook(
		otredis.NewHook(tracer, conf.Strings("redis.addrs"),
			conf.Int("redis.database")))
	return client, func() {
		if err := client.Close(); err != nil {
			level.Error(logging).Log("err", err.Error())
		}
	}
}

func provideDialector(conf contract.ConfigReader) (gorm.Dialector, error) {
	databaseType := conf.String("gorm.database")
	if databaseType == "mysql" {
		return mysql.Open(conf.String("gorm.dsn")), nil
	}
	if databaseType == "sqlite" {
		return sqlite.Open(conf.String("gorm.dsn")), nil
	}
	return nil, fmt.Errorf("unknow database type %s", databaseType)
}

func provideGormConfig(l log.Logger, conf contract.ConfigReader) *gorm.Config {
	return &gorm.Config{
		Logger: &logging.GormLogAdapter{l},
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: conf.String("name") + "_", // 表名前缀，`User` 的表名应该是 `t_users`
		},
	}
}

func provideGormDB(dialector gorm.Dialector, config *gorm.Config, tracer opentracing.Tracer) (*gorm.DB, func(), error) {
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, nil, err
	}
	otgorm.AddGormCallbacks(db, tracer)
	return db, func() {
		if sqlDb, err := db.DB(); err == nil {
			sqlDb.Close()
		}
	}, nil
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

type overallMiddleware func(endpoints svc.Endpoints) svc.Endpoints

func provideModule(db *gorm.DB, tracer opentracing.Tracer, logger log.Logger, middleware overallMiddleware, server kitty.AppServer, appName contract.AppName) *Module {
	return &Module{
		appName:   appName,
		db:        db,
		logger:    logger,
		tracer:    tracer,
		endpoints: middleware(NewEndpoints(server)),
	}
}

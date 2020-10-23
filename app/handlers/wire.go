//+build wireinject

package handlers

import (
	"github.com/Reasno/kitty/app/repository"
	kittyhttp "github.com/Reasno/kitty/pkg/http"
	"github.com/Reasno/kitty/pkg/sms"
	pb "github.com/Reasno/kitty/proto"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
)

var DbSet = wire.NewSet(
	provideDialector,
	provideGormConfig,
	provideGormDB,
)

var OpenTracingSet = wire.NewSet(
	provideJaegerLogAdatper,
	provideOpentracing,
)

func InjectDb() (*gorm.DB, error) {
	panic(wire.Build(ProvideLogger, DbSet))
}

func injectAppServer() (pb.AppServer, func(), error) {
	panic(wire.Build(
		ProvideLogger,
		provideSmsConfig,
		DbSet,
		OpenTracingSet,
		provideHttpClient,
		provideRedis,
		sms.NewTransport,
		repository.NewUserRepo,
		repository.NewCodeRepo,
		wire.Struct(new(appService), "*"),
		wire.Bind(new(redis.Cmdable), new(redis.UniversalClient)),
		wire.Bind(new(Sender), new(*sms.Transport)),
		wire.Bind(new(kittyhttp.Doer), new(*kittyhttp.Client)),
		wire.Bind(new(pb.AppServer), new(appService)),
		wire.Bind(new(UserRepository), new(*repository.UserRepo)),
		wire.Bind(new(CodeRepository), new(*repository.CodeRepo)),
	))
}

func InjectOpentracingTracer() opentracing.Tracer {
	panic(wire.Build(ProvideLogger, OpenTracingSet))
}

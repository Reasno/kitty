//+build wireinject

package handlers

import (
	"github.com/Reasno/kitty/app/repository"
	"github.com/Reasno/kitty/pkg/sms"
	pb "github.com/Reasno/kitty/proto"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
)

var DbSet = wire.NewSet(ProvideLogger, provideDialector, provideGormConfig, gorm.Open)

func InjectDb() (*gorm.DB, error) {
	panic(wire.Build(DbSet))
}

func injectAppServer() (pb.AppServer, error) {
	panic(wire.Build(
		provideSmsConfig,
		DbSet,
		provideRedis,
		sms.NewSender,
		repository.NewUserRepo,
		repository.NewCodeRepo,
		wire.Struct(new(appService), "log", "cr", "ur", "sender"),
		wire.Bind(new(redis.Cmdable), new(*redis.Client)),
		wire.Bind(new(Sender), new(*sms.Sender)),
		wire.Bind(new(pb.AppServer), new(appService)),
		wire.Bind(new(UserRepository), new(*repository.UserRepo)),
		wire.Bind(new(CodeRepository), new(*repository.CodeRepo)),
	))
}

func InjectOpentracingTracer() opentracing.Tracer {
	panic(wire.Build(
		ProvideLogger,
		provideJaegerLogAdatper,
		provideOpentracing,
	))
}

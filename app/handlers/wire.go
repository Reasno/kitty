//+build wireinject

package handlers

import (
	"github.com/Reasno/kitty/app/repository"
	"github.com/Reasno/kitty/pkg/contract"
	kittyhttp "github.com/Reasno/kitty/pkg/http"
	"github.com/Reasno/kitty/pkg/ots3"
	"github.com/Reasno/kitty/pkg/sms"
	"github.com/Reasno/kitty/pkg/wechat"
	pb "github.com/Reasno/kitty/proto"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var DbSet = wire.NewSet(
	provideDialector,
	provideGormConfig,
	provideGormDB,
)

var OpenTracingSet = wire.NewSet(
	provideJaegerLogAdapter,
	provideOpentracing,
)

var AppServerSet = wire.NewSet(
	provideLogger,
	provideSmsConfig,
	DbSet,
	OpenTracingSet,
	provideHttpClient,
	provideUploadManager,
	provideRedis,
	provideWechatConfig,
	wechat.NewTransport,
	sms.NewTransport,
	repository.NewUserRepo,
	repository.NewCodeRepo,
	repository.NewFileRepo,
	wire.Struct(new(appService), "*"),
	wire.Bind(new(redis.Cmdable), new(redis.UniversalClient)),
	wire.Bind(new(contract.SmsSender), new(*sms.Transport)),
	wire.Bind(new(contract.Uploader), new(*ots3.Manager)),
	wire.Bind(new(contract.HttpDoer), new(*kittyhttp.Client)),
	wire.Bind(new(pb.AppServer), new(appService)),
	wire.Bind(new(UserRepository), new(*repository.UserRepo)),
	wire.Bind(new(CodeRepository), new(*repository.CodeRepo)),
	wire.Bind(new(FileRepository), new(*repository.FileRepo)),
)

func injectModule(reader contract.ConfigReader) (*AppModule, func(), error) {
	panic(wire.Build(
		AppServerSet,
		provideSecurityConfig,
		provideHistogramMetrics,
		provideEndpointsMiddleware,
		provideModule))
}

func injectTestDb(conf contract.ConfigReader) (*gorm.DB, func(), error) {
	panic(wire.Build(
		provideLogger,
		DbSet,
	))
}

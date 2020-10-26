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
	"github.com/spf13/viper"
	"gorm.io/gorm"
)
var ConfigSet = wire.NewSet(provideConfig, wire.Bind(new(contract.ConfigReader), new(*viper.Viper)))

var DbSet = wire.NewSet(
	provideDialector,
	provideGormConfig,
	provideGormDB,
)

var OpenTracingSet = wire.NewSet(
	provideJaegerLogAdatper,
	provideOpentracing,
)

func injectDb() (*gorm.DB, error) {
	panic(wire.Build(ConfigSet, provideLogger, DbSet))
}

var AppServerSet = wire.NewSet(
	ConfigSet,
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
	wire.Struct(new(appService), "*"),
	wire.Bind(new(redis.Cmdable), new(redis.UniversalClient)),
	wire.Bind(new(contract.SmsSender), new(*sms.Transport)),
	wire.Bind(new(contract.Uploader), new(*ots3.Manager)),
	wire.Bind(new(contract.HttpDoer), new(*kittyhttp.Client)),
	wire.Bind(new(pb.AppServer), new(appService)),
	wire.Bind(new(UserRepository), new(*repository.UserRepo)),
	wire.Bind(new(CodeRepository), new(*repository.CodeRepo)),
)

func injectModule() (*AppModule, func(), error){
	panic(wire.Build(AppServerSet, provideSecurityConfig, provideHistogramMetrics, provideEndpointsMiddleware, provideModule))
}

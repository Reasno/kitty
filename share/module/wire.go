//+build wireinject

package module

import (
	"github.com/go-kit/kit/log"
	"github.com/google/wire"
	"glab.tagtic.cn/ad_gains/kitty/app/module"
	"glab.tagtic.cn/ad_gains/kitty/app/repository"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/event"
	"glab.tagtic.cn/ad_gains/kitty/pkg/invitecode"
	kittyhttp "glab.tagtic.cn/ad_gains/kitty/pkg/khttp"
	kclient "glab.tagtic.cn/ad_gains/kitty/pkg/kkafka/client"
	"glab.tagtic.cn/ad_gains/kitty/pkg/ots3"
	"glab.tagtic.cn/ad_gains/kitty/share/handlers"
	"glab.tagtic.cn/ad_gains/kitty/share/internal"
	"glab.tagtic.cn/ad_gains/kitty/share/listener"
)

var ShareServiceSet = wire.NewSet(
	module.DbSet,
	module.OpenTracingSet,
	module.NameAndEnvSet,
	module.ProvideSecurityConfig,
	module.ProvideKafkaFactory,
	module.ProvideHistogramMetrics,
	module.ProvideHttpClient,
	module.ProvideUploadManager,
	repository.NewUserRepo,
	repository.NewRelationRepo,
	repository.NewFileRepo,
	provideTokenizer,
	providePublisherOptions,
	provideInvitationCodeBus,
	provideDispatcher,
	internal.NewXTaskRequester,
	handlers.NewShareService,
	handlers.ProvideShareServer,
	wire.Struct(new(internal.InvitationManagerFactory), "*"),
	wire.Struct(new(internal.InvitationManagerFacade), "*"),
	wire.Bind(new(handlers.UserRepository), new(*repository.UserRepo)),
	wire.Bind(new(internal.RelationRepository), new(*repository.RelationRepo)),
	wire.Bind(new(handlers.InvitationManager), new(*internal.InvitationManagerFacade)),
	wire.Bind(new(contract.Uploader), new(*ots3.Manager)),
	wire.Bind(new(contract.HttpDoer), new(*kittyhttp.Client)),
	wire.Bind(new(internal.EncodeDecoder), new(*invitecode.Tokenizer)),
	wire.Bind(new(contract.Dispatcher), new(*event.Dispatcher)),
	wire.Bind(new(listener.InvitationCodeBus), new(*kclient.DataStore)),
)

func injectModule(reader contract.ConfigReader, logger log.Logger, dynConf config.DynamicConfigReader) (*Module, func(), error) {
	panic(wire.Build(
		ShareServiceSet,
		provideEndpointsMiddleware,
		provideEndpoints,
		provideHttp,
		provideGrpc,
		provideKafkaServer,
		provideProducerMiddleware,
		provideModule,
	))
}

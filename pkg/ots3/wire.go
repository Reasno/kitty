//+build wireinject

package ots3

import (
	"net/url"

	"github.com/go-kit/kit/log"
	"github.com/google/wire"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

func injectModule(conf contract.ConfigReader, logger log.Logger) *Module {
	panic(wire.Build(
		provideUploadManager,
		provideSecurityConfig,
		MakeUploadEndpoint,
		MakeHttpHandler,
		Middleware,
		wire.Struct(new(Module), "*"),
		wire.Struct(new(UploadService), "*"),
		config.ProvideEnv,
		wire.Bind(new(contract.Env), new(config.Env)),
		wire.Bind(new(contract.Uploader), new(*UploadService)),
	))
}

func InjectClientUploader(uri *url.URL) *ClientUploader {
	panic(wire.Build(NewClient, NewClientUploader))
}

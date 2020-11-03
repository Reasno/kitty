//+build wireinject

package ots3

import (
	"github.com/Reasno/kitty/pkg/contract"
	"github.com/go-kit/kit/log"
	"github.com/google/wire"
)

func injectModule(conf contract.ConfigReader, logger log.Logger) *Module {
	panic(wire.Build(
		provideUploadManager,
		MakeUploadEndpoint,
		MakeHttpHandler,
		wire.Struct(new(Module), "*"),
		wire.Struct(new(UploadService), "*"),
		wire.Bind(new(contract.Uploader), new(*UploadService)),
	))
}

func InjectClientUploader(conf contract.ConfigReader) *ClientUploader {
	panic(wire.Build(NewClient, NewClientUploader))
}

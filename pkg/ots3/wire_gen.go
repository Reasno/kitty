// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package ots3

import (
	"github.com/Reasno/kitty/pkg/contract"
	"github.com/go-kit/kit/log"
)

// Injectors from wire.go:

func injectModule(conf contract.ConfigReader, logger log.Logger) *Module {
	manager := provideUploadManager(conf)
	uploadService := &UploadService{
		logger: logger,
		s3:     manager,
	}
	endpoint := MakeUploadEndpoint(uploadService)
	handler := MakeHttpHandler(endpoint)
	module := &Module{
		handler: handler,
	}
	return module
}

func InjectClientUploader(conf contract.ConfigReader) *ClientUploader {
	client := NewClient(conf)
	clientUploader := NewClientUploader(client)
	return clientUploader
}

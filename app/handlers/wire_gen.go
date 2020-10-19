// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package handlers

import (
	"github.com/Reasno/kitty/proto"
)

// Injectors from wire.go:

func injectAppServer() kitty.AppServer {
	logger := provideLogger()
	handlersAppService := appService{
		log: logger,
	}
	return handlersAppService
}

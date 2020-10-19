//+build wireinject

package handlers

import (
	pb "github.com/Reasno/kitty/proto"
	"github.com/google/wire"
)

func injectAppServer() pb.AppServer {
	panic(wire.Build(
		provideLogger,
		wire.Struct(new(appService), "log"),
		wire.Bind(new(pb.AppServer), new(appService)),
	))
}

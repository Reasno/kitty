//+build wireinject

package handlers

import (
	pb "github.com/Reasno/kitty/proto"
	"github.com/google/wire"
	"github.com/opentracing/opentracing-go"
)

func injectAppServer() pb.AppServer {
	panic(wire.Build(
		ProvideLogger,
		wire.Struct(new(appService), "log"),
		wire.Bind(new(pb.AppServer), new(appService)),
	))
}

func InjectOpentracingTracer() opentracing.Tracer {
	panic(wire.Build(
		ProvideLogger,
		provideJaegerLogAdatper,
		provideOpentracing,
		))
}

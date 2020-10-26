package container

import (
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type ModuleContainer struct {
	HttpProviders []func(router *mux.Router)
	GrpcProviders []func(server *grpc.Server)
	CloserProviders []func()
	RunProviders []RunPair
	MigrationProvider []func() error
}

func NewModuleContainer() ModuleContainer {
	return ModuleContainer{
		HttpProviders:    []func(router *mux.Router){},
		GrpcProviders:    []func(server *grpc.Server){},
		CloserProviders:  []func(){},
		RunProviders:     []RunPair{},
		MigrationProvider: []func() error{},
	}
}

type RunPair struct {
	Loop func() error
	Exit func(error)
}

type HttpProvider interface {
	ProvideHttp(router *mux.Router)
}

type GrpcProvider interface {
	ProvideGrpc(server *grpc.Server)
}

type CloserProvider interface {
	ProvideCloser()
}

type RunProvider interface {
	ProvideRunLoop() error
	ProvideRunExit(error)
}

type MigrationProvider interface {
	ProvideMigration() error
}

type HttpFunc func(router *mux.Router)

func (h HttpFunc) ProvideHttp(router *mux.Router) {
	h(router)
}

func (s *ModuleContainer) Register(app interface{})  {
	if p, ok := app.(HttpProvider); ok {
		s.HttpProviders = append(s.HttpProviders, p.ProvideHttp)
	}
	if p, ok := app.(GrpcProvider); ok {
		s.GrpcProviders = append(s.GrpcProviders, p.ProvideGrpc)
	}
	if p, ok := app.(RunProvider); ok {
		s.RunProviders = append(s.RunProviders, RunPair{p.ProvideRunLoop, p.ProvideRunExit})
	}
	if p, ok := app.(MigrationProvider); ok {
		s.MigrationProvider = append(s.MigrationProvider, p.ProvideMigration)
	}
}

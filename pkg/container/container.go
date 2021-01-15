package container

import (
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
)

type ModuleContainer struct {
	HttpProviders     []func(router *mux.Router)
	GrpcProviders     []func(server *grpc.Server)
	CloserProviders   []func()
	RunProviders      []func(g *run.Group)
	MigrationProvider []Migrations
	SeedProvider      []func() error
	CronProviders     []func(crontab *cron.Cron)
}

func NewModuleContainer() ModuleContainer {
	return ModuleContainer{
		HttpProviders:     []func(router *mux.Router){},
		GrpcProviders:     []func(server *grpc.Server){},
		CloserProviders:   []func(){},
		RunProviders:      []func(g *run.Group){},
		MigrationProvider: []Migrations{},
		SeedProvider:      []func() error{},
	}
}

type Migrations struct {
	Migrate  func() error
	Rollback func(flag string) error
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
	ProvideRunGroup(group *run.Group)
}

type MigrationProvider interface {
	ProvideMigration() error
	ProvideRollback(flag string) error
}

type SeedProvider interface {
	ProvideSeed() error
}

type CronProvider interface {
	ProvideCron(crontab *cron.Cron)
}

type HttpFunc func(router *mux.Router)

func (h HttpFunc) ProvideHttp(router *mux.Router) {
	h(router)
}

func (s *ModuleContainer) Register(app interface{}) {
	if p, ok := app.(HttpProvider); ok {
		s.HttpProviders = append(s.HttpProviders, p.ProvideHttp)
	}
	if p, ok := app.(GrpcProvider); ok {
		s.GrpcProviders = append(s.GrpcProviders, p.ProvideGrpc)
	}
	if p, ok := app.(RunProvider); ok {
		s.RunProviders = append(s.RunProviders, p.ProvideRunGroup)
	}
	if p, ok := app.(CloserProvider); ok {
		s.CloserProviders = append(s.CloserProviders, p.ProvideCloser)
	}
	if p, ok := app.(MigrationProvider); ok {
		s.MigrationProvider = append(s.MigrationProvider, Migrations{p.ProvideMigration, p.ProvideRollback})
	}
	if p, ok := app.(SeedProvider); ok {
		s.SeedProvider = append(s.SeedProvider, p.ProvideSeed)
	}
	if p, ok := app.(CronProvider); ok {
		s.CronProviders = append(s.CronProviders, p.ProvideCron)
	}
}

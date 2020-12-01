package module

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/rule/service"
)

type Module struct {
	repository service.Repository
	endpoints  Endpoints
	close      func()
}

func New(moduleConfig contract.ConfigReader, logger log.Logger) *Module {
	module, cleanup, err := injectModule(moduleConfig, logger)
	if err != nil {
		panic(err)
	}
	module.close = cleanup
	return module
}

func (m *Module) ProvideHttp(router *mux.Router) {
	router.PathPrefix("/rule/").Handler(http.StripPrefix("/rule", MakeHTTPHandler(m.endpoints)))
}

func (m *Module) ProvideCloser() {
	m.close()
}

func (m *Module) ProvideRunGroup(group *run.Group) {
	ctx, cancel := context.WithCancel(context.Background())
	group.Add(func() error {
		return m.repository.WatchConfigUpdate(ctx)
	}, func(err error) {
		cancel()
	})
}

func (m *Module) GetRepository() service.Repository {
	return m.repository
}

package rule

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/klog"
)

type Module struct {
	repository Repository
	endpoints  Endpoints
	close      func()
}

func New(moduleConfig contract.ConfigReader, logger log.Logger) *Module {
	module, cleanup, err := injectModule(setUp(moduleConfig, logger))
	if err != nil {
		panic(err)
	}
	module.close = cleanup
	return module
}

func setUp(moduleConfig contract.ConfigReader, logger log.Logger) (contract.ConfigReader, log.Logger) {
	appLogger := log.With(logger, "module", config.ProvideAppName(moduleConfig).String())
	appLogger = level.NewFilter(logger, klog.LevelFilter(moduleConfig.String("level")))
	return moduleConfig, appLogger
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

func (m *Module) GetRepository() Repository {
	return m.repository
}

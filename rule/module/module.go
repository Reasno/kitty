package module

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	stdopentracing "github.com/opentracing/opentracing-go"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kerr"
	"glab.tagtic.cn/ad_gains/kitty/pkg/khttp"
	"glab.tagtic.cn/ad_gains/kitty/rule/service"
)

type Module struct {
	tracer     stdopentracing.Tracer
	repository service.Repository
	endpoints  Endpoints
	logger     log.Logger
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
	router.PathPrefix("/rule/").Handler(http.StripPrefix("/rule", MakeHTTPHandler(m.endpoints,
		httptransport.ServerBefore(
			opentracing.HTTPToContext(m.tracer, "app", m.logger),
			jwt.HTTPToContext(),
			khttp.IpToContext(),
		),
		httptransport.ServerErrorEncoder(kerr.ErrorEncoder))))
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

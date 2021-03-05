package module

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/sony/gobreaker"
	dmpclient "glab.tagtic.cn/ad_gains/kitty/dmp/svc/client/grpc"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kmiddleware"
	kitty "glab.tagtic.cn/ad_gains/kitty/proto"
	"glab.tagtic.cn/ad_gains/kitty/rule/repository"
	"glab.tagtic.cn/ad_gains/kitty/rule/service"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func provideEtcdClient(conf contract.ConfigReader) (*clientv3.Client, func(), error) {
	ctx, cancel := context.WithCancel(context.Background())
	client, err := clientv3.New(clientv3.Config{
		Endpoints: conf.Strings("etcd.addrs"),
		Context:   ctx,
	})
	return client, cancel, err
}

func provideRepository(client *clientv3.Client, logger log.Logger) (service.Repository, error) {
	return repository.NewRepository(client, logger)
}

func provideHistogramMetrics(appName contract.AppName, env contract.Env) metrics.Histogram {
	var his metrics.Histogram = prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: appName.String(),
		Subsystem: env.String(),
		Name:      "request_duration_seconds",
		Help:      "Total time spent serving requests.",
	}, []string{"module", "method"})
	return his
}

func provideModule(
	repository service.Repository,
	endpoints Endpoints,
	tracer stdopentracing.Tracer,
	logger log.Logger,
) *Module {
	// TODO: add middleware
	return &Module{
		repository: repository,
		endpoints:  endpoints,
		tracer:     tracer,
		logger:     logger,
	}
}

func provideDmpServer(conf contract.ConfigReader, tracer stdopentracing.Tracer, logger log.Logger, env contract.Env) (kitty.DmpServer, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)
	dmpAddr := conf.String("dmpAddr")
	if env.IsLocal() {
		conn, err = grpc.Dial(dmpAddr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		conn, err = grpc.Dial(dmpAddr, grpc.WithInsecure())
	}
	if err != nil {
		return nil, err
	}
	return NewDmpServer(DmpOption{
		Conn:   conn,
		Tracer: tracer,
		Logger: logger,
		Env:    env,
	})
}

type DmpOption struct {
	Conn   *grpc.ClientConn
	Tracer stdopentracing.Tracer
	Logger log.Logger
	Env    contract.Env
}

func NewDmpServer(opt DmpOption) (kitty.DmpServer, error) {
	endpoints, err := dmpclient.New(opt.Conn,
		grpctransport.ClientBefore(jwt.ContextToGRPC()),
		grpctransport.ClientBefore(opentracing.ContextToGRPC(opt.Tracer, opt.Logger)),
	)
	if err != nil {
		return nil, errors.Wrap(err, "fails to create dmp server")
	}

	endpoints.WrapAllExcept(kmiddleware.NewTimeoutMiddleware(500 * time.Millisecond))
	endpoints.WrapAllExcept(circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{})))
	if opt.Tracer != nil {
		endpoints.WrapAllLabeledExcept(
			kmiddleware.NewClientServerMiddleware(opt.Tracer, opt.Env.String()),
		)
	}
	return endpoints, nil
}

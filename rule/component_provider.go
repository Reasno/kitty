package rule

import (
	"context"
	"github.com/Reasno/kitty/pkg/contract"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"go.etcd.io/etcd/clientv3"
)

func provideEtcdClient(conf contract.ConfigReader) (*clientv3.Client, func(), error) {
	ctx, cancel := context.WithCancel(context.Background())
	client, err := clientv3.New(clientv3.Config{
		Endpoints:            conf.Strings("etcd.addrs"),
		Context:              ctx,
	})
	return client, cancel, err
}

func provideRepository(client *clientv3.Client, logger log.Logger, conf contract.ConfigReader) (*repository, error) {
	var ruleMapping = make(map[string]string)
	for k, v := range conf.Get("rules").(map[string]interface{}) {
		ruleMapping[k] = v.(string)
	}
	return NewRepository(client, logger, ruleMapping)
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

func provideModule(repository Repository, endpoints Endpoints) *Module {
	// TODO: add middleware
	return &Module{
		repository: repository,
		endpoints: endpoints,
	}
}

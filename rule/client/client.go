package client

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/pkg/errors"
	kconf "glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/rule"
	"go.etcd.io/etcd/clientv3"
)

type DynamicConfig struct {
	repository rule.Repository
	logger     log.Logger
}

func (d *DynamicConfig) Of(ruleName string, payload *rule.Payload) (*kconf.KoanfAdapter, error) {
	compiled := d.repository.GetCompiled(ruleName)
	calculated, err := rule.Calculate(compiled, payload, d.logger)
	if err != nil {
		return nil, err
	}
	c := koanf.New(".")
	err = c.Load(confmap.Provider(calculated, ""), yaml.Parser())
	if err != nil {
		return nil, errors.Wrap(err, "cannot load from map")
	}
	adapter := kconf.NewKoanfAdapter(c)
	return adapter, nil
}

type Option func(*config)

type config struct {
	ctx    context.Context
	client *clientv3.Client
	logger log.Logger
}

func WithClient(client *clientv3.Client) Option {
	return func(c *config) {
		c.client = client
	}
}

func WithLogger(logger log.Logger) Option {
	return func(c *config) {
		c.logger = logger
	}
}

func WithContext(ctx context.Context) Option {
	return func(c *config) {
		c.ctx = ctx
	}
}

func NewDynamicConfig(ctx context.Context, opt ...Option) (*DynamicConfig, error) {
	c := config{}
	for _, o := range opt {
		o(&c)
	}
	if c.client == nil {
		client, err := clientv3.New(clientv3.Config{
			Endpoints: []string{"etcd-1:2379", "etcd-2:2379", "etcd-3:2379"},
			Context:   ctx,
		})
		if err != nil {
			return nil, errors.Wrap(err, "Failed to connect to ETCD")
		}
		c.client = client
	}
	if c.logger == nil {
		c.logger = log.NewNopLogger()
	}
	repository, err := rule.NewRepository(c.client, c.logger)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create repository")
	}
	go func() {
		level.Error(c.logger).Log("err", repository.WatchConfigUpdate(ctx))
	}()
	return &DynamicConfig{repository: repository, logger: c.logger}, nil
}

package client

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/rule"
	"go.etcd.io/etcd/clientv3"
)

type DynamicConfig struct {
	repository rule.Repository
}

func (d DynamicConfig) String(s string) string {
	panic("implement me")
}

func (d DynamicConfig) Int(s string) int {
	panic("implement me")
}

func (d DynamicConfig) Strings(s string) []string {
	panic("implement me")
}

func (d DynamicConfig) Bool(s string) bool {
	panic("implement me")
}

func (d DynamicConfig) Get(s string) interface{} {
	panic("implement me")
}

func (d DynamicConfig) Float64(s string) float64 {
	panic("implement me")
}

func (d DynamicConfig) Cut(s string) contract.ConfigReader {
	panic("implement me")
}

type Option func(*config)

type config struct {
	ctx context.Context
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
	return &DynamicConfig{repository: repository}, nil
}

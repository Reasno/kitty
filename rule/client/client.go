package client

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/pkg/errors"
	kconf "glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/rule"
	"go.etcd.io/etcd/clientv3"
)

type DynamicConfig struct {
	repository rule.Repository
	logger     log.Logger
}

type ofRule struct {
	d        *DynamicConfig
	ruleName string
}

func (r *ofRule) Tenant(tenant *kconf.Tenant) (contract.ConfigReader, error) {
	var pl = rule.FromTenant(tenant)
	compiled := r.d.repository.GetCompiled(r.ruleName)

	calculated, err := rule.Calculate(compiled, pl, r.d.logger)
	if err != nil {
		return nil, err
	}

	c := koanf.New(".")
	err = c.Load(confmap.Provider(calculated, ""), nil)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load from map")
	}

	adapter := kconf.NewKoanfAdapter(c)
	return adapter, nil
}

func (d *DynamicConfig) WithRule(ruleName string) *ofRule {
	return &ofRule{
		ruleName: ruleName,
		d:        d,
	}
}

func (d *DynamicConfig) Watch(ctx context.Context) error {
	return d.repository.WatchConfigUpdate(ctx)
}

type Option func(*config)

type config struct {
	ctx         context.Context
	client      *clientv3.Client
	repo        rule.Repository
	logger      log.Logger
	listOfRules []string
}

func WithClient(client *clientv3.Client) Option {
	return func(c *config) {
		c.client = client
	}
}

func WithRepository(repository rule.Repository) Option {
	return func(c *config) {
		c.repo = repository
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

func WithListOfRules(listOfRules []string) Option {
	return func(c *config) {
		c.listOfRules = listOfRules
	}
}

func NewDynamicConfig(ctx context.Context, opt ...Option) (*DynamicConfig, error) {
	c := config{}
	for _, o := range opt {
		o(&c)
	}

	if c.logger == nil {
		c.logger = log.NewNopLogger()
	}
	if c.repo == nil {
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
		var mp = make(map[string]string)
		for _, v := range c.listOfRules {
			mp[v] = rule.OtherConfigPathPrefix + "/" + v
		}
		var err error
		c.repo, err = NewRepository(c.client, c.logger, mp)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create repository")
		}
	}

	return &DynamicConfig{repository: c.repo, logger: c.logger}, nil
}

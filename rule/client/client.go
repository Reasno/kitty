package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	kconf "glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
	"glab.tagtic.cn/ad_gains/kitty/rule/dto"
	"glab.tagtic.cn/ad_gains/kitty/rule/entity"
	"glab.tagtic.cn/ad_gains/kitty/rule/module"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type RuleEngine struct {
	env        contract.Env
	tracer     opentracing.Tracer
	dmpConn    *grpc.ClientConn
	repository Repository
	logger     log.Logger
}

type ofRule struct {
	d        *RuleEngine
	ruleName string
}

type Repository interface {
	GetCompiled(ruleName string) entity.Ruler
	WatchConfigUpdate(ctx context.Context) error
}

func (r *ofRule) Tenant(tenant *kconf.Tenant) (contract.ConfigReader, error) {
	var payload = dto.FromTenant(tenant)
	return r.Payload(payload)
}

func (r *ofRule) Payload(pl *dto.Payload) (contract.ConfigReader, error) {
	var compiled entity.Ruler

	parts := strings.Split(pl.PackageName, ".")
	if len(parts) > 0 {
		codeName := parts[len(parts)-1]
		compiled = r.d.repository.GetCompiled(codeName + "-" + r.ruleName)
	}
	// 兼容之前的情况，去商业化平台中心查配置
	if compiled == nil {
		compiled = r.d.repository.GetCompiled(r.ruleName)
	}
	if compiled == nil {
		return nil, fmt.Errorf("no suitable configuration found for %s", r.ruleName)
	}

	if compiled.ShouldEnrich() {
		endpoints, err := module.NewDmpServer(module.DmpOption{
			Conn:   r.d.dmpConn,
			Tracer: r.d.tracer,
			Logger: r.d.logger,
			Env:    r.d.env,
		})
		if err != nil {
			return nil, errors.Wrap(err, "unable to create dmp server")
		}
		if pl.Context == nil {
			pl.Context = context.Background()
		}
		resp, err := endpoints.UserMore(pl.Context, &pb.DmpReq{
			UserId:      pl.UserId,
			PackageName: pl.PackageName,
			Suuid:       pl.Suuid,
			Channel:     pl.Channel,
		})
		if err != nil {
			level.Warn(r.d.logger).Log("err", err)
		}
		if resp == nil {
			resp = &pb.DmpResp{}
		}
		pl.DMP = *resp
	}

	calculated, err := entity.Calculate(compiled, pl)
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

func (d *RuleEngine) Of(ruleName string) *ofRule {
	return &ofRule{
		ruleName: ruleName,
		d:        d,
	}
}

func (d *RuleEngine) Watch(ctx context.Context) error {
	return d.repository.WatchConfigUpdate(ctx)
}

type Option func(*config)

type config struct {
	ctx         context.Context
	env         contract.Env
	tracer      opentracing.Tracer
	dmpAddr     string
	client      *clientv3.Client
	repo        Repository
	logger      log.Logger
	listOfRules []string
	rulePrefix  string
	ruleRegexp  *regexp.Regexp
}

func WithClient(client *clientv3.Client) Option {
	return func(c *config) {
		c.client = client
	}
}

func WithRepository(repository Repository) Option {
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

func WithRulePrefix(prefix string) Option {
	return func(c *config) {
		c.rulePrefix = prefix
	}
}

func WithRuleRegexp(regexp *regexp.Regexp) Option {
	return func(c *config) {
		c.ruleRegexp = regexp
	}
}

func WithDMPAddr(dmpAddr string) Option {
	return func(c *config) {
		c.dmpAddr = dmpAddr
	}
}

func WithEnv(env contract.Env) Option {
	return func(c *config) {
		c.env = env
	}
}

func Rule(rule string) Option {
	return func(c *config) {
		c.listOfRules = append(c.listOfRules, rule)
	}
}

func NewRuleEngine(opt ...Option) (*RuleEngine, error) {
	c := config{
		ctx:         context.Background(),
		logger:      log.NewNopLogger(),
		listOfRules: make([]string, 0),
		env:         kconf.Env("production"),
		dmpAddr:     "xtasks.ad:8181",
	}
	for _, o := range opt {
		o(&c)
	}

	if c.repo == nil {
		if c.client == nil {
			client, err := clientv3.New(clientv3.Config{
				Endpoints: []string{"etcd-1:2379", "etcd-2:2379", "etcd-3:2379"},
				Context:   c.ctx,
			})
			if err != nil {
				return nil, errors.Wrap(err, "Failed to connect to ETCD")
			}
			c.client = client
		}

		var err error
		c.repo, err = NewRepositoryWithConfig(c.client, c.logger, RepositoryConfig{
			Prefix:      c.rulePrefix,
			Regex:       c.ruleRegexp,
			ListOfRules: c.listOfRules,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create repository")
		}
	}
	var (
		err  error
		conn *grpc.ClientConn
	)
	if c.env.IsLocal() {
		conn, err = grpc.Dial(c.dmpAddr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		conn, err = grpc.Dial(c.dmpAddr, grpc.WithInsecure())
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial dmp server")
	}

	return &RuleEngine{
		repository: c.repo,
		logger:     c.logger,
		tracer:     c.tracer,
		env:        c.env,
		dmpConn:    conn,
	}, nil
}

package rule

import (
	"bytes"
	"context"
	"github.com/Reasno/kitty/pkg/kerr"
	"github.com/Reasno/kitty/rule/msg"
	"github.com/antonmedv/expr"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"io"
)

var ErrDataHasChanged = errors.New(msg.ErrorRulesHasChanged)

type Service interface {
	CalculateRules(ctx context.Context, ruleName string, payload *Payload) (Data, error)
	GetRules(ctx context.Context, ruleName string) ([]byte, error)
	UpdateRules(ctx context.Context, ruleName string, content []byte, dryRun bool) error
	Preflight(ctx context.Context, ruleName string, hash string) error
}

type Repository interface {
	GetCompiled(ruleName string) []Rule
	GetRaw(ctx context.Context, key string) (value []byte, e error)
	SetRaw(ctx context.Context, key string, value string) error
	IsNewest(ctx context.Context, key, value string) (bool, error)
	WatchConfigUpdate(ctx context.Context) error
}

type service struct {
	logger log.Logger
	repo   Repository
}

func NewService(logger log.Logger, repo Repository) *service {
	return &service{logger: logger, repo: repo}
}

func (r *service) CalculateRules(ctx context.Context, ruleName string, payload *Payload) (Data, error) {
	for _, rule := range r.repo.GetCompiled(ruleName) {
		output, err := expr.Run(rule.program, payload)
		if err != nil {
			return nil, errors.Wrap(err, msg.ErrorRules)
		}
		if !output.(bool) {
			level.Debug(r.logger).Log("msg", "negative: "+rule.If)
			continue
		}
		level.Debug(r.logger).Log("msg", "positive: "+rule.If)
		return rule.Then, nil
	}
	return Data{}, nil
}

func (r *service) GetRules(ctx context.Context, ruleName string) ([]byte, error) {
	return r.repo.GetRaw(ctx, ruleName)
}

func (r *service) UpdateRules(ctx context.Context, ruleName string, content []byte, dryRun bool) error {
	var (
		value string
		buf   bytes.Buffer
		err   error
		tee   io.Reader
	)
	reader := bytes.NewReader(content)
	tee = io.TeeReader(reader, &buf)
	err = validateRules(tee)
	var invalid ErrInvalidRules
	if errors.As(err, &invalid) {
		return kerr.InvalidArgumentErr(invalid)
	}
	if err != nil {
		return err
	}
	if dryRun {
		return nil
	}
	value = buf.String()
	err = r.repo.SetRaw(ctx, ruleName, value)
	return err
}

func (r *service) Preflight(ctx context.Context, ruleName string, hash string) error {
	ok, err := r.repo.IsNewest(ctx, ruleName, hash)
	if err != nil {
		return errors.Wrap(err, msg.ErrorETCD)
	}
	if !ok {
		return ErrDataHasChanged
	}
	return nil
}

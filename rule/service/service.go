//go:generate mockery --name=Repository

package service

import (
	"bytes"
	"context"
	"io"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kerr"
	"glab.tagtic.cn/ad_gains/kitty/rule/dto"
	"glab.tagtic.cn/ad_gains/kitty/rule/entity"
	"glab.tagtic.cn/ad_gains/kitty/rule/msg"
)

var ErrDataHasChanged = errors.New(msg.ErrorRulesHasChanged)

type Service interface {
	CalculateRules(ctx context.Context, ruleName string, payload *dto.Payload) (dto.Data, error)
	GetRules(ctx context.Context, ruleName string) ([]byte, error)
	UpdateRules(ctx context.Context, ruleName string, content []byte, dryRun bool) error
	Preflight(ctx context.Context, ruleName string, hash string) error
}

type Repository interface {
	GetCompiled(ruleName string) []entity.Rule
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

func (r *service) CalculateRules(ctx context.Context, ruleName string, payload *dto.Payload) (dto.Data, error) {
	rules := r.repo.GetCompiled(ruleName)
	return entity.Calculate(rules, payload, r.logger)
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
	err = entity.ValidateRules(tee)
	var invalid *entity.ErrInvalidRules
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

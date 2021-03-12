//go:generate mockery --name=Repository
//go:generate mockery --name=DmpServer

package service

import (
	"bytes"
	"context"
	"io"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kerr"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
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
	GetCompiled(ruleName string) entity.Ruler
	GetRaw(ctx context.Context, key string) (value []byte, e error)
	SetRaw(ctx context.Context, key string, value string) error
	IsNewest(ctx context.Context, key, value string) (bool, error)
	WatchConfigUpdate(ctx context.Context) error
	ValidateRules(ruleName string, reader io.Reader) error
}

type service struct {
	dmpServerProd pb.DmpServer
	dmpServerDev  pb.DmpServer
	logger        log.Logger
	repo          Repository
}

func NewService(logger log.Logger, repo Repository) *service {
	return &service{logger: logger, repo: repo}
}

func (r *service) CalculateRules(ctx context.Context, ruleName string, payload *dto.Payload) (dto.Data, error) {
	rules := r.repo.GetCompiled(ruleName)
	if rules == nil {
		return nil, errors.New("rule not found")
	}
	if rules.ShouldEnrich() {
		resp, err := r.dmp(ruleName).UserMore(ctx, &pb.DmpReq{
			UserId:      payload.UserId,
			PackageName: payload.PackageName,
			Suuid:       payload.Suuid,
			Channel:     payload.Channel,
		})
		if err != nil {
			level.Warn(r.logger).Log("err", errors.Wrap(err, "dmp server error"))
		}
		if resp == nil {
			resp = &pb.DmpResp{}
		}
		payload.DMP = *resp
	}
	return entity.Calculate(rules, payload)
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
	err = r.repo.ValidateRules(ruleName, tee)
	var invalid *entity.ErrInvalidRules
	if errors.As(err, &invalid) {
		return kerr.InvalidArgumentErr(invalid, msg.ErrorRules)
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

func (r *service) dmp(rule string) pb.DmpServer {
	if len(rule) > 5 && rule[len(rule)-6:] == "-prod" {
		return r.dmpServerProd
	}
	return r.dmpServerDev
}

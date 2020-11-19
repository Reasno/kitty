package internal

import (
	"context"

	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

type InvitationManagerFacade struct {
	Name    contract.AppName
	Factory InvitationManagerFactory
	DynConf config.DynamicConfigReader
}

func (im *InvitationManagerFacade) AddToken(ctx context.Context, userId uint64, token string) error {
	m, err := im.getManager(ctx)
	if err != nil {
		return err
	}
	return m.AddToken(ctx, userId, token)
}

func (im *InvitationManagerFacade) ClaimReward(ctx context.Context, masterId uint64, apprenticeId uint64) error {
	m, err := im.getManager(ctx)
	if err != nil {
		return err
	}
	return m.ClaimReward(ctx, masterId, apprenticeId)
}

func (im *InvitationManagerFacade) AdvanceStep(ctx context.Context, apprenticeId uint64, eventName string) error {
	m, err := im.getManager(ctx)
	if err != nil {
		return err
	}
	return m.AdvanceStep(ctx, apprenticeId, eventName)
}

func (im *InvitationManagerFacade) ListApprentices(ctx context.Context, masterId uint64, depth int) ([]RelationWithRewardAmount, error) {
	m, err := im.getManager(ctx)
	if err != nil {
		return nil, err
	}
	return m.ListApprentices(ctx, masterId, depth)
}

func (im *InvitationManagerFacade) GetToken(ctx context.Context) string {
	m, _ := im.getManager(ctx)
	return m.GetToken(ctx)
}

func (im *InvitationManagerFacade) GetUrl(ctx context.Context) string {
	m, _ := im.getManager(ctx)
	return m.GetUrl(ctx)
}

func (im *InvitationManagerFacade) getManager(ctx context.Context) (*InvitationManager, error) {
	tenant := config.GetTenant(ctx)
	conf, err := im.DynConf.Tenant(tenant)
	if err != nil {
		return nil, errors.Wrap(err, "no configuration found for invitation tenant")
	}
	return im.Factory.NewManager(conf.Cut("share")), nil
}

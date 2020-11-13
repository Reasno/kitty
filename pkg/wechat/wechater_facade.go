package wechat

import (
	"context"

	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
)

type WechaterFacade struct {
	factory *WechaterFactory
	dynConf config.DynamicConfigReader
}

func (w *WechaterFacade) GetLoginResponse(ctx context.Context, code string) (result *WxLoginResult, err error) {
	wechater, err := w.getRealWechater(ctx)
	if err != nil {
		return nil, err
	}
	return wechater.GetLoginResponse(ctx, code)
}

func (w *WechaterFacade) GetUserInfoResult(ctx context.Context, wxLoginResult *WxLoginResult) (*WxUserInfoResult, error) {
	wechater, err := w.getRealWechater(ctx)
	if err != nil {
		return nil, err
	}
	return wechater.GetUserInfoResult(ctx, wxLoginResult)
}

func NewWechaterFacade(factory *WechaterFactory, reader config.DynamicConfigReader) *WechaterFacade {
	return &WechaterFacade{factory: factory, dynConf: reader}
}

func (w *WechaterFacade) getRealWechater(ctx context.Context) (Wechater, error) {
	tenant := config.GetTenant(ctx)
	conf, err := w.dynConf.Tenant(tenant)
	if err != nil {
		return nil, errors.Wrap(err, "no configuration found for sms tenant")
	}
	return w.factory.GetTransportByConf(conf), nil
}

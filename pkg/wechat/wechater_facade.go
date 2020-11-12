package wechat

import (
	"context"

	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
)

type WechaterFacade struct {
	factory *WechaterFactory
	dynConf config.DynamicConfigReader
}

func (w *WechaterFacade) GetLoginResponse(ctx context.Context, code string) (result *WxLoginResult, err error) {
	wechater := w.getRealWechater(ctx)
	return wechater.GetLoginResponse(ctx, code)
}

func (w *WechaterFacade) GetUserInfoResult(ctx context.Context, wxLoginResult *WxLoginResult) (*WxUserInfoResult, error) {
	wechater := w.getRealWechater(ctx)
	return wechater.GetUserInfoResult(ctx, wxLoginResult)
}

func NewWechaterFacade(factory *WechaterFactory, reader config.DynamicConfigReader) *WechaterFacade {
	return &WechaterFacade{factory: factory, dynConf: reader}
}

func (w *WechaterFacade) getRealWechater(ctx context.Context) Wechater {
	packageName := config.GetTenant(ctx).PackageName
	return w.factory.GetTransport(packageName)
}

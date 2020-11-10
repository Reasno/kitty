package wechat

import (
	"context"
)

type Wechater interface {
	GetLoginResponse(ctx context.Context, code string) (result *WxLoginResult, err error)
	GetUserInfoResult(ctx context.Context, wxLoginResult *WxLoginResult) (*WxUserInfoResult, error)
}

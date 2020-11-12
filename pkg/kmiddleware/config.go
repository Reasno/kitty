package kmiddleware

import (
	"context"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	kittyjwt "glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
)

func NewConfigMiddleware() endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if claim, ok := ctx.Value(jwt.JWTClaimsContextKey).(*kittyjwt.Claim); ok {
				t := NewTenantFromClaim(claim)
				t.Ip, _ = ctx.Value(contract.IpKey).(string)
				ctx = context.WithValue(ctx, config.TenantKey, t)
			}
			return e(ctx, request)
		}
	}
}

func NewTenantFromClaim(claim *kittyjwt.Claim) *config.Tenant {
	return &config.Tenant{
		Channel:     claim.Channel,
		VersionCode: claim.VersionCode,
		Suuid:       claim.Suuid,
		UserId:      claim.UserId,
		PackageName: claim.PackageName,
	}
}

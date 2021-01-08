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
			var tenant config.Tenant

			if claim, ok := ctx.Value(jwt.JWTClaimsContextKey).(*kittyjwt.Claim); ok {
				tenant = NewTenantFromClaim(claim)
			}
			if ip, ok := ctx.Value(contract.IpKey).(string); ok {
				tenant.Ip = ip
			}
			if t, ok := request.(interface {
				GetChannel() string
			}); ok {
				tenant.Channel = t.GetChannel()
			}
			if t, ok := request.(interface {
				GetUserId() uint64
			}); ok {
				tenant.UserId = t.GetUserId()
			}
			if t, ok := request.(interface {
				GetSuuid() string
			}); ok {
				tenant.Suuid = t.GetSuuid()
			}
			if t, ok := request.(interface {
				GetPackageName() string
			}); ok {
				tenant.PackageName = t.GetPackageName()
			}
			if t, ok := request.(interface {
				GetVersionCode() string
			}); ok {
				tenant.VersionCode = t.GetVersionCode()
			}
			ctx = context.WithValue(ctx, config.TenantKey, &tenant)
			return e(ctx, request)
		}
	}
}

func NewTenantFromClaim(claim *kittyjwt.Claim) config.Tenant {
	return config.Tenant{
		Channel:     claim.Channel,
		VersionCode: claim.VersionCode,
		Suuid:       claim.Suuid,
		UserId:      claim.UserId,
		PackageName: claim.PackageName,
	}
}

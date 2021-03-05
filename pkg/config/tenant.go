package config

import (
	"context"

	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

var TenantKey = struct{}{}

type Tenant struct {
	Channel     string          `json:"channel"`
	VersionCode string          `json:"version_code"`
	Os          uint8           `json:"os"`
	UserId      uint64          `json:"user_id"`
	Imei        string          `json:"imei"`
	Idfa        string          `json:"idfa"`
	Oaid        string          `json:"oaid"`
	Suuid       string          `json:"suuid"`
	Mac         string          `json:"mac"`
	AndroidId   string          `json:"android_id"`
	PackageName string          `json:"package_name"`
	Ip          string          `json:"ip"`
	Context     context.Context `json:"-"`
}

func GetTenant(ctx context.Context) *Tenant {
	if c, ok := ctx.Value(TenantKey).(*Tenant); ok {
		c.Context = ctx
		return c
	}
	return &Tenant{Context: context.Background()}
}

type DynamicConfigReader interface {
	Tenant(tenant *Tenant) (contract.ConfigReader, error)
}

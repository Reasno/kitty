package rule

import (
	"encoding/json"

	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	jwt2 "glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
)

type Payload struct {
	Channel     string `json:"channel"`
	VersionCode string `json:"version_code"`
	Os          uint8  `json:"os"`
	UserId      uint64 `json:"user_id"`
	Imei        string `json:"imei"`
	Idfa        string `json:"idfa"`
	Oaid        string `json:"oaid"`
	Suuid       string `json:"suuid"`
	Mac         string `json:"mac"`
	AndroidId   string `json:"android_id"`
	PackageName string `json:"package_name"`
	Ip          string `json:"ip"`
}

func FromClaim(claim jwt2.Claim) *Payload {
	return &Payload{
		Channel:     claim.Channel,
		VersionCode: claim.VersionCode,
		Suuid:       claim.Suuid,
		UserId:      claim.UserId,
	}
}

func FromTenant(tenant *config.Tenant) *Payload {
	return &Payload{
		Channel:     tenant.Channel,
		VersionCode: tenant.VersionCode,
		Os:          tenant.Os,
		UserId:      tenant.UserId,
		Imei:        tenant.Imei,
		Idfa:        tenant.Idfa,
		Oaid:        tenant.Oaid,
		Suuid:       tenant.Suuid,
		Mac:         tenant.Mac,
		AndroidId:   tenant.AndroidId,
		PackageName: tenant.PackageName,
		Ip:          tenant.Ip,
	}
}

func (p *Payload) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

type Data map[string]interface{}

type Response struct {
	Code    uint `json:"code"`
	Message uint `json:"message"`
	Data    Data `json:"data"`
}

func (p Response) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

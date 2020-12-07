package dto

import (
	"encoding/json"
	"strconv"

	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	jwt2 "glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
)

type Payload struct {
	Channel     string                 `json:"channel" schema:"channel"`
	VersionCode int                    `json:"version_code" schema:"version_code"`
	Os          uint8                  `json:"os" schema:"os"`
	UserId      uint64                 `json:"user_id" schema:"user_id"`
	Imei        string                 `json:"imei" schema:"imei"`
	Idfa        string                 `json:"idfa" schema:"idfa"`
	Oaid        string                 `json:"oaid" schema:"oaid"`
	Suuid       string                 `json:"suuid" schema:"suuid"`
	Mac         string                 `json:"mac" schema:"mac"`
	AndroidId   string                 `json:"android_id" schema:"android_id"`
	PackageName string                 `json:"package_name" schema:"package_name"`
	Ip          string                 `json:"ip" schema:"ip"`
	Q           map[string][]string    `json:"-" schema:"-"`
	B           map[string]interface{} `json:"-" schema:"-"`
}

func FromClaim(claim jwt2.Claim) *Payload {
	versionCode, _ := strconv.Atoi(claim.VersionCode)
	return &Payload{
		Channel:     claim.Channel,
		VersionCode: versionCode,
		Suuid:       claim.Suuid,
		UserId:      claim.UserId,
	}
}

func FromTenant(tenant *config.Tenant) *Payload {
	versionCode, _ := strconv.Atoi(tenant.VersionCode)
	return &Payload{
		Channel:     tenant.Channel,
		VersionCode: versionCode,
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

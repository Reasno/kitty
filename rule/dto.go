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

func (p *Payload) FromClaim(claim jwt2.Claim) {
	p.Channel = claim.Channel
	p.VersionCode = claim.VersionCode
	p.Suuid = claim.Suuid
	p.UserId = claim.UserId
}

func (p *Payload) FromTenant(tenant *config.Tenant) {
	p.Channel = tenant.Channel
	p.Os = tenant.Os
	p.UserId = tenant.UserId
	p.Oaid = tenant.Oaid
	p.Idfa = tenant.Idfa
	p.Mac = tenant.Mac
	p.AndroidId = tenant.AndroidId
	p.Mac = tenant.Mac
	p.Channel = tenant.Channel
	p.VersionCode = tenant.VersionCode
	p.Suuid = tenant.Suuid
	p.Imei = tenant.Imei
	p.Ip = tenant.Ip
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

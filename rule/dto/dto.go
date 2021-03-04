package dto

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	jwt2 "glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
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
	DMP         pb.DmpResp             `json:"-" schema:"-"`
	Context     context.Context        `json:"-" schema:"-"`
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
		Context:     tenant.Context,
	}
}

func (p *Payload) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func (p Payload) Now() time.Time {
	return time.Now()
}

func (p Payload) Date(s string) time.Time {
	date, err := time.ParseInLocation("2006-01-02", s, time.Local)
	if err != nil {
		panic(err)
	}
	return date
}

func (p Payload) DaysAgo(s string) int {
	return int(time.Now().Sub(p.DateTime(s)).Hours() / 24)
}

func (p Payload) HoursAgo(s string) int {
	return int(time.Now().Sub(p.DateTime(s)).Hours())
}

func (p Payload) DateTime(s string) time.Time {
	date, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	if err != nil {
		panic(err)
	}
	return date
}

func (p Payload) IsBefore(s string) bool {
	var (
		t   time.Time
		err error
	)
	if len(s) == 10 {
		t, err = time.ParseInLocation("2006-01-02", s, time.Local)
	} else {
		t, err = time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	}
	if err != nil {
		panic(err)
	}
	return time.Now().Before(t)
}

func (p Payload) IsAfter(s string) bool {
	var (
		t   time.Time
		err error
	)
	if len(s) == 10 {
		t, err = time.ParseInLocation("2006-01-02", s, time.Local)
	} else {
		t, err = time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	}
	if err != nil {
		panic(err)
	}
	return time.Now().After(t)
}

func (p Payload) IsBetween(begin string, end string) bool {
	return p.IsAfter(begin) && p.IsBefore(end)
}

func (p Payload) IsWeekday(day int) bool {
	return time.Now().Weekday() == time.Weekday(day)
}

func (p Payload) IsWeekend() bool {
	if weekday := time.Now().Weekday(); weekday == 0 || weekday == 6 {
		return true
	}
	return false
}

func (p Payload) IsToday(s string) bool {
	return time.Now().Format("2006-01-02") == s
}

func (p Payload) IsHourRange(begin int, end int) bool {
	now := time.Now().Hour()
	return now >= begin && now <= end
}

func (p Payload) IsBlackListed() bool {
	return p.DMP.BlackType == pb.DmpResp_BLACK
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

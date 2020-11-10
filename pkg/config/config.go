package config

import (
	"context"

	"github.com/knadh/koanf"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

type Env string
type AppName string

var TenantKey = struct{}{}

func (a AppName) String() string {
	return string(a)
}

func (e Env) IsLocal() bool {
	return e == "local"
}

func (e Env) IsTesting() bool {
	return e == "testing"
}

func (e Env) IsDev() bool {
	return e == "dev"
}

func (e Env) IsProd() bool {
	return e == "prod"
}

func (e Env) String() string {
	return string(e)
}

func ProvideEnv(conf contract.ConfigReader) Env {
	return Env(conf.String("env"))
}

func ProvideAppName(conf contract.ConfigReader) AppName {
	return AppName(conf.String("name"))
}

type KoanfAdapter struct {
	k *koanf.Koanf
}

func (k *KoanfAdapter) Cut(s string) contract.ConfigReader {
	cut := k.k.Cut("global")
	cut.Merge(k.k.Cut(s))
	return NewKoanfAdapter(cut)
}

func NewKoanfAdapter(k *koanf.Koanf) *KoanfAdapter {
	return &KoanfAdapter{k}
}

func (k *KoanfAdapter) String(s string) string {
	return k.k.String(s)
}

func (k *KoanfAdapter) Int(s string) int {
	return k.k.Int(s)
}

func (k *KoanfAdapter) Strings(s string) []string {
	return k.k.Strings(s)
}

func (k *KoanfAdapter) Bool(s string) bool {
	return k.k.Bool(s)
}

func (k *KoanfAdapter) Get(s string) interface{} {
	return k.k.Get(s)
}

func (k *KoanfAdapter) Float64(s string) float64 {
	return k.k.Float64(s)
}

type Tenant struct {
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

func GetTenant(ctx context.Context) *Tenant {
	if c, ok := ctx.Value(TenantKey).(*Tenant); ok {
		return c
	}
	return &Tenant{}
}

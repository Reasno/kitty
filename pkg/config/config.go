package config

import (
	"github.com/knadh/koanf"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

//func ProvideChildConfig(root *viper.Viper, nodes ...string) (*viper.Viper, error) {
//	var err error
//	v := viper.New()
//	for _, n := range nodes {
//		vSettings, ok := root.AllSettings()[n].(map[string]interface{})
//		if !ok {
//			return nil, fmt.Errorf("%s settings not found", n)
//		}
//		err = v.MergeConfigMap(vSettings)
//		if err != nil {
//			return nil, fmt.Errorf("config not merged: %w", err)
//		}
//	}
//	return v, nil
//}

type Env string
type AppName string

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
	cut := k.k.Cut(s)
	cut.Merge(k.k.Cut("global"))
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

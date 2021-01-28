package module

import (
	"github.com/go-kit/kit/log"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	kitty_log "glab.tagtic.cn/ad_gains/kitty/pkg/klog"
)

func ProvideConfig(cfgFile string) (contract.ConfigReader, error) {
	k := koanf.New(".")
	if cfgFile == "" {
		cfgFile = "./config/kitty.yaml"
	}

	err := k.Load(file.Provider("./config/kitty.yaml"), yaml.Parser())
	if err != nil {
		_ = k.Load(rawbytes.Provider([]byte(`
global:
  version: 0.2.0
  env: local
  http:
    addr: :8080
  grpc:
    addr: :9090
  security:
    enable: true
    kid: kitty
    key: zxcvb0997zSDvHSD
  level: debug
app:
  name: kitty
  redis:
    addrs:
      - 127.0.0.1:6379
    database: 0
  gorm:
    database: mysql
    dsn: root@tcp(127.0.0.1:3306)/kitty?charset=utf8mb4&parseTime=True&loc=Local
  jaeger:
    sampler:
      type: 'const'
      param: 1
    log:
      enable: true
  sms:
    sendUrl: "http://hy.mix2.zthysms.com/sendSms.do"
    balanceUrl: "http://hy.mix2.zthysms.com/balance.do"
    username: ""
    password: ""
    tag: ""
  wechat:
    wechatAccessTokenUrl: https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code
    wechatGetUserInfoUrl: https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s
    appId:
    appSecret:
  s3:
    accessKey:
    accessSecret:
    bucket: ad-material
    endpoint: http://minio.xg.tagtic.cn
    region: cn-foshan-1
    cdnUrl: http://ad-static-xg.tagtic.cn/ad-material/%s

`)), yaml.Parser())
	}

	return config.NewKoanfAdapter(k), nil
}

func ProvideLogger(conf contract.ConfigReader) log.Logger {
	logger := kitty_log.NewLogger(config.Env(conf.String("global.env")))
	return logger
}

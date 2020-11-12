package sms

import (
	"strings"
	"sync"

	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

type SenderFactory struct {
	mutex sync.Mutex
	cache map[string]contract.SmsSender
	conf  contract.ConfigReader
	doer  contract.HttpDoer
}

func NewTransportFactory(conf contract.ConfigReader, doer contract.HttpDoer) *SenderFactory {
	return &SenderFactory{conf: conf, doer: doer, cache: make(map[string]contract.SmsSender)}
}

func (t *SenderFactory) getSmsConfig(conf contract.ConfigReader, name string) *TransportConfig {
	conf = conf.Cut(name)
	return &TransportConfig{
		Tag:        conf.String("sms.tag"),
		SendUrl:    conf.String("sms.sendUrl"),
		BalanceUrl: conf.String("sms.balanceUrl"),
		UserName:   conf.String("sms.username"),
		Password:   conf.String("sms.password"),
		Client:     t.doer,
	}
}

func (t *SenderFactory) GetTransport(name string) contract.SmsSender {
	return t.GetTransportWithConf(name, t.conf)
}

func (t *SenderFactory) GetTransportWithConf(name string, conf contract.ConfigReader) contract.SmsSender {
	name = strings.ReplaceAll(name, ".", "_")
	if name == "" {
		name = "default"
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if tsp, ok := t.cache[name]; ok {
		return tsp
	}
	// Currently we only have one kind of sender.
	t.cache[name] = NewTransport(t.getSmsConfig(conf, name))
	return t.cache[name]
}

func (t *SenderFactory) GetTransportByConf(conf contract.ConfigReader) contract.SmsSender {
	// Currently we only have one kind of sender.
	return NewTransport(&TransportConfig{
		Tag:        conf.String("sms.tag"),
		SendUrl:    conf.String("sms.sendUrl"),
		BalanceUrl: conf.String("sms.balanceUrl"),
		UserName:   conf.String("sms.username"),
		Password:   conf.String("sms.password"),
		Client:     t.doer,
	})
}

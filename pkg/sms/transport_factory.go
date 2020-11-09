package sms

import (
	"strings"
	"sync"

	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

type TransportFactory struct {
	mutex sync.Mutex
	cache map[string]*Transport
	conf  contract.ConfigReader
	doer contract.HttpDoer
}

func NewTransportFactory(conf contract.ConfigReader, doer contract.HttpDoer) *TransportFactory {
	return &TransportFactory{conf: conf, doer: doer, cache: make(map[string]*Transport)}
}

func (t *TransportFactory) getSmsConfig(conf contract.ConfigReader, name string) *TransportConfig {
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

func (t *TransportFactory) GetTransport(name string) *Transport {
	name = strings.ReplaceAll(name, ".", "_")
	if name == "" {
		name = "default"
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if tsp, ok := t.cache[name]; ok {
		return tsp
	}
	t.cache[name] = NewTransport(t.getSmsConfig(t.conf, name))
	return t.cache[name]
}

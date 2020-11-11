package sms

import (
	"context"

	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

type SenderFacade struct {
	factory *SenderFactory
	dynConf config.DynamicConfigReader
}

func NewSenderFacade(factory *SenderFactory, dynConf config.DynamicConfigReader) *SenderFacade {
	return &SenderFacade{factory: factory, dynConf: dynConf}
}

func (s *SenderFacade) Send(ctx context.Context, mobile, content string) error {
	sender, err := s.getRealSender(ctx)
	if err != nil {
		return err
	}
	return sender.Send(ctx, mobile, content)
}

func (s *SenderFacade) getRealSender(ctx context.Context) (contract.SmsSender, error) {
	tenant := config.GetTenant(ctx)
	conf, err := s.dynConf.Tenant(tenant)
	if err != nil {
		return nil, errors.Wrap(err, "no configuration found for sms tenant")
	}
	return s.factory.GetTransportByConf(conf), nil
}

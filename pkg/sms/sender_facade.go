package sms

import (
	"context"

	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

type SenderFacade struct {
	factory *SenderFactory
}

func NewSenderFacade(factory *SenderFactory) *SenderFacade {
	return &SenderFacade{factory: factory}
}

func (s *SenderFacade) Send(ctx context.Context, mobile, content string) error {
	sender := s.getRealSender(ctx)
	return sender.Send(ctx, mobile, content)
}

func (s *SenderFacade) getRealSender(ctx context.Context) contract.SmsSender {
	packageName := config.GetTenant(ctx).PackageName
	return s.factory.GetTransport(packageName)
}

package handlers

import (
	"context"

	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

type MonitoredAppService struct {
	eventBus EventBus
	pb.AppServer
}

type EventBus interface {
	Emit(ctx context.Context, event string, tenant *config.Tenant) error
}

func (m MonitoredAppService) Login(ctx context.Context, request *pb.UserLoginRequest) (*pb.UserInfoReply, error) {
	resp, err := m.AppServer.Login(ctx, request)
	if err != nil {
		return resp, err
	}

	claim := config.Tenant{
		PackageName: request.PackageName,
		UserId:      resp.Data.Id,
		Suuid:       request.Device.Suuid,
		Channel:     request.Channel,
		VersionCode: request.VersionCode,
	}

	if resp.Data.IsNew {
		_ = m.eventBus.Emit(ctx, "new_user", &claim)
	}

	return resp, err
}

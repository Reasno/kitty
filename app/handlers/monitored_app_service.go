package handlers

import (
	"context"

	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

type MonitoredAppService struct {
	userBus  UserBus
	eventBus EventBus
	pb.AppServer
}

type UserBus interface {
	Emit(ctx context.Context, info contract.Marshaller) error
}
type EventBus interface {
	Emit(ctx context.Context, event string, tenant *config.Tenant) error
}

func (m MonitoredAppService) Login(ctx context.Context, request *pb.UserLoginRequest) (*pb.UserInfoReply, error) {
	resp, err := m.AppServer.Login(ctx, request)
	if err != nil {
		return resp, err
	}
	// emit new user registration event

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

	// emit User
	m.emitUser(ctx, resp)
	return resp, err
}

func (m MonitoredAppService) Bind(ctx context.Context, request *pb.UserBindRequest) (*pb.UserInfoReply, error) {
	resp, err := m.AppServer.Bind(ctx, request)
	if err != nil {
		return resp, err
	}
	m.emitUser(ctx, resp)
	return resp, err
}

func (m MonitoredAppService) Unbind(ctx context.Context, request *pb.UserUnbindRequest) (*pb.UserInfoReply, error) {
	resp, err := m.AppServer.Unbind(ctx, request)
	if err != nil {
		return resp, err
	}
	m.emitUser(ctx, resp)
	return resp, err
}

func (m MonitoredAppService) UpdateInfo(ctx context.Context, request *pb.UserInfoUpdateRequest) (*pb.UserInfoReply, error) {
	resp, err := m.AppServer.UpdateInfo(ctx, request)
	if err != nil {
		return resp, err
	}
	m.emitUser(ctx, resp)
	return resp, err
}

func (m MonitoredAppService) Refresh(ctx context.Context, request *pb.UserRefreshRequest) (*pb.UserInfoReply, error) {
	resp, err := m.AppServer.Refresh(ctx, request)
	if err != nil {
		return resp, err
	}
	m.emitUser(ctx, resp)
	return resp, err
}

func (m MonitoredAppService) SoftDelete(ctx context.Context, in *pb.UserSoftDeleteRequest) (*pb.UserInfoReply, error) {
	resp, err := m.AppServer.SoftDelete(ctx, in)
	if err != nil {
		return resp, err
	}
	m.emitUser(ctx, resp)
	return resp, err
}

func (m MonitoredAppService) emitUser(ctx context.Context, resp *pb.UserInfoReply) {
	_ = m.userBus.Emit(ctx, resp.Data)
}

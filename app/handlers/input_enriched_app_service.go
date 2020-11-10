package handlers

import (
	"context"

	stdtracing "github.com/opentracing/opentracing-go"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

type InputEnrichedAppService struct {
	pb.AppServer
}

func (s InputEnrichedAppService) GetCode(ctx context.Context, in *pb.GetCodeRequest) (*pb.GenericReply, error) {
	ctx = context.WithValue(ctx, config.TenantKey, &config.Tenant{
		PackageName: in.PackageName,
	})
	return s.AppServer.GetCode(ctx, in)
}

func (s InputEnrichedAppService) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserInfoReply, error) {
	if in.Channel == "" {
		in.Channel = "N/A"
	}
	if in.PackageName == "" {
		in.PackageName = "N/A"
	}
	if in.VersionCode == "" {
		in.VersionCode = "N/A"
	}
	if in.Device == nil {
		in.Device = &pb.Device{}
	}
	if in.Device.Suuid == "" {
		in.Device.Suuid = "N/A"
	}
	ctx = context.WithValue(ctx, config.TenantKey, &config.Tenant{
		PackageName: in.PackageName,
		Suuid:       in.Device.Suuid,
		VersionCode: in.VersionCode,
		Channel:     in.Channel,
		Os:          uint8(in.Device.Os),
		Idfa:        in.Device.Idfa,
		Oaid:        in.Device.Oaid,
		Mac:         in.Device.Mac,
		AndroidId:   in.Device.AndroidId,
	})
	span := stdtracing.SpanFromContext(ctx)
	span.SetTag("package.name", in.PackageName)
	span.SetTag("suuid", in.Device.Suuid)
	resp, err := s.AppServer.Login(ctx, in)
	if err == nil {
		span.SetTag("user.id", resp.Data.Id)
	}
	return resp, err
}

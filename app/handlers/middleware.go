package handlers

import (
	"context"
	"time"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/opentracing/opentracing-go"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

type MonitoredAppService struct {
	userBus  UserBus
	eventBus EventBus
	pb.AppServer
}

func NewMonitoredAppService(userBus UserBus, eventBus EventBus, appServer *AppService) *MonitoredAppService {
	return &MonitoredAppService{userBus: userBus, eventBus: eventBus, AppServer: appServer}
}

type UserBus interface {
	Emit(ctx context.Context, info contract.Marshaller) error
}
type EventBus interface {
	Emit(ctx context.Context, event string) error
}

func (m MonitoredAppService) Login(ctx context.Context, request *pb.UserLoginRequest) (*pb.UserInfoReply, error) {
	resp, err := m.AppServer.Login(ctx, request)
	if err != nil {
		return resp, err
	}
	// emit new user registration event
	span := opentracing.SpanFromContext(ctx)
	go func() {
		ctx := opentracing.ContextWithSpan(context.Background(), span)
		ctx = context.WithValue(ctx, jwt.JWTClaimsContextKey, &kjwt.Claim{
			StandardClaims: stdjwt.StandardClaims{},
			PackageName:    request.PackageName,
			UserId:         resp.Data.Id,
			Suuid:          request.Device.Suuid,
			Channel:        request.Channel,
			VersionCode:    request.VersionCode,
			Wechat:         request.Wechat,
			Mobile:         request.Mobile,
		})

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if resp.Data.IsNew {
			_ = m.eventBus.Emit(ctx, "new_user")
		}
	}()

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

func (m MonitoredAppService) emitUser(ctx context.Context, resp *pb.UserInfoReply) {
	span := opentracing.SpanFromContext(ctx)
	go func() {
		ctx := opentracing.ContextWithSpan(context.Background(), span)
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_ = m.userBus.Emit(ctx, resp.Data)
	}()
}

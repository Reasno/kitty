package handlers

import (
	"context"
	"time"

	"github.com/Reasno/kitty/pkg/contract"
	"github.com/Reasno/kitty/pkg/kjwt"
	pb "github.com/Reasno/kitty/proto"
	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/opentracing/opentracing-go"
)

type monitoredAppService struct {
	userBus  UserBus
	eventBus EventBus
	appService
}

type UserBus interface {
	Emit(ctx context.Context, info contract.Marshaller) error
}
type EventBus interface {
	Emit(ctx context.Context, event string) error
}

func (m monitoredAppService) Login(ctx context.Context, request *pb.UserLoginRequest) (*pb.UserInfoReply, error) {
	resp, err := m.appService.Login(ctx, request)
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
			err = m.eventBus.Emit(ctx, "new_user")
			m.appService.warn(err)
		}
	}()

	// emit User
	m.emitUser(ctx, resp)
	return resp, err
}

func (m monitoredAppService) Bind(ctx context.Context, request *pb.UserBindRequest) (*pb.UserInfoReply, error) {
	resp, err := m.appService.Bind(ctx, request)
	if err != nil {
		return resp, err
	}
	m.emitUser(ctx, resp)
	return resp, err
}

func (m monitoredAppService) Unbind(ctx context.Context, request *pb.UserUnbindRequest) (*pb.UserInfoReply, error) {
	resp, err := m.appService.Unbind(ctx, request)
	if err != nil {
		return resp, err
	}
	m.emitUser(ctx, resp)
	return resp, err
}

func (m monitoredAppService) UpdateInfo(ctx context.Context, request *pb.UserInfoUpdateRequest) (*pb.UserInfoReply, error) {
	resp, err := m.appService.UpdateInfo(ctx, request)
	if err != nil {
		return resp, err
	}
	m.emitUser(ctx, resp)
	return resp, err
}

func (m monitoredAppService) Refresh(ctx context.Context, request *pb.UserRefreshRequest) (*pb.UserInfoReply, error) {
	resp, err := m.appService.Refresh(ctx, request)
	if err != nil {
		return resp, err
	}
	m.emitUser(ctx, resp)
	return resp, err
}

func (m monitoredAppService) emitUser(ctx context.Context, resp *pb.UserInfoReply) {
	span := opentracing.SpanFromContext(ctx)
	go func() {
		ctx := opentracing.ContextWithSpan(context.Background(), span)
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err := m.userBus.Emit(ctx, resp.Data)
		m.appService.warn(err)
	}()
}

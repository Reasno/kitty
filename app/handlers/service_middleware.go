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
		err := m.userBus.Emit(ctx, resp.Data)
		m.appService.warn(err)
	}()
	return resp, err
}

func (m monitoredAppService) Bind(ctx context.Context, request *pb.UserBindRequest) (*pb.UserInfoReply, error) {
	resp, err := m.appService.Bind(ctx, request)
	if err != nil {
		return resp, err
	}
	span := opentracing.SpanFromContext(ctx)
	go func() {
		ctx := opentracing.ContextWithSpan(context.Background(), span)
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err := m.userBus.Emit(ctx, resp.Data)
		m.appService.warn(err)
	}()
	return resp, err
}

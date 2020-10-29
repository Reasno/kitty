package handlers

import (
	"context"
	pb "github.com/Reasno/kitty/proto"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

func newLoginToBindMiddleware(bind endpoint.Endpoint) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if ctx.Value(jwt.JWTTokenContextKey) == nil {
				return e(ctx, request)
			}
			loginRequest, ok := request.(pb.UserLoginRequest)
			if !ok {
				return e(ctx, request)
			}
			if len(loginRequest.Mobile) <= 0 && len(loginRequest.Wechat) <= 0 {
				return e(ctx, request)
			}
			bindReq := pb.UserBindRequest{
				Mobile: loginRequest.Mobile,
				Code:   loginRequest.Code,
				Wechat: loginRequest.Wechat,
			}
			return bind(ctx, bindReq)
		}

	}

}

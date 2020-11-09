package kmiddleware

import (
	"context"
	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	kittyjwt "glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Claim struct {
	stdjwt.StandardClaims
	Uid int64
}

type SecurityConfig struct {
	Enable bool
	JwtKey string
	JwtId  string
}

func NewAuthenticationMiddleware(securityConfig *SecurityConfig) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		if !securityConfig.Enable {
			return e
		}
		kf := func(token *stdjwt.Token) (interface{}, error) {
			return []byte(securityConfig.JwtKey), nil
		}
		e = jwt.NewParser(kf, stdjwt.SigningMethodHS256, kittyjwt.ClaimFactory)(e)
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = e(ctx, request)
			if err != nil {
				err = status.Error(codes.Unauthenticated, err.Error())
			}
			return response, err
		}
	}
}

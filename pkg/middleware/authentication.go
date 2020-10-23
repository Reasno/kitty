package middleware

import (
	"context"
	kittyjwt "github.com/Reasno/kitty/pkg/jwt"
	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/spf13/viper"
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
		if !viper.GetBool("security.enable") {
			return e
		}
		kf := func(token *stdjwt.Token) (interface{}, error) {
			return viper.GetString("security.key"), nil
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

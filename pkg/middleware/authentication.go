package middleware

import (
	"context"
	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

func NewAuthenticationMiddleware() endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		if ! viper.GetBool("security.enable") {
			return e
		}
		kf := func(token *stdjwt.Token) (interface{}, error) {
			return viper.GetString("security.key"), nil
		}
		e = jwt.NewParser(kf, stdjwt.SigningMethodHS256, jwt.StandardClaimsFactory)(e)
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = e(ctx, request)
			if err != nil {
				err = status.Error(codes.Unauthenticated, err.Error())
			}
			return response, err
		}
	}
}

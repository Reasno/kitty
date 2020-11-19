package kmiddleware

import (
	"context"
	"errors"

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
			return response, unwrap(err)
		}
	}
}

func unwrap(err error) error {
	if errors.Is(err, jwt.ErrTokenInvalid) {
		err = status.Error(codes.Unauthenticated, err.Error())
		return err
	}
	if errors.Is(err, jwt.ErrTokenExpired) {
		err = status.Error(codes.Unauthenticated, err.Error())
		return err
	}
	if errors.Is(err, jwt.ErrTokenContextMissing) {
		err = status.Error(codes.Unauthenticated, err.Error())
		return err
	}
	if errors.Is(err, jwt.ErrTokenNotActive) {
		err = status.Error(codes.Unauthenticated, err.Error())
		return err
	}
	if errors.Is(err, jwt.ErrUnexpectedSigningMethod) {
		err = status.Error(codes.Unauthenticated, err.Error())
		return err
	}
	if errors.Is(err, jwt.ErrTokenMalformed) {
		err = status.Error(codes.Unauthenticated, err.Error())
		return err
	}
	return err
}

func NewOptionalAuthenticationMiddleware(securityConfig *SecurityConfig) endpoint.Middleware {
	return func(plain endpoint.Endpoint) endpoint.Endpoint {
		kf := func(token *stdjwt.Token) (interface{}, error) {
			return []byte(securityConfig.JwtKey), nil
		}
		auth := jwt.NewParser(kf, stdjwt.SigningMethodHS256, kittyjwt.ClaimFactory)(plain)
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			_, ok := ctx.Value(jwt.JWTTokenContextKey).(string)
			if !ok {
				return plain(ctx, request)
			}
			response, err = auth(ctx, request)
			return response, unwrap(err)
		}
	}
}

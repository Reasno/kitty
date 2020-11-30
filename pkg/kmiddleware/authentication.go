package kmiddleware

import (
	"context"
	"errors"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kerr"
	kittyjwt "glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
)

type Claim struct {
	stdjwt.StandardClaims
	Uid int64
}

type SecurityConfig struct {
	JwtKey string
	JwtId  string
}

var TrustedTransportKey = struct{}{}

func NewAuthenticationMiddleware(securityConfig *SecurityConfig) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		kf := func(token *stdjwt.Token) (interface{}, error) {
			return []byte(securityConfig.JwtKey), nil
		}
		e := jwt.NewParser(kf, stdjwt.SigningMethodHS256, kittyjwt.ClaimFactory)(next)
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if b, ok := ctx.Value(TrustedTransportKey).(bool); ok && b {
				return next(ctx, request)
			}
			response, err = e(ctx, request)
			return response, wrap(err)
		}
	}
}

func wrap(err error) error {
	if errors.Is(err, jwt.ErrTokenInvalid) {
		return kerr.UnauthenticatedErr(err, err.Error())
	}
	if errors.Is(err, jwt.ErrTokenExpired) {
		return kerr.UnauthenticatedErr(err, err.Error())
	}
	if errors.Is(err, jwt.ErrTokenContextMissing) {
		return kerr.UnauthenticatedErr(err, err.Error())
	}
	if errors.Is(err, jwt.ErrTokenNotActive) {
		return kerr.UnauthenticatedErr(err, err.Error())
	}
	if errors.Is(err, jwt.ErrUnexpectedSigningMethod) {
		return kerr.UnauthenticatedErr(err, err.Error())
	}
	if errors.Is(err, jwt.ErrTokenMalformed) {
		return kerr.UnauthenticatedErr(err, err.Error())
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
			return response, wrap(err)
		}
	}
}

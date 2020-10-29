package jwt

import (
	"context"
	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"time"
)

type Claim struct {
	stdjwt.StandardClaims
	UserId   uint64
	Suuid string
	Channel string
	VersionCode string
	Wechat string
	Mobile string
}

func NewClaim(uid uint64, issuer, suuid, channel, versionCode, wechat, mobile string, ttl time.Duration) *Claim {
	return &Claim{
		StandardClaims: stdjwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
		UserId:   uid,
		Suuid: suuid,
		Channel: channel,
		VersionCode: versionCode,
		Wechat: wechat,
		Mobile: mobile,
	}
}

func GetClaim(ctx context.Context) *Claim {
	return ctx.Value(jwt.JWTClaimsContextKey).(*Claim)
}

func ClaimFactory() stdjwt.Claims {
	return &Claim{}
}

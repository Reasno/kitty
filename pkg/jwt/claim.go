package jwt

import (
	"context"
	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"time"
)

type Claim struct {
	stdjwt.StandardClaims
	Uid   uint64
	Suuid string
	Channel string
	VersionCode string
	Wechat string
}

func NewClaim(uid uint64, issuer, suuid, channel, versionCode, wechat string, ttl time.Duration) *Claim {
	return &Claim{
		StandardClaims: stdjwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
		Uid:   uid,
		Suuid: suuid,
		Channel: channel,
		VersionCode: versionCode,
		Wechat: wechat,
	}
}

func GetClaim(ctx context.Context) *Claim {
	claim := ctx.Value(jwt.JWTClaimsContextKey).(Claim)
	return &claim
}

func ClaimFactory() stdjwt.Claims {
	return &Claim{}
}

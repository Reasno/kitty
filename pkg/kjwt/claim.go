package kjwt

import (
	"context"
	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"time"
)

type Claim struct {
	stdjwt.StandardClaims
	PackageName string
	UserId      uint64
	Suuid       string
	Channel     string
	VersionCode string
	Wechat      string
	Mobile      string
}

func NewClaim(uid uint64, issuer, suuid, channel, versionCode, wechat, mobile, packageName string, ttl time.Duration) *Claim {
	return &Claim{
		StandardClaims: stdjwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
		UserId:      uid,
		Suuid:       suuid,
		Channel:     channel,
		VersionCode: versionCode,
		Wechat:      wechat,
		Mobile:      mobile,
		PackageName: packageName,
	}
}

func GetClaim(ctx context.Context) *Claim {
	if c, ok := ctx.Value(jwt.JWTClaimsContextKey).(*Claim); ok {
		return c
	}
	return &Claim{}
}

func ClaimFactory() stdjwt.Claims {
	return &Claim{}
}

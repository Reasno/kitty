package kjwt

import (
	"context"
	"time"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
)

type Claim struct {
	stdjwt.StandardClaims
	PackageName  string `json:"PackageName,omitempty"`
	UserId       uint64 `json:"UserId,omitempty"`
	Suuid        string `json:"Suuid,omitempty"`
	Channel      string `json:"Channel,omitempty"`
	VersionCode  string `json:"VersionCode,omitempty"`
	Wechat       string `json:"Wechat,omitempty"`
	Mobile       string `json:"Mobile,omitempty"`
	ThirdPartyId string `json:"ThirdPartyId,omitempty"`
}

func NewClaim(uid uint64, issuer, suuid, channel, versionCode, wechat, mobile, packageName, thirdPartyId string, ttl time.Duration) *Claim {
	return &Claim{
		StandardClaims: stdjwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
		UserId:       uid,
		Suuid:        suuid,
		Channel:      channel,
		VersionCode:  versionCode,
		Wechat:       wechat,
		Mobile:       mobile,
		PackageName:  packageName,
		ThirdPartyId: thirdPartyId,
	}
}

func NewAdminClaim(issuer string, ttl time.Duration) *Claim {
	return &Claim{
		StandardClaims: stdjwt.StandardClaims{
			Audience:  "admin",
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
	}
}

func ClaimFromContext(ctx context.Context) *Claim {
	if c, ok := ctx.Value(jwt.JWTClaimsContextKey).(*Claim); ok {
		return c
	}
	return &Claim{}
}

func ClaimFactory() stdjwt.Claims {
	return &Claim{}
}

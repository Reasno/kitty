package jwt

import (
	stdjwt "github.com/dgrijalva/jwt-go"
	"time"
)

type Claim struct {
	stdjwt.StandardClaims
	Uid   uint64
	Suuid string
}

func NewClaim(uid uint64, issuer, suuid string, ttl time.Duration) *Claim {
	return &Claim{
		StandardClaims: stdjwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
		Uid:   uid,
		Suuid: suuid,
	}
}

func ClaimFactory() stdjwt.Claims {
	return &Claim{}
}

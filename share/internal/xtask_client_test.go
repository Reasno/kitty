package internal

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/stretchr/testify/assert"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract/mocks"
)

func TestEndpoints_Request(t *testing.T) {
	conf := mocks.ConfigReader{}
	conf.On("String", "xtask.url").Return("http://120.31.70.243:8989/xtasks/score/madd")

	e, err := NewXTaskRequester(&conf, http.DefaultClient)
	assert.NoError(t, err)
	ctx := context.WithValue(context.Background(), jwt.JWTTokenContextKey, "eyJhbGciOiJIUzI1NiIsImtpZCI6ImtpdHR5IiwidHlwIjoiSldUIn0.eyJleHAiOjE2NDE3ODAyMTAsImlhdCI6MTYwNTc4MDIxMCwiaXNzIjoic2lnbkNtZCIsIlBhY2thZ2VOYW1lIjoiY29tLmRvbmV3cy53d3ciLCJVc2VySWQiOjF9._pwH3bVLc8tWDDvqtQGlIHCK-sITrnWH6oaym8GSGk8")
	resp, err := e.Request(ctx, &XTaskRequest{
		"Test",
		1000,
		"1",
		"1",
	})
	assert.NoError(t, err)
	fmt.Printf("%+v\n", resp)
}

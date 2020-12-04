package internal

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/stretchr/testify/assert"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract/mocks"
)

var useXtask bool

func init() {
	flag.BoolVar(&useXtask, "xtask", false, "use xtask for testing")
}

func TestEndpoints_Request(t *testing.T) {
	if !useXtask {
		t.Skip("need xtask to test enpoints")
	}
	conf := mocks.ConfigReader{}
	conf.On("String", "xtask.url").Return("https://xtasks.dev.tagtic.cn/xtasks/score/madd")

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

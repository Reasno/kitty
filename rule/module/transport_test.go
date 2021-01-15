package module

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"glab.tagtic.cn/ad_gains/kitty/rule/dto"
)

func TestDecodePayload(t *testing.T) {
	cases := []struct {
		name    string
		request *http.Request
		asserts func(t *testing.T, payload *dto.Payload)
	}{
		{
			"decode query",
			func() *http.Request {
				r, _ := http.NewRequest("GET", "http://example.org?foo=bar&foo=baz", nil)
				return r
			}(),
			func(t *testing.T, payload *dto.Payload) {
				assert.Contains(t, payload.Q["foo"], "bar")
				assert.Contains(t, payload.Q["foo"], "baz")
			},
		},
		{
			"decode body",
			func() *http.Request {
				r, _ := http.NewRequest("POST", "http://example.org", strings.NewReader(`{"foo":"bar"}`))
				return r
			}(),
			func(t *testing.T, payload *dto.Payload) {
				assert.Contains(t, payload.B["foo"], "bar")
			},
		},
	}

	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			p := dto.Payload{}
			err := decodePayload(&p, cc.request)
			assert.NoError(t, err)
			cc.asserts(t, &p)
		})
	}
}

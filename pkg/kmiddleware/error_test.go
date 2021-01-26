package kmiddleware

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kerr"
	"testing"
)

func TestNewErrorMarshallerMiddleware(t *testing.T) {
	mw := NewErrorMarshallerMiddleware(false)
	e1 := func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, errors.New("foo")
	}
	e2 := func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, kerr.NotFoundErr(errors.New("bar"), "")
	}
	e3 := func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, kerr.NotFoundErr(kerr.InvalidArgumentErr(errors.New("bar"), ""), "")
	}
	e4 := func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, errors.Wrap(kerr.NotFoundErr(errors.New("foo"), ""), "bar")
	}
	cases := []endpoint.Endpoint{e1, e2, e3, e4}
	for _, c := range cases {
		cc := c
		t.Run("", func(t *testing.T) {
			_, err := mw(cc)(nil, nil)
			if _, ok := err.(kerr.ServerError); !ok {
				t.Fail()
			}
		})
	}
}

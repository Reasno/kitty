package kmiddleware

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"glab.tagtic.cn/ad_gains/kitty/app/msg"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kerr"
)

type validator interface {
	Validate() error
}

func NewValidationMiddleware() endpoint.Middleware {
	return func(in endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			if t, ok := req.(validator); ok {
				err = t.Validate()
				if err != nil {
					return nil, kerr.InvalidArgumentErr(err, msg.InvalidParams)
				}
			}
			resp, err = in(ctx, req)
			return
		}
	}
}

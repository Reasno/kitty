package kmiddleware

import (
	"context"
	"github.com/Reasno/kitty/pkg/kerr"
	"github.com/go-kit/kit/endpoint"
)

func NewErrorMarshallerMiddleware() endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = e(ctx, request)
			if _, ok := err.(kerr.ServerError); err != nil && !ok {
				err = kerr.UnknownErr(err)
			}
			return response, err
		}
	}
}

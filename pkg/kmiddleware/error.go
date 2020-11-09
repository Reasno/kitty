package kmiddleware

import (
	"context"
	"errors"
	"github.com/Reasno/kitty/pkg/kerr"
	"github.com/go-kit/kit/endpoint"
)

func NewErrorMarshallerMiddleware() endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = e(ctx, request)

			if err != nil {
				var serverError kerr.ServerError
				if !errors.As(err, &serverError) {
					serverError = kerr.UnknownErr(err)
				}
				// Brings kerr.SeverError to the uppermost level
				return response, serverError
			}

			return response, nil
		}
	}
}

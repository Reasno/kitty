package kmiddleware

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kerr"
)

func NewErrorMarshallerMiddleware() endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			//defer func() {
			//	if er := recover(); er != nil {
			//		err = fmt.Errorf("panic: %s", er)
			//	}
			//}()
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

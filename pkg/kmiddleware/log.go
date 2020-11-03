package kmiddleware

import (
	"context"
	"fmt"
	log2 "github.com/Reasno/kitty/pkg/klog"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

func NewLoggingMiddleware(logger log.Logger, printTrace bool) endpoint.Middleware {
	return func(endpoint endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			l := log2.WithContext(level.Debug(logger), ctx)
			defer l.Log("req", request, "response", response)
			response, err = endpoint(ctx, request)
			if err != nil {
				l.Log("err", err)
				if err, ok := err.(interface{ StackTrace() errors.StackTrace }); printTrace && ok {
					fmt.Printf("\n%+v\n\n", err)
				}
			}
			return response, err
		}
	}
}

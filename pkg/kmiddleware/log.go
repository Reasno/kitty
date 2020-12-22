package kmiddleware

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/klog"
)

func NewLoggingMiddleware(logger log.Logger, printTrace bool) endpoint.Middleware {
	return func(endpoint endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			l := klog.WithContext(level.Debug(logger), ctx)
			response, err = endpoint(ctx, request)
			if err != nil {
				l.Log("err", err.Error())
				if stacktracer, ok := err.(interface{ StackTrace() errors.StackTrace }); printTrace && ok {
					fmt.Printf("\n%+v\n\n", stacktracer.StackTrace())
				}
			}
			l.Log("request", fmt.Sprintf("%+v", request), "response", fmt.Sprintf("%+v", response))
			return response, err
		}
	}
}

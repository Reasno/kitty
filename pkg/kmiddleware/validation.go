package kmiddleware

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
					return nil, status.Error(codes.InvalidArgument, err.Error())
				}
			}
			resp, err = in(ctx, req)
			return
		}
	}
}

type ValidationError struct {
	err error
}

func (ve ValidationError) StatusCode() int {
	return 400
}
func (ve ValidationError) GRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, ve.err.Error())
}
func (ve ValidationError) Error() string {
	return ve.err.Error()
}

package middleware

import (
	"context"
	"github.com/go-kit/kit/endpoint"
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
					return nil, newValidationError(err)
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
func (ve ValidationError) BusinessCode() int {
	return 100
}
func (ve ValidationError) Error() string {
	return ve.err.Error()
}

func newValidationError(err error) ValidationError {
	return ValidationError{err}
}

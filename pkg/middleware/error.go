package middleware

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewErrorMashallerMiddleware() endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = e(ctx, request)
			if err != nil {
				err = newJsonError(err)
			}
			return response, err
		}
	}
}

func newJsonError(e error) JsonError {
	s, ok := status.FromError(e)
	if !ok {
		s = status.New(codes.Unknown, e.Error())
	}
	return JsonError{e, s}
}

type JsonError struct {
	error `json:"message"`
	status *status.Status
}

type jsonRep struct {
	Code codes.Code `json:"code"`
	Message string `json:"message"`
	Details interface{} `json:"details"`
}

func (e JsonError) MarshalJSON() ([]byte, error) {
	r := jsonRep{
		e.status.Code(),
		e.status.Message(),
		e.status.Details(),
	}
	return json.Marshal(r)
}

func (e JsonError) GRPCStatus() *status.Status {
	return e.status
}

// StatusCode Implements https status
func (e JsonError) StatusCode() int {
	switch e.status.Code() {
	case codes.OK:
		return 200
	case codes.Canceled:
		return 499
	case codes.Unknown:
		return 500
	case codes.InvalidArgument:
		return 400
	case codes.DeadlineExceeded:
		return 504
	case codes.NotFound:
		return 404
	case codes.AlreadyExists:
		return 409
	case codes.PermissionDenied:
		return 403
	case codes.ResourceExhausted:
		return 429
	case codes.FailedPrecondition:
		return 400
	case codes.Aborted:
		return 409
	case codes.OutOfRange:
		return 400
	case codes.Unimplemented:
		return 501
	case codes.DataLoss:
		return 500
	case codes.Unauthenticated:
		return 401
	default:
		return 500
	}
}


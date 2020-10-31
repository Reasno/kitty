package kerr

import (
	"encoding/json"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func err(code codes.Code, e error) ServerError {
	return ServerError{e, status.New(code, redact(e)), uint32(code)}
}

func UnknownErr(e error) ServerError {
	s, ok := status.FromError(e)
	if !ok {
		s = status.New(codes.Unknown, redact(e))
	}
	return ServerError{e, s, uint32(s.Code())}
}

func InvalidArgumentErr(e error) ServerError {
	return err(codes.InvalidArgument, e)
}

func NotFoundErr(e error) ServerError {
	return err(codes.NotFound, e)
}

func InternalErr(e error) ServerError {
	return err(codes.Internal, e)
}

func UnauthorizedErr(e error) ServerError {
	return err(codes.Unauthenticated, e)
}

func CustomErr(code uint32, e error) ServerError {
	return ServerError{e, status.New(codes.Internal, redact(e)), code}
}

func redact(err error) string {
	return strings.Split(err.Error(), ":")[0]
}

type ServerError struct {
	err        error
	status     *status.Status
	customCode uint32
}

type jsonRep struct {
	Code    uint32      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

func (e ServerError) MarshalJSON() ([]byte, error) {
	r := jsonRep{
		e.customCode,
		e.status.Message(),
		e.status.Details(),
	}
	return json.Marshal(r)
}

func (e ServerError) Error() string {
	return e.err.Error()
}

func (e ServerError) GRPCStatus() *status.Status {
	return e.status
}

// StatusCode Implements https status
func (e ServerError) StatusCode() int {
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

// Unwrap implements go's standard errors.Unwrap() interface
func (e ServerError) Unwrap() error {
	return e.err
}

// StackTrace implements the interface of errors.Wrap()
func (e ServerError) StackTrace() errors.StackTrace {
	if err, ok := e.err.(stackTracer); ok {
		return err.StackTrace()
	}
	return errors.Wrap(e.err, "").(stackTracer).StackTrace()
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

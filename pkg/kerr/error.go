package kerr

import (
	"context"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
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

func ResourceExhaustedErr(e error) ServerError {
	return err(codes.ResourceExhausted, e)
}

func FailedPreconditionErr(e error) ServerError {
	return err(codes.FailedPrecondition, e)
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
	return 200
	//switch e.status.Code() {
	//case codes.OK:
	//	return 200
	//case codes.Canceled:
	//	return 499
	//case codes.Unknown:
	//	return 500
	//case codes.InvalidArgument:
	//	return 400
	//case codes.DeadlineExceeded:
	//	return 504
	//case codes.NotFound:
	//	return 404
	//case codes.AlreadyExists:
	//	return 409
	//case codes.PermissionDenied:
	//	return 403
	//case codes.ResourceExhausted:
	//	return 429
	//case codes.FailedPrecondition:
	//	return 400
	//case codes.Aborted:
	//	return 409
	//case codes.OutOfRange:
	//	return 400
	//case codes.Unimplemented:
	//	return 501
	//case codes.DataLoss:
	//	return 500
	//case codes.Unauthenticated:
	//	return 401
	//default:
	//	return 500
	//}
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
	return errors.WithStack(e.err).(stackTracer).StackTrace()
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// ErrorEncoder writes the error to the ResponseWriter, by default a content
// type of application/json, a body of json with key "error" and the value
// error.Error(), and a status code of 500. If the error implements Headerer,
// the provided headers will be applied to the response. If the error
// implements json.Marshaler, and the marshaling succeeds, the JSON encoded
// form of the error will be used. If the error implements StatusCoder, the
// provided StatusCode will be used instead of 500.
func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	const contentType = "application/json; charset=utf-8"
	body, _ := json.Marshal(errorWrapper{Message: err.Error(), Code: 2})
	if marshaler, ok := err.(json.Marshaler); ok {
		if jsonBody, marshalErr := marshaler.MarshalJSON(); marshalErr == nil {
			body = jsonBody
		}
	}
	w.Header().Set("Content-Type", contentType)
	if headerer, ok := err.(httptransport.Headerer); ok {
		for k := range headerer.Headers() {
			w.Header().Set(k, headerer.Headers().Get(k))
		}
	}
	code := http.StatusInternalServerError
	if sc, ok := err.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	w.Write(body)
}

type errorWrapper struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

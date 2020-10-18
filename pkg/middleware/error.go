package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
)

func NewErrorMashallerMiddleware() endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = e(ctx, request)
			if err != nil {
				rerr := &JsonError{error: err, Code: -1}
				if t, ok := err.(BusinessCoder); ok {
					rerr = rerr.WithBusinessCode(t.BusinessCode())
				}
				if t, ok := err.(Detailer); ok {
					rerr = rerr.WithDetail(t.Detail())
				}
				err = *rerr
			}
			return response, err
		}
	}
}

type BusinessCoder interface{
	BusinessCode() int
}

type Detailer interface{
	Detail() interface{}
}

type JsonError struct {
	error `json:"message"`
	Code int `json:"code"`
	Detail interface{} `json:"detail"`
}

func (e *JsonError) WithDetail(detail interface{}) *JsonError  {
	e.Detail = detail
	return e
}

func (e *JsonError) WithBusinessCode(code int) *JsonError {
	e.Code = code
	return e
}


func (e JsonError) MarshalJSON() ([]byte, error) {
	detail, err := json.Marshal(e.Detail)
	if err != nil || len(detail) <= 0 {
		detail = []byte("[]")
	}
	str := fmt.Sprintf("{\"code\": %d, \"message\":\"%s\", \"detail\": %s}", e.Code, e.Error(), detail)
	return []byte(str), nil
}

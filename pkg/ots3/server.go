package ots3

import (
	"context"
	"fmt"
	"github.com/Reasno/kitty/pkg/contract"
	"github.com/Reasno/kitty/pkg/kerr"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	httptransport "github.com/go-kit/kit/transport/http"
	"io"
	"net/http"
)

type UploadService struct {
	logger log.Logger
	s3     *Manager
}

func (s *UploadService) Upload(ctx context.Context, data io.Reader) (url string, err error) {
	defer func() {
		if closer, ok := data.(io.ReadCloser); ok {
			closer.Close()
		}
	}()
	defer level.Info(s.logger).Log("msg", fmt.Sprintf("file uploaded to %s", url))
	return s.s3.Upload(ctx, data)
}

type Request struct {
	data io.Reader
}

type Response struct {
	Data struct {
		Url string `json:"url"`
	} `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

func MakeUploadEndpoint(uploader contract.Uploader) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*Request)
		resp, err := uploader.Upload(ctx, req.data)
		if err != nil {
			return nil, kerr.InternalErr(err)
		}
		return &Response{
			Code: 0,
			Data: struct {
				Url string `json:"url"`
			}{Url: resp},
		}, nil
	}
}

func MakeHttpHandler(endpoint endpoint.Endpoint) http.Handler {
	server := httptransport.NewServer(
		endpoint,
		decodeRequest,
		httptransport.EncodeJSONResponse,
	)
	return server
}

func decodeRequest(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
	return &Request{
		data: request2.Body,
	}, nil
}

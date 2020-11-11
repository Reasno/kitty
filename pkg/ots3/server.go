package ots3

import (
	"context"
	"fmt"
	"mime"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	httptransport "github.com/go-kit/kit/transport/http"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kerr"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kmiddleware"

	"io"
	"net/http"
)

type UploadService struct {
	logger log.Logger
	s3     *Manager
}

func (s *UploadService) Upload(ctx context.Context, reader io.Reader) (newUrl string, err error) {
	defer func() {
		if closer, ok := reader.(io.ReadCloser); ok {
			closer.Close()
		}
	}()
	newUrl, err = s.s3.Upload(ctx, reader)
	level.Info(s.logger).Log("msg", fmt.Sprintf("file uploaded to %s", newUrl))
	return newUrl, err
}

type Request struct {
	contentType string
	data        io.Reader
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
		ext, _ := mime.ExtensionsByType(req.contentType)
		if len(ext) == 0 {
			ext = []string{""}
		}
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

func Middleware(logger log.Logger, env contract.Env) endpoint.Middleware {
	l := kmiddleware.NewLoggingMiddleware(logger, env.IsLocal())
	e := kmiddleware.NewErrorMarshallerMiddleware()
	return endpoint.Chain(e, l)
}

func MakeHttpHandler(endpoint endpoint.Endpoint, middleware endpoint.Middleware) http.Handler {
	server := httptransport.NewServer(
		middleware(endpoint),
		decodeRequest,
		httptransport.EncodeJSONResponse,
	)
	return server
}

func decodeRequest(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
	var ContentType = http.CanonicalHeaderKey("Content-Type")
	return &Request{
		contentType: request2.Header.Get(ContentType),
		data:        request2.Body,
	}, nil
}

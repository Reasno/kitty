package ots3

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/auth/jwt"
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
			return nil, kerr.InternalErr(err, "上传失败")
		}
		return &Response{
			Code: 0,
			Data: struct {
				Url string `json:"url"`
			}{Url: resp},
		}, nil
	}
}

func Middleware(logger log.Logger, env contract.Env, config *kmiddleware.SecurityConfig) endpoint.Middleware {
	l := kmiddleware.NewLoggingMiddleware(logger, env.IsLocal())
	e := kmiddleware.NewErrorMarshallerMiddleware()
	a := kmiddleware.NewAuthenticationMiddleware(config)
	return endpoint.Chain(e, l, a)
}

func provideSecurityConfig(conf contract.ConfigReader) *kmiddleware.SecurityConfig {
	return &kmiddleware.SecurityConfig{
		JwtKey: conf.String("security.key"),
		JwtId:  conf.String("security.kid"),
	}
}

func MakeHttpHandler(endpoint endpoint.Endpoint, middleware endpoint.Middleware) http.Handler {
	server := httptransport.NewServer(
		middleware(endpoint),
		decodeRequest,
		httptransport.EncodeJSONResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
	)
	return server
}

func decodeRequest(_ context.Context, request2 *http.Request) (request interface{}, err error) {
	return &Request{
		data: request2.Body,
	}, nil
}

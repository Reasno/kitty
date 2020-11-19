package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kmiddleware"
)

type XTaskClient struct {
	Doer contract.HttpDoer
	Conf contract.ConfigReader
}

type XTaskRequest struct {
	ScoreDesc  string `json:"score_desc"`
	ScoreValue int    `json:"score_value"`
	TaskId     string `json:"task_id"`
	UniqueId   string `json:"unique_id"`
}

type XTaskResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Icon         string `json:"icon"`
		Name         string `json:"name"`
		CurrentScore int    `json:"current_score"`
		TodayScore   int    `json:"today_score"`
		TotalScore   int    `json:"total_score"`
		Money        string `json:"money"`
	} `json:"data"`
}

type Endpoints struct {
	RequestEndpoint endpoint.Endpoint
}

func (e Endpoints) Request(ctx context.Context, dto *XTaskRequest) (*XTaskResponse, error) {
	resp, err := e.RequestEndpoint(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp.(*XTaskResponse), nil
}

func NewXTaskRequester(reader contract.ConfigReader, doer contract.HttpDoer) (XTaskRequester, error) {
	instance := reader.String("xtask.url")
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	uri, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}
	RequestEndpoint := httptransport.NewClient(
		"POST",
		uri,
		httptransport.EncodeJSONRequest,
		func(ctx context.Context, response2 *http.Response) (response interface{}, err error) {
			if response2.StatusCode >= 400 {
				return nil, errors.Wrapf(ErrFailedXtaskRequest, "status: %d", response2.StatusCode)
			}
			var xTaskResp XTaskResponse
			err = json.NewDecoder(response2.Body).Decode(&xTaskResp)
			if err != nil {
				return nil, errors.Wrap(err, "xtask decode body failure")
			}
			return &xTaskResp, nil
		},
		httptransport.SetClient(doer),
		httptransport.ClientBefore(func(ctx context.Context, request *http.Request) context.Context {
			token, _ := ctx.Value(jwt.JWTTokenContextKey).(string)
			request.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))
			return ctx
		}),
	).Endpoint()

	mw := kmiddleware.Retry(3, 3*time.Second)
	return Endpoints{RequestEndpoint: mw(RequestEndpoint)}, nil
}

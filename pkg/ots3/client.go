package ots3

import (
	"context"
	"encoding/json"
	"github.com/Reasno/kitty/pkg/contract"
	httptransport "github.com/go-kit/kit/transport/http"
	"io/ioutil"
	"net/http"
	"net/url"
)

func NewClient(conf contract.ConfigReader) *httptransport.Client {
	var u, _ = url.Parse(conf.String("url"))
	return httptransport.NewClient("POST", u, encodeClientRequest, decodeClientResponse)
}

func decodeClientResponse(_ context.Context, response2 *http.Response) (response interface{}, err error) {
	defer response2.Body.Close()
	b, err := ioutil.ReadAll(response2.Body)
	if err != nil {
		return nil, err
	}
	var resp Response
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func encodeClientRequest(ctx context.Context, request *http.Request, i interface{}) error {
	request.Header.Set("Content-Type", "application/octet-stream")
	request.Body = i.(Request).data
	return nil
}

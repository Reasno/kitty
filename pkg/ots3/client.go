package ots3

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewClient(uri *url.URL) *httptransport.Client {
	return httptransport.NewClient("POST", uri, encodeClientRequest, decodeClientResponse)
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
	if input, ok := i.(Request).data.(io.ReadCloser); ok {
		request.Body = input
		return nil
	}
	request.Body = ioutil.NopCloser(i.(Request).data)
	return nil
}

type ClientUploader struct {
	endpoint endpoint.Endpoint
}

func (c ClientUploader) Upload(ctx context.Context, reader io.Reader) (newUrl string, err error) {
	resp, err := c.endpoint(ctx, Request{data: reader})
	if err != nil {
		return "", err
	}
	return resp.(Response).Data.Url, err
}

func NewClientUploader(client *httptransport.Client) *ClientUploader {
	return &ClientUploader{endpoint: client.Endpoint()}
}

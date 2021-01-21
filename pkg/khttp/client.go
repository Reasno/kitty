package khttp

import (
	"bytes"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Client struct {
	tracer     opentracing.Tracer
	underlying contract.HttpDoer
}

func NewClient(tracer opentracing.Tracer) *Client {
	baseClient := &http.Client{Transport: &nethttp.Transport{}}
	return &Client{tracer, baseClient}
}

func NewClientWithDoer(tracer opentracing.Tracer, doer contract.HttpDoer) *Client {
	return &Client{tracer, doer}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req, tracer := nethttp.TraceRequest(c.tracer, req)
	defer tracer.Finish()
	response, err := c.underlying.Do(req)
	if err != nil {
		return response, err
	}
	length, _ := strconv.Atoi(req.Header.Get(http.CanonicalHeaderKey("Content-Length")))
	if length > 1000 {
		return response, err
	}
	var buf bytes.Buffer
	byt, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response, errors.Wrap(err, "cannot read response body")
	}
	tracer.Span().LogKV("response", string(byt))
	buf.Write(byt)
	response.Body = ioutil.NopCloser(&buf)
	return response, err
}

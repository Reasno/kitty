package khttp

import (
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"net/http"
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
	return c.underlying.Do(req)
}

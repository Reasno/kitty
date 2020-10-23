package http

import (
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type Client struct {
	tracer opentracing.Tracer
	http.Client
}

func NewClient( tracer opentracing.Tracer) *http.Client {
	return &http.Client{Transport:  &nethttp.Transport{}}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req, tracer := nethttp.TraceRequest(c.tracer, req)
	defer tracer.Finish()
	return c.Client.Do(req)
}

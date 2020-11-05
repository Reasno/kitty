package kkafka

import (
	"context"
	"net"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/segmentio/kafka-go"
)

type Transport struct {
	underlying kafka.RoundTripper
	tracer opentracing.Tracer
}

func NewTransport(underlying kafka.RoundTripper, tracer opentracing.Tracer) *Transport {
	return &Transport{
		underlying: underlying,
		tracer:     tracer,
	}
}

func (t *Transport) RoundTrip(ctx context.Context, addr net.Addr, request kafka.Request) (kafka.Response, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "kafka")
	defer span.Finish()
	span.LogFields(log.String("addr", addr.String()))
	resp, err := t.underlying.RoundTrip(ctx, addr, request)
	if err != nil {
		span.SetTag("error", err.Error())
	}
	return resp, err
}

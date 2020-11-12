package kkafka

import (
	"context"
	"net"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/segmentio/kafka-go"
)

type Transport struct {
	underlying kafka.RoundTripper
	tracer     opentracing.Tracer
}

func NewTransport(underlying kafka.RoundTripper, tracer opentracing.Tracer) *Transport {
	return &Transport{
		underlying: underlying,
		tracer:     tracer,
	}
}

func (t *Transport) RoundTrip(ctx context.Context, addr net.Addr, request kafka.Request) (kafka.Response, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "kafka")
	defer span.Finish()
	ext.SpanKind.Set(span, ext.SpanKindProducerEnum)
	ext.PeerAddress.Set(span, addr.String())
	resp, err := t.underlying.RoundTrip(ctx, addr, request)
	if err != nil {
		span.LogKV("error", err.Error())
		ext.Error.Set(span, true)
	}
	return resp, err
}

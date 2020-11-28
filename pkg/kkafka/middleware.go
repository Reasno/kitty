package kkafka

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type Middleware func(h Handler) Handler

func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}

func ErrorLogMiddleware(logger log.Logger) Middleware {
	return func(h Handler) Handler {
		return HandleFunc(func(ctx context.Context, msg kafka.Message) error {
			err := h.Handle(ctx, msg)
			if err != nil {
				level.Warn(logger).Log("err", err.Error(), "topic", msg.Topic)
				return err
			}
			return nil
		})
	}
}

func TracingConsumerMiddleware(tracer opentracing.Tracer, opName string) Middleware {
	return func(h Handler) Handler {
		return HandleFunc(func(ctx context.Context, msg kafka.Message) error {
			carrier := getCarrier(&msg)
			spanContext, err := tracer.Extract(opentracing.TextMap, carrier)
			span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, opName, opentracing.FollowsFrom(spanContext))
			defer span.Finish()

			ext.SpanKind.Set(span, ext.SpanKindConsumerEnum)
			ext.PeerService.Set(span, "kafka")
			span.SetTag("topic", msg.Topic)
			span.SetTag("partition", msg.Partition)
			span.SetTag("offset", msg.Offset)

			err = h.Handle(ctx, msg)
			if err != nil {
				span.LogKV("error", err.Error())
				ext.Error.Set(span, true)
				// This is user's fault, so we are not returning any error here
				return err
			}
			return nil
		})
	}
}

func TracingProducerMiddleware(tracer opentracing.Tracer, opName string) Middleware {
	return func(h Handler) Handler {
		return HandleFunc(func(ctx context.Context, msg kafka.Message) error {
			span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, opName)
			defer span.Finish()
			ext.SpanKind.Set(span, ext.SpanKindProducerEnum)

			carrier := make(opentracing.TextMapCarrier)
			err := tracer.Inject(span.Context(), opentracing.TextMap, carrier)
			if err != nil {
				return errors.Wrap(err, "unable to inject tracing context")
			}

			var header kafka.Header
			for k, v := range carrier {
				header.Key = k
				header.Value = []byte(v)
			}
			msg.Headers = append(msg.Headers, header)

			err = h.Handle(ctx, msg)
			if err != nil {
				span.LogKV("error", err.Error())
				ext.Error.Set(span, true)
				// This is user's fault, so we are not returning any error here
				return err
			}

			return nil
		})
	}
}

func getCarrier(msg *kafka.Message) opentracing.TextMapCarrier {

	var mapCarrier = make(opentracing.TextMapCarrier)
	if msg.Headers != nil {
		for _, v := range msg.Headers {
			mapCarrier[v.Key] = string(v.Value)
		}
	}
	return mapCarrier
}

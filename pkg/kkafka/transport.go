package kkafka

import (
	"context"
	"net"
	"strings"

	"github.com/oklog/run"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type Transport struct {
	underlying kafka.RoundTripper
	tracer     opentracing.Tracer
	topic      string
}

func NewTransport(underlying kafka.RoundTripper, tracer opentracing.Tracer, topic string) *Transport {
	return &Transport{
		underlying: underlying,
		tracer:     tracer,
		topic:      topic,
	}
}

func (t *Transport) RoundTrip(ctx context.Context, addr net.Addr, request kafka.Request) (kafka.Response, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "kafka.transport")
	defer span.Finish()
	ext.SpanKind.Set(span, ext.SpanKindProducerEnum)
	ext.PeerAddress.Set(span, addr.String())
	ext.MessageBusDestination.Set(span, t.topic)
	resp, err := t.underlying.RoundTrip(ctx, addr, request)
	if err != nil {
		span.LogKV("error", err.Error())
		ext.Error.Set(span, true)
	}
	return resp, err
}

type HandleFunc func(ctx context.Context, msg kafka.Message) error

type sub struct {
	*kafka.Reader
	tracer      opentracing.Tracer
	opName      string
	parallelism int
}

// RunOnce 执行一次fetch-handle-commit。如果handle出错不会返回error，fetch或commit出错会返回error。
// 可以套在RunGroup里并行消费
func (s *sub) runOnce(ctx context.Context, fn HandleFunc) error {
	msg, err := s.FetchMessage(ctx)
	carrier := s.getCarrier(&msg)
	spanContext, _ := s.tracer.Extract(opentracing.TextMap, carrier)
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, s.opName, opentracing.FollowsFrom(spanContext))
	defer span.Finish()
	ext.SpanKind.Set(span, ext.SpanKindConsumerEnum)
	ext.PeerAddress.Set(span, strings.Join(s.Reader.Config().Brokers, ","))
	ext.PeerService.Set(span, "kafka")
	span.SetTag("topic", s.Reader.Config().Topic)
	span.SetTag("partition", s.Reader.Config().Partition)
	span.SetTag("offset", s.Reader.Offset())
	span.SetTag("lag", s.Reader.Lag())
	if err != nil {
		span.LogKV("fetch.error", err.Error())
		ext.Error.Set(span, true)
		return err
	}
	err = fn(ctx, msg)
	if err != nil {
		span.LogKV("handle.error", err.Error())
		ext.Error.Set(span, true)
		// This is user's fault, so we are not returning any error here
		return nil
	}
	err = s.CommitMessages(ctx, msg)
	if err != nil {
		span.LogKV("commit.error", err.Error())
		ext.Error.Set(span, true)
		return err
	}
	return nil
}

func (s *sub) Subscribe(ctx context.Context, fn HandleFunc) error {
	var g run.Group
	ctx, cancel := context.WithCancel(ctx)
	for i := 0; i < s.parallelism; i++ {
		g.Add(func() error {
			return s.runOnce(ctx, fn)
		}, func(err error) {
			cancel()
		})
	}
	return g.Run()
}

func (s *sub) getCarrier(msg *kafka.Message) opentracing.TextMapCarrier {
	var mapCarrier = make(opentracing.TextMapCarrier)
	if msg.Headers != nil {
		for _, v := range msg.Headers {
			mapCarrier[v.Key] = string(v.Value)
		}
	}
	return mapCarrier
}

type pub struct {
	*kafka.Writer
	topic  string
	tracer opentracing.Tracer
	opName string
}

func (p *pub) Publish(ctx context.Context, msgs ...kafka.Message) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, p.tracer, p.opName)
	defer span.Finish()
	ext.SpanKind.Set(span, ext.SpanKindProducerEnum)
	ext.MessageBusDestination.Set(span, p.topic)
	carrier := make(opentracing.TextMapCarrier)
	err := p.tracer.Inject(span.Context(), opentracing.TextMap, carrier)
	if err != nil {
		return errors.Wrap(err, "unable to inject tracing context")
	}
	var header kafka.Header
	for k, v := range carrier {
		header.Key = k
		header.Value = []byte(v)
	}
	for i := range msgs {
		msgs[i].Headers = append(msgs[i].Headers, header)
	}
	return p.Writer.WriteMessages(ctx, msgs...)
}

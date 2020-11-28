package kkafka

import (
	"context"

	"net"

	"github.com/oklog/run"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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
	span.LogKV("api", request.ApiKey())
	resp, err := t.underlying.RoundTrip(ctx, addr, request)
	if err != nil {
		span.LogKV("error", err.Error())
		ext.Error.Set(span, true)
	}
	return resp, err
}

type HandleFunc func(ctx context.Context, msg kafka.Message) error

func (h HandleFunc) Handle(ctx context.Context, msg kafka.Message) error {
	return h(ctx, msg)
}

type Handler interface {
	Handle(ctx context.Context, msg kafka.Message) error
}

type Subscriber struct {
	reader      *kafka.Reader
	handler     Handler
	parallelism int
}

func (s *Subscriber) ServeOnce(ctx context.Context) error {
	msg, err := s.reader.ReadMessage(ctx)
	if err != nil {
		return err
	}
	// User space error will not result in a transport error
	_ = s.handler.Handle(ctx, msg)
	return nil
}

func (s *Subscriber) Serve(ctx context.Context) error {
	var (
		g  run.Group
		ch chan kafka.Message
	)
	ctx, cancel := context.WithCancel(ctx)
	g.Add(func() error {
		for {
			msg, err := s.reader.ReadMessage(ctx)
			if err != nil {
				return err
			}
			ch <- msg
		}
	}, func(err error) {
		cancel()
	})

	for i := 0; i < s.parallelism; i++ {
		g.Add(func() error {
			for {
				select {
				case msg := <-ch:
					_ = s.handler.Handle(ctx, msg)
				case <-ctx.Done():
					return nil
				}
			}
		}, func(err error) {
			cancel()
		})
	}
	return g.Run()
}

type pub struct {
	*kafka.Writer
}

func (p *pub) Handle(ctx context.Context, msg kafka.Message) error {
	return p.Writer.WriteMessages(ctx, msg)
}

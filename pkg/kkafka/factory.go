package kkafka

import (
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"glab.tagtic.cn/ad_gains/kitty/pkg/klog"
)

type KafkaProducerFactory struct {
	tracer  opentracing.Tracer
	mutex   sync.Mutex
	cache   map[string]*kafka.Writer
	brokers []string
	closers []func() error
	logger  log.Logger
}

func NewKafkaProducerFactory(brokers []string, logger log.Logger) *KafkaProducerFactory {
	return &KafkaProducerFactory{
		cache:   map[string]*kafka.Writer{},
		brokers: brokers,
		closers: []func() error{},
		logger:  logger,
	}
}

func NewKafkaProducerFactoryWithTracer(brokers []string, logger log.Logger, tracer opentracing.Tracer) *KafkaProducerFactory {
	return &KafkaProducerFactory{
		tracer:  tracer,
		cache:   map[string]*kafka.Writer{},
		brokers: brokers,
		closers: []func() error{},
		logger:  logger,
	}
}

func (k *KafkaProducerFactory) Writer(topic string) *kafka.Writer {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	if w, ok := k.cache[topic]; ok {
		return w
	}
	writer := &kafka.Writer{
		Addr:        kafka.TCP(k.brokers...),
		Topic:       topic,
		Balancer:    &kafka.LeastBytes{},
		Logger:      klog.KafkaLogAdapter{level.Debug(k.logger)},
		ErrorLogger: klog.KafkaLogAdapter{level.Warn(k.logger)},
	}
	if k.tracer != nil {
		writer.Transport = NewTransport(kafka.DefaultTransport, k.tracer)
	}
	k.cache[topic] = writer
	k.closers = append(k.closers, writer.Close)
	return writer
}

func (k *KafkaProducerFactory) Close() error {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	for _, v := range k.closers {
		err := v()
		if err != nil {
			return err
		}
	}
	return nil
}

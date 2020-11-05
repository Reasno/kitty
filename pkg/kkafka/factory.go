package kkafka

import (
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
)

type KafkaProducerFactory struct {
	tracer opentracing.Tracer
	mutex sync.Mutex
	cache   map[string]*kafka.Writer
	brokers []string
	closers []func() error
}

func NewKafkaProducerFactory(brokers []string) *KafkaProducerFactory {
	return &KafkaProducerFactory{
		nil,
		sync.Mutex{},
		map[string]*kafka.Writer{},
		brokers,
		[]func() error{},
	}
}

func NewKafkaProducerFactoryWithTracer(brokers []string, tracer opentracing.Tracer) *KafkaProducerFactory {
	return &KafkaProducerFactory{
		tracer,
		sync.Mutex{},
		map[string]*kafka.Writer{},
		brokers,
		[]func() error{},
	}
}

func (k *KafkaProducerFactory) Writer(topic string) *kafka.Writer {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	if w, ok := k.cache[topic]; ok {
		return w
	}
	writer := &kafka.Writer{
		Addr:         kafka.TCP(k.brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
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

package kkafka

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"glab.tagtic.cn/ad_gains/kitty/pkg/klog"
)

type KafkaFactory struct {
	tracer  opentracing.Tracer
	mutex   sync.Mutex
	brokers []string
	closers []func() error
	logger  log.Logger
}

func NewKafkaProducerFactory(brokers []string, logger log.Logger) *KafkaFactory {
	return &KafkaFactory{
		brokers: brokers,
		closers: []func() error{},
		logger:  logger,
	}
}

func NewKafkaProducerFactoryWithTracer(brokers []string, logger log.Logger, tracer opentracing.Tracer) *KafkaFactory {
	return &KafkaFactory{
		tracer:  tracer,
		brokers: brokers,
		closers: []func() error{},
		logger:  logger,
	}
}

func (k *KafkaFactory) Writer(topic string) Publisher {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	writer := &kafka.Writer{
		Addr:        kafka.TCP(k.brokers...),
		Topic:       topic,
		Balancer:    &kafka.LeastBytes{},
		Logger:      klog.KafkaLogAdapter{Logging: level.Debug(k.logger)},
		ErrorLogger: klog.KafkaLogAdapter{Logging: level.Warn(k.logger)},
	}

	k.closers = append(k.closers, writer.Close)
	if k.tracer != nil {
		writer.Transport = NewTransport(kafka.DefaultTransport, k.tracer, topic)
	}
	return &pub{
		Writer: writer,
		topic:  topic,
		tracer: k.tracer,
		opName: "kafka.publish",
	}
}

type readerConfig struct {
	groupId     string
	parallelism int
}

type readerOpt func(config *readerConfig)

func WithGroup(group string) readerOpt {
	return func(config *readerConfig) {
		config.groupId = group
	}
}

func WithParallelism(parallelism int) readerOpt {
	return func(config *readerConfig) {
		config.parallelism = parallelism
	}
}

func (k *KafkaFactory) Reader(topic string, opt ...readerOpt) Subscriber {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	var config = readerConfig{
		groupId:     "",
		parallelism: 1,
	}
	for _, o := range opt {
		o(&config)
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     k.brokers,
		Topic:       topic,
		GroupID:     config.groupId,
		Logger:      klog.KafkaLogAdapter{Logging: level.Debug(k.logger)},
		ErrorLogger: klog.KafkaLogAdapter{Logging: level.Warn(k.logger)},
	})

	k.closers = append(k.closers, reader.Close)

	snk := wrap(reader, k.tracer, config.parallelism)
	return snk
}

func wrap(reader *kafka.Reader, tracer opentracing.Tracer, parallelism int) *sub {
	return &sub{
		Reader:      reader,
		tracer:      tracer,
		opName:      "kafka",
		parallelism: parallelism,
	}
}

func (k *KafkaFactory) Close() error {
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

type Subscriber interface {
	Subscribe(ctx context.Context, fn HandleFunc) error
}

type Publisher interface {
	Publish(ctx context.Context, msg ...kafka.Message) error
}

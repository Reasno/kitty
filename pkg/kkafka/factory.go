package kkafka

import (
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

func NewKafkaFactory(brokers []string, logger log.Logger, tracer opentracing.Tracer) *KafkaFactory {
	return &KafkaFactory{
		tracer:  tracer,
		brokers: brokers,
		closers: []func() error{},
		logger:  logger,
	}
}

func (k *KafkaFactory) MakeHandler(topic string) Handler {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	writer := &kafka.Writer{
		Addr:        kafka.TCP(k.brokers...),
		Topic:       topic,
		Balancer:    &kafka.LeastBytes{},
		Logger:      klog.KafkaLogAdapter{Logging: level.Debug(k.logger)},
		ErrorLogger: klog.KafkaLogAdapter{Logging: level.Warn(k.logger)},
		BatchSize:   1,
	}

	k.closers = append(k.closers, writer.Close)
	if k.tracer != nil {
		writer.Transport = NewTransport(kafka.DefaultTransport, k.tracer, topic)
	}
	return &pub{
		Writer: writer,
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

func (k *KafkaFactory) MakeSub(topic string, handler Handler, opt ...readerOpt) *Subscriber {
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

	return &Subscriber{
		reader:      reader,
		handler:     handler,
		parallelism: config.parallelism,
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

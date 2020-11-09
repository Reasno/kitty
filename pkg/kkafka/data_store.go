package kkafka

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

type DataStore struct {
	Factory *KafkaProducerFactory
	Topic   string
	once    sync.Once
	writer  *kafka.Writer
}

func (e *DataStore) Emit(ctx context.Context, marshaller contract.Marshaller) error {
	e.once.Do(func() {
		e.writer = e.Factory.Writer(e.Topic)
	})
	b, err := marshaller.Marshal()
	if err != nil {
		return errors.Wrap(err, "unable to marshal pb")
	}
	return e.writer.WriteMessages(ctx, kafka.Message{
		Value: b,
	})
}

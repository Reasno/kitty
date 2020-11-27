package kkafka

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

type DataStore struct {
	Factory   *KafkaFactory
	Topic     string
	once      sync.Once
	publisher Publisher
}

func (e *DataStore) Emit(ctx context.Context, marshaller contract.Marshaller) error {
	e.once.Do(func() {
		e.publisher = e.Factory.Writer(e.Topic)
	})
	b, err := marshaller.Marshal()
	if err != nil {
		return errors.Wrap(err, "unable to marshal pb")
	}
	return e.publisher.Publish(ctx, kafka.Message{
		Value: b,
	})
}

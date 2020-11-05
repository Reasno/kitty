package kkafka

import (
	"context"

	"github.com/Reasno/kitty/pkg/contract"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type DataStore struct {
	Factory *KafkaProducerFactory
	Topic   string
}

func (e *DataStore) Emit(ctx context.Context, marshaller contract.Marshaller) error {
	b, err := marshaller.Marshal()
	if err != nil {
		return errors.Wrap(err, "unable to marshal pb")
	}
	return e.Factory.Writer(e.Topic).WriteMessages(ctx, kafka.Message{
		Value: b,
	})
}

package svc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/segmentio/kafka-go"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kkafka"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

func DecodeTaskRequest(ctx context.Context, msg *kafka.Message) (interface{}, error) {
	var taskEvent pb.TaskEvent
	err := taskEvent.Unmarshal(msg.Value)
	if err != nil {
		return nil, err
	}
	return &taskEvent, nil
}

func DecodeSignRequest(ctx context.Context, msg *kafka.Message) (interface{}, error) {
	var signEvent pb.SignEvent
	err := signEvent.Unmarshal(msg.Value)
	if err != nil {
		return nil, err
	}
	return &signEvent, nil
}

func provideTaskSubscriber(endpoint endpoint.Endpoint, options ...kkafka.SubscriberOption) kkafka.Handler {
	return kkafka.NewSubscriber(
		endpoint,
		DecodeTaskRequest,
		options...,
	)
}

func provideSignSubscriber(endpoint endpoint.Endpoint, options ...kkafka.SubscriberOption) kkafka.Handler {
	return kkafka.NewSubscriber(
		endpoint,
		DecodeSignRequest,
		options...,
	)
}

func MakeKafkaServer(endpoints Endpoints, factory *kkafka.KafkaFactory, conf contract.ConfigReader, options ...kkafka.SubscriberOption) kkafka.Server {
	group := conf.String("kafka.groupId")

	task := provideTaskSubscriber(endpoints.PushTaskEventEndpoint, options...)
	sign := provideSignSubscriber(endpoints.PushSignEventEndpoint, options...)

	return kkafka.NewMux(
		factory.MakeKafkaServer(conf.String("kafka.taskEventBus"), task, kkafka.WithGroup(group)),
		factory.MakeKafkaServer(conf.String("kafka.signEventBus"), sign, kkafka.WithGroup(group)),
	)
}

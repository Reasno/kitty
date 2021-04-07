package svc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/segmentio/kafka-go"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kkafka"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

func DecodeBindAdRequest(ctx context.Context, msg *kafka.Message) (interface{}, error) {
	var UserBindAd pb.UserBindAdRequest
	err := UserBindAd.Unmarshal(msg.Value)
	if err != nil {
		return nil, err
	}
	return &UserBindAd, nil
}

func provideBindAdSubscriber(endpoint endpoint.Endpoint, options ...kkafka.SubscriberOption) kkafka.Handler {
	return kkafka.NewSubscriber(
		endpoint,
		DecodeBindAdRequest,
		options...,
	)
}

func MakeKafkaServer(endpoints Endpoints, factory *kkafka.KafkaFactory, conf contract.ConfigReader, options ...kkafka.SubscriberOption) kkafka.Server {
	group := conf.String("kafka.groupId")

	sign := provideBindAdSubscriber(endpoints.BindAdEndpoint, options...)

	return kkafka.NewMux(
		factory.MakeKafkaServer(conf.String("kafka.bindAd"), sign, kkafka.WithGroup(group)),
	)
}

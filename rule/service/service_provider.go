package service

import (
	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

func ProvideService(logger log.Logger, repo Repository, dmpServer pb.DmpServer, redisClient redis.UniversalClient) Service {
	return &service{logger: logger, repo: repo, dmpServer: dmpServer, redisClient: redisClient}
}

package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

type UniqueID struct {
	redis redis.UniversalClient
	key   string
}

func NewUniqueID(redis redis.UniversalClient, conf contract.ConfigReader) *UniqueID {
	return &UniqueID{
		redis: redis,
		key:   conf.String("incrKey"),
	}
}

func (u *UniqueID) ID(ctx context.Context) (uint, error) {
	id, err := u.redis.Incr(ctx, u.key).Uint64()
	if err != nil {
		return 0, errors.Wrap(err, "cannot create id from redis")
	}
	return uint(id), err
}

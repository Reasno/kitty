package repository

import (
	"context"
	"github.com/Reasno/kitty/pkg/otredis"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"time"
)

const key = "CodeRepo"

type CodeRepo struct {
	client redis.Cmdable
	km     *otredis.KeyManager
}

func NewCodeRepo(cmdable redis.Cmdable, km *otredis.KeyManager) *CodeRepo {
	km.Add(key)
	return &CodeRepo{cmdable, km}
}

func (c *CodeRepo) AddCode(ctx context.Context, mobile string) (code string, err error) {
	n := rand.Intn(1_000_000)
	code = pad(n)
	_, err = c.client.Set(ctx, c.km.Key(mobile), code, 15*time.Minute).Result()
	if err != nil {
		return "", errors.Wrap(err, "cannot persist code in redis")
	}
	return code, nil
}

func (c *CodeRepo) CheckCode(ctx context.Context, mobile, code string) (bool, error) {
	value, err := c.client.Get(ctx, c.km.Key(mobile)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "cannot query code in redis")
	}
	return value == code, nil
}

func (c *CodeRepo) DeleteCode(ctx context.Context, mobile string) (err error) {
	_, err = c.client.Del(ctx, c.km.Key(mobile)).Result()
	return err
}

func pad(n int) string {
	s := strconv.Itoa(n)
	for len(s) < 6 {
		s = "0" + s
	}
	return s
}

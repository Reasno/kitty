package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/rule/entity"
	"glab.tagtic.cn/ad_gains/kitty/rule/module"
	"glab.tagtic.cn/ad_gains/kitty/rule/msg"
	repository2 "glab.tagtic.cn/ad_gains/kitty/rule/repository"
	"go.etcd.io/etcd/clientv3"
)

// repository 专门为客户端提供的 repository，不具备自举性，可以只watch需要的规则
type repository struct {
	client     *clientv3.Client
	logger     log.Logger
	containers map[string]repository2.Container
	prefix     string
	rwLock     sync.RWMutex
}

func NewRepository(client *clientv3.Client, logger log.Logger, activeContainers map[string]string) (*repository, error) {
	var (
		err   error
		value []byte
	)

	var repo = &repository{
		client:     client,
		logger:     logger,
		containers: make(map[string]repository2.Container),
		prefix:     "",
		rwLock:     sync.RWMutex{},
	}

	// 填充所有容器
	for k, v := range activeContainers {
		repo.containers[k] = repository2.Container{DbKey: v, Name: k, RuleSet: []entity.Rule{}}
	}

	// 依次拉取规则
	var count = 0
	for k, v := range repo.containers {
		count++
		value, err = repo.getRawRuleSetFromDbKey(context.Background(), v.DbKey)
		if err != nil {
			level.Warn(logger).Log("err", errors.Wrap(err, fmt.Sprintf(msg.ErrorInvalidConfig, v.Name)))
			value = []byte("{}")
		}
		v.RuleSet = entity.NewRules(bytes.NewReader(value), logger)
		repo.containers[k] = v
	}

	level.Info(logger).Log("msg", fmt.Sprintf("%d rules have been added", count))

	// 自动搜索共同前缀
	repo.prefix = module.Prefix(module.DbKeys(repo.containers))

	return repo, nil
}

func (r *repository) updateRuleSetByDbKey(dbKey string, rules []entity.Rule) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	for i, v := range r.containers {
		if dbKey == r.containers[i].DbKey {
			v.RuleSet = rules
			r.containers[i] = v
		}
	}
}

func (r *repository) WatchConfigUpdate(ctx context.Context) error {
	level.Info(r.logger).Log("msg", "listening to etcd changes: "+strings.Join(r.client.Endpoints(), ","))
	rch := r.client.Watch(ctx, r.prefix, clientv3.WithPrefix())
	for {
		select {
		case wresp := <-rch:
			if wresp.Err() != nil {
				return wresp.Err()
			}
			for _, ev := range wresp.Events {
				rules := entity.NewRules(bytes.NewReader(ev.Kv.Value), r.logger)
				r.updateRuleSetByDbKey(string(ev.Kv.Key), rules)
				level.Info(r.logger).Log("msg", fmt.Sprintf("配置已更新 %+v", ev.Kv))
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *repository) GetRaw(ctx context.Context, name string) (value []byte, e error) {
	panic("not implemented")
}

func (r *repository) getRawRuleSetFromDbKey(ctx context.Context, dbKey string) (value []byte, e error) {
	resp, err := r.client.Get(ctx, dbKey)
	if err != nil {
		return nil, errors.Wrapf(err, msg.ErrorGetKeyFromETCD, dbKey)
	}
	for _, ev := range resp.Kvs {
		return ev.Value, nil
	}
	return nil, err
}

func (r *repository) SetRaw(ctx context.Context, name string, value string) error {
	panic("not implemented")
}

// IsNewest 传入的内容是否和ETCD中的最新版本一致
func (r *repository) IsNewest(ctx context.Context, key, value string) (bool, error) {
	panic("not implemented")
}

func (r *repository) GetCompiled(ruleName string) []entity.Rule {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()
	if c, ok := r.containers[ruleName]; ok {
		return c.RuleSet
	}
	panic(fmt.Sprintf("unregistered rule %s", ruleName))
}

func (r *repository) ValidateRules(ruleName string, reader io.Reader) error {
	panic("implement me")
}

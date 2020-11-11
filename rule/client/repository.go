package client

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"strings"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/rule"
	"glab.tagtic.cn/ad_gains/kitty/rule/msg"
	"go.etcd.io/etcd/clientv3"
	"gopkg.in/yaml.v3"
)

// repository 专门为客户端提供的 repository，不具备自举性，可以只watch需要的规则
type repository struct {
	client     *clientv3.Client
	logger     log.Logger
	containers map[string]rule.Container
	prefix     string
	mapping    map[string]string
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
		containers: make(map[string]rule.Container),
		prefix:     "",
		mapping:    activeContainers,
		rwLock:     sync.RWMutex{},
	}

	// 填充所有容器
	for k, v := range activeContainers {
		repo.containers[k] = rule.Container{DbKey: v, Name: k, RuleSet: []rule.Rule{}}
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
		v.RuleSet = rule.NewRules(bytes.NewReader(value), logger)
		repo.containers[k] = v
	}

	level.Info(logger).Log("msg", fmt.Sprintf("%d rules have been added", count))

	// 自动搜索共同前缀
	repo.prefix = rule.Prefix(rule.DbKeys(repo.containers))

	return repo, nil
}

func (r *repository) updateRuleSetByDbKey(dbKey string, rules []rule.Rule) {
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
				rules := rule.NewRules(bytes.NewReader(ev.Kv.Value), r.logger)
				r.updateRuleSetByDbKey(string(ev.Kv.Key), rules)
				level.Info(r.logger).Log("msg", fmt.Sprintf("配置已更新 %+v", ev.Kv))
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *repository) GetRaw(ctx context.Context, name string) (value []byte, e error) {
	dbKey, ok := r.mapping[name]
	if !ok {
		return nil, fmt.Errorf("unknown rule set %s", name)
	}
	return r.getRawRuleSetFromDbKey(ctx, dbKey)
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
	dbKey, ok := r.mapping[name]
	if !ok {
		return fmt.Errorf("unknown rule set %s", name)
	}
	return r.setRawRuleSetFromDbKey(ctx, dbKey, value)
}

func (r *repository) setRawRuleSetFromDbKey(ctx context.Context, dbKey string, value string) error {
	_, err := r.client.Put(ctx, dbKey, value)
	return err
}

// getRev 返回目标值的版本号
func (r *repository) getRev(ctx context.Context, key string) (int64, error) {
	resp, err := r.client.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return resp.Header.Revision, nil
}

func (r *repository) validateAllRules() error {
	todo := rule.DbKeys(r.containers)
	for _, v := range todo {
		var tmp []rule.Rule
		value, err := r.getRawRuleSetFromDbKey(context.Background(), v)
		if err != nil {
			return errors.Wrapf(err, msg.ErrorInvalidConfig, v)
		}
		err = yaml.Unmarshal(value, &tmp)
		if err != nil {
			return errors.Wrapf(err, msg.ErrorInvalidConfig, v)
		}
		for i := range tmp {
			if err := tmp[i].Compile(); err != nil {
				return errors.Wrapf(err, msg.ErrorInvalidConfig, tmp[i].If)
			}
		}
	}
	return nil
}

// IsNewest 传入的内容是否和ETCD中的最新版本一致
func (r *repository) IsNewest(ctx context.Context, key, value string) (bool, error) {
	v, err := r.client.Get(ctx, key)
	if err != nil {
		return false, err
	}
	for _, ev := range v.Kvs {
		return getMd5(ev.Value) == value, nil
	}
	return true, nil
}

func (r *repository) GetCompiled(ruleName string) []rule.Rule {
	r.rwLock.RLock()
	defer r.rwLock.Unlock()
	if c, ok := r.containers[ruleName]; ok {
		return c.RuleSet
	}
	return []rule.Rule{}
}

func getMd5(orig []byte) string {
	m := md5.New()
	m.Write(orig)
	return fmt.Sprintf("%x", m.Sum(nil))
}

package client

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/rule/entity"
	"glab.tagtic.cn/ad_gains/kitty/rule/module"
	"glab.tagtic.cn/ad_gains/kitty/rule/msg"
	repository2 "glab.tagtic.cn/ad_gains/kitty/rule/repository"
	"go.etcd.io/etcd/client/v3"
)

// repository 专门为客户端提供的 repository，不具备自举性，可以只watch需要的规则
type repository struct {
	client     *clientv3.Client
	logger     log.Logger
	containers map[string]repository2.Container
	prefix     string
	rwLock     sync.RWMutex
	regexp     *regexp.Regexp
	limit      int64
}

type RepositoryConfig struct {
	Prefix      string
	Regex       *regexp.Regexp
	ListOfRules []string
	Limit       int64
}

func NewRepositoryWithConfig(client *clientv3.Client, logger log.Logger, config RepositoryConfig) (*repository, error) {

	var repo = &repository{
		client:     client,
		logger:     logger,
		containers: make(map[string]repository2.Container),
		prefix:     repository2.OtherConfigPathPrefix,
		rwLock:     sync.RWMutex{},
		regexp:     config.Regex,
		limit:      1000,
	}

	if config.Limit != 0 {
		repo.limit = config.Limit
	}

	if config.Prefix != "" {
		repo.prefix = repository2.OtherConfigPathPrefix + "/" + config.Prefix
	}

	// 自动搜索共同前缀
	if len(config.ListOfRules) != 0 {
		repo.prefix = repository2.OtherConfigPathPrefix + "/" + module.Prefix(config.ListOfRules)
	}

	// 第一次拉取配置
	configMap, err := repo.getRawRuleSetsFromPrefix(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorRules)
	}

	var count = 0
	for k, v := range configMap {
		name := dbKeyToName(k)
		if config.Regex != nil && !config.Regex.MatchString(name) {
			continue
		}
		count++
		c := repository2.Container{DbKey: k, Name: name, RuleSet: nil}
		c.RuleSet = entity.NewRules(bytes.NewReader(v), logger)
		repo.containers[dbKeyToName(k)] = c
	}

	level.Info(logger).Log("msg", fmt.Sprintf("%d rules have been added", count))

	return repo, nil
}

// NewRepository creates a new repo
// Deprecated: don't use
func NewRepository(client *clientv3.Client, logger log.Logger, activeContainers map[string]string) (*repository, error) {

	var repo = &repository{
		client:     client,
		logger:     logger,
		containers: make(map[string]repository2.Container),
		prefix:     repository2.OtherConfigPathPrefix,
		rwLock:     sync.RWMutex{},
	}

	// 填充所有容器
	for k, v := range activeContainers {
		repo.containers[k] = repository2.Container{DbKey: v, Name: k, RuleSet: nil}
	}

	// 自动搜索共同前缀
	if len(repo.containers) != 0 {
		repo.prefix = module.Prefix(module.DbKeys(repo.containers))
	}

	// 第一次拉取配置
	configMap, err := repo.getRawRuleSetsFromPrefix(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorRules)
	}

	var count = 0
	for k, v := range configMap {
		count++
		c := repository2.Container{DbKey: k, Name: dbKeyToName(k), RuleSet: nil}
		c.RuleSet = entity.NewRules(bytes.NewReader(v), logger)
		repo.containers[dbKeyToName(k)] = c
	}

	level.Info(logger).Log("msg", fmt.Sprintf("%d rules have been added", count))

	return repo, nil
}

func (r *repository) updateRuleSetByDbKey(dbKey string, rules entity.Ruler) {
	name := dbKeyToName(dbKey)

	if r.regexp != nil && r.regexp.MatchString(name) {
		c := repository2.Container{DbKey: dbKey, Name: name, RuleSet: rules}

		r.rwLock.Lock()
		r.containers[name] = c
		r.rwLock.Unlock()
		return
	}

	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	if v, ok := r.containers[name]; ok {
		v.RuleSet = rules

		r.containers[name] = v
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

func (r *repository) getRawRuleSetsFromPrefix(ctx context.Context) (value map[string][]byte, e error) {
	value = make(map[string][]byte)
	key := r.prefix
	for {
		resp, err := r.client.Get(ctx, key, clientv3.WithRange(clientv3.GetPrefixRangeEnd(r.prefix)), clientv3.WithLimit(r.limit))
		if err != nil {
			return nil, errors.Wrapf(err, "prefix not found %s", r.prefix)
		}
		for _, ev := range resp.Kvs {
			value[string(ev.Key)] = ev.Value
		}
		if !resp.More {
			return value, err
		}
		// move to next key
		key = string(append(resp.Kvs[len(resp.Kvs)-1].Key, 0))
	}

}

func (r *repository) GetCompiled(ruleName string) entity.Ruler {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()
	if c, ok := r.containers[ruleName]; ok {
		return c.RuleSet
	}
	return nil
}

func dbKeyToName(dbKey string) string {
	return strings.Replace(dbKey, repository2.OtherConfigPathPrefix+"/", "", 1)
}

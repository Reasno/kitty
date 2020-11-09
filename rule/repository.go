//go:generate mockery --name=Repository
package rule

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"strings"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/knadh/koanf"
	kyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/rule/msg"
	"go.etcd.io/etcd/clientv3"
	"gopkg.in/yaml.v3"
)

type repository struct {
	client     *clientv3.Client
	logger     log.Logger
	containers map[string]Container
	prefix     string
	mapping    map[string]string
	rwLock     sync.RWMutex
}

type Container struct {
	RuleSet []Rule
	DbKey   string
	Name    string
}

const CentralConfigPath = "/central-config"
const OtherConfigPathPrefix = "/monetization"

func NewRepository(client *clientv3.Client, logger log.Logger) (*repository, error) {

	var repo = &repository{
		client:     client,
		logger:     logger,
		containers: make(map[string]Container),
		prefix:     "",
		mapping:    nil,
		rwLock:     sync.RWMutex{},
	}

	// 读取配置中心
	activeContainers, err := repo.readCentralConfig()
	if err != nil {
		return nil, err
	}

	repo.resetActiveContainers(activeContainers)

	return repo, nil
}

func (r *repository) updateRuleSetByDbKey(dbKey string, rules []Rule) {
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
	centerCh := r.client.Watch(ctx, CentralConfigPath)
	rch := r.client.Watch(ctx, r.prefix, clientv3.WithPrefix())
	for {
		select {
		case cresp := <-centerCh:
			if cresp.Err() != nil {
				return cresp.Err()
			}
			for _, ev := range cresp.Events {
				activeContainers, err := r.readCentralConfig()
				if err != nil {
					return err
				}
				r.resetActiveContainers(activeContainers)
				level.Info(r.logger).Log("msg", fmt.Sprintf("中心配置已更新 %+v", ev.Kv))
			}
		case wresp := <-rch:
			if wresp.Err() != nil {
				return wresp.Err()
			}
			for _, ev := range wresp.Events {
				rules := NewRules(bytes.NewReader(ev.Kv.Value), r.logger)
				r.updateRuleSetByDbKey(string(ev.Kv.Key), rules)
				level.Info(r.logger).Log("msg", fmt.Sprintf("配置已更新 %+v", ev.Kv))
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *repository) readCentralConfig() (map[string]string, error) {

	b, err := r.getRawRuleSetFromDbKey(context.Background(), CentralConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get central config from repository")
	}

	var centralRules CentralRules
	c := koanf.New(".")
	err = c.Load(rawbytes.Provider(b), kyaml.Parser())
	if err != nil {
		return nil, errors.Wrap(err, "Unable to load central config")
	}

	err = c.Unmarshal("", &centralRules)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to parse central config")
	}

	var activeContainers = make(map[string]string)
	for _, v := range centralRules.Rule.List {
		collect(activeContainers, v.Path, r.prefix)
		for _, v := range v.Children {
			collect(activeContainers, v.Path, r.prefix)
		}
	}
	activeContainers["central-config"] = CentralConfigPath
	return activeContainers, nil
}

func collect(containers map[string]string, path string, p string) {
	containers[path[1:]+"-prod"] = p + path + "-prod"
	containers[path[1:]+"-dev"] = p + path + "-dev"
	containers[path[1:]+"-testing"] = p + path + "-testing"
	containers[path[1:]+"-local"] = p + path + "-local"
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
	todo := dbKeys(r.containers)
	for _, v := range todo {
		var tmp []Rule
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

func (r *repository) GetCompiled(ruleName string) []Rule {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()
	if c, ok := r.containers[ruleName]; ok {
		return c.RuleSet
	}
	return []Rule{}
}

func (r *repository) resetActiveContainers(activeContainers map[string]string) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	r.mapping = activeContainers

	// 填充所有容器
	for k, v := range activeContainers {
		r.containers[k] = Container{DbKey: v, Name: k, RuleSet: []Rule{}}
	}

	// 依次拉取规则
	var count = 0
	for k, v := range r.containers {
		count++
		value, err := r.getRawRuleSetFromDbKey(context.Background(), v.DbKey)
		if err != nil {
			level.Warn(r.logger).Log("err", errors.Wrap(err, fmt.Sprintf(msg.ErrorInvalidConfig, v.Name)))
			value = []byte("{}")
		}
		v.RuleSet = NewRules(bytes.NewReader(value), r.logger)
		r.containers[k] = v
	}
	level.Info(r.logger).Log("msg", fmt.Sprintf("%d rules have been added", count))

	// 自动搜索共同前缀
	r.prefix = OtherConfigPathPrefix
}

func getMd5(orig []byte) string {
	m := md5.New()
	m.Write(orig)
	return fmt.Sprintf("%x", m.Sum(nil))
}

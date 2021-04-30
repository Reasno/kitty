package repository

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/knadh/koanf"
	kyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/rule/entity"
	"glab.tagtic.cn/ad_gains/kitty/rule/msg"
	"go.etcd.io/etcd/clientv3"
)

type repository struct {
	client         *clientv3.Client
	logger         log.Logger
	containers     map[string]Container
	rwLock         sync.RWMutex
	updateChan     chan struct{}
	watchReadyChan chan struct{}
}

type Container struct {
	RuleSet entity.Ruler
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
		rwLock:     sync.RWMutex{},
		updateChan: nil,
	}

	// 读取配置中心
	activeContainers, err := repo.readCentralConfig()
	if err != nil {
		return nil, err
	}

	repo.resetActiveContainers(activeContainers)

	return repo, nil
}

func (r *repository) updateRuleSetByDbKey(dbKey string, rules entity.Ruler) {
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
	rch := r.client.Watch(ctx, OtherConfigPathPrefix, clientv3.WithPrefix())

	if r.watchReadyChan != nil {
		r.watchReadyChan <- struct{}{}
	}

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
			if r.updateChan != nil {
				r.updateChan <- struct{}{}
			}
		case wresp := <-rch:
			if wresp.Err() != nil {
				return wresp.Err()
			}
			for _, ev := range wresp.Events {
				rules := entity.NewRules(bytes.NewReader(ev.Kv.Value), r.logger)
				r.updateRuleSetByDbKey(string(ev.Kv.Key), rules)
				level.Info(r.logger).Log("msg", fmt.Sprintf("配置已更新 %+v", ev.Kv))
			}
			if r.updateChan != nil {
				r.updateChan <- struct{}{}
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

	return readCentralConfigBytes(b)
}

func readCentralConfigBytes(b []byte) (map[string]string, error) {

	var (
		err          error
		centralRules entity.CentralRules
	)
	c := koanf.New(".")
	err = c.Load(rawbytes.Provider(b), kyaml.Parser())
	if err != nil {
		return nil, errors.Wrap(err, "unable to load central config")
	}

	err = c.Unmarshal("", &centralRules)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse central config")
	}
	if len(centralRules.Rule.List) == 0 {
		return nil, errors.New("failed to unmarshal central rule")
	}

	var activeContainers = make(map[string]string)
	for _, v := range centralRules.Rule.List {
		if strings.Count(v.Path, "/") >= 2 {
			return nil, errors.Wrapf(err, "subpath are not allowed: %s", v.Path)
		}
		collect(activeContainers, v.Path, v.Tabs, OtherConfigPathPrefix)
		for _, v := range v.Children {
			if strings.Count(v.Path, "/") >= 2 {
				return nil, errors.Wrapf(err, "subpath are not allowed: %s", v.Path)
			}
			collect(activeContainers, v.Path, v.Tabs, OtherConfigPathPrefix)
		}
	}
	return activeContainers, nil
}

func collect(containers map[string]string, path string, tabs []string, p string) {
	if len(tabs) == 0 {
		tabs = []string{"prod", "dev", "testing", "local"}
	}
	if len(path) > 1 && path[0] == '/' {
		for _, t := range tabs {
			containers[path[1:]+"-"+t] = p + path + "-" + t
		}
	}
}

func (r *repository) GetRaw(ctx context.Context, name string) (value []byte, e error) {
	c, ok := r.containers[name]
	if !ok {
		return nil, fmt.Errorf("unknown rule set %s", name)
	}
	return r.getRawRuleSetFromDbKey(ctx, c.DbKey)
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
	c, ok := r.containers[name]
	if !ok {
		return fmt.Errorf("unknown rule set %s", name)
	}
	return r.setRawRuleSetFromDbKey(ctx, c.DbKey, value)
}

func (r *repository) setRawRuleSetFromDbKey(ctx context.Context, dbKey string, value string) error {
	_, err := r.client.Put(ctx, dbKey, value)
	return err
}

// getRev 返回目标值的版本号
//func (r *repository) getRev(ctx context.Context, key string) (int64, error) {
//	resp, err := r.client.Get(ctx, key)
//	if err != nil {
//		return 0, err
//	}
//	return resp.Header.Revision, nil
//}

func (r *repository) ValidateRules(ruleName string, reader io.Reader) error {
	if ruleName == CentralConfigPath[1:] {
		byt, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		_, err = readCentralConfigBytes(byt)
		return err
	}
	return entity.ValidateRules(reader)
}

// IsNewest 传入的内容是否和ETCD中的最新版本一致
func (r *repository) IsNewest(ctx context.Context, key, value string) (bool, error) {
	c, ok := r.containers[key]
	if !ok {
		return false, fmt.Errorf("unknown rule set %s", key)
	}

	v, err := r.client.Get(ctx, c.DbKey)
	if err != nil {
		return false, err
	}
	for _, ev := range v.Kvs {
		return getMd5(ev.Value) == value, nil
	}
	return true, nil
}

func (r *repository) GetCompiled(ruleName string) entity.Ruler {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()
	if c, ok := r.containers[ruleName]; ok {
		return c.RuleSet
	}
	return nil
}

func (r *repository) resetActiveContainers(activeContainers map[string]string) {

	// 更新容器
	newContainers := make(map[string]Container, len(activeContainers)+1)
	for k, v := range activeContainers {
		newContainers[k] = Container{DbKey: v, Name: k, RuleSet: nil}
	}
	newContainers["central-config"] = Container{DbKey: CentralConfigPath, Name: "central-config", RuleSet: nil}

	// 依次拉取规则
	var count = 0
	for k, v := range newContainers {
		count++
		value, err := r.getRawRuleSetFromDbKey(context.Background(), v.DbKey)
		if err != nil {
			level.Warn(r.logger).Log("err", errors.Wrap(err, fmt.Sprintf(msg.ErrorInvalidConfig, v.Name)))
			value = []byte("{}")
		}
		v.RuleSet = entity.NewRules(bytes.NewReader(value), r.logger)
		newContainers[k] = v
	}

	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	r.containers = newContainers
	level.Info(r.logger).Log("msg", fmt.Sprintf("%d rules have been added", count))
}

func getMd5(orig []byte) string {
	m := md5.New()
	m.Write(orig)
	return fmt.Sprintf("%x", m.Sum(nil))
}

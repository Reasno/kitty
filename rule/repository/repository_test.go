package repository

import (
	"context"
	"flag"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/client/v3"
)

var useEtcd bool

func init() {
	flag.BoolVar(&useEtcd, "etcd", false, "use local mysql for testing")
}

func readFiles(name string) string {
	byt, _ := ioutil.ReadFile("testdata/" + name)
	return string(byt)
}

func TestRepository_WatchConfigUpdate(t *testing.T) {
	if !useEtcd {
		t.Skip("test dynamic config requires etcd")
	}
	var (
		foobar                 = readFiles("foobar")
		foobaz                 = readFiles("foobaz")
		fooqux                 = readFiles("fooqux")
		configCentralManyLines = readFiles("config_central_many_lines")
		configCentralFewLines  = readFiles("config_central_few_lines")
	)
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		t.Fatal(err)
	}
	client.Delete(context.Background(), OtherConfigPathPrefix+"/kitty-testing")
	client.Delete(context.Background(), OtherConfigPathPrefix+"/egg-testing")
	client.Delete(context.Background(), CentralConfigPath)
	client.Put(context.Background(), CentralConfigPath, configCentralFewLines)
	client.Put(context.Background(), OtherConfigPathPrefix+"/kitty-testing", foobar)
	client.Put(context.Background(), OtherConfigPathPrefix+"/egg-testing", fooqux)
	repo, err := NewRepository(client, log.NewNopLogger())
	if err != nil {
		t.Fatal(err)
	}
	repo.updateChan = make(chan struct{})
	repo.watchReadyChan = make(chan struct{})
	watchCxt, cancel := context.WithCancel(context.Background())
	defer cancel()
	go repo.WatchConfigUpdate(watchCxt)

	type caseList []struct {
		name    string
		dbKey   string
		dataFoo interface{}
	}

	cases := caseList{
		{
			"central-config",
			CentralConfigPath,
			nil,
		},
		{
			"kitty-testing",
			OtherConfigPathPrefix + "/kitty-testing",
			"bar",
		},
		{
			"egg-testing",
			OtherConfigPathPrefix + "/egg-testing",
			"n/a",
		},
	}
	for _, c := range cases {
		if _, ok := repo.containers[c.name]; ok {
			assert.Equal(t, repo.containers[c.name].DbKey, c.dbKey)
			assert.Equal(t, repo.containers[c.name].Name, c.name)
			//assert.Equal(t, repo.containers[c.name].RuleSet.(*entity.AdvancedRuleCollection).items[0].then["foo"], c.dataFoo)
			continue
		}
		t.Fail()
	}

	// 等待watch准备就绪后再继续测试
	<-repo.watchReadyChan

	client.Put(context.Background(), OtherConfigPathPrefix+"/kitty-testing", foobaz)
	<-repo.updateChan

	cases = caseList{
		{
			"central-config",
			CentralConfigPath,
			nil,
		},
		{
			"kitty-testing",
			OtherConfigPathPrefix + "/kitty-testing",
			"baz",
		},
		{
			"egg-testing",
			OtherConfigPathPrefix + "/egg-testing",
			"n/a",
		},
	}
	for _, c := range cases {
		if _, ok := repo.containers[c.name]; ok {
			assert.Equal(t, repo.containers[c.name].DbKey, c.dbKey)
			assert.Equal(t, repo.containers[c.name].Name, c.name)
			//assert.Equal(t, repo.containers[c.name].RuleSet[0].Then["foo"], c.dataFoo)
			continue
		}
		t.Fail()
	}

	client.Put(context.Background(), CentralConfigPath, configCentralManyLines)
	<-repo.updateChan
	cases = caseList{
		{
			"central-config",
			CentralConfigPath,
			nil,
		},
		{
			"kitty-testing",
			OtherConfigPathPrefix + "/kitty-testing",
			"baz",
		},
		{
			"egg-testing",
			OtherConfigPathPrefix + "/egg-testing",
			"qux",
		},
	}
	for _, c := range cases {
		if _, ok := repo.containers[c.name]; ok {
			assert.Equal(t, repo.containers[c.name].DbKey, c.dbKey)
			assert.Equal(t, repo.containers[c.name].Name, c.name)
			continue
		}
		t.Fail()
	}
	client.Put(context.Background(), OtherConfigPathPrefix+"/egg-testing", foobar)
	<-repo.updateChan
	cases = caseList{
		{
			"central-config",
			CentralConfigPath,
			nil,
		},
		{
			"kitty-testing",
			OtherConfigPathPrefix + "/kitty-testing",
			"baz",
		},
		{
			"egg-testing",
			OtherConfigPathPrefix + "/egg-testing",
			"bar",
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.dbKey, repo.containers[c.name].DbKey)
		assert.Equal(t, c.name, repo.containers[c.name].Name)
	}

	client.Put(context.Background(), CentralConfigPath, configCentralFewLines)
	<-repo.updateChan

	cases = caseList{
		{
			"central-config",
			CentralConfigPath,
			nil,
		},
		{
			"kitty-testing",
			OtherConfigPathPrefix + "/kitty-testing",
			"baz",
		},
		{
			"egg-testing",
			OtherConfigPathPrefix + "/egg-testing",
			"n/a",
		},
	}
	for _, c := range cases {
		if _, ok := repo.containers[c.name]; ok {
			assert.Equal(t, repo.containers[c.name].DbKey, c.dbKey)
			assert.Equal(t, repo.containers[c.name].Name, c.name)
			continue
		}
		t.Fail()
	}
}

func TestRepository_IsNewest(t *testing.T) {
	if !useEtcd {
		t.Skip("test dynamic config requires etcd")
	}
	var (
		foobar = readFiles("foobar")
		foobaz = readFiles("foobaz")
	)
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		t.Fatal(err)
	}
	client.Put(context.Background(), OtherConfigPathPrefix+"/kitty-testing", foobar)
	repo, err := NewRepository(client, log.NewNopLogger())
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		data string
		ok   bool
	}{
		{
			foobar,
			true,
		},
		{
			foobaz,
			false,
		},
	}
	for _, c := range cases {
		cc := c
		t.Run("", func(t *testing.T) {
			ok, err := repo.IsNewest(context.Background(), "kitty-testing", getMd5([]byte(cc.data)))
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, ok, cc.ok)
		})
	}
}

func TestRepository_GetRaw(t *testing.T) {
	if !useEtcd {
		t.Skip("test dynamic config requires etcd")
	}
	var (
		foobar                = readFiles("foobar")
		foobaz                = readFiles("foobaz")
		configCentralFewLines = readFiles("config_central_few_lines")
	)
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		t.Fatal(err)
	}
	client.Put(context.Background(), OtherConfigPathPrefix+"/kitty-testing", foobar)
	client.Put(context.Background(), OtherConfigPathPrefix+"/kitty-local", foobaz)
	client.Put(context.Background(), CentralConfigPath, configCentralFewLines)
	repo, err := NewRepository(client, log.NewNopLogger())
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		key    string
		data   string
		hasErr bool
	}{
		{
			"kitty-testing",
			foobar,
			false,
		},
		{
			"kitty-local",
			foobaz,
			false,
		},
		{
			"whatever",
			"",
			true,
		},
	}
	for _, c := range cases {
		cc := c
		t.Run(cc.key, func(t *testing.T) {
			data, err := repo.GetRaw(context.Background(), cc.key)
			assert.Equal(t, err != nil, cc.hasErr)
			assert.Equal(t, string(data), cc.data)
		})
	}
}

func TestRepository_SetRaw(t *testing.T) {
	if !useEtcd {
		t.Skip("test dynamic config requires etcd")
	}
	var (
		foobar                = readFiles("foobar")
		configCentralFewLines = readFiles("config_central_few_lines")
	)
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		t.Fatal(err)
	}
	client.Put(context.Background(), CentralConfigPath, configCentralFewLines)
	repo, err := NewRepository(client, log.NewNopLogger())
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		key    string
		data   string
		hasErr bool
	}{
		{
			"kitty-testing",
			foobar,
			false,
		},
		{
			"whatever",
			"",
			true,
		},
	}
	for _, c := range cases {
		cc := c
		t.Run(cc.key, func(t *testing.T) {
			err := repo.SetRaw(context.Background(), cc.key, cc.data)
			assert.Equal(t, err != nil, cc.hasErr)
		})
	}
}

func TestRepository_ValidateRules(t *testing.T) {
	if !useEtcd {
		t.Skip("test dynamic config requires etcd")
	}
	var (
		foobar                = readFiles("foobar")
		configCentralFewLines = readFiles("config_central_few_lines")
		configCentralBad      = readFiles("config_central_bad")
		configCentralSubpath  = readFiles("config_central_subpath")
	)
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		t.Fatal(err)
	}

	repo, err := NewRepository(client, log.NewNopLogger())
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		key    string
		data   string
		hasErr bool
	}{
		{
			"kitty-testing",
			foobar,
			false,
		},
		{
			"central-config",
			configCentralFewLines,
			false,
		},
		{
			"central-config",
			configCentralBad,
			true,
		},
		{
			"central-config",
			configCentralSubpath,
			true,
		},
	}
	for _, c := range cases {
		cc := c
		t.Run(cc.key, func(t *testing.T) {
			err := repo.ValidateRules(cc.key, strings.NewReader(cc.data))
			assert.Equal(t, err != nil, cc.hasErr)
		})
	}
}

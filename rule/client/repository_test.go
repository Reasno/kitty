package client

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/client/v3"
)

func TestGetRawRuleSetsFromPrefix(t *testing.T) {
	if !useEtcd {
		t.Skip("test dynamic config requires etcd")
	}
	client, _ := clientv3.New(clientv3.Config{
		Endpoints: []string{"etcd-1:2379"},
	})
	repo, _ := NewRepositoryWithConfig(client, log.NewNopLogger(), RepositoryConfig{
		Prefix: "",
		Limit:  100,
	})
	v1, e := repo.getRawRuleSetsFromPrefix(context.Background())
	assert.NoError(t, e)
	assert.Less(t, 100, len(v1))

	repo, _ = NewRepositoryWithConfig(client, log.NewNopLogger(), RepositoryConfig{
		Prefix: "",
		Limit:  10000,
	})
	v2, e := repo.getRawRuleSetsFromPrefix(context.Background())
	assert.NoError(t, e)
	assert.Less(t, len(v2), 10000)

	assert.Equal(t, len(v1), len(v2))
}

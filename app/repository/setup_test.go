package repository

import (
	"context"
	"flag"
	"github.com/go-redis/redis/v8"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	mc "glab.tagtic.cn/ad_gains/kitty/pkg/contract/mocks"
	"testing"

	"github.com/go-gormigrate/gormigrate/v2"
	"glab.tagtic.cn/ad_gains/kitty/app/entity"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var m *gormigrate.Gormigrate
var db *gorm.DB

var useMysql bool

func init() {
	flag.BoolVar(&useMysql, "mysql", false, "use local mysql for testing")
}

func setUp(t *testing.T) {
	var err error
	if !useMysql {
		db, err = gorm.Open(sqlite.Open(":memory:?cache=shared"), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: "test_", // 表名前缀，`User` 的表名应该是 `test_users`
			},
		})
	} else {
		db, err = gorm.Open(mysql.Open("root@tcp(127.0.0.1:3306)/kitty?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: "test_", // 表名前缀，`User` 的表名应该是 `test_users`
			},
		})
	}

	if err != nil {
		t.Fatal("failed to connect database")
	}
	db.Set("redis", getRedis())
	db.Set("incrKey", getConf().String("incrKey"))
	m = ProvideMigrator(db, config.AppName("test"))
	err = m.Migrate()
	if err != nil {
		tearDown()
		t.Fatal("failed migration")
	}
}

func tearDown() {
	db.Migrator().DropTable(&entity.Device{}, &entity.Relation{}, &entity.User{}, &entity.OrientationStep{}, "test_migrations")
}

func user(id uint) entity.User {
	return entity.User{
		Model: gorm.Model{
			ID: id,
		},
	}
}

func getConf() contract.ConfigReader {
	conf := &mc.ConfigReader{}
	conf.On("String", "incrKey").Return("kitty-users-id", nil)
	return conf
}

func getRedis() redis.UniversalClient {
	rds := redis.NewUniversalClient(
		&redis.UniversalOptions{
			Addrs: []string{"127.0.0.1:6379"},
		})
	rds.FlushAll(context.Background())

	return rds
}

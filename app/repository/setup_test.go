package repository

import (
	"flag"
	"testing"

	"github.com/go-gormigrate/gormigrate/v2"
	"glab.tagtic.cn/ad_gains/kitty/app/entity"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
		db, err = gorm.Open(sqlite.Open(":memory:?cache=shared"), &gorm.Config{})
	} else {
		db, err = gorm.Open(mysql.Open("root@tcp(127.0.0.1:3306)/kitty?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	}

	if err != nil {
		t.Fatal("failed to connect database")
	}
	m = ProvideMigrator(db, config.AppName("test"))
	err = m.Migrate()
	if err != nil {
		t.Fatal("failed migration")
	}
}

func tearDown() {
	db.Migrator().DropTable("devices", "users", "test_migrations")
}

func user(id uint) entity.User {
	return entity.User{
		Model: gorm.Model{
			ID: id,
		},
	}
}

package repository

import (
	"context"
	"database/sql"
	"github.com/Reasno/kitty/app/entity"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

var repo *UserRepo
var m *gormigrate.Gormigrate
var db *gorm.DB

func setUp(t *testing.T) {
	var err error
	///db, err = gorm.Open(sqlite.Open(":memory:?cache=shared"), &gorm.Config{})
	db, err = gorm.Open(mysql.Open("root@tcp(127.0.0.1:3306)/kitty?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		t.Fatal("failed to connect database")
	}
	m = ProvideMigrator(db)
	err = m.Migrate()
	if err != nil {
		t.Fatal("failed migration")
	}
	repo = NewUserRepo(db)
}

func tearDown() {
	db.Migrator().DropTable("devices", "users", "migrations")
}

func TestGetFromWechat(t *testing.T) {
	setUp(t)
	defer tearDown()
	ctx := context.Background()
	u, err := repo.GetFromWechat(ctx, "", "foo", &entity.Device{Suuid: "bar"}, entity.User{UserName: "baz"})
	if err != nil {
		t.Fatal(err)
	}
	if u.WechatOpenId.String != "foo" {
		t.Fatalf("want foo, got %v", u.WechatOpenId)
	}
	if u.Devices[0].Suuid != "bar" {
		t.Fatalf("want bar, got %s", u.Devices[0].Suuid)
	}
	if u.CommonSUUID != "bar" {
		t.Fatalf("want bar, got %s", u.CommonSUUID)
	}
	if u.UserName != "baz" {
		t.Fatalf("want baz, got %s", u.UserName)
	}
	u2, err := repo.GetFromWechat(ctx, "", "foo", &entity.Device{Suuid: "bar2"}, entity.User{UserName: "baz2"})
	if err != nil {
		t.Fatal(err)
	}
	if u2.WechatOpenId.String != "foo" {
		t.Fatalf("want foo, got %v", u2.WechatOpenId)
	}
	if u2.Devices[0].Suuid != "bar2" {
		t.Fatalf("want bar2, got %s", u2.Devices[0].Suuid)
	}
	if u2.CommonSUUID != "bar" {
		t.Fatalf("want bar, got %s", u2.CommonSUUID)
	}
	if u2.UserName != "baz" {
		t.Fatalf("want baz, got %s", u2.Devices[0].Suuid)
	}
}

func TestGetFromMobile(t *testing.T) {
	setUp(t)
	defer tearDown()
	ctx := context.Background()
	u, err := repo.GetFromMobile(ctx, "", "110", &entity.Device{Suuid: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if u.Mobile.String != "110" {
		t.Fatalf("want 110, got %v", u.WechatOpenId)
	}
	if u.Devices[0].Suuid != "bar" {
		t.Fatalf("want bar, got %s", u.Devices[0].Suuid)
	}
	if u.CommonSUUID != "bar" {
		t.Fatalf("want bar, got %s", u.CommonSUUID)
	}
	u2, err := repo.GetFromMobile(ctx, "", "110", &entity.Device{Suuid: "bar2"})
	if err != nil {
		t.Fatal(err)
	}
	if u2.Mobile.String != "110" {
		t.Fatalf("want foo, got %v", u2.Mobile)
	}
	if u2.Devices[0].Suuid != "bar2" {
		t.Fatalf("want bar2, got %s", u2.Devices[0].Suuid)
	}
	if u2.CommonSUUID != "bar" {
		t.Fatalf("want bar, got %s", u2.CommonSUUID)
	}
}

func TestGetFromDevice(t *testing.T) {
	setUp(t)
	defer tearDown()
	ctx := context.Background()
	u, err := repo.GetFromDevice(ctx, "", "110", &entity.Device{Suuid: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if u.CommonSUUID != "110" {
		t.Fatalf("want 110, got %s", u.CommonSUUID)
	}
	if u.Devices[0].Suuid != "bar" {
		t.Fatalf("want bar, got %s", u.Devices[0].Suuid)
	}
	u2, err := repo.GetFromDevice(ctx, "", "110", &entity.Device{Suuid: "bar2"})
	if err != nil {
		t.Fatal(err)
	}
	if u2.CommonSUUID != "110" {
		t.Fatalf("want foo, got %s", u2.CommonSUUID)
	}
	if u2.Devices[0].Suuid != "bar2" {
		t.Fatalf("want bar2, got %s", u2.Devices[0].Suuid)
	}
}

func TestGetSave(t *testing.T) {
	setUp(t)
	defer tearDown()
	ctx := context.Background()
	user := entity.User{}
	user.ID = 50
	err := repo.Save(ctx, &user)
	if err != nil {
		t.Fatal(err)
	}
	u, err := repo.Get(ctx, 50)
	if err != nil {
		t.Fatal(err)
	}
	if u.ID != user.ID {
		t.Fatalf("want %d, go %d", user.ID, u.ID)
	}
}

func TestUniqueConstraint(t *testing.T) {
	setUp(t)
	defer tearDown()
	ctx := context.Background()
	user := entity.User{
		Mobile: sql.NullString{"110", true},
	}
	err := repo.Save(ctx, &user)
	if err != nil {
		t.Fatal(err)
	}
	user2 := entity.User{
		Mobile: sql.NullString{"110", true},
	}
	err = repo.Save(ctx, &user2)
	if err == nil {
		t.Fatal(err)
	}
	user3 := entity.User{
		WechatOpenId: sql.NullString{"110", true},
	}
	err = repo.Save(ctx, &user3)
	if err != nil {
		t.Fatal(err)
	}
	user4 := entity.User{
		WechatOpenId: sql.NullString{"110", true},
	}
	err = repo.Save(ctx, &user4)
	if err == nil {
		t.Fatal(err)
	}
}

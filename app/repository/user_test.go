package repository

import (
	"context"
	"github.com/Reasno/kitty/app/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

var repo UserRepo
func setUp(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal("failed to connect database")
	}
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Device{})
	repo = UserRepo{db}
}

func TestGetFromWechat(t *testing.T) {
	setUp(t)
	ctx := context.Background()
	u, err := repo.GetFromWechat(ctx, "foo", &entity.Device{Suuid: "bar"}, entity.User{UserName: "baz"})
	if err != nil {
		t.Fatal(err)
	}
	if u.WechatOpenId != "foo" {
		t.Fatalf("want foo, got %s", u.WechatOpenId)
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
	u2, err := repo.GetFromWechat(ctx, "foo", &entity.Device{Suuid: "bar2"}, entity.User{UserName: "baz2"})
	if err != nil {
		t.Fatal(err)
	}
	if u2.WechatOpenId != "foo" {
		t.Fatalf("want foo, got %s", u2.WechatOpenId)
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
	ctx := context.Background()
	u, err := repo.GetFromMobile(ctx, "110", &entity.Device{Suuid: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if u.Mobile != "110" {
		t.Fatalf("want 110, got %s", u.WechatOpenId)
	}
	if u.Devices[0].Suuid != "bar" {
		t.Fatalf("want bar, got %s", u.Devices[0].Suuid)
	}
	if u.CommonSUUID != "bar" {
		t.Fatalf("want bar, got %s", u.CommonSUUID)
	}
	u2, err := repo.GetFromMobile(ctx, "110", &entity.Device{Suuid: "bar2"})
	if err != nil {
		t.Fatal(err)
	}
	if u2.Mobile != "110" {
		t.Fatalf("want foo, got %s", u2.Mobile)
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
	ctx := context.Background()
	u, err := repo.GetFromDevice(ctx, "110", &entity.Device{Suuid: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if u.CommonSUUID != "110" {
		t.Fatalf("want 110, got %s", u.CommonSUUID)
	}
	if u.Devices[0].Suuid != "bar" {
		t.Fatalf("want bar, got %s", u.Devices[0].Suuid)
	}
	u2, err := repo.GetFromDevice(ctx, "110", &entity.Device{Suuid: "bar2"})
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

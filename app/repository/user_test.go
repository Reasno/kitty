package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"glab.tagtic.cn/ad_gains/kitty/app/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TestGetFromWechat(t *testing.T) {
	if !useMysql {
		t.Skip("GetFromWechat tests must be run on mysql")
	}
	setUp(t)
	defer tearDown()
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), getRedis(), getConf())
	ctx := context.Background()
	u, err := userRepo.GetFromWechat(ctx, "", "foo", &entity.Device{Suuid: "bar"}, entity.User{UserName: "baz"})
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
	if !u.IsNew {
		t.Fatalf("user must be new")
	}
	u2, err := userRepo.GetFromWechat(ctx, "", "foo", &entity.Device{Suuid: "bar2"}, entity.User{UserName: "baz2"})
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
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), getRedis(), getConf())
	ctx := context.Background()
	u, err := userRepo.GetFromMobile(ctx, "", "110", &entity.Device{Suuid: "bar"})
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
	u2, err := userRepo.GetFromMobile(ctx, "", "110", &entity.Device{Suuid: "bar2"})
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
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), getRedis(), getConf())
	ctx := context.Background()
	u, err := userRepo.GetFromDevice(ctx, "", "110", &entity.Device{Suuid: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if u.CommonSUUID != "110" {
		t.Fatalf("want 110, got %s", u.CommonSUUID)
	}
	if u.Devices[0].Suuid != "bar" {
		t.Fatalf("want bar, got %s", u.Devices[0].Suuid)
	}
	u2, err := userRepo.GetFromDevice(ctx, "", "110", &entity.Device{Suuid: "bar2"})
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
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), getRedis(), getConf())
	ctx := context.Background()
	user := entity.User{}
	user.ID = 50
	err := userRepo.Save(ctx, &user)
	if err != nil {
		t.Fatal(err)
	}
	u, err := userRepo.Get(ctx, 50)
	if err != nil {
		t.Fatal(err)
	}
	if u.ID != user.ID {
		t.Fatalf("want %d, go %d", user.ID, u.ID)
	}
}

func TestUserRepo_Delete(t *testing.T) {
	setUp(t)
	defer tearDown()
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), getRedis(), getConf())
	ctx := context.Background()
	user := entity.User{Model: gorm.Model{ID: uint(1)}, UserName: "hello"}
	_ = userRepo.Save(ctx, &user)

	err := userRepo.Delete(ctx, 1)
	assert.NoError(t, err)
	_, err = userRepo.Get(ctx, 1)
	assert.True(t, errors.Is(err, ErrRecordNotFound))
}

func TestUserRepo_GetAll(t *testing.T) {
	setUp(t)
	defer tearDown()

	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), getRedis(), getConf())
	ctx := context.Background()
	for i := 1; i < 5; i++ {
		user := entity.User{Model: gorm.Model{ID: uint(i)}, UserName: "hello"}
		_ = userRepo.Save(ctx, &user)
	}
	users, err := userRepo.GetAll(ctx, clause.Where{Exprs: []clause.Expression{clause.IN{
		Column: "id",
		Values: []interface{}{1, 2, 3, 4},
	}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 4 {
		t.Fatal("there should be four users")
	}
	users, err = userRepo.GetAll(ctx, clause.Where{Exprs: []clause.Expression{clause.Like{
		Column: "user_name",
		Value:  "%ell%",
	}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 4 {
		t.Fatal("there should be four users")
	}
	users, err = userRepo.GetAll(ctx, clause.Where{Exprs: []clause.Expression{clause.Gt{
		Column: "created_at",
		Value:  time.Unix(500, 0),
	}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 4 {
		t.Fatal("there should be four users")
	}
}

func TestUserRepo_Count(t *testing.T) {
	setUp(t)
	defer tearDown()

	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), getRedis(), getConf())
	ctx := context.Background()
	for i := 1; i < 5; i++ {
		user := entity.User{Model: gorm.Model{ID: uint(i)}, UserName: "hello"}
		_ = userRepo.Save(ctx, &user)
	}
	count, err := userRepo.Count(ctx, clause.Where{Exprs: []clause.Expression{clause.IN{
		Column: "id",
		Values: []interface{}{1, 2, 3, 4},
	}}})
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Fatal("there should be four users")
	}
}

func TestUniqueConstraint(t *testing.T) {
	if !useMysql {
		t.Skip("unique constraints tests must be run on mysql")
	}
	setUp(t)
	defer tearDown()
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), getRedis(), getConf())
	ctx := context.Background()
	user := entity.User{
		Mobile: sql.NullString{String: "110", Valid: true},
	}
	err := userRepo.Save(ctx, &user)
	if err != nil {
		t.Fatal(err)
	}
	user2 := entity.User{
		Mobile: sql.NullString{String: "110", Valid: true},
	}
	err = userRepo.Save(ctx, &user2)
	if err != ErrAlreadyBind {
		t.Fatal(err)
	}
	user3 := entity.User{
		WechatOpenId: sql.NullString{String: "110", Valid: true},
	}
	err = userRepo.Save(ctx, &user3)
	if err != nil {
		t.Fatal(err)
	}
	user4 := entity.User{
		WechatOpenId: sql.NullString{String: "110", Valid: true},
	}
	err = userRepo.Save(ctx, &user4)
	if err != ErrAlreadyBind {
		t.Fatal(err)
	}
}

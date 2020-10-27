package repository

import (
	"context"
	"github.com/Reasno/kitty/app/entity"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func (r *UserRepo) Update(ctx context.Context, id uint, user entity.User) (newUser *entity.User, err error) {
	var (
		u entity.User
	)
	r.db.WithContext(ctx).Model(entity.User{}).Where("id = ?", id).Updates(user)
	r.db.WithContext(ctx).First(&u, "id = ?", id)
	return &u, nil
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) GetFromWechat(ctx context.Context, wechat string, device *entity.Device, wechatUser entity.User) (*entity.User, error) {
	var (
		u entity.User
	)
	wechatUser.CommonSUUID = device.Suuid
	r.db.WithContext(ctx).Where(entity.User{WechatOpenId: wechat}).Attrs(wechatUser).FirstOrCreate(&u)
	u.AddNewDevice(device)
	r.db.WithContext(ctx).Save(device)
	return &u, nil
}

func (r *UserRepo) GetFromMobile(ctx context.Context, mobile string, device *entity.Device) (*entity.User, error) {
	var (
		u entity.User
	)
	r.db.WithContext(ctx).Where(entity.User{Mobile: mobile}).Attrs(entity.User{CommonSUUID: device.Suuid}).FirstOrCreate(&u)
	u.AddNewDevice(device)
	r.db.WithContext(ctx).Save(device)
	return &u, nil
}

func (r *UserRepo) GetFromDevice(ctx context.Context, suuid string, device *entity.Device) (*entity.User, error) {
	var (
		u entity.User
	)
	r.db.WithContext(ctx).Where(entity.User{CommonSUUID: suuid}).FirstOrCreate(&u)
	u.AddNewDevice(device)
	r.db.WithContext(ctx).Save(device)
	return &u, nil
}

func (r *UserRepo) Get(ctx context.Context, id uint) (*entity.User, error) {
	var (
		u entity.User
	)
	r.db.WithContext(ctx).First(&u, "id = ?", id)
	return &u, nil
}

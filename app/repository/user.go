package repository

import (
	"context"
	"github.com/Reasno/kitty/app/entity"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) GetFromWechat(ctx context.Context, wechat string, device *entity.Device) (*entity.User, error) {
	var (
		u entity.User
	)
	r.db.WithContext(ctx).Where(entity.User{Wechat: wechat}).FirstOrCreate(&u)
	u.AddNewDevice(device)
	r.db.WithContext(ctx).Save(device)
	return &u, nil
}

func (r *UserRepo) GetFromMobile(ctx context.Context, mobile string, device *entity.Device) (*entity.User, error) {
	var (
		u entity.User
	)
	r.db.WithContext(ctx).Where(entity.User{Mobile: mobile}).FirstOrCreate(&u)
	u.AddNewDevice(device)
	r.db.WithContext(ctx).Save(device)
	return &u, nil
}

package repository

import (
	"context"
	"database/sql"
	"github.com/Reasno/kitty/app/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

const emsg = "UserRepo"

func (r *UserRepo) Save(ctx context.Context, user *entity.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return errors.Wrap(err, emsg)
	}
	return nil
}

func (r *UserRepo) Update(ctx context.Context, id uint, user entity.User) (newUser *entity.User, err error) {
	var (
		u entity.User
	)
	err = r.db.WithContext(ctx).Model(entity.User{}).Where("id = ?", id).Updates(user).Error
	if err != nil {
		return nil, errors.Wrap(err, emsg)
	}
	err = r.db.WithContext(ctx).First(&u, "id = ?", id).Error
	if err != nil {
		return nil, errors.Wrap(err, emsg)
	}
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
	err := r.db.WithContext(ctx).Where(entity.User{WechatOpenId: sql.NullString{String: wechat, Valid: true}}).Attrs(wechatUser).FirstOrCreate(&u).Error
	if err != nil {
		return nil, errors.Wrap(err, emsg)
	}
	u.AddNewDevice(device)
	r.db.WithContext(ctx).Save(device)
	return &u, nil
}

func (r *UserRepo) GetFromMobile(ctx context.Context, mobile string, device *entity.Device) (*entity.User, error) {
	var (
		u entity.User
	)
	err := r.db.WithContext(ctx).Where(entity.User{Mobile: sql.NullString{String: mobile, Valid: true}}).Attrs(entity.User{CommonSUUID: device.Suuid}).FirstOrCreate(&u).Error
	if err != nil {
		return nil, errors.Wrap(err, emsg)
	}
	u.AddNewDevice(device)
	r.db.WithContext(ctx).Save(device)
	return &u, nil
}

func (r *UserRepo) GetFromDevice(ctx context.Context, suuid string, device *entity.Device) (*entity.User, error) {
	var (
		err error
		u entity.User
	)

	err = r.db.WithContext(ctx).Where(entity.User{CommonSUUID: suuid}).FirstOrCreate(&u).Error
	if err != nil {
		return nil, errors.Wrap(err, emsg)
	}
	u.AddNewDevice(device)
	r.db.WithContext(ctx).Save(device)
	return &u, nil
}

func (r *UserRepo) Get(ctx context.Context, id uint) (*entity.User, error) {
	var (
		u entity.User
	)
	if err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		return nil, errors.Wrap(err, emsg)
	}
	return &u, nil
}

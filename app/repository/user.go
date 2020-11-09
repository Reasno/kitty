package repository

import (
	"context"
	"database/sql"
	"github.com/Reasno/kitty/app/entity"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

var ErrAlreadyBind = errors.New("third party account is bound to another user")
var ErrRecordNotFound = errors.New("record not found")

const emsg = "UserRepo"

func (r *UserRepo) Save(ctx context.Context, user *entity.User) error {
	if err := r.db.Save(user).Error; err != nil {
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				return ErrAlreadyBind
			}
		}
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
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				return nil, ErrAlreadyBind
			}
		}
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

func (r *UserRepo) GetFromWechat(ctx context.Context, packageName, wechat string, device *entity.Device, wechatUser entity.User) (*entity.User, error) {
	var (
		u entity.User
	)
	wechatUser.CommonSUUID = device.Suuid
	wechatUser.PackageName = packageName
	wechatUser.WechatOpenId = sql.NullString{wechat, true}
	err := r.db.WithContext(ctx).Where("package_name = ? and wechat_open_id = ?", packageName, wechat).Attrs(wechatUser).FirstOrCreate(&u).Error
	if err != nil {
		return nil, errors.Wrap(err, emsg)
	}
	u.AddNewDevice(device)
	r.db.WithContext(ctx).Save(device)
	return &u, nil
}

func (r *UserRepo) GetFromMobile(ctx context.Context, packageName, mobile string, device *entity.Device) (*entity.User, error) {
	var (
		u entity.User
	)
	err := r.db.WithContext(ctx).Where("package_name = ? and mobile = ?", packageName, mobile).Attrs(entity.User{CommonSUUID: device.Suuid, PackageName: packageName, Mobile: sql.NullString{mobile, true}}).FirstOrCreate(&u).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, emsg)
	}
	u.AddNewDevice(device)
	r.db.WithContext(ctx).Save(device)
	return &u, nil
}

func (r *UserRepo) GetFromDevice(ctx context.Context, packageName, suuid string, device *entity.Device) (*entity.User, error) {
	var (
		err error
		u   entity.User
	)

	err = r.db.WithContext(ctx).Where("package_name = ? and common_s_uuid = ?", packageName, suuid).Attrs(entity.User{PackageName: packageName, CommonSUUID: suuid}).FirstOrCreate(&u).Error
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
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, errors.Wrap(err, emsg)
	}
	return &u, nil
}

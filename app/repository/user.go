package repository

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/app/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepo struct {
	fr *FileRepo
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

func (r *UserRepo) UpdateCallback(ctx context.Context, id uint, f func(user *entity.User) error) (err error) {
	var u entity.User
	return r.db.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)
		err := tx.Model(entity.User{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&u).Error
		if err != nil {
			return errors.Wrap(err, emsg)
		}
		err = f(&u)
		if err != nil {
			return err
		}

		err = tx.Save(u).Error
		if err != nil {
			if err, ok := err.(*mysql.MySQLError); ok {
				if err.Number == 1062 {
					return ErrAlreadyBind
				}
			}
			return errors.Wrap(err, emsg)
		}
		return nil
	})
}

func NewUserRepo(db *gorm.DB, fr *FileRepo) *UserRepo {
	return &UserRepo{fr, db}
}

func (r *UserRepo) GetFromWechat(ctx context.Context, packageName, wechat string, device *entity.Device, wechatUser entity.User) (*entity.User, error) {
	var (
		u entity.User
	)

	wechatUser.CommonSUUID = device.Suuid
	wechatUser.PackageName = packageName
	wechatUser.WechatOpenId = sql.NullString{String: wechat, Valid: true}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		defer func() {
			u.AddNewDevice(device)
			tx.WithContext(ctx).Save(device)
		}()
		err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where(
			"package_name = ? and wechat_open_id = ?", packageName, wechat,
		).First(&u).Error

		if err == nil {
			return nil
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		if wechatUser.HeadImg != "" {
			wechatUser.HeadImg, _ = r.fr.UploadFromUrl(ctx, wechatUser.HeadImg)
		}
		if err := tx.WithContext(ctx).Create(&wechatUser).Error; err != nil {
			return err
		}
		u = wechatUser
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, emsg)
	}
	return &u, nil
}

func (r *UserRepo) GetFromMobile(ctx context.Context, packageName, mobile string, device *entity.Device) (*entity.User, error) {
	var (
		u entity.User
	)
	err := r.db.WithContext(ctx).Where("package_name = ? and mobile = ?", packageName, mobile).Attrs(entity.User{CommonSUUID: device.Suuid, PackageName: packageName, Mobile: sql.NullString{String: mobile, Valid: true}}).FirstOrCreate(&u).Error
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

func (r *UserRepo) GetAll(ctx context.Context, ids ...uint) ([]entity.User, error) {
	var (
		u []entity.User
	)
	if err := r.db.WithContext(ctx).Find(&u, ids).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, errors.Wrap(err, emsg)
	}
	return u, nil
}

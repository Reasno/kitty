package repository

import (
	"database/sql"
	"fmt"

	"github.com/Reasno/kitty/app/entity"
	"github.com/Reasno/kitty/pkg/contract"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func ProvideMigrator(db *gorm.DB, appName contract.AppName) *gormigrate.Gormigrate {
	return gormigrate.New(db, &gormigrate.Options{
		TableName: fmt.Sprintf("%s_migrations", appName.String()),
	}, []*gormigrate.Migration{
		{
			ID: "202010280100",
			Migrate: func(db *gorm.DB) error {
				type Device struct {
					gorm.Model
					UserID    uint
					Os        uint8
					Imei      string
					Idfa      string
					Oaid      string
					Suuid     string
					Mac       string
					AndroidId string
					// 仅供数据库去重使用，应用不应依赖该字段，以免去重条件发生变化
					Hash string `gorm:"type:varchar(255);uniqueIndex:hash_index,sort:desc"`
				}
				type User struct {
					gorm.Model
					UserName      string         `json:"user_name" gorm:"default:游客"`
					WechatOpenId  sql.NullString `json:"wechat_openid" gorm:"type:varchar(255);uniqueIndex:wechat_openid_index"`
					WechatUnionId sql.NullString `json:"wechat_unionid"`
					HeadImg       string         `json:"head_img"`
					Gender        int            `json:"gender"`
					Birthday      string         `json:"birthday" gorm:"default:2000-01-01"`
					Mobile        sql.NullString `json:"mobile" gorm:"type:varchar(255);uniqueIndex:mobile_index"`
					CommonSUUID   string         `json:"common_suuid"`
					Devices       []Device
					Channel       string `json:"channel"`
					VersionCode   string `json:"version_code"`
					InviteCode    string `json:"invite_code"`
					PackageName   string `gorm:"type:varchar(255);uniqueIndex:mobile_index,priority:1;uniqueIndex:wechat_openid_index,priority:1"`
				}
				return db.AutoMigrate(
					&User{},
					&Device{},
				)
			},
			Rollback: func(db *gorm.DB) error {
				type User struct{}
				type Device struct{}
				return db.Migrator().DropTable(&User{}, &Device{})
			},
		},
		{
			ID: "202011030100",
			Migrate: func(db *gorm.DB) error {
				type User struct {
					ThirdPartyId string
				}
				if !db.Migrator().HasColumn(&User{}, "ThirdPartyId") {
					return db.Migrator().AddColumn(&User{}, "ThirdPartyId")
				}
				return nil
			},
			Rollback: func(db *gorm.DB) error {
				type User struct {
					ThirdPartyId string
				}
				if db.Migrator().HasColumn(&entity.User{}, "ThirdPartyId") {
					return db.Migrator().DropColumn(&entity.User{}, "ThirdPartyId")
				}
				return nil
			},
		},
		{
			ID: "202011050100",
			Migrate: func(db *gorm.DB) error {
				type User struct {
					PackageName  string         `gorm:"type:varchar(255);uniqueIndex:mobile_index,priority:1;uniqueIndex:wechat_openid_index,priority:1;uniqueIndex:taobao_openid_index,priority:1"`
					TaobaoOpenId sql.NullString `json:"taobao_openid" gorm:"type:varchar(255);uniqueIndex:taobao_openid_index"`
				}
				if !db.Migrator().HasColumn(&User{}, "TaobaoOpenId") {
					err := db.Migrator().AddColumn(&User{}, "TaobaoOpenId")
					if err != nil {
						return err
					}
					err = db.Migrator().CreateIndex(&User{}, "taobao_openid_index")
					if err != nil {
						return err
					}
				}
				return nil
			},
			Rollback: func(db *gorm.DB) error {
				type User struct {
					PackageName  string         `gorm:"type:varchar(255);uniqueIndex:mobile_index,priority:1;uniqueIndex:wechat_openid_index,priority:1;uniqueIndex:taobao_openid_index,priority:1"`
					TaobaoOpenId sql.NullString `json:"taobao_openid" gorm:"type:varchar(255);uniqueIndex:taobao_openid_index"`
				}
				if db.Migrator().HasColumn(&entity.User{}, "TaobaoOpenId") {
					err := db.Migrator().DropIndex(&User{}, "taobao_openid_index")
					if err != nil {
						return err
					}
					return db.Migrator().DropColumn(&User{}, "TaobaoOpenId")
				}
				return nil
			},
		},
	})
}

package repository

import (
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
				return db.AutoMigrate(
					&entity.User{},
					&entity.Device{},
				)
			},
			Rollback: func(db *gorm.DB) error {
				return db.Migrator().DropTable(&entity.User{}, &entity.Device{})
			},
		},
		{
			ID: "202011030100",
			Migrate: func(db *gorm.DB) error {
				if !db.Migrator().HasColumn(&entity.User{}, "ThirdPartyId") {
					return db.Migrator().AddColumn(&entity.User{}, "ThirdPartyId")
				}
				return nil
			},
			Rollback: func(db *gorm.DB) error {
				if db.Migrator().HasColumn(&entity.User{}, "ThirdPartyId") {
					return db.Migrator().DropColumn(&entity.User{}, "ThirdPartyId")
				}
				return nil
			},
		},
		{
			ID: "202011050100",
			Migrate: func(db *gorm.DB) error {
				if !db.Migrator().HasColumn(&entity.User{}, "TaobaoOpenId") {
					err := db.Migrator().CreateIndex(&entity.User{}, "taobao_openid_index")
					if err != nil {
						return err
					}
					return db.Migrator().AddColumn(&entity.User{}, "TaobaoOpenId")
				}
				return nil
			},
			Rollback: func(db *gorm.DB) error {
				if db.Migrator().HasColumn(&entity.User{}, "TaobaoOpenId") {
					err := db.Migrator().DropIndex(&entity.User{}, "taobao_openid_index")
					if err != nil {
						return err
					}
					return db.Migrator().DropColumn(&entity.User{}, "TaobaoOpenId")
				}
				return nil
			},
		},
	})
}

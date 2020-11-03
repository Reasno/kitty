package repository

import (
	"github.com/Reasno/kitty/app/entity"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func ProvideMigrator(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
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
	})
}

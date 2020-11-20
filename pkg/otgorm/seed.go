package otgorm

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Seed struct {
	Name string
	Run  func(*gorm.DB) error
}

type Seeds struct {
	Db    *gorm.DB
	Seeds []Seed
}

func (s *Seeds) Seed() error {
	for _, ss := range s.Seeds {
		if err := ss.Run(s.Db); err != nil {
			return errors.Wrapf(err, "failed to run %s", ss.Name)
		}
	}
	return nil
}

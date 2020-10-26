package repository

import (
	"fmt"
	"github.com/Reasno/kitty/app/entity"
	"gorm.io/gorm"
	"testing"
)

func TestUserMd5(t *testing.T)  {
	device := &entity.Device{
		Model:     gorm.Model{},
		UserID:    0,
		Os:        0,
		Imei:      "",
		Idfa:      "",
		Oaid:      "",
		Suuid:     "",
		Mac:       "",
		AndroidId: "",
		Hash:      "",
	}
	fmt.Println(device.HashCode())
}

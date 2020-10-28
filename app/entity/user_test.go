package entity

import (
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestUserEquals(t *testing.T)  {
	device1 := Device{
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
	device2 := Device{
		Model:     gorm.Model{},
		UserID:    0,
		Os:        0,
		Imei:      "a",
		Idfa:      "",
		Oaid:      "",
		Suuid:     "",
		Mac:       "",
		AndroidId: "",
		Hash:      "",
	}
	device3 := Device{
		Model:     gorm.Model{
			UpdatedAt: time.Now(),
		},
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
	if device1.Equals(&device2) {
		t.Fatal("two devices should not equal")
	}
	if ! device1.Equals(&device3) {
		t.Fatal("two devices should equal")
	}
}

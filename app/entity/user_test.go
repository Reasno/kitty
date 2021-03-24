package entity

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestUserEquals(t *testing.T) {
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
		Model: gorm.Model{
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
	if !device1.Equals(&device3) {
		t.Fatal("two devices should equal")
	}
}

type mockID struct{}

func (m mockID) ID(ctx context.Context) (uint, error) {
	return 42, nil
}

func TestUserIDAssigner(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	user := User{}
	assert.Error(t, user.BeforeCreate(db))
	assert.NoError(t, user.BeforeCreate(db.Set("IDAssigner", mockID{})))
	assert.Equal(t, uint(42), user.ID)
}

package entity

import (
	"crypto/md5"
	"fmt"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User describes a user
type User struct {
	gorm.Model
	UserName      string `json:"user_name" gorm:"default:游客"`
	WechatOpenId  string `json:"wechat_openid"`
	WechatUnionId string `json:"wechat_unionid"`
	HeadImg       string `json:"head_img"`
	Gender        int    `json:"gender"`
	Birthday      string `json:"birthday" gorm:"default:2000-01-01"`
	Mobile        string `json:"mobile"`
	CommonSUUID string `json:"common_suuid"`
	Devices       []Device
	Channel       string `json:"channel"`
	VersionCode   string `json:"version_code"`
	InviteCode    string `json:"invite_code"`
}

func (user *User) HasDevice(device *Device) bool {
	for _, mine := range user.Devices {
		if mine.Equals(device) {
			return true
		}
	}
	return false
}

func (user *User) AddNewDevice(device *Device) {
	device.Hash = device.HashCode()
	device.UserID = user.ID
	user.Devices = append(user.Devices, *device)
}

// Device describes a device.
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

// HashCode 生成唯一键
func (my Device) HashCode() string {
	m := md5.New()
	m.Write(uint64ToBytes(my.ID))
	m.Write([]byte(my.Idfa))
	m.Write([]byte(my.Imei))
	m.Write([]byte(my.Oaid))
	m.Write([]byte(my.Suuid))
	m.Write([]byte(my.Mac))
	m.Write([]byte(my.AndroidId))
	return fmt.Sprintf("%x", m.Sum(nil))
}

func (my Device) Equals(that *Device) bool {
	return my.HashCode() == that.HashCode()
}

func uint64ToBytes(n uint) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
		byte(n >> 16),
		byte(n >> 24),
		byte(n >> 32),
		byte(n >> 40),
		byte(n >> 48),
		byte(n >> 56),
	}
}

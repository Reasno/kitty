package entity

import (
	"crypto/md5"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User describes a user
type User struct {
	gorm.Model
	UserName      string         `json:"user_name" gorm:"default:游客;type:varchar(30)"`
	WechatOpenId  sql.NullString `json:"wechat_openid" gorm:"type:varchar(255);uniqueIndex:wechat_openid_index"`
	WechatUnionId sql.NullString `json:"wechat_unionid"`
	HeadImg       string         `json:"head_img" gorm:"default:http://ad-static-xg.tagtic.cn/ad-material/file/0b8f18e1e666474291174ba316cccb51.png"`
	Gender        int            `json:"gender"`
	Birthday      string         `json:"birthday" gorm:"default:2000-01-01"`
	Mobile        sql.NullString `json:"mobile" gorm:"type:varchar(255);uniqueIndex:mobile_index"`
	CommonSUUID   string         `json:"common_suuid"`
	CommonSMID    string         `json:"common_smid"`
	Devices       []Device
	Channel       string `json:"channel"`
	VersionCode   string `json:"version_code"`
	InviteCode    string `json:"invite_code"`
	PackageName   string `gorm:"type:varchar(255);uniqueIndex:mobile_index,priority:1;uniqueIndex:wechat_openid_index,priority:1;uniqueIndex:taobao_openid_index,priority:1"`
	ThirdPartyId  string
	TaobaoOpenId  sql.NullString `json:"taobao_openid" gorm:"type:varchar(255);uniqueIndex:taobao_openid_index"`
	CampaignID    sql.NullString `gorm:"type:varchar(255)"`
	AID           sql.NullString `gorm:"type:varchar(255)"`
	CID           sql.NullString `gorm:"type:varchar(255)"`
	IsNew         bool           `gorm:"-"`
	WechatExtra   []byte         `gorm:"type:blob"`
	TaobaoExtra   []byte         `gorm:"type:blob"`
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

// BeforeCreate is a gorm hook
func (user *User) BeforeCreate(db *gorm.DB) (err error) {
	if user.ID != 0 {
		return
	}

	var (
		assigner IDAssigner
	)

	if v, ok := db.Get("IDAssigner"); ok {
		assigner, _ = v.(IDAssigner)
	}

	if assigner == nil {
		return errors.New("IDAssigner missing")
	}

	id, err := assigner.ID(db.Statement.Context)
	if err != nil {
		return errors.Wrap(err, "failed to assign ID in BeforeCreate hook")
	}

	user.ID = id
	return nil
}

// AfterCreate is a gorm hook
func (user *User) AfterCreate(tx *gorm.DB) (err error) {
	user.IsNew = true
	return
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
	SMID      string
	IP        string
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
	m.Write([]byte(my.SMID))
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

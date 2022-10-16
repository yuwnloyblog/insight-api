package dbs

import (
	"time"
)

type UserStatus int

const (
	UserStatus_FORBIDDEN  = -1
	UserStatus_NORMAL     = 0
	UserStatus_MONTH_PAY  = 101
	UserStatus_SEASON_PAY = 102
	UserStatus_YEAR_PAY   = 103
)

type UserDao struct {
	ID       int64  `gorm:"primary_key"`
	Account  string `gorm:"account"`
	Password string `gorm:"password"`
	NickName string `gorm:"nack_name"`
	Status   int    `gorm:"status"`
	Phone    string `gorm:"phone"`
	Email    string `gorm:"email"`
	//WxUnionid  string    `gorm:"wx_unionid"`
	WxOpenid   string    `gorm:"wx_openid"`
	Avator     string    `gorm:"avator"`
	CreateTime time.Time `gorm:"create_time"`
}

func (user UserDao) TableName() string {
	return "users"
}

func (user UserDao) Create(u UserDao) (int64, error) {
	err := db.Create(&u).Error
	return u.ID, err
}

func (user UserDao) FindById(id int64) (*UserDao, error) {
	var item UserDao
	err := db.Where("id = ?", id).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (user UserDao) FindByWxOpenid(openid string) (*UserDao, error) {
	var item UserDao
	err := db.Where("wx_openid=?", openid).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (user UserDao) Updates(id int64, upd map[string]interface{}) error {
	if id > 0 && len(upd) > 0 {
		return db.Model(&user).Where("id=?", id).Updates(upd).Error
	}
	return nil
}

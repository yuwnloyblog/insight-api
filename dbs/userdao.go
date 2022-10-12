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
	ID         int64     `gorm:"primary_key"`
	Account    string    `gorm:"account"`
	Password   string    `gorm:"password"`
	NickName   string    `gorm:"nack_name"`
	Status     int       `gorm:"status"`
	Phone      string    `gorm:"phone"`
	Email      string    `gorm:"email"`
	WxUnionId  string    `gorm:"wx_unionid"`
	Avator     string    `gorm:"avator"`
	CreateDate time.Time `gorm:"create_time"`
}

func (user UserDao) TableName() string {
	return "users"
}

func (user UserDao) Create(u UserDao) (int64, error) {
	err := db.Create(&u).Error
	return u.ID, err
}

func (user UserDao) FindByWxUnionId(unionId string) (*UserDao, error) {
	var item UserDao
	err := db.Where("wx_unionid=?", unionId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

package dbs

import (
	"time"
)

type UserDao struct {
	ID         int64     `gorm:"primary_key"`
	Name       string    `gorm:"name"`
	Phone      string    `gorm:"phone"`
	WxUnionId  string    `gorm:"wx_unionid"`
	CreateDate time.Time `gorm:"create_time"`
	Status     int       `gorm:"status"`
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

package dbs

import (
	"time"
)

type UserStatus int

const (
	UserStatus_FORBIDDEN    = -1
	UserStatus_NORMAL       = 0
	UserStatus_MONTH_PAY    = 101
	UserStatus_SEASON_PAY   = 102
	UserStatus_HALFYEAR_PAY = 103
	UserStatus_YEAR_PAY     = 104
)

type UserDao struct {
	ID       int64  `gorm:"primary_key"`
	Account  string `gorm:"account"`
	Password string `gorm:"password"`
	NickName string `gorm:"nack_name"`
	Status   int    `gorm:"status"`
	Phone    string `gorm:"phone"`
	Email    string `gorm:"email"`
	City     string `gorm:"city"`
	Country  string `gorm:"country"`
	Gender   int    `gorm:"gender"`
	Language string `gorm:"language"`
	Province string `gorm:"province"`
	//WxUnionid  string    `gorm:"wx_unionid"`
	WxOpenid      string    `gorm:"wx_openid"`
	Avatar        string    `gorm:"avatar"`
	CreateTime    time.Time `gorm:"create_time"`
	FreeCount     int       `gorm:"free_count"`
	LatestPayTime int64     `gorm:"latest_pay_time"`
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
	if !checkPaytime(item.LatestPayTime, item.Status) {
		item.Status = 0
	}
	return &item, nil
}

func (user UserDao) FindByWxOpenid(openid string) (*UserDao, error) {
	var item UserDao
	err := db.Where("wx_openid=?", openid).Take(&item).Error
	if err != nil {
		return nil, err
	}
	if !checkPaytime(item.LatestPayTime, item.Status) {
		item.Status = 0
	}
	return &item, nil
}

func checkPaytime(payTime int64, status int) bool {
	var deta int64 = 0
	if status == UserStatus_MONTH_PAY {
		deta = 31 * 24 * 3600 * 1000
	} else if status == UserStatus_SEASON_PAY {
		deta = 92 * 24 * 3600 * 1000
	} else if status == UserStatus_HALFYEAR_PAY {
		deta = 183 * 24 * 3600 * 1000
	} else if status == UserStatus_YEAR_PAY {
		deta = 365 * 24 * 3600 * 1000
	}

	return (time.Now().UnixMilli() - deta) < payTime
}

// func (user UserDao) Updates(id int64, upd map[string]interface{}) error {
// 	if id > 0 && len(upd) > 0 {
// 		return db.Model(&user).Where("id=?", id).Updates(upd).Error
// 	}
// 	return nil
// }

func (user UserDao) Updates(u UserDao) error {
	if u.ID > 0 {
		upd := map[string]interface{}{}
		if u.NickName != "" {
			upd["nick_name"] = u.NickName
		}
		if u.Avatar != "" {
			upd["avatar"] = u.Avatar
		}
		if u.Phone != "" {
			upd["phone"] = u.Phone
		}
		if u.City != "" {
			upd["city"] = u.City
		}
		if u.Country != "" {
			upd["country"] = u.Country
		}
		if u.Language != "" {
			upd["language"] = u.Language
		}
		if u.Province != "" {
			upd["province"] = u.Province
		}
		return db.Model(&user).Where("id=?", u.ID).Update(upd).Error
	}
	return nil
}
func (user UserDao) UpdateStatus(uid int64, status int) error {
	upd := map[string]interface{}{}
	upd["status"] = status
	return db.Model(&user).Where("id=?", uid).Update(upd).Error
}
func (user UserDao) UpdatePayStatus(wxOpenid string, status int, payTime int64) error {
	upd := map[string]interface{}{}
	upd["status"] = status
	upd["latest_pay_time"] = payTime
	return db.Model(&user).Where("wx_openid=?", wxOpenid).Update(upd).Error
}

func (user UserDao) UpdateFreeCount(uid int64, freeCount int) error {
	upd := map[string]interface{}{}
	upd["free_count"] = freeCount
	return db.Model(&user).Where("id=?", uid).Update(upd).Error
}

func (user UserDao) UpdatePhone(uid int64, phone string) error {
	upd := map[string]interface{}{}
	upd["phone"] = phone
	return db.Model(&user).Where("id=?", uid).Update(upd).Error
}

package dbs

import "time"

type PayStatus int

const (
	PayStatus_PrePay  = 0
	PayStatus_PaySucc = 1
	PayStatus_PayFail = 2
)

type UserOrderDao struct {
	UserOrderAddDao
	CreateTime time.Time `gorm:"create_time"`
	UpdateTime time.Time `gorm:"update_time"`
}

type UserOrderAddDao struct {
	ID          int64   `gorm:"primary_key"`
	OrderId     string  `gorm:"order_id"`
	Money       float64 `gorm:"money"`
	Status      int     `gorm:"status"`
	PaymentMode string  `gorm:"payment_mode"`
	FailReason  string  `gorm:"fail_reason"`
}

func (user UserOrderDao) TableName() string {
	return "user_orders"
}

func (user UserOrderAddDao) TableName() string {
	return "user_orders"
}

func (user UserOrderDao) Create(u UserOrderAddDao) (int64, error) {
	err := db.Create(&u).Error
	return u.ID, err
}

func (user UserOrderDao) FindByOrderId(orderId string) (*UserOrderDao, error) {
	var item UserOrderDao
	err := db.Where("order_id=?", orderId).Take(&item).Error
	return &item, err
}

func (user UserOrderDao) Updates(id int64, upd map[string]interface{}) error {
	if id > 0 && len(upd) > 0 {
		return db.Where("id=?", id).Updates(upd).Error
	}
	return nil
}

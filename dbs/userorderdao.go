package dbs

import "time"

type PayStatus int

const (
	PayStatus_PrePay  = 0
	PayStatus_PaySucc = 1
	PayStatus_PayFail = 2
)

type UserOrderDao struct {
	ID            int64     `gorm:"primary_key"`
	OrderNo       string    `gorm:"order_no"`
	Amount        int64     `gorm:"amount"`
	Status        int       `gorm:"status"`
	PaymentMode   string    `gorm:"payment_mode"`
	FailReason    string    `gorm:"fail_reason"`
	CreateTime    time.Time `gorm:"create_time"`
	UpdateTime    time.Time `gorm:"update_time"`
	UserId        int64     `gorm:"user_id"`
	Description   string    `gorm:"description"`
	FellowType    int       `gorm:"fellow_type"`
	WxNotify      string    `gorm:"wx_notify"`
	WxTransaction string    `gorm:"wx_transaction"`
}

func (user UserOrderDao) TableName() string {
	return "user_orders"
}

func (user UserOrderDao) Create(u UserOrderDao) (int64, error) {
	err := db.Create(&u).Error
	return user.ID, err
}

func (user UserOrderDao) FindByOrderNo(orderId string) (*UserOrderDao, error) {
	var item UserOrderDao
	err := db.Where("order_no=?", orderId).Take(&item).Error
	return &item, err
}

func (user UserOrderDao) Updates(id int64, upd map[string]interface{}) error {
	if id > 0 && len(upd) > 0 {
		return db.Model(&user).Where("id=?", id).Updates(upd).Error
	}
	return nil
}

func (user UserOrderDao) UpdatesByOrderNo(orderNo string, upd map[string]interface{}) error {
	if orderNo != "" && len(upd) > 0 {
		return db.Model(&user).Where("order_no=?", orderNo).Updates(upd).Error
	}
	return nil
}

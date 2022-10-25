package services

import (
	"insight-api/dbs"
	"time"
)

func SaveWxOrder(uid int64, userOrder UserOrder) (int64, error) {
	orderDao := dbs.UserOrderDao{}
	return orderDao.Create(dbs.UserOrderDao{
		OrderNo:     userOrder.OrderNo,
		Amount:      userOrder.Amount,
		Status:      0,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
		UserId:      uid,
		FellowType:  userOrder.FellowType,
		Description: userOrder.Description,
	})
}

func UpdateOrderNo(id int64, orderNo string) {
	orderDao := dbs.UserOrderDao{}
	orderDao.Updates(id, map[string]interface{}{
		"order_no": orderNo,
	})
}

func UpdateOrderStatus(orderNo string, wxOpenid string, status int, notify, transaction string) error {
	orderDao := dbs.UserOrderDao{}
	order, err := orderDao.FindByOrderNo(orderNo)
	if err == nil {
		// 更新用户会员状态
		err = UpdateUserPayStatusByWxOpenid(wxOpenid, order.FellowType, time.Now().UnixMilli())
		if err != nil {
			return err
		}
		RemoveUserFromCache(order.UserId)
	} else {
		return err
	}
	return orderDao.UpdatesByOrderNo(orderNo, map[string]interface{}{
		"status":         status,
		"wx_notify":      notify,
		"wx_transaction": transaction,
	})
}

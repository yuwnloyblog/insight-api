package services

import (
	"fmt"
	"insight-api/dbs"
	"insight-api/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

func RegisterOrLogin(user User) (string, error) {
	userDao := dbs.UserDao{}
	userdb, err := userDao.FindByWxUnionId(user.WxUnionId)
	var dbId int64
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//入库
			dbId, err = userDao.Create(dbs.UserDao{
				Name:       "",
				WxUnionId:  user.WxUnionId,
				CreateDate: time.Now(),
				Status:     0,
			})
			if err != nil {
				return "", &CommonError{
					Code:     10100,
					ErrorMsg: err.Error(),
				}
			}
		} else {
			return "", &CommonError{
				Code:     10100,
				ErrorMsg: err.Error(),
			}
		}
	} else {
		dbId = userdb.ID
	}
	if dbId > 0 {
		idStr, _ := utils.Encode(userdb.ID)
		return GetToken(idStr), nil
	} else {
		return "", &CommonError{
			Code:     10101,
			ErrorMsg: "no user id",
		}
	}
}

var SecureKey []byte = []byte("a2cdgxgghijk1mn0p")

func GetToken(id string) string {
	data := fmt.Sprintf("%s|%d", id, time.Now().UnixMilli())
	encrypted, _ := utils.AesEncrypt([]byte(data), SecureKey)
	return string(encrypted)
}

func ParseToken(token string) (string, int64, error) {
	bs, err := utils.AesDecrypt([]byte(token), SecureKey)
	if err != nil {
		return "", 0, err
	} else {
		ret := string(bs)
		arr := strings.Split(ret, "|")
		t, _ := utils.ParseInt64(arr[1])
		return arr[0], t, nil
	}
}

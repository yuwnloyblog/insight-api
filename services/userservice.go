package services

import (
	"encoding/base64"
	"fmt"
	"insight-api/dbs"
	"insight-api/utils"
	"strings"
	"time"
)

func RegisterOrLogin(user User) (string, *User, error) {
	userDao := dbs.UserDao{}
	retUser := &User{}
	userdb, err := userDao.FindByWxOpenid(user.WxOpenid)
	var dbId int64
	if err != nil {
		//入库
		dbId, err = userDao.Create(dbs.UserDao{
			NickName:   user.NickName,
			Avatar:     user.Avatar,
			WxOpenid:   user.WxOpenid,
			CreateTime: time.Now(),
			Status:     0,
		})
		if err != nil {
			return "", nil, GetError(ErrorCode_UserDbInsertFail)
		}
		retUser.NickName = user.NickName
		retUser.Avatar = user.Avatar
		retUser.Status = 0
	} else {
		dbId = userdb.ID
		retUser.NickName = userdb.NickName
		retUser.Avatar = userdb.Avatar
		retUser.Status = userdb.Status
	}
	if dbId > 0 {
		idStr, _ := utils.Encode(dbId)
		return GetToken(idStr), retUser, nil
	} else {
		return "", nil, GetError(ErrorCode_UserIdIs0)
	}
}

func UpdateUserInfo(uid int64, user User) error {
	userDao := dbs.UserDao{}
	err := userDao.Updates(dbs.UserDao{
		ID:       uid,
		NickName: user.NickName,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
	})
	return err
}

func GetUserInfByCache(uid int64) (*User, error) {
	key := fmt.Sprintf("user_%d", uid)
	uo, ok := utils.CacheGet(key)
	if ok {
		return uo.(*User), nil
	} else {
		user, err := GetUserInfo(uid)
		if err != nil {
			return nil, err
		}
		utils.CachePut(key, user)
		return user, nil
	}
}

func GetUserInfo(uid int64) (*User, error) {
	ud := dbs.UserDao{}
	userdb, err := ud.FindById(uid)
	if err != nil {
		return nil, err
	}
	idStr, _ := utils.Encode(userdb.ID)
	return &User{
		Id:       idStr,
		NickName: userdb.NickName,
		Phone:    userdb.Phone,
		Avatar:   userdb.Avatar,
		Status:   userdb.Status,
		WxOpenid: userdb.WxOpenid,
	}, nil
}

var SecureKey []byte = []byte("a2cdgxghijk1mn0p")

func GetToken(id string) string {
	data := fmt.Sprintf("%s|%d", id, time.Now().UnixMilli())
	encrypted, _ := utils.AesEncrypt([]byte(data), SecureKey)
	return base64.StdEncoding.EncodeToString(encrypted)
}

func ParseToken(token string) (string, int64, error) {
	encrypted, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", 0, err
	}
	bs, err := utils.AesDecrypt(encrypted, SecureKey)
	if err != nil {
		return "", 0, err
	} else {
		ret := string(bs)
		arr := strings.Split(ret, "|")
		t, _ := utils.ParseInt64(arr[1])
		return arr[0], t, nil
	}
}

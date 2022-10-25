package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"insight-api/configures"
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
			FreeCount:  10,
		})
		if err != nil {
			return "", nil, GetError(ErrorCode_UserDbInsertFail)
		}
		retUser.NickName = user.NickName
		retUser.Avatar = user.Avatar
		retUser.Status = 0
		//更新phone
		go func() {
			phone, err := CalculatePhone(user.Phone)
			if err == nil && phone != "" {
				userDao.UpdatePhone(dbId, phone)
			}
		}()
	} else {
		dbId = userdb.ID
		retUser.NickName = userdb.NickName
		retUser.Avatar = userdb.Avatar
		retUser.Status = userdb.Status
		if dbId > 0 && userdb.Phone == "" {
			//更新phone
			go func() {
				phone, err := CalculatePhone(user.Phone)
				if err == nil && phone != "" {
					userDao.UpdatePhone(dbId, phone)
				}
			}()
		}
	}
	if dbId > 0 {
		idStr, _ := utils.Encode(dbId)
		return GetToken(idStr), retUser, nil
	} else {
		return "", nil, GetError(ErrorCode_UserIdIs0)
	}
}

func UpdateUserStatus(uid int64, status int) error {
	userDao := dbs.UserDao{}
	return userDao.UpdateStatus(uid, status)
}

func UpdateUserPayStatusByWxOpenid(openId string, status int, payTime int64) error {
	userDao := dbs.UserDao{}
	return userDao.UpdatePayStatus(openId, status, payTime)
}

func UpdateUserInfo(uid int64, user User) error {
	userDao := dbs.UserDao{}
	err := userDao.Updates(dbs.UserDao{
		ID:       uid,
		NickName: user.NickName,
		Avatar:   user.Avatar,
		City:     user.City,
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

func RemoveUserFromCache(uid int64) bool {
	key := fmt.Sprintf("user_%d", uid)
	return utils.CacheRemove(key)
}

func GetUserInfo(uid int64) (*User, error) {
	ud := dbs.UserDao{}
	userdb, err := ud.FindById(uid)
	if err != nil {
		return nil, err
	}
	idStr, _ := utils.Encode(userdb.ID)
	return &User{
		Id:            idStr,
		NickName:      userdb.NickName,
		Phone:         userdb.Phone,
		Avatar:        userdb.Avatar,
		Status:        userdb.Status,
		WxOpenid:      userdb.WxOpenid,
		LatestPayTime: userdb.LatestPayTime,
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
func CheckFreeCount(uid int64) bool {
	userDao := dbs.UserDao{}
	user, err := userDao.FindById(uid)
	if err != nil {
		return false
	}
	if user.FreeCount <= 0 {
		return false
	}
	user.UpdateFreeCount(uid, user.FreeCount-1)
	return true
}

func CalculatePhone(wxPhoneCode string) (string, error) {
	//POST https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=ACCESS_TOKEN

	header := map[string]string{}
	body := fmt.Sprintf("{\"code\":\"%s\"}", wxPhoneCode)
	accessToken, err := GetWxAccessToken()
	if err != nil {
		fmt.Println("err")
	}
	resp, err := utils.HttpDo("POST", fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", accessToken), header, body)
	if err != nil {
		return "", err
	}
	fmt.Println(resp)
	var phoneResp WxPhoneResp
	err = json.Unmarshal([]byte(resp), &phoneResp)
	if err != nil {
		return "", err
	}
	return phoneResp.PhoneInfo.PhoneNumber, nil
}

func GetWxAccessToken() (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", configures.Config.Wx.AppId, configures.Config.Wx.Secret)
	resp, err := utils.HttpDo("GET", url, map[string]string{}, "")
	if err != nil {
		return "", err
	}
	var atResp AccessTokenResp
	err = json.Unmarshal([]byte(resp), &atResp)
	if err != nil {
		return "", err
	}

	return atResp.AccessToken, nil
}

type WxPhoneResp struct {
	ErrorCode int          `json:"errorcode"`
	ErrorMsg  string       `json:"errmsg"`
	PhoneInfo *WxPhoneInfo `json:"phone_info"`
}
type WxPhoneInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
}
type AccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

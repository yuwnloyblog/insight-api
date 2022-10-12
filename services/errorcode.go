package services

import "fmt"

type ErrorCode int

const (
	ErrorCode_Unknown = 10000

	//100xx 参数校验，登录和注册相关
	ErrorCode_NoWxJsCode     = 10001
	ErrorCode_WxLoginFail    = 10002
	ErrorCode_UserIdIs0      = 10003
	ErrorCode_WxLoginRespErr = 10004

	//101xx 用户，app，dev等数据相关
	ErrorCode_UserDbReadFail   = 10100
	ErrorCode_UserDbInsertFail = 10101
	ErrorCode_UserDbUpdateFail = 10102
	ErrorCode_UserDbDelFail    = 10103
)

var errMsgMap map[ErrorCode]string = map[ErrorCode]string{
	ErrorCode_Unknown: "unknown error",

	ErrorCode_NoWxJsCode:  "js_code is required",
	ErrorCode_WxLoginFail: "wx login failed",
}

type CommonError struct {
	Code     int    `json:"code"`
	ErrorMsg string `json:"msg"`
}

func (err *CommonError) Error() string {
	return fmt.Sprintf("%d:%s", err.Code, err.ErrorMsg)
}

func GetError(code ErrorCode) error {
	if msg, ok := errMsgMap[code]; ok {
		return &CommonError{
			Code:     int(code),
			ErrorMsg: msg,
		}
	}
	return &CommonError{
		Code:     int(ErrorCode_Unknown),
		ErrorMsg: "unknown error",
	}
}

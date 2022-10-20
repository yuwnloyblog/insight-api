package services

import "fmt"

type ErrorCode int

const (
	ErrorCode_Success = 0
	ErrorCode_Unknown = 10000

	//100xx 参数校验，登录和注册相关
	ErrorCode_NoWxJsCode     = 10001
	ErrorCode_WxLoginFail    = 10002
	ErrorCode_UserIdIs0      = 10003
	ErrorCode_WxLoginRespErr = 10004
	ErrorCode_TokenErr       = 10005
	ErrorCode_UidStrError    = 10006
	ErrorCode_NoUid          = 10007
	ErrorCode_NotLogin       = 10008
	ErrorCode_TokenOutDate   = 10009
	ErrorCode_NeedPay        = 10010
	ErrorCode_ParamErr       = 10011
	ErrorCode_PrepayCallErr  = 10012
	ErrorCode_PrepayParseErr = 10013

	//101xx 用户，app，dev等数据相关
	ErrorCode_UserDbReadFail   = 10100
	ErrorCode_UserDbInsertFail = 10101
	ErrorCode_UserDbUpdateFail = 10102
	ErrorCode_UserDbDelFail    = 10103
)

var errMsgMap map[ErrorCode]string = map[ErrorCode]string{
	ErrorCode_Success: "success",
	ErrorCode_Unknown: "unknown error",

	ErrorCode_NoWxJsCode:   "js_code is required",
	ErrorCode_WxLoginFail:  "wx login failed",
	ErrorCode_NotLogin:     "not login",
	ErrorCode_TokenOutDate: "session is out of date.",
	ErrorCode_NeedPay:      "need pay.",
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
		Code:     int(code),
		ErrorMsg: "unknown error",
	}
}

func GetSuccess() interface{} {
	return &CommonError{
		Code:     int(ErrorCode_Success),
		ErrorMsg: "success",
	}
}

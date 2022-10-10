package apis

import (
	"encoding/json"
	"fmt"
	"insight-api/configures"
	"insight-api/services"
	"insight-api/tools"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WxLogin(ctx *gin.Context) {
	jsCode := ctx.PostForm("js_code")

	if jsCode == "" {
		ctx.JSON(http.StatusBadRequest, services.CommonError{
			Code:     10000,
			ErrorMsg: "js_code is required",
		})
		return
	}

	appId := configures.Config.Wx.AppId
	secret := configures.Config.Wx.Secret

	header := map[string]string{}
	resp, err := tools.HttpDo("GET", fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appId, secret, jsCode), header, "")
	if err == nil {
		var wxResp WxLoginResp
		json.Unmarshal([]byte(resp), &wxResp)
		//入库
		services.RegisterOrLogin(services.User{
			WxUnionId: wxResp.UnionId,
		})
		ctx.JSON(http.StatusOK, wxResp)
	} else {
		ctx.JSON(http.StatusInternalServerError, services.CommonError{
			Code:     10001,
			ErrorMsg: "wx login failed",
		})
	}
}

type WxLoginResp struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrorCode  int    `json:"errcode"`
	ErrorMsg   string `json:"errmsg"`
}

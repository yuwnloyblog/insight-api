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
		ctx.JSON(http.StatusBadRequest, services.GetError(services.ErrorCode_NoWxJsCode))
		return
	}

	appId := configures.Config.Wx.AppId
	secret := configures.Config.Wx.Secret

	header := map[string]string{}
	resp, err := tools.HttpDo("GET", fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appId, secret, jsCode), header, "")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, services.GetError(services.ErrorCode_WxLoginFail))
		return
	}
	var wxResp services.WxLoginResp
	err = json.Unmarshal([]byte(resp), &wxResp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, services.GetError(services.ErrorCode_WxLoginRespErr))
		return
	}
	//入库
	token, err := services.RegisterOrLogin(services.User{
		WxUnionId: wxResp.UnionId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, services.LoginUserResp{
		Token:  token,
		WxResp: &wxResp,
	})
}

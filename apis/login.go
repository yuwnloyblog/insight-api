package apis

import (
	"encoding/json"
	"fmt"
	"insight-api/configures"
	"insight-api/services"
	"insight-api/tools"
	"insight-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func WxLogin(ctx *gin.Context) {
	jsCode := ctx.PostForm("js_code")

	if jsCode == "" {
		ctx.JSON(http.StatusBadRequest, services.GetError(services.ErrorCode_NoWxJsCode))
		return
	}

	// phoenCode := ctx.PostForm("phone_code")

	appId := configures.Config.Wx.AppId
	secret := configures.Config.Wx.Secret
	fmt.Println("appid:", appId, "secret:", secret)
	header := map[string]string{}
	wxUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appId, secret, jsCode)
	resp, err := tools.HttpDo("GET", wxUrl, header, "")
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
	if wxResp.ErrorCode > 0 {
		fmt.Println("wxUrl:", wxUrl)
		ctx.JSON(http.StatusOK, services.LoginUserResp{
			Token:  "",
			WxResp: &wxResp,
		})
		return
	}
	//入库
	token, u, err := services.RegisterOrLogin(services.User{
		WxOpenid: wxResp.OpenId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.Writer.Header().Set("X-Status", strconv.Itoa(u.Status))
	ctx.JSON(http.StatusOK, services.LoginUserResp{
		Token:    token,
		NickName: u.NickName,
		Avatar:   u.Avatar,
		Status:   u.Status,
		WxResp:   &wxResp,
	})
}

func print(a interface{}) string {
	ret, _ := json.Marshal(a)
	return string(ret)
}

func HandleToken(ctx *gin.Context) {
	isLogined := "0"
	token := ctx.Request.Header.Get("X-Token")
	if token != "" {
		uidStr, _, err := services.ParseToken(token)
		if err != nil {
			ctx.JSON(http.StatusForbidden, services.GetError(services.ErrorCode_TokenErr))
			ctx.Abort()
			return
		}
		uid, err := utils.Decode(uidStr)
		if err != nil {
			ctx.JSON(http.StatusForbidden, services.GetError(services.ErrorCode_UidStrError))
			ctx.Abort()
			return
		}
		isLogined = "1"
		user, err := services.GetUserInfByCache(uid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, services.GetError(services.ErrorCode_UserDbReadFail))
			ctx.Abort()
			return
		}
		ctx.Set("uid", uid)
		ctx.Set("status", user.Status)
		fmt.Println("uid:", uid, "status:", user.Status)

		ctx.Writer.Header().Set("X-Status", strconv.Itoa(user.Status))
	}
	ctx.Writer.Header().Set("X-IsLogined", isLogined)
}

func checkLogin(ctx *gin.Context) bool {
	uid := ctx.GetInt64("uid")
	if uid <= 0 {
		ctx.JSON(http.StatusUnauthorized, services.GetError(services.ErrorCode_NotLogin))
		return false
	}
	return true
}

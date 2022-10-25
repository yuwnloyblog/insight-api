package apis

import (
	"fmt"
	"insight-api/services"
	"insight-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WxPayNotify(ctx *gin.Context) {
	json, _ := ctx.GetRawData()
	err := services.HandlePayNotify(json)
	fmt.Println("PayErr:", err)
	ctx.JSON(http.StatusOK, services.PayNotifyResp{
		Code:    "SUCCESS",
		Message: "",
	})
}

func WxPay(ctx *gin.Context) {
	if !checkLogin(ctx) {
		return
	}

	amountStr := ctx.PostForm("amount")
	amountInt, err := utils.ParseInt(amountStr)
	if err != nil {
		amountInt = 1
	}

	typeStr := ctx.PostForm("fellow_type")
	fllowType, err := utils.ParseInt(typeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, services.GetError(services.ErrorCode_ParamErr))
		return
	}

	uid := ctx.GetInt64("uid")
	user, err := services.GetUserInfByCache(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, services.GetError(services.ErrorCode_UserDbReadFail))
		return
	}
	resp, _, err := services.WxPayCall(fllowType, amountInt, user.WxOpenid, uid)
	if err == nil {
		ctx.JSON(http.StatusOK, resp)
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, services.GetError(services.ErrorCode_PrepayCallErr))
		return
	}
}

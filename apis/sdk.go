package apis

import (
	"insight-api/dbs"
	"insight-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SdkList(ctx *gin.Context) {
	appIdStr := ctx.Query("app_id")

	if !checkLogin(ctx) {
		return
	}

	status := ctx.GetInt("status")
	if status != dbs.UserStatus_YEAR_PAY && status != dbs.UserStatus_HALFYEAR_PAY && status != dbs.UserStatus_SEASON_PAY && status != dbs.UserStatus_MONTH_PAY {
		ctx.JSON(http.StatusForbidden, services.GetError(services.ErrorCode_NeedPay))
		return
	}

	sdks := services.QuerySdksByAppIdStr(appIdStr)
	ctx.JSON(http.StatusOK, &services.Sdks{
		Items: sdks,
	})
}

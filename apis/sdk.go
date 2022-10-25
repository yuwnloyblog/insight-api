package apis

import (
	"insight-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SdkList(ctx *gin.Context) {
	appIdStr := ctx.Query("app_id")

	if !checkLogin(ctx) {
		return
	}

	if !checkPay(ctx) {
		return
	}

	sdks := services.QuerySdksByAppId(appIdStr)
	ctx.JSON(http.StatusOK, &services.Sdks{
		Items: sdks,
	})
}

func checkFreeCount(ctx *gin.Context) bool {
	uid := ctx.GetInt64("uid")
	if uid <= 0 {
		return false
	}
	return services.CheckFreeCount(uid)
}

package apis

import (
	"insight-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SdkList(ctx *gin.Context) {
	appIdStr := ctx.Query("app_id")

	sdks := services.QuerySdksByAppIdStr(appIdStr)
	ctx.JSON(http.StatusOK, &services.Sdks{
		Items: sdks,
	})
}

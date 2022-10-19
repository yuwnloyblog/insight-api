package apis

import (
	"insight-api/services"
	"insight-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AppList(ctx *gin.Context) {
	startStr := ctx.Query("start")
	devId := ctx.Query("dev_id")

	countStr := ctx.Query("count")
	count, err := utils.ParseInt(countStr)
	if err != nil {
		count = 50
	} else {
		if count > 50 || count <= 0 {
			count = 50
		}
	}
	keyword := ctx.Query("keyword")

	if startStr != "" && !checkLogin(ctx) {
		return
	}

	ctx.JSON(http.StatusOK, services.QueryApps(keyword, startStr, devId, count))
}

func AppInfo(ctx *gin.Context) {
	if !checkLogin(ctx) {
		return
	}
	appIdStr := ctx.Query("id")
	app := services.GetAppByIdStr(appIdStr)
	ctx.JSON(http.StatusOK, app)
}

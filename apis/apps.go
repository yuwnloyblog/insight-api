package apis

import (
	"insight-api/services"
	"insight-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AppList(ctx *gin.Context) {
	startStr := ctx.Query("start")

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

	ctx.JSON(http.StatusOK, services.QueryApps(keyword, startStr, count))
}

func AppInfo(ctx *gin.Context) {
	appIdStr := ctx.Query("id")
	app := services.GetAppByIdStr(appIdStr)
	ctx.JSON(http.StatusOK, app)
}

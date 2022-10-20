package apis

import (
	"insight-api/services"
	"insight-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AppList(ctx *gin.Context) {
	pageStr := ctx.Query("page")
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
	page, err := utils.ParseInt(pageStr)
	if err != nil {
		page = 1
	} else {
		if page <= 0 {
			page = 1
		}
	}

	keyword := ctx.Query("keyword")

	if page > 1 && !checkLogin(ctx) {
		return
	}

	if devId != "" && !checkPay(ctx) {
		return
	}

	ctx.JSON(http.StatusOK, services.QueryApps(keyword, devId, page, count))
}

func AppInfo(ctx *gin.Context) {
	if !checkLogin(ctx) {
		return
	}
	appIdStr := ctx.Query("id")
	app := services.GetAppByIdStr(appIdStr)
	ctx.JSON(http.StatusOK, app)
}

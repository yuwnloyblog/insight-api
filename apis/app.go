package apis

import (
	"insight-api/services"
	"insight-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AppList(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	devId := ctx.Query("dev_id")

	if devId != "" {
		devId = DecodeUuid(devId)
	}

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
	retApps := services.QueryAppInfos(keyword, devId, page, count)
	if retApps != nil && len(retApps.Items) > 0 {
		for _, app := range retApps.Items {
			if app.Developer.Id != "" {
				app.Developer.Id = EncodeUuid(app.Developer.Id)
			}
			if !strings.HasPrefix(app.LogoUrl, "https://file.lwoowl.cn") && !strings.HasPrefix(app.LogoUrl, "https://pp.myapp.com") {
				app.LogoUrl = ""
			}
		}
	}

	ctx.JSON(http.StatusOK, retApps)
}

func AppInfo(ctx *gin.Context) {
	if !checkLogin(ctx) {
		return
	}
	appIdStr := ctx.Query("id")
	appMap := services.GetAppByIdStr(appIdStr)

	for _, app := range appMap {
		if app.Developer != nil && app.Developer.Id != "" {
			app.Developer.Id = EncodeUuid(app.Developer.Id)
		}
		if !strings.HasPrefix(app.LogoUrl, "https://file.lwoowl.cn") && !strings.HasPrefix(app.LogoUrl, "https://pp.myapp.com") {
			app.LogoUrl = ""
		}
	}

	ctx.JSON(http.StatusOK, appMap)
}

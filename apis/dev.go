package apis

import (
	"insight-api/services"
	"insight-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeveloperList(ctx *gin.Context) {
	pageStr := ctx.Query("page")

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
	if err != nil || page <= 0 {
		page = 1
	}
	keyword := ctx.Query("keyword")

	if page > 1 && !checkLogin(ctx) {
		return
	}
	devs := services.QueryDevelopers(keyword, page, count)
	if devs != nil && len(devs.Items) > 0 {
		for _, dev := range devs.Items {
			if dev.Id != "" {
				dev.Id = EncodeUuid(dev.Id)
			}
		}
	}
	ctx.JSON(http.StatusOK, devs)
}

func DeveloperInfo(ctx *gin.Context) {
	if !checkLogin(ctx) {
		return
	}
	devIdStr := ctx.Query("id")
	if devIdStr != "" {
		devIdStr = DecodeUuid(devIdStr)
	}
	devloper := services.GetDeveloperById(devIdStr, "")
	if devloper != nil && devloper.Id != "" {
		devloper.Id = EncodeUuid(devloper.Id)
	}
	ctx.JSON(http.StatusOK, devloper)
}

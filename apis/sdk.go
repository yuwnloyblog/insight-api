package apis

import (
	"insight-api/services"
	"insight-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SdkList(ctx *gin.Context) {
	var sdks *services.Sdks
	appIdStr := ctx.Query("app_id")
	if appIdStr != "" {
		if !checkLogin(ctx) {
			return
		}

		if !checkPay(ctx) {
			return
		}
		appId, err := utils.Decode(appIdStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, services.GetError(services.ErrorCode_ParamErr))
			return
		}

		sdks = services.QuerySdksByAppId(appId)
	} else {
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
		if err != nil {
			page = 1
		} else {
			if page <= 0 {
				page = 1
			}
		}
		sdks = services.QuerySdksByPage(page, count)
	}

	if sdks != nil && len(sdks.Items) > 0 {
		for _, sdk := range sdks.Items {
			if sdk.Developer != nil && sdk.Developer.Id != "" {
				sdk.Developer.Id = EncodeUuid(sdk.Developer.Id)
			}
		}
	}
	ctx.JSON(http.StatusOK, sdks)
}

func checkFreeCount(ctx *gin.Context) bool {
	uid := ctx.GetInt64("uid")
	if uid <= 0 {
		return false
	}
	return services.CheckFreeCount(uid)
}

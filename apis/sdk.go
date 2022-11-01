package apis

import (
	"insight-api/services"
	"insight-api/utils"
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

	appId, err := utils.Decode(appIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, services.GetError(services.ErrorCode_ParamErr))
		return
	}

	sdks := services.QuerySdksByAppId(appId)
	if len(sdks) > 0 {
		for _, sdk := range sdks {
			if sdk.Developer != nil && sdk.Developer.Id != "" {
				sdk.Developer.Id = EncodeUuid(sdk.Developer.Id)
			}
		}
	}
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

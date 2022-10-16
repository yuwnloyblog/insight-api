package apis

import (
	"insight-api/services"
	"insight-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserInfoUpdate(ctx *gin.Context) {
	avator := ctx.PostForm("avator")
	nickname := ctx.PostForm("nickname")

	uid := ctx.GetInt64("uid")
	if uid <= 0 {
		ctx.JSON(http.StatusInternalServerError, services.GetError(services.ErrorCode_UidStrError))
		return
	}
	err := services.UpdateUserInfo(avator, nickname, uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, services.GetError(services.ErrorCode_UserDbUpdateFail))
		return
	}
	ctx.JSON(http.StatusOK, services.GetSuccess())
}

func GetUserInf(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, services.GetError(services.ErrorCode_NoUid))
		return
	}
	id, err := utils.Decode(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, services.GetError(services.ErrorCode_UidStrError))
		return
	}
	user, err := services.GetUserInfo(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, services.GetError(services.ErrorCode_UserDbReadFail))
		return
	}
	ctx.JSON(http.StatusOK, user)
}

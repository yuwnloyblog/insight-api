package apis

import (
	"insight-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostFeedback(ctx *gin.Context) {
	if !checkLogin(ctx) {
		return
	}
	uid := ctx.GetInt64("uid")
	feedback := ctx.PostForm("feedback")
	if feedback != "" && uid > 0 {
		services.PostFeedback(uid, feedback)
	}
	ctx.JSON(http.StatusOK, services.GetSuccess())
}

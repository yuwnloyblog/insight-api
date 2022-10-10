package apis

import (
	"insight-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeInfo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, services.GetHomeInfo())
}

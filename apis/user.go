package apis

import (
	"insight-api/services"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {

}

func HandleToken(ctx *gin.Context) {
	token := ctx.Request.Header.Get("X-Token")
	uid, _, err := services.ParseToken(token)
	if err != nil {
		// ctx.JSON(http.StatusUnauthorized, &services.CommonError{
		// 	Code:     10002,
		// 	ErrorMsg: "not login",
		// })
		// ctx.Abort()
		// return
	} else {
		ctx.Set("uid", uid)
	}
}

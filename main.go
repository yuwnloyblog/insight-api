package main

import (
	"insight-api/apis"
	"insight-api/configures"
	"insight-api/dbs"
	"insight-api/logs"
	"insight-api/tools"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	configures.InitConfigures()
	logs.InitLogs()
	dbs.InitMysql()

	// // router()
	tools.ReloadAppPic(147542)
	// tools.ReplaceIcon4AppInfo(1)

}
func router() {
	r := gin.Default()
	r.Use(apis.HandleToken)
	// r.Group("/", apis.HandleToken)
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello")
	})
	r.GET("/homeinfo", apis.HomeInfo)

	r.GET("/applist", apis.AppList)
	r.GET("/app/list", apis.AppList)
	r.GET("/appinfo", apis.AppInfo)
	r.GET("/app/info", apis.AppInfo)

	r.GET("/devlist", apis.DeveloperList)
	r.GET("/dev/list", apis.DeveloperList)
	r.GET("/devinfo", apis.DeveloperInfo)
	r.GET("/dev/info", apis.DeveloperInfo)

	r.GET("/sdklist", apis.SdkList)
	r.GET("/sdk/list", apis.SdkList)

	//user
	r.POST("/user/update", apis.UserInfoUpdate)
	r.GET("/user/info", apis.GetUserInfo)
	r.POST("/user/update_status", apis.UpdatePayStatus)

	r.POST("/wx_login", apis.WxLogin)
	r.POST("/user/wx_login", apis.WxLogin)
	r.POST("/user/wx_pay", apis.WxPay)
	r.POST("/user/wx_pay_notify", apis.WxPayNotify)

	r.Run(":8080")
}

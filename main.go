package main

import (
	"insight-api/apis"
	"insight-api/tools"

	"github.com/gin-gonic/gin"
)

func main() {
	// configures.InitConfigures()
	// dbs.InitMysql()

	//router()

	//tools.CatchDevelopers()
	tools.CatchSdks()
}

func router() {
	r := gin.Default()
	r.Group("/", apis.HandleToken)
	r.GET("/homeinfo", apis.HomeInfo)
	r.GET("/applist", apis.AppList)
	r.GET("/appinfo", apis.AppInfo)
	r.GET("/devlist", apis.DeveloperList)
	r.GET("/devinfo", apis.DeveloperInfo)
	r.GET("/sdklist", apis.SdkList)

	r.POST("/wx_login", apis.WxLogin)

	r.Run(":8080")
}

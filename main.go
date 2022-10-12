package main

import (
	"insight-api/apis"
	"insight-api/configures"
	"insight-api/dbs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	configures.InitConfigures()
	dbs.InitMysql()

	router()
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
	r.GET("/appinfo", apis.AppInfo)
	r.GET("/devlist", apis.DeveloperList)
	r.GET("/devinfo", apis.DeveloperInfo)
	r.GET("/sdklist", apis.SdkList)

	r.POST("/wx_login", apis.WxLogin)

	r.Run(":8080")
}

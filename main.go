package main

import (
	"fmt"
	"insight-api/apis"
	"insight-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// configures.InitConfigures()
	// logs.InitLogs()
	// dbs.InitMysql()

	// router()
	us := "da73f3aae4a24c4ba1628bd04e64599b"
	fmt.Println(us)
	val, _ := utils.PruneUuid(us)
	fmt.Println(val)
	v, err := utils.Parse2Uuid(val)
	fmt.Println(v, err)
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

	r.Run(":8080")
}

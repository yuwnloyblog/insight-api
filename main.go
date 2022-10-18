package main

import (
	"encoding/hex"
	"fmt"
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

	//router()
	// tools.CheckDev()
	// tools.PureDev()
	// name, err := tools.DownloadPicture("https://huawei-icon.oss-cn-hangzhou.aliyuncs.com/application/icon144/70cef5dad83849f8ba75ef7031f12c37.png")
	// fmt.Println(name, err)

	//tools.QiniuUpload("")
	//tools.DeleteFile("70cef5dad83849f8ba75ef7031f12c37.png")
	ret, err := tools.QiniuFetch("https://file.lwoowl.cn/devs/qiniu-x.png")
	fmt.Println(ret, err)

	dst := make([]byte, 16)
	src := []byte("11aabbccddeeff223344556677889900")
	r, e := hex.Decode(dst, src)
	fmt.Println(r, e, dst)
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

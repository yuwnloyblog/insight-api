package main

import (
	"fmt"
	"insight-api/apis"
	"insight-api/configures"
	"insight-api/dbs"
	"insight-api/logs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	configures.InitConfigures()
	logs.InitLogs()
	dbs.InitMysql()

	// router()
	// tools.ReloadAppPic(354417) //356940
	// id, err := utils.Decode("uKflLV6tU")
	// fmt.Println(id, err)
	// sdks := services.QuerySdksByAppId(id)
	// fmt.Println(len(sdks))
	appInfoDao := dbs.AppInfoDao{}
	start := time.Now()
	apps, err := appInfoDao.QueryListByPage("", "", "com.getui.msgpushx", "", 1, 10)
	fmt.Println(time.Now().Sub(start))
	fmt.Println(len(apps), err)
	// if err == nil {
	// 	for _, app := range apps {
	// 		fmt.Println(app)
	// 	}
	// }
}

func router() {
	r := gin.Default()
	r.Use(apis.HandleToken)
	// r.Group("/", apis.HandleToken)
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello")
	})
	r.GET("/homeinfo", apis.HomeInfo)

	r.GET("/app/list", apis.AppList)
	r.GET("/app/info", apis.AppInfo)

	r.GET("/dev/list", apis.DeveloperList)
	r.GET("/dev/info", apis.DeveloperInfo)

	r.GET("/sdk/list", apis.SdkList)

	//user
	r.POST("/user/update", apis.UserInfoUpdate)
	r.GET("/user/info", apis.GetUserInfo)
	r.POST("/user/update_status", apis.UpdatePayStatus)

	r.POST("/wx_login", apis.WxLogin)
	r.POST("/user/wx_login", apis.WxLogin)
	r.POST("/user/wx_pay", apis.WxPay)
	r.POST("/user/wx_pay_notify", apis.WxPayNotify)

	r.POST("/feedback/add", apis.PostFeedback)

	r.Run(":8080")
}

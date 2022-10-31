package tools

import (
	"fmt"
	"insight-api/dbs"
	"strings"
	"time"
)

func GenAppSdkRel() {
	appDao := dbs.AppDao{}
	appSdkRelDao := dbs.AppSdkRelDao{}
	start := int64(0)
	for {
		list, err := appDao.QueryList("", "", start, 1000)
		if err == nil && len(list) > 0 {
			for _, app := range list {
				start = app.ID
				if app.SdkUids != "" {
					sdkIds := strings.Split(app.SdkUids, ",")
					for _, sdkId := range sdkIds {
						err = appSdkRelDao.Create(dbs.AppSdkRelDao{
							AppId: app.ID,
							SdkId: sdkId,
						})
						if err != nil {
							fmt.Println("Fail", app.ID, sdkId)
						} else {
							fmt.Println("Success", app.ID, sdkId)
						}
						time.Sleep(5 * time.Millisecond)
					}
				} else {
					fmt.Println("Omit", app.ID)
				}
			}
		} else {
			fmt.Println("Finish!!!")
		}
	}
}

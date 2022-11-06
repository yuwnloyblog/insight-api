package tools

import (
	"fmt"
	"insight-api/dbs"
	"strings"
	"time"
)

func AddDowncount4Rel(s string) {
	appinfoDao := dbs.AppInfoDao{}
	appSdkDao := dbs.AppSdkRelDao{}
	start := s
	for {
		list, err := appinfoDao.QueryList(start, 100)
		if err == nil && len(list) > 0 {
			for _, app := range list {
				start = app.Id
				if app.DownloadCount > 0 {
					err = appSdkDao.UpdateByBundleId(app.Id, app.DownloadCount)
					if err == nil {
						fmt.Println("success ", app.Id)
					} else {
						fmt.Println("fail ", app.Id, err)
					}

					time.Sleep(5 * time.Millisecond)
				}
			}
		} else {
			fmt.Println("Finish!!!")
			break
		}
	}
}

func AddBundle4SdkRel(s int64) {
	appDao := dbs.AppDao{}
	appSdkRelDao := dbs.AppSdkRelDao{}
	start := s
	for {
		list, err := appDao.QueryList("", "", start, 1000)
		if err == nil && len(list) > 0 {
			for _, app := range list {
				start = app.ID
				if app.BundleId != "" {
					err = appSdkRelDao.UpdateByAppId(app.ID, app.BundleId)
					if err == nil {
						fmt.Println("success ", app.ID)
					} else {
						fmt.Println("Fail ", app.ID, err)
					}
					time.Sleep(5 * time.Millisecond)
				}
			}
		} else {
			fmt.Println("Finish!!!")
		}
	}
}

func GenAppSdkRel(s int64) {
	appDao := dbs.AppDao{}
	appSdkRelDao := dbs.AppSdkRelDao{}
	start := s
	for {
		list, err := appDao.QueryList("", "", start, 1000)
		if err == nil && len(list) > 0 {
			for _, app := range list {
				start = app.ID
				if app.SdkUids != "" {
					sdkIds := strings.Split(app.SdkUids, ",")

					if app.ID < 17050 {
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
						appSdkRelArr := []dbs.AppSdkRelDao{}
						for _, sdkId := range sdkIds {
							appSdkRelArr = append(appSdkRelArr, dbs.AppSdkRelDao{
								AppId: app.ID,
								SdkId: sdkId,
							})
						}
						if len(appSdkRelArr) > 0 {
							err = appSdkRelDao.BatchCreate(appSdkRelArr)
							if err != nil {
								fmt.Println("Fail", app.ID)
							} else {
								fmt.Println("Success", app.ID)
							}
						}

						time.Sleep(5 * time.Millisecond)
					}
				} else {
					fmt.Println("Omit", app.ID)
				}
			}
		} else {
			fmt.Println("Finish!!!")
			break
		}
	}
}

package tools

import (
	"fmt"
	"insight-api/dbs"
	"strings"
	"time"
)

func ReplaceIcon4AppInfo(startPage int) {
	appInfoDao := dbs.AppInfoDao{}
	pageIndex := startPage
	for {
		appinfos, err := appInfoDao.QueryListByPage("", "", pageIndex, 1000)
		if err == nil && len(appinfos) > 0 {
			for _, app := range appinfos {
				packageName := app.Id
				if packageName != "" {
					iconUrl, err := GetIconFromMyApp(packageName)
					if err == nil {
						appInfoDao.UpdateLogo(app.Id, iconUrl)
					} else {
						continue
					}
					time.Sleep(10 * time.Millisecond)
				}
			}
			pageIndex++
		} else {
			fmt.Println("Error_Quit", err)
			break
		}

	}
}

func ReplaceIcon4App(start int) {
	appDao := dbs.AppDao{}
	appInfoDao := dbs.AppInfoDao{}
	startIndex := start
	for {
		apps, err := appDao.QueryList("", "", int64(startIndex), 1000)
		if err == nil && len(apps) > 0 {
			for _, app := range apps {
				startIndex = int(app.ID)
				if app.LogoUrl == "" || strings.HasPrefix(app.LogoUrl, "https://file.lwoowl.cn") || strings.HasPrefix(app.LogoUrl, "https://pp.myapp.com") {
					fmt.Println("Omit_Logo ", app.ID, app.LogoUrl)
					continue
				}
				if app.Platform == "ios" {
					fmt.Println("Omit_Platform ", app.ID, app.Platform)
					continue
				}
				packageName := app.BundleId
				if packageName != "" {
					time.Sleep(5 * time.Millisecond)
					//从appinfo中查找
					appInfo, err := appInfoDao.FindById(app.BundleId)
					if err == nil && appInfo != nil {
						if strings.HasPrefix(appInfo.LogoUrl, "https://pp.myapp.com") {
							appDao.UpdateLogo(app.ID, appInfo.LogoUrl)
							fmt.Println(app.ID, app.BundleId, appInfo.LogoUrl)
							continue
						}
					}

					iconUrl, err := GetIconFromMyApp(packageName)
					if err == nil {
						appDao.UpdateLogo(app.ID, iconUrl)
						appInfoDao.UpdateLogo(app.BundleId, iconUrl)
						fmt.Println(app.ID, app.BundleId, iconUrl)
						continue
					}
				}
			}
		} else {
			fmt.Println("Error_Quite", err)
			break
		}
	}
}

package tools

import (
	"fmt"
	"insight-api/dbs"
	"math"
	"sort"
	"strings"
	"time"
)

func GenerateAppInfos() {
	appInfoDao := dbs.AppInfoDao{}
	appDao := dbs.AppDao{}
	start := int64(math.MaxInt32)
	for {
		apps, err := appDao.QueryList("", "", start, 100)
		if err == nil && len(apps) > 0 {
			for _, app := range apps {
				start = app.ID

				bundleId := app.BundleId
				appInfo := getAppInfoByBundleId(bundleId)
				if appInfo != nil {
					//更新
					f := UpdateAppInfo(app, appInfo)
					if f {
						fmt.Println("Update ", app.ID, appInfo.Id)
					} else {
						fmt.Println("Nothing ", app.ID, appInfo.Id)
					}
				} else {
					appInfoDao.Create(dbs.AppInfoDao{
						Id:                app.BundleId,
						Title:             app.Title,
						Website:           app.Website,
						Description:       app.Description,
						CreateTime:        app.CreateTime,
						LogoUrl:           app.LogoUrl,
						Category:          app.Category,
						LatestVersion:     app.LatestVersion,
						DownloadCount:     app.DownloadCount,
						DeveloperId:       app.DeveloperIdStr,
						DeveloperTitle:    app.DeveloperTitle,
						FirstReleaseDate:  app.FirstReleaseDate,
						LatestReleaseDate: app.LatestReleaseDate,
						Channels:          app.Channel,
					})
					fmt.Println("Add ", app.ID, app.BundleId)
				}
				time.Sleep(5 * time.Millisecond)
			}
		} else {
			break
		}
	}
}

func UpdateAppInfo(app *dbs.AppDao, appInfo *dbs.AppInfoDao) bool {
	if !strings.Contains(appInfo.Channels, app.Channel) {
		upd := map[string]interface{}{}

		appInfo.DownloadCount = appInfo.DownloadCount + app.DownloadCount
		upd["download_count"] = appInfo.DownloadCount

		channels := strings.Split(appInfo.Channels, ",")
		channels = append(channels, app.Channel)
		sort.Strings(channels)
		upd["channels"] = strings.Join(channels, ",")

		appInfoDao := dbs.AppInfoDao{}
		appInfoDao.Updates(appInfo.Id, upd)
		return true
	}
	return false
}

func getAppInfoByBundleId(bundleId string) *dbs.AppInfoDao {
	appInfoDao := dbs.AppInfoDao{}
	appInfo, err := appInfoDao.FindById(bundleId)
	if err != nil {
		return nil
	}
	return appInfo
}

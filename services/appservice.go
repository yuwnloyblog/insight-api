package services

import (
	"insight-api/dbs"
	"insight-api/utils"
	"math"
)

func QueryApps(keyword, startStr, devId string, count int) *Apps {
	appDao := dbs.AppDao{}
	start, err := utils.Decode(startStr)
	if err != nil {
		start = math.MaxInt32
	}
	appTables, err := appDao.QueryList(keyword, devId, start, count)
	apps := &Apps{}
	if err == nil {
		l := len(appTables)
		if l >= count {
			apps.HasMore = true
		}
		apps.Items = make([]*App, 0)
		for _, appTable := range appTables {
			isStr, _ := utils.Encode(appTable.ID)
			apps.Items = append(apps.Items, &App{
				Id:            isStr,
				Title:         appTable.Title,
				BundleId:      appTable.BundleId,
				Platform:      appTable.Platform,
				Channel:       appTable.Channel,
				Website:       appTable.Website,
				Description:   appTable.Description,
				ReleaseDate:   appTable.ReleaseDate.UnixMilli(),
				Developer:     GetDeveloperById(appTable.DeveloperIdStr, appTable.DeveloperTitle),
				Size:          appTable.Size,
				CreateTime:    appTable.CreateTime.UnixMilli(),
				LogoUrl:       appTable.LogoUrl,
				Category:      appTable.Category,
				LatestVersion: appTable.LatestVersion,
				CountryCode:   appTable.CountryCode,
			})
		}
	}
	return apps
}

func GetAppById(appId int64) *App {
	appDao := dbs.AppDao{}
	appdb, err := appDao.FindById(appId)
	if err == nil {
		idStr, _ := utils.Encode(appdb.ID)
		return &App{
			Id:            idStr,
			Title:         appdb.Title,
			BundleId:      appdb.BundleId,
			Platform:      appdb.Platform,
			Channel:       appdb.Channel,
			Website:       appdb.Website,
			Description:   appdb.Description,
			ReleaseDate:   appdb.ReleaseDate.UnixMilli(),
			Developer:     GetDeveloperById(appdb.DeveloperIdStr, appdb.DeveloperTitle),
			Size:          appdb.Size,
			CreateTime:    appdb.CreateTime.UnixMilli(),
			LogoUrl:       appdb.LogoUrl,
			Category:      appdb.Category,
			LatestVersion: appdb.LatestVersion,
			CountryCode:   appdb.CountryCode,

			Sdks: QuerySdksByAppId(appId),
		}
	}
	return nil
}

func GetAppByIdStr(appIdStr string) *App {
	appId, err := utils.Decode(appIdStr)
	if err == nil {
		return GetAppById(appId)
	}
	return nil
}

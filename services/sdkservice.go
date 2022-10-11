package services

import (
	"insight-api/dbs"
	"insight-api/utils"
)

func QuerySdksByAppId(appId int64) []*SdkInfo {
	sdkDao := dbs.SdkDao{}
	sdkdbs, err := sdkDao.QueryList(appId)
	if err == nil && len(sdkdbs) > 0 {
		sdkInfos := make([]*SdkInfo, 0)
		for _, sdkdb := range sdkdbs {
			sdkInfos = append(sdkInfos, &SdkInfo{
				Id:        sdkdb.ID,
				Name:      sdkdb.Name,
				Platforms: sdkdb.Platforms,
				Category:  sdkdb.Category,
				Developer: &Developer{
					Id:    sdkdb.DeveloperId,
					Title: sdkdb.DeveloperName,
				},
				LogoUrl: sdkdb.LogoUrl,
			})
		}
		return sdkInfos
	}
	return []*SdkInfo{}
}

func QuerySdksByAppIdStr(appIdStr string) []*SdkInfo {
	appId, err := utils.Decode(appIdStr)
	if err == nil {
		return QuerySdksByAppId(appId)
	}
	return []*SdkInfo{}
}

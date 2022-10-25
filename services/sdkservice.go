package services

import (
	"insight-api/dbs"
)

func QuerySdksByAppId(appId string) []*SdkInfo {
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
					Title: sdkdb.DeveloperName,
				},
				LogoUrl: sdkdb.LogoUrl,
			})
		}
		return sdkInfos
	}
	return []*SdkInfo{}
}

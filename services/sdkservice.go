package services

import (
	"insight-api/dbs"
)

func QuerySdksByPage(page, count int) *Sdks {
	sdkDao := dbs.SdkDao{}
	sdkdbs, err := sdkDao.QueryListByPage(page, count)
	if err == nil && len(sdkdbs) > 0 {
		sdkInfos := []*SdkInfo{}
		for _, sdkdb := range sdkdbs {
			sdkInfos = append(sdkInfos, &SdkInfo{
				Id:        sdkdb.ID,
				Name:      sdkdb.Title,
				Platforms: sdkdb.Platforms,
				Category:  sdkdb.Category,
				Developer: &Developer{
					Title: sdkdb.DeveloperName,
				},
				LogoUrl:  sdkdb.LogoUrl,
				LogoUrl2: sdkdb.LogoUrl,
			})
			return &Sdks{
				Items: sdkInfos,
				PageInfo: &PageInfo{
					Page:  page,
					Count: count,
				},
			}
		}
	}
	return &Sdks{}
}
func QuerySdksByAppId(appId int64) *Sdks {
	sdkDao := dbs.SdkDao{}
	sdkdbs, err := sdkDao.QueryList(appId)
	if err == nil && len(sdkdbs) > 0 {
		sdkInfos := make([]*SdkInfo, 0)
		for _, sdkdb := range sdkdbs {
			sdkInfos = append(sdkInfos, &SdkInfo{
				Id:        sdkdb.ID,
				Name:      sdkdb.Title,
				Platforms: sdkdb.Platforms,
				Category:  sdkdb.Category,
				Developer: &Developer{
					Title: sdkdb.DeveloperName,
				},
				LogoUrl:  sdkdb.LogoUrl,
				LogoUrl2: sdkdb.LogoUrl,
			})
		}
		return &Sdks{Items: sdkInfos}
	}
	return &Sdks{}
}

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
			idStr, _ := utils.Encode(sdkdb.ID)
			sdkInfos = append(sdkInfos, &SdkInfo{
				Id:          idStr,
				Name:        sdkdb.Name,
				BundleId:    sdkdb.BundleId,
				Description: sdkdb.Description,
				LogoUrl:     sdkdb.LogoUrl,
				Classify:    sdkdb.Classify,
				Developer:   GetDeveloperById(sdkdb.DeveloperId),
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

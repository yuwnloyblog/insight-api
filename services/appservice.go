package services

import (
	"encoding/json"
	"insight-api/dbs"
	"insight-api/utils"
)

func QueryAppInfos(keyword, devId, sdkId, notSdkId string, page, count int) *Apps {
	appInfoDao := dbs.AppInfoDao{}
	appInfos, err := appInfoDao.QueryListByPage(keyword, devId, sdkId, notSdkId, page, count)
	apps := &Apps{
		PageInfo: &PageInfo{
			Page:  page,
			Count: count,
		},
	}
	if err == nil {
		l := len(appInfos)
		if l >= count {
			apps.HasMore = true
		}
		apps.Items = make([]*App, 0)
		for _, appInfo := range appInfos {
			apps.Items = append(apps.Items, &App{
				Id:          appInfo.Id,
				Title:       appInfo.Title,
				Channel:     appInfo.Channels,
				Website:     appInfo.Website,
				Description: appInfo.Description,
				Developer: &Developer{
					Title: appInfo.DeveloperTitle,
				},
				CreateTime:    appInfo.CreateTime.UnixMilli(),
				LogoUrl:       appInfo.LogoUrl,
				Category:      appInfo.Category,
				LatestVersion: appInfo.LatestVersion,
			})
		}
	}
	return apps
}
func QueryAppInfosByIndex(keyword string, page, count int) *Apps {
	query := IndexQuery{
		Query: keyword,
		Page:  page,
		Limit: count,
	}
	apps := &Apps{
		PageInfo: &PageInfo{
			Page:  page,
			Count: count,
		},
	}
	bs, _ := json.Marshal(query)
	headers := map[string]string{}
	headers["Content-Type"] = "application/json"
	body := string(bs)
	ret, err := utils.HttpDo("POST", "http://127.0.0.1:5678/api/query?database=appinfos", headers, body)
	if err == nil {
		var resp IndexResp
		err = json.Unmarshal([]byte(ret), &resp)
		if err == nil && resp.Data != nil {
			apps.Items = []*App{}
			for _, doc := range resp.Data.Documents {
				if doc.Document != nil {
					doc.Document.Title = doc.Text
					apps.Items = append(apps.Items, doc.Document)
				}
			}
		}
	}
	return apps
}

func GetAppByIdStr(appIdStr string) map[string]*App {
	appDao := dbs.AppDao{}
	apps, err := appDao.FindByBundleId(appIdStr)
	if err != nil || len(apps) <= 0 {
		return map[string]*App{}
	}
	ret := map[string]*App{}
	for _, app := range apps {
		idStr, _ := utils.Encode(app.ID)
		appItem := &App{
			Id:                idStr,
			Title:             app.Title,
			BundleId:          app.BundleId,
			Platform:          app.Platform,
			Channel:           app.Channel,
			Website:           app.Website,
			Description:       app.Description,
			Developer:         GetDeveloperById(app.DeveloperIdStr, app.DeveloperTitle),
			Size:              app.Size,
			CreateTime:        app.CreateTime.UnixMilli(),
			LogoUrl:           app.LogoUrl,
			Category:          app.Category,
			LatestVersion:     app.LatestVersion,
			CountryCode:       app.CountryCode,
			LatestReleaseDate: app.LatestReleaseDate,
			FirstReleaseDate:  app.FirstReleaseDate,
		}
		if app.SdkUids != "" {
			ret[app.Channel] = appItem
		}
	}
	return ret
}

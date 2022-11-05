package tools

import (
	"encoding/json"
	"fmt"
	"insight-api/dbs"
	"insight-api/services"
	"insight-api/utils"
	"time"
)

func AddAppIndex(s string, si int64) {
	start := s
	sIndex := si
	appInfoDao := dbs.AppInfoDao{}

	for {
		apps, err := appInfoDao.QueryList(start, 10)
		fmt.Println("xxx:", err)
		if err == nil && len(apps) > 0 {
			indexApp := []AppIndex{}
			for _, app := range apps {
				index := AppIndex{
					Id:   sIndex,
					Text: app.Title,
					Document: services.App{
						Id:          app.Id,
						Channel:     app.Channels,
						Website:     app.Website,
						Description: app.Description,
						Developer: &services.Developer{
							Title: app.DeveloperTitle,
						},
						CreateTime:    app.CreateTime.UnixMilli(),
						LogoUrl:       app.LogoUrl,
						Category:      app.Category,
						LatestVersion: app.LatestVersion,
					},
				}
				indexApp = append(indexApp, index)
				start = app.Id
				sIndex++
			}
			if len(indexApp) > 0 {
				CreateIndex(indexApp)
				//fmt.Println(len(indexApp))
			}
			time.Sleep(5 * time.Millisecond)
		} else {
			fmt.Println("Finish!!!")
			break
		}
	}
}

func CreateIndex(indexArr []AppIndex) error {
	headers := map[string]string{}
	headers["Content-Type"] = "application/json"
	bs, err := json.Marshal(indexArr)
	if err != nil {
		return err
	}
	body := string(bs)
	ret, err := utils.HttpDo("POST", "http://127.0.0.1:5678/api/index/batch?database=appinfos", headers, body)
	if err != nil {
		return err
	}
	fmt.Println(indexArr[0].Id, ret, err)
	return nil
}

type AppIndex struct {
	Id       int64        `json:"id"`
	Text     string       `json:"text"`
	Document services.App `json:"document"`
}

/*
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
*/

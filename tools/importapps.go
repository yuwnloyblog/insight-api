package tools

import (
	"encoding/json"
	"fmt"
	"insight-api/dbs"
	"os"
	"strings"
	"time"
)

func ImportDeveloperFromFile() {
	content, err := os.ReadFile("/Users/bytedance/Documents/workspace/insight-api/tools/applist.json")
	if err != nil {
		panic(err)
	}

	var data AppCommonData
	err = json.Unmarshal(content, &data)

	if err == nil {
		devDao := dbs.DeveloperDao{}
		for _, appResp := range data.Data.Gross {
			devDao.Create(dbs.DeveloperDao{
				Title:       appResp.Publisher.Title,
				FoundedTime: time.Now(),
				CreateTime:  time.Now(),
			})
		}
	}
}

func ImportFromFile() {
	content, err := os.ReadFile("/Users/bytedance/Documents/workspace/insight-api/tools/applist.json")
	if err != nil {
		panic(err)
	}

	var data AppCommonData
	err = json.Unmarshal(content, &data)

	if err == nil {
		fmt.Println(data)
		appDao := dbs.AppDao{}
		for _, appResp := range data.Data.Gross {
			channel, platform, bundleId := ParseUid(appResp.Uid)
			appDao.Create(dbs.AppDao{
				Title:         appResp.Title,
				BundleId:      bundleId,
				Platform:      platform,
				Channel:       channel,
				Category:      appResp.Category,
				DownloadCount: int64(appResp.DownloadCount),
				ReleaseDate:   time.Now(),
				CreateTime:    time.Now(),
			})
		}
	} else {
		fmt.Println(err)
	}
}

func ParseUid(uid string) (channel, platform, bundleId string) {
	arr := strings.Split(uid, ":")
	if len(arr) == 3 {
		channel = arr[0]
		platform = arr[1]
		bundleId = arr[2]
	}
	return
}

type AppResp struct {
	Order         int       `json:"order"`
	Title         string    `json:"title"`
	LogoUrl       string    `json:"logoURL"`
	Channel       string    `json:"channel"`
	Uid           string    `json:"uid"`
	Category      string    `json:"category"`
	DownloadCount int       `json:"downloadCount"`
	Publisher     Publisher `json:"publisher"`
}

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
	NextPage   int `json:"nextPage"`
	TotalLimit int `json:"totalLimit"`
}
type AppListResp struct {
	Gross      []*AppResp `json:"gross"`
	UpdatedAt  string     `json:"updatedAt"`
	Pagination Pagination `json:"pagination"`
}

type AppCommonData struct {
	Data AppListResp `json:"data"`
}

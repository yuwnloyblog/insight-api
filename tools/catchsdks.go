package tools

import (
	"encoding/json"
	"fmt"
	"insight-api/dbs"
	"insight-api/utils"
	"strings"
	"time"
)

func CatchSdks() {
	header := map[string]string{}
	header["cookie"] = `_ga_MF5DDQ4TF9=GS1.1.1665500668.1.0.1665500668.0.0.0; _ga=GA1.1.390781912.1665500669; fork_session=935b65f690f18816; mp_f39218c184e602888e26ea53b00435dd_mixpanel={"distinct_id": "yuhongda0315@163.com","$device_id": "183c7936b51c35-05dc43c8eb10c2-19525635-13c680-183c7936b528a2","$initial_referrer": "$direct","$initial_referring_domain": "$direct","g_version": "1.20.0","$user_id": "yuhongda0315@163.com","g_team_name": "grtd"}; intercom-session-d1key3b8=ZDJPK1pNVlh4SkhiRTNrWGhnTjJzcXJISUJRbUtxdDFwOGk2Z2VhSDFLcThaWmExSVhIaWNSTjNqRlFBZGZCdS0tYjAvdkt2Q0RMYm1GeUFiaXI1Vm95UT09--436f32fdda0315f5905ca221319dc3b5dd672979; _ga_G812B88X1Y=GS1.1.1665500669.1.1.1665500693.0.0.0`
	header["content-type"] = "application/json"
	header["content-length"] = "39"

	index := 1
	sdkDao := dbs.SdkDao{}
	for {
		url := fmt.Sprintf("https://api.app.forkai.cn/webapi/sdks/search?page=%d", index)
		ret, err := utils.HttpDo("POST", url, header, ``)
		if err == nil {
			var data SdkResp
			err = json.Unmarshal([]byte(ret), &data)
			if err == nil {
				c := 0
				if len(data.Data.Services) > 0 {
					c = len(data.Data.Services)
					for _, sdk := range data.Data.Services {
						sdkDao.Create(dbs.SdkDao{
							ID:             sdk.Uid,
							Name:           sdk.Name,
							Platforms:      strings.Join(sdk.Platforms, ","),
							Category:       sdk.Category,
							DeveloperName:  sdk.PublisherName,
							DeveloperId:    sdk.PublisherUid,
							LogoUrl:        sdk.LogoUrl,
							AppCount:       sdk.InstalledProductCount,
							DeveloperCount: sdk.PublisherCount,
							WebsiteCount:   sdk.InstalledWebsiteCount,
						})
					}
					fmt.Println("page:", index, "count:", c)
				} else {
					break
				}
			}
		} else {
			fmt.Println("http err:", err)
		}
		index++
		time.Sleep(1 * time.Second)
	}
}

type SdkResp struct {
	Data *SdkCommonData `json:"data"`
}

type SdkCommonData struct {
	Total      int         `json:"total"`
	Services   []*Sdk      `json:"services"`
	Pagination *Pagination `json:"pagination"`
}
type Sdk struct {
	Uid                   string   `json:"uid"`
	Name                  string   `json:"name"`
	Platforms             []string `json:"platforms"`
	Category              string   `json:"category"`
	PublisherName         string   `json:"publisherName"`
	PublisherUid          string   `json:"publisherUID"`
	LogoUrl               string   `json:"logoURL"`
	InstalledProductCount int      `json:"installedProductCount"`
	PublisherCount        int      `json:"publisherCount"`
	InstalledWebsiteCount int      `json:"installedWebsiteCount"`
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
	NextPage int `json:"nextPage"`
}

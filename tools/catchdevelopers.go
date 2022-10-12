package tools

import (
	"encoding/json"
	"fmt"
	"insight-api/dbs"
	"strconv"
	"strings"
	"time"
)

var header map[string]string

func CatchDevelopers() {
	header := map[string]string{}
	header["cookie"] = `_ga_MF5DDQ4TF9=GS1.1.1665500668.1.0.1665500668.0.0.0; _ga=GA1.1.390781912.1665500669; fork_session=574227d113785cce; mp_f39218c184e602888e26ea53b00435dd_mixpanel={"distinct_id": "yuhongda0315@163.com","$device_id": "183c7f26f01ab7-0c84907d721d63-19525635-13c680-183c7f26f02c6e","$initial_referrer": "$direct","$initial_referring_domain": "$direct","$user_id": "yuhongda0315@163.com","g_team_name": "于洪达的个人团队"}; intercom-session-d1key3b8=aUNzektsaTBwSEhrczZOWm9PR2lQdzloMURSNzhGbmp6VWd6OVdHY1F5UEhTdWJPOXhLd2xsZnFmUTFJaGlJMi0taDdoTjh2bjgycmVROUVLQ0pSakxPUT09--fdb3b4eb2b924d5ae4ebe35332c0bbaa6d0f8882; _ga_G812B88X1Y=GS1.1.1665506619.2.1.1665506901.0.0.0`
	header["content-type"] = "application/json"
	header["content-length"] = "39"
	index := 1
	devDao := dbs.DeveloperDao{}
	for {
		url := fmt.Sprintf("https://api.app.forkai.cn/webapi/publishers/search?page=%d", index)
		ret, err := HttpDo("POST", url, header, "")
		if err == nil {
			var data DeveloperCommonData
			err = json.Unmarshal([]byte(ret), &data)
			if err == nil {
				c := 0
				if len(data.Data.Publishers) > 0 {
					c = len(data.Data.Publishers)
					for _, pub := range data.Data.Publishers {
						devDao.Create(dbs.DeveloperDao{
							ID:             pub.Uid,
							Title:          pub.Title,
							Industry:       pub.Industry,
							LogoUrl:        pub.LogoUrl,
							FoundedYear:    strconv.Itoa(pub.FoundedYear),
							AddressArea:    strings.Join(pub.AddressArea, "-"),
							FinancingRound: pub.FinancingRound,
							CreateTime:     time.Now(),
							AppCount:       pub.ProductCount,
							WebsiteCount:   pub.WebsiteCount,
							DownloadCount:  pub.DownloadCount,
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

type DeveloperCommonData struct {
	Data PublishersResp `json:"data"`
}

type PublishersResp struct {
	Publishers []*Publisher `json:"publishers"`
}

type Publisher struct {
	Uid            string   `json:"uid"`
	Title          string   `json:"title"`
	LogoUrl        string   `json:"logoURL"`
	Industry       string   `json:"industry"`
	FoundedYear    int      `json:"foundedYear"`
	AddressArea    []string `json:"addressArea"`
	FinancingRound string   `json:"financingRound"`
	ProductCount   int      `json:"productCount"`
	WebsiteCount   int      `json:"websiteCount"`
	DownloadCount  int64    `json:"downloadCount"`

	InstalledPublisherCount int `json:"installedPublisherCount"`
	ServicesCount           int `json:"servicesCount"`
	InstalledProductCount   int `json:"installedProductCount"`
	InstalledWebsiteCount   int `json:"installedWebsiteCount"`
}

package tools

import (
	"encoding/json"
	"fmt"
	"insight-api/dbs"
	"insight-api/utils"
	"strconv"
	"strings"
	"time"
)

func CatchServicers() {
	header := map[string]string{}
	header["cookie"] = `_ga_MF5DDQ4TF9=GS1.1.1665500668.1.0.1665500668.0.0.0; _ga=GA1.1.390781912.1665500669; fork_session=574227d113785cce; mp_f39218c184e602888e26ea53b00435dd_mixpanel={"distinct_id": "yuhongda0315@163.com","$device_id": "183c7f26f01ab7-0c84907d721d63-19525635-13c680-183c7f26f02c6e","$initial_referrer": "$direct","$initial_referring_domain": "$direct","$user_id": "yuhongda0315@163.com","g_team_name": "于洪达的个人团队"}; intercom-session-d1key3b8=aUNzektsaTBwSEhrczZOWm9PR2lQdzloMURSNzhGbmp6VWd6OVdHY1F5UEhTdWJPOXhLd2xsZnFmUTFJaGlJMi0taDdoTjh2bjgycmVROUVLQ0pSakxPUT09--fdb3b4eb2b924d5ae4ebe35332c0bbaa6d0f8882; _ga_G812B88X1Y=GS1.1.1665506619.2.1.1665506901.0.0.0`
	header["content-type"] = "application/json"
	header["content-length"] = "39"
	index := 1
	servicerDao := dbs.ServiceProviderDao{}
	for {
		url := fmt.Sprintf("https://api.app.forkai.cn/webapi/service-providers/search?page=%d", index)
		ret, err := utils.HttpDo("POST", url, header, "")
		if err == nil {
			var data ServiceProviderCommonData
			err = json.Unmarshal([]byte(ret), &data)
			if err == nil {
				c := 0
				if len(data.Data.Providers) > 0 {
					c = len(data.Data.Providers)
					for _, pub := range data.Data.Providers {
						servicerDao.Create(dbs.ServiceProviderDao{
							ID:             pub.Uid,
							Title:          pub.Title,
							Industry:       pub.Industry,
							LogoUrl:        pub.LogoUrl,
							FoundedYear:    strconv.Itoa(pub.FoundedYear),
							AddressArea:    strings.Join(pub.AddressArea, "-"),
							FinancingRound: pub.FinancingRound,
							CreateTime:     time.Now(),

							AppCount:       pub.InstalledProductCount,
							WebsiteCount:   pub.InstalledWebsiteCount,
							ServiceCount:   pub.ServicesCount,
							DeveloperCount: pub.InstalledPublisherCount,
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

func CatchServiceProvider(devId string, headers map[string]string) *Publisher {
	url := fmt.Sprintf("https://api.app.forkai.cn/webapi/service-providers/%s", devId)
	ret, err := utils.HttpDo("GET", url, headers, "")
	if err == nil {
		var providerResp ProviderResp
		err = json.Unmarshal([]byte(ret), &providerResp)
		if err == nil && providerResp.Data != nil && providerResp.Data.Provider != nil {
			return providerResp.Data.Provider
		}
	}
	return nil
}

type ProviderResp struct {
	Data *ProviderCommon `json:"data"`
}
type ProviderCommon struct {
	Provider *Publisher `json:"provider"`
}
type ServiceProviderCommonData struct {
	Data ServiceProviderResp `json:"data"`
}

type ServiceProviderResp struct {
	Providers []*Publisher `json:"providers"`
}

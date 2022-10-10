package tools

import (
	"encoding/json"
	"fmt"
	"insight-api/dbs"
	"strings"
	"time"
)

func CatchApps() {
	header2 := map[string]string{}
	header2["cookie"] = `_ga=GA1.1.1684076364.1665304345; fork_session=3f66ad476497a91d; mp_f39218c184e602888e26ea53b00435dd_mixpanel={"distinct_id": "yuhongda0315@163.com","$device_id": "183bccf8d0e1b17-01c76202b5b0f9-1a525635-1d73c0-183bccf8d0f17e9","$initial_referrer": "$direct","$initial_referring_domain": "$direct","g_version": "1.20.0","$user_id": "yuhongda0315@163.com","g_team_name": "grtd"}; _ga_G812B88X1Y=GS1.1.1665377700.5.1.1665378795.0.0.0; intercom-session-d1key3b8=Ky9BUmxneTZoOXNHck1DZjhwVksyNVZ0bHlFWTgyUXVOUlAvRTNzVmZLcEtZa2xPUTZZcnJMMVphTVZBYnZUaC0tZld0ME81cWdnZWJ5SmtjcjdsUlM0QT09--68f67736b41940e646a3b205874e449262fa030c`
	header2["content-type"] = "application/json"
	header2["content-length"] = "39"
	index := 1
	appDao := dbs.AppDao{}
	for {
		url := fmt.Sprintf("https://api.app.forkai.cn/webapi/products/query?page=%d", index)
		ret, err := HttpDo("POST", url, header2, `{"channel":"huawei","countryCode":"CN"}`)
		if err == nil {
			var data ProductsResp
			err = json.Unmarshal([]byte(ret), &data)
			if err == nil {
				c := 0
				if len(data.Data.Products) > 0 {
					c = len(data.Data.Products)
					for _, p := range data.Data.Products {

						devIdStr := ""
						devTitle := ""
						if p.Publisher != nil {
							devIdStr = p.Publisher.Uid
							devTitle = p.Publisher.Title
						}
						channel, platform, bundle := ParseUid(p.Uid)
						appDao.Create(dbs.AppDao{
							Title:             p.Title,
							LogoUrl:           p.LogoUrl,
							DeveloperIdStr:    devIdStr,
							DeveloperTitle:    devTitle,
							Channel:           channel,
							Platform:          platform,
							BundleId:          bundle,
							Uid:               p.Uid,
							Category:          p.Category,
							DownloadCount:     p.Latest1yDownloads,
							FirstReleaseDate:  p.FirstReleaseDate,
							LatestReleaseDate: p.LatestReleaseDate,
							Website:           p.Homepage,
							Size:              p.Size,
							Paid:              ParsePaid(p.Paid),
							LatestVersion:     p.Version,
							SdkUids:           strings.Join(p.SdkUids, ","),
							SdkDevs:           strings.Join(p.SdkProviders, ","),
							CloudServices:     strings.Join(p.CloudServiceUids, ","),
							CloudServiceDevs:  strings.Join(p.CloudServiceProviders, ","),

							ReleaseDate: time.Now(),
							CreateTime:  time.Now(),
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
func ParsePaid(paid bool) string {
	if paid {
		return "1"
	}
	return "0"
}

type ProductsResp struct {
	Data ProductCommonData `json:"data"`
}

type ProductCommonData struct {
	Total      int        `json:"total"`
	Products   []*Product `json:"products"`
	Pagination Pagination `json:"pagination"`
}

type Product struct {
	Title                 string     `json:"title"`
	LogoUrl               string     `json:"logoURL"`
	Publisher             *Publisher `json:"publisher"`
	Channel               string     `json:"channel"`
	Uid                   string     `json:"uid"`
	BundleId              string     `json:"bundleID"`
	Category              string     `json:"category"`
	Latest1yDownloads     int64      `json:"latest1yDownloads"`
	Score                 float64    `json:"score"`
	FirstReleaseDate      string     `json:"firstReleaseDate"`
	LatestReleaseDate     string     `json:"latestReleaseDate"`
	Homepage              string     `json:"homepage"`
	Size                  int64      `json:"size"`
	Paid                  bool       `json:"paid"`
	Version               string     `json:"version"`
	CountryCode           string     `json:"countryCode"`
	SdkUids               []string   `json:"sdkUIDs"`
	SdkProviders          []string   `json:"sdkProviders"`
	CloudServiceUids      []string   `json:"cloudServiceUIDs"`
	CloudServiceProviders []string   `json:"cloudServiceProviders"`
}

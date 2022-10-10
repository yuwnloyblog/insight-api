package tools

import (
	"encoding/json"
	"fmt"
	"insight-api/dbs"
	"strings"
	"time"
)

var header map[string]string

func init() {
	header := map[string]string{}
	header["cookie"] = `_ga=GA1.1.1684076364.1665304345; fork_session=3f66ad476497a91d; mp_f39218c184e602888e26ea53b00435dd_mixpanel={"distinct_id": "yuhongda0315@163.com","$device_id": "183bccf8d0e1b17-01c76202b5b0f9-1a525635-1d73c0-183bccf8d0f17e9","$initial_referrer": "$direct","$initial_referring_domain": "$direct","g_version": "1.20.0","$user_id": "yuhongda0315@163.com","g_team_name": "grtd"}; _ga_G812B88X1Y=GS1.1.1665377700.5.1.1665378795.0.0.0; intercom-session-d1key3b8=Ky9BUmxneTZoOXNHck1DZjhwVksyNVZ0bHlFWTgyUXVOUlAvRTNzVmZLcEtZa2xPUTZZcnJMMVphTVZBYnZUaC0tZld0ME81cWdnZWJ5SmtjcjdsUlM0QT09--68f67736b41940e646a3b205874e449262fa030c`
}

func CatchDevelopers() {
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
				fmt.Println("len:", len(data.Data.Publishers), ret)
				if len(data.Data.Publishers) > 0 {
					c = len(data.Data.Publishers)
					for _, pub := range data.Data.Publishers {
						devDao.Create(dbs.DeveloperDao{
							Title:          pub.Title,
							Trade:          pub.Industry,
							LogoUrl:        pub.LogoUrl,
							AddressArea:    strings.Join(pub.AddressArea, "-"),
							FoundedTime:    foundYearFormat(pub.FoundedYear),
							FinancingRound: pub.FinancingRound,
							CreateTime:     time.Now(),
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
		time.Sleep(5 * time.Second)
	}
}

func foundYearFormat(year int) time.Time {
	if year <= 0 {
		year = 1970
	}
	format := "2006-01-02"
	t, _ := time.Parse(format, fmt.Sprintf("%d-01-01", year))
	return t
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
}

package tools

import (
	"encoding/json"
	"fmt"
	"insight-api/dbs"
	"math"
	"strconv"
	"strings"
	"time"
)

func FillDev() {
	start := math.MaxInt64
	appdao := dbs.AppDao{}
	devdao := dbs.DeveloperDao{}
	for {
		apps, err := appdao.QueryList("", "", int64(start), 100)
		if err == nil && len(apps) > 0 {
			for _, app := range apps {
				start = int(app.ID)
				devId := app.DeveloperIdStr
				dev, _ := devdao.FindById(devId)
				if dev == nil {
					fmt.Println("dev_id:", devId)
					QueryDeveloper(devId)
					time.Sleep(200 * time.Millisecond)
				}
			}
		} else {
			break
		}
	}
}

func QueryDeveloperOnly(id string) {
	url := fmt.Sprintf("https://api.app.forkai.cn/webapi/publishers/%s", id)
	ret, err := HttpDo("GET", url, headers, "")
	if err == nil {
		fmt.Println(ret)
		var resp PublisherResp
		err = json.Unmarshal([]byte(ret), &resp)
		if err == nil && resp.Data != nil && resp.Data.Publisher != nil {
			d := resp.Data.Publisher
			bs, _ := json.Marshal(d)
			fmt.Println(string(bs))
		}
	}
}
func QueryDeveloper(id string) {
	url := fmt.Sprintf("https://api.app.forkai.cn/webapi/publishers/%s", id)
	ret, err := HttpDo("GET", url, headers, "")
	if err == nil {
		var resp PublisherResp
		err = json.Unmarshal([]byte(ret), &resp)
		if err == nil && resp.Data != nil && resp.Data.Publisher != nil {
			d := resp.Data.Publisher
			devDao := dbs.DeveloperDao{}
			devDao.Create(dbs.DeveloperDao{
				ID:          d.Uid,
				Title:       d.Title,
				LogoUrl:     d.LogoUrl,
				Industry:    d.Industry,
				FoundedYear: strconv.Itoa(d.FoundedYear),
				Address:     d.Address,
				Website:     d.Website,
				Email:       d.Email,
				Description: d.Description,
				AddressArea: strings.Join(d.AddressArea, "-"),
				CreateTime:  time.Now(),
			})
		}
	}
}

type PublisherResp struct {
	Data *PublisherCommonData `json:"data"`
}
type PublisherCommonData struct {
	Publisher *Publisher `json:"publisher"`
}

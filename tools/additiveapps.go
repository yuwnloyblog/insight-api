package tools

import (
	"encoding/json"
	"fmt"
	"insight-api/dbs"
	"insight-api/utils"
	"math/rand"
	"strings"
	"time"
)

func AdditiveAppsScanByChannel(channel, countryCode, updateScope string) {
	headers := map[string]string{}
	headers["cookie"] = `_ga=GA1.1.1684076364.1665304345; fork_session=3f66ad476497a91d; mp_f39218c184e602888e26ea53b00435dd_mixpanel={"distinct_id": "yuhongda0315@163.com","$device_id": "183bccf8d0e1b17-01c76202b5b0f9-1a525635-1d73c0-183bccf8d0f17e9","$initial_referrer": "$direct","$initial_referring_domain": "$direct","g_version": "1.20.0","$user_id": "yuhongda0315@163.com","g_team_name": "grtd"}; _ga_G812B88X1Y=GS1.1.1665377700.5.1.1665378795.0.0.0; intercom-session-d1key3b8=Ky9BUmxneTZoOXNHck1DZjhwVksyNVZ0bHlFWTgyUXVOUlAvRTNzVmZLcEtZa2xPUTZZcnJMMVphTVZBYnZUaC0tZld0ME81cWdnZWJ5SmtjcjdsUlM0QT09--68f67736b41940e646a3b205874e449262fa030c`
	headers["content-type"] = "application/json"
	rand.Seed(time.Now().Unix())
	pageIndex := 1

	for {
		body := fmt.Sprintf(`{
			"channel": "%s",
			"countryCode": "%s",
			"updateDateStart": "%s",
			"updateDateEnd": "now"
		  }`, channel, countryCode, updateScope)
		url := fmt.Sprintf("https://api.app.forkai.cn/webapi/products/query?page=%d", pageIndex)
		ret, err := utils.HttpDo("POST", url, headers, body)
		if err == nil {
			var data ProductsResp
			err = json.Unmarshal([]byte(ret), &data)
			if err == nil {
				if len(data.Data.Products) > 0 {
					for _, prod := range data.Data.Products {
						AppCompareAndSave(prod)
						time.Sleep(time.Second * time.Duration(rand.Intn(3)+1))
					}
				} else {
					fmt.Printf("Finish [%s_%s_%s]\n", channel, countryCode, updateScope)
					break
				}
			} else {
				fmt.Printf("RespParseErr [%s_%s_%s] [%s] %s\n", channel, countryCode, updateScope, err.Error(), url)
				break
			}
		} else {
			fmt.Printf("HttpErr [%s_%s_%s] [%s] %s\n", channel, countryCode, updateScope, err.Error(), url)
			break
		}
	}
}

func AppCompareAndSave(prod *Product) {
	appDao := dbs.AppDao{}
	chgDao := dbs.ChangeLogDao{}
	var addSdks []string
	var delSdks []string
	var addServices []string
	var delServices []string
	app, err := appDao.FindByUid(prod.Uid)
	if err != nil { //新的，写入
		devIdStr := ""
		devTitle := ""
		if prod.Publisher != nil {
			devIdStr = prod.Publisher.Uid
			devTitle = prod.Publisher.Title
		}
		channel, platform, bundle := ParseUid(prod.Uid)
		appDao.Create(dbs.AppDao{
			Title:             prod.Title,
			LogoUrl:           GetIconUrl(prod.BundleId, prod.LogoUrl, "apps"), // 更换获取logo的方式
			DeveloperIdStr:    devIdStr,
			DeveloperTitle:    devTitle,
			Channel:           channel,
			Platform:          platform,
			BundleId:          bundle,
			Uid:               prod.Uid,
			Category:          prod.Category,
			DownloadCount:     prod.Latest1yDownloads,
			FirstReleaseDate:  prod.FirstReleaseDate,
			LatestReleaseDate: prod.LatestReleaseDate,
			Website:           prod.Homepage,
			Size:              prod.Size,
			Paid:              ParsePaid(prod.Paid),
			LatestVersion:     prod.Version,
			SdkUids:           strings.Join(prod.SdkUids, ","),
			SdkDevs:           strings.Join(prod.SdkProviders, ","),
			CloudServices:     strings.Join(prod.CloudServiceUids, ","),
			CloudServiceDevs:  strings.Join(prod.CloudServiceProviders, ","),
			ReleaseDate:       time.Now(),
			CreateTime:        time.Now(),
		})
		addSdks = prod.SdkUids
		addServices = prod.CloudServiceUids
		delSdks = []string{}
		delServices = []string{}
		//TODO save prod.BundleId
	} else { //更新
		addSdks, delSdks = compareIds(strings.Split(app.SdkUids, ","), prod.SdkUids)
		addServices, delServices = compareIds(strings.Split(app.CloudServices, ","), prod.CloudServiceUids)
		//更新app
		upd := map[string]interface{}{}
		if app.Title != prod.Title {
			upd["title"] = prod.Title
		}
		if prod.Publisher != nil && app.DeveloperIdStr != prod.Publisher.Uid {
			upd["developer_id_str"] = prod.Publisher.Uid
			upd["developer_title"] = prod.Publisher.Title
		}
		if app.Category != prod.Category {
			upd["category"] = prod.Category
		}
		if prod.Latest1yDownloads > app.DownloadCount {
			upd["download_count"] = prod.Latest1yDownloads
		}
		if app.FirstReleaseDate != prod.FirstReleaseDate {
			upd["first_release_date"] = prod.FirstReleaseDate
		}
		if app.LatestReleaseDate != prod.LatestReleaseDate {
			upd["latest_release_date"] = prod.LatestReleaseDate
		}
		if app.Website != prod.Homepage {
			upd["website"] = prod.Homepage
		}
		if app.Size != prod.Size {
			upd["size"] = prod.Size
		}
		if app.LatestVersion != prod.Version {
			upd["latest_version"] = prod.Version
		}
		if len(addSdks) > 0 || len(delSdks) > 0 {
			upd["sdk_uids"] = strings.Join(prod.SdkUids, ",")
			upd["sdk_devs"] = strings.Join(prod.SdkProviders, ",")
		}
		if len(addServices) > 0 || len(delServices) > 0 {
			upd["cloud_services"] = strings.Join(prod.CloudServiceUids, ",")
			upd["cloud_service_devs"] = strings.Join(prod.CloudServiceProviders, ",")
		}
		appDao.Updates(app.ID, upd)
		if len(upd) > 0 {
			//TODO save prod.BundleId
		}
	}
	//增加change log
	if len(addSdks) > 0 || len(delSdks) > 0 || len(addServices) > 0 || len(delServices) > 0 {
		chgDao.Create(dbs.ChangeLogDao{
			Uid:           prod.Uid,
			ChangeVersion: prod.Version,
			CreateTime:    time.Now(),
			ChangeTime:    prod.LatestReleaseDate,
			AddSdks:       strings.Join(addSdks, ","),
			AddServices:   strings.Join(addServices, ","),
			DelSdks:       strings.Join(delSdks, ","),
			DelServices:   strings.Join(delServices, ","),
		})
	}
	//检查开发者
	if prod.Publisher != nil {
		devIdStr := prod.Publisher.Uid
		CheckDeveloper(devIdStr)
	}
	//检查SDK uids
	if len(addSdks) > 0 {
		CheckSkdUids(addSdks)
	}
	//检查services uids
	if len(addServices) > 0 {
		CheckServiceUids(addServices)
	}
}

func compareIds(oldIds, newIds []string) ([]string, []string) {
	addIds := []string{}
	delIds := []string{}
	oldMap := map[string]bool{}
	for _, id := range oldIds {
		if _, ok := oldMap[id]; !ok {
			oldMap[id] = true
		}
	}
	for _, id := range newIds {
		if _, ok := oldMap[id]; !ok {
			addIds = append(addIds, id)
		}
	}
	newMap := map[string]bool{}
	for _, id := range newIds {
		if _, ok := newMap[id]; !ok {
			newMap[id] = true
		}
	}
	for _, id := range oldIds {
		if _, ok := newMap[id]; !ok {
			delIds = append(delIds, id)
		}
	}
	return addIds, delIds
}

func CheckDeveloper(devId string) {

}

func CheckSkdUids(sdkUids []string) {

}

func CheckServiceUids(serviceUids []string) {

}

func RefreshAppInfos(appInfoIds []string) {

}

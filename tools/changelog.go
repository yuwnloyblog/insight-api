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

func ChangeLogs() {
	var headers map[string]string = map[string]string{}
	headers["cookie"] = `_ga=GA1.1.390781912.1665500669; _ga_MF5DDQ4TF9=GS1.1.1667011701.5.0.1667011701.0.0.0; fork_session=2a4a8f7c9120f145; mp_f39218c184e602888e26ea53b00435dd_mixpanel={"distinct_id": "yuhongda0315@163.com","$device_id": "18422433f8a19e-0c8efb2c2808bf-19525635-13c680-18422433f8b533","$initial_referrer": "$direct","$initial_referring_domain": "$direct","$user_id": "yuhongda0315@163.com","g_team_name": "grtd"}; intercom-session-d1key3b8=c21zdnNsSm9Ub2FiTjJsTi9JVnJZKzlZckNkRkFiNFFxblR2RmlzVWh6QjFva2JPaEpYbGdqcDB4NTZublcrLy0tTkRYc3NCalJEbmZYd2NVbFpoS01YZz09--0b9aef4a431ce97df1df5f2e1081b41afd317947; _ga_G812B88X1Y=GS1.1.1667021202.12.1.1667022410.0.0.0`
	headers["content-type"] = "application/json"
	rand.Seed(time.Now().Unix())
	pageReleaseId := int64(0)
	pageUid := ""
	pageChgTime := ""
	chgDao := dbs.ChangeLogDao{}
	for {
		url := "https://api.app.forkai.cn/webapi/products/changelogs?pageSize=10&startTime=2022-09-29T00:00:00%2B08:00&endTime=2022-10-29T23:59:59%2B08:00&" + fmt.Sprintf("pageMarkCountryCode=CN&pageMarkReleaseID=%d&pageMarkUID=%s&pageMarkChangeTime=%s", pageReleaseId, pageUid, pageChgTime)

		ret, err := utils.HttpDo("GET", url, headers, "")
		if err == nil {
			var chgResp ChangelogsResp
			err = json.Unmarshal([]byte(ret), &chgResp)
			if err == nil && chgResp.Data != nil {
				for _, chg := range chgResp.Data.Changelogs {
					pageReleaseId = chg.ReleaseId
					pageUid = chg.Uid
					pageChgTime = strings.ReplaceAll(chg.ChangeTime, "+", "%2B")
					chgDao.Create(dbs.ChangeLogDao{
						Uid:           chg.Uid,
						ChangeVersion: chg.Version,
						CreateTime:    time.Now(),
						ChangeTime:    chg.ChangeTime,
						AddSdks:       strings.Join(chg.InstalledSdkUids, ","),
						AddServices:   strings.Join(chg.InstalledServicesIds, ","),
						DelSdks:       strings.Join(chg.UninstallSdkUids, ","),
						DelServices:   strings.Join(chg.UninstalledServicesIds, ","),
					})
					fmt.Println("change_logs:", chg.Uid, chg.ChangeTime)
					CheckApp(chg.Uid, headers)
					CheckSkdUids(chg.InstalledSdkUids, headers)
					CheckServiceUids(chg.InstalledServicesIds, headers)

				}
				time.Sleep(1 * time.Second)
				fmt.Println("*********************")
			} else {
				break
			}
		} else {
			break
		}
	}

}

func CheckApp(uid string, headers map[string]string) {
	prod := CatchApp(uid, headers)
	if prod != nil {
		appDao := dbs.AppDao{}
		app, err := appDao.FindByUid(uid)
		if err == nil && app != nil { //更新
			// 更新app
			upd := map[string]interface{}{}

			ReloadPic(prod.LogoUrl, "apps", app.ID)

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
			upd["sdk_uids"] = strings.Join(prod.SdkUids, ",")
			upd["sdk_devs"] = strings.Join(prod.SdkProviders, ",")

			upd["cloud_services"] = strings.Join(prod.CloudServiceUids, ",")
			upd["cloud_service_devs"] = strings.Join(prod.CloudServiceProviders, ",")
			logoUrl, logo_err := ReloadPic(prod.LogoUrl, "apps", app.ID)
			if logo_err == nil {
				upd["logo_url"] = logoUrl
			} else {
				logoUrl = ""
			}
			err = appDao.Updates(app.ID, upd)
			fmt.Println("Update_App. ", app.ID, prod.Uid, prod.Title, err)
			AddorUpdateAppInfo(prod, logoUrl)
		} else { //新的写入
			devIdStr := ""
			devTitle := ""
			if prod.Publisher != nil {
				devIdStr = prod.Publisher.Uid
				devTitle = prod.Publisher.Title
			}
			channel, platform, bundle := ParseUid(prod.Uid)

			nId, err := appDao.Create(dbs.AppDao{
				Title: prod.Title,
				//LogoUrl:           logoUrl,
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
			fmt.Println("Add_App. ", nId, prod.Uid, prod.Title, err)
			logoUrl, url_err := ReloadPic(prod.LogoUrl, "apps", nId)
			if url_err != nil {
				logoUrl = prod.LogoUrl
			}
			if err == nil && nId > 0 {
				appDao.UpdateLogo(nId, logoUrl)
			}
			AddorUpdateAppInfo(prod, logoUrl) //更新appinfo
			CheckDeveloper(devIdStr, headers)
		}
	} else {
		fmt.Println("Not Catch App. ", uid)
	}
}

type ChangelogsResp struct {
	Data *ChangelogsCommon `json:"data"`
}
type ChangelogsCommon struct {
	Changelogs []*Changelog `json:"changelogs"`
}

/*
 */
type Changelog struct {
	ChangeTime             string   `json:"changeTime"`
	Uid                    string   `json:"uid"`
	Name                   string   `json:"name"`
	BundleId               string   `json:"bundleID"`
	Version                string   `json:"version"`
	ReleaseId              int64    `json:"releaseID"`
	Channels               []string `json:"channels"`
	Category               string   `json:"category"`
	LogoUrl                string   `json:"logoURL"`
	Starred                bool     `json:"starred"`
	InstalledSdkUids       []string `json:"installedSDKUIDs"`
	UninstallSdkUids       []string `json:"uninstalledSDKUIDs"`
	InstalledServicesIds   []string `json:"installedCloudServiceUIDs"`
	UninstalledServicesIds []string `json:"uninstalledCloudServiceUIDs"`
	FirstSdkVersion        bool     `json:"firstSDKVersion"`
	FirstServiceVersion    bool     `json:"firstCloudServiceVersion"`
	CountryCode            string   `json:"countryCode"`
}

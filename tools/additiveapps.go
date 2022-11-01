package tools

import (
	"encoding/json"
	"fmt"
	"insight-api/dbs"
	"insight-api/utils"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func AdditiveAppsScanByChannel(channel, countryCode, updateScope string) {
	var headers map[string]string = map[string]string{}
	headers["cookie"] = `_ga=GA1.1.390781912.1665500669; _ga_MF5DDQ4TF9=GS1.1.1667011701.5.0.1667011701.0.0.0; fork_session=2c549e6faccb141b; mp_f39218c184e602888e26ea53b00435dd_mixpanel={"distinct_id": "yuhongda0315@163.com","$device_id": "18421a402baccf-0be1423c25754e-19525635-13c680-18421a402bbd03","$user_id": "yuhongda0315@163.com","$initial_referrer": "$direct","$initial_referring_domain": "$direct","g_team_name": "grtd"}; intercom-session-d1key3b8=dEM3SVBUbDdzYnErVUNuTDRTVVY1MS96Wk9RSkJnZFRzd3dFcVhSVExpZy93TzJ3RlFhUkZkSVY4TEcrL0JCKy0taFZMenFyNm10N1Q1WGN3RnkvQ29aZz09--abcfbd8c9c4d6f456f2c787613beee4bd814061a; _ga_G812B88X1Y=GS1.1.1667011692.10.1.1667014932.0.0.0`
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
						AppCompareAndSave(prod, headers)
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

func AppCompareAndSave(prod *Product, headers map[string]string) {
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
		logoUrl, url_err := ReloadPic(prod.LogoUrl, "apps", nId)
		if url_err != nil {
			logoUrl = prod.LogoUrl
		}
		if err == nil && nId > 0 {
			appDao.UpdateLogo(nId, logoUrl)
		}

		addSdks = prod.SdkUids
		addServices = prod.CloudServiceUids
		delSdks = []string{}
		delServices = []string{}
		//TODO save prod.BundleId
		AddorUpdateAppInfo(prod, logoUrl)
	} else { //更新
		addSdks, delSdks = compareIds(strings.Split(app.SdkUids, ","), prod.SdkUids)
		addServices, delServices = compareIds(strings.Split(app.CloudServices, ","), prod.CloudServiceUids)
		//更新app
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
			AddorUpdateAppInfo(prod, app.LogoUrl)
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
		CheckDeveloper(devIdStr, headers)
	}
	//检查SDK uids
	if len(addSdks) > 0 {
		CheckSkdUids(addSdks, headers)
	}
	//检查services uids
	if len(addServices) > 0 {
		CheckServiceUids(addServices, headers)
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

func CheckDeveloper(devId string, headers map[string]string) {
	devDao := dbs.DeveloperDao{}
	dev, err := devDao.FindById(devId)
	if err == nil && dev != nil {

	} else {
		publisher := CatchDeveloper(devId, headers)
		if publisher != nil {
			devDao.Create(dbs.DeveloperDao{
				ID:             devId,
				Title:          publisher.Title,
				LogoUrl:        GenerateDevLogoUrl(devId, publisher.LogoUrl), //GetIconUrl("", publisher.LogoUrl, "devs"),
				Industry:       publisher.Industry,
				FoundedYear:    strconv.Itoa(publisher.FoundedYear),
				Address:        publisher.Address,
				Website:        publisher.Website,
				Email:          publisher.Email,
				Description:    publisher.Description,
				AddressArea:    strings.Join(publisher.AddressArea, "-"),
				CreateTime:     time.Now(),
				FinancingRound: publisher.FinancingRound,
				DownloadCount:  publisher.DownloadCount,
				AppCount:       publisher.ProductCount,
				WebsiteCount:   publisher.WebsiteCount,
			})
		}
	}
}

func CheckSkdUids(sdkUids []string, headers map[string]string) {
	sdkDao := dbs.SdkDao{}
	for _, sdkuid := range sdkUids {
		dbSdk, err := sdkDao.FindById(sdkuid)
		if err == nil && dbSdk != nil {

		} else {
			sdk := CatchSdk(sdkuid, headers)
			if sdk != nil {
				icon, err := ReloadSdkPic(sdk.LogoUrl, "sdks", sdk.Uid)
				if err != nil {
					icon = sdk.LogoUrl
				}
				sdkDao.Create(dbs.SdkDao{
					ID:            sdk.Uid,
					Title:         sdk.Title,
					Platforms:     strings.Join(sdk.Platforms, ","),
					DeveloperId:   sdk.PublisherUid,
					DeveloperName: sdk.PublisherName,
					Category:      sdk.Category,
					LogoUrl:       icon,
				})
				//TODO sdk provider
				CheckServiceProvider(sdk.PublisherUid, headers)
			}
		}
	}
}

func CheckServiceUids(serviceUids []string, headers map[string]string) {

}

func CheckServiceProvider(devId string, headers map[string]string) {
	providerDao := dbs.ServiceProviderDao{}
	dbProvider, err := providerDao.FindById(devId)
	if err == nil && dbProvider != nil {

	} else {
		provider := CatchServiceProvider(devId, headers)
		if provider != nil {
			providerDao.Create(dbs.ServiceProviderDao{
				ID:             provider.Uid,
				Title:          provider.Title,
				LogoUrl:        GenerateDevLogoUrl(provider.Uid, provider.LogoUrl),
				Industry:       provider.Industry,
				FoundedYear:    strconv.Itoa(provider.FoundedYear),
				AddressArea:    strings.Join(provider.AddressArea, "-"),
				CreateTime:     time.Now(),
				FinancingRound: provider.FinancingRound,
				Address:        provider.Address,
				Website:        provider.Website,
				Description:    provider.Description,
				Email:          provider.Email,
			})
		}
	}
}

func AddorUpdateAppInfo(prod *Product, logoUrl string) {
	bundleId := prod.BundleId
	appinfoDao := dbs.AppInfoDao{}
	appinfo, err := appinfoDao.FindById(bundleId)
	if err == nil && appinfo != nil { //更新
		appDao := dbs.AppDao{}
		apps, err := appDao.QueryListByBundleId(bundleId)
		if err == nil {
			totalCount := 0
			var logurl string = ""
			channels := []string{}
			for _, a := range apps {
				if logurl == "" {
					logurl = a.LogoUrl
				}
				channels = append(channels, a.Channel)
				totalCount += int(a.DownloadCount)
			}
			upd := map[string]interface{}{}
			if logurl != "" {
				upd["logo_url"] = logurl
			}
			if totalCount > 0 {
				upd["download_count"] = totalCount
			}
			if len(channels) > 0 {
				upd["channels"] = strings.Join(channels, ",")
			}

			if len(upd) > 0 {
				appinfoDao.Updates(bundleId, upd)
			}
		}
	} else {
		appinfoDao.Create(dbs.AppInfoDao{
			Id:                prod.BundleId,
			Title:             prod.Title,
			Website:           prod.Homepage,
			CreateTime:        time.Now(),
			LogoUrl:           logoUrl,
			Category:          prod.Category,
			LatestVersion:     prod.Version,
			DeveloperId:       prod.Publisher.Uid,
			DeveloperTitle:    prod.Publisher.Title,
			FirstReleaseDate:  prod.FirstReleaseDate,
			LatestReleaseDate: prod.LatestReleaseDate,
			Channels:          prod.Channel,
		})
	}
}

func GenerateDevLogoUrl(devId, oldUrl string) string {
	newUrl, err := ReloadDevPic(oldUrl, "devs", devId)
	if err != nil {
		return oldUrl
	}
	return newUrl
}

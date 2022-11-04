package tools

import (
	"fmt"
	"insight-api/dbs"
	"strings"
	"time"
)

func XiufuSdkuids() {
	var headers map[string]string = map[string]string{}
	headers["cookie"] = `_ga=GA1.1.390781912.1665500669; _ga_MF5DDQ4TF9=GS1.1.1667011701.5.0.1667011701.0.0.0; fork_session=2a4a8f7c9120f145; mp_f39218c184e602888e26ea53b00435dd_mixpanel={"distinct_id": "yuhongda0315@163.com","$device_id": "18422433f8a19e-0c8efb2c2808bf-19525635-13c680-18422433f8b533","$initial_referrer": "$direct","$initial_referring_domain": "$direct","$user_id": "yuhongda0315@163.com","g_team_name": "grtd"}; intercom-session-d1key3b8=c21zdnNsSm9Ub2FiTjJsTi9JVnJZKzlZckNkRkFiNFFxblR2RmlzVWh6QjFva2JPaEpYbGdqcDB4NTZublcrLy0tTkRYc3NCalJEbmZYd2NVbFpoS01YZz09--0b9aef4a431ce97df1df5f2e1081b41afd317947; _ga_G812B88X1Y=GS1.1.1667021202.12.1.1667022410.0.0.0`
	headers["content-type"] = "application/json"
	appDao := dbs.AppDao{}

	page := 1
	for {
		apps, err := appDao.QueryNoSdks((page-1)*1000, 1000)
		page++
		if err == nil && len(apps) > 0 {
			for _, app := range apps {
				prod := CatchAppBySearch(app.Uid, app.Channel, app.Title, headers)
				if prod != nil && len(prod.SdkUids) > 0 {
					upd := map[string]interface{}{}
					upd["sdk_uids"] = strings.Join(prod.SdkUids, ",")
					upd["sdk_devs"] = strings.Join(prod.SdkProviders, ",")
					upd["cloud_services"] = strings.Join(prod.CloudServiceUids, ",")
					upd["cloud_service_devs"] = strings.Join(prod.CloudServiceProviders, ",")

					appDao.Updates(app.ID, upd)
					fmt.Println("update ", app.ID)
				} else {
					fmt.Println("noupdate ", app.ID)
				}
				time.Sleep(1 * time.Second)
			}
		} else {
			break
		}
	}
}

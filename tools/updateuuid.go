package tools

import (
	"fmt"
	"insight-api/dbs"
	"insight-api/utils"
	"math"
	"strings"
	"time"
)

func UpdateAppInfoDevId(s string) {
	appInfoDao := dbs.AppInfoDao{}
	start := s
	for {
		appInfos, err := appInfoDao.QueryList(start, 100)
		if err == nil && len(appInfos) > 0 {
			for _, app := range appInfos {
				start = app.Id

				newDevId := strings.ReplaceAll(app.DeveloperId, "-", "")
				newDevId, err = utils.PruneUuid(newDevId)
				if err != nil {
					fmt.Println("Err id:", app.Id)
					continue
				}
				upd := map[string]interface{}{}
				upd["developer_id"] = newDevId

				err := appInfoDao.Updates(app.Id, upd)
				if err != nil {
					fmt.Println("Fail ", app.Id, newDevId)
				} else {
					fmt.Println("Success ", app.Id, newDevId)
				}

				time.Sleep(50 * time.Millisecond)
			}
		}
	}
}

func UpdateAppDevId(start int64) {
	appDao := dbs.AppDao{}
	if start <= 0 {
		start = int64(math.MaxInt32)
	}
	for {
		apps, err := appDao.QueryList("", "", int64(start), 100)
		if err == nil && len(apps) > 0 {
			for _, app := range apps {
				start = app.ID
				upd := map[string]interface{}{}
				if app.DeveloperIdStr != "" && len(app.DeveloperIdStr) > 26 {
					time.Sleep(5 * time.Millisecond)

					newDevId := strings.ReplaceAll(app.DeveloperIdStr, "-", "")
					newDevId, err = utils.PruneUuid(newDevId)
					if err != nil {
						fmt.Println("Err id:", app.ID)
						continue
					}
					upd["developer_id_str"] = newDevId
				}

				//sdkDevids
				sdkDevids, err := TransDevIds(app.SdkDevs)
				if err == nil {
					upd["sdk_devs"] = sdkDevids
				}
				//cloudDevids
				cloudDevids, err := TransDevIds(app.CloudServiceDevs)
				if err == nil {
					upd["cloud_service_devs"] = cloudDevids
				}
				if len(upd) > 0 {
					//更新数据库
					appDao.Updates(app.ID, upd)
					fmt.Println("id:", app.ID, upd)
					time.Sleep(5 * time.Millisecond)
				}
				fmt.Println("id:", app.ID)
			}
		}
	}
}

func TransDevIds(devids string) (string, error) {
	if len(devids) <= 0 {
		return "", fmt.Errorf("err")
	}
	devArr := strings.Split(devids, ",")
	newDevArr := make([]string, 0)
	for _, devId := range devArr {
		newDevId := strings.ReplaceAll(devId, "-", "")
		newDevId, err := utils.PruneUuid(newDevId)
		if err != nil {
			return "", err
		}
		newDevArr = append(newDevArr, newDevId)
	}
	return strings.Join(newDevArr, ","), nil
}

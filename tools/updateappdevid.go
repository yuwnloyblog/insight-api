package tools

import (
	"fmt"
	"insight-api/dbs"
	"insight-api/utils"
	"math"
	"strings"
	"time"
)

func UpdateAppDevId(start int64) {
	appDao := dbs.AppDao{}
	devDao := dbs.DeveloperDao{}
	if start <= 0 {
		start = int64(math.MaxInt32)
	}
	for {
		apps, err := appDao.QueryList("", "", int64(start), 100)
		if err == nil && len(apps) > 0 {
			for _, app := range apps {
				start = app.ID
				if app.DeveloperIdStr != "" && len(app.DeveloperIdStr) > 26 {
					time.Sleep(50 * time.Millisecond)

					newIdStr := strings.ReplaceAll(app.DeveloperIdStr, "-", "")
					newIdStr, err = utils.PruneUuid(newIdStr)
					if err != nil {
						fmt.Println("Err id:", app.ID)
						continue
					}

					//更新数据库
					err = appDao.UpdateDevIdStr(app.ID, newIdStr)
					fmt.Println(err)
					err = devDao.UpdateId(app.DeveloperIdStr, newIdStr)
					fmt.Println(err)
					fmt.Println("id:", app.ID, "old_dev_id:", app.DeveloperIdStr, "new_dev_id:", newIdStr)
				}
			}
		}
	}
}

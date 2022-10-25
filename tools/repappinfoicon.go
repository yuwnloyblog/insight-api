package tools

import (
	"fmt"
	"insight-api/dbs"
	"time"
)

func ReplaceIcon4AppInfo(startPage int) {
	appDao := dbs.AppDao{}

	pageIndex := startPage
	for {
		apps, err := appDao.QueryListByPage("", "", pageIndex, 100)
		if err == nil && len(apps) > 0 {
			for _, app := range apps {
				packageName := app.BundleId
				if packageName != "" {
					iconUrl, err := GetIconFromMyApp(packageName)
					if err == nil {
						appDao.UpdateLogo(app.ID, iconUrl)
						fmt.Println(pageIndex, app.ID, app.BundleId, iconUrl)
					} else {
						continue
					}
					time.Sleep(50 * time.Millisecond)
				}
			}
			pageIndex++
		} else {
			break
		}
	}
}

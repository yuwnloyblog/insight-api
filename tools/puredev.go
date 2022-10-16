package tools

import (
	"fmt"
	"insight-api/dbs"
	"math"
	"time"
)

func PureDev() {
	devdao := dbs.DeveloperDao{}
	start := ""
	for {
		devs, err := devdao.QueryList("", start, 100)
		if err == nil && len(devs) > 0 {
			for _, dev := range devs {
				start = dev.ID
				isHas := HasApps(dev.ID)
				if !isHas {
					fmt.Println("dev_id:", dev.ID)
					devdao.Delete(dev.ID)
					time.Sleep(10 * time.Millisecond)
				}
			}
		} else {
			break
		}
	}
}

func HasApps(devId string) bool {
	appdao := dbs.AppDao{}
	apps, err := appdao.QueryList("", devId, math.MaxInt64, 1)
	if err == nil && len(apps) > 0 {
		return true
	}
	return false
}

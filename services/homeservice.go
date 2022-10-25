package services

import (
	"encoding/json"
	"insight-api/dbs"
	"insight-api/utils"
)

func GetHomeInfo() *HomeInfo {
	cacheKey := "homeinfo"
	if cacheObj, ok := utils.CacheGet(cacheKey); ok {
		return cacheObj.(*HomeInfo)
	}
	confDao := dbs.ConfDao{}
	var retHomeinfo HomeInfo
	conf, err := confDao.FindConfById("homeinfo")
	if err == nil && conf.Value != "" {
		err = json.Unmarshal([]byte(conf.Value), &retHomeinfo)
		if err == nil {
			utils.CachePut(cacheKey, &retHomeinfo)
			return &retHomeinfo
		}
	}
	retHomeinfo = HomeInfo{
		TotalAppCount:       808038,
		TotalDeveloperCount: 923604,
		Sotre: map[string]int{
			"apple":       379842,
			"huawei":      161159,
			"xiaomi":      66858,
			"a360":        66281,
			"google_play": 133898,
		},
		Platform: map[string]int{
			"android": 428196,
			"ios":     379842,
		},
	}
	utils.CachePut(cacheKey, &retHomeinfo)
	return &retHomeinfo
}

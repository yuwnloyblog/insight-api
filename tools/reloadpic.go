package tools

import (
	"fmt"
	"insight-api/dbs"
	"insight-api/utils"
	"math"
	"strings"
	"time"
)

func ReloadAppPic(start int64) {
	appDao := dbs.AppDao{}
	if start <= 0 {
		start = int64(math.MaxInt32)
	}
	for {
		apps, err := appDao.QueryList("", "", int64(start), 100)
		if err == nil && len(apps) > 0 {
			for _, app := range apps {
				start = app.ID
				if app.LogoUrl != "" && !strings.HasPrefix(app.LogoUrl, "https://file.lwoowl.cn") {
					time.Sleep(50 * time.Millisecond)
					nf, err := ReloadPic(app.LogoUrl, "apps", app.ID)
					if err != nil {
						fmt.Println(err.Error(), app.ID)
						continue
					}
					nf = "https://file.lwoowl.cn/" + nf
					//更新数据库
					// appDao.UpdateLogo(app.ID, nf)
					fmt.Println("id:", app.ID, "old_url:", app.LogoUrl, "new_url:", nf)
				}
			}
		}
	}
}
func ReloadPic(url, middlePath string, id int64) (string, error) {
	filename, err := GetFileNameFromUrl(url)
	if err != nil {
		return "", fmt.Errorf("Err_Url %s %s", url, err.Error())
	}
	err = DownloadPicture(url, filename)
	if err != nil {
		return "", fmt.Errorf("Err_Download %s %s", url, err.Error())

	}
	tail := GetFileTail(filename)

	if err != nil {
		return "", fmt.Errorf("Err_HandleUuid %s %s", url, err.Error())
	}
	fn, _ := utils.Encode(id)
	newFileName := fmt.Sprintf("%s/%s", middlePath, fn)

	if tail != "" {
		newFileName = newFileName + "." + tail
	}
	err = QiniuUpload(filename, newFileName)
	DeleteFile(filename)
	if err != nil {
		return "", fmt.Errorf("Err_Upload %s %s", url, err.Error())
	}
	return newFileName, nil
}

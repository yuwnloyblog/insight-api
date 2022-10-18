package tools

import (
	"fmt"
	"insight-api/dbs"
	"math"
	"time"
)

func ReloadAppPic(t string) {
	appDao := dbs.AppDao{}
	start := int64(math.MaxInt32)
	for {
		apps, err := appDao.QueryList("", "", int64(start), 100)
		if err == nil && len(apps) > 0 {
			for _, app := range apps {
				start = app.ID
				if app.LogoUrl != "" {
					time.Sleep(50 * time.Millisecond)
					err = ReloadPic(app.LogoUrl, "apps")
					if err != nil {
						fmt.Println(err.Error(), app.ID)
						continue
					}
				}
			}
		}
	}
}
func ReloadPic(url, middlePath string) error {
	filename, err := GetFileNameFromUrl(url)
	if err != nil {
		return fmt.Errorf("Err_Url %s %s", url, err.Error())
	}
	err = DownloadPicture(url, filename)
	if err != nil {
		return fmt.Errorf("Err_Download %s %s", url, err.Error())

	}
	newFileName := fmt.Sprintf("/%s/%s", middlePath, GetClearUuid())
	tail := GetFileTail(filename)
	if tail != "" {
		newFileName = newFileName + "." + tail
	}
	err = QiniuUpload(filename, newFileName)
	if err != nil {
		return fmt.Errorf("Err_Upload %s %s", url, err.Error())
	}
	DeleteFile(filename)
	return nil
}

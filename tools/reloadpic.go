package tools

import (
	"fmt"
	"insight-api/dbs"
	"insight-api/utils"
	"strings"
	"time"
)

func UpdateAppLogo() {
	appDao := dbs.AppDao{}
	appInfoDao := dbs.AppInfoDao{}
	start := 0
	for {
		apps, err := appDao.QueryList("", "", int64(start), 1000)
		if err == nil && len(apps) > 0 {
			for _, app := range apps {
				start = int(app.ID)
				if app.ID < 354417 {
					if app.LogoUrl != "" && !strings.HasPrefix(app.LogoUrl, "https://file.lwoowl.cn") && !strings.HasPrefix(app.LogoUrl, "https://pp.myapp.com") {
						iconUrl, err := ReloadPicNoUpload(app.LogoUrl, "apps", app.ID)
						if err == nil {
							iconUrl = "https://file.lwoowl.cn/" + iconUrl
							appDao.UpdateLogo(app.ID, iconUrl)
							fmt.Println(app.ID, app.LogoUrl, iconUrl)
							appinfo, err := appInfoDao.FindById(app.BundleId)
							if err == nil {
								if appinfo != nil && !strings.HasPrefix(appinfo.LogoUrl, "https://file.lwoowl.cn") && !strings.HasPrefix(appinfo.LogoUrl, "https://pp.myapp.com") {
									appinfo.UpdateLogo(appinfo.Id, iconUrl)
								}
							}
						}
					}
					time.Sleep(5 * time.Millisecond)
				} else {
					return
				}
			}
		}
	}
}

func ReloadAppPic(start int64) {
	appDao := dbs.AppDao{}
	appInfoDao := dbs.AppInfoDao{}
	for {
		apps, err := appDao.QueryList("", "", int64(start), 100)
		if err == nil && len(apps) > 0 {
			for _, app := range apps {
				start = app.ID
				time.Sleep(10 * time.Millisecond)

				if app.LogoUrl != "" && !strings.HasPrefix(app.LogoUrl, "https://file.lwoowl.cn") && !strings.HasPrefix(app.LogoUrl, "https://pp.myapp.com") {

					iconUrl, err := GetIconFromMyApp(app.BundleId)
					if err != nil {
						nf, err := ReloadPic(app.LogoUrl, "apps", app.ID)
						if err != nil {
							fmt.Println(err.Error(), app.ID)
							continue
						}
						iconUrl = "https://file.lwoowl.cn/" + nf
					}
					if iconUrl != "" {
						appDao.UpdateLogo(app.ID, iconUrl)
						appInfoDao.UpdateLogo(app.BundleId, iconUrl)
						fmt.Println(app.ID, app.LogoUrl, iconUrl)
					}
				}
			}
		}
	}
}

func ReloadSdkPic(url, middlePath, sdkId string) (string, error) {
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
	fn := strings.ReplaceAll(sdkId, ".", "")
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

func ReloadDevPic(url, middlePath, devId string) (string, error) {
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
	fn, _ := utils.PruneUuid(devId)
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
	newFileName := fmt.Sprintf("https://file.lwoowl.cn/%s/%s", middlePath, fn)

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
func ReloadPicNoUpload(url, middlePath string, id int64) (string, error) {
	filename, err := GetFileNameFromUrl(url)
	if err != nil {
		return "", fmt.Errorf("Err_Url %s %s", url, err.Error())
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
	return newFileName, nil
}

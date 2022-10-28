package tools

import (
	"fmt"
	"insight-api/dbs"
	"insight-api/utils"
	"strings"
	"time"
)

func ReloadSdkLogo(start string) {
	sdkDao := dbs.SdkDao{}
	StartIndex := start
	for {
		sdks, err := sdkDao.QueryAllList(StartIndex, 1000)
		if err == nil && len(sdks) > 0 {
			for _, sdk := range sdks {
				StartIndex = sdk.ID
				if sdk.LogoUrl != "" && !strings.HasPrefix(sdk.LogoUrl, "https://file.lwoowl.cn") && !strings.HasPrefix(sdk.LogoUrl, "https://pp.myapp.com") {
					nf, err := ReloadSdkPic(sdk.LogoUrl, "sdks")
					if err != nil {
						fmt.Println(err.Error(), sdk.ID)
						continue
					}
					nf = "https://file.lwoowl.cn/" + nf
					//更新数据库
					sdk.UpdateLogo(sdk.ID, nf)
					fmt.Println("id:", sdk.ID, "old_url:", sdk.LogoUrl, "new_url:", nf)
					time.Sleep(5 * time.Millisecond)
				}
			}
		} else {
			fmt.Println("Quit", err)
			break
		}
	}
}

func ReloadSdkPic(url, middlePath string) (string, error) {
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
	fn, _ := utils.PruneUuid(utils.GetUuid())
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

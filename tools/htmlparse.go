package tools

import (
	"fmt"
	"insight-api/utils"

	"github.com/PuerkitoBio/goquery"
)

func GetIconFromMyApp(packageName string) (string, error) {
	doc, err := goquery.NewDocument(fmt.Sprintf("https://sj.qq.com/appdetail/%s", packageName))
	if err != nil {
		return "", err
	}
	se := doc.Find("img.jsx-2463046046")
	val, ex := se.Attr("src")
	if ex {
		return val, nil
	}
	return "", fmt.Errorf("Can not get Icon.[%s]", packageName)
}

func UploadPic(url, middlePath string) (string, error) {
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
	fn, _ := utils.PruneUuid(utils.GetClearUuid())
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

func GetIconUrl(packageName, logoUrl, middlePath string) string {
	url, err := GetIconFromMyApp(packageName)
	if err == nil && url != "" {
		return url
	}
	url, err = UploadPic(logoUrl, middlePath)
	if err == nil && url != "" {
		return url
	}
	return logoUrl
}

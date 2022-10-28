package tools

import (
	"bufio"
	"fmt"
	"insight-api/dbs"
	"insight-api/utils"
	"io"
	"os"
	"strings"
)

func TmpUpdate() {
	sdkDao := dbs.SdkDao{}
	lines, err := ReadLine("sdk")
	if err == nil {
		for _, line := range lines {
			arr := strings.Split(line, " ")
			id := arr[1]
			oldUrl := arr[3]
			nf, err := ReloadSdkPicTmp(oldUrl, "sdks")
			if err != nil {
				fmt.Println(err.Error(), id)
				continue
			}
			nf = "https://file.lwoowl.cn/" + nf
			//更新数据库
			sdkDao.UpdateLogo(id, nf)
			fmt.Println("id:", id, "old_url:", oldUrl, "new_url:", nf)
		}
	}
}

func ReadLine(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	var result []string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				return result, nil
			}
			return nil, err
		}
		result = append(result, line)
	}
}

func ReloadSdkPicTmp(url, middlePath string) (string, error) {
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

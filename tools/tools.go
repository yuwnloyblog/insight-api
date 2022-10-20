package tools

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

var (
	accessKey = "ey1SqHmKV9HSI0-2IHtV8KysxfKtJBNzRQYKvAUP"
	secretKey = "FZomdOtDJpS368hWGLvJAK-q6QwTLUd9_l8A2oKc"
	bucket    = "insight-file"
)

func QiniuFetch(url string) (string, error) {
	mac := auth.New(accessKey, secretKey)

	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(mac, &cfg)

	resURL := url
	// 指定保存的key
	fetchRet, err := bucketManager.Fetch(resURL, bucket, "qiniu-x.png")
	if err != nil {
		return "", err
	} else {
		return fetchRet.String(), nil
	}
}

func QiniuUpload(needUploadFile, uploadFilePath string) error {
	localFile := needUploadFile
	key := uploadFilePath

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := auth.New(accessKey, secretKey)

	cfg := storage.Config{}
	// 空间对应的机房
	region, err := storage.GetRegion(accessKey, bucket)
	if err != nil {
		region = &storage.ZoneHuabei
	}
	cfg.Zone = region //&storage.ZoneHuabei
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	//构建代理client对象
	client := http.Client{}

	resumeUploader := storage.NewResumeUploaderEx(&cfg, &storage.Client{Client: &client})
	upToken := putPolicy.UploadToken(mac)

	ret := storage.PutRet{}

	err = resumeUploader.PutFile(context.Background(), &ret, upToken, key, localFile, nil)
	if err != nil {
		return err
	}
	return nil
	//fmt.Println(ret.Key, ret.Hash)
}

// https://huawei-icon.oss-cn-hangzhou.aliyuncs.com/application/icon144/70cef5dad83849f8ba75ef7031f12c37.png
// https://company-icon.oss-cn-hangzhou.aliyuncs.com/1d1ba71b33d61b416c21600c7d1984c5.jpg
// 保存图片
func DownloadPicture(url, filename string) error {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, body, 0755)
	if err != nil {
		return err
	}
	return nil
}

func GetFileNameFromUrl(url string) (string, error) {
	arr := strings.Split(url, "/")
	if len(arr) <= 0 {
		return "", errors.New("url is not validate")
	}
	fileName := arr[len(arr)-1]
	fmt.Println(arr[len(arr)-1])
	return fileName, nil
}

func DeleteFile(file string) error {
	return os.Remove(file)
}

func GetFileTail(file string) string {
	arr := strings.Split(file, ".")
	len := len(arr)
	if len > 0 {
		return arr[len-1]
	}
	return ""
}

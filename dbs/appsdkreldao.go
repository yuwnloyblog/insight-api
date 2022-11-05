package dbs

import (
	"bytes"
	"fmt"
)

type AppSdkRelDao struct {
	ID               int    `gorm:"primary_key"`
	AppId            int64  `gorm:"app_id"`
	SdkId            string `gorm:"sdk_id"`
	AppBundleId      string `gorm:"app_bundle_id"`
	AppDownloadCount int64  `gorm:"app_download_count"`
}

func (sdk AppSdkRelDao) TableName() string {
	return "app_sdk_rel"
}

func (app AppSdkRelDao) Create(item AppSdkRelDao) error {
	err := db.Create(&item).Error
	return err
}

func (app AppSdkRelDao) BatchCreate(items []AppSdkRelDao) error {
	var buffer bytes.Buffer
	sql := "insert into `" + app.TableName() + "` (`app_id`,`sdk_id`) values "
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, item := range items {
		if i == len(items)-1 {
			buffer.WriteString(fmt.Sprintf("(%d,'%s');", item.AppId, item.SdkId))
		} else {
			buffer.WriteString(fmt.Sprintf("(%d,'%s'),", item.AppId, item.SdkId))
		}
	}
	return db.Exec(buffer.String()).Error
}

func (app AppSdkRelDao) UpdateByAppId(appId int64, bundleId string) error {
	upd := map[string]interface{}{}
	upd["app_bundle_id"] = bundleId
	return db.Model(&app).Where("app_id=?", appId).Update(upd).Error
}

package dbs

import (
	"gorm.io/gorm"
)

type SdkDao struct {
	ID          int64  `gorm:"primary_key"`
	Name        string `gorm:"name"`
	BundleId    string `gorm:"bundle_id"`
	Description string `gorm:"description"`
	LogoUrl     string `gorm:"logo_url"`
	Classify    string `gorm:"classify"`
	DeveloperId int64  `gorm:"developer_id"`
}

func (app SdkDao) TableName() string {
	return "sdks"
}

func (app SdkDao) FindById(id int64) (*SdkDao, error) {
	var appItem SdkDao
	err := db.Where("id=?", id).Take(&appItem).Error
	if err != nil {
		return nil, err
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &appItem, nil
}

func (app SdkDao) QueryList(appId int64) ([]*SdkDao, error) {
	var items []*SdkDao
	err := db.Raw("SELECT * FROM sdks LEFT JOIN app_sdks on app_sdks.sdk_id = sdks.id where app_sdks.app_id=?", appId).Scan(&items).Error
	return items, err
}

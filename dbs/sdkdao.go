package dbs

import (
	"gorm.io/gorm"
)

type SdkDao struct {
	ID                    string `gorm:"id"`
	Name                  string `gorm:"name"`
	Platforms             string `gorm:"platforms"`
	Category              string `gorm:"category"`
	DeveloperName         string `gorm:"dev_name"`
	DeveloperId           string `gorm:"dev_id"`
	LogoUrl               string `gorm:"logo_url"`
	InstalledAppCount     int    `gorm:"app_count"`
	DeveloperCount        int    `gorm:"dev_count"`
	InstalledWebsiteCount int    `gorm:"web_count"`
}

func (sdk SdkDao) TableName() string {
	return "sdks"
}

func (sdk SdkDao) FindById(id int64) (*SdkDao, error) {
	var appItem SdkDao
	err := db.Where("id=?", id).Take(&appItem).Error
	if err != nil {
		return nil, err
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &appItem, nil
}

func (sdk SdkDao) QueryList(appId int64) ([]*SdkDao, error) {
	var items []*SdkDao
	err := db.Raw("SELECT * FROM sdks LEFT JOIN app_sdks on app_sdks.sdk_id = sdks.id where app_sdks.app_id=?", appId).Scan(&items).Error
	return items, err
}

func (sdk SdkDao) Create(s SdkDao) error {
	return db.Create(&s).Error
}

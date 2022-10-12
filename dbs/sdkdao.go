package dbs

import (
	"strings"

	"gorm.io/gorm"
)

type SdkDao struct {
	ID             string `gorm:"id"`
	Name           string `gorm:"name"`
	Platforms      string `gorm:"platforms"`
	Category       string `gorm:"category"`
	DeveloperName  string `gorm:"developer_name"`
	DeveloperId    string `gorm:"developer_id"`
	LogoUrl        string `gorm:"logo_url"`
	AppCount       int    `gorm:"app_count"`
	DeveloperCount int    `gorm:"developer_count"`
	WebsiteCount   int    `gorm:"website_count"`
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
	appDao := AppDao{}
	app, err := appDao.FindById(appId)
	if err != nil {
		return nil, err
	}
	ids := strings.Split(app.SdkUids, ",")
	if len(ids) > 0 {
		var items []*SdkDao
		err = db.Raw("SELECT * FROM sdks where id in (?)", ids).Scan(&items).Error
		return items, err
	}
	return nil, nil
	// err := db.Raw("SELECT * FROM sdks LEFT JOIN app_sdks on app_sdks.sdk_id = sdks.id where app_sdks.app_id=?", appId).Scan(&items).Error
	// return items, err
}

func (sdk SdkDao) Create(s SdkDao) error {
	return db.Create(&s).Error
}

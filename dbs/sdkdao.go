package dbs

import (
	"strings"
)

type SdkDao struct {
	ID             string `gorm:"id"`
	Title          string `gorm:"title"`
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

func (sdk SdkDao) FindById(id string) (*SdkDao, error) {
	var appItem SdkDao
	err := db.Where("id=?", id).Take(&appItem).Error
	if err != nil {
		return nil, err
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
		err = db.Debug().Raw("SELECT * FROM sdks where id in (?)", ids).Scan(&items).Error
		return items, err
	}
	return nil, nil
	// err := db.Raw("SELECT * FROM sdks LEFT JOIN app_sdks on app_sdks.sdk_id = sdks.id where app_sdks.app_id=?", appId).Scan(&items).Error
	// return items, err
}

func (sdk SdkDao) QueryAllList(start string, count int) ([]*SdkDao, error) {
	var items []*SdkDao
	if start != "" {
		err := db.Where("id>?", start).Order("id asc").Limit(count).Find(&items).Error
		return items, err
	} else {
		err := db.Order("id asc").Limit(count).Find(&items).Error
		return items, err
	}
}
func (sdk SdkDao) QueryListByPage(keyword string, page, count int) ([]*SdkDao, error) {
	var items []*SdkDao
	whereStr := ""
	args := []interface{}{}
	if keyword != "" {
		if whereStr != "" {
			whereStr = whereStr + " AND "
		}
		whereStr = whereStr + " title like ? "
		args = append(args, "%"+keyword+"%")
	}
	err := db.Debug().Where(whereStr, args).Order("app_count desc").Limit(count).Offset((page - 1) * count).Find(&items).Error
	return items, err
}

func (sdk SdkDao) Create(s SdkDao) error {
	return db.Create(&s).Error
}

func (sdk SdkDao) UpdateLogo(id string, logo string) error {
	upd := map[string]interface{}{}
	upd["logo_url"] = logo
	return db.Model(&sdk).Where("id=?", id).Update(upd).Error
}

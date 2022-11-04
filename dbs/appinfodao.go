package dbs

import (
	"time"
)

type AppInfoDao struct {
	Id                string    `gorm:"id"`
	Title             string    `grom:"title"`
	Website           string    `gorm:"website"`
	Description       string    `gorm:"description"`
	CreateTime        time.Time `gorm:"time"`
	LogoUrl           string    `gorm:"logo_url"`
	Category          string    `gorm:"category"`
	LatestVersion     string    `gorm:"latest_version"`
	DownloadCount     int64     `gorm:"download_count"`
	DeveloperId       string    `gorm:"developer_id"`
	DeveloperTitle    string    `gorm:"developer_title"`
	FirstReleaseDate  string    `grom:"first_release_date"`
	LatestReleaseDate string    `gorm:"latest_release_date"`
	Channels          string    `gorm:"channels"`
}

func (app AppInfoDao) TableName() string {
	return "appinfos"
}

func (app AppInfoDao) FindById(id string) (*AppInfoDao, error) {
	var appItem AppInfoDao
	err := db.Where("id=?", id).Take(&appItem).Error
	if err != nil {
		return nil, err
	}
	return &appItem, nil
}
func (app AppInfoDao) QueryList(start string, count int) ([]*AppInfoDao, error) {
	var items []*AppInfoDao
	err := db.Where("id>?", start).Order("id asc").Find(&items).Error
	return items, err
}
func (app AppInfoDao) QueryListByPage(keyword, devId, sdkId, notSdkId string, page, count int) ([]*AppInfoDao, error) {
	var items []*AppInfoDao
	whereStr := ""
	args := []interface{}{}
	if keyword != "" {
		if whereStr != "" {
			whereStr = whereStr + " AND "
		}
		whereStr = whereStr + " title like ? "
		args = append(args, "%"+keyword+"%")
	}
	if devId != "" {
		if whereStr != "" {
			whereStr = whereStr + " AND "
		}
		whereStr = whereStr + " developer_id = ? "
		args = append(args, devId)
	}

	if sdkId != "" || notSdkId != "" {
		if whereStr != "" {
			whereStr = whereStr + " AND "
		}
		if sdkId != "" {
			whereStr = whereStr + " sdk_id=? "
		} else {
			whereStr = whereStr + " sdk_id!=? "
		}
		args = append(args, sdkId)

		sql := "SELECT * FROM appinfos LEFT JOIN app_sdk_rel on app_sdk_rel.app_bundle_id=appinfos.id"
		if whereStr != "" {
			sql = sql + " WHERE " + whereStr
		}
		err := db.Raw(sql, args...).Order("download_count desc").Limit(count).Offset((page - 1) * count).Find(&items).Error
		return items, err
	}

	err := db.Where(whereStr, args...).Order("download_count desc").Limit(count).Offset((page - 1) * count).Find(&items).Error
	return items, err
}

func (app AppInfoDao) Create(item AppInfoDao) error {
	err := db.Create(&item).Error
	return err
}

func (app AppInfoDao) Updates(id string, upd map[string]interface{}) error {
	if len(upd) > 0 {
		return db.Model(&app).Where("id=?", id).Update(upd).Error
	}
	return nil
}

func (app AppInfoDao) UpdateLogo(id string, url string) error {
	upd := map[string]interface{}{}
	upd["logo_url"] = url
	return db.Model(&app).Where("id=?", id).Update(upd).Error
}

func (app AppInfoDao) UpdateDevIdStr(id string, developerId string) error {
	upd := map[string]interface{}{}
	upd["developer_id"] = developerId
	return db.Model(&app).Where("id=?", id).Update(upd).Error
}

package dbs

import (
	"time"
)

type AppDao struct {
	ID                int64     `gorm:"primary_key"`
	Title             string    `gorm:"title"`
	BundleId          string    `gorm:"bundle_id"` //包名
	Platform          string    `gorm:"platform"`  //平台，iOS/Android
	Channel           string    `gorm:"channel"`   //苹果，华为，小米等
	DeveloperIdStr    string    `gorm:"developer_id_str"`
	DeveloperTitle    string    `gorm:"developer_title"`
	Uid               string    `gorm:"uid"`
	Category          string    `gorm:"category"`
	DownloadCount     int64     `gorm:"download_count"`
	FirstReleaseDate  string    `gorm:"first_release_date"`
	LatestReleaseDate string    `gorm:"latest_release_date"`
	Size              int64     `gorm:"size"`
	Paid              string    `gorm:"paid"`
	LatestVersion     string    `gorm:"latest_version"`
	CountryCode       string    `gorm:"country_code"`
	SdkUids           string    `gorm:"sdk_uids"`
	SdkDevs           string    `gorm:"sdk_devs"`
	CloudServices     string    `gorm:"cloud_services"`
	CloudServiceDevs  string    `gorm:"cloud_service_devs"`
	Website           string    `gorm:"website"`
	Description       string    `gorm:"description"`
	ReleaseDate       time.Time `gorm:"release_date"`
	DeveloperId       int64     `gorm:"developer_id"`
	LogoUrl           string    `gorm:"logo_url"`
	CreateTime        time.Time `gorm:"create_time"`
}

func (app AppDao) TableName() string {
	return "apps"
}

func (app AppDao) FindById(id int64) (*AppDao, error) {
	var appItem AppDao
	err := db.Where("id=?", id).Take(&appItem).Error
	if err != nil {
		return nil, err
	}
	return &appItem, nil
}

func (app AppDao) QueryList(keyword, devId string, start int64, count int) ([]*AppDao, error) {
	var items []*AppDao

	whereStr := "id < ?"
	args := []interface{}{}
	args = append(args, start)
	if keyword != "" {
		whereStr = whereStr + " AND title like ?"
		args = append(args, "%"+keyword+"%")
	}
	if devId != "" {
		whereStr = whereStr + " AND developer_id_str = ?"
		args = append(args, devId)
	}
	err := db.Where(whereStr, args...).Order("id desc").Limit(count).Find(&items).Error
	return items, err

}

func (app AppDao) QueryListByPage(keyword, devId string, page, count int) ([]*AppDao, error) {
	var items []*AppDao
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
		whereStr = whereStr + " developer_id_str = ? "
		args = append(args, devId)
	}
	err := db.Where(whereStr, args...).Order("download_count desc").Limit(count).Offset((page - 1) * count).Find(&items).Error
	return items, err
}

func (app AppDao) Create(item AppDao) error {
	err := db.Create(&item).Error
	return err
}

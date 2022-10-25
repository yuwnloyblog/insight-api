package dbs

import "time"

type ChangeLogDao struct {
	ID            int64     `gorm:"primary_key"`
	Uid           string    `gorm:"uid"`
	ChangeVersion string    `gorm:"change_version"`
	CreateTime    time.Time `gorm:"change_time"`
	ChangeTime    string    `gorm:"change_time"`
	AddSdks       string    `gorm:"add_sdks"`
	DelSdks       string    `gorm:"del_sdks"`
	AddServices   string    `gorm:"add_services"`
	DelServices   string    `gorm:"del_services"`
	CountryCode   string    `gorm:"country_code"`
}

func (chg ChangeLogDao) TableName() string {
	return "changelogs"
}

func (chg ChangeLogDao) Create(item ChangeLogDao) error {
	err := db.Create(&item).Error
	return err
}

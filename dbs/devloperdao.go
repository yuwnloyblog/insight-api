package dbs

import (
	"time"
)

type DeveloperDao struct {
	ID             int64     `gorm:"primary_key"`
	Title          string    `gorm:"title"`
	LogoUrl        string    `gorm:"logo_url"`
	Description    string    `gorm:"description"`
	Trade          string    `gorm:"trade"`
	Address        string    `gorm:"address"`
	AddressArea    string    `gorm:"address_area"`
	Website        string    `gorm:"website"`
	Contact        string    `gorm:"contact"`
	FoundedTime    time.Time `gorm:"founded_time"`
	CreateTime     time.Time `gorm:"create_time"`
	FinancingRound string    `gorm:"financing_round"`
}

func (dev DeveloperDao) TableName() string {
	return "developers"
}

func (dev DeveloperDao) FindById(id int64) (*DeveloperDao, error) {
	var appItem DeveloperDao
	err := db.Where("id=?", id).Take(&appItem).Error
	if err != nil {
		return nil, err
	}
	//  else if err == gorm.ErrRecordNotFound {
	// 	return nil, nil
	// }
	return &appItem, nil
}

func (dev DeveloperDao) QueryList(keyword string, start int64, count int) ([]*DeveloperDao, error) {
	var items []*DeveloperDao
	err := db.Where("id < ?", start).Order("id desc").Limit(count).Find(&items).Error
	return items, err
}

func (dev DeveloperDao) Create(d DeveloperDao) error {
	return db.Create(&d).Error
}

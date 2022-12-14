package dbs

import (
	"time"
)

type DeveloperDao struct {
	ID             string    `gorm:"id"`
	Title          string    `gorm:"title"`
	LogoUrl        string    `gorm:"logo_url"`
	Industry       string    `gorm:"industry"`
	FoundedYear    string    `gorm:"founded_year"`
	Address        string    `gorm:"address"`
	Website        string    `grom:"website"`
	Email          string    `gorm:"email"`
	Description    string    `gorm:"description"`
	AddressArea    string    `gorm:"address_area"`
	CreateTime     time.Time `gorm:"create_time"`
	FinancingRound string    `gorm:"financing_round"`
	AppCount       int       `gorm:"app_count"`
	WebsiteCount   int       `gorm:"website_count"`
	DownloadCount  int64     `gorm:"download_count"`
}

func (dev DeveloperDao) TableName() string {
	return "developers"
}

func (dev DeveloperDao) FindById(id string) (*DeveloperDao, error) {
	var appItem DeveloperDao
	err := db.Where("id=?", id).Take(&appItem).Error
	if err != nil {
		return nil, err
	}
	return &appItem, nil
}

func (dev DeveloperDao) QueryList(keyword string, start string, count int) ([]*DeveloperDao, error) {
	var items []*DeveloperDao
	if start == "" {
		start = "ffffffff-ffff-ffff-ffff-ffffffffffff"
	}
	if keyword != "" {
		err := db.Where("title like ? AND id < ?", "%"+keyword+"%", start).Order("id desc").Limit(count).Find(&items).Error
		return items, err
	} else {
		err := db.Where("id < ?", start).Order("id desc").Limit(count).Find(&items).Error
		return items, err
	}
}

func (dev DeveloperDao) QueryListByPage(keyword string, page, count int) ([]*DeveloperDao, error) {
	var items []*DeveloperDao
	if keyword != "" {
		err := db.Where("title like ? ", "%"+keyword+"%").Order("download_count desc").Limit(count).Offset((page - 1) * count).Find(&items).Error
		return items, err
	} else {
		err := db.Order("download_count desc").Limit(count).Offset((page - 1) * count).Find(&items).Error
		return items, err
	}
}

func (dev DeveloperDao) Create(d DeveloperDao) error {
	return db.Create(&d).Error
}

func (dev DeveloperDao) Delete(id string) error {
	return db.Where("id=?", id).Delete(&DeveloperDao{}).Error
}

func (dev DeveloperDao) UpdateId(id, newId string) error {
	// upd := map[string]interface{}{}
	// upd["id"] = "nihao"
	// return db.Debug().Model(&dev).Where("id=?", id).Update(upd).Error
	return db.Debug().Exec("update developers set id=? where id=?", newId, id).Error
}

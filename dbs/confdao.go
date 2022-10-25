package dbs

import "time"

type ConfDao struct {
	ID         string    `gorm:"id"`
	Value      string    `gorm:"value"`
	CreateTime time.Time `gorm:"create_time"`
}

func (fd ConfDao) TableName() string {
	return "configures"
}

func (conf ConfDao) FindConfById(id string) (*ConfDao, error) {
	var item ConfDao
	err := db.Where("id=?", id).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

package dbs

import "time"

type FeedbackDao struct {
	ID         int64     `gorm:"primary_key"`
	UserId     int64     `gorm:"user_id"`
	Feedback   string    `gorm:"feedback"`
	CreateTime time.Time `gorm:"create_time"`
}

func (fd FeedbackDao) TableName() string {
	return "feedbacks"
}

func (fd FeedbackDao) Create(f FeedbackDao) (int64, error) {
	err := db.Create(&f).Error
	return f.ID, err
}

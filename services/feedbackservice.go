package services

import (
	"insight-api/dbs"
	"time"
)

func PostFeedback(uid int64, feedback string) error {
	fdDao := dbs.FeedbackDao{}
	_, err := fdDao.Create(dbs.FeedbackDao{
		UserId:     uid,
		Feedback:   feedback,
		CreateTime: time.Now(),
	})
	return err
}

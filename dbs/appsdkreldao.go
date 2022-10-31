package dbs

type AppSdkRelDao struct {
	ID    int    `gorm:"primary_key"`
	AppId int64  `gorm:"app_id"`
	SdkId string `gorm:"sdk_id"`
}

func (sdk AppSdkRelDao) TableName() string {
	return "app_sdk_rel"
}

func (app AppSdkRelDao) Create(item AppSdkRelDao) error {
	err := db.Create(&item).Error
	return err
}

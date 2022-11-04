package dbs

import "time"

/*
"installedPublisherCount": 5760,

	    "uid": "a5f7d926-4d00-472b-9d3e-6b6c26204291",
	    "foundedYear": 2019,
	    "addressArea": [
	        "中国",
	        "贵州"
	    ],
	    "financingRound": "",
	    "servicesCount": 10,
	    "installedProductCount": 18447,
	    "installedWebsiteCount": 169,
	    "logoURL": "https://company-icon.oss-cn-hangzhou.aliyuncs.com/b65057e7488480e2acf5439078fa90d8.jpg",
	    "title": "华为云计算技术有限公司",
	    "industry": "信息传输、软件和信息技术服务业"
	},
*/
type ServiceProviderDao struct {
	ID             string    `gorm:"id"`
	Title          string    `gorm:"title"`
	LogoUrl        string    `gorm:"logo_url"`
	Industry       string    `gorm:"industry"`
	FoundedYear    string    `gorm:"founded_time"`
	AddressArea    string    `gorm:"address_area"`
	CreateTime     time.Time `gorm:"create_time"`
	FinancingRound string    `gorm:"financing_round"`
	AppCount       int       `gorm:"app_count"`
	WebsiteCount   int       `gorm:"website_count"`
	ServiceCount   int       `gorm:"service_count"`
	DeveloperCount int       `gorm:"developer_count"`

	Address     string `gorm:"address"`
	Website     string `gorm:"website"`
	Description string `gorm:"description"`
	Email       string `gorm:"email"`
}

func (dev ServiceProviderDao) TableName() string {
	return "service_providers"
}

func (dev ServiceProviderDao) FindById(id string) (*ServiceProviderDao, error) {
	var appItem ServiceProviderDao
	err := db.Where("id=?", id).Take(&appItem).Error
	if err != nil {
		return nil, err
	}
	//  else if err == gorm.ErrRecordNotFound {
	// 	return nil, nil
	// }
	return &appItem, nil
}

func (dev ServiceProviderDao) QueryListByPage(keyword string, page, count int) ([]*ServiceProviderDao, error) {
	var items []*ServiceProviderDao
	err := db.Order("app_count desc").Limit(count).Offset((page - 1) * count).Find(&items).Error
	return items, err
}

func (dev ServiceProviderDao) Create(d ServiceProviderDao) error {
	return db.Create(&d).Error
}

package services

func GetHomeInfo() *HomeInfo {
	return &HomeInfo{
		TotalAppCount:       10000,
		TotalDeveloperCount: 100,
	}
}

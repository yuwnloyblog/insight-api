package services

func GetHomeInfo() *HomeInfo {

	return &HomeInfo{
		TotalAppCount:       808038,
		TotalDeveloperCount: 923604,
		Sotre: map[string]int{
			"apple":       379842,
			"huawei":      161159,
			"xiaomi":      66858,
			"a360":        66281,
			"google_play": 133898,
		},
		Platform: map[string]int{
			"android": 428196,
			"ios":     379842,
		},
	}
}

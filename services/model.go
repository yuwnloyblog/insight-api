package services

import "fmt"

type HomeInfo struct {
	TotalAppCount       int            `json:"app_count"`       //总应用数量
	TotalDeveloperCount int            `json:"developer_count"` //总开发者数量
	Sotre               map[string]int `json:"store"`
	History             map[string]int `json:"history"`
}

type Apps struct {
	Items   []*App `json:"items"`
	HasMore bool   `json:"has_more"`
}

type App struct {
	Id            string     `json:"id"`
	Title         string     `json:"title"`
	BundleId      string     `json:"bundle_id"` //包名
	Platform      string     `json:"platform"`  //平台，iOS/Android
	Channel       string     `json:"channel"`   //苹果，华为，小米等
	Website       string     `json:"website"`
	Description   string     `json:"desc"`
	ReleaseDate   int64      `json:"release_date"`
	Developer     *Developer `json:"developer"`
	Size          int64      `json:"size"`
	CreateTime    int64      `json:"create_time"`
	LogoUrl       string     `json:"logo_url"`
	Category      string     `json:"category"`
	LatestVersion string     `json:"latest_version"`
	CountryCode   string     `json:"country_code"`

	Sdks []*SdkInfo `json:"sdks"`
}

type Developers struct {
	Items   []*Developer `json:"items"`
	HasMore bool         `json:"has_more"`
}

type Developer struct {
	Id             string `json:"id"`
	Title          string `json:"title"`
	Trade          string `json:"trade"`
	FoundedTime    string `json:"founded_time"`
	AddressArea    string `json:"address_area"`
	FinancingRound string `json:"financing_round"`
	LogoUrl        string `json:"logo_url"`
}

type SdkInfo struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Platforms string     `json:"platform"`
	Category  string     `json:"category"`
	Developer *Developer `json:"devloper"`
	LogoUrl   string     `json:"logoUrl"`
}

type User struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	WxUnionId string `json:"wx_unionid"`
	Phone     string `json:"phone"`
}

type CommonError struct {
	Code     int    `json:"code"`
	ErrorMsg string `json:"msg"`
}

func (err *CommonError) Error() string {
	return fmt.Sprintf("%d:%s", err.Code, err.ErrorMsg)
}

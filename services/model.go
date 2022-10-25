package services

type HomeInfo struct {
	TotalAppCount       int            `json:"app_count"`       //总应用数量
	TotalDeveloperCount int            `json:"developer_count"` //总开发者数量
	Sotre               map[string]int `json:"store"`
	Platform            map[string]int `json:"platform"`
}

type Apps struct {
	Items    []*App    `json:"items"`
	PageInfo *PageInfo `json:"page_info"`
	HasMore  bool      `json:"has_more"`
}
type PageInfo struct {
	Page  int `json:"page"`
	Count int `json:"count"`
}

type App struct {
	Id                string     `json:"id"`
	Title             string     `json:"title"`
	BundleId          string     `json:"bundle_id,omitempty"` //包名
	Platform          string     `json:"platform,omitempty"`  //平台，iOS/Android
	Channel           string     `json:"channel"`             //苹果，华为，小米等
	Website           string     `json:"website"`
	Description       string     `json:"desc"`
	LatestReleaseDate string     `json:"latest_release_date,omitempty"`
	FirstReleaseDate  string     `json:"first_release_date,omitempty"`
	Developer         *Developer `json:"developer"`
	Size              int64      `json:"size"`
	CreateTime        int64      `json:"create_time"`
	LogoUrl           string     `json:"logo_url"`
	Category          string     `json:"category"`
	LatestVersion     string     `json:"latest_version,omitempty"`
	CountryCode       string     `json:"country_code,omitempty"`

	Sdks []*SdkInfo `json:"sdks,omitempty"`
}

type Developers struct {
	Items    []*Developer `json:"items"`
	PageInfo *PageInfo    `json:"page_info"`
	HasMore  bool         `json:"has_more"`
}

type Developer struct {
	Id             string `json:"id,omitempty"`
	Title          string `json:"title"`
	Trade          string `json:"trade,omitempty"`
	FoundedTime    string `json:"founded_time,omitempty"`
	AddressArea    string `json:"address_area,omitempty"`
	FinancingRound string `json:"financing_round,omitempty"`
	LogoUrl        string `json:"logo_url,omitempty"`
	AppCount       int    `json:"app_count,omitempty"`
}

type Sdks struct {
	Items []*SdkInfo `json:"items"`
}
type SdkInfo struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Platforms string     `json:"platform"`
	Category  string     `json:"category"`
	Developer *Developer `json:"devloper"`
	LogoUrl   string     `json:"logoUrl"`
}

type WxLoginResp struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrorCode  int    `json:"errcode"`
	ErrorMsg   string `json:"errmsg"`
}

type User struct {
	Id            string `json:"id"`
	NickName      string `json:"nick_name"`
	WxOpenid      string `json:"wx_openid"`
	Phone         string `json:"phone"`
	Avatar        string `json:"avatar"`
	Status        int    `json:"status"`
	City          string `json:"city,omitempty"`
	Country       string `json:"country,omitempty"`
	Language      string `json:"language,omitempty"`
	Province      string `json:"province,omitempty"`
	LatestPayTime int64  `json:"latest_pay_time"`
}

type LoginUserResp struct {
	Token    string       `json:"token"`
	NickName string       `json:"nick_name"`
	Avatar   string       `json:"avatar"`
	Status   int          `json:"status"`
	WxResp   *WxLoginResp `json:"wx_resp"`
}

type PayNotifyResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type PayBody struct {
	Mchid       string  `json:"mchid"`
	OutTradeNo  string  `json:"out_trade_no"`
	Appid       string  `json:"appid"`
	Description string  `json:"description"`
	NotifyUrl   string  `json:"notify_url"`
	Amount      *Amount `json:"amount"`
	Payer       *Payer  `json:"payer"`
}
type Amount struct {
	Total    int    `json:"total"`
	Currency string `json:"currency"`
}
type Payer struct {
	Openid string `json:"openid"`
}
type PayResp struct {
	PrepayId string `json:"prepay_id"`
}

type UserOrder struct {
	OrderNo     string `json:"order_no"`
	Amount      int64  `json:"amount"`
	FellowType  int    `json:"fellow_type"`
	Description string `json:"description"`
}

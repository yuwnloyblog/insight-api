package apis

import (
	"context"
	"encoding/json"
	"fmt"
	"insight-api/configures"
	"insight-api/services"
	"insight-api/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	payUtils "github.com/wechatpay-apiv3/wechatpay-go/utils"
)

func WxPayNotify(ctx *gin.Context) {
	bs, err := ctx.GetRawData()
	fmt.Println(string(bs), err)
	ctx.JSON(http.StatusOK, PayNotifyResp{
		Code:    "SUCCESS",
		Message: "",
	})
}

type PayNotifyResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WxPay(ctx *gin.Context) {
	if !checkLogin(ctx) {
		return
	}

	amountStr := ctx.PostForm("amount")
	amountInt, err := utils.ParseInt64(amountStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, services.GetError(services.ErrorCode_ParamErr))
		return
	}

	uid := ctx.GetInt64("uid")
	user, err := services.GetUserInfByCache(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, services.GetError(services.ErrorCode_UserDbReadFail))
		return
	}

	var (
		mchID                      string = "1632714503"                               // 商户号
		mchCertificateSerialNumber string = "6172DED25454955F959D2871961712DB0C757C69" // 商户证书序列号
		mchAPIv3Key                string = "szxidatodayHoogmsofthello7890357"         // 商户APIv3密钥
	)

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := payUtils.LoadPrivateKeyWithPath("conf/apiclient_key.pem")
	if err != nil {
		log.Fatal("load merchant private key error")
	}

	payCtx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(payCtx, opts...)
	if err != nil {
		log.Fatalf("new wechat pay client err:%s", err)
	}

	svc := jsapi.JsapiApiService{Client: client}
	// 得到prepay_id，以及调起支付所需的参数和签名
	resp, result, err := svc.PrepayWithRequestPayment(ctx,
		jsapi.PrepayRequest{
			Appid:       core.String(configures.Config.Wx.AppId),
			Mchid:       core.String(mchID),
			Description: core.String("隼见-会员-月付"),
			OutTradeNo:  core.String(generateTradeNo()),
			Attach:      core.String("自定义数据说明"),
			NotifyUrl:   core.String("https://api.lwoowl.cn/user/wx_pay_notify"),
			Amount: &jsapi.Amount{
				Total: core.Int64(amountInt),
			},
			Payer: &jsapi.Payer{
				Openid: core.String(user.WxOpenid),
			},
		},
	)

	if err == nil {
		bs, _ := json.Marshal(resp)
		fmt.Println(string(bs))
		bs, _ = json.Marshal(result)
		fmt.Println("result:", string(bs))
		ctx.JSON(http.StatusOK, resp)
		return
	} else {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, services.GetError(services.ErrorCode_PrepayCallErr))
	}
}

func generateTradeNo() string {
	str := fmt.Sprintf("%s%s", utils.TimeFormat(time.Now()), utils.GetClearUuid())
	return str[:20]
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

package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"insight-api/configures"
	"insight-api/dbs"
	"insight-api/utils"
	"log"
	"math/rand"
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	payUtils "github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var (
	mchID                      string = "1632714503"                               // 商户号
	mchCertificateSerialNumber string = "6172DED25454955F959D2871961712DB0C757C69" // 商户证书序列号
	mchAPIv3Key                string = "szxidatodayHoogmsofthello7890357"         // 商户APIv3密钥
)
var wxPayClient *core.Client

func getWxPayClient(ctx context.Context) (*core.Client, error) {
	if wxPayClient != nil {
		return wxPayClient, nil
	}
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := payUtils.LoadPrivateKeyWithPath("conf/apiclient_key.pem")
	if err != nil {
		log.Fatal("load merchant private key error")
		return nil, err
	}

	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Fatalf("new wechat pay client err:%s", err)
		return nil, err
	}
	return client, nil
}

func HandlePayNotify(jsonBs []byte) error {
	notifyReqStr := string(jsonBs)
	var notifyReq notify.Request
	err := json.Unmarshal(jsonBs, &notifyReq)
	if err != nil {
		return err
	}
	// 处理通知内容
	if notifyReq.EventType == "TRANSACTION.SUCCESS" && notifyReq.Resource != nil {
		//解密数据
		transactionStr, err := payUtils.DecryptAES256GCM(mchAPIv3Key, notifyReq.Resource.AssociatedData, notifyReq.Resource.Nonce, notifyReq.Resource.Ciphertext)
		if err != nil {
			return err
		}
		//解析transaction
		var transaction payments.Transaction
		err = json.Unmarshal([]byte(transactionStr), &transaction)
		if err != nil {
			return err
		}
		//更新订单状态
		UpdateOrderStatus(*transaction.OutTradeNo, *transaction.Payer.Openid, 1, notifyReqStr, transactionStr)
	} else {
		return errors.New("支付不成功！" + notifyReq.EventType)
	}
	return nil
}

func WxPayCall(fllowType, _amount int, wxopenid string, uid int64) (*jsapi.PrepayWithRequestPaymentResponse, *core.APIResult, error) {
	ctx := context.Background()
	client, err := getWxPayClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	desc, amountInt := getFllowParam(fllowType, _amount)
	svc := jsapi.JsapiApiService{Client: client}

	orderNo := GenerateTradeNo(uid)
	maxTryCount := 3
	for {
		//存储订单
		_, err = SaveWxOrder(uid, UserOrder{
			OrderNo:     orderNo,
			Amount:      amountInt,
			FellowType:  fllowType,
			Description: desc,
		})
		if err == nil {
			break
		} else {
			orderNo = GenerateTradeNo(uid)
		}
		if maxTryCount <= 0 {
			return nil, nil, err
		}
		maxTryCount--
	}

	// 得到prepay_id，以及调起支付所需的参数和签名
	resp, result, err := svc.PrepayWithRequestPayment(ctx,
		jsapi.PrepayRequest{
			Appid:       core.String(configures.Config.Wx.AppId),
			Mchid:       core.String(mchID),
			Description: core.String(desc),
			OutTradeNo:  core.String(orderNo),
			Attach:      core.String(""),
			NotifyUrl:   core.String("https://api.lwoowl.cn/user/wx_pay_notify"),
			Amount: &jsapi.Amount{
				Total: core.Int64(amountInt),
			},
			Payer: &jsapi.Payer{
				Openid: core.String(wxopenid),
			},
		},
	)

	return resp, result, err
}

func getFllowParam(fllowType, amount int) (string, int64) {
	if fllowType == dbs.UserStatus_MONTH_PAY {
		return "隼见-会员-月付", 3000
	} else if fllowType == dbs.UserStatus_SEASON_PAY {
		return "隼见-会员-季付", 8500
	} else if fllowType == dbs.UserStatus_HALFYEAR_PAY {
		return "隼见-会员-半年付", 16000
	} else if fllowType == dbs.UserStatus_YEAR_PAY {
		return "隼见-会员-年付", 30000
	} else {
		return "隼见-会员-捐赠", int64(amount)
	}
}

func GenerateTradeNo(uid int64) string {
	uidSit := uid
	uidSit = uidSit % 1000000000
	rand.Seed(time.Now().UnixMilli())
	str := fmt.Sprintf("%s%09d%05d", utils.TimeFormat(time.Now()), uidSit, rand.Intn(100000))
	return str
}

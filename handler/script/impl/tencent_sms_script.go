package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	tErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"

	"github.com/ecodeclub/notify-go/common/domain/sms"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	smsSrv "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type TencentSmsScript struct {
	account TencentAccount
}

// 腾讯云账号样例：{"url":"sms.tencentcloudapi.com","region":"ap-guangzhou","secretId":"AKIDhDxxxxxxxx1WljQq","secretKey":"B4hwww39yxxxrrrrgxyi","smsSdkAppId":"1423123125","templateId":"1182097","signName":"Java3y公众号","supplierId":10,"supplierName":"腾讯云","scriptName":"TencentSmsScript"}
// 云片账号配置样例：{"url":"https://sms.yunpian.com/v2/sms/tpl_batch_send.json","apikey":"caffff8234234231b5cd7","tpl_id":"523333332","supplierId":20,"supplierName":"云片","scriptName":"YunPianSmsScript"}
func (t *TencentSmsScript) Send(ctx context.Context, smsParam sms.SmsParam) []sms.SMSRecord {
	var records []sms.SMSRecord
	err := json.Unmarshal([]byte(smsParam.AccountConfig), &t.account)
	if err != nil {
		return records
	}

	err = t.sendReq(ctx, smsParam)
	if err != nil {
	}
	return records
}

func (t *TencentSmsScript) sendReq(ctx context.Context, smsParam sms.SmsParam) error {
	credential := common.NewCredential(t.account.SecretID, t.account.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = t.account.URL
	client, _ := smsSrv.NewClient(credential, t.account.Region, cpf)
	response, err := client.SendSmsWithContext(ctx, &smsSrv.SendSmsRequest{
		PhoneNumberSet:   common.StringPtrs(smsParam.Phones),
		SmsSdkAppId:      common.StringPtr(t.account.SmsSdkAppID),
		SignName:         common.StringPtr(t.account.SignName),
		TemplateId:       common.StringPtr(t.account.TemplateID),
		TemplateParamSet: common.StringPtrs([]string{smsParam.Content}),
	})
	var tencentCloudSDKError *tErrors.TencentCloudSDKError
	if errors.As(err, &tencentCloudSDKError) {
		fmt.Printf("An API error has returned: %s", err)
	}

	// 输出json格式的字符串回包
	fmt.Printf("%s", response.ToJsonString())
	return err
}

func (t *TencentSmsScript) Pull(ctx context.Context, id int) []sms.SMSRecord {
	var records []sms.SMSRecord
	return records
}

type TencentAccount struct {
	URL         string `json:"url"`
	Region      string `json:"region"`
	SecretID    string `json:"secretId"`
	SecretKey   string `json:"secretKey"`
	SmsSdkAppID string `json:"smsSdkAppId"`
	TemplateID  string `json:"templateId"`
	SignName    string `json:"signName"`
	sms.SmsAccount
}

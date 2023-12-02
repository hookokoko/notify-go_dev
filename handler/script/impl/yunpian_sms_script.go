package impl

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/ecodeclub/notify-go/common/domain/sms"
	"github.com/ecodeclub/notify-go/pkg/ral"
)

type YunpianSmsScript struct {
	account YunpianAccount
	client  *ral.Client
}

func (y *YunpianSmsScript) Send(ctx context.Context, smsParam sms.SmsParam) []sms.SMSRecord {
	var records []sms.SMSRecord
	err := json.Unmarshal([]byte(smsParam.AccountConfig), &y.account)
	if err != nil {
		return records
	}

	err = y.sendReq(ctx, smsParam)
	if err != nil {
	}
	return records
}

func (y *YunpianSmsScript) sendReq(ctx context.Context, smsParam sms.SmsParam) error {
	req := ral.Request{
		Header: map[string]string{
			"content-type": "application/json;charset=utf-8",
		},
		Body: map[string]string{
			"apikey": y.account.ApiKey,
			"tpl_id": y.account.TplId,
			"mobile": strings.Join(smsParam.Phones, ","),
			"text":   smsParam.Content,
		},
	}
	var resp map[string]any
	err := y.client.Ral(ctx, "BatchSend", req, &resp, map[string]any{})
	return err
}

func (y *YunpianSmsScript) Pull(ctx context.Context, id int) []sms.SMSRecord {
	var records []sms.SMSRecord
	return records
}

// 账号参数示例
// {"url":"https://sms.yunpian.com/v2/sms/tpl_batch_send.json","apikey":"caffff8234234231b5cd7","tpl_id":"523333332","supplierId":20,"supplierName":"云片","scriptName":"YunPianSmsScript"}
type YunpianAccount struct {
	ApiKey string
	TplId  string
	Url    string
	sms.SmsAccount
}

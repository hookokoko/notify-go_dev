package impl

import (
	"context"
	"encoding/json"

	"github.com/ecodeclub/notify-go/common/domain/sms"
)

type AliSmsScript struct {
	account AliAccount
}

func (a *AliSmsScript) Send(ctx context.Context, smsParam sms.SmsParam) []sms.SMSRecord {
	var records []sms.SMSRecord
	err := json.Unmarshal([]byte(smsParam.AccountConfig), &a.account)
	if err != nil {
		return records
	}

	err = a.sendReq(ctx, smsParam)
	if err != nil {
	}
	return records
}

func (a *AliSmsScript) sendReq(ctx context.Context, smsParam sms.SmsParam) error {
	var err error
	return err
}

func (a *AliSmsScript) Pull(ctx context.Context, id int) []sms.SMSRecord {
	//TODO implement me
	panic("implement me")
}

type AliAccount struct {
	sms.SmsAccount
	UserName string
	Password string
}

package impl

import (
	"context"

	"github.com/ecodeclub/notify-go/common/domain/sms"
)

type TencentSmsScript struct {
}

func (t *TencentSmsScript) Send(ctx context.Context, smsParam sms.SmsParam) []sms.SMSRecord {
	//TODO implement me
	panic("implement me")
}

func (t *TencentSmsScript) Pull(ctx context.Context, id int) []sms.SMSRecord {
	//TODO implement me
	panic("implement me")
}

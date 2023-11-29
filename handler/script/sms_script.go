package script

import (
	"context"

	"github.com/ecodeclub/notify-go/common/domain/sms"
	"github.com/ecodeclub/notify-go/handler/script/impl"
)

var SmsScriptHandler = map[string]SMSScript{
	"Tencent": &impl.TencentSmsScript{},
}

type SMSScript interface {
	// Send 发送短信
	// @param smsParam
	// @return 渠道商发送接口返回值
	Send(ctx context.Context, smsParam sms.SmsParam) []sms.SMSRecord

	// Pull 拉取回执
	// @param id 渠道账号的ID
	// @return 渠道商回执接口返回值
	Pull(ctx context.Context, id int) []sms.SMSRecord
}

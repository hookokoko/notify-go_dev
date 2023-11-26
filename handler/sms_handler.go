package handler

import (
	"context"
	"fmt"

	"github.com/ecodeclub/notify-go/common/domain"
	"github.com/ecodeclub/notify-go/common/enum/channel_type"
	"github.com/ecodeclub/notify-go/repo"
	"github.com/ecodeclub/notify-go/storage"
	"gorm.io/gorm"
)

func init() {
	once.Do(func() {
		handlerHolder.put(channel_type.SMS.String(), NewSMSHandler(storage.MysqlDB()))
	})
}

type SMSHandler struct {
	channelCode channel_type.ChannelType
	*baseHandler
	// 其他字段
	channelAccountDao repo.IChannelAccountDao
}

func NewSMSHandler(db *gorm.DB) *SMSHandler {
	return nil
}

func (s *SMSHandler) Do(ctx context.Context, taskInfo domain.TaskInfo) bool {
	// 调用一些baseHandler的方法

	ok := s.handle(ctx, taskInfo)
	return ok
}

func (s *SMSHandler) handle(ctx context.Context, taskInfo domain.TaskInfo) bool {
	fmt.Println(">>>>>>>> sms handler")

	return true
}

//func (s *SMSHandler) getAccountConfig(ctx context.Context, sendAccount int) MailAccountConfig {
//	var c *MailAccountConfig
//	account, err := s.channelAccountDao.FindById(ctx, sendAccount)
//	if err != nil {
//		fmt.Println("get account error:", err)
//		return *c
//	}
//	err = json.Unmarshal([]byte(account.AccountConfig), c)
//	if err != nil {
//		fmt.Println("unmarshal account config error:", err)
//		return *c
//	}
//	return *c
//}

type SmsParma struct {
	MessageTemplateId int      `json:"message_template_id"`
	Phones            []string `json:"phones"`
	SendAccountId     int      `json:"send_account_id"`
	ScriptName        string   `json:"script_name"`
	Content           string   `json:"content"`
}

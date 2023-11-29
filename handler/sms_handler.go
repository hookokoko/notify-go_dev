package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/ecodeclub/notify-go/common/domain"
	"github.com/ecodeclub/notify-go/common/domain/sms"
	"github.com/ecodeclub/notify-go/common/enum/channel_type"
	"github.com/ecodeclub/notify-go/common/model"
	"github.com/ecodeclub/notify-go/handler/script"
	"github.com/ecodeclub/notify-go/repo"
	"gorm.io/gorm"
)

const AutoFlowRule = iota

//func init() {
//	once.Do(func() {
//		handlerHolder.put(channel_type.SMS.String(), NewSMSHandler(storage.MysqlDB()))
//	})
//}

type SMSHandler struct {
	channelCode channel_type.ChannelType
	*baseHandler
	// 其他字段
	channelAccountDao repo.IChannelAccountDao
}

func NewSMSHandler(db *gorm.DB) *SMSHandler {
	return &SMSHandler{}
}

func (s *SMSHandler) Do(ctx context.Context, taskInfo domain.TaskInfo) bool {
	// 调用一些baseHandler的方法

	ok := s.handle(ctx, taskInfo)
	return ok
}

func (s *SMSHandler) handle(ctx context.Context, taskInfo domain.TaskInfo) bool {
	fmt.Println(">>>>>>>> sms handler")
	smsParam := sms.SmsParam{
		Phones:            taskInfo.Receiver,
		Content:           s.getSMSContent(taskInfo),
		MessageTemplateId: taskInfo.MessageTemplateId,
	}

	messageTypeSmsConfig := s.loadBalance(s.getMessageTypeSmsConfig(ctx, taskInfo))
	smsParam.SendAccountId = messageTypeSmsConfig.SendAccount
	smsParam.ScriptName = messageTypeSmsConfig.ScriptName

	if handler, ok := script.SmsScriptHandler[messageTypeSmsConfig.ScriptName]; ok {
		records := handler.Send(ctx, smsParam)
		fmt.Println(records)
		return true
	}

	return false
}

func (s *SMSHandler) loadBalance(smsConfigs []sms.MessageTypeSmsConfig) sms.MessageTypeSmsConfig {
	if len(smsConfigs) == 0 {
		return sms.MessageTypeSmsConfig{}
	}

	weightSum := 0
	for _, smsConfig := range smsConfigs {
		weightSum += smsConfig.Weights
	}

	// 随机生成一个在[1, weightSum]区间的整数
	idx := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(weightSum) + 1
	for i := 0; i < len(smsConfigs); i++ {
		idx -= smsConfigs[i].Weights
		if idx <= 0 {
			return smsConfigs[i]
		}
	}

	return sms.MessageTypeSmsConfig{}
}

func (s *SMSHandler) getSMSContent(taskInfo domain.TaskInfo) string {
	var smsContent model.SmsContentModel
	err := json.Unmarshal([]byte(taskInfo.Content), &smsContent)
	if err != nil {
		fmt.Println("json.Unmarshal err:", err)
		return ""
	}
	if smsContent.Url != "" {
		return fmt.Sprintf("%s %s", smsContent.Content, smsContent.Url)
	}
	return smsContent.Content
}

func (s *SMSHandler) getMessageTypeSmsConfig(ctx context.Context, taskInfo domain.TaskInfo) []sms.MessageTypeSmsConfig {
	var (
		smsConfig  []sms.MessageTypeSmsConfig
		smsAccount sms.SmsAccount
	)

	if taskInfo.SendAccount != AutoFlowRule {
		channelAccount, err := s.channelAccountDao.FindById(ctx, taskInfo.SendAccount)
		if err != nil {
			fmt.Println("smsHandler.getMessageTypeSmsConfig err:", err)
			return smsConfig
		}
		if channelAccount.IsDeleted != 1 {
			if err := json.Unmarshal([]byte(channelAccount.AccountConfig), &smsAccount); err != nil {
				fmt.Println("json.Unmarshal err:", err)
				return smsConfig
			}
		}
		smsConfig = []sms.MessageTypeSmsConfig{
			{
				Weights:     100,
				SendAccount: taskInfo.SendAccount,
				ScriptName:  smsAccount.ScriptName,
			},
		}
		return smsConfig
	}
	var res []map[string]any
	str := `[
{"message_type_10":[{"weights":99,"scriptName":"TencentSmsScript"},{"weights":1,"scriptName":"YunPianSmsScript"}]},
{"message_type_20":[{"weights":99,"scriptName":"TencentSmsScript"},{"weights":1,"scriptName":"YunPianSmsScript"}]},
{"message_type_30":[{"weights":20,"scriptName":"TencentSmsScript"}]},
{"message_type_40":[{"weights":20,"scriptName":"TencentSmsScript"}]},
{"message_type_50":[{"weights":20,"scriptName":"LinTongSmsScript"}]}]`
	err := json.Unmarshal([]byte(str), &res)
	if err != nil {
		fmt.Println("json.Unmarshal err:", err)
		return smsConfig
	}

	for _, item := range res {
		for k, v := range item {
			if fmt.Sprintf("message_type_%d", taskInfo.MsgType) == k {
				b, _ := json.Marshal(v)
				err := json.Unmarshal(b, &smsConfig)
				if err != nil {
					fmt.Println("json.Unmarshal err:", err)
				}
			}
		}
	}
	return smsConfig
}

package notify_go_dev

import (
	"context"
	"fmt"
	"strings"

	"github.com/ecodeclub/notify-go/common/domain"
	"github.com/ecodeclub/notify-go/common/pipeline"
	"github.com/ecodeclub/notify-go/repo"
	"github.com/pborman/uuid"
	"gorm.io/gorm"
)

// 前置检查
var preCheck pipeline.Filter[*Delivery] = func(next pipeline.HandlerFunc[*Delivery]) pipeline.HandlerFunc[*Delivery] {
	return func(ctx context.Context, deli *Delivery) error {
		fmt.Println("pre check begin")
		/*
			1. 没有传入 消息模板Id 或者 messageParam
			2. 空的receiver
			3. receiver 大于100的请求
		*/
		err := next(ctx, deli)
		fmt.Println("pre check end")
		return err
	}
}

// AssembleFilterBuilder 组装发送参数filter的构建
type AssembleFilterBuilder struct {
	templateDao repo.ITemplateDao
}

func NewAssembleFilterBuilder(db *gorm.DB) *AssembleFilterBuilder {
	return &AssembleFilterBuilder{
		templateDao: repo.NewTemplateDao(db),
	}
}

func (af *AssembleFilterBuilder) Build() pipeline.Filter[*Delivery] {
	return func(next pipeline.HandlerFunc[*Delivery]) pipeline.HandlerFunc[*Delivery] {
		return func(ctx context.Context, object *Delivery) error {
			fmt.Println("assemble begin")
			template, err := af.templateDao.FindById(ctx, object.MessageTemplateId)
			if err != nil {
				return err
			}
			for _, messageParam := range object.MessageParamList {
				taskInfo := domain.TaskInfo{
					BizId:             messageParam.BizId,
					MessageId:         uuid.NewUUID().String(), // 生成唯一id
					MessageTemplateId: object.MessageTemplateId,
					BusinessId:        0, // TODO 根据模版类型和当天日期生成
					Receiver:          strings.Split(messageParam.Receiver, ","),
					IdType:            template.IDType,
					SendChannel:       template.SendChannel,
					TemplateType:      template.TemplateType,
					MsgType:           template.MsgType,
					ShieldType:        template.ShieldType,
					SendAccount:       template.SendAccount,
					Content:           ReplacePlaceHolder(template.MsgContent, messageParam.Variables),
				}
				object.TaskInfo = append(object.TaskInfo, taskInfo)
			}
			err = next(ctx, object)
			fmt.Println("assemble end")
			return err
		}
	}
}

// 后置检查
var afterCheck pipeline.Filter[*Delivery] = func(next pipeline.HandlerFunc[*Delivery]) pipeline.HandlerFunc[*Delivery] {
	return func(ctx context.Context, deli *Delivery) error {
		fmt.Println("after check begin")
		/*
			1. 利用正则过滤掉不合法的接收者
		*/
		err := next(ctx, deli)
		fmt.Println("after check end")
		return err
	}
}

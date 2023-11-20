package experience

import (
	"context"
	"fmt"
	"strings"

	"github.com/ecodeclub/ekit/bean/option"
	"github.com/ecodeclub/notify-go/experience/common/domain"
	"github.com/ecodeclub/notify-go/experience/common/pipeline"
	"github.com/pborman/uuid"
)

// Delivery {"code":"send","messageParam": {"bizId":null,"extra":null,"receiver":"123@qq.com","variables":null},"messageTemplateId":17,"recallMessageId":null}
type Delivery struct {
	// 消息模板Id
	MessageTemplateId int `json:"message_template_id"`
	// 请求参数
	MessageParamList []MessageParam `json:"message_param_list"`
	// 发送任务的信息
	TaskInfo    []domain.TaskInfo `json:"task_info"`
	sendService SendService
}

type MessageParam struct {
	//业务消息发送Id, 用于链路追踪, 若不存在, austin 则生成一个消息Id
	BizId string `json:"biz_id"`
	// 接收者。多个用 逗号 分隔开【不能大于100个】必传
	Receiver string `json:"receiver"`
	// 消息内容中的可变部分(占位符替换) 可选
	Variables map[string]string `json:"variables"`
	// 扩展参数 可选
	Extra map[string]string `json:"extra"`
}

func WithSendService(srv SendService) option.Option[Delivery] {
	return func(d *Delivery) {
		d.sendService = srv
	}
}

func NewDelivery(messageParamList []MessageParam, templateId int, opts ...option.Option[Delivery]) Delivery {
	d := Delivery{
		MessageTemplateId: templateId,
		MessageParamList:  messageParamList,
		TaskInfo:          make([]domain.TaskInfo, 0),
		sendService:       NewDefaultSendImpl(),
	}
	option.Apply[Delivery](&d, opts...)
	return d
}

func (deli *Delivery) Send(ctx context.Context, filters ...pipeline.Filter[*Delivery]) error {
	filters = append(filters, preCheck, assemble, afterCheck)
	sendHandlerFunc := NewSendFuncBuilder().Build()
	processPipeline := pipeline.FilterChain[*Delivery](filters...).Then(sendHandlerFunc)
	err := processPipeline(ctx, deli)
	return err
}

// 前置检查
var preCheck pipeline.Filter[*Delivery] = func(next pipeline.HandlerFunc[*Delivery]) pipeline.HandlerFunc[*Delivery] {
	return func(ctx context.Context, deli *Delivery) error {
		fmt.Println("pre check begin")
		err := next(ctx, deli)
		fmt.Println("pre check end")
		return err
	}
}

// 组装发送参数
var assemble pipeline.Filter[*Delivery] = func(next pipeline.HandlerFunc[*Delivery]) pipeline.HandlerFunc[*Delivery] {
	return func(ctx context.Context, deli *Delivery) error {
		fmt.Println("assemble begin")
		template := MessageTemplate{
			ID:          111233,
			Name:        "测试模版名",
			IDType:      50,
			SendChannel: 30,
			MsgContent:  `${name}先生/女士,你好，这是一条短信内容哦`,
		} // TODO getFrom(deli.MessageTemplateId)
		// 所有的taskinfo一起序列化，一起发送
		for _, messageParam := range deli.MessageParamList {
			taskInfo := domain.TaskInfo{
				BizId:             messageParam.BizId,
				MessageId:         uuid.NewUUID().String(), // 生成唯一id
				MessageTemplateId: deli.MessageTemplateId,
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
			deli.TaskInfo = append(deli.TaskInfo, taskInfo)
		}
		err := next(ctx, deli)
		fmt.Println("assemble end")
		return err
	}
}

// 后置检查
var afterCheck pipeline.Filter[*Delivery] = func(next pipeline.HandlerFunc[*Delivery]) pipeline.HandlerFunc[*Delivery] {
	return func(ctx context.Context, deli *Delivery) error {
		fmt.Println("after check begin")
		err := next(ctx, deli)
		fmt.Println("after check end")
		return err
	}
}

type SendFuncBuilder struct {
}

func NewSendFuncBuilder() *SendFuncBuilder {
	return &SendFuncBuilder{}
}

func (sb *SendFuncBuilder) Build() pipeline.HandlerFunc[*Delivery] {
	return func(ctx context.Context, object *Delivery) error {
		fmt.Println("send begin")
		err := object.sendService.Send(ctx, object.TaskInfo)
		fmt.Println("send end")
		return err
	}
}

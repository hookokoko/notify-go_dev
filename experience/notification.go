// Copyright 2021 ecodeclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package experience

import (
	"context"
	"github.com/ecodeclub/ekit/bean/option"
	"github.com/pborman/uuid"
	"strings"
)

// Delivery {"code":"send","messageParam": {"bizId":null,"extra":null,"receiver":"123@qq.com","variables":null},"messageTemplateId":17,"recallMessageId":null}
type Delivery struct {
	// 消息模板Id
	MessageTemplateId int `json:"message_template_id"`
	// 请求参数
	MessageParamList []MessageParam `json:"message_param_list"`
	// 发送任务的信息
	TaskInfo []TaskInfo `json:"task_info"`
	// 发送方式，默认异步
	IsSync      bool
	SendService SendService
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

type TaskInfo struct {
	//业务消息发送Id, 用于链路追踪, 若不存在, 则使用 messageId
	BizId string `json:"biz_id"`
	//消息唯一Id(数据追踪使用)
	MessageId string `json:"message_id"`
	//消息模板Id
	MessageTemplateId int `json:"message_template_id"`
	//业务Id(数据追踪使用)
	BusinessId int `json:"business_id"`
	//接收者
	Receiver []string `json:"receiver"`
	//发送的Id类型
	IdType int `json:"id_type"`
	//发送渠道
	SendChannel int `json:"send_channel"`
	//模板类型
	TemplateType int `json:"template_type"`
	//消息类型
	MsgType int `json:"msg_type"`
	// 屏蔽类型
	ShieldType int `json:"shield_type"`
	// 发送文案, 似乎不需要转成对应channel的struct
	// message_template表存储的content是JSON(所有内容都会塞进去)
	Content string `json:"content"`
	// 发送账号（邮件下可有多个发送账号、短信可有多个发送账号..）
	SendAccount int `json:"send_account"`
}

func WithIsSync(isSync bool) option.Option[Delivery] {
	return func(d *Delivery) {
		d.IsSync = isSync
	}
}

func NewDelivery(messageParamList []MessageParam, opts ...option.Option[Delivery]) Delivery {
	d := Delivery{
		MessageParamList: messageParamList,
		TaskInfo:         make([]TaskInfo, 0),
	}
	option.Apply[Delivery](&d, opts...)

	if d.IsSync {
		d.SendService = NewSendSyncService()
	} else {
		d.SendService = NewSendAsyncService()
	}

	return d
}

type SendFunc func(context.Context, *Delivery) error

type SendMiddleware func(SendFunc) SendFunc

func (deli *Delivery) Send(ctx context.Context, mls ...SendMiddleware) error {
	// 前置检查
	var preCheck SendMiddleware = func(next SendFunc) SendFunc {
		return func(ctx context.Context, deli *Delivery) error {
			// 这里可以干点啥
			err := next(ctx, deli)
			// 这里也可以干点啥
			return err
		}
	}

	// 组装发送参数
	var assemble SendMiddleware = func(next SendFunc) SendFunc {
		return func(ctx context.Context, deli *Delivery) error {
			template := MessageTemplate{} // TODO getFrom(deli.MessageTemplateId)
			// 所有的taskinfo一起序列化，一起发送
			for _, messageParam := range deli.MessageParamList {
				taskInfo := TaskInfo{
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
			return next(ctx, deli)
		}
	}

	// 后置检查
	var afterCheck SendMiddleware = func(next SendFunc) SendFunc {
		return func(ctx context.Context, deli *Delivery) error {
			return next(ctx, deli)
		}
	}

	// 消息实际发送
	var send SendFunc = func(ctx context.Context, deli *Delivery) error {
		err := deli.SendService.Send(ctx, deli.TaskInfo)
		return err
	}

	// 构造调用链
	// (...(preCheck(assemble(afterCheck(send)))))
	totalMls := append([]SendMiddleware{afterCheck, assemble, preCheck}, mls...)

	// 组装其他中间件
	for i := len(totalMls) - 1; i > 0; i-- {
		send = mls[i](send)
	}

	return send(ctx, deli)
}

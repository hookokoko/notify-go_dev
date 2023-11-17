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
	"fmt"
	"strings"
)

func ReplacePlaceHolder(tpl string, variables map[string]string) string {
	for k, v := range variables {
		tpl = strings.ReplaceAll(tpl, fmt.Sprintf("${%s}", k), v)
	}
	return tpl
}

type MessageTemplate struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`             // 模板标题
	AuditStatus    int    `json:"audit_status"`     // 当前消息审核状态: 10.待审核 20.审核成功 30.被拒绝
	FlowId         int    `json:"flow_id"`          // 工单ID（审核模板走工单）
	MsgStatus      int    `json:"msg_status"`       // 消息状态
	CronTaskId     int    `json:"cron_task_id"`     // 定时任务Id(由xxl-job返回)
	IDType         int    `json:"id_type"`          // 消息的发送ID类型: 10. userId 20.did 30.手机号 40.openId 50.email 60.企业微信userId
	SendChannel    int    `json:"send_channel"`     // 消息发送渠道: 10.IM 20.Push 30.短信 40.Email 50.公众号 60.小程序 70.企业微信
	TemplateType   int    `json:"template_type"`    // 模板类型 10.运营类 20.技术类接口调用
	ShieldType     int    `json:"shield_type"`      // 屏蔽类型 10.夜间不屏蔽 20.夜间屏蔽 30.夜间屏蔽(次日早上9点发送)
	MsgType        int    `json:"msg_type"`         // 消息类型 10.通知类消息 20.营销类消息 30.验证码类消息
	ExpectPushTime string `json:"expect_push_time"` // 推送消息的时间, 0: 立即发送, else: crontab 表达式
	MsgContent     string `json:"msg_content"`      // 消息内容 占位符用{$var}表示
	SendAccount    int    `json:"send_account"`     // 发送账号 一个渠道下可存在多个账号
	Creator        string `json:"creator"`          // 创建者
	Updater        string `json:"updater"`          // 更新者
	Auditor        string `json:"auditor"`          // 审核人
	Team           string `json:"team"`             // 业务方团队
	Proposer       string `json:"proposer"`         // 业务方
	IsDeleted      int    `json:"is_deleted"`       // 是否删除: 0.不删除 1.删除
	Created        int32  `json:"created"`          // 创建时间,单位 s
	Updated        int32  `json:"updated"`          // 更新时间,单位 s
}

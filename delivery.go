package notify_go_dev

import (
	"github.com/ecodeclub/notify-go/common/domain"
)

// Delivery {"code":"send","messageParam": {"bizId":null,"extra":null,"receiver":"123@qq.com","variables":null},"messageTemplateId":17,"recallMessageId":null}
type Delivery struct {
	// 消息模板Id
	MessageTemplateId int `json:"message_template_id"`
	// 请求参数
	MessageParamList []MessageParam `json:"message_param_list"`
	// 发送任务的信息
	TaskInfo []domain.TaskInfo `json:"task_info"`
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

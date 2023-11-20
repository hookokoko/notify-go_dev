package domain

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

package sms

type SmsParam struct {
	MessageTemplateId int      `json:"message_template_id"`
	Phones            []string `json:"phones"`
	SendAccountId     int      `json:"send_account_id"`
	ScriptName        string   `json:"script_name"`
	Content           string   `json:"content"`
}

type MessageTypeSmsConfig struct {
	// 权重(决定着流量的占比)
	Weights int `json:"weights"`
	//短信模板若指定了账号，则该字段有值
	SendAccount int `json:"sendAccount"`
	//script名称
	ScriptName string `json:"scriptName"`
}

type SmsAccount struct {
	// 标识渠道商Id
	SupplierId int
	// 标识渠道商名字
	SupplierName string
	//【重要】类名，定位到具体的处理"下发"/"回执"逻辑, 依据ScriptName对应具体的某一个短信账号
	ScriptName string
}

type SMSRecord struct {
	ID                int64  `gorm:"column:id;primaryKey" json:"id"`
	MessageTemplateID int64  `gorm:"column:message_template_id;not null" json:"message_template_id"` // 消息模板ID
	Phone             int64  `gorm:"column:phone;not null" json:"phone"`                             // 手机号
	SupplierID        int32  `gorm:"column:supplier_id;not null" json:"supplier_id"`                 // 发送短信渠道商的ID
	SupplierName      string `gorm:"column:supplier_name;not null" json:"supplier_name"`             // 发送短信渠道商的名称
	MsgContent        string `gorm:"column:msg_content;not null" json:"msg_content"`                 // 短信发送的内容
	SeriesID          string `gorm:"column:series_id;not null" json:"series_id"`                     // 下发批次的ID
	ChargingNum       int32  `gorm:"column:charging_num;not null" json:"charging_num"`               // 计费条数
	ReportContent     string `gorm:"column:report_content;not null" json:"report_content"`           // 回执内容
	Status            int32  `gorm:"column:status;not null" json:"status"`                           // 短信状态： 10.发送 20.成功 30.失败
	SendDate          int32  `gorm:"column:send_date;not null" json:"send_date"`                     // 发送日期：20211112
	CreateAt          int64  `gorm:"column:create_at;not null" json:"create_at"`                     // 创建时间
	UpdateAt          int64  `gorm:"column:update_at;not null" json:"update_at"`                     // 更新时间
}

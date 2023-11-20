package model

// SmsContentModel 短信内容模型
type SmsContentModel struct {
	Content string `json:"content"`
	Url     string `json:"url"`
}

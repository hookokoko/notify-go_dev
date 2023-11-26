package model

// PushContentModel 通知栏内容模型
type PushContentModel struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Url     string `json:"url"`
}

package model

type EmailContentModel struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Url     string `json:"url"`
}

package handler

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/smtp"
	"sync"

	"github.com/ecodeclub/notify-go/common/domain"
	"github.com/ecodeclub/notify-go/common/enum/channel_type"
	"github.com/ecodeclub/notify-go/repo"
	"github.com/jordan-wright/email"
	"gorm.io/gorm"
)

var once *sync.Once

//func init() {
//	once.Do(func() {
//		handlerHolder.put(channel_type.EMAIL.String(), NewEmailHandler(storage.MysqlDB()))
//	})
//}

type EmailHandler struct {
	channelCode channel_type.ChannelType
	*baseHandler
	// 其他字段
	email             *email.Email
	channelAccountDao repo.IChannelAccountDao
}

func NewEmailHandler(db *gorm.DB) *EmailHandler {
	return &EmailHandler{
		channelCode:       channel_type.EMAIL,
		baseHandler:       newBaseHandler(handlerHolder, nil, nil),
		email:             email.NewEmail(),
		channelAccountDao: repo.NewChannelAccountDao(db),
	}
}

func (e *EmailHandler) Do(ctx context.Context, taskInfo domain.TaskInfo) bool {
	// 调用一些baseHandler的方法
	ok := e.handle(ctx, taskInfo)
	return ok
}

func (e *EmailHandler) handle(ctx context.Context, taskInfo domain.TaskInfo) bool {
	fmt.Println(">>>>>>>> email handler")
	cfg := e.getAccountConfig(ctx, taskInfo.SendAccount)

	e.email = &email.Email{
		From:        cfg.From,
		To:          taskInfo.Receiver,
		Subject:     "",
		HTML:        []byte(taskInfo.Content),
		Attachments: nil, // TODO
	}

	err := e.email.SendWithTLS(cfg.SmtpHostAddr,
		smtp.PlainAuth("", cfg.SmtpUser, cfg.SmtpPwd, cfg.SmtpHost),
		&tls.Config{ServerName: cfg.SmtpHost},
	)
	if err != nil {
		fmt.Println("send email error:", err)
		return false
	}
	return true
}

func (e *EmailHandler) getAccountConfig(ctx context.Context, sendAccount int) MailAccountConfig {
	var c *MailAccountConfig
	account, err := e.channelAccountDao.FindById(ctx, sendAccount)
	if err != nil {
		fmt.Println("get account error:", err)
		return *c
	}
	err = json.Unmarshal([]byte(account.AccountConfig), c)
	if err != nil {
		fmt.Println("unmarshal account config error:", err)
		return *c
	}
	return *c
}

type MailAccountConfig struct {
	SmtpHostAddr string `json:"smtp_host_addr"`
	SmtpHost     string `json:"smtp_host"`
	SmtpPort     int    `json:"smtp_port"`
	SmtpUser     string `json:"smtp_user"`
	SmtpPwd      string `json:"smtp_pwd"`
	From         string `json:"from"`
}

package handler

import (
	"context"
	"fmt"

	"github.com/ecodeclub/notify-go/experience/common/domain"
	"github.com/ecodeclub/notify-go/experience/common/enum/channel_type"
)

func init() {
	// TODO once
	handlerHolder.put(channel_type.EMAIL.String(), NewEmailHandler())
}

type EmailHandler struct {
	channelCode channel_type.ChannelType
	*baseHandler
	// 其他字段
}

func NewEmailHandler() *EmailHandler {
	return &EmailHandler{
		channelCode: channel_type.EMAIL,
		baseHandler: newBaseHandler(handlerHolder, nil, nil),
	}
}

func (e *EmailHandler) Do(ctx context.Context, taskInfo domain.TaskInfo) bool {
	// 调用一些baseHandler的方法
	ok := e.handle(ctx, taskInfo)
	return ok
}

func (e *EmailHandler) handle(ctx context.Context, taskInfo domain.TaskInfo) bool {
	fmt.Println(">>>>>>>> email handler")
	return true
}

func (e *EmailHandler) getAccountConfig() bool {
	return true
}

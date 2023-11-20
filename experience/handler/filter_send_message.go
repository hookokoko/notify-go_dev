package handler

import (
	"context"

	"github.com/ecodeclub/notify-go/experience/common/domain"
	"github.com/ecodeclub/notify-go/experience/common/enum/channel_type"
	"github.com/ecodeclub/notify-go/experience/common/pipeline"
)

//type SendMessageFilter struct{}
//
//func (f *SendMessageFilter) Process(ctx context.Context, taskInfo *domain.TaskInfo) error {
//	h := handlerHolder.route(channel_type.ChannelType(taskInfo.SendChannel).String())
//	h.Do(ctx, *taskInfo)
//	return nil
//}
//
//func (f *SendMessageFilter) Before(ctx context.Context, taskInfo *domain.TaskInfo) error {
//	return nil
//}
//
//func (f *SendMessageFilter) After(ctx context.Context, taskInfo *domain.TaskInfo) error {
//	return nil
//}

type SendMsgHandlerBuilder struct {
}

func NewSendMsgHandlerBuilder() *SendMsgHandlerBuilder {
	return &SendMsgHandlerBuilder{}
}

func (f *SendMsgHandlerBuilder) Build() pipeline.HandlerFunc[*domain.TaskInfo] {
	return func(ctx context.Context, object *domain.TaskInfo) error {
		h := handlerHolder.route(channel_type.ChannelType(object.SendChannel).String())
		h.Do(ctx, *object)
		return nil
	}
}

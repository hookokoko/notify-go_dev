package handler

import (
	"context"

	"github.com/ecodeclub/notify-go/common/domain"
	"github.com/ecodeclub/notify-go/common/enum/channel_type"
	"github.com/ecodeclub/notify-go/common/pipeline"
)

type SendMsgHandlerBuilder struct {
}

func NewSendMsgHandlerBuilder() *SendMsgHandlerBuilder {
	return &SendMsgHandlerBuilder{}
}

func (f *SendMsgHandlerBuilder) Build() pipeline.HandlerFunc[*domain.TaskInfo] {
	return func(ctx context.Context, object *domain.TaskInfo) error {
		h := handlerHolder.route(channel_type.ChannelType(object.SendChannel).String())
		// TODO 如果h的限流参数不为空, do sth
		h.Do(ctx, *object)
		return nil
	}
}

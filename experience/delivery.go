package experience

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit/bean/option"
	"github.com/ecodeclub/notify-go/experience/common/domain"
	"github.com/ecodeclub/notify-go/experience/common/pipeline"
	"gorm.io/gorm"
)

type Notification struct {
	db          *gorm.DB
	sendService SendService
	Delivery
}

func WithSendService(srv SendService) option.Option[Notification] {
	return func(d *Notification) {
		d.sendService = srv
	}
}

func NewNotification(messageParamList []MessageParam, templateId int, db *gorm.DB,
	opts ...option.Option[Notification]) Notification {
	n := Notification{
		Delivery: Delivery{
			MessageTemplateId: templateId,
			MessageParamList:  messageParamList,
			TaskInfo:          make([]domain.TaskInfo, 0),
		},
		sendService: NewDefaultSendImpl(),
		db:          db,
	}

	option.Apply[Notification](&n, opts...)
	return n
}

func (ni *Notification) Send(ctx context.Context, filters ...pipeline.Filter[*Delivery]) error {
	filters = append(filters,
		preCheck,
		NewAssembleFilterBuilder(ni.db).Build(),
		afterCheck,
	)
	send := NewSendFuncBuilder(ni.sendService).Build()
	processPipeline := pipeline.FilterChain[*Delivery](filters...).Then(send)
	err := processPipeline(ctx, &ni.Delivery)
	return err
}

type SendFuncBuilder struct {
	service SendService
}

func NewSendFuncBuilder(srv SendService) *SendFuncBuilder {
	return &SendFuncBuilder{
		service: srv,
	}
}

func (sb *SendFuncBuilder) Build() pipeline.HandlerFunc[*Delivery] {
	return func(ctx context.Context, object *Delivery) error {
		fmt.Println("send begin")
		err := sb.service.Send(ctx, object.TaskInfo)
		fmt.Println("send end")
		return err
	}
}

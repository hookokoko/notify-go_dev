package handler

import (
	"context"
	"fmt"
	"slices"

	"github.com/ecodeclub/notify-go/common/domain"
	"github.com/ecodeclub/notify-go/common/pipeline"
)

type DiscardFilterBuilder struct {
	discardMessageIds []int // TODO 配置文件读取
}

func NewDiscardFilterBuilder(discardMessageIds []int) *DiscardFilterBuilder {
	return &DiscardFilterBuilder{
		discardMessageIds: discardMessageIds,
	}
}

func (dfb *DiscardFilterBuilder) Build() pipeline.Filter[*domain.TaskInfo] {
	return func(next pipeline.HandlerFunc[*domain.TaskInfo]) pipeline.HandlerFunc[*domain.TaskInfo] {
		return func(ctx context.Context, object *domain.TaskInfo) error {
			// task discard begin
			exist := slices.Contains(dfb.discardMessageIds, object.MessageTemplateId)
			if exist {
				fmt.Printf("[%d]消息被丢弃", object.MessageTemplateId)
				return nil
			}
			err := next(ctx, object)
			// task discard end
			return err
		}
	}
}

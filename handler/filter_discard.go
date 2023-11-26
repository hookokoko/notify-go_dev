package handler

import (
	"context"
	"fmt"

	"github.com/ecodeclub/notify-go/common/domain"
	"github.com/ecodeclub/notify-go/common/pipeline"
)

//type DiscardFilter struct{}

//func (f *DiscardFilter) Process(ctx context.Context, taskInfo *domain.TaskInfo) error {
//	return nil
//}
//
//func (f *DiscardFilter) Before(ctx context.Context, taskInfo *domain.TaskInfo) error {
//	return nil
//}
//
//func (f *DiscardFilter) After(ctx context.Context, taskInfo *domain.TaskInfo) error {
//	return nil
//}

type DiscardFilterBuilder struct{}

func NewDiscardFilterBuilder() *DiscardFilterBuilder {
	return &DiscardFilterBuilder{}
}

func (dfb *DiscardFilterBuilder) Build() pipeline.Filter[*domain.TaskInfo] {
	return func(next pipeline.HandlerFunc[*domain.TaskInfo]) pipeline.HandlerFunc[*domain.TaskInfo] {
		return func(ctx context.Context, object *domain.TaskInfo) error {
			fmt.Println("task discard begin")
			err := next(ctx, object)
			fmt.Println("task discard end")
			return err
		}
	}
}

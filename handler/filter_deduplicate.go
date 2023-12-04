package handler

import (
	"context"

	"github.com/ecodeclub/notify-go/common/domain"
	"github.com/ecodeclub/notify-go/common/pipeline"
)

type DeduplicateFilterBuilder struct {
}

func NewDeduplicateFilterBuilder() *DeduplicateFilterBuilder {
	return &DeduplicateFilterBuilder{}
}

func (dfb *DeduplicateFilterBuilder) Build() pipeline.Filter[*domain.TaskInfo] {
	return func(next pipeline.HandlerFunc[*domain.TaskInfo]) pipeline.HandlerFunc[*domain.TaskInfo] {
		return func(ctx context.Context, object *domain.TaskInfo) error {
			// task deduplicate begin
			err := next(ctx, object)
			// task deduplicate end
			return err
		}
	}
}

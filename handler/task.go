package handler

import (
	"context"

	"github.com/ecodeclub/notify-go/common/domain"
	"github.com/ecodeclub/notify-go/common/pipeline"
	"github.com/panjf2000/ants/v2"
)

type TaskRunner func(taskInfo *domain.TaskInfo) bool

type TaskExecutor struct {
	pool *ants.Pool
}

func NewTaskExecutor(size int) *TaskExecutor {
	pool, err := ants.NewPool(size)
	if err != nil {
		return nil
	}
	return &TaskExecutor{
		pool: pool,
	}
}

func (te *TaskExecutor) run(info *domain.TaskInfo) error {
	ctx := context.TODO()
	//processController := pipeline.NewPipelineController[*domain.TaskInfo](
	//	"task",
	//	&DiscardFilter{},
	//	&SendMessageFilter{},
	//)
	//err := processController.Process(ctx, info)

	send := NewSendMsgHandlerBuilder().Build()
	processPipeline := pipeline.FilterChain[*domain.TaskInfo](
		NewDiscardFilterBuilder().Build(),
	).Then(send)

	err := processPipeline(ctx, info)
	return err
}

package experience

import (
	"context"
	"github.com/ecodeclub/notify-go/experience/handler"

	"github.com/ecodeclub/notify-go/experience/common/domain"
)

// SendService 不同发送方式的抽象
type SendService interface {
	Send(ctx context.Context, taskList []domain.TaskInfo) error
}

type DefaultSendImpl struct {
}

func NewDefaultSendImpl() *DefaultSendImpl {
	return &DefaultSendImpl{}
}

func (n *DefaultSendImpl) Send(ctx context.Context, taskList []domain.TaskInfo) error {
	service := handler.NewConsumeService()
	err := service.Consume2Send(taskList)
	return err
}

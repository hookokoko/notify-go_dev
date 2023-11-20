package experience

import (
	"context"

	"github.com/ecodeclub/notify-go/experience/common/domain"
)

// SendService 不同channel发送方式的抽象
type SendService interface {
	Send(ctx context.Context, taskList []domain.TaskInfo) error
}

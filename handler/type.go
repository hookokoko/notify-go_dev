package handler

import (
	"context"

	"github.com/ecodeclub/notify-go/common/domain"
)

const (
	SendSuccess AnchorState = iota
	SendFail
)

type Handler interface {
	Do(ctx context.Context, taskInfo domain.TaskInfo) bool
	//Recall()
}

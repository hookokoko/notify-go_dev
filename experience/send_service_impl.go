package experience

import (
	"context"
	"encoding/json"

	"github.com/ecodeclub/notify-go/experience/common/domain"
	"github.com/ecodeclub/notify-go/experience/handler"
	"github.com/ecodeclub/notify-go/experience/pkg/mq/kafka"
)

type KafkaSendImpl struct {
	client kafka.Kafka
}

func NewKafkaSendService(client kafka.Kafka) *KafkaSendImpl {
	return &KafkaSendImpl{
		client: client,
	}
}

func (k *KafkaSendImpl) Send(ctx context.Context, taskList []domain.TaskInfo) error {
	topic := "test-topic"
	tagId := ""
	body, err := json.Marshal(taskList)
	if err != nil {
		return err
	}
	err = k.client.Produce(ctx, topic, body, tagId)
	return err
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

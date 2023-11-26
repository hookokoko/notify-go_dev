package kafka

import (
	"context"
	"encoding/json"

	"github.com/ecodeclub/notify-go/common/domain"
	"github.com/ecodeclub/notify-go/handler"
	"github.com/ecodeclub/notify-go/pkg/log"
	"github.com/ecodeclub/notify-go/pkg/mq/kafka"
)

type MqConsumeServiceImpl struct {
	service handler.ConsumeService
	client  kafka.Kafka
}

func NewMqConsumeServiceImpl(service handler.ConsumeService, client kafka.Kafka) MqConsumeServiceImpl {
	return MqConsumeServiceImpl{
		service: service,
		client:  client,
	}
}

func (k *MqConsumeServiceImpl) Consume(ctx context.Context, topic string, groupId string) {
	defer k.stop()
	k.client.Consume(ctx, topic, groupId, func(msg []byte) {
		var taskInfoList []domain.TaskInfo
		err := json.Unmarshal(msg, &taskInfoList)
		if err != nil {
			panic(err)
		}
		err = k.service.Consume2Send(taskInfoList)
		if err != nil {
			log.Default().Error("消费", "error", err)
		}
	})
}

func (k *MqConsumeServiceImpl) stop() {
	k.client.Stop()
}

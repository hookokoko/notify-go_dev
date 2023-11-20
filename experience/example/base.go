package main

import (
	"context"
	"time"

	"github.com/ecodeclub/notify-go/experience"
	"github.com/ecodeclub/notify-go/experience/storage"
)

func main() {
	//go consumer.Consume()
	//<-time.Tick(3 * time.Second)
	params := []experience.MessageParam{
		{
			BizId:     "test biz",
			Receiver:  "648646891@qq.com",
			Variables: map[string]string{"name": "chenxx"},
		},
	}
	//d := experience.NewDelivery(params, 123,
	//	experience.WithSendService(experience.NewKafkaSendService(kafka.Kafka{Hosts: []string{"127.0.0.1:9092"}})))

	d := experience.NewNotification(params, 123, storage.MysqlDB())

	err := d.Send(context.Background())
	if err != nil {
		panic(err)
	}
	<-time.Tick(3 * time.Second)
}

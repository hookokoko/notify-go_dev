package consumer

import (
	"context"

	"github.com/ecodeclub/notify-go/experience/consumer/kafka"
	"github.com/ecodeclub/notify-go/experience/handler"
	client "github.com/ecodeclub/notify-go/experience/pkg/mq/kafka"
)

// TODO
// 这里需要做的事情：
// 1. 获取到消费的kafka topic、消费者组的设置
// 2. 其他消息队列的消费者启动参数构造

func Consume() {
	mqc := kafka.NewMqConsumeServiceImpl(handler.NewConsumeService(), client.NewKafka([]string{"127.0.0.1:9092"}))
	mqc.Consume(context.Background(), "test-topic", "EMAIL")
}

/*
需要明确的时，topic是共用一个，还是更细致的划分？
不管怎样，Consume要多启动几个，每个代表要处理的一类消息，每个启动的consume属于不同的消费者组。
实现起来，有两种方式：
1. 共用一个topic，根据tag id处理自己想要的消息，其余直接放弃处理
2. 分不同的topic，不用tag id
*/

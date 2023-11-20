package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/ecodeclub/notify-go/pkg/log"
)

type Kafka struct {
	Hosts []string
}

func NewKafka(hosts []string) Kafka {
	return Kafka{
		Hosts: hosts,
	}
}

func (k *Kafka) Produce(ctx context.Context, topic string, body []byte, tagId string) error {
	logger := log.FromContext(ctx)
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(k.Hosts, config)
	if err != nil {
		logger.Error("创建生产者出错", "err", err)
		return err
	}
	defer producer.AsyncClose()

	producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(tagId),
		Value: sarama.ByteEncoder(body),
	}

	select {
	case <-producer.Successes():
	case msgErr := <-producer.Errors():
		logger.Error("发送kafka消息失败", "msgErr", msgErr)
	}

	return nil
}

func (k *Kafka) Consume(ctx context.Context, topic string, groupId string, handler func(msg []byte)) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Consumer.Return.Errors = true
	cg, err := sarama.NewConsumerGroup(k.Hosts, groupId, saramaCfg)
	if err != nil {
		log.FromContext(ctx).Error("创建消费者出错", "err", err)
		return
	}
	defer func() { _ = cg.Close() }()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err = cg.Consume(ctx, []string{topic}, k.wrapSaramaHandler(handler)); err != nil {
				log.FromContext(ctx).Error("消费出错", "err", err)
			}
		}
	}
}

func (k *Kafka) wrapSaramaHandler(handler func(msg []byte)) sarama.ConsumerGroupHandler {
	return &consumerGroupHandler{
		handler: handler,
	}
}

type consumerGroupHandler struct {
	handler func(msg []byte)
}

func (c *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		msg := msg
		go func() {
			c.handler(msg.Value)
			session.MarkMessage(msg, "")
		}()
	}
	return nil
}

func (c *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error { return nil }

func (c *consumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error { return nil }

func (k *Kafka) Stop() {}

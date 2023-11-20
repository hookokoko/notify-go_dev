package kafka

import (
	"context"
	"testing"
)

func Test_Produce(t *testing.T) {
	k := &Kafka{
		Hosts: []string{"127.0.0.1:9092"},
	}
	if err := k.Produce(context.Background(), "test", []byte("test"), ""); err != nil {
		t.Error(err)
	}
}

func Test_Consume(t *testing.T) {
	k := &Kafka{
		Hosts: []string{"127.0.0.1:9092"},
	}
	k.Consume(context.Background(), "test", "test", func(msg []byte) {
		t.Log(string(msg))
	})
}

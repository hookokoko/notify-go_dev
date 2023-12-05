package handler

import (
	"context"
	"fmt"
	"github.com/ecodeclub/notify-go/common/domain"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRds(t *testing.T) {
	rds := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rds.FlushDB(context.Background())
	rds.Set(context.Background(), "key1", "val1", 0)
	rds.Set(context.Background(), "key2", "val2", 0)
	res := rds.MGet(context.TODO(), "key1", "key2", "key3").Val()
	fmt.Printf("%+v\n", res)
	fmt.Println(res[2] == nil)
}

func TestDuplicateService_duplicate(t *testing.T) {
	configs := map[string]deduplicateConfig{
		"deduplication_10": {
			Num: 5,
		},
		"deduplication_20": {
			Num:  1,
			Time: 300,
		},
	}

	rds := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	tests := []struct {
		name         string
		taskInfo     *domain.TaskInfo
		config       deduplicateConfig
		wantErr      error
		wantReceiver []string
	}{
		{
			name: "test1",
			taskInfo: &domain.TaskInfo{
				MessageTemplateId: 111111,
				SendChannel:       10,
				Receiver: []string{
					"123",
				},
			},
			config:  configs["deduplication_10"],
			wantErr: nil,
			wantReceiver: []string{
				"123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := DuplicateService{
				rds:    rds,
				config: configs,
			}
			err := ds.duplicate(context.TODO(), tt.taskInfo)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantReceiver, tt.taskInfo.Receiver)
		})
	}
}

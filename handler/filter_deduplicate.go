package handler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ecodeclub/notify-go/common/domain"
	"github.com/ecodeclub/notify-go/common/pipeline"
)

// 配置样例：{"deduplication_10":{"num":1,"time":300},"deduplication_20":{"num":5}}
// num = -1 不限制次数
// time = 0 不限制时间, 秒
type deduplicationStrategy struct {
	Num  int `json:"num"`
	Time int `json:"time"`
	DeduplicationKey
}

type DeduplicationKey interface {
	Key(receiver string, info domain.TaskInfo) string
}

// ContentDeduplicationKey 指定时间内相同的内容去重
// key ==> MD5(templateId + receiver + content)
type ContentDeduplicationKey struct{}

func (cdk ContentDeduplicationKey) Key(receiver string, info domain.TaskInfo) string {
	return fmt.Sprintf("deduplication_%s_%d_%s", receiver, info.MessageTemplateId, info.Content)
}

type DeduplicationFilterBuilder struct {
	service DuplicationService
}

// FrequencyDeduplicationKey 指定时间内超过指定频次进行拦截
// key ==> templateId + receiver + sendChannel
type FrequencyDeduplicationKey struct{}

func (fdk FrequencyDeduplicationKey) Key(receiver string, info domain.TaskInfo) string {
	return fmt.Sprintf("deduplication_%s_%d_%d", receiver, info.MessageTemplateId, info.SendChannel)
}

type DuplicationService struct {
	rds    *redis.Client
	config []deduplicationStrategy
}

func (ds DuplicationService) duplicate(ctx context.Context, taskInfo *domain.TaskInfo) (err error) {
	for _, cfg := range ds.config {
		filteredReceiver := make([]string, 0, len(taskInfo.Receiver))
		readyPutRedisReceiver := make(map[string]int, len(taskInfo.Receiver))
		keys := ds.deduplicationKeys(cfg, *taskInfo)
		rdsRes := ds.rds.MGet(ctx, keys...).Val()
		for idx, rcv := range taskInfo.Receiver {
			cnt := 0
			if rdsRes[idx] != nil {
				cnt, err = strconv.Atoi(rdsRes[idx].(string))
				if err != nil {
					return
				}
			}
			if cnt < cfg.Num {
				filteredReceiver = append(filteredReceiver, rcv)
				readyPutRedisReceiver[keys[idx]] = cnt + 1
			}
		}
		for k, cnt := range readyPutRedisReceiver {
			err = ds.rds.Set(ctx, k, cnt, time.Second*time.Duration(cfg.Time)).Err()
			if err != nil {
				return
			}
		}
		taskInfo.Receiver = filteredReceiver
	}
	return
}

func (ds DuplicationService) deduplicationKeys(s deduplicationStrategy, taskInfo domain.TaskInfo) []string {
	keys := make([]string, 0, len(taskInfo.Receiver))
	for _, rcv := range taskInfo.Receiver {
		keys = append(keys, s.Key(rcv, taskInfo))
	}
	return keys
}

func NewDeduplicationFilterBuilder(rds *redis.Client, config ...deduplicationStrategy) *DeduplicationFilterBuilder {
	return &DeduplicationFilterBuilder{
		service: DuplicationService{rds: rds, config: config},
	}
}

func (dfb *DeduplicationFilterBuilder) Build() pipeline.Filter[*domain.TaskInfo] {
	return func(next pipeline.HandlerFunc[*domain.TaskInfo]) pipeline.HandlerFunc[*domain.TaskInfo] {
		return func(ctx context.Context, object *domain.TaskInfo) error {
			// task Deduplication begin
			if err := dfb.service.duplicate(ctx, object); err != nil {
				return err
			}
			err := next(ctx, object)
			// task Deduplication end
			return err
		}
	}
}

package handler

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"

	"github.com/ecodeclub/notify-go/common/domain"
	"github.com/ecodeclub/notify-go/common/pipeline"
)

// 配置样例：{"deduplication_10":{"num":1,"time":300},"deduplication_20":{"num":5}}
// num = -1 不限制次数
// time = 0 不限制时间, 秒
type deduplicateConfig struct {
	Num  int `json:"num"`
	Time int `json:"time"`
	DeduplicateKey
}

type DeduplicateKey interface {
	Key(receiver string, info domain.TaskInfo) string
}

// ContentDeduplicateKey 指定时间内相同的内容去重
// key ==> MD5(templateId + receiver + content)
type ContentDeduplicateKey struct{}

func (cdk ContentDeduplicateKey) Key(receiver string, info domain.TaskInfo) string {
	return fmt.Sprintf("deduplication_%d", time.Now().Unix())
}

type DeduplicateFilterBuilder struct {
	service DuplicateService
}

// FrequencyDeduplicateKey 指定时间内超过指定频次进行拦截
// key ==> templateId + receiver + sendChannel
type FrequencyDeduplicateKey struct{}

func (fdk FrequencyDeduplicateKey) Key(receiver string, info domain.TaskInfo) string {
	return fmt.Sprintf("deduplication_%d", time.Now().Unix())
}

type DuplicateService struct {
	rds    *redis.Client
	config map[string]deduplicateConfig
}

func (ds DuplicateService) duplicate(ctx context.Context, taskInfo *domain.TaskInfo) (err error) {
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

func (ds DuplicateService) deduplicationKeys(cfg deduplicateConfig, taskInfo domain.TaskInfo) []string {
	keys := make([]string, 0, len(taskInfo.Receiver))
	for _, rcv := range taskInfo.Receiver {
		keys = append(keys, cfg.Key(rcv, taskInfo))
	}
	return keys
}

func NewDeduplicateFilterBuilder(rds *redis.Client) *DeduplicateFilterBuilder {
	config := make(map[string]deduplicateConfig)
	config["deduplication_10"] = deduplicateConfig{Num: 1, Time: 300}
	config["deduplication_20"] = deduplicateConfig{Num: 5}
	return &DeduplicateFilterBuilder{
		service: DuplicateService{rds: rds, config: config},
	}
}

func (dfb *DeduplicateFilterBuilder) Build() pipeline.Filter[*domain.TaskInfo] {
	return func(next pipeline.HandlerFunc[*domain.TaskInfo]) pipeline.HandlerFunc[*domain.TaskInfo] {
		return func(ctx context.Context, object *domain.TaskInfo) error {
			// task deduplicate begin
			if err := dfb.service.duplicate(ctx, object); err != nil {
				return err
			}
			err := next(ctx, object)
			// task deduplicate end
			return err
		}
	}
}

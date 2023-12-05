package handler

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/ecodeclub/notify-go/common/domain"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestDuplicateService_duplicate(t *testing.T) {
	mockRds := miniredis.RunT(t)
	rds := redis.NewClient(&redis.Options{Addr: mockRds.Addr()})
	tests := []struct {
		name        string
		taskInfo    *domain.TaskInfo
		strategyArr []deduplicationStrategy
		wantErr     error
	}{
		{
			name: "同一用户同渠道5min内超过3次发送",
			taskInfo: &domain.TaskInfo{
				MessageTemplateId: 111111,
				SendChannel:       10,
				Receiver:          []string{"123"},
			},
			strategyArr: []deduplicationStrategy{{Num: 3, Time: 300, DeduplicationKey: FrequencyDeduplicationKey{}}},
			wantErr:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := DuplicationService{
				rds:    rds,
				config: tt.strategyArr,
			}

			var err error
			taskInfo1 := *tt.taskInfo
			taskInfo2 := *tt.taskInfo
			taskInfo3 := *tt.taskInfo
			taskInfo4 := *tt.taskInfo

			// 第1次发送
			err = duplicationN(1, ds, &taskInfo1)
			assert.Equal(t, tt.wantErr, err)
			mockRds.CheckGet(t, tt.strategyArr[0].Key(taskInfo1.Receiver[0], taskInfo1), "1")
			assert.Equal(t, 1, len(taskInfo1.Receiver))

			// 刚好第n次发送
			err = duplicationN(tt.strategyArr[0].Num, ds, &taskInfo2)
			assert.Equal(t, tt.wantErr, err)
			mockRds.CheckGet(t, tt.strategyArr[0].Key(taskInfo1.Receiver[0], taskInfo1), strconv.Itoa(tt.strategyArr[0].Num))
			assert.Equal(t, 0, len(taskInfo2.Receiver))

			// 第n+100次发送
			err = duplicationN(tt.strategyArr[0].Num+100, ds, &taskInfo3)
			assert.Equal(t, tt.wantErr, err)
			mockRds.CheckGet(t, tt.strategyArr[0].Key(taskInfo1.Receiver[0], taskInfo1), strconv.Itoa(tt.strategyArr[0].Num))
			assert.Equal(t, 0, len(taskInfo3.Receiver))

			// 超过过期时间发送
			mockRds.FastForward(time.Duration(tt.strategyArr[0].Time) * time.Second)
			err = duplicationN(tt.strategyArr[0].Num-1, ds, &taskInfo4)
			assert.Equal(t, tt.wantErr, err)
			mockRds.CheckGet(t, tt.strategyArr[0].Key(taskInfo1.Receiver[0], taskInfo1), strconv.Itoa(tt.strategyArr[0].Num-1))
			assert.Equal(t, 1, len(taskInfo4.Receiver))
		})
	}
}

func duplicationN(n int, ds DuplicationService, taskInfo *domain.TaskInfo) (err error) {
	for i := 0; i < n; i++ {
		err = ds.duplicate(context.TODO(), taskInfo)
		if err != nil {
			return
		}
	}
	return
}

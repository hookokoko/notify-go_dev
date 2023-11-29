package handler

import (
	"log"
	"testing"

	"github.com/ecodeclub/notify-go/common/domain/sms"
)

func TestLoadBalance(t *testing.T) {
	// 创建测试用例
	testCases := []struct {
		name   string
		config []sms.MessageTypeSmsConfig
	}{
		{
			name:   "EmptyConfigs",
			config: []sms.MessageTypeSmsConfig{},
		},
		{
			name: "SingleConfig",
			config: []sms.MessageTypeSmsConfig{
				{Weights: 100},
			},
		},
		{
			name: "MultipleConfigs",
			config: []sms.MessageTypeSmsConfig{
				{Weights: 30},
				{Weights: 50},
				{Weights: 20},
			},
		},
	}

	// 遍历测试用例
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			handler := &SMSHandler{}
			res := make(map[int]int)
			for i := 0; i < 1000; i++ {
				res[handler.loadBalance(testCase.config).Weights]++
			}
			log.Println(res)
		})
	}
}

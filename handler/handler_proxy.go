// Copyright 2021 ecodeclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"github.com/ecodeclub/notify-go/common/domain"
	"github.com/ecodeclub/notify-go/common/enum/channel_type"
	"github.com/ecodeclub/notify-go/pkg/log"
)

type ConsumeService struct {
	executors map[string]*TaskExecutor
}

func NewConsumeService() ConsumeService {
	service := ConsumeService{
		executors: make(map[string]*TaskExecutor),
	}
	for _, channel := range channel_type.Values() {
		executor := NewTaskExecutor(10)
		if executor == nil {
			log.Default().Error("创建协程池失败", "channel", channel)
		}
		service.executors[channel] = executor
	}
	return service
}

func (c ConsumeService) Consume2Send(taskInfoList []domain.TaskInfo) error {
	groupId := channel_type.ChannelType(taskInfoList[0].SendChannel).String()
	for _, info := range taskInfoList {
		err := c.executors[groupId].run(&info)
		if err != nil {
			log.Default().Error("我执行失败了...", "body", info)
		}
	}
	return nil
}

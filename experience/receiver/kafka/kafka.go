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

package kafka

import (
	"context"
	"encoding/json"
	"github.com/ecodeclub/notify-go/experience"
	"github.com/ecodeclub/notify-go/experience/receiver"
)

type SendMqServiceImpl struct {
	service receiver.ConsumeService
}

func (k *SendMqServiceImpl) Send(ctx context.Context, topic string, body []byte, tagId string) error {
	//TODO implement me
	panic("implement me")
}

func (k *SendMqServiceImpl) Consumer(msg []byte) {
	var taskInfoList []experience.TaskInfo
	err := json.Unmarshal(msg, &taskInfoList)
	if err != nil {
		panic(err)
	}
	k.service.Consume2Send(taskInfoList)
}

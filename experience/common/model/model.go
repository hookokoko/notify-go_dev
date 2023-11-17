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

package model

import (
	"github.com/ecodeclub/notify-go/experience/common/enum"
	"reflect"
)

// GetChanelModelClassByCode 通过code获取class
func GetChanelModelClassByCode(code int) reflect.Type {
	for _, channelType := range enum.ChannelTypes {
		if code == channelType.Code {
			return channelType.ContentModelClass
		}
	}
	return nil
}

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

package enum

import (
	"github.com/ecodeclub/notify-go/experience/common/model"
	"reflect"
)

// ChannelType 发送渠道类型枚举
type ChannelType struct {
	Code              int
	Description       string
	ContentModelClass reflect.Type
	CodeEn            string
	AccessTokenPrefix string
	AccessTokenExpire int64
}

var ChannelTypes = []ChannelType{
	{
		Code:              20,
		Description:       "push(通知栏)",
		ContentModelClass: reflect.TypeOf(model.PushContentModel{}),
		CodeEn:            "push",
		AccessTokenPrefix: "ge_tui_access_token_",
		AccessTokenExpire: 3600 * 24,
	},
	{
		Code:              30,
		Description:       "sms(短信)",
		ContentModelClass: reflect.TypeOf(model.SmsContentModel{}),
		CodeEn:            "sms",
	},
	{
		Code:              40,
		Description:       "email(邮件)",
		ContentModelClass: reflect.TypeOf(model.EmailContentModel{}),
		CodeEn:            "email",
	},
}

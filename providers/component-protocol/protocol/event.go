// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protocol

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// eventConvert .
// 前端触发的事件转换，如果是组件自身的事件，则透传；
// 否则, (1) 组件名为空，界面刷新：InitializeOperation
//
//	(2) 通过协议定义的Rending触发的事件：RenderingOperation
func eventConvert(receiver string, event cptype.ComponentEvent) cptype.ComponentEvent {
	if receiver == event.Component {
		return event
	} else if event.Component != "" {
		return cptype.ComponentEvent{Operation: cptype.RenderingOperation}
	} else {
		return cptype.ComponentEvent{Operation: cptype.InitializeOperation}
	}
}

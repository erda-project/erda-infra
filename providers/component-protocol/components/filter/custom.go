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

package filter

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// CustomData is kanban's custom data handler.
type CustomData interface {
	CustomDataPtr() interface{}
	EncodeFromCustomData(customDataPtr interface{}, stdDataPtr *Data)
	DecodeToCustomData(stdDataPtr *Data, customDataPtr interface{})
}

// CustomState is kanban's custom state handler.
type CustomState interface {
	CustomStatePtr() interface{}
	EncodeFromCustomState(customStatePtr interface{}, stdStatePtr *State)
	DecodeToCustomState(stdStatePtr *State, customStatePtr interface{})
}

// CustomInParams is kanban's custom inParams handler.
type CustomInParams interface {
	CustomInParamsPtr() interface{}
	EncodeFromCustomInParams(customInParamsPtr interface{}, stdInParamsPtr *cptype.ExtraMap)
	DecodeToCustomInParams(stdInParamsPtr *cptype.ExtraMap, customInParamsPtr interface{})
}

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

package kanban

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

// CustomState .
type CustomState struct {
	Name string `json:"name,omitempty"`
}

// CustomInParams .
type CustomInParams struct {
	ProjectID uint64 `json:"projectID,omitempty"`
}

// CustomStatePtr .
func (p *provider) CustomStatePtr() interface{} {
	if p.StatePtr == nil {
		p.StatePtr = &CustomState{}
	}
	return p.StatePtr
}

// EncodeFromCustomState .
func (p *provider) EncodeFromCustomState(customStatePtr interface{}, stdStatePtr *cptype.ExtraMap) {
	cputil.MustObjJSONTransfer(customStatePtr, stdStatePtr)
}

// DecodeToCustomState .
func (p *provider) DecodeToCustomState(stdStatePtr *cptype.ExtraMap, customStatePtr interface{}) {
	cputil.MustObjJSONTransfer(stdStatePtr, customStatePtr)
}

// CustomInParamsPtr .
func (p *provider) CustomInParamsPtr() interface{} {
	if p.InParamsPtr == nil {
		p.InParamsPtr = &CustomInParams{}
	}
	return p.InParamsPtr
}

// EncodeFromCustomInParams .
func (p *provider) EncodeFromCustomInParams(customInParamsPtr interface{}, stdInParamsPtr *cptype.ExtraMap) {
	cputil.MustObjJSONTransfer(customInParamsPtr, stdInParamsPtr)
}

// DecodeToCustomInParams .
func (p *provider) DecodeToCustomInParams(stdInParamsPtr *cptype.ExtraMap, customInParamsPtr interface{}) {
	cputil.MustObjJSONTransfer(stdInParamsPtr, customInParamsPtr)
}

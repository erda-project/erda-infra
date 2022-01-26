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

package impl

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/components/complexgraph"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

// DecodeData .
func (d *DefaultComplexGraph) DecodeData(compData cptype.ComponentData, modelDataPtr interface{}) {
	cputil.MustObjJSONTransfer(compData, modelDataPtr.(*complexgraph.Data))
	if custom, ok := d.Impl.(complexgraph.CustomData); ok {
		custom.DecodeToCustomData(d.StdDataPtr, custom.CustomDataPtr())
	}
}

// DecodeState .
func (d *DefaultComplexGraph) DecodeState(compState cptype.ComponentState, modelStatePtr interface{}) {
	cputil.MustObjJSONTransfer(compState, modelStatePtr.(*cptype.ExtraMap))
	if custom, ok := d.Impl.(complexgraph.CustomState); ok {
		custom.DecodeToCustomState(d.StdStatePtr, custom.CustomStatePtr())
	}
}

// DecodeInParams .
func (d *DefaultComplexGraph) DecodeInParams(compInParams cptype.InParams, modelInParamsPtr interface{}) {
	cputil.MustObjJSONTransfer(compInParams, modelInParamsPtr.(*cptype.ExtraMap))
	if custom, ok := d.Impl.(complexgraph.CustomInParams); ok {
		custom.DecodeToCustomInParams(d.StdInParamsPtr, custom.CustomInParamsPtr())
	}
}

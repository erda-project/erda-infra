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
	"github.com/erda-project/erda-infra/providers/component-protocol/components/kanban"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

// EncodeData .
func (d *DefaultKanban) EncodeData(srcStructPtr interface{}, dstRawPtr *cptype.ComponentData) {
	// custom -> struct
	if custom, ok := d.Impl.(kanban.CustomData); ok {
		custom.EncodeFromCustomData(custom.CustomDataPtr(), srcStructPtr.(*kanban.Data))
	}
	// struct -> raw
	cputil.MustObjJSONTransfer(srcStructPtr, dstRawPtr)
}

// EncodeState .
func (d *DefaultKanban) EncodeState(srcStructPtr interface{}, dstRawPtr *cptype.ComponentState) {
	if custom, ok := d.Impl.(kanban.CustomState); ok {
		custom.EncodeFromCustomState(custom.CustomStatePtr(), srcStructPtr.(*cptype.ExtraMap))
	}
	cputil.MustObjJSONTransfer(srcStructPtr, dstRawPtr)
}

// EncodeInParams .
func (d *DefaultKanban) EncodeInParams(srcStructPtr interface{}, dstRawPtr *cptype.InParams) {
	if custom, ok := d.Impl.(kanban.CustomInParams); ok {
		custom.EncodeFromCustomInParams(custom.CustomInParamsPtr(), srcStructPtr.(*cptype.ExtraMap))
	}
	cputil.MustObjJSONTransfer(srcStructPtr, dstRawPtr)
}

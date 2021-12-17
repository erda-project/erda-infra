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

// DefaultKanban .
type DefaultKanban struct {
	Impl kanban.IKanban

	*StdStructuredPtr
}

// StdStructuredPtr .
type StdStructuredPtr struct {
	StdInParamsPtr *cptype.ExtraMap
	StdDataPtr     *kanban.Data
	StdStatePtr    *cptype.ExtraMap
}

// DataPtr .
func (s *StdStructuredPtr) DataPtr() interface{} { return s.StdDataPtr }

// StatePtr .
func (s *StdStructuredPtr) StatePtr() interface{} { return s.StdStatePtr }

// InParamsPtr .
func (s *StdStructuredPtr) InParamsPtr() interface{} { return s.StdInParamsPtr }

// RegisterCompStdOps .
func (d *DefaultKanban) RegisterCompStdOps() (opFuncs map[cptype.OperationKey]cptype.OperationFunc) {
	return map[cptype.OperationKey]cptype.OperationFunc{
		// kanban-level
		kanban.OpBoardCreate{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterBoardCreateOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &kanban.OpBoardCreate{}).(*kanban.OpBoardCreate))(sdk)
		},

		// board-level
		kanban.OpBoardLoadMore{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterBoardLoadMoreOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &kanban.OpBoardLoadMore{}).(*kanban.OpBoardLoadMore))(sdk)
		},
		kanban.OpBoardUpdate{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterBoardUpdateOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &kanban.OpBoardUpdate{}).(*kanban.OpBoardUpdate))(sdk)
		},
		kanban.OpBoardDelete{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterBoardDeleteOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &kanban.OpBoardDelete{}).(*kanban.OpBoardDelete))(sdk)
		},

		// card-level
		kanban.OpCardMoveTo{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterCardMoveToOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &kanban.OpCardMoveTo{}).(*kanban.OpCardMoveTo))(sdk)
		},
	}
}

// Initialize .
func (d *DefaultKanban) Initialize(sdk *cptype.SDK) {}

// Finalize .
func (d *DefaultKanban) Finalize(sdk *cptype.SDK) {}

// SkipOp providers default impl for user.
func (d *DefaultKanban) SkipOp(sdk *cptype.SDK) bool { return !d.Impl.Visible(sdk) }

// Visible .
func (d *DefaultKanban) Visible(sdk *cptype.SDK) bool { return true }

// BeforeHandleOp providers default impl for user.
func (d *DefaultKanban) BeforeHandleOp(sdk *cptype.SDK) {}

// AfterHandleOp providers default impl for user.
func (d *DefaultKanban) AfterHandleOp(sdk *cptype.SDK) {}

// StdStructuredPtrCreator .
func (d *DefaultKanban) StdStructuredPtrCreator() func() cptype.IStdStructuredPtr {
	return func() cptype.IStdStructuredPtr {
		return &StdStructuredPtr{
			StdInParamsPtr: &cptype.ExtraMap{},
			StdStatePtr:    &cptype.ExtraMap{},
			StdDataPtr:     &kanban.Data{},
		}
	}
}

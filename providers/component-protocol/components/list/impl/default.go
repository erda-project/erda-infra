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
	"github.com/erda-project/erda-infra/providers/component-protocol/components/list"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

// DefaultList .
type DefaultList struct {
	Impl list.IList

	*StdStructuredPtr
}

// StdStructuredPtr .
type StdStructuredPtr struct {
	StdInParamsPtr *cptype.ExtraMap
	StdDataPtr     *list.Data
	StdStatePtr    *cptype.ExtraMap
}

// DataPtr .
func (s *StdStructuredPtr) DataPtr() interface{} { return s.StdDataPtr }

// StatePtr .
func (s *StdStructuredPtr) StatePtr() interface{} { return s.StdStatePtr }

// InParamsPtr .
func (s *StdStructuredPtr) InParamsPtr() interface{} { return s.StdInParamsPtr }

// RegisterCompStdOps .
func (d *DefaultList) RegisterCompStdOps() (opFuncs map[cptype.OperationKey]cptype.OperationFunc) {
	return map[cptype.OperationKey]cptype.OperationFunc{
		// list level
		list.OpChangePage{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterChangePage(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &list.OpChangePage{}).(*list.OpChangePage))
		},

		// item level
		list.OpItemStar{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterItemStarOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &list.OpItemStar{}).(*list.OpItemStar))
		},
		list.OpItemClickGoto{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterItemClickGotoOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &list.OpItemClickGoto{}).(*list.OpItemClickGoto))
		},
		list.OpItemClick{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterItemClickOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &list.OpItemClick{}).(*list.OpItemClick))
		},
	}
}

// RegisterCompNonStdOps .
func (d *DefaultList) RegisterCompNonStdOps() (opFuncs map[cptype.OperationKey]cptype.OperationFunc) {
	return nil
}

// Initialize .
func (d *DefaultList) Initialize(sdk *cptype.SDK) {}

// Finalize .
func (d *DefaultList) Finalize(sdk *cptype.SDK) {}

// SkipOp providers default impl for user.
func (d *DefaultList) SkipOp(sdk *cptype.SDK) bool { return !d.Impl.Visible(sdk) }

// Visible .
func (d *DefaultList) Visible(sdk *cptype.SDK) bool { return true }

// BeforeHandleOp providers default impl for user.
func (d *DefaultList) BeforeHandleOp(sdk *cptype.SDK) {}

// AfterHandleOp providers default impl for user.
func (d *DefaultList) AfterHandleOp(sdk *cptype.SDK) {}

// StdStructuredPtrCreator .
func (d *DefaultList) StdStructuredPtrCreator() func() cptype.IStdStructuredPtr {
	return func() cptype.IStdStructuredPtr {
		return &StdStructuredPtr{
			StdInParamsPtr: &cptype.ExtraMap{},
			StdStatePtr:    &cptype.ExtraMap{},
			StdDataPtr:     &list.Data{},
		}
	}
}

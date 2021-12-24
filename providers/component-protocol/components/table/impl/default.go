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
	"github.com/erda-project/erda-infra/providers/component-protocol/components/table"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

// DefaultTable .
type DefaultTable struct {
	Impl table.ITable

	*StdStructuredPtr
}

// StdStructuredPtr .
type StdStructuredPtr struct {
	StdInParamsPtr *cptype.ExtraMap
	StdDataPtr     *table.Data
	StdStatePtr    *cptype.ExtraMap
}

// DataPtr .
func (s *StdStructuredPtr) DataPtr() interface{} { return s.StdDataPtr }

// StatePtr .
func (s *StdStructuredPtr) StatePtr() interface{} { return s.StdStatePtr }

// InParamsPtr .
func (s *StdStructuredPtr) InParamsPtr() interface{} { return s.StdInParamsPtr }

// RegisterCompStdOps .
func (d *DefaultTable) RegisterCompStdOps() (opFuncs map[cptype.OperationKey]cptype.OperationFunc) {
	return map[cptype.OperationKey]cptype.OperationFunc{
		// table-level
		table.OpTableChangePage{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterTableChangePageOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &table.OpTableChangePage{}).(*table.OpTableChangePage))(sdk)
		},
		table.OpTableChangeSort{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterTableSortOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &table.OpTableChangeSort{}).(*table.OpTableChangeSort))(sdk)
		},
		table.OpBatchRowsHandle{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterBatchRowsHandleOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &table.OpBatchRowsHandle{}).(*table.OpBatchRowsHandle))(sdk)
		},
		// row-level
		table.OpRowSelect{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterRowSelectOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &table.OpRowSelect{}).(*table.OpRowSelect))(sdk)
		},
		table.OpRowAdd{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterRowAddOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &table.OpRowAdd{}).(*table.OpRowAdd))(sdk)
		},
		table.OpRowEdit{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterRowEditOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &table.OpRowEdit{}).(*table.OpRowEdit))(sdk)
		},
		table.OpRowDelete{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterRowDeleteOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &table.OpRowDelete{}).(*table.OpRowDelete))(sdk)
		},
		// cell-level
	}
}

// RegisterCompNonStdOps .
func (d *DefaultTable) RegisterCompNonStdOps() (opFuncs map[cptype.OperationKey]cptype.OperationFunc) {
	return nil
}

// Initialize .
func (d *DefaultTable) Initialize(sdk *cptype.SDK) {}

// Finalize .
func (d *DefaultTable) Finalize(sdk *cptype.SDK) {}

// SkipOp providers default impl for user.
func (d *DefaultTable) SkipOp(sdk *cptype.SDK) bool { return !d.Impl.Visible(sdk) }

// Visible .
func (d *DefaultTable) Visible(sdk *cptype.SDK) bool { return true }

// BeforeHandleOp providers default impl for user.
func (d *DefaultTable) BeforeHandleOp(sdk *cptype.SDK) {}

// AfterHandleOp providers default impl for user.
func (d *DefaultTable) AfterHandleOp(sdk *cptype.SDK) {}

// StdStructuredPtrCreator .
func (d *DefaultTable) StdStructuredPtrCreator() func() cptype.IStdStructuredPtr {
	return func() cptype.IStdStructuredPtr {
		return &StdStructuredPtr{
			StdInParamsPtr: &cptype.ExtraMap{},
			StdStatePtr:    &cptype.ExtraMap{},
			StdDataPtr:     &table.Data{},
		}
	}
}

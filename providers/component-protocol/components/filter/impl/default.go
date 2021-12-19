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
	"github.com/erda-project/erda-infra/providers/component-protocol/components/filter"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

// DefaultFilter .
type DefaultFilter struct {
	Impl filter.IFilter

	*StdStructuredPtr
}

// StdStructuredPtr .
type StdStructuredPtr struct {
	StdInParamsPtr *cptype.ExtraMap
	StdDataPtr     *filter.Data
	StdStatePtr    *cptype.ExtraMap
}

// DataPtr .
func (s *StdStructuredPtr) DataPtr() interface{} { return s.StdDataPtr }

// StatePtr .
func (s *StdStructuredPtr) StatePtr() interface{} { return s.StdStatePtr }

// InParamsPtr .
func (s *StdStructuredPtr) InParamsPtr() interface{} { return s.StdInParamsPtr }

// RegisterCompStdOps .
func (d *DefaultFilter) RegisterCompStdOps() (opFuncs map[cptype.OperationKey]cptype.OperationFunc) {
	return map[cptype.OperationKey]cptype.OperationFunc{
		filter.OpFilter{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterFilterOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &filter.OpFilter{}).(*filter.OpFilter))(sdk)
		},
		filter.OpFilterItemSave{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterFilterItemSaveOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &filter.OpFilterItemSave{}).(*filter.OpFilterItemSave))(sdk)
		},
		filter.OpFilterItemDelete{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterFilterItemDeleteOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &filter.OpFilterItemDelete{}).(*filter.OpFilterItemDelete))(sdk)
		},
	}
}

// Initialize .
func (d *DefaultFilter) Initialize(sdk *cptype.SDK) {}

// Finalize .
func (d *DefaultFilter) Finalize(sdk *cptype.SDK) {}

// SkipOp providers default impl for user.
func (d *DefaultFilter) SkipOp(sdk *cptype.SDK) bool { return !d.Impl.Visible(sdk) }

// BeforeHandleOp providers default impl for user.
func (d *DefaultFilter) BeforeHandleOp(sdk *cptype.SDK) {}

// AfterHandleOp providers default impl for user.
func (d *DefaultFilter) AfterHandleOp(sdk *cptype.SDK) {}

// StdStructuredPtrCreator .
func (d *DefaultFilter) StdStructuredPtrCreator() func() cptype.IStdStructuredPtr {
	return func() cptype.IStdStructuredPtr {
		return &StdStructuredPtr{
			StdInParamsPtr: &cptype.ExtraMap{},
			StdStatePtr:    &cptype.ExtraMap{},
			StdDataPtr:     &filter.Data{},
		}
	}
}

// Visible .
func (d *DefaultFilter) Visible(sdk *cptype.SDK) bool { return true }

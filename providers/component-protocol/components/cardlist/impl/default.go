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
	"github.com/erda-project/erda-infra/providers/component-protocol/components/cardlist"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

// DefaultCardList .
type DefaultCardList struct {
	Impl cardlist.ICardList

	*StdStructuredPtr
}

// StdStructuredPtr .
type StdStructuredPtr struct {
	StdInParamsPtr *cptype.ExtraMap
	StdDataPtr     *cardlist.Data
	StdStatePtr    *cptype.ExtraMap
}

// DataPtr .
func (s *StdStructuredPtr) DataPtr() interface{} { return s.StdDataPtr }

// StatePtr .
func (s *StdStructuredPtr) StatePtr() interface{} { return s.StdStatePtr }

// InParamsPtr .
func (s *StdStructuredPtr) InParamsPtr() interface{} { return s.StdInParamsPtr }

// RegisterCompStdOps .
func (d *DefaultCardList) RegisterCompStdOps() (opFuncs map[cptype.OperationKey]cptype.OperationFunc) {
	return map[cptype.OperationKey]cptype.OperationFunc{
		cardlist.OpCardListGoto{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterCardListOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &cardlist.OpCardListGoto{}).(*cardlist.OpCardListGoto))(sdk)
		},
		cardlist.OpCardListIconGoto{}.OpKey(): func(sdk *cptype.SDK) {
			d.Impl.RegisterCardListIconOp(*cputil.MustObjJSONTransfer(sdk.Event.OperationData, &cardlist.OpCardListIconGoto{}).(*cardlist.OpCardListIconGoto))(sdk)
		},
	}
}

// RegisterCompNonStdOps .
func (d *DefaultCardList) RegisterCompNonStdOps() (opFuncs map[cptype.OperationKey]cptype.OperationFunc) {
	return nil
}

// Initialize .
func (d *DefaultCardList) Initialize(sdk *cptype.SDK) {}

// Finalize .
func (d *DefaultCardList) Finalize(sdk *cptype.SDK) {}

// SkipOp providers default impl for user.
func (d *DefaultCardList) SkipOp(sdk *cptype.SDK) bool { return !d.Impl.Visible(sdk) }

// BeforeHandleOp providers default impl for user.
func (d *DefaultCardList) BeforeHandleOp(sdk *cptype.SDK) {}

// AfterHandleOp providers default impl for user.
func (d *DefaultCardList) AfterHandleOp(sdk *cptype.SDK) {}

// StdStructuredPtrCreator .
func (d *DefaultCardList) StdStructuredPtrCreator() func() cptype.IStdStructuredPtr {
	return func() cptype.IStdStructuredPtr {
		return &StdStructuredPtr{
			StdInParamsPtr: &cptype.ExtraMap{},
			StdStatePtr:    &cptype.ExtraMap{},
			StdDataPtr:     &cardlist.Data{},
		}
	}
}

// Visible .
func (d *DefaultCardList) Visible(sdk *cptype.SDK) bool { return true }

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
	"github.com/erda-project/erda-infra/providers/component-protocol/components/indicatorcard"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// DefaultIndicatorCard default indicator card component
type DefaultIndicatorCard struct {
	Impl indicatorcard.IIndicatorCard
	*StdStructuredPtr
}

// Initialize .
func (d *DefaultIndicatorCard) Initialize(sdk *cptype.SDK) {
}

// Visible .
func (d *DefaultIndicatorCard) Visible(sdk *cptype.SDK) bool {
	return true
}

// StdStructuredPtr .
type StdStructuredPtr struct {
	StdInParamsPtr *cptype.ExtraMap
	StdDataPtr     *indicatorcard.Data
	StdStatePtr    *cptype.ExtraMap
}

// DataPtr .
func (s *StdStructuredPtr) DataPtr() interface{} { return s.StdDataPtr }

// StatePtr .
func (s *StdStructuredPtr) StatePtr() interface{} { return s.StdStatePtr }

// InParamsPtr .
func (s *StdStructuredPtr) InParamsPtr() interface{} { return s.StdInParamsPtr }

// RegisterCompStdOps .
func (d *DefaultIndicatorCard) RegisterCompStdOps() (opFuncs map[cptype.OperationKey]cptype.OperationFunc) {
	return map[cptype.OperationKey]cptype.OperationFunc{}
}

// RegisterCompNonStdOps .
func (d *DefaultIndicatorCard) RegisterCompNonStdOps() (opFuncs map[cptype.OperationKey]cptype.OperationFunc) {
	return nil
}

// Finalize .
func (d *DefaultIndicatorCard) Finalize(sdk *cptype.SDK) {}

// SkipOp providers default impl for user.
func (d *DefaultIndicatorCard) SkipOp(sdk *cptype.SDK) bool { return !d.Impl.Visible(sdk) }

// BeforeHandleOp providers default impl for user.
func (d *DefaultIndicatorCard) BeforeHandleOp(sdk *cptype.SDK) {}

// AfterHandleOp providers default impl for user.
func (d *DefaultIndicatorCard) AfterHandleOp(sdk *cptype.SDK) {}

// StdStructuredPtrCreator .
func (d *DefaultIndicatorCard) StdStructuredPtrCreator() func() cptype.IStdStructuredPtr {
	return func() cptype.IStdStructuredPtr {
		return &StdStructuredPtr{
			StdInParamsPtr: &cptype.ExtraMap{},
			StdStatePtr:    &cptype.ExtraMap{},
			StdDataPtr:     &indicatorcard.Data{},
		}
	}
}

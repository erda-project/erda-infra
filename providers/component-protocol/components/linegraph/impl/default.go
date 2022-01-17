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
	"github.com/erda-project/erda-infra/providers/component-protocol/components/defaults"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/linegraph"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// DefaultLineGraph default line graph component
type DefaultLineGraph struct {
	defaults.DefaultImpl
	Impl linegraph.ILineGraph
	*StdStructuredPtr
}

// Initialize .
func (d *DefaultLineGraph) Initialize(sdk *cptype.SDK) {
}

// Visible .
func (d *DefaultLineGraph) Visible(sdk *cptype.SDK) bool {
	return true
}

// StdStructuredPtr .
type StdStructuredPtr struct {
	StdInParamsPtr *cptype.ExtraMap
	StdDataPtr     *linegraph.Data
	StdStatePtr    *cptype.ExtraMap
}

// DataPtr .
func (s *StdStructuredPtr) DataPtr() interface{} { return s.StdDataPtr }

// StatePtr .
func (s *StdStructuredPtr) StatePtr() interface{} { return s.StdStatePtr }

// InParamsPtr .
func (s *StdStructuredPtr) InParamsPtr() interface{} { return s.StdInParamsPtr }

// RegisterCompStdOps .
func (d *DefaultLineGraph) RegisterCompStdOps() (opFuncs map[cptype.OperationKey]cptype.OperationFunc) {
	return map[cptype.OperationKey]cptype.OperationFunc{}
}

// RegisterCompNonStdOps .
func (d *DefaultLineGraph) RegisterCompNonStdOps() (opFuncs map[cptype.OperationKey]cptype.OperationFunc) {
	return nil
}

// Finalize .
func (d *DefaultLineGraph) Finalize(sdk *cptype.SDK) {}

// SkipOp providers default impl for user.
func (d *DefaultLineGraph) SkipOp(sdk *cptype.SDK) bool { return !d.Impl.Visible(sdk) }

// BeforeHandleOp providers default impl for user.
func (d *DefaultLineGraph) BeforeHandleOp(sdk *cptype.SDK) {}

// AfterHandleOp providers default impl for user.
func (d *DefaultLineGraph) AfterHandleOp(sdk *cptype.SDK) {}

// StdStructuredPtrCreator .
func (d *DefaultLineGraph) StdStructuredPtrCreator() func() cptype.IStdStructuredPtr {
	return func() cptype.IStdStructuredPtr {
		return &StdStructuredPtr{
			StdInParamsPtr: &cptype.ExtraMap{},
			StdStatePtr:    &cptype.ExtraMap{},
			StdDataPtr:     &linegraph.Data{},
		}
	}
}

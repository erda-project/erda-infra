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

package defaults

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// FieldActualImplRef .
const (
	FieldActualImplRef = "ActualImplRef"
)

// DefaultImpl .
type DefaultImpl struct {
	// ActualImplRef inject by framework according to field FieldActualImplRef
	ActualImplRef cptype.IComponent
}

// RegisterInitializeOp .
func (d *DefaultImpl) RegisterInitializeOp() (opFunc cptype.OperationFunc) {
	return nil
}

// RegisterRenderingOp .
func (d *DefaultImpl) RegisterRenderingOp() (opFunc cptype.OperationFunc) {
	return d.ActualImplRef.RegisterInitializeOp()
}

// RegisterInitializeOpV2 .
func (d *DefaultImpl) RegisterInitializeOpV2() (opFunc cptype.EnhancedOperationFunc) {
	return func(sdk *cptype.SDK) cptype.IStdStructuredPtr {
		d.ActualImplRef.RegisterInitializeOp()(sdk)
		return nil
	}
}

// RegisterRenderingOpV2 .
func (d *DefaultImpl) RegisterRenderingOpV2() (opFunc cptype.EnhancedOperationFunc) {
	return func(sdk *cptype.SDK) cptype.IStdStructuredPtr {
		d.ActualImplRef.RegisterInitializeOp()(sdk)
		return nil
	}
}

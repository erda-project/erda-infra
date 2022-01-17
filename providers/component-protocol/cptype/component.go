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

package cptype

const (
	// FieldImplForInject used to inject `Impl` field by framework.
	FieldImplForInject = "Impl"
)

// IComponent combines IDefaultComponent and BaseOperationRegister.
type IComponent interface {
	IDefaultComponent
	CompBaseOperationRegister
}

// CompBaseOperationRegister includes Initialize & Rendering Op.
type CompBaseOperationRegister interface {
	RegisterInitializeOp() (opFunc OperationFunc)
	RegisterRenderingOp() (opFunc OperationFunc)
	RegisterInitializeOpV2() (opFunc EnhancedOperationFunc)
	RegisterRenderingOpV2() (opFunc EnhancedOperationFunc)
}

// CompOperationRegister includes a component's all operations.
type CompOperationRegister interface {
	CompBaseOperationRegister
	CompNonBaseOperationRegister
}

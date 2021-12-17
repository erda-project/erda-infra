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

// IDefaultComponent std component's default impl should implement this interface.
type IDefaultComponent interface {
	CompStdStructuredPtrCreator
	CompStdOperationRegister
	CompFrameworkSteper
	CompEncoder
	CompDecoder
}

// CompStdStructuredPtrCreator return creator of StdStructuredPtr.
type CompStdStructuredPtrCreator interface {
	StdStructuredPtrCreator() func() IStdStructuredPtr
}

// IStdStructuredPtr represents std structured pointer type.
type IStdStructuredPtr interface {
	DataPtr() interface{}
	StatePtr() interface{}
	InParamsPtr() interface{}
}

// CompStdOperationRegister register a component's all custom operations to standard cptype.OperationFunc,
// and then used by framework.
type CompStdOperationRegister interface {
	RegisterCompStdOps() (opFuncs map[OperationKey]OperationFunc)
}

// CompFrameworkSteper represents all component steps played in framework.
type CompFrameworkSteper interface {
	Initialize(sdk *SDK)
	Finalize(sdk *SDK)
	SkipOp(sdk *SDK) bool
	Visible(sdk *SDK) bool
	BeforeHandleOp(sdk *SDK)
	AfterHandleOp(sdk *SDK)
}

// CompEncoder is a protocol-level interface, convert structured-struct to raw-cp-result.
// encode std-struct-ptr to raw-ptr.
type CompEncoder interface {
	EncodeData(srcStdStructPtr interface{}, dstRawPtr *ComponentData)
	EncodeState(srcStdStructPtr interface{}, dstRawPtr *ComponentState)
	EncodeInParams(srcStdStructPtr interface{}, dstRawPtr *InParams)
}

// CompDecoder is a protocol-level interface, convert raw-cp-result to structured-struct.
// decode raw-ptr to std-struct-ptr.
type CompDecoder interface {
	DecodeData(srcRawPtr ComponentData, dstStdStructPtr interface{})
	DecodeState(srcRawPtr ComponentState, dstStdStructPtr interface{})
	DecodeInParams(srcRawPtr InParams, dstStdStructPtr interface{})
}

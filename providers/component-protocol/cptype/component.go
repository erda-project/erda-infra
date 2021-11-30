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

// IComponent .
type IComponent interface {
	Initialize(sdk *SDK)
	Finalize(sdk *SDK)
	SkipOp(sdk *SDK) bool
	Visible(sdk *SDK) bool
	BeforeHandleOp(sdk *SDK)
	AfterHandleOp(sdk *SDK)
	StdStructuredPtrCreator() func() IStdStructuredPtr

	OperationsRegister
	Encoder
	Decoder
}

// IStdStructuredPtr .
type IStdStructuredPtr interface {
	DataPtr() interface{}
	StatePtr() interface{}
	InParamsPtr() interface{}
}

// OperationsRegister .
type OperationsRegister interface {
	RegisterInitializeOp() (opFunc OperationFunc)
	RegisterRenderingOp() (opFunc OperationFunc)
	RegisterCompStdOps() (opFuncs map[OperationKey]OperationFunc)
}

// Encoder is a protocol-level interface, convert structured-struct to raw-cp-result.
// encode std-struct-ptr to raw-ptr.
type Encoder interface {
	EncodeData(srcStdStructPtr interface{}, dstRawPtr *ComponentData)
	EncodeState(srcStdStructPtr interface{}, dstRawPtr *ComponentState)
	EncodeInParams(srcStdStructPtr interface{}, dstRawPtr *InParams)
}

// Decoder is a protocol-level interface, convert raw-cp-result to structured-struct.
// decode raw-ptr to std-struct-ptr.
type Decoder interface {
	DecodeData(srcRawPtr ComponentData, dstStdStructPtr interface{})
	DecodeState(srcRawPtr ComponentState, dstStdStructPtr interface{})
	DecodeInParams(srcRawPtr InParams, dstStdStructPtr interface{})
}

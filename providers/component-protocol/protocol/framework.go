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

package protocol

import (
	"context"
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

const (
	fieldStdStructuredPtr = "StdStructuredPtr"
)

// FRAMEWORK .
type FRAMEWORK struct {
	IC cptype.IComponent
}

// Render .
func (F FRAMEWORK) Render(ctx context.Context, c *cptype.Component, scenario cptype.Scenario, event cptype.ComponentEvent, gs *cptype.GlobalStateData) (err error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("component %s render panic: %v", c.Name, r)
			logrus.Error(msg)
			debug.PrintStack()
			err = fmt.Errorf(msg)
		}
	}()
	sdk := cputil.SDK(ctx)
	sdk.GlobalState = gs
	sdk.Comp = ensureCompFields(sdk, c)
	sdk.Event = event
	// structured comp ptr
	stdStructuredCompPtr := F.IC.StdStructuredPtrCreator()()
	universalStdPtr := F.makeUniversalPtr(stdStructuredCompPtr)
	// register operations
	F.registerOperations(sdk)
	// init
	F.IC.Initialize(sdk)
	// decoder
	F.IC.DecodeInParams(sdk.InParams, universalStdPtr.StdInParamsPtr)
	F.IC.DecodeState(c.State, universalStdPtr.StdStatePtr)
	F.IC.DecodeData(c.Data, universalStdPtr.StdDataPtr)
	//F.IC.DecodeOperations(c.Operations, universalStdPtr.StdOperationsPtr)
	// visible
	visible := F.IC.Visible(sdk)
	defer F.setVisible(sdk, visible)
	// handle op
	if !F.IC.SkipOp(sdk) {
		F.IC.BeforeHandleOp(sdk)
		F.handleOp(sdk, stdStructuredCompPtr, universalStdPtr)
		F.IC.AfterHandleOp(sdk)
		// encoder
		F.IC.EncodeData(universalStdPtr.StdDataPtr, &sdk.Comp.Data)
		F.IC.EncodeState(universalStdPtr.StdStatePtr, &sdk.Comp.State)
		F.IC.EncodeInParams(universalStdPtr.StdInParamsPtr, &sdk.InParams)
		// flat extra
		F.flatExtra(sdk.Comp)
		// finalize
		F.IC.Finalize(sdk)
	}
	return nil
}

func (F FRAMEWORK) setVisible(sdk *cptype.SDK, visible bool) {
	sdk.Comp.Options.Visible = visible
	// compatible with old version
	sdk.Comp.Props["visible"] = visible
}

func ensureCompFields(sdk *cptype.SDK, comp *cptype.Component) *cptype.Component {
	if sdk.InParams == nil {
		sdk.InParams = make(cptype.InParams)
	}
	if sdk.GlobalState == nil {
		sdk.GlobalState = &cptype.GlobalStateData{}
	}
	if sdk.CompOpFuncs == nil {
		sdk.CompOpFuncs = make(map[cptype.OperationKey]cptype.OperationFunc)
	}
	if comp.Data == nil {
		comp.Data = make(cptype.ComponentData)
	}
	if comp.State == nil {
		comp.State = make(cptype.ComponentState)
	}
	if comp.Props == nil {
		comp.Props = make(cptype.ComponentProps)
	}
	if comp.Operations == nil {
		comp.Operations = make(cptype.ComponentOperations)
	}
	if comp.Options == nil {
		comp.Options = &cptype.ComponentOptions{}
	}
	return comp
}

func (F FRAMEWORK) flatExtra(comp *cptype.Component) {
	if !comp.Options.FlatExtra {
		return
	}
	m := make(map[string]interface{})
	cputil.MustObjJSONTransfer(comp, &m)
	cputil.MustFlatMapMeta(m, comp.Options.RemoveExtraAfterFlat)
	cputil.MustObjJSONTransfer(&m, comp)
}

func (F FRAMEWORK) registerOperations(sdk *cptype.SDK) {
	sdk.RegisterOperation(cptype.InitializeOperation, F.IC.RegisterInitializeOp())
	sdk.RegisterOperation(cptype.RenderingOperation, F.IC.RegisterRenderingOp())
	// comp standard ops
	for opKey, opFunc := range F.IC.RegisterCompStdOps() {
		sdk.RegisterOperation(opKey, opFunc)
	}
}

func (F FRAMEWORK) handleOp(sdk *cptype.SDK, stdPtr cptype.IStdStructuredPtr, universalPtr *cptype.UniversalStructuredCompPtr) {
	op := sdk.Event.Operation
	opFunc, ok := sdk.CompOpFuncs[op]
	if !ok {
		panic(fmt.Errorf("component [%s] not supported operation [%s]", sdk.Comp.Name, op))
	}
	// do op
	opFunc(sdk)

	// ensure structured ptr
	setUniversalPtr(stdPtr, universalPtr)
}

func (F FRAMEWORK) makeUniversalPtr(stdPtr cptype.IStdStructuredPtr) *cptype.UniversalStructuredCompPtr {
	// set std ptr
	reflect.ValueOf(F.IC).Elem().FieldByName(fieldStdStructuredPtr).Set(reflect.ValueOf(stdPtr))

	// universal
	universalPtr := &cptype.UniversalStructuredCompPtr{}
	setUniversalPtr(stdPtr, universalPtr)
	return universalPtr
}

func setUniversalPtr(stdPtr cptype.IStdStructuredPtr, universalPtr *cptype.UniversalStructuredCompPtr) {
	universalPtr.StdDataPtr = stdPtr.DataPtr()
	universalPtr.StdStatePtr = stdPtr.StatePtr()
	universalPtr.StdInParamsPtr = stdPtr.InParamsPtr()
}

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
	originSDK := cputil.SDK(ctx)
	originSDK.GlobalState = gs
	sdk := originSDK.Clone()
	sdk.Comp = ensureCompFields(sdk, c)
	sdk.Event = event
	// structured comp ptr
	stdStructuredCompPtr := F.IC.StdStructuredPtrCreator()()
	F.injectStdStructurePtr(stdStructuredCompPtr)
	sdk.StdStructuredPtr = stdStructuredCompPtr
	// register operations
	F.registerOperations(sdk)
	// init
	F.IC.Initialize(sdk)
	// decoder
	F.IC.DecodeInParams(sdk.InParams, stdStructuredCompPtr.InParamsPtr())
	F.IC.DecodeState(c.State, stdStructuredCompPtr.StatePtr())
	F.IC.DecodeData(c.Data, stdStructuredCompPtr.DataPtr())
	// visible
	visible := F.IC.Visible(sdk)
	defer F.setVisible(sdk, visible)
	// handle op
	if !F.IC.SkipOp(sdk) {
		F.IC.BeforeHandleOp(sdk)
		F.handleOp(sdk, stdStructuredCompPtr)
		F.IC.AfterHandleOp(sdk)
		// encoder
		ensureCompFieldsBeforeEncode(sdk)
		F.IC.EncodeData(stdStructuredCompPtr.DataPtr(), &sdk.Comp.Data)
		F.IC.EncodeState(stdStructuredCompPtr.StatePtr(), &sdk.Comp.State)
		F.IC.EncodeInParams(stdStructuredCompPtr.InParamsPtr(), &sdk.InParams)
		// global state
		sdk.MergeClonedGlobalState()
		// flat extra
		F.flatExtra(sdk.Comp)
		// finalize
		F.IC.Finalize(sdk)
	}
	return nil
}

func (F FRAMEWORK) setVisible(sdk *cptype.SDK, visible bool) {
	sdk.Comp.Options.Visible = visible
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

// ensureCompFieldsBeforeEncode ensure data/state/inParams to empty to avoid `omitempty` issue.
func ensureCompFieldsBeforeEncode(sdk *cptype.SDK) {
	sdk.Comp.Data = cptype.ComponentData{}
	sdk.Comp.State = cptype.ComponentState{}
	sdk.InParams = cptype.InParams{}
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
	// comp non-standard ops
	for opKey, opFunc := range F.IC.RegisterCompNonStdOps() {
		sdk.RegisterOperation(opKey, opFunc)
	}
}

func (F FRAMEWORK) handleOp(sdk *cptype.SDK, stdPtr cptype.IStdStructuredPtr) {
	op := sdk.Event.Operation
	opFunc, ok := sdk.CompOpFuncs[op]
	if !ok {
		panic(fmt.Errorf("component [%s] not supported operation [%s]", sdk.Comp.Name, op))
	}
	// nil opFunc equals to empty op func
	if opFunc == nil {
		return
	}
	// do op
	stdResp := opFunc(sdk)
	// if stdResp is not nil, set stdPtr to stdResp
	if stdResp != nil {
		reflect.ValueOf(stdPtr).Elem().Set(reflect.ValueOf(stdResp).Elem())
	}
}

func (F FRAMEWORK) injectStdStructurePtr(stdPtr cptype.IStdStructuredPtr) {
	// inject std ptr
	reflect.ValueOf(F.IC).Elem().FieldByName(fieldStdStructuredPtr).Set(reflect.ValueOf(stdPtr))
}

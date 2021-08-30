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

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// CompRenderSpec .
type CompRenderSpec struct {
	// Scenario key
	Scenario string `json:"scenario"`
	// CompName is component name
	CompName string `json:"name"`
	// RenderC used to created component
	RenderC RenderCreator
}

// RenderCreator .
type RenderCreator func() CompRender

// CompRender .
type CompRender interface {
	Render(ctx context.Context, c *cptype.Component, scenario cptype.Scenario, event cptype.ComponentEvent, gs *cptype.GlobalStateData) error
}

// MustRegisterComponent .
func MustRegisterComponent(r *CompRenderSpec) {
	if err := RegisterComponent(r); err != nil {
		panic(err)
	}
}

// RegisterComponent register a component under scenario
func RegisterComponent(r *CompRenderSpec) error {
	if r == nil {
		return fmt.Errorf("register request is empty")
	}
	if r.Scenario == "" {
		// use default scenario
		r.Scenario = cptype.DefaultComponentNamespace
	}
	if r.CompName == "" {
		return fmt.Errorf("component name is empty")
	}

	logrus.Infof("begin register component, scenario: %s, component: %s", r.Scenario, r.CompName)
	// if scenario not exit, crate it
	if _, ok := ScenarioRenders[r.Scenario]; !ok {
		s := make(ScenarioRender)
		ScenarioRenders[r.Scenario] = &s
	}
	// if compName key not exist, create it and the CompRenderSpec
	s := ScenarioRenders[r.Scenario]
	if _, ok := (*s)[r.CompName]; !ok {
		(*s)[r.CompName] = r
	} else {
		err := fmt.Errorf("register render failed, component [%s] already exist", r.CompName)
		return err
	}
	logrus.Infof("register component render success, scenario: %s, component: %s", r.Scenario, r.CompName)
	return nil
}

type emptyComp struct{}

// Render .
func (ca *emptyComp) Render(ctx context.Context, c *cptype.Component, scenario cptype.Scenario, event cptype.ComponentEvent, gs *cptype.GlobalStateData) error {
	return nil
}

var emptyRenderFunc = func() CompRender { return &emptyComp{} }

// getCompRender .
func getCompRender(ctx context.Context, r ScenarioRender, compName, typ string) (*CompRenderSpec, error) {
	if len(r) == 0 {
		return nil, fmt.Errorf(i18n(ctx, "scenario.render.is.empty"))
	}
	if compName == "" {
		return nil, fmt.Errorf(i18n(ctx, "component.name.is.empty"))
	}
	var c *CompRenderSpec
	if _, ok := r[compName]; !ok {
		// component not exist
		return nil, fmt.Errorf(i18n(ctx, "${component %s missing renderCreator}", compName))
	}
	c = r[compName]
	if c == nil {
		return nil, fmt.Errorf(i18n(ctx, "component.render.is.empty"))
	}
	return c, nil
}

// protoCompStateRending .
func protoCompStateRending(ctx context.Context, p *cptype.ComponentProtocol, r cptype.RendingItem) error {
	if p == nil {
		return fmt.Errorf(i18n(ctx, "protocol.empty"))
	}
	pc, err := getProtoComp(ctx, p, r.Name)
	if err != nil {
		logrus.Errorf("failed to get protocol component, err: %v", err)
		return err
	}
	// inParams
	inParams := ctx.Value(cptype.GlobalInnerKeyCtxSDK).(*cptype.SDK).InParams
	for _, state := range r.State {
		// parse state bound info
		stateFrom, stateFromKey, err := parseStateBound(state.Value)
		if err != nil {
			logrus.Errorf("failed to parse component state bound, component: %s, state bound: %#v", pc.Name, state)
			return err
		}
		var stateFromValue interface{}
		switch stateFrom {
		case cptype.InParamsStateBindingKey: // {{ inParams.key }} 表示从 inParams 绑定
			stateFromValue = getProtoInParamsValue(inParams, stateFromKey)
		default: // 否则，从其他组件 state 绑定
			// get bound key value
			stateFromValue, err = getProtoCompStateValue(ctx, p, stateFrom, stateFromKey)
			if err != nil {
				logrus.Errorf("failed to get component state value, component: %s, key: %s", stateFrom, stateFromKey)
				return err
			}
		}
		// set component state value
		err = setCompStateValueFromComps(pc, state.Name, stateFromValue)
		if err != nil {
			logrus.Errorf("failed to set component state, component: %s, state key: %s, value: %#v", pc.Name, state.Name, stateFromValue)
			return err
		}
	}
	return nil
}

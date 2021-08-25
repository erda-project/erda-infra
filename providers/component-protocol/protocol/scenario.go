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
	"fmt"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// ScenarioRender is a group of component renders.
// key: componentName
// value: componentRender
type ScenarioRender map[string]*CompRenderSpec

// ScenarioRenders contains all scenario renders.
var ScenarioRenders = make(map[string]*ScenarioRender)

// getScenarioRenders .
func getScenarioRenders(scenario string) (*ScenarioRender, error) {
	renders := ScenarioRenders[scenario]
	if renders == nil {
		renders = &ScenarioRender{}
	}

	defaultRenders := ScenarioRenders[cptype.DefaultComponentNamespace]
	if defaultRenders == nil {
		defaultRenders = &ScenarioRender{}
	}

	// append component render from default component namespace
	p, ok := defaultProtocols[scenario]
	if !ok {
		return nil, fmt.Errorf("failed to get scenario renders, default protocol not exist, scenario: %s", scenario)
	}
	for compName := range p.Components {
		// skip if scenario-level component render exist
		_, renderExist := (*renders)[compName]
		if renderExist {
			continue
		}
		// if component render not exist, add render from default component namespace
		defaultCompRender, dOK := (*defaultRenders)[compName]
		if dOK {
			(*renders)[compName] = defaultCompRender
			continue
		}
		// build-in component renders
		switch compName {
		case "Container", "LRContainer", "RowContainer", "SplitPage", "Popover", "Title", "Drawer":
			(*renders)[compName] = &CompRenderSpec{RenderC: emptyRenderFunc}
			continue
		}
		return nil, fmt.Errorf("failed to found component render neither in scenario renders nor default renders, "+
			"sceneario: %s, component: %s", scenario, compName)
	}

	return renders, nil
}

// getScenarioKey get scenario key from protocol.
// return scenarioType if not empty.
func getScenarioKey(req cptype.Scenario) (string, error) {
	if req.ScenarioType == "" && req.ScenarioKey == "" {
		return "", fmt.Errorf("scenario.is.empty")
	}
	if req.ScenarioType != "" {
		return req.ScenarioType, nil
	}
	return req.ScenarioKey, nil
}

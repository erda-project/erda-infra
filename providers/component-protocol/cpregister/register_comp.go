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

package cpregister

import (
	"reflect"

	"github.com/erda-project/erda-infra/providers/component-protocol/components/defaults"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/protocol"
)

// RegisterComponent register legacy component which implements protocol.RenderCreator
func RegisterComponent(scenario, componentName string, componentCreator protocol.ComponentCreator) {
	protocol.MustRegisterComponent(&protocol.CompRenderSpec{
		Scenario: scenario,
		CompName: componentName,
		RenderC:  nil,
		Creator: func() cptype.IComponent {
			compInstance := componentCreator()
			ref := reflect.ValueOf(compInstance)
			ref.Elem().FieldByName(cptype.FieldImplForInject).Set(ref)
			ref.Elem().FieldByName(defaults.FieldActualImplRef).Set(ref)
			return compInstance
		},
	})
}

// RegisterLegacyComponent register legacy component which implements protocol.RenderCreator
// For most scenarios, you should use RegisterComponent.
func RegisterLegacyComponent(scenario, componentName string, componentCreator protocol.RenderCreator) {
	protocol.MustRegisterComponent(&protocol.CompRenderSpec{
		Scenario: scenario,
		CompName: componentName,
		RenderC:  componentCreator,
		Creator:  nil,
	})
}

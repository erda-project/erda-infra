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

package base

import (
	"reflect"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/component-protocol/cpregister"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/protocol"
)

type creators struct {
	RenderCreator    protocol.RenderCreator
	ComponentCreator protocol.ComponentCreator
}

// InitProvider register component as provider to scenario-namespace.
func InitProvider(scenario, compName string) {
	InitProviderWithCreator(scenario, compName, nil)
}

// InitProviderToDefaultNamespace register component as provider to default-namespace.
func InitProviderToDefaultNamespace(compName string, creator servicehub.Creator) {
	initProviderToNamespace(defaultComponentProviderNamespace, compName, creator)
}

// InitProviderWithCreator register component as provider with custom providerCreator to scenario-namespace.
func InitProviderWithCreator(scenario, compName string, creator servicehub.Creator) {
	initProviderToNamespace(scenario, compName, creator)
}

// initProviderToNamespace register component as provider to specific namespace.
func initProviderToNamespace(scenario, compName string, creator servicehub.Creator) {
	// generate std providerName
	providerName := MakeComponentProviderName(scenario, compName)
	if creator == nil {
		creator = func() servicehub.Provider { return &DefaultProvider{} }
	}
	// register to servicehub
	servicehub.Register(providerName, &servicehub.Spec{Creator: creator})
	// add to explicit provider creator map for hubListener to auto register to hub.config
	cpregister.AllExplicitProviderCreatorMap[providerName] = creator

	// generate creators compatible for IComponent and old CompRender
	creators := func() creators {
		switch creator().(type) {
		case cptype.IComponent:
			return creators{ComponentCreator: func() cptype.IComponent {
				rr := creator().(cptype.IComponent)
				ref := reflect.ValueOf(rr)
				ref.Elem().FieldByName("Impl").Set(ref)
				return rr
			}}
		case protocol.CompRender:
			return creators{RenderCreator: func() protocol.CompRender { return creator().(protocol.CompRender) }}
		default:
			return creators{RenderCreator: func() protocol.CompRender { return &DefaultProvider{} }}
		}
	}()

	// register protocol comp
	protocol.MustRegisterComponent(&protocol.CompRenderSpec{
		Scenario: scenario,
		CompName: compName,
		RenderC:  creators.RenderCreator,
		Creator:  creators.ComponentCreator,
	})
}

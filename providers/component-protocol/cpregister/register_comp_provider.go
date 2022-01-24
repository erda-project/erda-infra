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

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

// ComponentCreatorAndProvider used for RegisterProviderComponent.
type ComponentCreatorAndProvider interface {
	cptype.IComponent
	servicehub.Provider
}

// Option represents options when register.
type Option struct {
	providerSpec *servicehub.Spec
}

// OpFunc do function for option.
type OpFunc func(*Option)

// WithCustomProviderSpecButCreator use custom spec as base but ignore creator.
func (o *Option) WithCustomProviderSpecButCreator(spec *servicehub.Spec) {
	o.providerSpec = spec
}

// RegisterProviderComponent register a provider to component. Normally this provider will use servicehub's auto dependency injection feature.
func RegisterProviderComponent(scenario, componentName string, providerPtr ComponentCreatorAndProvider, opFuncs ...OpFunc) {
	// handle options
	opt := Option{
		providerSpec: &servicehub.Spec{},
	}
	for _, opFunc := range opFuncs {
		opFunc(&opt)
	}

	// generate provider name
	providerName := cputil.MakeComponentProviderName(scenario, componentName)

	// register component
	RegisterComponent(scenario, componentName, func() cptype.IComponent {
		newProviderPtr := reflect.New(reflect.TypeOf(providerPtr).Elem())
		newProviderPtr.Elem().Set(reflect.ValueOf(providerPtr).Elem())
		copied := newProviderPtr.Interface()
		return copied.(ComponentCreatorAndProvider)
	})

	// register as provider
	opt.providerSpec.Creator = func() servicehub.Provider { return providerPtr }
	servicehub.Register(providerName, opt.providerSpec)

	// mark for auto servicehub config adding
	AllExplicitProviderCreatorMap[providerName] = nil
}

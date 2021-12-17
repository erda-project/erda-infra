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

package filter

import (
	"reflect"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/filter/impl"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/protocol"
)

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	p.DefaultFilter = impl.DefaultFilter{}
	v := reflect.ValueOf(p)
	v.Elem().FieldByName("Impl").Set(v)
	compName := "filter"
	if ctx.Label() != "" {
		compName = ctx.Label()
	}
	protocol.MustRegisterComponent(&protocol.CompRenderSpec{
		Scenario: "filter-demo",
		CompName: compName,
		Creator:  func() cptype.IComponent { return p },
	})
	return nil
}

// Provide .
func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	return p
}

func init() {
	servicehub.Register("component-protocol.components.filter-demo.filter", &servicehub.Spec{
		Creator: func() servicehub.Provider { return &provider{} },
	})
}

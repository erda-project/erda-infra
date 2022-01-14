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

package linegraph

import (
	"reflect"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/commodel/data-structure"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/linegraph"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/linegraph/impl"
	"github.com/erda-project/erda-infra/providers/component-protocol/cpregister/base"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/protocol"
)

type provider struct {
	impl.DefaultLineGraph
}

// RegisterInitializeOp .
func (p *provider) RegisterInitializeOp() (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) {

		// Demo case 1
		d := linegraph.New("line graph demo")
		d.SetXAxis("Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday")
		d.SetYAxis("Dimension", 1, 2, 3, 4, 5, 6, 7)
		d.SetXOptions(&linegraph.Options{
			Structure: &structure.DataStructure{Type: structure.String},
		})
		d.SetYOptions([]*linegraph.Options{
			{Dimension: "Dimension", Structure: &structure.DataStructure{Type: structure.Number}},
			{Dimension: "Dimension2", Structure: &structure.DataStructure{Type: structure.Number}},
		}...)
		d.SetYAxis("Dimension2", 7, 6, 5, 4, 3, 2, 1)

		// Demo case 2
		d = linegraph.NewDataBuilder().
			WithTitle("line graph demo").
			WithXAxis("Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday").
			WithXOptions(linegraph.NewOptionsBuilder().WithType(structure.String).Build()).
			WithYAxis("Dimension", 1, 2, 3, 4, 5, 6, 7).
			WithYAxis("Dimension2", 7, 6, 5, 4, 3, 2, 1).
			WithYOptions([]*linegraph.Options{
				linegraph.NewOptionsBuilder().WithDimension("Dimension").WithType(structure.Number).Build(),
				linegraph.NewOptionsBuilder().WithDimension("Dimension2").WithType(structure.Number).Build(),
			}...).Build()

		p.StdDataPtr = d
	}
}

// RegisterRenderingOp .
func (p *provider) RegisterRenderingOp() (opFunc cptype.OperationFunc) {
	return p.RegisterInitializeOp()
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	p.DefaultLineGraph = impl.DefaultLineGraph{}
	v := reflect.ValueOf(p)
	v.Elem().FieldByName("Impl").Set(v)
	compName := "linegraph"
	if ctx.Label() != "" {
		compName = ctx.Label()
	}
	protocol.MustRegisterComponent(&protocol.CompRenderSpec{
		Scenario: "linegraph-demo",
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
	base.InitProviderWithCreator("linegraph-demo", "linegraph", func() servicehub.Provider { return &provider{} })
}

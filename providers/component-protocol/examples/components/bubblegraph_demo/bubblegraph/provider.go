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

package bubblegraph

import (
	"reflect"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/bubblegraph"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/bubblegraph/impl"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/protocol"
)

type provider struct {
	impl.DefaultBubbleGraph
}

// RegisterInitializeOp .
func (p *provider) RegisterInitializeOp() (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) {
		data := bubblegraph.NewDataBuilder().
			WithTitle("test bubble graph component").
			WithBubble(bubblegraph.NewBubbleBuilder().
				WithX(&bubblegraph.Axis{Value: 1}).
				WithY(&bubblegraph.Axis{Value: 100}).
				WithSize(&bubblegraph.BubbleSize{Value: 10}).
				WithGroup("test group").
				WithDimension("test dimension").
				Build()).
			WithBubble(bubblegraph.NewBubbleBuilder().
				WithX(&bubblegraph.Axis{Value: 2}).
				WithY(&bubblegraph.Axis{Value: 200}).
				WithSize(&bubblegraph.BubbleSize{Value: 20}).
				WithGroup("test group").
				WithDimension("test dimension").
				Build()).
			WithBubble(bubblegraph.NewBubbleBuilder().
				WithX(&bubblegraph.Axis{Value: 3}).
				WithY(&bubblegraph.Axis{Value: 300}).
				WithSize(&bubblegraph.BubbleSize{Value: 30}).
				WithGroup("test group").
				WithDimension("test dimension").
				Build()).
			WithBubble(bubblegraph.NewBubbleBuilder().
				WithX(&bubblegraph.Axis{Value: 4}).
				WithY(&bubblegraph.Axis{Value: 400}).
				WithSize(&bubblegraph.BubbleSize{Value: 40}).
				WithGroup("test group").
				WithDimension("test dimension").
				Build()).
			WithBubble(bubblegraph.NewBubbleBuilder().
				WithX(&bubblegraph.Axis{Value: 5}).
				WithY(&bubblegraph.Axis{Value: 400}).
				WithSize(&bubblegraph.BubbleSize{Value: 50}).
				WithGroup("test group").
				WithDimension("test dimension").
				Build()).
			WithBubble(bubblegraph.NewBubbleBuilder().
				WithX(&bubblegraph.Axis{Value: 6}).
				WithY(&bubblegraph.Axis{Value: 600}).
				WithSize(&bubblegraph.BubbleSize{Value: 60}).
				WithGroup("test group").
				WithDimension("test dimension").
				Build()).
			WithBubble(bubblegraph.NewBubbleBuilder().
				WithX(&bubblegraph.Axis{Value: 7}).
				WithY(&bubblegraph.Axis{Value: 700}).
				WithSize(&bubblegraph.BubbleSize{Value: 70}).
				WithGroup("test group").
				WithDimension("test dimension").
				Build()).
			Build()

		p.StdDataPtr = data
	}
}

// RegisterRenderingOp .
func (p *provider) RegisterRenderingOp() (opFunc cptype.OperationFunc) {
	return p.RegisterInitializeOp()
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	p.DefaultBubbleGraph = impl.DefaultBubbleGraph{}
	v := reflect.ValueOf(p)
	v.Elem().FieldByName("Impl").Set(v)
	compName := "bubblegraph"
	if ctx.Label() != "" {
		compName = ctx.Label()
	}
	protocol.MustRegisterComponent(&protocol.CompRenderSpec{
		Scenario: "bubblegraph-demo",
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
	servicehub.Register("component-protocol.components.bubblegraph-demo", &servicehub.Spec{
		Creator: func() servicehub.Provider { return &provider{} },
	})
}

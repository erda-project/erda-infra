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
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/bubblegraph"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/bubblegraph/impl"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/commodel/data-structure"
	"github.com/erda-project/erda-infra/providers/component-protocol/cpregister/base"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

type provider struct {
	impl.DefaultBubbleGraph
}

// RegisterInitializeOp .
func (p *provider) RegisterInitializeOp() (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) {
		data := bubblegraph.NewDataBuilder().
			WithTitle("test bubble graph component").
			WithXOptions(bubblegraph.NewOptionsBuilder().WithType(structure.Number).Build()).
			WithYOptions(bubblegraph.NewOptionsBuilder().WithType(structure.Number).Build()).
			WithBubble(bubblegraph.NewBubbleBuilder().
				WithValueX(1).
				WithValueY(100).
				WithValueSize(10).
				WithGroup("test group").
				WithDimension("test dimension").
				Build()).
			WithBubble(bubblegraph.NewBubbleBuilder().
				WithValueX(2).
				WithValueY(200).
				WithValueSize(20).
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

func init() {
	base.InitProviderWithCreator("bubblegraph-demo", "bubblegraph", func() servicehub.Provider { return &provider{} })
}

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

package indicatorcard

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/components/indicatorcard"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/indicatorcard/impl"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

type provider struct {
	impl.DefaultIndicatorCard
}

// RegisterInitializeOp .
func (p *provider) RegisterInitializeOp() (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) {
		data := indicatorcard.Data{
			[]*indicatorcard.IndicatorCard{
				{
					Value: 1,
					Title: "test_title",
					Tips:  "this is a demo",
					Unit:  "ms",
				},
				{
					Value: 2,
					Title: "test_title2",
					Tips:  "this is a demo",
					Unit:  "",
				},
				{
					Value: 3,
					Title: "test_title3",
					Tips:  "this is a demo",
					Unit:  "s",
				},
			},
		}
		p.StdDataPtr = &data
	}
}

// RegisterRenderingOp .
func (p *provider) RegisterRenderingOp() (opFunc cptype.OperationFunc) {
	return p.RegisterInitializeOp()
}

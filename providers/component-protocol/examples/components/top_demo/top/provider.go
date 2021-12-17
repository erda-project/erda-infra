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

package top

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/components/topn"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/topn/impl"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

type provider struct {
	impl.DefaultTop
}

// RegisterInitializeOp .
func (p *provider) RegisterInitializeOp() (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) {
		data := topn.Data{
			List: []topn.Record{
				{
					Title: "rps-max-top5",
					Items: []topn.Item{
						{
							ID:    "id2",
							Name:  "name2",
							Value: 1.2,
							Unit:  "req/s",
						},
					},
				},
				{
					Title: "rps-max-top5",
					Items: []topn.Item{
						{
							ID:    "id2",
							Name:  "name2",
							Value: 1.2,
							Unit:  "req/s",
						},
					},
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

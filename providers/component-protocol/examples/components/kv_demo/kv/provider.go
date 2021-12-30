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

package kv

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/components/kv"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/kv/impl"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

type provider struct {
	impl.DefaultKV
}

// RegisterInitializeOp .
func (p *provider) RegisterInitializeOp() (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) {
		data := kv.Data{
			List: []*kv.KV{
				{
					Key:    "key",
					SubKey: "ms",
					Value:  1,
					Tip:    "tip",
				},
				{
					Key:    "key2",
					SubKey: "",
					Value:  2,
					Tip:    "tip",
				},
				{
					Key:    "key3",
					SubKey: "s",
					Value:  1,
					Tip:    "tip",
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

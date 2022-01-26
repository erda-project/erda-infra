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

package paralleldemo

import (
	"fmt"

	"github.com/labstack/gommon/random"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/kv"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/kv/impl"
	"github.com/erda-project/erda-infra/providers/component-protocol/cpregister"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
	"github.com/erda-project/erda-infra/providers/i18n"
)

type comp struct {
	impl.DefaultKV

	InParams *InParams

	Tran i18n.Translator `translator:"hello"`
}

func (c *comp) Init(ctx servicehub.Context) error {
	c.InParams = &InParams{}
	return nil
}

// InParams .
type InParams struct {
	StartTime uint64 `json:"startTime,omitempty"`
	EndTime   uint64 `json:"endTime,omitempty"`
}

func (c *comp) RegisterInitializeOp() (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) cptype.IStdStructuredPtr {
		fmt.Printf("%p\n", c.InParams)
		if c.InParams.StartTime == 0 {
			panic("missing InParams.StartTime")
		}
		if c.InParams.EndTime == 0 {
			panic("missing InParams.EndTime")
		}
		fmt.Println("startTime", c.InParams.StartTime)
		fmt.Println("endTime", c.InParams.EndTime)
		return &impl.StdStructuredPtr{
			StdDataPtr: &kv.Data{List: []*kv.KV{{
				Key:   sdk.Comp.Name,
				Value: random.String(4),
			}}},
		}
	}
}

func (c *comp) CustomInParamsPtr() interface{} {
	return c.InParams
}

func (c *comp) EncodeFromCustomInParams(customInParamsPtr interface{}, stdInParamsPtr *cptype.ExtraMap) {
	cputil.MustObjJSONTransfer(customInParamsPtr, stdInParamsPtr)
}

func (c *comp) DecodeToCustomInParams(stdInParamsPtr *cptype.ExtraMap, customInParamsPtr interface{}) {
	cputil.MustObjJSONTransfer(stdInParamsPtr, customInParamsPtr)
}

func init() {
	cpregister.RegisterProviderComponent("parallel-demo", "kv", &comp{})
}

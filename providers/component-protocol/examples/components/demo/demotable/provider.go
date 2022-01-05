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

package demotable

import (
	"reflect"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/component-protocol/cpregister/base"
	"github.com/erda-project/erda-infra/providers/component-protocol/protocol"
)

// Interface export ability for demotable
type Interface interface {
	protocol.CompRender
}

type config struct {
	Scenario string
	Name     string
}

type provider struct {
	Cfg *config
}

func init() {
	base.InitProviderWithCreator("demo", "table", func() servicehub.Provider {
		interfaceType := reflect.TypeOf((*Interface)(nil)).Elem()
		return &servicehub.Spec{
			Types: []reflect.Type{interfaceType},
			ConfigFunc: func() interface{} {
				return &config{}
			},
			Creator: func() servicehub.Provider {
				return &provider{}
			},
		}
	})
}

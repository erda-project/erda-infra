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

package kanban

import (
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/component-protocol/cpregister"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

func init() {
	// register provider
	servicehub.Register("your-provider-name", &servicehub.Spec{
		Dependencies: []string{"i18n"},
		Creator:      func() servicehub.Provider { return &component{} },
	})
}

func (c *component) Init(ctx servicehub.Context) error {
	// register component
	cpregister.RegisterComponent("kanban-demo", "xxx", func() cptype.IComponent { return c })
	cpregister.RegisterComponent("kanban-demo", "sssss", func() cptype.IComponent { return c })
	cpregister.RegisterComponent("kanban-demo", "kanban", func() cptype.IComponent { return c })
	return nil
}

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

package posthook

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// OnlyReturnRenderingComps only return rendering components.
func OnlyReturnRenderingComps(renderingItems []cptype.RendingItem, req *cptype.ComponentProtocol) {
	if req.Options != nil && req.Options.ReturnAllComponents {
		return
	}
	// init new components map
	onlyReturnComps := make(map[string]*cptype.Component, len(req.Components))
	// construct map for easy use
	renderingItemByName := map[string]struct{}{}
	for _, item := range renderingItems {
		renderingItemByName[item.Name] = struct{}{}
	}
	// only return rendering components
	for name := range req.Components {
		if _, ok := renderingItemByName[name]; ok {
			onlyReturnComps[name] = req.Components[name]
		}
	}
	req.Components = onlyReturnComps
}

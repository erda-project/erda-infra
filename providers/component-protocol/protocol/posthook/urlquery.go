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
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

// HandleURLQuery .
func HandleURLQuery(renderingItems []cptype.RendingItem, req *cptype.ComponentProtocol) {
	for _, comp := range req.Components {
		if comp.Options == nil {
			continue
		}
		urlQuery := comp.Options.URLQuery
		if len(urlQuery) == 0 {
			continue
		}

		// we set urlQuery to comp state currently
		if comp.State == nil {
			comp.State = make(cptype.ComponentState)
		}
		comp.State[cputil.MakeCompURLQueryKey(comp.Name)] = urlQuery
	}
}

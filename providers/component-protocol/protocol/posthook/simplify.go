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

// SimplifyProtocol simplify protocol if could.
func SimplifyProtocol(renderingItems []cptype.RendingItem, req *cptype.ComponentProtocol) {
	for _, comp := range req.Components {
		simplifyComp(comp)
	}
}

func simplifyComp(comp *cptype.Component) {
	if len(comp.Data) == 0 {
		comp.Data = nil
	}
	if len(comp.State) == 0 {
		comp.State = nil
	}
	if len(comp.Operations) == 0 {
		comp.Operations = nil
	}
	if comp.Options != nil {
		if !comp.Options.Visible &&
			!comp.Options.AsyncAtInit &&
			!comp.Options.FlatExtra &&
			!comp.Options.RemoveExtraAfterFlat &&
			len(comp.Options.UrlQuery) == 0 &&
			(comp.Options.ContinueRender == nil || len(comp.Options.ContinueRender.OpKey) == 0) {
			comp.Options = nil
		}
	}
	if len(comp.Props) == 0 {
		comp.Props = nil
	}
}

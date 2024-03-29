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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

func TestHandleContinueRender(t *testing.T) {
	renderingItems := []cptype.RendingItem{{Name: "page"}, {Name: "list"}}
	p := cptype.ComponentProtocol{
		Components: map[string]*cptype.Component{
			"page": {
				Options: &cptype.ComponentOptions{
					ContinueRender: &cptype.ContinueRender{
						OpKey: "loadPageMore",
					},
				},
			},
			"list": {
				Options: &cptype.ComponentOptions{
					ContinueRender: &cptype.ContinueRender{
						OpKey: "loadListDetail",
					},
				},
			},
		},
		Options: nil,
	}
	HandleContinueRender(renderingItems, &p)
	assert.NotNil(t, p.Options)
	assert.True(t, len(p.Options.ParallelContinueRenders) == 2)
	assert.True(t, p.Options.ParallelContinueRenders["page"].OpKey == "loadPageMore")
	assert.True(t, p.Options.ParallelContinueRenders["list"].OpKey == "loadListDetail")
}

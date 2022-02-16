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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

func Test_simplifyComp(t *testing.T) {
	c1 := &cptype.Component{
		Data: map[string]interface{}{},
		State: map[string]interface{}{
			"A": "B",
		},
		Options: &cptype.ComponentOptions{
			Visible:              false,
			AsyncAtInit:          false,
			ContinueRender:       &cptype.ContinueRender{OpKey: ""},
			FlatExtra:            false,
			RemoveExtraAfterFlat: false,
		},
	}
	simplifyComp(c1)
	assert.Nil(t, c1.Data)
	assert.NotNil(t, c1.State)
	assert.Nil(t, c1.Options)

	b, err := json.Marshal(c1)
	assert.NoError(t, err)
	fmt.Println(string(b))
}

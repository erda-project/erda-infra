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

package protocol

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

func Test_doParallelRendering(t *testing.T) {
	py := `
hierarchy:
  root: page
  structure:
    page:
      - filter
      - grid
      - chartList
    grid:
      - topN_1
      - topN_2
    chartList:
      - chart_1
      - chart_2
  parallel:
    page:
      - filter
      - grid
    grid:
      - topN_1
      - topN_2
    chartList:
      - chart_1
      - chart_2
`
	var p cptype.ComponentProtocol
	assert.NoError(t, yaml.Unmarshal([]byte(py), &p))
	orders, err := calculateDefaultRenderOrderByHierarchy(&p)
	assert.NoError(t, err)
	var compRenderings []cptype.RendingItem
	for _, order := range orders {
		compRenderings = append(compRenderings, cptype.RendingItem{Name: order})
	}
	node, err := parseParallelRendering(&p, compRenderings)
	assert.NoError(t, err)
	//spew.Dump(node)
	fmt.Println(node.String())
}

func Test_removeOneNode(t *testing.T) {
	nodes := []*Node{{Name: "n1"}, {Name: "n2"}, {Name: "n3"}}
	removeOneNode(&nodes, "n2")
	assert.Equal(t, 2, len(nodes))
}

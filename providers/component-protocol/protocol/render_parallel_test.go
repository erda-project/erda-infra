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
      left:
       leftContent
      right:
       rightContent
    leftContent:
      - head
      - workTabs
      - workContainer
      - messageTabs
      - messageContainer
    rightContent:
      - userProfile
    workContainer:
      - workCards
      - workList
    messageContainer:
      - messageList
    workList:
      filter:
      - workListFilter
  parallel:
    page:
      - leftContent
      - rightContent
    leftContent:
      - head
      - workTabs
      - messageTabs
    workContainer:
      - workCards
      - workList
    messageContainer:
      - messageList
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
	fmt.Println(node.String())

	//assert.NoError(t, renderFromNode(context.Background(), nil, nil, node))
}

func Test_removeOneNode(t *testing.T) {
	nodes := []*Node{{Name: "n1"}, {Name: "n2"}, {Name: "n3"}}
	removeOneNode(&nodes, "n2")
	assert.Equal(t, 2, len(nodes))
}

func TestNode_calcRenderableNextNodes(t *testing.T) {
	pageNode := &Node{Name: "page", doneNextNodesByName: map[string]*Node{}}
	filterNode := &Node{Name: "filter", doneNextNodesByName: map[string]*Node{}}
	gridNode := &Node{Name: "grid", doneNextNodesByName: map[string]*Node{}}
	chartNode := &Node{Name: "chart", doneNextNodesByName: map[string]*Node{}}

	pageNode.NextNodes = []*Node{filterNode, gridNode, chartNode}
	filterNode.Parallel = true
	gridNode.Parallel = true
	chartNode.Parallel = false

	renderableNodes := pageNode.calcRenderableNextNodes()
	assert.Equal(t, 2, len(renderableNodes))
	fmt.Println(renderableNodes)

	pageNode.doneNextNodesByName = map[string]*Node{filterNode.Name: filterNode, gridNode.Name: gridNode}
	renderableNodes = pageNode.calcRenderableNextNodes()
	assert.Equal(t, 1, len(renderableNodes))
	fmt.Println(renderableNodes)

	//////

	pageNode.NextNodes = []*Node{filterNode, gridNode, chartNode}
	pageNode.doneNextNodesByName = map[string]*Node{}
	filterNode.Parallel = true
	gridNode.Parallel = false
	chartNode.Parallel = true

	renderableNodes = pageNode.calcRenderableNextNodes()
	assert.Equal(t, 1, len(renderableNodes))
	fmt.Println(renderableNodes)

	pageNode.doneNextNodesByName = map[string]*Node{filterNode.Name: filterNode}
	renderableNodes = pageNode.calcRenderableNextNodes()
	assert.Equal(t, 2, len(renderableNodes))
	fmt.Println(renderableNodes)
}

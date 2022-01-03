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
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

// Node used for parallel rendering.
type Node struct {
	Name     string
	Parallel bool

	NextNodes       []*Node
	nextNodesByName map[string]*Node
	PreviousNode    *Node

	BindingStates []cptype.RendingState

	doneNextNodesByName map[string]*Node
}

func (n *Node) toRendingItem() cptype.RendingItem {
	return cptype.RendingItem{Name: n.Name, State: n.BindingStates}
}

func printIndent(w io.Writer, repeat int) {
	if repeat == 0 {
		repeat = 1
	}
	for i := 0; i < repeat; i++ {
		fmt.Fprintf(w, " ")
	}
}
func printNode(w io.Writer, n *Node) {
	if n.Parallel {
		fmt.Fprintf(w, "[P] %s\n", n.Name)
	} else {
		fmt.Fprintf(w, "[S] %s\n", n.Name)
	}
}

func (n *Node) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "root: %s\n", n.Name)
	depth := 1
	n.printNexts(w, depth)
	return w.String()
}
func (n *Node) printNexts(w io.Writer, depth int) {
	for _, next := range n.NextNodes {
		printIndent(w, depth*2)
		printNode(w, next)
		next.printNexts(w, depth+1)
	}
}

func makeSerialNode(item cptype.RendingItem) *Node {
	return &Node{Name: item.Name, Parallel: false, BindingStates: item.State, doneNextNodesByName: map[string]*Node{}}
}
func (n *Node) addNext(next *Node) {
	// set next
	n.NextNodes = append(n.NextNodes, next)
	if n.nextNodesByName == nil {
		n.nextNodesByName = make(map[string]*Node)
	}
	n.nextNodesByName[next.Name] = next
	// set previous
	next.PreviousNode = n
}
func removeOneNode(nodes *[]*Node, removeNodeName string) {
	index := -1
	for i, node := range *nodes {
		if node.Name == removeNodeName {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	*nodes = append((*nodes)[:index], (*nodes)[index+1:]...)
}
func (n *Node) cutOffPrevious() {
	previousNode := n.PreviousNode
	if previousNode == nil {
		return
	}

	// cut off from node's previous node
	n.PreviousNode = nil

	// cut off from previous node's next node
	delete(previousNode.nextNodesByName, n.Name)
	removeOneNode(&previousNode.NextNodes, n.Name)
}
func (n *Node) linkSubParallelNode(subNode *Node) {
	// set parallel to true
	subNode.Parallel = true
	// link as serial
	n.linkSubSerialNode(subNode)
}

func (n *Node) linkSubSerialNode(subNode *Node) {
	// find index that subNode should be put into
	subNodeIndex := -1
	for i, nextNode := range n.NextNodes {
		if nextNode.Name == subNode.Name {
			subNodeIndex = i
			break
		}
	}
	// not found, append to end
	if subNodeIndex == -1 {
		// first drop subNode's original link
		subNode.cutOffPrevious()
		// then add to node's nextNodes
		n.addNext(subNode)
	}
}

func parseParallelRendering(p *cptype.ComponentProtocol, compRenderingItems []cptype.RendingItem) (*Node, error) {
	if len(compRenderingItems) == 0 {
		return nil, nil
	}

	// link all nodes according to compRenderingItem's serial-order
	nodesMap := make(map[string]*Node)
	var rootNode *Node
	var lastNode *Node
	for _, item := range compRenderingItems {
		// make new serial node
		node := makeSerialNode(item)
		// add to nodes map
		nodesMap[node.Name] = node
		// set root node
		if lastNode == nil {
			rootNode = node
		} else {
			// link node with previous
			lastNode.addNext(node)
		}
		// set current node as lastNode
		lastNode = node
		continue
	}

	// link again according to hierarchy structure
	for nodeName, v := range p.Hierarchy.Structure {
		var subCompNames []string
		if err := cputil.ObjJSONTransfer(&v, &subCompNames); err != nil {
			continue
		}
		node := nodesMap[nodeName]
		for _, subNodeName := range subCompNames {
			// set subNode's previous again
			subNode := nodesMap[subNodeName]
			node.linkSubSerialNode(subNode)
		}
	}

	// link all nodes again according to hierarchy.Parallel definition
	parallelDef := p.Hierarchy.Parallel
	if parallelDef == nil {
		return rootNode, nil
	}
	for parentNodeName, subParallelNodeNames := range parallelDef {
		// check firstly
		parentNode, ok := nodesMap[parentNodeName]
		if !ok {
			return nil, fmt.Errorf("invalid parallel definition, %s not exist", parentNodeName)
		}
		for _, subNodeName := range subParallelNodeNames {
			// check firstly
			subNode, ok := nodesMap[subNodeName]
			if !ok {
				return nil, fmt.Errorf("invalid parallel definition, %s not exist", subNodeName)
			}
			// link parent and sub node
			parentNode.linkSubParallelNode(subNode)
		}
	}

	return rootNode, nil
}

func renderFromNode(ctx context.Context, req *cptype.ComponentProtocolRequest, sr ScenarioRender, node *Node) error {
	// render itself
	if err := renderOneNode(ctx, req, sr, node); err != nil {
		return err
	}

	// continue render until done
	i := 0
	for {
		if i > 50 {
			return fmt.Errorf("abnormal render next nodes, over 50 times, force stop")
		}
		if len(node.doneNextNodesByName) == len(node.NextNodes) {
			break
		}
		// render next nodes
		if err := node.renderNextNodes(ctx, req, sr); err != nil {
			return err
		}
		i++
	}

	return nil
}

func (n *Node) renderNextNodes(ctx context.Context, req *cptype.ComponentProtocolRequest, sr ScenarioRender) error {
	// render next nodes
	renderableNodes := n.calcRenderableNextNodes()
	if len(renderableNodes) == 0 {
		return nil
	}
	printRenderableNodes(renderableNodes)
	var wg sync.WaitGroup
	var errorMsgs []string
	for _, nextNode := range renderableNodes {
		wg.Add(1)
		go func(nextNode *Node) {
			logrus.Infof("begin render node: %s", nextNode.Name)
			defer logrus.Infof("end render node: %s", nextNode.Name)
			defer wg.Done()

			if err := renderFromNode(ctx, req, sr, nextNode); err != nil {
				errorMsgs = append(errorMsgs, err.Error())
			}
		}(nextNode)
	}
	wg.Wait()
	if len(errorMsgs) > 0 {
		return fmt.Errorf(strings.Join(errorMsgs, ", "))
	}
	return nil
}

func printRenderableNodes(nodes []*Node) {
	var nodeNames []string
	for _, node := range nodes {
		nodeNames = append(nodeNames, node.Name)
	}
	switch len(nodeNames) {
	case 0:
		return
	case 1:
		logrus.Infof("[S] serial renderable node: %s", strings.Join(nodeNames, ", "))
	default:
		logrus.Infof("[P] parallel renderable nodes: %s", strings.Join(nodeNames, ", "))
	}
}

func (n *Node) calcRenderableNextNodes() []*Node {
	if n.doneNextNodesByName == nil {
		n.doneNextNodesByName = make(map[string]*Node)
	}
	var renderableNodes []*Node
	defer func() {
		for _, node := range renderableNodes {
			n.doneNextNodesByName[node.Name] = node
		}
	}()
	// get from nextNodes by order
	for _, next := range n.NextNodes {
		// skip already done
		if _, done := n.doneNextNodesByName[next.Name]; done {
			continue
		}
		// add if empty
		if len(renderableNodes) == 0 {
			renderableNodes = append(renderableNodes, next)
			continue
		}

		// if first is serial, stop until next serial node come(exclude)
		// s->p->s => s->p
		// s->s    => s
		// p->p->s => p->p
		// p->s    => p
		if !next.Parallel {
			return renderableNodes
		}

		// add
		renderableNodes = append(renderableNodes, next)
	}
	return renderableNodes
}

func renderOneNode(ctx context.Context, req *cptype.ComponentProtocolRequest, sr ScenarioRender, node *Node) error {
	return renderOneComp(ctx, req, sr, node.toRendingItem())
}

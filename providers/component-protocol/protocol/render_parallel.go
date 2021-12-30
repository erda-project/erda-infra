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
	"fmt"
	"io"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

type Node struct {
	Name     string
	Parallel bool

	NextNodes       []*Node
	nextNodesByName map[string]*Node
	PreviousNode    *Node
}

func printNewLines(w io.Writer, repeat int) {
	if repeat == 0 {
		repeat = 1
	}
	for i := 0; i < repeat; i++ {
		fmt.Fprintln(w)
	}
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
	fmt.Fprintf(w, "name: %s\n", n.Name)
	depth := 1
	n.printNexts(w, depth)
	return w.String()
}
func (n *Node) printNexts(w io.Writer, depth int) {
	for i, next := range n.NextNodes {
		printIndent(w, depth*i*2)
		printNode(w, next)
		next.printNexts(w, depth+1)
	}
}

func makeSerialNode(item cptype.RendingItem) *Node {
	return &Node{Name: item.Name, Parallel: false}
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

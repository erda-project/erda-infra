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

package graph

// reference http://dnaeon.github.io/dependency-graph-resolution-algorithm-in-go/

import (
	"errors"
	"fmt"
)

// Node represents a single node in the graph with it's dependencies
type Node struct {
	// Name of the node
	Name string

	// Dependencies of the node
	Deps []string
}

func (n *Node) String() string {
	if len(n.Deps) > 0 {
		return fmt.Sprintf("%s -> %s", n.Name, n.Deps)
	}
	return n.Name
}

// NewNode creates a new node
func NewNode(name string, deps ...string) *Node {
	n := &Node{
		Name: name,
		Deps: deps,
	}
	return n
}

// Graph dependency graph
type Graph []*Node

// Display the dependency graph
func (g Graph) Display() {
	for _, node := range g {
		if len(node.Deps) <= 0 {
			fmt.Println(node.Name)
		} else {
			for _, dep := range node.Deps {
				fmt.Printf("%s -> %s\n", node.Name, dep)
			}
		}
	}
}

// Resolve the dependency graph
func Resolve(graph Graph) (Graph, error) {
	// A map containing the node names and the actual node object
	nodeNames := make(map[string]*Node)

	// A map containing the nodes and their dependencies
	nodeDependencies := make(map[string]map[string]struct{})

	// Populate the maps
	for _, node := range graph {
		nodeNames[node.Name] = node

		dependencySet := make(map[string]struct{})
		for _, dep := range node.Deps {
			dependencySet[dep] = struct{}{}
		}
		nodeDependencies[node.Name] = dependencySet
	}

	// Iteratively find and remove nodes from the graph which have no dependencies.
	// If at some point there are still nodes in the graph and we cannot find
	// nodes without dependencies, that means we have a circular dependency
	var resolved Graph
	for len(nodeDependencies) != 0 {
		// Get all nodes from the graph which have no dependencies
		readySet := make(map[string]struct{})
		for name, deps := range nodeDependencies {
			if len(deps) == 0 {
				readySet[name] = struct{}{}
			}
		}

		// If there aren't any ready nodes, then we have a cicular dependency
		if len(readySet) == 0 {
			var g Graph
			for name := range nodeDependencies {
				g = append(g, nodeNames[name])
			}
			return g, errors.New("Circular dependency found")
		}

		// Remove the ready nodes and add them to the resolved graph
		for name := range readySet {
			delete(nodeDependencies, name)
			resolved = append(resolved, nodeNames[name])
		}

		// Also make sure to remove the ready nodes from the
		// remaining node dependencies as well
		for name, deps := range nodeDependencies {
			diff := make(map[string]struct{})
			for dep := range deps {
				if _, ok := readySet[dep]; !ok {
					diff[dep] = struct{}{}
				}
			}
			nodeDependencies[name] = diff
		}
	}

	return resolved, nil
}

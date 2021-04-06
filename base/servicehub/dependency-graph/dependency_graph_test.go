// Copyright 2021 Terminus
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

import "fmt"

func Example_circular() {
	node1 := NewNode("node1", "node2")
	node2 := NewNode("node2", "node3")
	node3 := NewNode("node3", "node1")
	var g Graph
	g = append(g, node1, node2, node3)
	g, err := Resolve(g)
	if err != nil {
		fmt.Println(err)
		// g.Display()
		return
	}
	fmt.Println("OK")
	// Output:
	// Circular dependency found
}

func Example_ok() {
	node1 := NewNode("node1", "node2")
	node2 := NewNode("node2", "node3")
	node3 := NewNode("node3")
	var g Graph
	g = append(g, node1, node2, node3)
	g, err := Resolve(g)
	if err != nil {
		fmt.Println(err)
		g.Display()
		return
	}
	fmt.Println("OK")
	g.Display()
	// Output:
	// OK
	// node3
	// node2 -> node3
	// node1 -> node2
}

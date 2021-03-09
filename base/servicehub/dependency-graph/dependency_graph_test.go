// Author: recallsong
// Email: songruiguo@qq.com

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

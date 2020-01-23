package main

import (
	"fmt"
	"github.com/heimdalr/dag"
)

// data structure that will be used as vertex in the graph
type myVertex struct {
	value int
}

// implement the Vertex's interface method String()
func (v myVertex) String() string {
	return fmt.Sprintf("%d", v.value)
}

// implement the Vertex's interface method Id()
func (v myVertex) Id() string {
	return fmt.Sprintf("%d", v.value)
}

func main() {

	// initialize a new graph
	d := dag.NewDAG()

	// init three vertices
	v1 := &myVertex{1}
	v2 := &myVertex{2}
	v3 := &myVertex{3}

	// add the above vertices and connect them with two edges
	_ = d.AddEdge(v1, v2)
	_ = d.AddEdge(v1, v3)

	// describe the graph
	fmt.Print(d.String())
}

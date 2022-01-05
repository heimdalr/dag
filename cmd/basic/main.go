package main

import (
	"fmt"
	"github.com/heimdalr/dag"
)

func main() {

	// initialize a new graph
	d := dag.NewDAG()

	// init three vertices
	v1, _ := d.AddVertex(1)
	v2, _ := d.AddVertex(2)
	v3, _ := d.AddVertex(3)

	// add the above vertices and connect them with two edges
	_ = d.AddEdge(v1, v2)
	_ = d.AddEdge(v1, v3)

	// describe the graph
	fmt.Print(d.String())
}

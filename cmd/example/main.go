package main

import (
	"fmt"
	"github.com/heimdalr/dag"
	"os"
)

// data structure that will be used as vertex in the graph
type myVertex struct {
	Label string
}

// implement the Vertex interface
func (v myVertex) String() string {
	return v.Label
}

func main() {

	// initialize a new graph
	d := dag.NewDAG()

	os.Getenv("goo")

	// init three vertices
	v1 := &myVertex{"1"}
	v2 := &myVertex{"2"}
	v3 := &myVertex{"3"}

	// add the above vertices and connect them with two edges
	_ = d.AddEdge(v1, v2)
	_ = d.AddEdge(v1, v3)

	// describe the graph
	fmt.Print(d.String())
}
package dag_test

import (
	"fmt"
	"github.com/heimdalr/dag"
)

type foobar struct {
	a string
	b string
}

func Example() {

	// initialize a new graph
	d := dag.NewDAG()

	// init three vertices
	v1, _ := d.AddVertex(1)
	v2, _ := d.AddVertex(2)
	v3, _ := d.AddVertex(foobar{a: "foo", b: "bar"})

	// add the above vertices and connect them with two edges
	_ = d.AddEdge(v1, v2)
	_ = d.AddEdge(v1, v3)

	// describe the graph
	fmt.Print(d.String())

	// Unordered output:
	// DAG Vertices: 3 - Edges: 2
	// Vertices:
	//   1
	//   2
	//   {foo bar}
	// Edges:
	//   1 -> 2
	//   1 -> {foo bar}
}

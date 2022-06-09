package main

import (
	"fmt"
	"github.com/heimdalr/dag"
)

func main() {

	// initialize a new graph
	d := dag.NewDAG()

	// init three vertices
	v0, _ := d.AddVertex(0)
	v1, _ := d.AddVertex(1)
	v2, _ := d.AddVertex(2)
	v3, _ := d.AddVertex(3)
	v4, _ := d.AddVertex(4)

	// add the above vertices and connect them with two edges
	_ = d.AddEdge(v0, v1)
	_ = d.AddEdge(v0, v3)
	_ = d.AddEdge(v1, v2)
	_ = d.AddEdge(v2, v4)
	_ = d.AddEdge(v3, v4)

	// worker function
	flowCallback := func(d *dag.DAG, id string, parentResults []dag.FlowResult) interface{} {

		v, _ := d.GetVertex(id)
		fmt.Printf("%v based on:\n", v)
		for _, r := range parentResults {
			p, _ := d.GetVertex(r.ID)
			fmt.Printf("- %v\n", p)
		}
		return nil
	}

	d.DescendantsFlow(v0, nil, flowCallback)

}

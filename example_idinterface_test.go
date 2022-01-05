package dag_test

import (
	"fmt"
	"github.com/heimdalr/dag"
)

type idVertex struct {
	id  string
	msg string
}

func (v idVertex) ID() string {
	return v.id
}

func ExampleIDInterface() {

	// initialize a new graph
	d := dag.NewDAG()

	// init three vertices
	id, _ := d.AddVertex(idVertex{id: "1", msg: "foo"})
	fmt.Printf("id of vertex is %s\n", id)
	v, _ := d.GetVertex(id)
	fmt.Printf("%s", v)

	// Output:
	// id of vertex is 1
	// {1 foo}
}

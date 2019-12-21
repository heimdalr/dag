package dag

import (
	"testing"
)

type testVertex struct {
	Label string
}

func (v testVertex) String() string {
	return v.Label
}

func TestDAG(t *testing.T) {
	dag := NewDAG()
	if dag.GetOrder() != 0 {
		t.Fatalf("DAG number of vertices expected to be 0 but got %dag", dag.GetOrder())
	}
}

func Test_AddVertex(t *testing.T) {
	dag := NewDAG()
	var v Vertex = testVertex{"1"}
	err := dag.AddVertex(&v)
	if err != nil {
		t.Fatalf("Can't add vertex to DAG: %s", err)
	}
	if dag.GetOrder() != 1 {
		t.Fatalf("DAG number of vertices expected to be 1 but got %d", dag.GetOrder())
	}
	err2 := dag.AddVertex(&v)
	if err2 == nil {
		t.Fatal("Expected to see a duplicate entry error")
	}
	var v2 Vertex = testVertex{"2"}
	err3 := dag.AddVertex(&v2)
	if err3 != nil {
		t.Fatal("Did not expect to see a duplicate entry error")
	}

}

func Test_DeleteVertex(t *testing.T) {
	dag := NewDAG()
	var v Vertex = testVertex{"1"}
	_ = dag.AddVertex(&v)
	dag.DeleteVertex(&v)
	if dag.GetOrder() != 0 {
		t.Fatalf("DAG number of vertices expected to be 0 but got %d", dag.GetOrder())
	}
	dag.DeleteVertex(&v)
}

func Test_AddEdge(t *testing.T) {
	dag := NewDAG()
	var v1 Vertex = testVertex{"1"}
	var v2 Vertex = testVertex{"2"}
	_ = dag.AddVertex(&v1)
	_ = dag.AddVertex(&v2)
	errEdge2 := dag.AddEdge(&v1, &v2)
	if errEdge2 != nil {
		t.Fatalf("Can't add edge to DAG: %s", errEdge2)
	}
}

func Test_Ancestors(t *testing.T) {
	dag := NewDAG()
	var v1 Vertex = testVertex{"1"}
	var v2 Vertex = testVertex{"2"}
	var v3 Vertex = testVertex{"3"}
	var v4 Vertex = testVertex{"4"}
	_ = dag.AddVertex(&v1)
	_ = dag.AddVertex(&v2)
	_ = dag.AddVertex(&v3)
	_ = dag.AddVertex(&v4)
	_ = dag.AddEdge(&v1, &v2)
	_ = dag.AddEdge(&v2, &v3)
	_ = dag.AddEdge(&v2, &v4)
	//fmt.Print(dag)
	ancestors, err := dag.GetAncestors(&v4)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if len(ancestors) != 2 {
		t.Fatalf("DAG number of getAncestorsAux expected to be 2 but got %d", len(ancestors))
	}
}

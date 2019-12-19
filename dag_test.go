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
	if dag.Order() != 0 {
		t.Fatalf("DAG number of vertices expected to be 0 but got %dag", dag.Order())
	}
}

func TestDAG_AddVertex(t *testing.T) {
	dag := NewDAG()
	var v Vertex = testVertex{"1"}
	err := dag.AddVertex("1", &v)
	if err != nil {
		t.Fatalf("Can't add vertex to DAG: %s", err)
	}
	if dag.Order() != 1 {
		t.Fatalf("DAG number of vertices expected to be 1 but got %d", dag.Order())
	}
	err2 := dag.AddVertex("1", &v)
	if err2 == nil {
		t.Fatal("Expected to see a duplicate entry error")
	}
	err3 := dag.AddVertex("2", &v)
	if err3 != nil {
		t.Fatal("Did not expect to see a duplicate entry error")
	}

}

func TestDAG_DeleteVertex(t *testing.T) {
	dag := NewDAG()
	var v Vertex = testVertex{"1"}
	_ = dag.AddVertex("1", &v)
	dag.DeleteVertex("1")
	if dag.Order() != 0 {
		t.Fatalf("DAG number of vertices expected to be 0 but got %d", dag.Order())
	}
	dag.DeleteVertex("1")
}

func TestDAG_AddEdge(t *testing.T) {
	dag := NewDAG()
	var v Vertex = testVertex{"1"}
	_ = dag.AddVertex("1", &v)
	errEdge1 := dag.AddEdge("1", "2")
	if errEdge1 == nil {
		t.Fatal("Expected to see a missing vertex error")
	}
	_ = dag.AddVertex("2", &v)
	errEdge2 := dag.AddEdge("1", "2")
	if errEdge2 != nil {
		t.Fatalf("Can't add edge to DAG: %s", errEdge2)
	}
}

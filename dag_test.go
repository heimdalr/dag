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

func makeVertex(label string) *Vertex {
	var v Vertex = testVertex{label}
	return &v
}

func TestNewDAG(t *testing.T) {
	dag := NewDAG()
	if order := dag.GetOrder(); order != 0 {
		t.Errorf("GetOrder() = %d, want 0", order)
	}
	if size := dag.GetSize(); size != 0 {
		t.Errorf("GetSize() = %d, want 0", size)
	}
}

func TestAddVertex(t *testing.T) {
	dag := NewDAG()
	v := makeVertex("1")
	err := dag.AddVertex(v)
	if err != nil {
		t.Error(err)
	}
	if order := dag.GetOrder(); order != 1 {
		t.Errorf("GetOrder() = %d, want 1", order)
	}
	if size := dag.GetSize(); size != 0 {
		t.Errorf("GetSize() = %d, want 0", size)
	}
	if leafs := len(dag.GetLeafs()); leafs != 1 {
		t.Errorf("GetLeafs() = %d, want 1", leafs)
	}
	if roots := len(dag.GetRoots()); roots != 1 {
		t.Errorf("GetLeafs() = %d, want 1", roots)
	}
	if vertices := len(dag.GetVertices()); vertices != 1 {
		t.Errorf("GetVertices() = %d, want 1", vertices)
	}
	if ptr := dag.GetVertices()[0]; ptr != v {
		t.Errorf("GetVertices()[0] = %p, want %p", ptr, &v)
	}
}

func TestDeleteVertex(t *testing.T) {
	dag := NewDAG()
	v := makeVertex("1")
	_ = dag.AddVertex(v)
	dag.DeleteVertex(v)
	if order := dag.GetOrder(); order != 0 {
		t.Errorf("GetOrder() = %d, want 0", order)
	}
	if size := dag.GetSize(); size != 0 {
		t.Errorf("GetSize() = %d, want 0", size)
	}
	if leafs := len(dag.GetLeafs()); leafs != 0 {
		t.Errorf("GetLeafs() = %d, want 0", leafs)
	}
	if roots := len(dag.GetRoots()); roots != 0 {
		t.Errorf("GetLeafs() = %d, want 0", roots)
	}
	if vertices := len(dag.GetVertices()); vertices != 0 {
		t.Errorf("GetVertices() = %d, want 0", vertices)
	}
}

func TestAddEdge(t *testing.T) {
	dag := NewDAG()
	src := makeVertex("src")
	dst := makeVertex("dst")
	_ = dag.AddVertex(src)
	_ = dag.AddVertex(dst)
	err := dag.AddEdge(src, dst)
	if err != nil {
		t.Error(err)
	}
	children, errChildren := dag.GetChildren(src)
	if errChildren != nil {
		t.Error(errChildren)
	}
	if length := len(children); length != 1 {
		t.Errorf("GetChildren() = %d, want 1", length)
	}
	parents, errParents := dag.GetParents(dst)
	if errParents != nil {
		t.Error(errParents)
	}
	if length := len(parents); length != 1 {
		t.Errorf("GetParents() = %d, want 1", length)
	}
	if leafs := len(dag.GetLeafs()); leafs != 1 {
		t.Errorf("GetLeafs() = %d, want 1", leafs)
	}
	if roots := len(dag.GetRoots()); roots != 1 {
		t.Errorf("GetLeafs() = %d, want 1", roots)
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


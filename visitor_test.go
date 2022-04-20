package dag

import (
	"testing"

	"github.com/go-test/deep"
)

type testVisitor struct {
	Values []string
}

func (pv *testVisitor) Visit(v Vertexer) {
	_, value := v.Vertex()
	pv.Values = append(pv.Values, value.(string))
}

func getTestWalkDAG() *DAG {
	dag := NewDAG()
	v1, v2, v3, v4, v5 := "1", "2", "3", "4", "5"
	_ = dag.AddVertexByID(v1, "v1")
	_ = dag.AddVertexByID(v2, "v2")
	_ = dag.AddVertexByID(v3, "v3")
	_ = dag.AddVertexByID(v4, "v4")
	_ = dag.AddVertexByID(v5, "v5")
	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v2, v3)
	_ = dag.AddEdge(v2, v4)
	_ = dag.AddEdge(v4, v5)
	return dag
}

func TestDFSWalk(t *testing.T) {
	dag := getTestWalkDAG()

	pv := &testVisitor{}
	DFSWalk(dag, pv)

	expected := []string{"v1", "v2", "v3", "v4", "v5"}
	actual := pv.Values
	if deep.Equal(expected, actual) != nil {
		t.Errorf("DFSWalk() = %v, want %v", actual, expected)
	}
}

func TestBFSWalk(t *testing.T) {
	dag := getTestWalkDAG()
	pv := &testVisitor{}
	BFSWalk(dag, pv)

	expected := []string{"v1", "v2", "v3", "v4", "v5"}
	actual := pv.Values
	if deep.Equal(expected, actual) != nil {
		t.Errorf("BFSWalk() = %v, want %v", actual, expected)
	}
}

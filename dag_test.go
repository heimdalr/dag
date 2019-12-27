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
	//v2 := testVertex{label}.(Vertex)
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

func TestDAG_AddVertex(t *testing.T) {
	dag := NewDAG()
	dag.AddVertex(nil)
	v := makeVertex("1")
	dag.AddVertex(v)
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
	if !dag.GetVertices()[v] {
		t.Errorf("GetVertices()[v] = %t, want true", dag.GetVertices()[v])
	}
}

func TestDAG_DeleteVertex(t *testing.T) {
	dag := NewDAG()
	v := makeVertex("1")
	dag.AddVertex(v)
	dag.DeleteVertex(nil)
	if order := dag.GetOrder(); order != 1 {
		t.Errorf("GetOrder() = %d, want 1", order)
	}
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

func TestDAG_AddEdge(t *testing.T) {
	dag := NewDAG()
	_ = dag.AddEdge(nil, nil)
	src := makeVertex("src")
	dst := makeVertex("dst")
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
	if err := dag.AddEdge(src, src); err != nil {
		t.Error("AddEdge(x, x) expected to not return an error")
	}
}

func TestDAG_AddEdgeSafe(t *testing.T) {
	dag := NewDAG()
	src := makeVertex("src")
	dst := makeVertex("dst")
	loopErr := dag.AddEdgeSafe(src, src)
	if loopErr == nil {
		t.Error("AddEdgeSafe(src, src) expected error")
	} else {
		if _, ok := loopErr.(LoopError); !ok {
			t.Errorf("AddEdgeSafe(src, src) expected LoopError, got %T", loopErr)
		}
		expectedText := "loop between 'src' and 'src'"
		if text := loopErr.Error(); text != expectedText {
			t.Errorf("AddEdgeSafe(src, src) = \"%s\", want \"%s\"", text, expectedText)
		}
	}
	if err := dag.AddEdgeSafe(src, dst); err != nil {
		t.Errorf("AddEdgeSafe(x, y) unexpected error: %v", err)
	}
	if err := dag.AddEdgeSafe(dst, src); err == nil {
		t.Errorf("AddEdgeSafe(y, x) expected error: %v", err)
	}
}

func TestDAG_DeleteEdge(t *testing.T) {
	dag := NewDAG()
	src := makeVertex("src")
	dst := makeVertex("dst")
	_ = dag.AddEdge(src, dst)
	if size := dag.GetSize(); size != 1 {
		t.Errorf("GetSize() = %d, want 1", size)
	}
	dag.DeleteEdge(src, nil)
	dag.DeleteEdge(src, dst)
	if size := dag.GetSize(); size != 0 {
		t.Errorf("GetSize() = %d, want 0", size)
	}
	dag.DeleteEdge(src, dst)
}

func TestDAG_GetChildren(t *testing.T) {
	dag := NewDAG()
	v1 := makeVertex("1")
	v2 := makeVertex("2")
	v3 := makeVertex("3")
	v4 := makeVertex("4")
	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v2, v3)
	_, errUnknown := dag.GetChildren(v4)
	if errUnknown == nil {
		t.Errorf("GetChildren(v4) expected error")
	} else {
		if _, ok := errUnknown.(VertexUnknownError); !ok {
			t.Errorf("GetChildren(v4) expected VertexUnknownError, got %T", errUnknown)
		}
		expectedText := "'4' is unknown"
		if text := errUnknown.Error(); text != expectedText {
			t.Errorf("GetChildren(v4) = \"%s\", want \"%s\"", errUnknown, expectedText)
		}
	}
	children, errChildren := dag.GetChildren(v1)
	if errChildren != nil {
		t.Error(errChildren)
	}
	if length := len(children); length != 1 {
		t.Errorf("GetChildren() = %d, want 1", length)
	}
	if truth := children[v2]; !truth {
		t.Errorf("GetChildren()[v2] = %t, want true", truth)
	}
	if truth := children[v3]; truth {
		t.Errorf("GetChildren()[v3] = %t, want false", truth)
	}

}

func TestDAG_GetParents(t *testing.T) {
	dag := NewDAG()
	v1 := makeVertex("1")
	_, errUnknown := dag.GetParents(v1)
	if errUnknown == nil {
		t.Errorf("GetParents(v1) expected error")
	} else {
		if _, ok := errUnknown.(VertexUnknownError); !ok {
			t.Errorf("GetParents(v1) expected VertexUnknownError, got %T", errUnknown)
		}
		expectedText := "'1' is unknown"
		if text := errUnknown.Error(); text != expectedText {
			t.Errorf("GetParents(v1) = \"%s\", want \"%s\"", errUnknown, expectedText)
		}
	}
}

func TestDAG_GetDescendants(t *testing.T) {
	dag := NewDAG()
	v1 := makeVertex("1")
	_, errUnknown := dag.GetDescendants(v1)
	if errUnknown == nil {
		t.Errorf("GetDescendants(v1) expected error")
	} else {
		if _, ok := errUnknown.(VertexUnknownError); !ok {
			t.Errorf("GetDescendants(v1) expected VertexUnknownError, got %T", errUnknown)
		}
		expectedText := "'1' is unknown"
		if text := errUnknown.Error(); text != expectedText {
			t.Errorf("GetDescendants(v1) = \"%s\", want \"%s\"", errUnknown, expectedText)
		}
	}
	v2 := makeVertex("2")
	v3 := makeVertex("3")
	v4 := makeVertex("4")
	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v2, v3)
	_ = dag.AddEdge(v2, v4)
	desc, err := dag.GetDescendants(v1)
	if err != nil {
		t.Errorf("GetDescendants(v1) unexpected error: %s", err)
	}
	if len(desc) != 3 {
		t.Errorf("GetDescendants(v1) = %d, want 3", len(desc))
	}
	if desc, _ := dag.GetDescendants(v4); len(desc) != 0 {
		t.Errorf("GetDescendants(v4) = %d, want 0", len(desc))
	}
	if desc, _ := dag.GetDescendants(v2); len(desc) != 2 {
		t.Errorf("GetDescendants(v2) = %d, want 2", len(desc))
	}
}

func TestGetAncestors(t *testing.T) {
	dag := NewDAG()
	v1 := makeVertex("1")
	v2 := makeVertex("2")
	v3 := makeVertex("3")
	v4 := makeVertex("4")
	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v2, v3)
	_ = dag.AddEdge(v2, v4)
	ancestors, err := dag.GetAncestors(v4)
	if err != nil {
		t.Errorf("GetAncestors(v4) unexpected error: %s", err)
	}
	if len(ancestors) != 2 {
		t.Errorf("GetAncestors(v4) = %d, want 2", len(ancestors))
	}
	if ancestors, _ := dag.GetAncestors(v1); len(ancestors) != 0 {
		t.Errorf("GetAncestors(v1) = %d, want 0", len(ancestors))
	}
	if ancestors, _ := dag.GetAncestors(v2); len(ancestors) != 1 {
		t.Errorf("GetAncestors(v2) = %d, want 1", len(ancestors))
	}
	v5 := makeVertex("5")
	if _, err := dag.GetAncestors(v5); err == nil {
		t.Error("GetAncestors(v2) expected to return an error")
	}

}

func TestDAG_String(t *testing.T) {
	dag := NewDAG()
	v1 := makeVertex("1")
	v2 := makeVertex("2")
	v3 := makeVertex("3")
	v4 := makeVertex("4")
	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v2, v3)
	_ = dag.AddEdge(v2, v4)
	expected := "DAG Vertices: 4 - Edges: 3"
	s := dag.String()
	if s[:len(expected)] != expected {
		t.Errorf("String() = \"%s\", want \"%s\"", s, expected)
	}
}

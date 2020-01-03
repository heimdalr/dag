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

	// add a single vertex and inspect the graph
	v := &testVertex{"1"}
	_ = dag.AddVertex(v)
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

	// duplicate
	errDuplicate := dag.AddVertex(v)
	if errDuplicate == nil {
		t.Errorf("AddVertex(v) = nil, want %T", VertexDuplicateError{v})
	}
	if _, ok := errDuplicate.(VertexDuplicateError); !ok {
		t.Errorf("AddVertex(v) expected VertexDuplicateError, got %T", errDuplicate)
	}

	// nil
	errNil:= dag.AddVertex(nil)
	if errNil == nil {
		t.Errorf("AddVertex(nil) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNil.(VertexNilError); !ok {
		t.Errorf("AddVertex(nil) expected VertexNilError, got %T", errNil)
	}


}

func TestDAG_DeleteVertex(t *testing.T) {
	dag := NewDAG()
	v1 := &testVertex{"1"}
	v2 := &testVertex{"2"}
	v3 := &testVertex{"3"}
	_ = dag.AddVertex(v1)

	// delete a single vertex and inspect the graph
	_ = dag.DeleteVertex(v1)
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

	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v2, v3)
	if order := dag.GetOrder(); order != 3 {
		t.Errorf("GetOrder() = %d, want 3", order)
	}
	if size := dag.GetSize(); size != 2 {
		t.Errorf("GetSize() = %d, want 2", size)
	}
	if leafs := len(dag.GetLeafs()); leafs != 1 {
		t.Errorf("GetLeafs() = %d, want 1", leafs)
	}
	if roots := len(dag.GetRoots()); roots != 1 {
		t.Errorf("GetLeafs() = %d, want 1", roots)
	}
	if vertices := len(dag.GetVertices()); vertices != 3 {
		t.Errorf("GetVertices() = %d, want 3", vertices)
	}
	if vertices, _ := dag.GetDescendants(v1); len(vertices) != 2 {
		t.Errorf("GetDescendants(v1) = %d, want 2", len(vertices) )
	}
	if vertices, _ := dag.GetAncestors(v3); len(vertices) != 2 {
		t.Errorf("GetAncestors(v3) = %d, want 2", len(vertices) )
	}

	_ = dag.DeleteVertex(v2)
	if order := dag.GetOrder(); order != 2 {
		t.Errorf("GetOrder() = %d, want 2", order)
	}
	if size := dag.GetSize(); size != 0 {
		t.Errorf("GetSize() = %d, want 0", size)
	}
	if leafs := len(dag.GetLeafs()); leafs != 2 {
		t.Errorf("GetLeafs() = %d, want 2", leafs)
	}
	if roots := len(dag.GetRoots()); roots != 2 {
		t.Errorf("GetLeafs() = %d, want 2", roots)
	}
	if vertices := len(dag.GetVertices()); vertices != 2 {
		t.Errorf("GetVertices() = %d, want 2", vertices)
	}
	if vertices, _ := dag.GetDescendants(v1); len(vertices) != 0 {
		t.Errorf("GetDescendants(v1) = %d, want 0", len(vertices) )
	}
	if vertices, _ := dag.GetAncestors(v3); len(vertices) != 0 {
		t.Errorf("GetAncestors(v3) = %d, want 0", len(vertices) )
	}


	// unknown
	foo := &testVertex{"foo"}
	errUnknown := dag.DeleteVertex(foo)
	if errUnknown == nil {
		t.Errorf("DeleteVertex(foo) = nil, want %T", VertexUnknownError{foo})
	}
	if _, ok := errUnknown.(VertexUnknownError); !ok {
		t.Errorf("DeleteVertex(foo) expected VertexUnknownError, got %T", errUnknown)
	}

	// nil
	errNil:= dag.DeleteVertex(nil)
	if errNil == nil {
		t.Errorf("DeleteVertex(nil) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNil.(VertexNilError); !ok {
		t.Errorf("DeleteVertex(nil) expected VertexNilError, got %T", errNil)
	}
}

func TestDAG_AddEdge(t *testing.T) {
	dag := NewDAG()
	v0 := &testVertex{"v0"}
	v1 := &testVertex{"v1"}
	v2 := &testVertex{"v2"}
	v3 := &testVertex{"v3"}

	// add a single edge and inspect the graph
	_ = dag.AddEdge(v1, v2)
	if children, _ := dag.GetChildren(v1); len(children) != 1 {
		t.Errorf("GetChildren(v1) = %d, want 1", len(children))
	}
	if parents, _ := dag.GetParents(v2); len(parents) != 1 {
		t.Errorf("GetParents(v2) = %d, want 1", len(parents))
	}
	if leafs := len(dag.GetLeafs()); leafs != 1 {
		t.Errorf("GetLeafs() = %d, want 1", leafs)
	}
	if roots := len(dag.GetRoots()); roots != 1 {
		t.Errorf("GetLeafs() = %d, want 1", roots)
	}
	if vertices, _ := dag.GetDescendants(v1); len(vertices) != 1 {
		t.Errorf("GetDescendants(v1) = %d, want 1", len(vertices) )
	}
	if vertices, _ := dag.GetAncestors(v2); len(vertices) != 1 {
		t.Errorf("GetAncestors(v2) = %d, want 1", len(vertices) )
	}

	err := dag.AddEdge(v2, v3)
	if err != nil {
		t.Fatal(err)
	}
	if vertices, _ := dag.GetDescendants(v1); len(vertices) != 2 {
		t.Errorf("GetDescendants(v1) = %d, want 2", len(vertices) )
	}
	if vertices, _ := dag.GetAncestors(v3); len(vertices) != 2 {
		t.Errorf("GetAncestors(v3) = %d, want 2", len(vertices) )
	}

	_ = dag.AddEdge(v0, v1)
	if vertices, _ := dag.GetDescendants(v0); len(vertices) != 3 {
		t.Errorf("GetDescendants(v0) = %d, want 3", len(vertices) )
	}
	if vertices, _ := dag.GetAncestors(v3); len(vertices) != 3 {
		t.Errorf("GetAncestors(v3) = %d, want 3", len(vertices) )
	}

	// loop
	errLoopSrcSrc := dag.AddEdge(v1, v1)
	if errLoopSrcSrc == nil {
		t.Errorf("AddEdge(v1, v1) = nil, want %T", EdgeLoopError{v1, v1})
	}
	if _, ok := errLoopSrcSrc.(EdgeLoopError); !ok {
		t.Errorf("AddEdge(v1, v1) expected EdgeLoopError, got %T", errLoopSrcSrc)
	}
	errLoopDstSrc := dag.AddEdge(v2, v1)
	if errLoopDstSrc == nil {
		t.Errorf("AddEdge(v2, v1) = nil, want %T", EdgeLoopError{v2, v1})
	}
	if _, ok := errLoopDstSrc.(EdgeLoopError); !ok {
		t.Errorf("AddEdge(v2, v1) expected EdgeLoopError, got %T", errLoopDstSrc)
	}

	// duplicate
	errDuplicate := dag.AddEdge(v1, v2)
	if errDuplicate == nil {
		t.Errorf("AddEdge(v1, v2) = nil, want %T", EdgeDuplicateError{v1, v2})
	}
	if _, ok := errDuplicate.(EdgeDuplicateError); !ok {
		t.Errorf("AddEdge(v1, v2) expected EdgeDuplicateError, got %T", errDuplicate)
	}

	// nil
	errNilSrc := dag.AddEdge(nil, v2)
	if errNilSrc == nil {
		t.Errorf("AddEdge(nil, v2) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNilSrc.(VertexNilError); !ok {
		t.Errorf("AddEdge(nil, v2) expected VertexNilError, got %T", errNilSrc)
	}
	errNilDst := dag.AddEdge(v1, nil)
	if errNilDst == nil {
		t.Errorf("AddEdge(v1, nil) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNilDst.(VertexNilError); !ok {
		t.Errorf("AddEdge(v1, nil) expected VertexNilError, got %T", errNilDst)
	}
}

func TestDAG_DeleteEdge(t *testing.T) {
	dag := NewDAG()
	src := &testVertex{"src"}
	dst := &testVertex{"dst"}
	_ = dag.AddEdge(src, dst)
	if size := dag.GetSize(); size != 1 {
		t.Errorf("GetSize() = %d, want 1", size)
	}
	_ = dag.DeleteEdge(src, dst)
	if size := dag.GetSize(); size != 0 {
		t.Errorf("GetSize() = %d, want 0", size)
	}

	// unknown
	errUnknown := dag.DeleteEdge(src, dst)
	if errUnknown == nil {
		t.Errorf("DeleteEdge(src, dst) = nil, want %T", EdgeUnknownError{})
	}
	if _, ok := errUnknown.(EdgeUnknownError); !ok {
		t.Errorf("DeleteEdge(src, dst) expected EdgeUnknownError, got %T", errUnknown)
	}

	// nil
	errNilSrc := dag.DeleteEdge(nil, dst)
	if errNilSrc == nil {
		t.Errorf("DeleteEdge(nil, dst) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNilSrc.(VertexNilError); !ok {
		t.Errorf("DeleteEdge(nil, dst) expected VertexNilError, got %T", errNilSrc)
	}
	errNilDst := dag.DeleteEdge(src, nil)
	if errNilDst == nil {
		t.Errorf("DeleteEdge(src, nil) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNilDst.(VertexNilError); !ok {
		t.Errorf("DeleteEdge(src, nil) expected VertexNilError, got %T", errNilDst)
	}
	
	// unknown
	foo := &testVertex{"foo"}
	errUnknownSrc := dag.DeleteEdge(foo, dst)
	if errUnknownSrc == nil {
		t.Errorf("DeleteEdge(foo, dst) = nil, want %T", VertexUnknownError{})
	}
	if _, ok := errUnknownSrc.(VertexUnknownError); !ok {
		t.Errorf("DeleteEdge(foo, dst) expected VertexUnknownError, got %T", errUnknownSrc)
	}
	errUnknownDst := dag.DeleteEdge(src, foo)
	if errUnknownDst == nil {
		t.Errorf("DeleteEdge(src, foo) = nil, want %T", VertexUnknownError{})
	}
	if _, ok := errUnknownDst.(VertexUnknownError); !ok {
		t.Errorf("DeleteEdge(src, foo) expected VertexUnknownError, got %T", errUnknownDst)
	}
}

func TestDAG_GetChildren(t *testing.T) {
	dag := NewDAG()
	v1 := &testVertex{"1"}
	v2 := &testVertex{"2"}
	v3 := &testVertex{"3"}
	v4 := &testVertex{"4"}
	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v1, v3)

	children, _ := dag.GetChildren(v1)
	if length := len(children); length != 2 {
		t.Errorf("GetChildren() = %d, want 2", length)
	}
	if truth := children[v2]; !truth {
		t.Errorf("GetChildren()[v2] = %t, want true", truth)
	}
	if truth := children[v3]; !truth {
		t.Errorf("GetChildren()[v3] = %t, want true", truth)
	}

	// nil
	_, errNil:= dag.GetChildren(nil)
	if errNil == nil {
		t.Errorf("GetChildren(nil) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNil.(VertexNilError); !ok {
		t.Errorf("GetChildren(nil) expected VertexNilError, got %T", errNil)
	}

	// unknown
	_, errUnknown := dag.GetChildren(v4)
	if errUnknown == nil {
		t.Errorf("GetChildren(v4) = nil, want %T", VertexUnknownError{v4})
	}
	if _, ok := errUnknown.(VertexUnknownError); !ok {
		t.Errorf("GetChildren(v4) expected VertexUnknownError, got %T", errUnknown)
	}
}

func TestDAG_GetParents(t *testing.T) {
	dag := NewDAG()
	v1 := &testVertex{"1"}
	v2 := &testVertex{"2"}
	v3 := &testVertex{"3"}
	v4 := &testVertex{"4"}
	_ = dag.AddEdge(v1, v3)
	_ = dag.AddEdge(v2, v3)

	parents, _ := dag.GetParents(v3)
	if length := len(parents); length != 2 {
		t.Errorf("GetParents(v3) = %d, want 2", length)
	}
	if truth := parents[v1]; !truth {
		t.Errorf("GetParents(v3)[v1] = %t, want true", truth)
	}
	if truth := parents[v2]; !truth {
		t.Errorf("GetParents(v3)[v2] = %t, want true", truth)
	}

	// nil
	_, errNil:= dag.GetParents(nil)
	if errNil == nil {
		t.Errorf("GetParents(nil) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNil.(VertexNilError); !ok {
		t.Errorf("GetParents(nil) expected VertexNilError, got %T", errNil)
	}

	// unknown
	_, errUnknown := dag.GetParents(v4)
	if errUnknown == nil {
		t.Errorf("GetParents(v4) = nil, want %T", VertexUnknownError{v4})
	}
	if _, ok := errUnknown.(VertexUnknownError); !ok {
		t.Errorf("GetParents(v4) expected VertexUnknownError, got %T", errUnknown)
	}

}

func TestDAG_GetDescendants(t *testing.T) {
	dag := NewDAG()
	v1 := &testVertex{"1"}
	v2 := &testVertex{"2"}
	v3 := &testVertex{"3"}
	v4 := &testVertex{"4"}
	v5 := &testVertex{"5"}
	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v2, v3)
	_ = dag.AddEdge(v2, v4)

	if desc, _ := dag.GetDescendants(v1); len(desc) != 3 {
		t.Errorf("GetDescendants(v1) = %d, want 3", len(desc))
	}
	if desc, _ := dag.GetDescendants(v2); len(desc) != 2 {
		t.Errorf("GetDescendants(v2) = %d, want 2", len(desc))
	}
	if desc, _ := dag.GetDescendants(v3); len(desc) != 0 {
		t.Errorf("GetDescendants(v4) = %d, want 0", len(desc))
	}
	if desc, _ := dag.GetDescendants(v4); len(desc) != 0 {
		t.Errorf("GetDescendants(v4) = %d, want 0", len(desc))
	}

	// nil
	_, errNil:= dag.GetDescendants(nil)
	if errNil == nil {
		t.Errorf("GetDescendants(nil) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNil.(VertexNilError); !ok {
		t.Errorf("GetDescendants(nil) expected VertexNilError, got %T", errNil)
	}

	// unknown
	_, errUnknown := dag.GetDescendants(v5)
	if errUnknown == nil {
		t.Errorf("GetDescendants(v5) = nil, want %T", VertexUnknownError{v5})
	}
	if _, ok := errUnknown.(VertexUnknownError); !ok {
		t.Errorf("GetDescendants(v5) expected VertexUnknownError, got %T", errUnknown)
	}
}

func TestDAG_GetAncestors(t *testing.T) {
	dag := NewDAG()
	v0 := &testVertex{"0"}
	v1 := &testVertex{"1"}
	v2 := &testVertex{"2"}
	v3 := &testVertex{"3"}
	v4 := &testVertex{"4"}
	v5 := &testVertex{"5"}
	v6 := &testVertex{"6"}
	v7 := &testVertex{"7"}

	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v2, v3)
	_ = dag.AddEdge(v2, v4)
	_ = dag.AddVertex(v0)
	_ = dag.AddVertex(v5)
	_ = dag.AddVertex(v6)
	_ = dag.AddVertex(v7)

	if ancestors, _ := dag.GetAncestors(v4); len(ancestors) != 2 {
		t.Errorf("GetAncestors(v4) = %d, want 2", len(ancestors))
	}
	if ancestors, _ := dag.GetAncestors(v3); len(ancestors) != 2 {
		t.Errorf("GetAncestors(v3) = %d, want 2", len(ancestors))
	}
	if ancestors, _ := dag.GetAncestors(v2); len(ancestors) != 1 {
		t.Errorf("GetAncestors(v2) = %d, want 1", len(ancestors))
	}
	if ancestors, _ := dag.GetAncestors(v1); len(ancestors) != 0 {
		t.Errorf("GetAncestors(v1) = %d, want 0", len(ancestors))
	}

	_ = dag.AddEdge(v3, v5)
	_ = dag.AddEdge(v4, v6)

	if ancestors, _ := dag.GetAncestors(v4); len(ancestors) != 2 {
		t.Errorf("GetAncestors(v4) = %d, want 2", len(ancestors))
	}
	if ancestors, _ := dag.GetAncestors(v7); len(ancestors) != 0 {
		t.Errorf("GetAncestors(v4) = %d, want 7", len(ancestors))
	}
	_ = dag.AddEdge(v5, v7)
	if ancestors, _ := dag.GetAncestors(v7); len(ancestors) != 4 {
		t.Errorf("GetAncestors(v7) = %d, want 4", len(ancestors))
	}
	_ = dag.AddEdge(v0, v1)
	if ancestors, _ := dag.GetAncestors(v7); len(ancestors) != 5 {
		t.Errorf("GetAncestors(v7) = %d, want 5", len(ancestors))
	}

	// nil
	_, errNil:= dag.GetAncestors(nil)
	if errNil == nil {
		t.Errorf("GetAncestors(nil) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNil.(VertexNilError); !ok {
		t.Errorf("GetAncestors(nil) expected VertexNilError, got %T", errNil)
	}

	// unknown
	foo := &testVertex{"foo"}
	_, errUnknown := dag.GetAncestors(foo)
	if errUnknown == nil {
		t.Errorf("GetAncestors(foo) = nil, want %T", VertexUnknownError{foo})
	}
	if _, ok := errUnknown.(VertexUnknownError); !ok {
		t.Errorf("GetAncestors(foo) expected VertexUnknownError, got %T", errUnknown)
	}

}

func TestDAG_String(t *testing.T) {
	dag := NewDAG()
	v1 := &testVertex{"1"}
	v2 := &testVertex{"2"}
	v3 := &testVertex{"3"}
	v4 := &testVertex{"4"}
	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v2, v3)
	_ = dag.AddEdge(v2, v4)
	expected := "DAG Vertices: 4 - Edges: 3"
	s := dag.String()
	if s[:len(expected)] != expected {
		t.Errorf("String() = \"%s\", want \"%s\"", s, expected)
	}
}


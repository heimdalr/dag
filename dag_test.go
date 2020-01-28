package dag

import (
	"fmt"
	"testing"
)

type iVertex struct {
	value int
}

// implement the Vertex's interface method String()
func (v iVertex) String() string {
	return fmt.Sprintf("%d", v.value)
}

// implement the Vertex's interface method Id()
func (v iVertex) Id() string {
	return fmt.Sprintf("%d", v.value)
}

type sVertex struct {
	value string
}

// implement the Vertex's interface method String()
func (v sVertex) String() string {
	return fmt.Sprintf("%s", v.value)
}

// implement the Vertex's interface method Id()
func (v sVertex) Id() string {
	return fmt.Sprintf("%s", v.value)
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
	v := &iVertex{1}
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
	errNil := dag.AddVertex(nil)
	if errNil == nil {
		t.Errorf("AddVertex(nil) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNil.(VertexNilError); !ok {
		t.Errorf("AddVertex(nil) expected VertexNilError, got %T", errNil)
	}

}

func TestDAG_GetVertex(t *testing.T) {
	dag := NewDAG()
	v1 := &iVertex{1}
	id := v1.String()
	_ = dag.AddVertex(v1)
	if v, _ := dag.GetVertex(id); v != v1 {
		t.Errorf("GetVertex() = %s, want %s", v.String(), v1.String())
	}

	// unknown
	_, errUnknown := dag.GetVertex("foo")
	if errUnknown == nil {
		t.Errorf("DeleteVertex(\"foo\") = nil, want %T", IdUnknownError{"foo"})
	}
	if _, ok := errUnknown.(IdUnknownError); !ok {
		t.Errorf("DeleteVertex(\"foo\") expected IdUnknownError, got %T", errUnknown)
	}

	// nil
	_, errNil := dag.GetVertex("")
	if errNil == nil {
		t.Errorf("DeleteVertex(\"\") = nil, want %T", IdEmptyError{})
	}
	if _, ok := errNil.(IdEmptyError); !ok {
		t.Errorf("DeleteVertex(\"\") expected IdEmptyError, got %T", errNil)
	}
}

func TestDAG_DeleteVertex(t *testing.T) {
	dag := NewDAG()
	v1 := &iVertex{1}
	v2 := &iVertex{2}
	v3 := &iVertex{3}
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
		t.Errorf("GetDescendants(v1) = %d, want 2", len(vertices))
	}
	if vertices, _ := dag.GetAncestors(v3); len(vertices) != 2 {
		t.Errorf("GetAncestors(v3) = %d, want 2", len(vertices))
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
		t.Errorf("GetDescendants(v1) = %d, want 0", len(vertices))
	}
	if vertices, _ := dag.GetAncestors(v3); len(vertices) != 0 {
		t.Errorf("GetAncestors(v3) = %d, want 0", len(vertices))
	}

	// unknown
	foo := &iVertex{-1}
	errUnknown := dag.DeleteVertex(foo)
	if errUnknown == nil {
		t.Errorf("DeleteVertex(foo) = nil, want %T", VertexUnknownError{foo})
	}
	if _, ok := errUnknown.(VertexUnknownError); !ok {
		t.Errorf("DeleteVertex(foo) expected VertexUnknownError, got %T", errUnknown)
	}

	// nil
	errNil := dag.DeleteVertex(nil)
	if errNil == nil {
		t.Errorf("DeleteVertex(nil) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNil.(VertexNilError); !ok {
		t.Errorf("DeleteVertex(nil) expected VertexNilError, got %T", errNil)
	}
}

func TestDAG_AddEdge(t *testing.T) {
	dag := NewDAG()
	v0 := &iVertex{0}
	v1 := &iVertex{1}
	v2 := &iVertex{2}
	v3 := &iVertex{3}

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
		t.Errorf("GetDescendants(v1) = %d, want 1", len(vertices))
	}
	if vertices, _ := dag.GetAncestors(v2); len(vertices) != 1 {
		t.Errorf("GetAncestors(v2) = %d, want 1", len(vertices))
	}

	err := dag.AddEdge(v2, v3)
	if err != nil {
		t.Fatal(err)
	}
	if vertices, _ := dag.GetDescendants(v1); len(vertices) != 2 {
		t.Errorf("GetDescendants(v1) = %d, want 2", len(vertices))
	}
	if vertices, _ := dag.GetAncestors(v3); len(vertices) != 2 {
		t.Errorf("GetAncestors(v3) = %d, want 2", len(vertices))
	}

	_ = dag.AddEdge(v0, v1)
	if vertices, _ := dag.GetDescendants(v0); len(vertices) != 3 {
		t.Errorf("GetDescendants(v0) = %d, want 3", len(vertices))
	}
	if vertices, _ := dag.GetAncestors(v3); len(vertices) != 3 {
		t.Errorf("GetAncestors(v3) = %d, want 3", len(vertices))
	}

	// loop
	errLoopSrcSrc := dag.AddEdge(v1, v1)
	if errLoopSrcSrc == nil {
		t.Errorf("AddEdge(v1, v1) = nil, want %T", SrcDstEqualError{v1, v1})
	}
	if _, ok := errLoopSrcSrc.(SrcDstEqualError); !ok {
		t.Errorf("AddEdge(v1, v1) expected SrcDstEqualError, got %T", errLoopSrcSrc)
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
	v0 := &iVertex{0}
	v1 := &iVertex{1}
	_ = dag.AddEdge(v0, v1)
	if size := dag.GetSize(); size != 1 {
		t.Errorf("GetSize() = %d, want 1", size)
	}
	_ = dag.DeleteEdge(v0, v1)
	if size := dag.GetSize(); size != 0 {
		t.Errorf("GetSize() = %d, want 0", size)
	}

	// unknown
	errUnknown := dag.DeleteEdge(v0, v1)
	if errUnknown == nil {
		t.Errorf("DeleteEdge(v0, v1) = nil, want %T", EdgeUnknownError{})
	}
	if _, ok := errUnknown.(EdgeUnknownError); !ok {
		t.Errorf("DeleteEdge(v0, v1) expected EdgeUnknownError, got %T", errUnknown)
	}

	// nil
	errNilSrc := dag.DeleteEdge(nil, v1)
	if errNilSrc == nil {
		t.Errorf("DeleteEdge(nil, v1) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNilSrc.(VertexNilError); !ok {
		t.Errorf("DeleteEdge(nil, v1) expected VertexNilError, got %T", errNilSrc)
	}
	errNilDst := dag.DeleteEdge(v0, nil)
	if errNilDst == nil {
		t.Errorf("DeleteEdge(v0, nil) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNilDst.(VertexNilError); !ok {
		t.Errorf("DeleteEdge(v0, nil) expected VertexNilError, got %T", errNilDst)
	}

	// unknown
	foo := &iVertex{-1}
	errUnknownSrc := dag.DeleteEdge(foo, v1)
	if errUnknownSrc == nil {
		t.Errorf("DeleteEdge(foo, v1) = nil, want %T", VertexUnknownError{})
	}
	if _, ok := errUnknownSrc.(VertexUnknownError); !ok {
		t.Errorf("DeleteEdge(foo, v1) expected VertexUnknownError, got %T", errUnknownSrc)
	}
	errUnknownDst := dag.DeleteEdge(v0, foo)
	if errUnknownDst == nil {
		t.Errorf("DeleteEdge(v0, foo) = nil, want %T", VertexUnknownError{})
	}
	if _, ok := errUnknownDst.(VertexUnknownError); !ok {
		t.Errorf("DeleteEdge(v0, foo) expected VertexUnknownError, got %T", errUnknownDst)
	}
}

func TestDAG_GetChildren(t *testing.T) {
	dag := NewDAG()
	v1 := &iVertex{1}
	v2 := &iVertex{2}
	v3 := &iVertex{3}
	v4 := &iVertex{4}
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
	_, errNil := dag.GetChildren(nil)
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
	v1 := &iVertex{1}
	v2 := &iVertex{2}
	v3 := &iVertex{3}
	v4 := &iVertex{4}
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
	_, errNil := dag.GetParents(nil)
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
	v1 := &iVertex{1}
	v2 := &iVertex{2}
	v3 := &iVertex{3}
	v4 := &iVertex{4}
	v5 := &iVertex{5}
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
	_, errNil := dag.GetDescendants(nil)
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

func Equal(a, b []Vertex) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestDAG_GetOrderedDescendants(t *testing.T) {
	dag := NewDAG()
	v1 := &iVertex{1}
	v2 := &iVertex{2}
	v3 := &iVertex{3}
	v4 := &iVertex{4}
	v5 := &iVertex{5}
	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v2, v3)
	_ = dag.AddEdge(v2, v4)

	if desc, _ := dag.GetOrderedDescendants(v1); len(desc) != 3 {
		t.Errorf("GetOrderedDescendants(v1) = %d, want 3", len(desc))
	}
	if desc, _ := dag.GetOrderedDescendants(v2); len(desc) != 2 {
		t.Errorf("GetOrderedDescendants(v2) = %d, want 2", len(desc))
	}
	if desc, _ := dag.GetOrderedDescendants(v3); len(desc) != 0 {
		t.Errorf("GetOrderedDescendants(v4) = %d, want 0", len(desc))
	}
	if desc, _ := dag.GetOrderedDescendants(v4); len(desc) != 0 {
		t.Errorf("GetOrderedDescendants(v4) = %d, want 0", len(desc))
	}
	if desc, _ := dag.GetOrderedDescendants(v1); !Equal(desc, []Vertex{v2, v3, v4}) && !Equal(desc, []Vertex{v2, v4, v3}) {
		t.Errorf("GetOrderedDescendants(v4) = %v, want %v or %v", desc, []Vertex{v2, v3, v4}, []Vertex{v2, v4, v3})
	}

	// nil
	_, errNil := dag.GetOrderedDescendants(nil)
	if errNil == nil {
		t.Errorf("GetOrderedDescendants(nil) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNil.(VertexNilError); !ok {
		t.Errorf("GetOrderedDescendants(nil) expected VertexNilError, got %T", errNil)
	}

	// unknown
	_, errUnknown := dag.GetOrderedDescendants(v5)
	if errUnknown == nil {
		t.Errorf("GetOrderedDescendants(v5) = nil, want %T", VertexUnknownError{v5})
	}
	if _, ok := errUnknown.(VertexUnknownError); !ok {
		t.Errorf("GetOrderedDescendants(v5) expected VertexUnknownError, got %T", errUnknown)
	}
}

func TestDAG_GetAncestors(t *testing.T) {
	dag := NewDAG()
	v0 := &iVertex{0}
	v1 := &iVertex{1}
	v2 := &iVertex{2}
	v3 := &iVertex{3}
	v4 := &iVertex{4}
	v5 := &iVertex{5}
	v6 := &iVertex{6}
	v7 := &iVertex{7}

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
	_, errNil := dag.GetAncestors(nil)
	if errNil == nil {
		t.Errorf("GetAncestors(nil) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNil.(VertexNilError); !ok {
		t.Errorf("GetAncestors(nil) expected VertexNilError, got %T", errNil)
	}

	// unknown
	foo := &iVertex{-1}
	_, errUnknown := dag.GetAncestors(foo)
	if errUnknown == nil {
		t.Errorf("GetAncestors(foo) = nil, want %T", VertexUnknownError{foo})
	}
	if _, ok := errUnknown.(VertexUnknownError); !ok {
		t.Errorf("GetAncestors(foo) expected VertexUnknownError, got %T", errUnknown)
	}

}

func TestDAG_GetOrderedAncestors(t *testing.T) {
	dag := NewDAG()
	v1 := &iVertex{1}
	v2 := &iVertex{2}
	v3 := &iVertex{3}
	v4 := &iVertex{4}
	v5 := &iVertex{5}
	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v2, v3)
	_ = dag.AddEdge(v2, v4)

	if desc, _ := dag.GetOrderedAncestors(v4); len(desc) != 2 {
		t.Errorf("GetOrderedAncestors(v4) = %d, want 2", len(desc))
	}
	if desc, _ := dag.GetOrderedAncestors(v2); len(desc) != 1 {
		t.Errorf("GetOrderedAncestors(v2) = %d, want 1", len(desc))
	}
	if desc, _ := dag.GetOrderedAncestors(v1); len(desc) != 0 {
		t.Errorf("GetOrderedAncestors(v1) = %d, want 0", len(desc))
	}
	if desc, _ := dag.GetOrderedAncestors(v4); !Equal(desc, []Vertex{v2, v1}) {
		t.Errorf("GetOrderedAncestors(v4) = %v, want %v", desc, []Vertex{v2, v1})
	}

	// nil
	_, errNil := dag.GetOrderedAncestors(nil)
	if errNil == nil {
		t.Errorf("GetOrderedAncestors(nil) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNil.(VertexNilError); !ok {
		t.Errorf("GetOrderedAncestors(nil) expected VertexNilError, got %T", errNil)
	}

	// unknown
	_, errUnknown := dag.GetOrderedAncestors(v5)
	if errUnknown == nil {
		t.Errorf("GetOrderedAncestors(v5) = nil, want %T", VertexUnknownError{v5})
	}
	if _, ok := errUnknown.(VertexUnknownError); !ok {
		t.Errorf("GetOrderedAncestors(v5) expected VertexUnknownError, got %T", errUnknown)
	}
}

func TestDAG_AncestorsWalker(t *testing.T) {
	dag := NewDAG()
	v1 := &iVertex{1}
	v2 := &iVertex{2}
	v3 := &iVertex{3}
	v4 := &iVertex{4}
	v5 := &iVertex{5}
	v6 := &iVertex{6}
	v7 := &iVertex{7}
	v8 := &iVertex{8}
	v9 := &iVertex{9}
	v10 := &iVertex{10}
	foo := &iVertex{-1}

	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v1, v3)
	_ = dag.AddEdge(v2, v4)
	_ = dag.AddEdge(v2, v5)
	_ = dag.AddEdge(v4, v6)
	_ = dag.AddEdge(v5, v6)
	_ = dag.AddEdge(v6, v7)
	_ = dag.AddEdge(v7, v8)
	_ = dag.AddEdge(v7, v9)
	_ = dag.AddEdge(v8, v10)
	_ = dag.AddEdge(v9, v10)

	vertices, _, _ := dag.AncestorsWalker(v10)
	var ancestors []Vertex
	for v := range vertices {
		ancestors = append(ancestors, v)
	}
	exp1 := []Vertex{v9, v8, v7, v6, v4, v5, v2, v1}
	exp2 := []Vertex{v8, v9, v7, v6, v4, v5, v2, v1}
	exp3 := []Vertex{v9, v8, v7, v6, v5, v4, v2, v1}
	exp4 := []Vertex{v8, v9, v7, v6, v5, v4, v2, v1}
	if !(Equal(ancestors, exp1) || Equal(ancestors, exp2) || Equal(ancestors, exp3) || Equal(ancestors, exp4)) {
		t.Errorf("AncestorsWalker(v10) = %v, want %v, %v, %v, or %v ", ancestors, exp1, exp2, exp3, exp4)
	}

	// nil
	_, _, errNil := dag.AncestorsWalker(nil)
	if errNil == nil {
		t.Errorf("AncestorsWalker(nil) = nil, want %T", VertexNilError{})
	}
	if _, ok := errNil.(VertexNilError); !ok {
		t.Errorf("AncestorsWalker(nil) expected VertexNilError, got %T", errNil)
	}

	// unknown
	_, _, errUnknown := dag.AncestorsWalker(foo)
	if errUnknown == nil {
		t.Errorf("AncestorsWalker(foo) = nil, want %T", VertexUnknownError{foo})
	}
	if _, ok := errUnknown.(VertexUnknownError); !ok {
		t.Errorf("AncestorsWalker(foo) expected VertexUnknownError, got %T", errUnknown)
	}
}

func TestDAG_AncestorsWalkerSignal(t *testing.T) {
	dag := NewDAG()

	v1 := &iVertex{1}
	v2 := &iVertex{2}
	v3 := &iVertex{3}
	v4 := &iVertex{4}
	v5 := &iVertex{5}
	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v2, v3)
	_ = dag.AddEdge(v2, v4)
	_ = dag.AddEdge(v4, v5)

	var ancestors []Vertex
	vertices, signal, _ := dag.AncestorsWalker(v5)
	for v := range vertices {
		ancestors = append(ancestors, v)
		if v == v2 {
			signal <- true
			break
		}
	}
	if !Equal(ancestors, []Vertex{v4, v2}) {
		t.Errorf("AncestorsWalker(v4) = %v, want %v", ancestors, []Vertex{v4, v2})
	}

}

func TestDAG_ReduceTransitively(t *testing.T) {
	dag := NewDAG()
	accountCreate := &sVertex{"AccountCreate"}
	projectCreate := &sVertex{"ProjectCreate"}
	networkCreate := &sVertex{"NetworkCreate"}
	contactCreate := &sVertex{"ContactCreate"}
	authzCreate := &sVertex{"AuthzCreate"}
	mailSend := &sVertex{"MailSend"}

	_ = dag.AddEdge(accountCreate, projectCreate)
	_ = dag.AddEdge(accountCreate, networkCreate)
	_ = dag.AddEdge(accountCreate, contactCreate)
	_ = dag.AddEdge(accountCreate, authzCreate)
	_ = dag.AddEdge(accountCreate, mailSend)

	_ = dag.AddEdge(projectCreate, mailSend)
	_ = dag.AddEdge(networkCreate, mailSend)
	_ = dag.AddEdge(contactCreate, mailSend)
	_ = dag.AddEdge(authzCreate, mailSend)

	if order := dag.GetOrder(); order != 6 {
		t.Errorf("GetOrder() = %d, want 6", order)
	}
	if size := dag.GetSize(); size != 9 {
		t.Errorf("GetSize() = %d, want 9", size)
	}
	if isEdge, _ := dag.IsEdge(accountCreate, mailSend); !isEdge {
		t.Errorf("IsEdge(accountCreate, mailSend) = %t, want %t", isEdge, true)
	}

	dag.ReduceTransitively()

	if order := dag.GetOrder(); order != 6 {
		t.Errorf("GetOrder() = %d, want 6", order)
	}
	if size := dag.GetSize(); size != 8 {
		t.Errorf("GetSize() = %d, want 8", size)
	}
	if isEdge, _ := dag.IsEdge(accountCreate, mailSend); isEdge {
		t.Errorf("IsEdge(accountCreate, mailSend) = %t, want %t", isEdge, false)
	}

	ordered, _ := dag.GetOrderedDescendants(accountCreate)
	length := len(ordered)
	if length != 5 {
		t.Errorf("length(ordered) = %d, want 5", length)
	}
	last := ordered[length-1]
	if last != mailSend {
		t.Errorf("ordered[length-1]) = %s, want %s", last, mailSend.String())
	}
}

func TestDAG_String(t *testing.T) {
	dag := NewDAG()
	v1 := &iVertex{1}
	v2 := &iVertex{2}
	v3 := &iVertex{3}
	v4 := &iVertex{4}
	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v2, v3)
	_ = dag.AddEdge(v2, v4)
	expected := "DAG Vertices: 4 - Edges: 3"
	s := dag.String()
	if s[:len(expected)] != expected {
		t.Errorf("String() = \"%s\", want \"%s\"", s, expected)
	}
}

func TestErrors(t *testing.T) {
	v1 := &iVertex{1}
	v2 := &iVertex{2}

	tests := []struct {
		want string
		err  error
	}{
		{"don't know what to do with 'nil'", VertexNilError{}},
		{"'1' is already known", VertexDuplicateError{v1}},
		{"'1' is unknown", VertexUnknownError{v1}},
		{"edge between '1' and '2' is already known", EdgeDuplicateError{v1, v2}},
		{"edge between '1' and '2' is unknown", EdgeUnknownError{v1, v2}},
		{"edge between '1' and '2' would create a loop", EdgeLoopError{v1, v2}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt.err), func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Example() {

	// initialize a new graph
	d := NewDAG()

	// init three vertices
	v1 := &iVertex{1}
	v2 := &iVertex{2}
	v3 := &iVertex{3}

	// add the above vertices and connect them with two edges
	_ = d.AddEdge(v1, v2)
	_ = d.AddEdge(v1, v3)

	// describe the graph
	fmt.Print(d.String())

	// Unordered output:
	// DAG Vertices: 3 - Edges: 2
	// Vertices:
	//   2
	//   3
	//   1
	// Edges:
	//   1 -> 2
	//   1 -> 3
}

func ExampleDAG_AncestorsWalker() {
	dag := NewDAG()

	v1 := &iVertex{1}
	v2 := &iVertex{2}
	v3 := &iVertex{3}
	v4 := &iVertex{4}
	v5 := &iVertex{5}
	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v2, v3)
	_ = dag.AddEdge(v2, v4)
	_ = dag.AddEdge(v4, v5)

	var ancestors []Vertex
	vertices, signal, _ := dag.AncestorsWalker(v5)
	for v := range vertices {
		ancestors = append(ancestors, v)
		if v == v2 {
			signal <- true
			break
		}
	}
	fmt.Printf("%v", ancestors)

	// Output:
	//   [4 2]
}

func TestLarge(t *testing.T) {
	d := NewDAG()
	root := &iVertex{1}
	levels := 7
	branches := 8

	expectedVertexCount, _ := largeAux(d, levels, branches, root)
	expectedVertexCount++
	vertexCount := len(d.GetVertices())
	if vertexCount != expectedVertexCount {
		t.Errorf("GetVertices() = %d, want %d", vertexCount, expectedVertexCount)
	}

	descendants, _ := d.GetDescendants(root)
	descendantsCount := len(descendants)
	expectedDescendantsCount := vertexCount - 1
	if descendantsCount != expectedDescendantsCount {
		t.Errorf("GetDescendants(root) = %d, want %d", descendantsCount, expectedDescendantsCount)
	}

	_, _ = d.GetDescendants(root)

	children, _ := d.GetChildren(root)
	childrenCount := len(children)
	expectedChildrenCount := branches
	if childrenCount != expectedChildrenCount {
		t.Errorf("GetChildren(root) = %d, want %d", childrenCount, expectedChildrenCount)
	}

	/*
		var childList []Vertex
		for x := range children {
			childList = append(childList, x)
		}
		_ = d.DeleteEdge(root, childList[0])
	*/
}

func largeAux(d *DAG, level int, branches int, parent *iVertex) (int, int) {
	var vertexCount int
	var edgeCount int
	if level > 1 {
		if branches < 1 || branches > 9 {
			panic("number of branches must be between 1 and 9")
		}
		for i := 1; i <= branches; i++ {
			value := (*parent).value*10 + i
			child := &iVertex{value}
			vertexCount++
			err := d.AddEdge(parent, child)
			edgeCount++
			if err != nil {
				panic(err)
			}
			childVertexCount, childEdgeCount := largeAux(d, level-1, branches, child)
			vertexCount += childVertexCount
			edgeCount += childEdgeCount
		}
	}
	return vertexCount, edgeCount
}

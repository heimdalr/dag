package dag

import (
	"encoding/json"
	"testing"
)

type testNonComparableVertexType struct {
	ID                 string            `json:"i"`
	NotComparableField map[string]string `json:"v"`
}

func TestOverrideVertexHashFunOption(t *testing.T) {
	dag := NewDAG()
	/*     1    4
	 *     |\  /
	 *     | 2
	 *     |/
	 *     3
	 */

	dag.Options(Options{
		VertexHashFunc: func(v interface{}) interface{} {
			return v.(testNonComparableVertexType).ID
		}})

	testVertex1 := testNonComparableVertexType{
		ID:                 "1",
		NotComparableField: map[string]string{"not": "comparable"},
	}
	vertexId1, err := dag.addVertex(testVertex1)
	if err != nil {
		t.Errorf("Should create a vertex with a not comparable field when a correct VertexHashFunc option is set")
	}

	testVertex2 := testNonComparableVertexType{
		ID:                 "2",
		NotComparableField: map[string]string{"stillNot": "comparable"},
	}
	vertexId2, err := dag.addVertex(testVertex2)
	if err != nil {
		t.Errorf("Should create a vertex with a not comparable field when a correct VertexHashFunc option is set")
	}
	err = dag.AddEdge(vertexId1, vertexId2)
	if err != nil {
		t.Errorf("Should create an edge between vertices with not comparable fields when a correct VertexHashFunc option is set")
	}

	testVertex3 := testNonComparableVertexType{
		ID:                 "3",
		NotComparableField: map[string]string{"stillNot": "comparable"},
	}
	vertexId3, err := dag.addVertex(testVertex3)
	if err != nil {
		t.Errorf("Should create a vertex with a not comparable field when a correct VertexHashFunc option is set")
	}

	err = dag.AddEdge(vertexId1, vertexId3)
	if err != nil {
		t.Errorf("Should create an edge between vertices with not comparable fields when a correct VertexHashFunc option is set")
	}
	err = dag.AddEdge(vertexId2, vertexId3)
	if err != nil {
		t.Errorf("Should create an edge between vertices with not comparable fields when a correct VertexHashFunc option is set")
	}

	testVertex4 := testNonComparableVertexType{
		ID:                 "4",
		NotComparableField: map[string]string{"stillNot": "comparable"},
	}
	vertexId4, err := dag.addVertex(testVertex4)
	if err != nil {
		t.Errorf("Should create a vertex with a not comparable field when a correct VertexHashFunc option is set")
	}

	err = dag.AddEdge(vertexId4, vertexId2)
	if err != nil {
		t.Errorf("Should create an edge between vertices with not comparable fields when a correct VertexHashFunc option is set")
	}

	isEdge, err := dag.IsEdge(vertexId1, vertexId2)
	if !isEdge || err != nil {
		t.Errorf("Should return true for edge between vertices with not comparable fields when a correct VertexHashFunc option is set")
	}

	err = dag.DeleteEdge(vertexId1, vertexId3)
	if err != nil {
		t.Errorf("Should delete an edge between vertices with not comparable fields when a correct VertexHashFunc option is set")
	}

	isEdge, err = dag.IsEdge(vertexId1, vertexId3)
	if isEdge || err != nil {
		t.Errorf("Should return false for edge between vertices with not comparable fields when a correct VertexHashFunc option is set")
	}

	err = dag.DeleteVertex(vertexId2)
	if err != nil {
		t.Errorf("Should delete a vertex with not comparable fields when a correct VertexHashFunc option is set")
	}

	vertexId2, err = dag.addVertex(testVertex2)
	if err != nil {
		t.Errorf("Should create a vertex with a not comparable field when a correct VertexHashFunc option is set")
	}

	_ = dag.AddEdge(vertexId1, vertexId2)
	_ = dag.AddEdge(vertexId2, vertexId3)
	err = dag.AddEdge(vertexId4, vertexId2)
	if err != nil {
		t.Errorf("Should create an edge between vertices with not comparable fields when a correct VertexHashFunc option is set")
	}

	roots := dag.GetRoots()
	if len(roots) != 2 {
		t.Errorf("Should return 2 roots")
	}
	for rootId := range roots {
		if isRoot, err := dag.IsRoot(rootId); !isRoot || err != nil {
			t.Errorf("Should return true for root")
		}
	}

	leaves := dag.GetLeaves()
	if len(leaves) != 1 {
		t.Errorf("Should return 1 leaf")
	}
	for leafId := range leaves {
		if isLeaf, err := dag.IsLeaf(leafId); !isLeaf || err != nil {
			t.Errorf("Should return true for leaf")
		}
	}

	vertex2Parents, err := dag.GetParents(vertexId2)
	if len(vertex2Parents) != 2 || err != nil {
		t.Errorf("Should return 2 parents for vertex 2")
	}

	vertex2Children, err := dag.GetChildren(vertexId2)
	if len(vertex2Children) != 1 || err != nil {
		t.Errorf("Should return 1 child for vertex 2")
	}

	vertex3Ancestors, err := dag.GetAncestors(vertexId3)
	if len(vertex3Ancestors) != 3 || err != nil {
		t.Errorf("Should return 3 ancestors for vertex 3, received %d", len(vertex3Ancestors))
	}

	vertex3OrderedAncestors, err := dag.GetOrderedAncestors(vertexId3)
	if len(vertex3OrderedAncestors) != 3 || err != nil {
		t.Errorf("Should return 3 ancestors for vertex 3, received %d", len(vertex3OrderedAncestors))
	}

	vertex4Descendants, err := dag.GetDescendants(vertexId4)
	if len(vertex4Descendants) != 2 || err != nil {
		t.Errorf("Should return 2 descendants for vertex 4, received %d", len(vertex4Descendants))
	}

	vertex4OrderedDescendants, err := dag.GetOrderedDescendants(vertexId4)
	if len(vertex4OrderedDescendants) != 2 || err != nil {
		t.Errorf("Should return 2 descendants for vertex 4, received %d", len(vertex4OrderedDescendants))
	}

	_, _, err = dag.GetDescendantsGraph(vertexId1)
	if err != nil {
		t.Errorf("Should return a string representation of the descendants graph")
	}

	_, _, err = dag.GetAncestorsGraph(vertexId1)
	if err != nil {
		t.Errorf("Should return a string representation of the ancestors graph")
	}

	_, err = dag.Copy()
	if err != nil {
		t.Errorf("Should return a copy of the DAG")
	}

	dagString := dag.String()
	if dagString == "" {
		t.Errorf("Should return a string representation of the DAG")
	}

	dag.ReduceTransitively()
	dag.FlushCaches()
	dag.DescendantsWalker(vertexId1) // nolint:errcheck

	mv := newMarshalVisitor(dag)
	dag.DFSWalk(mv)
	dag.BFSWalk(mv)
	dag.OrderedWalk(mv)

	_, err = dag.MarshalJSON()
	if err != nil {
		t.Error(err)
	}

	data, err := json.Marshal(dag)
	if err != nil {
		t.Error(err)
	}

	var wd testNonComparableStorableDAG
	_, err = UnmarshalJSON(data, &wd, dag.options)
	if err != nil {
		t.Fatal(err)
	}
}

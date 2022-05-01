package dag

import (
	"encoding/json"
	"errors"
)

// MarshalJSON returns the JSON encoding of DAG.
//
// It traverses the DAG using the Depth-First-Search algorithm
// and uses an internal structure to store vertices and edges.
func (d *DAG) MarshalJSON() ([]byte, error) {
	mv := newMarshalVisitor(d)
	DFSWalk(d, mv)
	return json.Marshal(mv.storableDAG)
}

// UnmarshalJSON is an informative method. See the UnmarshalJSON function below.
func (d *DAG) UnmarshalJSON(data []byte) error {
	return errors.New("This method is not supported, request function UnmarshalJSON instead")
}

// UnmarshalJSON parses the JSON-encoded data that defined by StorableDAG.
// It returns a new DAG defined by the vertices and edges of wd.
// If the internal structure of data and wd do not match,
// then deserialization will fail and return json eror
//
// Because the vertex data passed in by the user is an interface{},
// it does not indicate a specific structure, so it cannot be deserialized.
// And this function needs to pass in a clear DAG structure.
//
// Example:
// dag := NewDAG()
// data, err := json.Marshal(d)
// if err != nil {
//     panic(err)
// }
// var wd YourStorableDAG
// restoredDag, err := UnmarshalJSON(data, &wd)
// if err != nil {
//     panic(err)
// }
//
// For more specific information please read the test code.
func UnmarshalJSON(data []byte, wd StorableDAG) (*DAG, error) {
	err := json.Unmarshal(data, &wd)
	if err != nil {
		return nil, err
	}
	dag := NewDAG()
	for _, v := range wd.Vertices() {
		dag.AddVertexByID(v.Vertex())
	}
	for _, e := range wd.Edges() {
		dag.AddEdge(e.Edge())
	}
	return dag, nil
}

type marshalVisitor struct {
	d *DAG
	storableDAG
}

func newMarshalVisitor(d *DAG) *marshalVisitor {
	return &marshalVisitor{d: d}
}

func (mv *marshalVisitor) Visit(v Vertexer) {
	mv.StorableVertices = append(mv.StorableVertices, v)

	srcID, _ := v.Vertex()
	children, _ := mv.d.getChildren(srcID)
	ids := vertexIDs(children)
	for _, dstID := range ids {
		e := storableEdge{SrcID: srcID, DstID: dstID}
		mv.StorableEdges = append(mv.StorableEdges, e)
	}
}

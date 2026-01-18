package dag

// Vertex describes the per-vertex behavior exposed to the graph.
type Vertex interface {
	ID() string
	Value() interface{}
}

// vertex implements the Vertex interface and owns vertex-level data.
type vertex struct {
	id  string
	val interface{}
}

func (v *vertex) ID() string {
	return v.id
}

func (v *vertex) Value() interface{} {
	return v.val
}


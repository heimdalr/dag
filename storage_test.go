package dag

type testVertex struct {
	WID string `json:"i"`
	Val string `json:"v"`
}

func (tv testVertex) ID() string {
	return tv.WID
}

func (tv testVertex) Vertex() (id string, value interface{}) {
	return tv.WID, tv.Val
}

type testStorableDAG struct {
	StorableVertices []testVertex   `json:"vs"`
	StorableEdges    []storableEdge `json:"es"`
}

func (g testStorableDAG) Vertices() []Vertexer {
	l := make([]Vertexer, 0, len(g.StorableVertices))
	for _, v := range g.StorableVertices {
		l = append(l, v)
	}
	return l
}

func (g testStorableDAG) Edges() []Edger {
	l := make([]Edger, 0, len(g.StorableEdges))
	for _, v := range g.StorableEdges {
		l = append(l, v)
	}
	return l
}

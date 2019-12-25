package dag

import (
	"fmt"
	"sync"
)

type Vertex interface {
	String() string
}

// Error type to describe the situation, that a given vertex does not exit in the graph.
type VertexUnknownError struct {
	v *Vertex
}

// Implements the `error` interface.
func (e VertexUnknownError) Error() string {
	return fmt.Sprintf("'%s' is unknown", (*e.v).String())
}

// Error type to describe loop errors (i.e. errors that where raised to prevent establishing loops in the graph).
type LoopError struct {
	src *Vertex
	dst *Vertex
}

// Implements the `error` interface.
func (e LoopError) Error() string {
	return fmt.Sprintf("loop between '%s' and '%s'", (*e.src).String(), (*e.dst).String())
}

// The DAG type implements a Directed Acyclic Graph.
type DAG struct {
	vertices     map[*Vertex]bool
	muVertices   sync.Mutex
	inboundEdge  map[*Vertex]map[*Vertex]bool
	outboundEdge map[*Vertex]map[*Vertex]bool
	muEdges      sync.Mutex
}

// Creates / initializes a new Directed Acyclic Graph or DAG.
func NewDAG() *DAG {
	d := &DAG{
		vertices:     make(map[*Vertex]bool),
		inboundEdge:  make(map[*Vertex]map[*Vertex]bool),
		outboundEdge: make(map[*Vertex]map[*Vertex]bool),
	}
	return d
}

// Add a vertex.
// For vertices that are part of an edge use AddEdge() instead.
func (d *DAG) AddVertex(v *Vertex) {
	if v == nil {
		return
	}
	d.muVertices.Lock()
	d.vertices[v] = true
	d.muVertices.Unlock()
}

// Delete a vertex including all inbound and outbound edges.
func (d *DAG) DeleteVertex(v *Vertex) {
	if v == nil {
		return
	}
	if _, ok := d.vertices[v]; ok {
		d.muEdges.Lock()
		delete(d.inboundEdge, v)
		delete(d.outboundEdge, v)
		d.muEdges.Unlock()
		d.muVertices.Lock()
		delete(d.vertices, v)
		d.muVertices.Unlock()
	}
}

func (d *DAG) addEdgeAux(src *Vertex, dst *Vertex, check bool) error {
	if src == nil || dst == nil {
		return nil
	}

	// ensure vertices
	d.muVertices.Lock()
	d.vertices[src] = true
	d.vertices[dst] = true
	d.muVertices.Unlock()

	// check for circles, iff desired
	if check {
		if src == dst {
			return LoopError{src, dst}
		}
		descendants, _ := d.GetDescendants(dst)
		if descendants[src] {
			return LoopError{src, dst}
		}
	}

	// test / compute edge nodes
	outbound, outboundExists := d.outboundEdge[src]
	inbound, inboundExists := d.inboundEdge[dst]

	d.muEdges.Lock()

	// add outbound
	if !outboundExists {
		newSet := make(map[*Vertex]bool)
		d.outboundEdge[src] = newSet
		outbound = newSet
	}
	outbound[dst] = true

	// add inbound
	if !inboundExists {
		newSet := make(map[*Vertex]bool)
		d.inboundEdge[dst] = newSet
		inbound = newSet
	}
	inbound[src] = true

	d.muEdges.Unlock()
	return nil
}

// Add an edge prevents circles
func (d *DAG) AddEdgeSafe(src *Vertex, dst *Vertex) error {
	return d.addEdgeAux(src, dst, true)
}

// Add an edge without checking for circles.
func (d *DAG) AddEdge(src *Vertex, dst *Vertex) error {
	return d.addEdgeAux(src, dst, false)

}

// Delete an edge.
func (d *DAG) DeleteEdge(src *Vertex, dst *Vertex) {

	// test / compute edge nodes
	_, outboundExists := d.outboundEdge[src][dst]
	_, inboundExists := d.inboundEdge[dst][src]

	if inboundExists || outboundExists {
		d.muEdges.Lock()

		// delete outbound
		if outboundExists {
			delete(d.inboundEdge[dst], dst)
		}

		// delete inbound
		if inboundExists {
			delete(d.outboundEdge[src], dst)
		}
		d.muEdges.Unlock()
	}
}

// Return the total number of vertices.
func (d *DAG) GetOrder() int {
	return len(d.vertices)
}

// Return the total number of edges.
func (d *DAG) GetSize() int {
	count := 0
	for _, value := range d.outboundEdge {
		count += len(value)
	}
	return count
}

// Return all vertices without children.
func (d *DAG) GetLeafs() map[*Vertex]bool {
	leafs := make(map[*Vertex]bool)
	for v := range d.vertices {
		dstIds, ok := d.outboundEdge[v]
		if !ok || len(dstIds) == 0 {
			leafs[v] = true
		}
	}
	return leafs
}

// Return all vertices without parents.
func (d *DAG) GetRoots() map[*Vertex]bool {
	roots := make(map[*Vertex]bool)
	for v := range d.vertices {
		srcIds, ok := d.inboundEdge[v]
		if !ok || len(srcIds) == 0 {
			roots[v] = true
		}
	}
	return roots
}

// Return all vertices.
func (d *DAG) GetVertices() map[*Vertex]bool {
	return d.vertices
}

// Return all children of the given vertex.
func (d *DAG) GetChildren(v *Vertex) (map[*Vertex]bool, error) {
	if _, ok := d.vertices[v]; !ok {
		return nil, VertexUnknownError{v}
	}
	return d.outboundEdge[v], nil
}

// Return all parents of the given vertex.
func (d *DAG) GetParents(v *Vertex) (map[*Vertex]bool, error) {
	if _, ok := d.vertices[v]; !ok {
		return nil, VertexUnknownError{v}
	}
	return d.inboundEdge[v], nil
}

func (d *DAG) getAncestorsAux(v *Vertex, ancestors map[*Vertex]bool, m sync.Mutex) {
	if parents, ok := d.inboundEdge[v]; ok {
		for parent := range parents {
			d.getAncestorsAux(parent, ancestors, m)
			m.Lock()
			ancestors[parent] = true
			m.Unlock()
		}
	}
}

// Return all Ancestors of the given vertex.
func (d *DAG) GetAncestors(v *Vertex) (map[*Vertex]bool, error) {
	if _, ok := d.vertices[v]; !ok {
		return nil, VertexUnknownError{v}
	}
	ancestors := make(map[*Vertex]bool)
	var m sync.Mutex
	d.getAncestorsAux(v, ancestors, m)
	return ancestors, nil
}

func (d *DAG) getDescendantsAux(v *Vertex, descendents map[*Vertex]bool, m sync.Mutex) {
	if children, ok := d.outboundEdge[v]; ok {
		for child := range children {
			d.getDescendantsAux(child, descendents, m)
			m.Lock()
			descendents[child] = true
			m.Unlock()
		}
	}
}

// Return all Descendants of the given vertex.
func (d *DAG) GetDescendants(v *Vertex) (map[*Vertex]bool, error) {
	if _, ok := d.vertices[v]; !ok {
		return nil, VertexUnknownError{v}
	}
	descendents := make(map[*Vertex]bool)
	var m sync.Mutex
	d.getDescendantsAux(v, descendents, m)
	return descendents, nil
}

func (d *DAG) String() string {
	result := fmt.Sprintf("DAG Vertices: %d - Edges: %d\n", d.GetOrder(), d.GetSize())
	result += fmt.Sprintf("Vertices:\n")
	for k := range d.vertices {
		result += fmt.Sprintf("  %v\n", (*k).String())
	}
	result += fmt.Sprintf("Edges:\n")
	for v, children := range d.outboundEdge {
		for child := range children {
			result += fmt.Sprintf("  %s -> %s\n", (*v).String(), (*child).String())
		}
	}
	return result
}

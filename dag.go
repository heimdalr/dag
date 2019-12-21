package dag

import (
	"errors"
	"fmt"
	"sync"
)

type Vertex interface {
	String() string
}

type id = *Vertex
type idSet = map[id]bool

// DAG type implements a Directed Acyclic Graph data structure.
type DAG struct {
	vertices     idSet
	muVertices   sync.Mutex
	inboundEdge  map[id]idSet
	outboundEdge map[id]idSet
	muEdges      sync.Mutex
}

// Creates a new Directed Acyclic Graph or DAG.
func NewDAG() *DAG {
	d := &DAG{
		vertices:     make(idSet),
		inboundEdge:  make(map[id]idSet),
		outboundEdge: make(map[id]idSet),
	}
	return d
}

// Add a vertex.
func (d *DAG) AddVertex(v *Vertex) error {
	if _, ok := d.vertices[v]; ok {
		return errors.New(fmt.Sprintf("duplicate %s", (*v).String()))
	}
	d.muVertices.Lock()
	d.vertices[v] = true
	d.muVertices.Unlock()
	return nil
}

// Delete a vertex including all inbound and outbound edges.
func (d *DAG) DeleteVertex(v *Vertex) {
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

// Add an edge, iff both vertices exist.
func (d *DAG) AddEdge(src *Vertex, dst *Vertex) error {

	// sanity checking
	if src == dst {
		return errors.New(fmt.Sprintf("src (%s) and dst (%s) must be different", (*src).String(), (*dst).String()))
	}
	if _, ok := d.vertices[src]; !ok {
		return errors.New(fmt.Sprintf("%s is unknown", (*src).String()))
	}
	if _, ok := d.vertices[dst]; !ok {
		return errors.New(fmt.Sprintf("%s is unknown", (*dst).String()))
	}

	// test / compute edge nodes
	outbound, outboundExists := d.outboundEdge[src]
	inbound, inboundExists := d.inboundEdge[dst]

	d.muEdges.Lock()

	// add outbound
	if !outboundExists {
		newSet := make(idSet)
		d.outboundEdge[src] = newSet
		outbound = newSet
	}
	outbound[dst] = true

	// add inbound
	if !inboundExists {
		newSet := make(idSet)
		d.inboundEdge[dst] = newSet
		inbound = newSet
	}
	inbound[src] = true

	d.muEdges.Unlock()
	return nil
}

// Delete an edge, iff such exists.
func (d *DAG) DeleteEdge(src *Vertex, dst *Vertex) error {

	// sanity checking
	if src == dst {
		return errors.New(fmt.Sprintf("src (%s) and dst (%s) must be different", (*src).String(), (*dst).String()))
	}
	if _, ok := d.vertices[src]; !ok {
		return errors.New(fmt.Sprintf("%s is unknown", (*src).String()))
	}
	if _, ok := d.vertices[dst]; !ok {
		return errors.New(fmt.Sprintf("%s is unknown", (*dst).String()))
	}

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

	return nil
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
func (d *DAG) GetLeafs() []*Vertex {
	var leafs []*Vertex
	for v := range d.vertices {
		dstIds, ok := d.outboundEdge[v]
		if !ok || len(dstIds) == 0 {
			leafs = append(leafs, v)
		}
	}
	return leafs
}

// Return all vertices without parents.
func (d *DAG) GetRoots() []*Vertex {
	var roots []*Vertex
	for v := range d.vertices {
		srcIds, ok := d.inboundEdge[v]
		if !ok || len(srcIds) == 0 {
			roots = append(roots, v)
		}
	}
	return roots
}

// Return all vertices.
func (d *DAG) GetVertices() []*Vertex {
	length := len(d.vertices)
	vertices := make([]*Vertex, length)
	i := 0
	for v := range d.vertices {
		vertices[i] = v
		i += 1
	}
	return vertices
}

// Return all children of the given vertex.
func (d *DAG) GetChildren(v *Vertex) ([]*Vertex, error) {
	if _, ok := d.vertices[v]; !ok {
		return nil, errors.New(fmt.Sprintf("%s is unknown", (*v).String()))
	}
	if children, ok := d.outboundEdge[v]; ok {
		result := make([]*Vertex, len(children))
		i := 0
		for child := range children {
			result[i] = child
			i += 1
		}
		return result, nil
	}
	return nil, nil
}

// Return all parents of the given vertex.
func (d *DAG) GetParents(v *Vertex) ([]*Vertex, error) {
	if _, ok := d.vertices[v]; !ok {
		return nil, errors.New(fmt.Sprintf("%s is unknown", (*v).String()))
	}
	if parents, ok := d.inboundEdge[v]; ok {
		result := make([]*Vertex, len(parents))
		i := 0
		for parent := range parents {
			result[i] = parent
			i += 1
		}
		return result, nil
	}
	return nil, nil
}

func (d *DAG) getAncestorsAux(v *Vertex) []*Vertex {
	var ancestors []*Vertex
	if parents, ok := d.inboundEdge[v]; ok {
		for parent := range parents {
			ancestors = append(ancestors, d.getAncestorsAux(parent)...)
			ancestors = append(ancestors, parent)
		}
	}
	return ancestors
}

// Return all Ancestors of the given vertex.
func (d *DAG) GetAncestors(v *Vertex) ([]*Vertex, error) {
	if _, ok := d.vertices[v]; !ok {
		return nil, errors.New(fmt.Sprintf("%s is unknown", (*v).String()))
	}
	return d.getAncestorsAux(v), nil
}

func (d *DAG) getDescendantsAux(v *Vertex) []*Vertex {
	var descendants []*Vertex
	if children, ok := d.outboundEdge[v]; ok {
		for child := range children {
			descendants = append(descendants, d.getDescendantsAux(child)...)
			descendants = append(descendants, child)
		}
	}
	return descendants
}

// Return all Ancestors of the given vertex.
func (d *DAG) GetDescendants(v *Vertex) ([]*Vertex, error) {
	if _, ok := d.vertices[v]; !ok {
		return nil, errors.New(fmt.Sprintf("%s is unknown", (*v).String()))
	}
	return d.getDescendantsAux(v), nil
}

func (d *DAG) String() string {
	result := fmt.Sprintf("DAG Vertices: %d - Edges: %d\n", d.GetOrder(), d.GetSize())
	result += fmt.Sprintf("Vertices:\n")
	for _, v := range d.vertices {
		result += fmt.Sprintf("  %v\n", v)
	}
	result += fmt.Sprintf("Edges:\n")
	for v, children := range d.outboundEdge {
		for child := range children {
			result += fmt.Sprintf("  %s -> %s\n", (*v).String(), (*child).String())
		}
	}
	return result
}

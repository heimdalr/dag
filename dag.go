// Package dag implements a Directed Acyclic Graph data structure and relevant methods.
package dag

import (
	"fmt"
	"sync"
)

// Interface for the nodes in the DAG.
type Vertex interface {
	String() string
}


// The DAG type implements a Directed Acyclic Graph.
type DAG struct {
	vertices         map[Vertex]bool
	verticesLocked   map[Vertex]*cMutex
	muVertices       sync.Mutex
	inboundEdge      map[Vertex]map[Vertex]bool
	outboundEdge     map[Vertex]map[Vertex]bool
	ancestorCache    map[Vertex]map[Vertex]bool
	descendantsCache map[Vertex]map[Vertex]bool
	muEdges          sync.Mutex
}

// Creates / initializes a new Directed Acyclic Graph or DAG.
func NewDAG() *DAG {
	return &DAG{
		vertices:         make(map[Vertex]bool),
		verticesLocked:   make(map[Vertex]*cMutex),
		inboundEdge:      make(map[Vertex]map[Vertex]bool),
		outboundEdge:     make(map[Vertex]map[Vertex]bool),
		ancestorCache:    make(map[Vertex]map[Vertex]bool),
		descendantsCache: make(map[Vertex]map[Vertex]bool),
	}
}

// Add a vertex.
// For vertices that are part of an edge use AddEdge() instead.
func (d *DAG) AddVertex(v Vertex) error {

	// sanity checking
	if v == nil {
		return VertexNilError{}
	}
	if _, exists := d.vertices[v]; exists {
		return VertexDuplicateError{v}
	}

	d.muVertices.Lock()
	d.vertices[v] = true
	d.muVertices.Unlock()

	return nil
}

// Delete a vertex including all inbound and outbound edges. Delete cached ancestors and descendants of relevant
// vertices.
func (d *DAG) DeleteVertex(v Vertex) error {

	// sanity checking
	if v == nil {
		return VertexNilError{}
	}
	if _, exists := d.vertices[v]; !exists {
		return VertexUnknownError{v}
	}

	// get descendents and ancestors as they are now
	descendants, _ := d.GetDescendants(v)
	ancestors, _ := d.GetAncestors(v)

	d.muEdges.Lock()

	// delete v in outbound edges of parents
	if _, exists := d.inboundEdge[v]; exists {
		for parent := range d.inboundEdge[v] {
			delete(d.outboundEdge[parent], v)
		}
	}

	// delete v in inbound edges of children
	if _, exists := d.outboundEdge[v]; exists {
		for child := range d.outboundEdge[v] {
			delete(d.inboundEdge[child], v)
		}
	}

	// delete in- and outbound of v itself
	delete(d.inboundEdge, v)
	delete(d.outboundEdge, v)

	// for v and all its descendants delete cached ancestors
	for descendant := range descendants {
		if _, exists := d.ancestorCache[descendant]; exists {
			delete(d.ancestorCache, descendant)
		}
	}
	delete(d.ancestorCache, v)

	// for v and all its ancestors delete cached descendants
	for ancestor := range ancestors {
		if _, exists := d.descendantsCache[ancestor]; exists {
			delete(d.descendantsCache, ancestor)
		}
	}
	delete(d.descendantsCache, v)

	d.muEdges.Unlock()

	d.muVertices.Lock()

	// delete v itself
	delete(d.vertices, v)

	d.muVertices.Unlock()

	return nil
}

// Add an edge while preventing circles.
func (d *DAG) AddEdge(src Vertex, dst Vertex) error {

	// sanity checking
	if src == nil {
		return VertexNilError{}
	}
	if dst == nil {
		return VertexNilError{}
	}

	// ensure vertices
	d.muVertices.Lock()
	d.vertices[src] = true
	d.vertices[dst] = true
	d.muVertices.Unlock()

	// test / compute edge nodes
	_, outboundExists := d.outboundEdge[src]
	_, inboundExists := d.inboundEdge[dst]

	// if the edge is already known, there is nothing else to do
	if outboundExists && d.outboundEdge[src][dst] && inboundExists && d.inboundEdge[dst][src] {
		return EdgeDuplicateError{src, dst}
	}

	// get descendents and ancestors as they are now
	descendants, _ := d.GetDescendants(dst)
	ancestors, _ := d.GetAncestors(src)

	// check for circles, iff desired
	if src == dst || descendants[src] {
		return EdgeLoopError{src, dst}
	}

	d.muEdges.Lock()

	// prepare d.outbound[src], iff needed
	if !outboundExists {
		d.outboundEdge[src] = make(map[Vertex]bool)
	}

	// dst is a child of src
	d.outboundEdge[src][dst] = true

	// prepare d.inboundEdge[dst], iff needed
	if !inboundExists {
		d.inboundEdge[dst] = make(map[Vertex]bool)
	}

	// src is a parent of dst
	d.inboundEdge[dst][src] = true

	// for dst and all its descendants delete cached ancestors
	for descendant := range descendants {
		if _, exists := d.ancestorCache[descendant]; exists {
			delete(d.ancestorCache, descendant)
		}
	}
	delete(d.ancestorCache, dst)

	// for src and all its ancestors delete cached descendants
	for ancestor := range ancestors {
		if _, exists := d.descendantsCache[ancestor]; exists {
			delete(d.descendantsCache, ancestor)
		}
	}
	delete(d.descendantsCache, src)


	d.muEdges.Unlock()

	return nil
}

// Delete an edge.
func (d *DAG) DeleteEdge(src Vertex, dst Vertex) error {

	// sanity checking
	if src == nil {
		return VertexNilError{}
	}
	if dst == nil {
		return VertexNilError{}
	}
	if _, ok := d.vertices[src]; !ok {
		return VertexUnknownError{src}
	}
	if _, ok := d.vertices[dst]; !ok {
		return VertexUnknownError{dst}
	}

	// test / compute edge nodes
	_, outboundExists := d.outboundEdge[src][dst]
	_, inboundExists := d.inboundEdge[dst][src]

	if !inboundExists || !outboundExists {
		return EdgeUnknownError{src, dst}
	}

	// get descendents and ancestors as they are now
	descendants, _ := d.GetDescendants(src)
	ancestors, _ := d.GetAncestors(dst)

	d.muEdges.Lock()

	// delete outbound
	if outboundExists {
		delete(d.outboundEdge[src], dst)
	}

	// delete inbound
	if inboundExists {
		delete(d.inboundEdge[dst], src)
	}

	// for src and all its descendants delete cached ancestors
	for descendant := range descendants {
		if _, exists := d.ancestorCache[descendant]; exists {
			delete(d.ancestorCache, descendant)
		}
	}
	delete(d.ancestorCache, src)

	// for dst and all its ancestors delete cached descendants
	for ancestor := range ancestors {
		if _, exists := d.descendantsCache[ancestor]; exists {
			delete(d.descendantsCache, ancestor)
		}
	}
	delete(d.descendantsCache, dst)

	d.muEdges.Unlock()

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
func (d *DAG) GetLeafs() map[Vertex]bool {
	leafs := make(map[Vertex]bool)
	for v := range d.vertices {
		dstIds, ok := d.outboundEdge[v]
		if !ok || len(dstIds) == 0 {
			leafs[v] = true
		}
	}
	return leafs
}

// Return all vertices without parents.
func (d *DAG) GetRoots() map[Vertex]bool {
	roots := make(map[Vertex]bool)
	for v := range d.vertices {
		srcIds, ok := d.inboundEdge[v]
		if !ok || len(srcIds) == 0 {
			roots[v] = true
		}
	}
	return roots
}

// Return all vertices.
func (d *DAG) GetVertices() map[Vertex]bool {
	return d.vertices
}

// Return all children of the given vertex.
func (d *DAG) GetChildren(v Vertex) (map[Vertex]bool, error) {

	// sanity checking
	if v == nil {
		return nil, VertexNilError{}
	}
	if _, ok := d.vertices[v]; !ok {
		return nil, VertexUnknownError{v}
	}

	return d.outboundEdge[v], nil
}

// Return all parents of the given vertex.
func (d *DAG) GetParents(v Vertex) (map[Vertex]bool, error) {
	// sanity checking
	if v == nil {
		return nil, VertexNilError{}
	}
	if _, ok := d.vertices[v]; !ok {
		return nil, VertexUnknownError{v}
	}

	return d.inboundEdge[v], nil
}

func (d *DAG) getAncestorsAux(v Vertex) map[Vertex]bool {
	d.ancestorCache[v] = make(map[Vertex]bool)
	if parents, ok := d.inboundEdge[v]; ok {
		for parent := range parents {
			if _, exists := d.ancestorCache[parent]; !exists {
				d.ancestorCache[parent] = d.getAncestorsAux(parent)
			}
			d.muEdges.Lock()
			for ancestor := range d.ancestorCache[parent] {
				d.ancestorCache[v][ancestor] = true

			}
			d.ancestorCache[v][parent] = true
			d.muEdges.Unlock()
		}
	}
	return d.ancestorCache[v]
}

// Return all Ancestors of the given vertex.
func (d *DAG) GetAncestors(v Vertex) (map[Vertex]bool, error) {
	// sanity checking
	if v == nil {
		return nil, VertexNilError{}
	}
	if _, ok := d.vertices[v]; !ok {
		return nil, VertexUnknownError{v}
	}

	if _, exists := d.ancestorCache[v]; !exists {
		return d.getAncestorsAux(v), nil
	}
	return d.ancestorCache[v], nil
}

func (d *DAG) getDescendantsAux(v Vertex) map[Vertex]bool {
	d.lockVertex(v)
	defer d.unlockVertex(v)
	if _, exists := d.descendantsCache[v]; exists {
		return d.descendantsCache[v]
	}
	vCache := make(map[Vertex]bool)
	if children, ok := d.outboundEdge[v]; ok {
		var waitGroup sync.WaitGroup
		waitGroup.Add(len(children))
		for child := range children {
			go func(child Vertex) {
				if _, exists := d.descendantsCache[child]; !exists {
					d.getDescendantsAux(child)
				}
				for descendant := range d.descendantsCache[child] {
					vCache[descendant] = true
				}
				vCache[child] = true
				waitGroup.Done()
			}(child)
		}
		waitGroup.Wait()
	}
	d.descendantsCache[v] = vCache
	return vCache
}

// Return all Descendants of the given vertex.
func (d *DAG) GetDescendants(v Vertex) (map[Vertex]bool, error) {
	// sanity checking
	if v == nil {
		return nil, VertexNilError{}
	}
	if _, ok := d.vertices[v]; !ok {
		return nil, VertexUnknownError{v}
	}
	return d.getDescendantsAux(v), nil
}

// Return a representation of the graph.
func (d *DAG) String() string {
	result := fmt.Sprintf("DAG Vertices: %d - Edges: %d\n", d.GetOrder(), d.GetSize())
	result += fmt.Sprintf("Vertices:\n")
	for k := range d.vertices {
		result += fmt.Sprintf("  %v\n", k.String())
	}
	result += fmt.Sprintf("Edges:\n")
	for v, children := range d.outboundEdge {
		for child := range children {
			result += fmt.Sprintf("  %s -> %s\n", v.String(), child.String())
		}
	}
	return result
}


type cMutex struct {
	mutex sync.Mutex
	count int
}




func (d *DAG) lockVertex(v Vertex) {

	// use the muVertices-mutex to flag that we are interested to lock v
	d.muVertices.Lock()

	// if there is no cMutex for v, create it
	if _, ok := d.verticesLocked[v]; !ok {
		d.verticesLocked[v] = new(cMutex)
	}

	// flag that we are interested in this cMutex (thus now one deletes it)
	d.verticesLocked[v].count++

	vMutex := d.verticesLocked[v].mutex

	// as the cMutex is there and we have flagged it, we can release the muVertices-mutex
	d.muVertices.Unlock()

	// and wait on the vertex specific mutex / lock it
	vMutex.Lock()
}

func (d *DAG) unlockVertex(v Vertex) {

	// use the muVertices-mutex to flag that we we want to release the lock for v
	d.muVertices.Lock()

	// unlock
	d.verticesLocked[v].mutex.Unlock()

	// only if the unlock succeeded, reduce the count
	d.verticesLocked[v].count--

	// if we where the last one interested in this cMutex
	if d.verticesLocked[v].count == 0 {
		delete(d.verticesLocked, v)
	}

	// release the muVertices-mutex
	d.muVertices.Unlock()
}
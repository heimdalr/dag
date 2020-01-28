// Package dag implements directed acyclic graphs (DAGs).
package dag

import (
	"fmt"
	"sync"
)

// Vertex is the interface to be implemented for the vertices of the DAG.
type Vertex interface {

	// Return a string representation of the vertex.
	String() string

	// Return the id of this vertex. This id must be unique and never change.
	Id() string
}

// DAG implements the data structure of the DAG.
type DAG struct {
	muDAG            sync.RWMutex
	vertices         map[Vertex]bool
	vertexIds        map[string]Vertex
	inboundEdge      map[Vertex]map[Vertex]bool
	outboundEdge     map[Vertex]map[Vertex]bool
	muCache          sync.RWMutex
	verticesLocked   *dMutex
	ancestorsCache   map[Vertex]map[Vertex]bool
	descendantsCache map[Vertex]map[Vertex]bool
}

// NewDAG creates / initializes a new DAG.
func NewDAG() *DAG {
	return &DAG{
		vertices:         make(map[Vertex]bool),
		vertexIds:        make(map[string]Vertex),
		inboundEdge:      make(map[Vertex]map[Vertex]bool),
		outboundEdge:     make(map[Vertex]map[Vertex]bool),
		verticesLocked:   newDMutex(),
		ancestorsCache:   make(map[Vertex]map[Vertex]bool),
		descendantsCache: make(map[Vertex]map[Vertex]bool),
	}
}

// AddVertex adds the vertex v to the DAG. AddVertex returns an error, if v is
// nil, v is already part of the graph, or the id of v is already part of the
// graph.
func (d *DAG) AddVertex(v Vertex) error {

	// sanity checking
	if v == nil {
		return VertexNilError{}
	}
	id := v.Id()
	d.muDAG.RLock()
	_, VertexExists := d.vertices[v]
	_, IdExists := d.vertexIds[id]
	d.muDAG.RUnlock()
	if VertexExists {
		return VertexDuplicateError{v}
	}
	if IdExists {
		return IdDuplicateError{v}
	}
	d.muDAG.Lock()
	d.vertices[v] = true
	d.vertexIds[id] = v
	d.muDAG.Unlock()

	return nil
}

// GetVertex returns a vertex by its id. GetVertex returns an error, if id is
// the empty string or unknown.
func (d *DAG) GetVertex(id string) (Vertex, error) {

	// sanity checking
	if id == "" {
		return nil, IdEmptyError{}
	}
	d.muDAG.RLock()
	v, IdExists := d.vertexIds[id]
	d.muDAG.RUnlock()
	if !IdExists {
		return nil, IdUnknownError{id}
	}

	return v, nil
}

// DeleteVertex deletes the vertex v. DeleteVertex also deletes all attached
// edges (inbound and outbound) as well as ancestor- and descendant-caches of
// related vertices. DeleteVertex returns an error, if v is nil or unknown.
func (d *DAG) DeleteVertex(v Vertex) error {

	if err := d.saneVertex(v); err != nil {
		return err
	}

	// get descendents and ancestors as they are now
	descendants, _ := d.GetDescendants(v)
	ancestors, _ := d.GetAncestors(v)

	d.muDAG.Lock()

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
		if _, exists := d.ancestorsCache[descendant]; exists {
			delete(d.ancestorsCache, descendant)
		}
	}
	delete(d.ancestorsCache, v)

	// for v and all its ancestors delete cached descendants
	for ancestor := range ancestors {
		if _, exists := d.descendantsCache[ancestor]; exists {
			delete(d.descendantsCache, ancestor)
		}
	}
	delete(d.descendantsCache, v)

	// delete v itself
	delete(d.vertices, v)
	delete(d.vertexIds, v.Id())

	d.muDAG.Unlock()

	return nil
}

// AddEdge adds an edge between src and dst. AddEdge returns an error, if src
// or dst are nil or if the edge would create a loop. AddEdge calls AddVertex,
// if src and/or dst are not yet known within the DAG.
func (d *DAG) AddEdge(src Vertex, dst Vertex) error {

	// sanity checking
	if src == nil || dst == nil {
		return VertexNilError{}
	}

	// check for equality
	if src == dst {
		return SrcDstEqualError{src, dst}
	}

	// ensure vertices
	if !d.vertices[src] {
		err := d.AddVertex(src)
		if err != nil {
			return err
		}
	}
	if !d.vertices[dst] {
		err := d.AddVertex(dst)
		if err != nil {
			return err
		}
	}

	// test / compute edge nodes and the edge itself
	d.muDAG.RLock()
	_, outboundExists := d.outboundEdge[src]
	_, inboundExists := d.inboundEdge[dst]
	edgeKnown := outboundExists && d.outboundEdge[src][dst] && inboundExists && d.inboundEdge[dst][src]
	d.muDAG.RUnlock()

	// if the edge is already known, there is nothing else to do
	if edgeKnown {
		return EdgeDuplicateError{src, dst}
	}

	// get descendents and ancestors as they are now
	descendants, _ := d.GetDescendants(dst)
	ancestors, _ := d.GetAncestors(src)

	// check for circles, iff desired
	if src == dst || descendants[src] {
		return EdgeLoopError{src, dst}
	}

	d.muDAG.Lock()

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
		if _, exists := d.ancestorsCache[descendant]; exists {
			delete(d.ancestorsCache, descendant)
		}
	}
	delete(d.ancestorsCache, dst)

	// for src and all its ancestors delete cached descendants
	for ancestor := range ancestors {
		if _, exists := d.descendantsCache[ancestor]; exists {
			delete(d.descendantsCache, ancestor)
		}
	}
	delete(d.descendantsCache, src)

	d.muDAG.Unlock()

	return nil
}

// DeleteEdge deletes an edge. DeleteEdge also deletes ancestor- and
// descendant-caches of related vertices. DeleteEdge returns an error, if src
// or dst are nil or unknown, or if there is no edge between src and dst.
func (d *DAG) DeleteEdge(src Vertex, dst Vertex) error {

	// sanity checking
	if src == nil || dst == nil {
		return VertexNilError{}
	}

	// check for equality
	if src == dst {
		return SrcDstEqualError{src, dst}
	}

	d.muDAG.RLock()
	if _, ok := d.vertices[src]; !ok {
		return VertexUnknownError{src}
	}
	if _, ok := d.vertices[dst]; !ok {
		return VertexUnknownError{dst}
	}

	// test / compute edge nodes
	_, outboundExists := d.outboundEdge[src][dst]
	_, inboundExists := d.inboundEdge[dst][src]
	d.muDAG.RUnlock()

	if !inboundExists || !outboundExists {
		return EdgeUnknownError{src, dst}
	}

	// get descendents and ancestors as they are now
	descendants, _ := d.GetDescendants(src)
	ancestors, _ := d.GetAncestors(dst)

	d.muDAG.Lock()

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
		if _, exists := d.ancestorsCache[descendant]; exists {
			delete(d.ancestorsCache, descendant)
		}
	}
	delete(d.ancestorsCache, src)

	// for dst and all its ancestors delete cached descendants
	for ancestor := range ancestors {
		if _, exists := d.descendantsCache[ancestor]; exists {
			delete(d.descendantsCache, ancestor)
		}
	}
	delete(d.descendantsCache, dst)

	d.muDAG.Unlock()

	return nil
}

// GetOrder returns the number of vertices in the graph.
func (d *DAG) GetOrder() int {
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()
	return len(d.vertices)
}

// GetSize returns the number of edges in the graph.
func (d *DAG) GetSize() int {
	count := 0
	d.muDAG.RLock()
	for _, value := range d.outboundEdge {
		count += len(value)
	}
	d.muDAG.RUnlock()
	return count
}

// GetLeafs returns all vertices without children.
func (d *DAG) GetLeafs() map[Vertex]bool {
	leafs := make(map[Vertex]bool)
	d.muDAG.RLock()
	for v := range d.vertices {
		dstIds, ok := d.outboundEdge[v]
		if !ok || len(dstIds) == 0 {
			leafs[v] = true
		}
	}
	d.muDAG.RUnlock()
	return leafs
}

// GetRoots returns all vertices without parents.
func (d *DAG) GetRoots() map[Vertex]bool {
	roots := make(map[Vertex]bool)
	d.muDAG.RLock()
	for v := range d.vertices {
		srcIds, ok := d.inboundEdge[v]
		if !ok || len(srcIds) == 0 {
			roots[v] = true
		}
	}
	d.muDAG.RUnlock()
	return roots
}

func copyMap(in map[Vertex]bool) map[Vertex]bool {
	out := make(map[Vertex]bool)
	for key, value := range in {
		out[key] = value
	}
	return out
}

// GetVertices returns all vertices.
func (d *DAG) GetVertices() map[Vertex]bool {
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()
	return copyMap(d.vertices)
}

// GetParents returns all parents of vertex v. GetParents returns an error,
// if v is nil or unknown.
func (d *DAG) GetParents(v Vertex) (map[Vertex]bool, error) {
	if err := d.saneVertex(v); err != nil {
		return nil, err
	}
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()
	return copyMap(d.inboundEdge[v]), nil
}

// GetChildren returns all children of vertex v. GetChildren returns an error,
// if v is nil or unknown.
func (d *DAG) GetChildren(v Vertex) (map[Vertex]bool, error) {

	if err := d.saneVertex(v); err != nil {
		return nil, err
	}
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()
	return copyMap(d.outboundEdge[v]), nil
}

// GetAncestors return all ancestors of the vertex v. GetAncestors returns an
// error, if v is nil or unknown.
//
// Note, in order to get the ancestors, GetAncestors populates the ancestor-
// cache as needed. Depending on order and size of the sub-graph of v this may
// take a long time and consume a lot of memory.
func (d *DAG) GetAncestors(v Vertex) (map[Vertex]bool, error) {

	if err := d.saneVertex(v); err != nil {
		return nil, err
	}
	return copyMap(d.getAncestors(v)), nil
}

func (d *DAG) getAncestors(v Vertex) map[Vertex]bool {

	// in the best case we have already a populated cache
	d.muCache.RLock()
	cache, exists := d.ancestorsCache[v]
	d.muCache.RUnlock()
	if exists {
		return cache
	}

	// lock this vertex to work on it exclusively
	d.verticesLocked.lock(v)
	defer d.verticesLocked.unlock(v)

	// now as we have locked this vertex, check (again) that no one has
	// meanwhile populated the cache
	d.muCache.RLock()
	cache, exists = d.ancestorsCache[v]
	d.muCache.RUnlock()
	if exists {
		return cache
	}

	// as there is no cache, we start from scratch and first of all collect
	// all ancestors locally
	cache = make(map[Vertex]bool)
	var mu sync.Mutex
	if parents, ok := d.inboundEdge[v]; ok {

		// for each parent collect its ancestors
		for parent := range parents {
			parentAncestors := d.getAncestors(parent)
			mu.Lock()
			for ancestor := range parentAncestors {
				cache[ancestor] = true
			}
			cache[parent] = true
			mu.Unlock()
		}
	}

	// remember the collected descendents
	d.muCache.Lock()
	d.ancestorsCache[v] = cache
	d.muCache.Unlock()
	return cache
}

// GetOrderedAncestors returns all ancestors of the vertex v in a breath-first
// order. Only the first occurrence of each vertex is returned.
// GetOrderedAncestors returns an error, if v is nil or unknown.
//
// Note, there is no order between sibling vertices. Two consecutive runs of
// GetOrderedAncestors may return different results.
func (d *DAG) GetOrderedAncestors(v Vertex) ([]Vertex, error) {

	vertices, _, err := d.AncestorsWalker(v)
	if err != nil {
		return nil, err
	}
	var ancestors []Vertex
	for v := range vertices {
		ancestors = append(ancestors, v)
	}
	return ancestors, nil
}

// AncestorsWalker returns a channel and subsequently returns / walks all
// ancestors of the vertex v in a breath first order. The second channel
// returned may be used to stop further walking. AncestorsWalker returns an
// error, if v is nil or unknown.
//
// Note, there is no order between sibling vertices. Two consecutive runs of
// AncestorsWalker may return different results.
func (d *DAG) AncestorsWalker(v Vertex) (chan Vertex, chan bool, error) {
	if err := d.saneVertex(v); err != nil {
		return nil, nil, err
	}
	vertices := make(chan Vertex)
	signal := make(chan bool, 1)
	go func() {
		d.muDAG.RLock()
		d.walkAncestors(v, vertices, signal)
		d.muDAG.RUnlock()
		close(vertices)
		close(signal)
	}()
	return vertices, signal, nil
}

func (d *DAG) walkAncestors(v Vertex, vertices chan Vertex, signal chan bool) {

	var fifo []Vertex
	visited := make(map[Vertex]bool)
	for parent := range d.inboundEdge[v] {
		visited[parent] = true
		fifo = append(fifo, parent)
	}
	for {
		if len(fifo) == 0 {
			return
		}
		top := fifo[0]
		fifo = fifo[1:]
		for parent := range d.inboundEdge[top] {
			if !visited[parent] {
				visited[parent] = true
				fifo = append(fifo, parent)
			}
		}
		select {
		case <-signal:
			return
		default:
			vertices <- top
		}
	}
}

// GetDescendants return all ancestors of the vertex v. GetDescendants returns
// an error, if v is nil or unknown.
//
// Note, in order to get the ancestors, GetDescendants populates the
// descendants-cache as needed. Depending on order and size of the sub-graph of
// v this may take a long time and consume a lot of memory.
func (d *DAG) GetDescendants(v Vertex) (map[Vertex]bool, error) {

	if err := d.saneVertex(v); err != nil {
		return nil, err
	}
	return copyMap(d.getDescendants(v)), nil
}

func (d *DAG) getDescendants(v Vertex) map[Vertex]bool {

	// in the best case we have already a populated cache
	d.muCache.RLock()
	cache, exists := d.descendantsCache[v]
	d.muCache.RUnlock()
	if exists {
		return cache
	}

	// lock this vertex to work on it exclusively
	d.verticesLocked.lock(v)
	defer d.verticesLocked.unlock(v)

	// now as we have locked this vertex, check (again) that no one has
	// meanwhile populated the cache
	d.muCache.RLock()
	cache, exists = d.descendantsCache[v]
	d.muCache.RUnlock()
	if exists {
		return cache
	}

	// as there is no cache, we start from scratch and first of all collect
	// all descendants locally
	cache = make(map[Vertex]bool)
	var mu sync.Mutex
	if children, ok := d.outboundEdge[v]; ok {

		// for each child use a goroutine to collect its descendants
		//var waitGroup sync.WaitGroup
		//waitGroup.Add(len(children))
		for child := range children {
			//go func(child Vertex, mu *sync.Mutex, cache map[Vertex]bool) {
			childDescendants := d.getDescendants(child)
			mu.Lock()
			for descendant := range childDescendants {
				cache[descendant] = true
			}
			cache[child] = true
			mu.Unlock()
			//waitGroup.Done()
			//}(child, &mu, cache)
		}
		//waitGroup.Wait()
	}

	// remember the collected descendents
	d.muCache.Lock()
	d.descendantsCache[v] = cache
	d.muCache.Unlock()
	return cache
}

// GetOrderedDescendants returns all descendants of the vertex v in a breath-
// first order. Only the first occurrence of each vertex is returned.
// GetOrderedDescendants returns an error, if v is nil or unknown.
//
// Note, there is no order between sibling vertices. Two consecutive runs of
// GetOrderedDescendants may return different results.
func (d *DAG) GetOrderedDescendants(v Vertex) ([]Vertex, error) {

	vertices, _, err := d.DescendantsWalker(v)
	if err != nil {
		return nil, err
	}
	var descendants []Vertex
	for v := range vertices {
		descendants = append(descendants, v)
	}
	return descendants, nil
}

// DescendantsWalker returns a channel and subsequently returns / walks all
// descendants of the vertex v in a breath first order. The second channel
// returned may be used to stop further walking. DescendantsWalker returns an
// error, if v is nil or unknown.
//
// Note, there is no order between sibling vertices. Two consecutive runs of
// DescendantsWalker may return different results.
func (d *DAG) DescendantsWalker(v Vertex) (chan Vertex, chan bool, error) {
	if err := d.saneVertex(v); err != nil {
		return nil, nil, err
	}
	vertices := make(chan Vertex)
	signal := make(chan bool, 1)
	go func() {
		d.muDAG.RLock()
		d.walkDescendants(v, vertices, signal)
		d.muDAG.RUnlock()
		close(vertices)
		close(signal)
	}()
	return vertices, signal, nil
}

func (d *DAG) walkDescendants(v Vertex, vertices chan Vertex, signal chan bool) {

	var fifo []Vertex
	visited := make(map[Vertex]bool)
	for child := range d.outboundEdge[v] {
		visited[child] = true
		fifo = append(fifo, child)
	}
	for {
		if len(fifo) == 0 {
			return
		}
		top := fifo[0]
		fifo = fifo[1:]
		for child := range d.outboundEdge[top] {
			if !visited[child] {
				visited[child] = true
				fifo = append(fifo, child)
			}
		}
		select {
		case <-signal:
			return
		default:
			vertices <- top
		}
	}
}

// String return a textual representation of the graph.
func (d *DAG) String() string {
	result := fmt.Sprintf("DAG Vertices: %d - Edges: %d\n", d.GetOrder(), d.GetSize())
	result += fmt.Sprintf("Vertices:\n")
	d.muDAG.RLock()
	for k := range d.vertices {
		result += fmt.Sprintf("  %v\n", k.String())
	}
	result += fmt.Sprintf("Edges:\n")
	for v, children := range d.outboundEdge {
		for child := range children {
			result += fmt.Sprintf("  %s -> %s\n", v.String(), child.String())
		}
	}
	d.muDAG.RUnlock()
	return result
}

func (d *DAG) saneVertex(v Vertex) error {
	// sanity checking
	if v == nil {
		return VertexNilError{}
	}
	d.muDAG.RLock()
	_, exists := d.vertices[v]
	d.muDAG.RUnlock()
	if !exists {
		return VertexUnknownError{v}
	}
	return nil
}

/***************************
********** Errors **********
****************************/

// VertexNilError is the error type to describe the situation, that a nil is
// given instead of a vertex.
type VertexNilError struct{}

// Implements the error interface.
func (e VertexNilError) Error() string {
	return fmt.Sprint("don't know what to do with 'nil'")
}

// IdEmptyError is the error type to describe the situation, that a nil is
// given instead of a vertex.
type IdEmptyError struct{}

// Implements the error interface.
func (e IdEmptyError) Error() string {
	return fmt.Sprint("don't know what to do with 'nil'")
}

// VertexDuplicateError is the error type to describe the situation, that a
// given vertex already exists in the graph.
type VertexDuplicateError struct {
	v Vertex
}

// Implements the error interface.
func (e VertexDuplicateError) Error() string {
	return fmt.Sprintf("'%s' is already known", e.v.String())
}

// IdDuplicateError is the error type to describe the situation, that a given
// vertex id already exists in the graph.
type IdDuplicateError struct {
	v Vertex
}

// Implements the error interface.
func (e IdDuplicateError) Error() string {
	return fmt.Sprintf("the id '%s' is already known", e.v.Id())
}

// VertexUnknownError is the error type to describe the situation, that a given
// vertex does not exit in the graph.
type VertexUnknownError struct {
	v Vertex
}

// Implements the error interface.
func (e VertexUnknownError) Error() string {
	return fmt.Sprintf("'%s' is unknown", e.v.String())
}

// IdUnknownError is the error type to describe the situation, that a given
// vertex does not exit in the graph.
type IdUnknownError struct {
	id string
}

// Implements the error interface.
func (e IdUnknownError) Error() string {
	return fmt.Sprintf("'%s' is unknown", e.id)
}

// EdgeDuplicateError is the error type to describe the situation, that an edge
// already exists in the graph.
type EdgeDuplicateError struct {
	src Vertex
	dst Vertex
}

// Implements the error interface.
func (e EdgeDuplicateError) Error() string {
	return fmt.Sprintf("edge between '%s' and '%s' is already known", e.src.String(), e.dst.String())
}

// EdgeUnknownError is the error type to describe the situation, that a given
// edge does not exit in the graph.
type EdgeUnknownError struct {
	src Vertex
	dst Vertex
}

// Implements the error interface.
func (e EdgeUnknownError) Error() string {
	return fmt.Sprintf("edge between '%s' and '%s' is unknown", e.src.String(), e.dst.String())
}

// EdgeLoopError is the error type to describe loop errors (i.e. errors that
// where raised to prevent establishing loops in the graph).
type EdgeLoopError struct {
	src Vertex
	dst Vertex
}

// Implements the error interface.
func (e EdgeLoopError) Error() string {
	return fmt.Sprintf("edge between '%s' and '%s' would create a loop", e.src.String(), e.dst.String())
}

// SrcDstEqualError is the error type to describe the situation, that src and
// dst are equal.
type SrcDstEqualError struct {
	src Vertex
	dst Vertex
}

// Implements the error interface.
func (e SrcDstEqualError) Error() string {
	return fmt.Sprintf("src ('%s') and dst ('%s') equal", e.src.String(), e.dst.String())
}

/***************************
********** dMutex **********
****************************/

type cMutex struct {
	mutex sync.Mutex
	count int
}

// Structure for dynamic mutexes.
type dMutex struct {
	mutexes     map[interface{}]*cMutex
	globalMutex sync.Mutex
}

// Initialize a new dynamic mutex structure.
func newDMutex() *dMutex {
	return &dMutex{
		mutexes: make(map[interface{}]*cMutex),
	}
}

// Get a lock for instance i
func (d *dMutex) lock(i interface{}) {

	// acquire global lock
	d.globalMutex.Lock()

	// if there is no cMutex for i, create it
	if _, ok := d.mutexes[i]; !ok {
		d.mutexes[i] = new(cMutex)
	}

	// increase the count in order to show, that we are interested in this
	// instance mutex (thus now one deletes it)
	d.mutexes[i].count++

	// remember the mutex for later
	mutex := &d.mutexes[i].mutex

	// as the cMutex is there, we have increased the count and we know the
	// instance mutex, we can release the global lock
	d.globalMutex.Unlock()

	// and wait on the instance mutex
	(*mutex).Lock()
}

// Release the lock for instance i.
func (d *dMutex) unlock(i interface{}) {

	// acquire global lock
	d.globalMutex.Lock()

	// unlock instance mutex
	d.mutexes[i].mutex.Unlock()

	// decrease the count, as we are no longer interested in this instance
	// mutex
	d.mutexes[i].count--

	// if we where the last one interested in this instance mutex delete the
	// cMutex
	if d.mutexes[i].count == 0 {
		delete(d.mutexes, i)
	}

	// release the global lock
	d.globalMutex.Unlock()
}

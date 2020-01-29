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

	d.muDAG.Lock()
	defer d.muDAG.Unlock()

	// sanity checking
	if v == nil {
		return VertexNilError{}
	}
	if _, exists := d.vertices[v]; exists {
		return VertexDuplicateError{v}
	}
	if _, exists := d.vertexIds[v.Id()]; exists {
		return IdDuplicateError{v}
	}
	d.addVertex(v)
	return nil
}

func (d *DAG) addVertex(v Vertex) {
	d.vertices[v] = true
	d.vertexIds[v.Id()] = v
}

// GetVertex returns a vertex by its id. GetVertex returns an error, if id is
// the empty string or unknown.
func (d *DAG) GetVertex(id string) (Vertex, error) {
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()

	if id == "" {
		return nil, IdEmptyError{}
	}

	v, IdExists := d.vertexIds[id]
	if !IdExists {
		return nil, IdUnknownError{id}
	}
	return v, nil
}

// DeleteVertex deletes the vertex v. DeleteVertex also deletes all attached
// edges (inbound and outbound) as well as ancestor- and descendant-caches of
// related vertices. DeleteVertex returns an error, if v is nil or unknown.
func (d *DAG) DeleteVertex(v Vertex) error {

	d.muDAG.Lock()
	defer d.muDAG.Unlock()

	if err := d.saneVertex(v); err != nil {
		return err
	}

	// get descendents and ancestors as they are now
	descendants := copyMap(d.getDescendants(v))
	ancestors := copyMap(d.getAncestors(v))

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

	return nil
}

// AddEdge adds an edge between src and dst. AddEdge returns an error, if src
// or dst are nil or if the edge would create a loop. AddEdge calls AddVertex,
// if src and/or dst are not yet known within the DAG.
func (d *DAG) AddEdge(src Vertex, dst Vertex) error {

	d.muDAG.Lock()
	defer d.muDAG.Unlock()

	if src == nil || dst == nil {
		return VertexNilError{}
	}
	if src == dst {
		return SrcDstEqualError{src, dst}
	}

	// ensure vertices
	if !d.vertices[src] {
		if _, idExists := d.vertexIds[src.Id()]; idExists {
			return IdDuplicateError{src}
		}
		d.addVertex(src)
	}
	if !d.vertices[dst] {
		if _, idExists := d.vertexIds[dst.Id()]; idExists {
			return IdDuplicateError{dst}
		}
		d.addVertex(dst)
	}

	// if the edge is already known, there is nothing else to do
	if d.isEdge(src, dst) {
		return EdgeDuplicateError{src, dst}
	}

	// get descendents and ancestors as they are now
	descendants := copyMap(d.getDescendants(dst))
	ancestors := copyMap(d.getAncestors(src))

	// check for circles, iff desired
	if src == dst || descendants[src] {
		return EdgeLoopError{src, dst}
	}

	// prepare d.outbound[src], iff needed
	if _, exists := d.outboundEdge[src]; !exists {
		d.outboundEdge[src] = make(map[Vertex]bool)
	}

	// dst is a child of src
	d.outboundEdge[src][dst] = true

	// prepare d.inboundEdge[dst], iff needed
	if _, exists := d.inboundEdge[dst]; !exists {
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

	return nil
}

// IsEdge returns true, if there exists an edge between src and dst. IsEdge
// returns false if there is no such edge. IsEdge returns an error, if src or
// dst are nil, unknown or the same.
func (d *DAG) IsEdge(src Vertex, dst Vertex) (bool, error) {
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()

	if err := d.saneVertex(src); err != nil {
		return false, err
	}
	if err := d.saneVertex(dst); err != nil {
		return false, err
	}
	if src == dst {
		return false, SrcDstEqualError{src, dst}
	}

	return d.isEdge(src, dst), nil
}

func (d *DAG) isEdge(src Vertex, dst Vertex) bool {

	_, outboundExists := d.outboundEdge[src]
	_, inboundExists := d.inboundEdge[dst]

	return outboundExists && d.outboundEdge[src][dst] &&
		inboundExists && d.inboundEdge[dst][src]
}

// DeleteEdge deletes an edge. DeleteEdge also deletes ancestor- and
// descendant-caches of related vertices. DeleteEdge returns an error, if src
// or dst are nil or unknown, or if there is no edge between src and dst.
func (d *DAG) DeleteEdge(src Vertex, dst Vertex) error {

	d.muDAG.Lock()
	defer d.muDAG.Unlock()

	if err := d.saneVertex(src); err != nil {
		return err
	}
	if err := d.saneVertex(dst); err != nil {
		return err
	}
	if src == dst {
		return SrcDstEqualError{src, dst}
	}
	if !d.isEdge(src, dst) {
		return EdgeUnknownError{src, dst}
	}

	// get descendents and ancestors as they are now
	descendants := copyMap(d.getDescendants(src))
	ancestors := copyMap(d.getAncestors(dst))

	// delete outbound and inbound
	delete(d.outboundEdge[src], dst)
	delete(d.inboundEdge[dst], src)

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
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()
	return d.getSize()
}

func (d *DAG) getSize() int {
	count := 0
	for _, value := range d.outboundEdge {
		count += len(value)
	}
	return count
}

// GetLeafs returns all vertices without children.
func (d *DAG) GetLeafs() map[Vertex]bool {
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()
	return d.getLeafs()
}

func (d *DAG) getLeafs() map[Vertex]bool {
	leafs := make(map[Vertex]bool)
	for v := range d.vertices {
		dstIds, ok := d.outboundEdge[v]
		if !ok || len(dstIds) == 0 {
			leafs[v] = true
		}
	}
	return leafs
}

// GetRoots returns all vertices without parents.
func (d *DAG) GetRoots() map[Vertex]bool {
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()
	return d.getRoots()
}

func (d *DAG) getRoots() map[Vertex]bool {
	roots := make(map[Vertex]bool)
	for v := range d.vertices {
		srcIds, ok := d.inboundEdge[v]
		if !ok || len(srcIds) == 0 {
			roots[v] = true
		}
	}
	return roots
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
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()
	if err := d.saneVertex(v); err != nil {
		return nil, err
	}
	return copyMap(d.inboundEdge[v]), nil
}

// GetChildren returns all children of vertex v. GetChildren returns an error,
// if v is nil or unknown.
func (d *DAG) GetChildren(v Vertex) (map[Vertex]bool, error) {
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()
	if err := d.saneVertex(v); err != nil {
		return nil, err
	}
	return copyMap(d.outboundEdge[v]), nil
}

// GetAncestors return all ancestors of the vertex v. GetAncestors returns an
// error, if v is nil or unknown.
//
// Note, in order to get the ancestors, GetAncestors populates the ancestor-
// cache as needed. Depending on order and size of the sub-graph of v this may
// take a long time and consume a lot of memory.
func (d *DAG) GetAncestors(v Vertex) (map[Vertex]bool, error) {
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()
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
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()
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
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()
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
// Note, in order to get the descendants, GetDescendants populates the
// descendants-cache as needed. Depending on order and size of the sub-graph of
// v this may take a long time and consume a lot of memory.
func (d *DAG) GetDescendants(v Vertex) (map[Vertex]bool, error) {
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()
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
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()
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
	d.muDAG.RLock()
	defer d.muDAG.RUnlock()
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

// ReduceTransitively transitively reduce the graph.
//
// Note, in order to do the reduction the descendant-cache of all vertices is
// populated (i.e. the transitive closure). Depending on order and size of DAG
// this may take a long time and consume a lot of memory.
func (d *DAG) ReduceTransitively() {

	d.muDAG.Lock()
	defer d.muDAG.Unlock()

	graphChanged := false

	// populate the descendents cache for all roots (i.e. the whole graph)
	for root := range d.getRoots() {
		_ = d.getDescendants(root)
	}

	// for each vertex
	for v := range d.vertices {

		// map of descendants of the children of v
		descendentsOfChildrenOfV := make(map[Vertex]bool)

		// for each child of v
		for childOfV := range d.outboundEdge[v] {

			// collect child descendants
			for descendent := range d.descendantsCache[childOfV] {
				descendentsOfChildrenOfV[descendent] = true
			}
		}

		// for each child of v
		for childOfV := range d.outboundEdge[v] {

			// remove the edge between v and child, iff child is a
			// descendants of any of the children of v
			if descendentsOfChildrenOfV[childOfV] {
				delete(d.outboundEdge[v], childOfV)
				delete(d.inboundEdge[childOfV], v)
				graphChanged = true
			}
		}
	}

	// flush the descendants- and ancestor cache if the graph has changed
	if graphChanged {
		d.flushCaches()
	}
}

// FlushCaches completely flushes the descendants- and ancestor cache.
//
// Note, the only reason to call this method is to free up memory.
// Otherwise the caches are automatically maintained.
func (d *DAG) FlushCaches() {
	d.muDAG.Lock()
	defer d.muDAG.Unlock()
	d.flushCaches()
}

func (d *DAG) flushCaches() {
	d.ancestorsCache = make(map[Vertex]map[Vertex]bool)
	d.descendantsCache = make(map[Vertex]map[Vertex]bool)
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
	_, exists := d.vertices[v]
	if !exists {
		return VertexUnknownError{v}
	}
	return nil
}

func copyMap(in map[Vertex]bool) map[Vertex]bool {
	out := make(map[Vertex]bool)
	for key, value := range in {
		out[key] = value
	}
	return out
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

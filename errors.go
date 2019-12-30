package dag

import "fmt"

// Error type to describe the situation, that a nil is given instead of a vertex.
type VertexNilError struct{}

// Implements the error interface.
func (e VertexNilError) Error() string {
	return fmt.Sprint("don't know what to do with 'nil'")
}

// Error type to describe the situation, that a given vertex already exists in the graph.
type VertexDuplicateError struct {
	v Vertex
}

// Implements the error interface.
func (e VertexDuplicateError) Error() string {
	return fmt.Sprintf("'%s' is already known", e.v.String())
}

// Error type to describe the situation, that a given vertex does not exit in the graph.
type VertexUnknownError struct {
	v Vertex
}

// Implements the error interface.
func (e VertexUnknownError) Error() string {
	return fmt.Sprintf("'%s' is unknown", e.v.String())
}

// Error type to describe the situation, that an edge already exists in the graph.
type EdgeDuplicateError struct {
	src Vertex
	dst Vertex
}

// Implements the error interface.
func (e EdgeDuplicateError) Error() string {
	return fmt.Sprintf("edge between '%s' and '%s' is already known", e.src.String(), e.dst.String())
}

// Error type to describe the situation, that a given edge does not exit in the graph.
type EdgeUnknownError struct {
	src Vertex
	dst Vertex
}

// Implements the error interface.
func (e EdgeUnknownError) Error() string {
	return fmt.Sprintf("edge between '%s' and '%s' is unknown", e.src.String(), e.dst.String())
}

// Error type to describe loop errors (i.e. errors that where raised to prevent establishing loops in the graph).
type EdgeLoopError struct {
	src Vertex
	dst Vertex
}

// Implements the error interface.
func (e EdgeLoopError) Error() string {
	return fmt.Sprintf("edge between '%s' and '%s' would create a loop", e.src.String(), e.dst.String())
}

package dag

// Options is the configuration for the DAG.
type Options struct {
	// VertexHashFunc is the function that calculates the hash value of a vertex.
	// This can be useful when the vertex contains not comparable types such as maps.
	// If VertexHashFunc is nil, the defaultVertexHashFunc is used.
	VertexHashFunc func(v interface{}) interface{}
}

func defaultVertexHashFunc(v interface{}) interface{} {
	return v
}

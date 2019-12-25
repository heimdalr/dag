# dag

<!--[![Build Status](https://travis-ci.org/go-redis/redis.png?branch=master)](https://travis-ci.org/go-redis/redis)-->
[![GoDoc](https://godoc.org/github.com/sebogh/dag?status.svg)](https://godoc.org/github.com/sebogh/dag) 

Yet another directed acyclic graph (DAG) implementation in golang.

## Quickstart

``` go

// data structure that will be used as vertex in the graph
type myVertex struct {
	Label string
}

// implement the Vertex interface 
func (v myVertex) String() string {
	return v.Label
}


func main() {

    // initialize a new graph
    dag := NewDAG()

    // init three vertices
    var v1 Vertex = myVertex{"1"}
    var v2 Vertex = myVertex{"2"}
    var v3 Vertex = myVertex{"3"}

    // add the above vertices and connect them with two edges
	_ = dag.AddEdge(&v1, &v2)
	_ = dag.AddEdge(&v1, &v3)

    // get children of &v1 -> set (map to bool) containing &v2 and &v3
    children, _ := dag.GetChildren(&v1)
}
```
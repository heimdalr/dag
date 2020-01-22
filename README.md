# dag

[![Build Status](https://travis-ci.org/heimdalr/dag.svg?branch=master)](https://travis-ci.org/heimdalr/dag)
[![codecov](https://codecov.io/gh/heimdalr/dag/branch/master/graph/badge.svg)](https://codecov.io/gh/heimdalr/dag)
[![GoDoc](https://godoc.org/github.com/heimdalr/dag?status.svg)](https://godoc.org/github.com/heimdalr/dag) 
[![Go Report Card](https://goreportcard.com/badge/github.com/heimdalr/dag)](https://goreportcard.com/report/github.com/heimdalr/dag)

Yet another directed acyclic graph (DAG) implementation in golang.

## Quickstart

Running: 

``` go
package main

import (
	"fmt"
	"github.com/heimdalr/dag"
)

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
	d := dag.NewDAG()

	// init three vertices
	v1 := &myVertex{"1"}
	v2 := &myVertex{"2"}
	v3 := &myVertex{"3"}

	// add the above vertices and connect them with two edges
	_ = d.AddEdge(v1, v2)
	_ = d.AddEdge(v1, v3)

	// describe the graph
	fmt.Print(d.String())
}
```

will result in something like:

```
DAG Vertices: 3 - Edges: 2
Vertices:
  2
  3
  1
Edges:
  1 -> 2
  1 -> 3
```

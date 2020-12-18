# dag

[![circleci](https://circleci.com/gh/heimdalr/dag/tree/master.svg?style=shield)](https://circleci.com/gh/heimdalr/dag)
[![codecov](https://codecov.io/gh/heimdalr/dag/branch/master/graph/badge.svg)](https://codecov.io/gh/heimdalr/dag)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/heimdalr/dag)](https://pkg.go.dev/github.com/heimdalr/dag)
[![Go Report Card](https://goreportcard.com/badge/github.com/heimdalr/dag)](https://goreportcard.com/report/github.com/heimdalr/dag)
[![Nutrition Facts](http://code.grevit.net/badge/O%2B%2B_S%2B%2B_I%2B_C%2B_E%2B%2B%2B_M_V%2B_PS%2B%2B_!D)](http://code.grevit.net/facts/O%2B%2B_S%2B%2B_I%2B_C%2B_E%2B%2B%2B_M_V%2B_PS%2B%2B_!D)

Implementation of directed acyclic graphs (DAGs).

The implementation is fast and thread-safe. It prevents adding cycles or 
duplicates and thereby always maintains a valid DAG. The implementation caches
 descendants and ancestors to speed up subsequent calls. 

<!--
github.com/heimdalr/dag:

3.770388s to add 597871 vertices and 597870 edges
1.578741s to get descendants
0.143887s to get descendants 2nd time
0.444065s to get descendants ordered
0.000008s to get children
1.301297s to transitively reduce the graph with caches poupulated
2.723708s to transitively reduce the graph without caches poupulated
0.168572s to delete an edge from the root


"github.com/hashicorp/terraform/dag":

3.195338s to add 597871 vertices and 597870 edges
1.121812s to get descendants
1.803096s to get descendants 2nd time
3.056972s to transitively reduce the graph
-->




## Quickstart

Running: 

``` go
package main

import (
	"fmt"
	"github.com/heimdalr/dag"
)

func main() {

	// initialize a new graph
	d := NewDAG()

	// init three vertices
	v1, _ := d.AddVertex(1)
	v2, _ := d.AddVertex(2)
	v3, _ := d.AddVertex(struct{a string; b string}{a: "foo", b: "bar"})

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
  1
  2
  {foo bar}
Edges:
  1 -> 2
  1 -> {foo bar}
```

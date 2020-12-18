package main

/*
import (
	"fmt"
	"github.com/hashicorp/terraform/dag"
	"math"
	"time"
)

type largeVertex struct {
	value int
}

// implement the Vertex's interface method String()
func (v largeVertex) String() string {
	return fmt.Sprintf("%d", v.value)
}

// implement the Vertex's interface method ID()
func (v largeVertex) ID() string {
	return fmt.Sprintf("%d", v.value)
}
*/

func main() {
	/*
		var d dag.AcyclicGraph

		root := d.Add(1)
		levels := 7
		branches := 9
		var start, end time.Time

		start = time.Now()
		largeAux(d, levels, branches, root)
		_ = d.Validate()
		end = time.Now()
		fmt.Printf("%fs to add %d vertices and %d edges\n", end.Sub(start).Seconds(), len(d.Vertices()), len(d.Edges()))
		expectedVertexCount := sum(0, levels-1, branches, pow)
		vertexCount := len(d.Vertices())
		if vertexCount != expectedVertexCount {
			panic(fmt.Sprintf("GetVertices() = %d, want %d", vertexCount, expectedVertexCount))
		}

		start = time.Now()
		descendants, _ := d.Descendents(root)
		end = time.Now()
		fmt.Printf("%fs to get descendants\n", end.Sub(start).Seconds())
		descendantsCount := descendants.Len()
		expectedDescendantsCount := vertexCount - 1
		if descendantsCount != expectedDescendantsCount {
			panic(fmt.Sprintf("GetDescendants(root) = %d, want %d", descendantsCount, expectedDescendantsCount))
		}

		start = time.Now()
		_, _ = d.Descendents(root)
		end = time.Now()
		fmt.Printf("%fs to get descendants 2nd time\n", end.Sub(start).Seconds())

		start = time.Now()
		d.TransitiveReduction()
		end = time.Now()
		fmt.Printf("%fs to transitively reduce the graph\n", end.Sub(start).Seconds())
	*/
}

/*
func largeAux(d dag.AcyclicGraph, level int, branches int, parent dag.Vertex) {
	if level > 1 {
		if branches < 1 || branches > 9 {
			panic("number of branches must be between 1 and 9")
		}
		for i := 1; i <= branches; i++ {
			value := parent.(int)*10 + i
			child := d.Add(value)
			d.Connect(dag.BasicEdge(child, parent))
			largeAux(d, level-1, branches, child)
		}
	}
}

func sum(x, y, branches int, fn interface{}) int {
	if x > y {
		return 0
	}
	f, ok := fn.(func(int, int) int)
	if !ok {
		panic("function no of correct tpye")
	}
	current := f(branches, x)
	rest := sum(x+1, y, branches, f)
	return current + rest
}

func pow(base int, exp int) int {
	pow := math.Pow(float64(base), float64(exp))
	return int(pow)
}*/

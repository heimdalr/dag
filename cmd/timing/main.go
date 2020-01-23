package main

import (
	"fmt"
	"github.com/heimdalr/dag"
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

// implement the Vertex's interface method Id()
func (v largeVertex) Id() string {
	return fmt.Sprintf("%d", v.value)
}

func main() {
	d := dag.NewDAG()
	root := &largeVertex{1}
	levels := 7
	branches := 9
	var start, end time.Time

	start = time.Now()
	largeAux(d, levels, branches, root)
	end = time.Now()
	fmt.Printf("%fs to add %d vertices and %d edges\n", end.Sub(start).Seconds(), d.GetOrder(), d.GetSize())
	expectedVertexCount := sum(0, levels-1, branches, pow)
	vertexCount := len(d.GetVertices())
	if vertexCount != expectedVertexCount {
		panic(fmt.Sprintf("GetVertices() = %d, want %d", vertexCount, expectedVertexCount))
	}

	start = time.Now()
	descendants, _ := d.GetDescendants(root)
	end = time.Now()
	fmt.Printf("%fs to get descendants\n", end.Sub(start).Seconds())
	descendantsCount := len(descendants)
	expectedDescendantsCount := vertexCount - 1
	if descendantsCount != expectedDescendantsCount {
		panic(fmt.Sprintf("GetDescendants(root) = %d, want %d", descendantsCount, expectedDescendantsCount))
	}

	start = time.Now()
	descendantsOrdered, _ := d.GetOrderedDescendants(root)
	end = time.Now()
	fmt.Printf("%fs to get descendants ordered\n", end.Sub(start).Seconds())
	descendantsOrderedCount := len(descendantsOrdered)
	if descendantsOrderedCount != expectedDescendantsCount {
		panic(fmt.Sprintf("GetOrderedDescendants(root) = %d, want %d", descendantsOrderedCount, expectedDescendantsCount))
	}

	start = time.Now()
	_, _ = d.GetDescendants(root)
	end = time.Now()
	fmt.Printf("%fs to get descendants 2nd time\n", end.Sub(start).Seconds())

	start = time.Now()
	children, _ := d.GetChildren(root)
	end = time.Now()
	fmt.Printf("%fs to get children\n", end.Sub(start).Seconds())
	childrenCount := len(children)
	expectedChildrenCount := branches
	if childrenCount != expectedChildrenCount {
		panic(fmt.Sprintf("GetChildren(root) = %d, want %d", childrenCount, expectedChildrenCount))
	}

	var childList []dag.Vertex
	for x := range children {
		childList = append(childList, x)
	}
	start = time.Now()
	if len(childList) > 0 {
		_ = d.DeleteEdge(root, childList[0])
	}
	end = time.Now()
	fmt.Printf("%fs to delete an edge from the root\n", end.Sub(start).Seconds())

}

func largeAux(d *dag.DAG, level int, branches int, parent *largeVertex) {
	if level > 1 {
		if branches < 1 || branches > 9 {
			panic("number of branches must be between 1 and 9")
		}
		for i := 1; i <= branches; i++ {
			value := (*parent).value*10 + i
			child := &largeVertex{value}
			err := d.AddEdge(parent, child)
			if err != nil {
				panic(err)
			}
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
}

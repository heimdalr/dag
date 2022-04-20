package dag

import (
	"encoding/json"
	"testing"

	"github.com/go-test/deep"
)

func getTestMarshalDAG() *DAG {
	dag := NewDAG()
	v1, v2, v3, v4, v5 := "1", "2", "3", "4", "5"
	_ = dag.AddVertexByID(v1, "v1")
	_ = dag.AddVertexByID(v2, "v2")
	_ = dag.AddVertexByID(v3, "v3")
	_ = dag.AddVertexByID(v4, "v4")
	_ = dag.AddVertexByID(v5, "v5")
	_ = dag.AddEdge(v1, v2)
	_ = dag.AddEdge(v2, v3)
	_ = dag.AddEdge(v2, v4)
	_ = dag.AddEdge(v4, v5)
	return dag
}

func TestMarshalUnmarshalJSON(t *testing.T) {
	d := getTestMarshalDAG()
	data, err := json.Marshal(d)
	if err != nil {
		t.Error(err)
	}

	expected := `{"vs":[{"i":"1","v":"v1"},{"i":"2","v":"v2"},{"i":"3","v":"v3"},{"i":"4","v":"v4"},{"i":"5","v":"v5"}],"es":[{"s":"1","d":"2"},{"s":"2","d":"3"},{"s":"2","d":"4"},{"s":"4","d":"5"}]}`
	actual := string(data)
	if deep.Equal(expected, actual) != nil {
		t.Errorf("Marshal() = %v, want %v", actual, expected)
	}

	d1 := &DAG{}
	errNotSupported := json.Unmarshal(data, d1)
	if errNotSupported == nil {
		t.Errorf("UnmarshalJSON() = nil, want %v", "This method is not supported")
	}

	var wd testStorableDAG
	dag, err := UnmarshalJSON(data, &wd)
	if err != nil {
		t.Fatal(err)
	}
	if deep.Equal(d, dag) != nil {
		t.Errorf("UnmarshalJSON() = %v, want %v", dag.String(), d.String())
	}
}

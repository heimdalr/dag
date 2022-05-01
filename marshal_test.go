package dag

import (
	"encoding/json"
	"testing"

	"github.com/go-test/deep"
)

func TestMarshalUnmarshalJSON(t *testing.T) {
	cases := []struct {
		dag      *DAG
		expected string
	}{
		{
			dag:      getTestWalkDAG(),
			expected: `{"vs":[{"i":"1","v":"v1"},{"i":"2","v":"v2"},{"i":"3","v":"v3"},{"i":"4","v":"v4"},{"i":"5","v":"v5"}],"es":[{"s":"1","d":"2"},{"s":"2","d":"3"},{"s":"2","d":"4"},{"s":"4","d":"5"}]}`,
		},
		{
			dag:      getTestWalkDAG2(),
			expected: `{"vs":[{"i":"1","v":"v1"},{"i":"3","v":"v3"},{"i":"5","v":"v5"},{"i":"2","v":"v2"},{"i":"4","v":"v4"}],"es":[{"s":"1","d":"3"},{"s":"3","d":"5"},{"s":"2","d":"3"},{"s":"4","d":"5"}]}`,
		},
		{
			dag:      getTestWalkDAG3(),
			expected: `{"vs":[{"i":"1","v":"v1"},{"i":"3","v":"v3"},{"i":"2","v":"v2"},{"i":"4","v":"v4"},{"i":"5","v":"v5"}],"es":[{"s":"1","d":"3"},{"s":"2","d":"3"},{"s":"4","d":"5"}]}`,
		},
		{
			dag:      getTestWalkDAG4(),
			expected: `{"vs":[{"i":"1","v":"v1"},{"i":"2","v":"v2"},{"i":"3","v":"v3"},{"i":"5","v":"v5"},{"i":"4","v":"v4"}],"es":[{"s":"1","d":"2"},{"s":"2","d":"3"},{"s":"2","d":"4"},{"s":"3","d":"5"}]}`,
		},
	}

	for _, c := range cases {
		testMarshalUnmarshalJSON(t, c.dag, c.expected)
	}
}

func testMarshalUnmarshalJSON(t *testing.T, d *DAG, expected string) {
	data, err := json.Marshal(d)
	if err != nil {
		t.Error(err)
	}

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

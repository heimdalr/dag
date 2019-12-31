package dag

import (
	"fmt"
	"testing"
)

func TestErrors(t *testing.T) {
	v1 := &testVertex{"1"}
	v2 := &testVertex{"2"}

	tests := []struct {
		want string
		err error
	}{
		{"don't know what to do with 'nil'", VertexNilError{}},
		{"'1' is already known", VertexDuplicateError{v1}},
		{"'1' is unknown", VertexUnknownError{v1}},
		{"edge between '1' and '2' is already known", EdgeDuplicateError{v1, v2}},
		{"edge between '1' and '2' is unknown", EdgeUnknownError{v1, v2}},
		{"edge between '1' and '2' would create a loop", EdgeLoopError{v1, v2}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt.err), func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
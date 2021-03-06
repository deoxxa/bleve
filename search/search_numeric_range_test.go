package search

import (
	"reflect"
	"testing"

	"github.com/couchbaselabs/bleve/numeric_util"
)

func TestSplitRange(t *testing.T) {
	min := numeric_util.Float64ToInt64(1.0)
	max := numeric_util.Float64ToInt64(5.0)
	ranges := splitInt64Range(min, max, 4)
	enumerated := ranges.Enumerate()
	if len(enumerated) != 135 {
		t.Errorf("expected 135 terms, got %d", len(enumerated))
	}

}

func TestIncrementBytes(t *testing.T) {
	tests := []struct {
		in  []byte
		out []byte
	}{
		{
			in:  []byte{0},
			out: []byte{1},
		},
		{
			in:  []byte{0, 0},
			out: []byte{0, 1},
		},
		{
			in:  []byte{0, 255},
			out: []byte{1, 0},
		},
	}

	for _, test := range tests {
		actual := incrementBytes(test.in)
		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("expected %#v, got %#v", test.out, actual)
		}
	}
}

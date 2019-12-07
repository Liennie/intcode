package util

import (
	"reflect"
	"testing"
)

func TestPermInts(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  [][]int
	}{
		{
			name:  "nil",
			input: nil,
			want:  [][]int{{}},
		},
		{
			name:  "0",
			input: []int{},
			want:  [][]int{{}},
		},
		{
			name:  "1",
			input: []int{1},
			want:  [][]int{{1}},
		},
		{
			name:  "2",
			input: []int{1, 2},
			want: [][]int{
				{1, 2},
				{2, 1},
			},
		},
		{
			name:  "2 rev",
			input: []int{2, 1},
			want: [][]int{
				{1, 2},
				{2, 1},
			},
		},
		{
			name:  "2 same",
			input: []int{1, 1},
			want: [][]int{
				{1, 1},
			},
		},
		{
			name:  "2/3 same",
			input: []int{1, 1, 2},
			want: [][]int{
				{1, 1, 2},
				{1, 2, 1},
				{2, 1, 1},
			},
		},
		{
			name:  "3",
			input: []int{1, 2, 3},
			want: [][]int{
				{1, 2, 3},
				{1, 3, 2},
				{2, 1, 3},
				{2, 3, 1},
				{3, 1, 2},
				{3, 2, 1},
			},
		},
		{
			name:  "3 rev",
			input: []int{3, 2, 1},
			want: [][]int{
				{1, 2, 3},
				{1, 3, 2},
				{2, 1, 3},
				{2, 3, 1},
				{3, 1, 2},
				{3, 2, 1},
			},
		},
		{
			name:  "3 rand",
			input: []int{2, 1, 3},
			want: [][]int{
				{1, 2, 3},
				{1, 3, 2},
				{2, 1, 3},
				{2, 3, 1},
				{3, 1, 2},
				{3, 2, 1},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := make([][]int, 0)
			for p := range PermInts(test.input...) {
				res = append(res, p)
			}

			if !reflect.DeepEqual(test.want, res) {
				t.Fail()

				t.Logf("\nExpected: %v\nActual:   %v", test.want, res)
			}
		})
	}
}

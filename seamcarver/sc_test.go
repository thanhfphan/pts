package seamcarver

import (
	"testing"
)

func TestRetrieveSeamPath(t *testing.T) {
	tests := []struct {
		name   string
		energy [][]float64
		want   []int
	}{
		{
			name: "simple case",
			energy: [][]float64{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			want: []int{0, 0, 0},
		},
		{
			name: "another simple case",
			energy: [][]float64{
				{10, 11, 12},
				{10, 1, 10},
				{11, 10, 12},
			},
			want: []int{0, 1, 1},
		},
		{
			name: "complex case",
			energy: [][]float64{
				{1, 2, 3, 4},
				{4, 3, 2, 1},
				{1, 2, 3, 4},
				{4, 3, 2, 1},
			},
			want: []int{0, 1, 0, 1},
		},
		{
			name: "single column",
			energy: [][]float64{
				{1},
				{2},
				{3},
			},
			want: []int{0, 0, 0},
		},
		{
			name: "single row",
			energy: [][]float64{
				{1, 2, 3},
			},
			want: []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := retrieveSeamPath(tt.energy)
			if !equal(got, tt.want) {
				t.Errorf("retrieveSeamPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

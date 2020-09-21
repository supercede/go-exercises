package sum_test

import (
	"testing"

	"github.com/supercede/go-exercises/tests/sum"
)

func TestInts(t *testing.T) {
	tt := []struct {
		name    string
		numbers []int
		result  int
	}{
		{"sum 1, 2, 5, 6, 7", []int{1, 2, 5, 6, 7}, 21},
		{"sum 1 to 5", []int{1, 2, 3, 4, 5}, 15},
		{"sum nil", nil, 0},
		{"sum empty list", []int{}, 0},
		{"sum 1 and -4", []int{1, -4}, -3},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sum := sum.Ints(tc.numbers...)
			if sum != tc.result {
				t.Errorf("Test failed. Expected %s to equal %v, instead got %v", tc.name, tc.result, sum)
			}
		})
	}
}

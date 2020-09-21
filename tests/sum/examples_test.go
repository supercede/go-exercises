package sum_test

import (
	"fmt"

	"github.com/supercede/go-exercises/tests/sum"
)

func ExampleInts() {
	s := sum.Ints(1, 2, 3, 4, 5)
	fmt.Println("The sum is", s)
	// Output:
	// The sum is 15
}

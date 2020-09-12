package quiz_test

import (
	"fmt"

	"github.com/supercede/go-exercises/csv-quiz/quiz"
)

func ExampleParseFile() {
	questions := quiz.ParseFile("examples.csv")
	fmt.Println(questions)
	// Output:
	// [{5+5 10} {7/2 3.5} {5^3 125}]
}

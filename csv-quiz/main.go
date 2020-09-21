package main

import (
	"flag"
	"fmt"

	"github.com/supercede/go-exercises/csv-quiz/quiz"
)

func main() {
	helpString := flag.String("csv", "problems.csv", "Add a valid csv file with the format 'question, answer'")
	flag.Parse()

	questions := quiz.ParseFile(*helpString)
	result := quiz.PrintQuestions(questions)
	fmt.Println(result)
}

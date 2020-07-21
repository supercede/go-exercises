package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

type problem struct {
	question, answer string
}

func handleError(e error, message string) {
	if e != nil {
		fmt.Printf(message)
		os.Exit(1)
	}
}

func parseFile(fileName string) []problem {
	f, err := os.Open(fileName)
	handleError(err, fmt.Sprintf("Couldn't find file: %s\n", fileName))

	file := csv.NewReader(f)
	records, err := file.ReadAll()

	handleError(err, fmt.Sprintf("Failed to parse file: %s\n", fileName))

	var questions = make([]problem, len(records))
	for i, r := range records {
		questions[i] = problem{
			r[0],
			strings.TrimSpace(r[1]),
		}
	}
	return questions
}

func printQuestions(q []problem) string {
	correct := 0
	for i, d := range q {
		fmt.Printf("Question %d: %v\n", i+1, d.question)

		var answer string
		fmt.Scanf("%s\n", &answer)

		if answer == d.answer {
			correct += 1
		}
	}
	return fmt.Sprintf("You got %d questions correctly out of %d", correct, len(q))
}

func main() {
	helpString := flag.String("csv", "problems.csv", "Add a valid csv file with the format 'question, answer'")
	flag.Parse()

	questions := parseFile(*helpString)
	result := printQuestions(questions)
	fmt.Println(result)
}

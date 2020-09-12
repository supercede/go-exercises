package quiz

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type problem struct {
	question, answer string
}

// Handles runtime errors
func handleError(e error, message string) {
	if e != nil {
		fmt.Println(message)
		os.Exit(1)
	}
}

// Parses a csv file into a slice of Problem type
func ParseFile(fileName string) []problem {
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

// Prints each question and parses the given answer
func PrintQuestions(q []problem) string {
	correct := 0
	for i, d := range q {
		fmt.Printf("Question %d: %v\n", i+1, d.question)

		var answer string
		fmt.Scanf("%s\n", &answer)

		if answer == d.answer {
			correct++
		}
	}
	return fmt.Sprintf("You got %d questions correctly out of %d", correct, len(q))
}

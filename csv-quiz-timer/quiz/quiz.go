package quiz

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type Problem struct {
	Question, Answer string
}

// Handles runtime errors
func handleError(e error, message string) {
	if e != nil {
		fmt.Println(message)
		os.Exit(1)
	}
}

// Parses a csv file into a slice of Problem type
func ParseFile(fileName string) []Problem {
	f, err := os.Open(fileName)
	handleError(err, fmt.Sprintf("Couldn't find file: %s\n", fileName))

	file := csv.NewReader(f)
	records, err := file.ReadAll()

	handleError(err, fmt.Sprintf("Failed to parse file: %s\n", fileName))

	var questions = make([]Problem, len(records))
	for i, r := range records {
		questions[i] = Problem{
			r[0],
			strings.TrimSpace(r[1]),
		}
	}
	return questions
}

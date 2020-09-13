package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/supercede/go-exercises/csv-quiz/quiz"
)

func main() {
	helpString := flag.String("csv", "problems.csv", "Add a valid csv file with the format 'question, answer'")
	flag.Parse()

	timeLimit := flag.Int("limit", 5, "Time limit for the quiz")

	questions := quiz.ParseFile(*helpString)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, d := range questions {
		fmt.Printf("Question %d: %v\n", i+1, d.Question)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)

			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Println("Time Up!")
			fmt.Printf("You got %d questions correctly out of %d", correct, len(questions))
			return
		case answer := <-answerCh:
			if answer == d.Answer {
				correct++
			}
		}
	}
	fmt.Printf("You got %d questions correctly out of %d", correct, len(questions))
}

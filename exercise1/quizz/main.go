package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// Return -filename and -timer values.
func parseFlags() (filename string, timer int) {
	flagFilename := flag.String("filename", "./problems.csv", "Name of the quizz CSV file.")
	flagTimer := flag.Int("timer", 30, "Time limit in seconds.")
	flag.Parse()
	return *flagFilename, *flagTimer
}

func askQuestions(quizzList []Quizz, scanner *bufio.Scanner, finished chan struct{}, correct *int) {
	for i, quizz := range quizzList {
		fmt.Printf("Question #%d: %s ? ", i+1, quizz.question)
		scanner.Scan()
		if scanner.Err() != nil {
			fmt.Printf("Error getting input: %v", scanner.Err())
			return
		}
		answer := strings.TrimSpace(scanner.Text())
		if answer == quizz.answer {
			*correct++
			fmt.Println("Correct answer!")
		} else {
			fmt.Println("Wrong answer!")
		}
	}
	finished <- struct{}{}
}

func waitQuizz(finished chan struct{}, timeLimit int) {
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	select {
	case <-finished:
		break
	case <-timer.C:
		fmt.Println("\nTime expired!")
		break
	}
}

func runQuizz(filename string, timeLimit int) {
	reader, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Can't open %s: %v\n", filename, err)
		return
	}
	scanner := bufio.NewScanner(os.Stdin)
	quizzList, _ := GetQuizzList(reader)
	finished := make(chan struct{})
	correct := 0

	fmt.Printf("Time for Quizz: %d\n", timeLimit)
	go askQuestions(quizzList, scanner, finished, &correct)
	waitQuizz(finished, timeLimit)
	fmt.Printf("Correct answers: %d/%d\n", correct, len(quizzList))
}

func main() {
	filename, timeLimit := parseFlags()
	runQuizz(filename, timeLimit)
}

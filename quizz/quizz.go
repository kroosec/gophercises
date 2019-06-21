package main

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"
)

// Quizz entry.
type Quizz struct {
	question string
	answer   string
}

// ErrorQuestion Errorneous question.
var ErrorQuestion = errors.New("Erroneous Quizz Question")

// GetQuizzList Returns a slice of Quizz{question, answer}, from a CSV-formatted input.
func GetQuizzList(quizzReader io.Reader) ([]Quizz, error) {
	var quizzList []Quizz

	reader := csv.NewReader(quizzReader)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(record) != 2 {
			return nil, ErrorQuestion
		}
		question := strings.TrimSpace(record[0])
		answer := strings.TrimSpace(record[1])
		quizzList = append(quizzList, Quizz{question, answer})
	}
	return quizzList, nil
}

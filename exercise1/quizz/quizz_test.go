package main

import (
	"strings"
	"testing"
)

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if got != want {
		t.Fatalf("want %v error, got %v", got, want)
	}
}

func assertEqualInt(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("want %d got %d", want, got)
	}
}

func assertEqualBool(t *testing.T, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("want %v got %v", want, got)
	}
}

func assertEqualStr(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("want %s got %s", want, got)
	}
}

func TestGetQuizzList(t *testing.T) {
	t.Run("Correct parsing of quizz questions", func(t *testing.T) {
		quizzString := `5+5,10
7+3,10
1+1,2
8+3,11
1+2 ,3
8+6, 14
3+1,4
1+4,5
5+1,6
2+3,5
3+3,6
2+4,6
5+2,7`
		quizzList, err := GetQuizzList(strings.NewReader(quizzString))
		assertError(t, err, nil)
		want := []Quizz{
			{"5+5", "10"},
			{"7+3", "10"},
			{"1+1", "2"},
			{"8+3", "11"},
			{"1+2", "3"},
			{"8+6", "14"},
			{"3+1", "4"},
			{"1+4", "5"},
			{"5+1", "6"},
			{"2+3", "5"},
			{"3+3", "6"},
			{"2+4", "6"},
			{"5+2", "7"},
		}
		assertEqualInt(t, len(quizzList), len(want))
		for i, quizz := range want {
			assertEqualStr(t, quizzList[i].question, quizz.question)
			assertEqualStr(t, quizzList[i].answer, quizz.answer)
		}
	})

	t.Run("Erroneous entry", func(t *testing.T) {
		quizzString := "5+5,10,18"
		_, err := GetQuizzList(strings.NewReader(quizzString))
		assertError(t, err, ErrorQuestion)
	})
}

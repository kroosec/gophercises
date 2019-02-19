package cyoa

import (
	"reflect"
	"testing"
	"strings"
)

var testAdventure string = `{
  "intro": {
    "title": "title1",
    "story": [
      "story11"
    ],
    "options": [
      {
        "text": "text11",
        "arc": "foo2"
      },
      {
        "text": "text12",
        "arc": "foo3"
      }
    ]
  }
}
`

func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("expected no error, got '%v'", got)
	}
}

func assertParsed(t *testing.T, got, want Story) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want '%v' got '%v'", want, got)
	}
}

func TestParseAdventure(t *testing.T) {
	want := Story{"intro": Arc{"title1", []string{"story11"}, []Option{{"text11", "foo2"}, {"text12", "foo3"}}}}
	result, err := ParseAdventureJson(strings.NewReader(testAdventure))
	assertNoError(t, err)
	assertParsed(t, result, want)
}

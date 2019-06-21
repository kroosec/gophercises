package cyoa

import (
	"encoding/json"
	"errors"
	"io"
)

type Story map[string]Arc

type Option struct {
	Text string
	Arc  string
}

type Arc struct {
	Title      string
	Paragraphs []string `json:"story"`
	Options    []Option
}

func ParseAdventureJson(r io.Reader) (Story, error) {
	var parsed Story
	d := json.NewDecoder(r)
	if err := d.Decode(&parsed); err != nil {
		return nil, err
	}
	if _, ok := parsed["intro"]; ok == false {
		return nil, errors.New("Missing intro")
	}
	return parsed, nil
}

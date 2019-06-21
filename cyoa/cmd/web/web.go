package main

import (
	"flag"
	"fmt"
	"github.com/kroosec/gophercises/cyoa"
	"html/template"
	"log"
	"net/http"
	"os"
)

type storyHandler struct {
	story cyoa.Story
	tmpl  *template.Template
}

func (s *storyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var arcName string
	if r.URL.Path == "/" {
		arcName = "intro"
	} else {
		arcName = r.URL.Path[1:]
	}
	if arc, ok := s.story[arcName]; ok {
		err := s.tmpl.Execute(w, arc)
		if err != nil {
			fmt.Printf("%v", err)
			http.Error(w, "Error", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Arc not found.", http.StatusNotFound)
	}
}

func main() {
	port := flag.Int("port", 1234, "Server port to listen on.")
	fileName := flag.String("filename", "./gopher.json", "Filename of the adventure JSON.")
	flag.Parse()

	reader, err := os.Open(*fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	story, err := cyoa.ParseAdventureJson(reader)
	if err != nil {
		fmt.Printf("Couldn't parse '%s': %v", *fileName, err)
		return
	}

	tmpl := template.Must(template.ParseFiles("template.html"))
	http.Handle("/", &storyHandler{story, tmpl})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

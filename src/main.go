package main

import (
	"embed"
	_ "embed" // required for go:embed
	"github.com/johnerator/src/types"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const (
	baseHTML = "html/templates/base.layout.html"
	characterHTML = "html/templates/character.html"
	footerHTML = "html/templates/footer.partial.html"
	homeHTML = "html/templates/home.page.html"
	index = "/"
	newCharacterHTML = "html/templates/new-character.html"
)

var n types.Names
var s types.Skills
var o types.Obsessions

//go:embed html/templates/*
var templates embed.FS

//go:embed assets/names.json
var namesJSON []byte

//go:embed assets/obsessions.json
var obsJSON []byte

//go:embed assets/skills.json
var skillsJSON []byte

func characterHandler(w http.ResponseWriter, r *http.Request) {
	var skillCount int
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "error creating character", 500)
	}
	if r.Form["skills"] != nil {
		skillCount, _ = strconv.Atoi(r.Form["skills"][0])
	}
	c := types.CreateNewCharacter(n, s, o, skillCount)
	parsedTemplate, err := template.ParseFS(templates, characterHTML)
	if err != nil {
		log.Fatal("unable to parse ", err)
	}
	err = parsedTemplate.Execute(w, c)
	if err != nil {
		log.Println("error creating john ", err)
	}
}

func newCharacter(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, err := template.ParseFS(templates, newCharacterHTML)
	if err != nil {
		log.Fatal("unable to parse ", err)
	}
	err = parsedTemplate.Execute(w, nil)
	if err != nil {
		log.Println("error creating john ", err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != index {
		http.NotFound(w, r)
		return
	}
	files := []string{
		homeHTML,
		baseHTML,
		footerHTML,
	}
	templates, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = templates.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}

func main() {
	n = types.LoadNames(namesJSON)
	s = types.LoadSkills(skillsJSON)
	o = types.LoadObsessions(obsJSON)
	http.HandleFunc("/character", characterHandler)
	http.HandleFunc("/", home)
	http.HandleFunc("/new-character", newCharacter)
	// start server
	log.Println("Listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

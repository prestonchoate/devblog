package main

import (
	"html/template"
	"log"
	"net/http"
)

type Post struct {
	Id          int
	Title       string
	Image       string
	Content     string
	Description string
}

type Project struct {
	ID          int
	Title       string
	Image       string
	Description string
	Link        string
}

var templates *template.Template
var templateDirs = []string{
	"templates/partials/*.html",
	"templates/*.html",
}

var posts = []Post{
	{Id: 1, Title: "Post 1", Image: "public/placeholder.svg", Content: "Content 1", Description: "Description 1"},
	{Id: 2, Title: "Post 2", Image: "public/placeholder.svg", Content: "Content 2", Description: "Description 2"},
	{Id: 3, Title: "Post 3", Image: "public/placeholder.svg", Content: "Content 3", Description: "Description 3"},
	{Id: 4, Title: "Post 4", Image: "public/placeholder.svg", Content: "Content 4", Description: "Description 4"},
	{Id: 5, Title: "Post 5", Image: "public/placeholder.svg", Content: "Content 5", Description: "Description 5"},
	{Id: 6, Title: "Post 6", Image: "public/placeholder.svg", Content: "Content 6", Description: "Description 6"},
	{Id: 7, Title: "Post 7", Image: "public/placeholder.svg", Content: "Content 7", Description: "Description 7"},
}

var projects = []Project{
	{ID: 1, Title: "Project 1", Image: "public/placeholder.svg", Description: "Description 1", Link: "https://github.com/prestonchoate"},
	{ID: 2, Title: "Project 2", Image: "public/placeholder.svg", Description: "Description 2", Link: "https://github.com/prestonchoate"},
	{ID: 3, Title: "Project 3", Image: "public/placeholder.svg", Description: "Description 3", Link: "https://github.com/prestonchoate"},
	{ID: 4, Title: "Project 4", Image: "public/placeholder.svg", Description: "Description 4", Link: "https://github.com/prestonchoate"},
	{ID: 5, Title: "Project 5", Image: "public/placeholder.svg", Description: "Description 5", Link: "https://github.com/prestonchoate"},
}

func loadTemplates() (templates *template.Template) {
	for _, dir := range templateDirs {
		var err error
		templates, err = templates.ParseGlob(dir)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println(templates.DefinedTemplates())
	return templates
}

func init() {
	templates = loadTemplates()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("homeHandler")
	data := map[string]interface{}{
		"title":    "Home",
		"posts":    posts,
		"projects": projects,
	}
	err := templates.ExecuteTemplate(w, "home_page", data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {

	fs := http.FileServer(http.Dir("./public/"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", homeHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

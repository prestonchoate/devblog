package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var templDirs = []string{
	"./templates/",
}

var templates *template.Template

func getTemplates() (templates *template.Template, err error) {
	var allFiles []string
	for _, dir := range templDirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Fatalf("Unable to read templates! Error: %v\n", err.Error())
		}
		for _, file := range files {
			filename := file.Name()
			if strings.HasSuffix(filename, ".tmpl") {
				filepath := filepath.Join(dir, filename)
				allFiles = append(allFiles, filepath)
			}
		}
	}
	templates, err = template.New("").ParseFiles(allFiles...)
	if err != nil {
		log.Fatalf("Unable to parse templates! Error: %v\n", err.Error())
	}
	return
}

func init() {
	templates, _ = getTemplates()
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

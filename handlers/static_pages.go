package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/prestonchoate/devblog/config"
	"github.com/prestonchoate/devblog/data"
)

// FileSystem custom file system handler
type FileSystem struct {
	Fs http.FileSystem
}

// Open opens file
func (fs FileSystem) Open(path string) (http.File, error) {
	f, err := fs.Fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.Fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	config := config.GetInstance()
	data := map[string]interface{}{
		"title": "About",
		"links": config.GetLinks(),
	}
	err := config.GetTemplates().ExecuteTemplate(w, "about_page", data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	config := config.GetInstance()
	postRepo, err := data.GetPostRepositoryInstance()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	projectRepo, err := data.GetProjectRepositoryInstance()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	// only send the first 3 posts and projects to the home page
	data := map[string]interface{}{
		"title":      "Home",
		"posts":      postRepo.GetPosts(3),
		"projects":   projectRepo.GetProjects(3),
		"links":      config.GetLinks(),
		"showHeader": true,
	}
	err = config.GetTemplates().ExecuteTemplate(w, "home_page", data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

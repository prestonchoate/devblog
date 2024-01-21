package handlers

import (
	"log"
	"net/http"

	"github.com/prestonchoate/devblog/config"
	"github.com/prestonchoate/devblog/data"
)

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
	postRepo := data.GetPostRepositoryInstance()
	projectRepo := data.GetProjectRepositoryInstance()
	// only send the first 3 posts and projects to the home page
	data := map[string]interface{}{
		"title":      "Home",
		"posts":      postRepo.GetPosts(3),
		"projects":   projectRepo.GetProjects(3),
		"links":      config.GetLinks(),
		"showHeader": true,
	}
	err := config.GetTemplates().ExecuteTemplate(w, "home_page", data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

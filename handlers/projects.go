package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/prestonchoate/devblog/config"
	"github.com/prestonchoate/devblog/data"
)

func HandleProjectListing(w http.ResponseWriter, r *http.Request) {
	c := config.GetInstance()
	projectRepo := data.GetProjectRepositoryInstance()
	data := map[string]interface{}{
		"title":    "Projects",
		"projects": projectRepo.GetProjects(0),
		"links":    c.GetLinks(),
	}
	err := c.GetTemplates().ExecuteTemplate(w, "project_listing", data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func HandleProjectView(w http.ResponseWriter, r *http.Request) {
	project := r.Context().Value("project").(*data.Project)
	if project == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	c := config.GetInstance()
	data := map[string]interface{}{
		"title":   project.Title,
		"project": project,
		"links":   c.GetLinks(),
	}
	err := c.GetTemplates().ExecuteTemplate(w, "project", data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ProjectCtx(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectIdParam := chi.URLParam(r, "projectId")
		projectId, err := strconv.Atoi(projectIdParam)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		project := data.GetProjectRepositoryInstance().GetProject(projectId)
		if project == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "project", project)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

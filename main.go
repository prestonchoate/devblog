package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prestonchoate/devblog/config"
	"github.com/prestonchoate/devblog/data"
	"github.com/prestonchoate/devblog/handlers"
)

func main() {
	_ = config.GetInstance()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	fileSystem := handlers.FileSystem{Fs: http.Dir("./public/")}
	fs := http.FileServer(fileSystem)
	r.Mount("/public/", http.StripPrefix(strings.TrimRight("/public", "/"), fs))

	r.Get("/", handlers.HomeHandler)
	r.Get("/about", handlers.AboutHandler)

	r.Route("/blog", func(r chi.Router) {
		r.Get("/", handlers.HandlePostListing)
		r.Route("/{slug}", func(r chi.Router) {
			r.Use(handlers.PostCtx)
			r.Get("/", handlers.HandlePostView)
		})
	})

	r.Route("/projects", func(r chi.Router) {
		r.Get("/", handlers.HandleProjectListing)
		r.Route("/{projectId}", func(r chi.Router) {
			r.Use(handlers.ProjectCtx)
			r.Get("/", handlers.HandleProjectView)
		})
	})

	r.Route("/dashboard", func(r chi.Router) {
		r.Get("/login", handlers.HandleAdminLoginPage)
		r.Route("/admin", func(r chi.Router) {
			r.Use(handlers.AdminCtx)
			r.Get("/", handlers.HandleAdminDashboard)
		})
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/login", handlers.HandleApiLogin)
	})

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		posts := data.GetCmsPosts()
		data, err := json.Marshal(posts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)

	})

	log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", r))
}

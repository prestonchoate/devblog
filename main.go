package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prestonchoate/devblog/config"
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
		r.Route("/{postId}", func(r chi.Router) {
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

	log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", r))
}
